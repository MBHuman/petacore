package wire

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"petacore/internal/logger"
	"petacore/internal/runtime/executor"
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/runtime/rsql/visitor"
	"petacore/internal/runtime/system"
	"petacore/internal/storage"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

// ─────────────────────────────────────────────
// Session
// ─────────────────────────────────────────────

type Session struct {
	preparedStatements map[string]*PreparedStatement
	portals            map[string]*PreparedStatement
	params             map[string]string
}

func NewSession() *Session {
	return &Session{
		preparedStatements: make(map[string]*PreparedStatement),
		portals:            make(map[string]*PreparedStatement),
		params:             make(map[string]string),
	}
}

// ─────────────────────────────────────────────
// PreparedStatement
// ─────────────────────────────────────────────

type FieldDescription struct {
	Name        string
	DataTypeOID uint32
}

type PreparedStatement struct {
	Query             string
	Stmt              statements.SQLStatement
	Params            []interface{}
	Columns           []FieldDescription
	ParamOIDs         []uint32
	ResultFormatCodes []int16
}

// ─────────────────────────────────────────────
// WireServer
// ─────────────────────────────────────────────

type WireServer struct {
	storage   *storage.DistributedStorageVClock
	listener  net.Listener
	port      string
	allocPool *pmem.ArenaPool
}

func NewWireServer(storage *storage.DistributedStorageVClock, port string) *WireServer {
	return &WireServer{
		storage:   storage,
		port:      port,
		allocPool: pmem.NewArenaPool(26 * 1024 * 1024),
	}
}

func (ws *WireServer) Start() error {
	listener, err := net.Listen("tcp", ":"+ws.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", ws.port, err)
	}
	ws.listener = listener
	go ws.acceptConnections()
	return nil
}

func (ws *WireServer) Stop() error {
	if ws.listener != nil {
		return ws.listener.Close()
	}
	return nil
}

func (ws *WireServer) acceptConnections() {
	for {
		conn, err := ws.listener.Accept()
		if err != nil {
			return
		}
		logger.Info("Accepted connection")
		go ws.handleConnection(conn)
	}
}

// ─────────────────────────────────────────────
// Connection handler
// ─────────────────────────────────────────────

func (ws *WireServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Minute))
	if tc, ok := conn.(*net.TCPConn); ok {
		tc.SetWriteBuffer(4 * 1024 * 1024) // 4MB вместо дефолтных ~128KB
		tc.SetReadBuffer(1 * 1024 * 1024)
		tc.SetNoDelay(false)
	}

	// Single large write buffer — all protocol writes go here.
	// We Flush() only at explicit synchronisation points, so the OS sees
	// one writev() per logical message group instead of one per field.
	w := bufio.NewWriterSize(conn, 256*1024)
	r := bufio.NewReaderSize(conn, 65536)

	session := NewSession()

	// ── Startup / SSL negotiation ─────────────────────────────────────
	startupLen, err := readInt32(r)
	if err != nil {
		logger.Errorf("read startup length: %v", err)
		return
	}
	if startupLen < 4 {
		logger.Errorf("startup length too small: %d", startupLen)
		return
	}
	startupBody := make([]byte, startupLen-4)
	if _, err := io.ReadFull(r, startupBody); err != nil {
		logger.Errorf("read startup body: %v", err)
		return
	}

	protoVersion := binary.BigEndian.Uint32(startupBody[:4])

	if protoVersion == 80877103 { // SSLRequest
		w.WriteByte('N') // no SSL
		w.Flush()
		// Re-read real startup
		startupLen, err = readInt32(r)
		if err != nil {
			logger.Errorf("read startup after SSL: %v", err)
			return
		}
		startupBody = make([]byte, startupLen-4)
		if _, err := io.ReadFull(r, startupBody); err != nil {
			logger.Errorf("read startup body after SSL: %v", err)
			return
		}
	}

	ws.sendStartup(w)
	w.Flush()

	// ── Main message loop ─────────────────────────────────────────────
	for {
		msgType, err := r.ReadByte()
		if err != nil {
			if err != io.EOF {
				logger.Errorf("read message type: %v", err)
			}
			return
		}
		msgLen, err := readInt32(r)
		if err != nil {
			logger.Errorf("read message length: %v", err)
			return
		}
		bodyLen := int(msgLen) - 4
		var body []byte
		if bodyLen > 0 {
			body = make([]byte, bodyLen)
			if _, err := io.ReadFull(r, body); err != nil {
				logger.Errorf("read message body: %v", err)
				return
			}
		}

		switch msgType {
		case 'Q': // Simple query
			query := cstring(body)
			logger.Debugf("Query: %s", query)
			ws.handleQuery(w, query, session)
			w.Flush()

		case 'P': // Parse
			ws.handleParse(w, body, session)
			w.Flush()

		case 'B': // Bind
			ws.handleBind(w, body, session)
			w.Flush()

		case 'D': // Describe
			ws.handleDescribe(w, body, session)
			w.Flush()

		case 'E': // Execute
			ws.handleExecute(w, body, session)
			w.Flush()

		case 'S': // Sync
			writeReadyForQuery(w)
			w.Flush()

		case 'H': // Flush
			w.Flush()

		case 'X': // Terminate
			return

		default:
			logger.Errorf("unsupported message type: %c", msgType)
			writeError(w, "0A000", "unsupported message type")
			writeReadyForQuery(w)
			w.Flush()
		}
	}
}

// ─────────────────────────────────────────────
// Startup
// ─────────────────────────────────────────────

func (ws *WireServer) sendStartup(w *bufio.Writer) {
	// 1. AuthenticationOk — ВСЕГДА первым
	writeMsg(w, 'R', func(b *msgBuf) {
		b.int32(0)
	})

	// 2. ParameterStatus — после аутентификации
	writeParameterStatus(w, "server_version", "13.0.0")
	writeParameterStatus(w, "client_encoding", "UTF8")
	writeParameterStatus(w, "server_encoding", "UTF8")
	writeParameterStatus(w, "standard_conforming_strings", "on")
	writeParameterStatus(w, "TimeZone", "UTC")
	writeParameterStatus(w, "integer_datetimes", "on")
	writeParameterStatus(w, "application_name", "")

	// 3. BackendKeyData
	writeMsg(w, 'K', func(b *msgBuf) {
		b.int32(1)
		b.int32(1)
	})

	// 4. ReadyForQuery — последним
	writeReadyForQuery(w)
}

// ─────────────────────────────────────────────
// Parse
// ─────────────────────────────────────────────

func (ws *WireServer) handleParse(w *bufio.Writer, body []byte, session *Session) {
	name, rest := readCString(body)
	query, rest := readCString(rest)
	_ = rest // paramOIDs from client — we infer our own

	logger.Debugf("Parse name=%q query=%q", name, query)

	allocator := ws.allocPool.Get()
	defer ws.allocPool.Put(allocator)

	var stmt statements.SQLStatement
	if strings.TrimSpace(query) == "" {
		stmt = &statements.EmptyStatement{}
	} else {
		var err error
		stmt, err = visitor.ParseSQL(allocator, query)
		if err != nil {
			writeError(w, "42601", err.Error())
			writeReadyForQuery(w)
			return
		}
	}

	paramCount := countParams(query)
	paramOIDs := make([]uint32, paramCount)
	for i := range paramOIDs {
		paramOIDs[i] = 23 // INT4
	}

	cols := ws.getFieldDescriptions(stmt)
	session.preparedStatements[name] = &PreparedStatement{
		Query:     query,
		Stmt:      stmt,
		ParamOIDs: paramOIDs,
		Columns:   cols,
	}

	// ParseComplete: '1' + int32(4)
	writeMsg(w, '1', func(b *msgBuf) {})
}

// ─────────────────────────────────────────────
// Bind
// ─────────────────────────────────────────────

func (ws *WireServer) handleBind(w *bufio.Writer, body []byte, session *Session) {
	portal, rest := readCString(body)
	stmtName, rest := readCString(rest)

	logger.Infof("Bind stmt=%q portal=%q", stmtName, portal)

	ps, ok := session.preparedStatements[stmtName]
	if !ok {
		writeError(w, "26000", "prepared statement does not exist")
		writeReadyForQuery(w)
		return
	}

	// Number of parameter format codes
	numFmtCodes := int(readInt16(rest))
	rest = rest[2:]
	paramFmtCodes := make([]int16, numFmtCodes)
	for i := 0; i < numFmtCodes; i++ {
		paramFmtCodes[i] = readInt16(rest)
		rest = rest[2:]
	}

	// Number of parameters
	numParams := int(readInt16(rest))
	rest = rest[2:]

	if numParams != len(ps.ParamOIDs) {
		writeError(w, "08P01", fmt.Sprintf("expected %d parameters, got %d", len(ps.ParamOIDs), numParams))
		writeReadyForQuery(w)
		return
	}

	params := make([]interface{}, numParams)
	for i := 0; i < numParams; i++ {
		paramLen := int(int32(binary.BigEndian.Uint32(rest)))
		rest = rest[4:]
		if paramLen == -1 {
			params[i] = nil
			continue
		}
		paramBytes := rest[:paramLen]
		rest = rest[paramLen:]

		paramFmt := int16(0)
		if len(paramFmtCodes) == 1 {
			paramFmt = paramFmtCodes[0]
		} else if i < len(paramFmtCodes) {
			paramFmt = paramFmtCodes[i]
		}

		if paramFmt == 0 {
			// Text format
			params[i] = string(paramBytes)
		} else {
			// Binary format - decode based on OID and length
			paramOID := ps.ParamOIDs[i]
			switch paramOID {
			case uint32(ptypes.PTypeBool): // BOOL
				if len(paramBytes) == 1 {
					params[i] = paramBytes[0] != 0
				} else {
					params[i] = string(paramBytes)
				}
			case uint32(ptypes.PTypeInt8), uint32(ptypes.PTypeInt2), uint32(ptypes.PTypeInt4): // INT8, INT2, INT4
				switch len(paramBytes) {
				case 2:
					v := int16(binary.BigEndian.Uint16(paramBytes))
					params[i] = strconv.FormatInt(int64(v), 10)
				case 4:
					v := int32(binary.BigEndian.Uint32(paramBytes))
					params[i] = strconv.FormatInt(int64(v), 10)
				case 8:
					v := int64(binary.BigEndian.Uint64(paramBytes))
					params[i] = strconv.FormatInt(v, 10)
				default:
					params[i] = string(paramBytes)
				}
			case uint32(ptypes.PTypeFloat4), uint32(ptypes.PTypeFloat8): // FLOAT4, FLOAT8
				switch len(paramBytes) {
				case 4:
					bits := binary.BigEndian.Uint32(paramBytes)
					params[i] = fmt.Sprintf("%g", math.Float32frombits(bits))
				case 8:
					bits := binary.BigEndian.Uint64(paramBytes)
					params[i] = fmt.Sprintf("%g", math.Float64frombits(bits))
				default:
					params[i] = string(paramBytes)
				}
			default:
				// For other types, keep as string
				params[i] = string(paramBytes)
			}
		}
	}

	// Result format codes
	numResFmt := int(readInt16(rest))
	rest = rest[2:]
	resFmtCodes := make([]int16, numResFmt)
	for i := 0; i < numResFmt; i++ {
		resFmtCodes[i] = readInt16(rest)
		rest = rest[2:]
	}

	session.portals[portal] = &PreparedStatement{
		Query:             ps.Query,
		Stmt:              ps.Stmt,
		Params:            params,
		Columns:           ps.Columns,
		ParamOIDs:         ps.ParamOIDs,
		ResultFormatCodes: resFmtCodes,
	}

	// BindComplete: '2' + int32(4)
	writeMsg(w, '2', func(b *msgBuf) {})
}

// ─────────────────────────────────────────────
// Describe
// ─────────────────────────────────────────────

func (ws *WireServer) handleDescribe(w *bufio.Writer, body []byte, session *Session) {
	objType := body[0]
	name := cstring(body[1:])

	switch objType {
	case 'S':
		ps, ok := session.preparedStatements[name]
		if !ok {
			writeError(w, "26000", "prepared statement does not exist")
			writeReadyForQuery(w)
			return
		}
		writeParameterDescription(w, ps.ParamOIDs)
		// Don't send RowDescription here — defer to Execute
		logger.Debugf("Describe stmt=%q cols=%d", name, len(ps.Columns))

	case 'P':
		_, ok := session.portals[name]
		if !ok {
			writeError(w, "26000", "portal does not exist")
			writeReadyForQuery(w)
			return
		}
		logger.Debugf("Describe portal=%q", name)
	}
}

// ─────────────────────────────────────────────
// Execute
// ─────────────────────────────────────────────

func (ws *WireServer) handleExecute(w *bufio.Writer, body []byte, session *Session) {
	portal := cstring(body)
	logger.Debugf("Execute portal=%q", portal)

	allocator := ws.allocPool.Get()

	ps, ok := session.portals[portal]
	if !ok {
		ws.allocPool.Put(allocator)
		writeError(w, "26000", "portal does not exist")
		writeReadyForQuery(w)
		return
	}

	stmt := ps.Stmt
	if len(ps.Params) > 0 {
		substituted := substituteParams(ps.Query, ps.Params)
		logger.Debugf("Substituted query: %s", substituted)
		var err error
		stmt, err = visitor.ParseSQL(allocator, substituted)
		if err != nil {
			ws.allocPool.Put(allocator)
			writeError(w, "42601", err.Error())
			writeReadyForQuery(w)
			return
		}
	}

	result, err := executor.ExecuteStatement(allocator, stmt, ws.storage, session.params)
	if err != nil {
		ws.allocPool.Put(allocator)
		writeError(w, "XX000", err.Error())
		writeReadyForQuery(w)
		return
	}

	switch ps.Stmt.(type) {
	case *statements.EmptyStatement:
		writeEmptyQueryResponse(w)
	case *statements.CreateTableStatement:
		writeCommandComplete(w, "CREATE TABLE")
	case *statements.InsertStatement:
		writeCommandComplete(w, "INSERT")
	case *statements.DropTableStatement:
		writeCommandComplete(w, "DROP TABLE")
	case *statements.TruncateTableStatement:
		writeCommandComplete(w, "TRUNCATE TABLE")
	case *statements.SetStatement:
		writeCommandComplete(w, "SET")
	case *statements.SelectStatement:
		ws.sendSelectResult(w, result, ps.Columns, ps.ResultFormatCodes)
	}

	ws.allocPool.Put(allocator)
}

// ─────────────────────────────────────────────
// Simple Query
// ─────────────────────────────────────────────

func (ws *WireServer) handleQuery(w *bufio.Writer, query string, session *Session) {
	q := strings.TrimSpace(query)
	if q == "" || q == ";" || strings.EqualFold(q, "select 1") || strings.EqualFold(q, "keep alive") {
		writeCommandComplete(w, "SELECT 1")
		writeReadyForQuery(w)
		return
	}

	allocator := ws.allocPool.Get()

	stmt, err := visitor.ParseSQL(allocator, query)
	if err != nil {
		ws.allocPool.Put(allocator)
		writeError(w, "42601", err.Error())
		writeReadyForQuery(w)
		return
	}

	result, err := executor.ExecuteStatement(allocator, stmt, ws.storage, session.params)
	if err != nil {
		ws.allocPool.Put(allocator)
		writeError(w, "XX000", err.Error())
		writeReadyForQuery(w)
		return
	}

	switch stmt.(type) {
	case *statements.CreateTableStatement:
		writeCommandComplete(w, "CREATE TABLE")
		ws.allocPool.Put(allocator)
	case *statements.InsertStatement:
		writeCommandComplete(w, "INSERT")
		ws.allocPool.Put(allocator)
	case *statements.DropTableStatement:
		writeCommandComplete(w, "DROP TABLE")
		ws.allocPool.Put(allocator)
	case *statements.TruncateTableStatement:
		writeCommandComplete(w, "TRUNCATE TABLE")
		ws.allocPool.Put(allocator)
	case *statements.SetStatement:
		writeCommandComplete(w, "SET")
		ws.allocPool.Put(allocator)
	case *statements.SelectStatement:
		cols := ws.getFieldDescriptions(stmt)
		ws.sendSelectResult(w, result, cols, nil)
		ws.allocPool.Put(allocator)
	}

	writeReadyForQuery(w)
}

// ─────────────────────────────────────────────
// SELECT result serialisation
// ─────────────────────────────────────────────

const flushThreshold = 512 * 1024 // 512KB

func (ws *WireServer) sendSelectResult(
	w *bufio.Writer,
	result *table.ExecuteResult,
	portalCols []FieldDescription,
	formatCodes []int16,
) {
	// Determine column metadata
	var colNames []string
	var colOIDs []ptypes.OID

	if result.Schema != nil && len(result.Schema.Fields) > 0 {
		for _, f := range result.Schema.Fields {
			colNames = append(colNames, f.Name)
			colOIDs = append(colOIDs, f.OID)
		}
	} else if len(portalCols) > 0 {
		for _, c := range portalCols {
			colNames = append(colNames, c.Name)
			colOIDs = append(colOIDs, ptypes.OID(c.DataTypeOID))
		}
	}

	nCols := len(colNames)

	// RowDescription
	writeMsg(w, 'T', func(b *msgBuf) {
		b.int16(int16(nCols))
		for i := 0; i < nCols; i++ {
			b.cstring(colNames[i])
			b.int32(0)                   // tableOID
			b.int16(0)                   // attNum
			b.uint32(uint32(colOIDs[i])) // dataTypeOID
			b.int16(-1)                  // dataTypeSize
			b.int32(-1)                  // typeModifier
			b.int16(0)                   // format (text)
		}
	})

	bp := rowsBufPool.Get().(*[]byte)
	buf := (*bp)[:0]

	for _, row := range result.Rows {
		buf = appendDataRow(buf, result, row, formatCodes)

		// Флашим кусками — не ждём пока накопится гигабайт
		if len(buf) >= flushThreshold {
			w.Write(buf)
			w.Flush()     // явный flush каждые 512KB
			buf = buf[:0] // сбрасываем без реаллокации
		}
	}

	if len(buf) > 0 {
		w.Write(buf)
	}

	*bp = buf
	rowsBufPool.Put(bp)

	writeCommandComplete(w, fmt.Sprintf("SELECT %d", len(result.Rows)))
}

var rowsBufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 256*1024) // 256KB начальный размер
		return &b
	},
}

func appendDataRow(dst []byte, result *table.ExecuteResult, row *ptypes.Row, formatCodes []int16) []byte {
	nCols := len(result.Schema.Fields)

	// Запоминаем позицию начала сообщения
	start := len(dst)

	// Заголовок: 'D' + placeholder length (заполним потом) + nCols
	dst = append(dst, 'D', 0, 0, 0, 0) // тип + длина
	dst = binary.BigEndian.AppendUint16(dst, uint16(nCols))

	for i := 0; i < nCols; i++ {
		var value interface{}
		if i < len(result.Schema.Fields) {
			buf, oid, err := result.Schema.GetField(row, i)
			if err == nil && buf != nil {
				if bt, err := serializers.DeserializeGeneric(buf, oid); err == nil && bt != nil {
					value = bt.IntoGo()
				}
			}
		}

		if value == nil {
			dst = binary.BigEndian.AppendUint32(dst, uint32(math.MaxUint32)) // -1 = NULL
			continue
		}

		fmtCode := int16(0)
		if len(formatCodes) == 1 {
			fmtCode = formatCodes[0]
		} else if i < len(formatCodes) {
			fmtCode = formatCodes[i]
		}

		colOID := result.Schema.Fields[i].OID
		encoded := encodeValue(value, colOID, fmtCode)

		dst = binary.BigEndian.AppendUint32(dst, uint32(len(encoded)))
		dst = append(dst, encoded...)
	}

	// Записываем итоговую длину сообщения (без байта типа)
	msgLen := len(dst) - start - 1 // -1 за тип 'D'
	binary.BigEndian.PutUint32(dst[start+1:], uint32(msgLen))

	return dst
}

// encodeValue serialises a Go value into PostgreSQL wire bytes (text or binary).
func encodeValue(value interface{}, colOID ptypes.OID, format int16) []byte {
	if format == 1 {
		return encodeBinary(value, colOID)
	}
	return encodeText(value, colOID)
}

func encodeText(value interface{}, colOID ptypes.OID) []byte {
	switch colOID {
	case ptypes.PTypeTimestamp, ptypes.PTypeTimestampz:
		switch ts := value.(type) {
		case int64:
			t := time.UnixMicro(ts)
			if colOID == ptypes.PTypeTimestampz {
				return []byte(t.Format("2006-01-02 15:04:05.999999-07"))
			}
			return []byte(t.UTC().Format("2006-01-02 15:04:05.999999"))
		case *time.Time:
			if ts != nil {
				if colOID == ptypes.PTypeTimestampz {
					return []byte(ts.Format("2006-01-02 15:04:05.999999-07"))
				}
				return []byte(ts.UTC().Format("2006-01-02 15:04:05.999999"))
			}
		}
	}

	switch v := value.(type) {
	case bool:
		if v {
			return []byte("t")
		}
		return []byte("f")
	case int:
		return strconv.AppendInt(nil, int64(v), 10)
	case int32:
		return strconv.AppendInt(nil, int64(v), 10)
	case int64:
		return strconv.AppendInt(nil, v, 10)
	case float32:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 32)
	case float64:
		return strconv.AppendFloat(nil, v, 'f', -1, 64)
	case string:
		return []byte(v)
	case []string:
		return []byte(formatArrayAsText(v))
	case []int32:
		return []byte(formatArrayAsText(v))
	case []int64:
		return []byte(formatArrayAsText(v))
	case []float64:
		return []byte(formatArrayAsText(v))
	case []bool:
		return []byte(formatArrayAsText(v))
	default:
		return []byte(fmt.Sprintf("%v", v))
	}
}

func encodeBinary(value interface{}, colOID ptypes.OID) []byte {
	switch colOID {
	case ptypes.PTypeBool:
		if b, ok := value.(bool); ok {
			if b {
				return []byte{1}
			}
			return []byte{0}
		}
	case ptypes.PTypeInt2, ptypes.PTypeInt4:
		v := toInt64(value)
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(int32(v)))
		return buf
	case ptypes.PTypeInt8:
		v := toInt64(value)
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v))
		return buf
	case ptypes.PTypeFloat8:
		var f float64
		switch vv := value.(type) {
		case float32:
			f = float64(vv)
		case float64:
			f = vv
		case string:
			f, _ = strconv.ParseFloat(vv, 64)
		}
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, math.Float64bits(f))
		return buf
	case ptypes.PTypeTimestamp, ptypes.PTypeTimestampz:
		var v64 int64
		switch vv := value.(type) {
		case int64:
			v64 = vv
		case int:
			v64 = int64(vv)
		case string:
			v64, _ = strconv.ParseInt(vv, 10, 64)
		}
		t := time.UnixMicro(v64)
		return []byte(fmt.Sprintf("%v", t))
	}
	// Fallback: text
	return encodeText(value, colOID)
}

func toInt64(v interface{}) int64 {
	switch vv := v.(type) {
	case int:
		return int64(vv)
	case int32:
		return int64(vv)
	case int64:
		return vv
	case string:
		n, _ := strconv.ParseInt(vv, 10, 64)
		return n
	}
	return 0
}

// ─────────────────────────────────────────────
// Field descriptions (for Describe/Parse)
// ─────────────────────────────────────────────

func (ws *WireServer) getFieldDescriptions(stmt statements.SQLStatement) []FieldDescription {
	sel, ok := stmt.(*statements.SelectStatement)
	if !ok || sel.Primary == nil {
		return nil
	}

	primary := sel.Primary
	var fields []FieldDescription

	colResolver := func(colName, alias, fnName string, isSelectAll bool) FieldDescription {
		name := alias
		if name == "" {
			if fnName != "" {
				name = fnName
			} else if colName != "" {
				name = colName
			} else {
				name = "?column?"
			}
		}
		oid := uint32(25) // TEXT
		if fnName != "" {
			if fn, ok := functions.GetRegisteredFunction(fnName); ok {
				oid = uint32(fn.GetFunction().ProRetType)
			}
		}
		return FieldDescription{Name: name, DataTypeOID: oid}
	}

	if primary.From == nil {
		for _, col := range primary.Columns {
			fnName := ""
			if col.Function != nil {
				fnName = col.Function.Name
			}
			fields = append(fields, colResolver(col.ColumnName, col.Alias, fnName, col.IsSelectAll))
		}
		return fields
	}

	tableName := primary.From.TableName
	if system.IsSystemTable(tableName) {
		for _, col := range primary.Columns {
			oid := uint32(25)
			name := col.ColumnName
			if col.ColumnName == "ssl" {
				oid = 16
			}
			fields = append(fields, FieldDescription{Name: name, DataTypeOID: oid})
		}
		return fields
	}

	// Normal table — schema not known statically
	return nil
}

// ─────────────────────────────────────────────
// Low-level protocol writers (zero extra alloc)
// ─────────────────────────────────────────────

// msgBuf accumulates a message body and writes it atomically.
type msgBuf struct {
	w    *bufio.Writer
	typ  byte
	data []byte
}

var writeBufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 8192)
		return &b
	},
}

// Собирает тип + длина + тело в один []byte → один w.Write()
func writeMsg(w *bufio.Writer, typ byte, fn func(b *msgBuf)) {
	bp := writeBufPool.Get().(*[]byte)
	buf := (*bp)[:0]

	// Резервируем место под заголовок (5 байт: тип + int32)
	buf = append(buf, typ, 0, 0, 0, 0)

	mb := &msgBuf{data: buf[5:]} // тело пишем после заголовка
	fn(mb)

	// Объединяем заголовок + тело
	full := append(buf[:5], mb.data...)
	// Записываем длину (len тела + 4 для самого int32)
	binary.BigEndian.PutUint32(full[1:5], uint32(len(mb.data)+4))

	w.Write(full)

	*bp = full[:0]
	writeBufPool.Put(bp)
}

func (b *msgBuf) raw(p []byte)     { b.data = append(b.data, p...) }
func (b *msgBuf) cstring(s string) { b.data = append(append(b.data, s...), 0) }
func (b *msgBuf) int16(v int16)    { b.data = binary.BigEndian.AppendUint16(b.data, uint16(v)) }
func (b *msgBuf) int32(v int32)    { b.data = binary.BigEndian.AppendUint32(b.data, uint32(v)) }
func (b *msgBuf) uint32(v uint32)  { b.data = binary.BigEndian.AppendUint32(b.data, v) }

func writeReadyForQuery(w *bufio.Writer) {
	// 'Z' + int32(5) + 'I'
	w.Write([]byte{'Z', 0, 0, 0, 5, 'I'})
}

func writeCommandComplete(w *bufio.Writer, tag string) {
	writeMsg(w, 'C', func(b *msgBuf) {
		b.cstring(tag)
	})
}

func writeEmptyQueryResponse(w *bufio.Writer) {
	w.Write([]byte{'I', 0, 0, 0, 4})
}

func writeError(w *bufio.Writer, code, msg string) {
	writeMsg(w, 'E', func(b *msgBuf) {
		b.data = append(b.data, 'S')
		b.cstring("ERROR")
		b.data = append(b.data, 'C')
		b.cstring(code)
		b.data = append(b.data, 'M')
		b.cstring(msg)
		b.data = append(b.data, 0) // terminator
	})
}

func writeParameterStatus(w *bufio.Writer, name, value string) {
	writeMsg(w, 'S', func(b *msgBuf) {
		b.cstring(name)
		b.cstring(value)
	})
}

func writeParameterDescription(w *bufio.Writer, oids []uint32) {
	writeMsg(w, 't', func(b *msgBuf) {
		b.int16(int16(len(oids)))
		for _, oid := range oids {
			b.uint32(oid)
		}
	})
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func readInt32(r *bufio.Reader) (int32, error) {
	var buf [4]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(buf[:])), nil
}

func readInt16(b []byte) int16 {
	return int16(binary.BigEndian.Uint16(b[:2]))
}

// cstring reads a null-terminated C string from a byte slice.
func cstring(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

// readCString returns (string, remaining bytes).
func readCString(b []byte) (string, []byte) {
	for i, c := range b {
		if c == 0 {
			return string(b[:i]), b[i+1:]
		}
	}
	return string(b), nil
}

var paramRe = regexp.MustCompile(`\$(\d+)`)

func countParams(query string) int {
	m := paramRe.FindAllStringSubmatch(query, -1)
	max := 0
	for _, mm := range m {
		n, _ := strconv.Atoi(mm[1])
		if n > max {
			max = n
		}
	}
	return max
}

func substituteParams(query string, params []interface{}) string {
	result := query
	for i, param := range params {
		placeholder := fmt.Sprintf("$%d", i+1)
		var replacement string
		if param == nil {
			replacement = "NULL"
		} else if b, ok := param.(bool); ok {
			// Handle bool specially
			if b {
				replacement = "true"
			} else {
				replacement = "false"
			}
		} else if s, ok := param.(string); ok {
			if _, err := strconv.Atoi(s); err == nil {
				replacement = s
			} else {
				replacement = fmt.Sprintf("'%v'", s)
			}
		} else {
			replacement = fmt.Sprintf("%v", param)
		}
		result = strings.ReplaceAll(result, placeholder, replacement)
	}
	return result
}

// ─────────────────────────────────────────────
// Array formatting
// ─────────────────────────────────────────────

func formatArrayAsText(value interface{}) string {
	switch arr := value.(type) {
	case []string:
		return pgArray(len(arr), func(i int) string {
			v := arr[i]
			if strings.ContainsAny(v, ",{}\" ") || v == "" {
				return `"` + strings.ReplaceAll(strings.ReplaceAll(v, `\`, `\\`), `"`, `\"`) + `"`
			}
			return v
		})
	case []int:
		return pgArray(len(arr), func(i int) string { return strconv.Itoa(arr[i]) })
	case []int32:
		return pgArray(len(arr), func(i int) string { return strconv.FormatInt(int64(arr[i]), 10) })
	case []int64:
		return pgArray(len(arr), func(i int) string { return strconv.FormatInt(arr[i], 10) })
	case []float32:
		return pgArray(len(arr), func(i int) string { return strconv.FormatFloat(float64(arr[i]), 'f', -1, 32) })
	case []float64:
		return pgArray(len(arr), func(i int) string { return strconv.FormatFloat(arr[i], 'f', -1, 64) })
	case []bool:
		return pgArray(len(arr), func(i int) string {
			if arr[i] {
				return "t"
			}
			return "f"
		})
	default:
		return fmt.Sprintf("%v", value)
	}
}

func pgArray(n int, elem func(int) string) string {
	if n == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.Grow(n*8 + 2)
	sb.WriteByte('{')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(elem(i))
	}
	sb.WriteByte('}')
	return sb.String()
}

// Suppress unused import if logger.Debugf uses zap fields indirectly
var _ = zap.Any
