// TODO пересмотреть wire сервер, разбить на файлы, слишком много строк кода
package wire

import (
	"fmt"
	"io"
	"net"
	"petacore/internal/logger"
	"petacore/internal/runtime/executor"
	"regexp"
	"strconv"

	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/runtime/rsql/visitor"
	"petacore/internal/runtime/system"
	"petacore/internal/storage"
	"strings"
	"time"

	"github.com/jackc/pgproto3/v2"
	"go.uber.org/zap"
)

// Session represents a client session with prepared statements and portals
type Session struct {
	preparedStatements map[string]*PreparedStatement
	portals            map[string]*PreparedStatement
	params             map[string]string
}

// NewSession creates a new session
func NewSession() *Session {
	return &Session{
		preparedStatements: make(map[string]*PreparedStatement),
		portals:            make(map[string]*PreparedStatement),
		params:             make(map[string]string),
	}
}

// WireServer представляет PostgreSQL wire protocol сервер
type WireServer struct {
	storage  *storage.DistributedStorageVClock
	listener net.Listener
	port     string
}

// NewWireServer создает новый wire сервер
func NewWireServer(storage *storage.DistributedStorageVClock, port string) *WireServer {
	return &WireServer{
		storage: storage,
		port:    port,
	}
}

// Start запускает сервер
func (ws *WireServer) Start() error {
	listener, err := net.Listen("tcp", ":"+ws.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", ws.port, err)
	}
	ws.listener = listener

	go ws.acceptConnections()
	return nil
}

// Stop останавливает сервер
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
			// Listener closed
			return
		}
		logger.Info("Accepted connection")
		go ws.handleConnection(conn)
	}
}

func (ws *WireServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	logger.Info("Accepted connection")

	session := NewSession()

	// Set timeout
	conn.SetDeadline(time.Now().Add(5 * time.Minute))

	// Create backend for message handling
	backend := pgproto3.NewBackend(pgproto3.NewChunkReader(conn), conn)

	// Read first message - could be SSL request or startup
	msg, err := backend.ReceiveStartupMessage()
	if err != nil {
		logger.Errorf("Error reading first message: %v", err)
		return
	}

	switch m := msg.(type) {
	case *pgproto3.SSLRequest:
		// We don't support SSL, send 'N'
		logger.Info("SSL request received, sending 'N'")
		conn.Write([]byte{'N'})
		// Now read startup message
		startupMessage, err := backend.ReceiveStartupMessage()
		if err != nil {
			logger.Errorf("Error reading startup message after SSL: %v", err)
			return
		}
		if sm, ok := startupMessage.(*pgproto3.StartupMessage); ok {
			ws.handleStartup(backend, sm, session)
		} else {
			logger.Errorf("Unexpected message after SSL: %T", startupMessage)
			return
		}
	case *pgproto3.StartupMessage:
		ws.handleStartup(backend, m, session)
	default:
		logger.Errorf("Unexpected first message: %T", m)
		return
	}

	// Main message loop
	for {
		msg, err := backend.Receive()
		if err != nil {
			if err == io.EOF {
				logger.Info("Connection closed by client")
				return
			}
			// Any other error (including unexpected EOF) means connection is broken
			logger.Errorf("Connection error: %v", err)
			return
		}

		switch msg := msg.(type) {
		case *pgproto3.Query:
			logger.Debugf("Query: %s", msg.String)
			ws.handleQuery(backend, msg.String, session)
		case *pgproto3.Parse:
			ws.handleParse(backend, msg, session)
		case *pgproto3.Bind:
			ws.handleBind(backend, msg, session)
		case *pgproto3.Describe:
			ws.handleDescribe(backend, msg, session)
		case *pgproto3.Execute:
			ws.handleExecute(backend, msg, session)
		case *pgproto3.Sync:
			backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		case *pgproto3.Flush:
			// Flush any pending output - in this implementation, Send() is synchronous
			logger.Info("Flush received")
		case *pgproto3.Terminate:
			logger.Info("Terminate received")
			return
		default:
			logger.Errorf("Unsupported message type: %T", msg)
			backend.Send(&pgproto3.ErrorResponse{
				Severity: "ERROR",
				Code:     "0A000",
				Message:  "unsupported message type",
			})
			backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		}
	}
}

func (ws *WireServer) handleStartup(backend *pgproto3.Backend, startupMessage *pgproto3.StartupMessage, session *Session) {
	// Send AuthenticationOk
	backend.Send(&pgproto3.AuthenticationOk{})

	// Send ParameterStatus
	backend.Send(&pgproto3.ParameterStatus{Name: "server_version", Value: "13.0.0"})
	backend.Send(&pgproto3.ParameterStatus{Name: "client_encoding", Value: "UTF8"})
	backend.Send(&pgproto3.ParameterStatus{Name: "server_encoding", Value: "UTF8"})
	backend.Send(&pgproto3.ParameterStatus{Name: "standard_conforming_strings", Value: "on"})
	backend.Send(&pgproto3.ParameterStatus{Name: "TimeZone", Value: "UTC"})
	backend.Send(&pgproto3.ParameterStatus{Name: "integer_datetimes", Value: "on"})
	backend.Send(&pgproto3.ParameterStatus{Name: "application_name", Value: ""})

	// Send BackendKeyData (dummy)
	backend.Send(&pgproto3.BackendKeyData{ProcessID: 1, SecretKey: 1})

	// Send ReadyForQuery
	backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
	logger.Info("Startup complete")
}

func (ws *WireServer) handleParse(backend *pgproto3.Backend, msg *pgproto3.Parse, session *Session) {
	logger.Debugf("Parse name: '%s', query: '%s'", msg.Name, msg.Query)
	var stmt statements.SQLStatement
	if strings.TrimSpace(msg.Query) == "" {
		logger.Debug("Empty query in Parse, creating EmptyStatement")
		stmt = &statements.EmptyStatement{}
	} else {
		var err error
		stmt, err = visitor.ParseSQL(msg.Query)
		if err != nil {
			logger.Errorf("Parse error: %v", err)
			backend.Send(&pgproto3.ErrorResponse{
				Severity: "ERROR",
				Code:     "42601",
				Message:  err.Error(),
			})
			backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
			return
		}
	}

	// Count parameters in query and determine their types
	paramCount := countParams(msg.Query)
	paramOIDs := ws.inferParameterTypes(stmt, paramCount)

	logger.Debug("paramsOIDs:", zap.Any("paramsOIDs", paramOIDs))

	session.preparedStatements[msg.Name] = &PreparedStatement{
		Query:     msg.Query,
		Stmt:      stmt,
		ParamOIDs: paramOIDs,
	}
	// Populate columns for description
	session.preparedStatements[msg.Name].Columns = ws.getFieldDescriptions(stmt)
	logger.Infof("Prepared statement '%s' created", msg.Name)
	backend.Send(&pgproto3.ParseComplete{})
}

// inferParameterTypes attempts to determine parameter types from statement context
func (ws *WireServer) inferParameterTypes(stmt statements.SQLStatement, paramCount int) []uint32 {
	paramOIDs := make([]uint32, paramCount)

	// Default to INT4 for all parameters as most common use case
	// This allows numeric IDs to work without type conversion
	for i := range paramOIDs {
		paramOIDs[i] = 23 // INT4
	}

	return paramOIDs
}

// colTypeToOID converts internal column type to PostgreSQL OID
func (ws *WireServer) colTypeToOID(colType table.ColType) uint32 {
	switch colType {
	case table.ColTypeInt:
		return 23 // INT4
	case table.ColTypeBigInt:
		return 20 // INT8
	case table.ColTypeFloat:
		return 701 // FLOAT8
	case table.ColTypeBool:
		return 16 // BOOL
	case table.ColTypeString:
		return 25 // TEXT
	default:
		return 25 // TEXT
	}
}

func (ws *WireServer) handleQuery(backend *pgproto3.Backend, query string, session *Session) {
	// Handle keep-alive queries and empty queries
	trimmedQuery := strings.TrimSpace(query)
	if trimmedQuery == "" || trimmedQuery == ";" || strings.ToLower(trimmedQuery) == "select 1" || strings.ToLower(trimmedQuery) == "keep alive" {
		// Send empty result for keep-alive queries
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("SELECT 1")})
		backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		return
	}

	// Parse and execute query
	stmt, err := visitor.ParseSQL(query)
	if err != nil {
		backend.Send(&pgproto3.ErrorResponse{
			Severity: "ERROR",
			Code:     "42601",
			Message:  err.Error(),
		})
		backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		return
	}

	logger.Debugf("Parsed statement %+v", stmt)

	result, err := executor.ExecuteStatement(stmt, ws.storage, session.params)
	if err != nil {
		backend.Send(&pgproto3.ErrorResponse{
			Severity: "ERROR",
			Code:     "XX000",
			Message:  err.Error(),
		})
		backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		return
	}

	// Send result based on statement type
	switch s := stmt.(type) {
	case *statements.CreateTableStatement:
		// DDL - send CommandComplete
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("CREATE TABLE")})
	case *statements.InsertStatement:
		// DML - send CommandComplete
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("INSERT")})
	case *statements.DropTableStatement:
		// DDL - send CommandComplete
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("DROP TABLE")})
	case *statements.TruncateTableStatement:
		// DDL - send CommandComplete
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("TRUNCATE TABLE")})
	case *statements.SelectStatement:
		ws.sendSelectResult(backend, s, result, true)
	case *statements.SetStatement:
		// SET - send CommandComplete
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("SET")})
		// case *statements.ShowStatement:
		// 	ws.sendDescribeResult(backend, result)
	}

	backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
}

func (ws *WireServer) handleDescribe(backend *pgproto3.Backend, msg *pgproto3.Describe, session *Session) {
	switch msg.ObjectType {
	case 'S':
		// Describe prepared statement
		ps, ok := session.preparedStatements[msg.Name]
		if !ok {
			backend.Send(&pgproto3.ErrorResponse{
				Severity: "ERROR",
				Code:     "26000",
				Message:  "prepared statement does not exist",
			})
			backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
			return
		}
		// Send ParameterDescription (empty for now)
		backend.Send(&pgproto3.ParameterDescription{ParameterOIDs: ps.ParamOIDs})
		// Send RowDescription
		backend.Send(&pgproto3.RowDescription{Fields: ps.Columns})
	case 'P':
		// Describe portal
		ps, ok := session.portals[msg.Name]
		if !ok {
			backend.Send(&pgproto3.ErrorResponse{
				Severity: "ERROR",
				Code:     "26000",
				Message:  "portal does not exist",
			})
			backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
			return
		}
		// Send RowDescription
		backend.Send(&pgproto3.RowDescription{Fields: ps.Columns})
	}
}

func (ws *WireServer) handleBind(backend *pgproto3.Backend, msg *pgproto3.Bind, session *Session) {
	logger.Infof("Bind prepared: '%s', portal: '%s'", msg.PreparedStatement, msg.DestinationPortal)
	ps, ok := session.preparedStatements[msg.PreparedStatement]
	if !ok {
		backend.Send(&pgproto3.ErrorResponse{
			Severity: "ERROR",
			Code:     "26000",
			Message:  "prepared statement does not exist",
		})
		backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		return
	}

	// Check parameter count
	expectedParams := len(ps.ParamOIDs)
	providedParams := len(msg.Parameters)
	if providedParams != expectedParams {
		backend.Send(&pgproto3.ErrorResponse{
			Severity: "ERROR",
			Code:     "08P01",
			Message:  fmt.Sprintf("Prepared statement \"%s\" requires %d parameters, but %d were provided", msg.PreparedStatement, expectedParams, providedParams),
		})
		backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		return
	}

	// Parse parameters from binary/text format
	params := make([]interface{}, len(msg.Parameters))
	for i, paramBytes := range msg.Parameters {
		if paramBytes == nil {
			params[i] = nil
		} else {
			// Determine format for this parameter
			format := int16(0) // default text
			if len(msg.ParameterFormatCodes) == 1 {
				// Single format applies to all parameters
				format = msg.ParameterFormatCodes[0]
			} else if i < len(msg.ParameterFormatCodes) {
				// Per-parameter format
				format = msg.ParameterFormatCodes[i]
			}

			if format == 0 {
				// Text format - just convert to string
				params[i] = string(paramBytes)
			} else {
				// Binary format - decode based on OID
				// For simplicity, decode as int32 if length is 4 bytes
				if len(paramBytes) == 4 {
					// INT4 binary format
					val := int32(paramBytes[0])<<24 | int32(paramBytes[1])<<16 | int32(paramBytes[2])<<8 | int32(paramBytes[3])
					params[i] = fmt.Sprintf("%d", val)
				} else {
					// Unknown binary format, treat as string
					params[i] = string(paramBytes)
				}
			}
		}
	}

	// Copy the prepared statement to the portal with result format codes and parameters
	portalPS := &PreparedStatement{
		Query:             ps.Query,
		Stmt:              ps.Stmt,
		Params:            params,
		Columns:           ps.Columns,
		ParamOIDs:         ps.ParamOIDs,
		ResultFormatCodes: msg.ResultFormatCodes,
	}
	session.portals[msg.DestinationPortal] = portalPS
	backend.Send(&pgproto3.BindComplete{})
}

var paramRe = regexp.MustCompile(`\$(\d+)`)

func countParams(query string) int {
	m := paramRe.FindAllStringSubmatch(query, -1)
	max := 0
	for _, mm := range m {
		// mm[1] = digits
		n, err := strconv.Atoi(mm[1])
		if err != nil {
			continue
		}
		if n > max {
			max = n
		}
	}
	return max
}

// TODO пересмотреть, сейчас несколько стратегий на FROM в select
func (ws *WireServer) getFieldDescriptions(stmt statements.SQLStatement) []pgproto3.FieldDescription {
	switch s := stmt.(type) {
	case *statements.SelectStatement:
		var fields []pgproto3.FieldDescription
		if s.From == nil {
			// TODO избавиться от хардкода
			return fields
		}

		tableName := s.From.TableName

		if tableName == "" {
			// System functions
			for _, col := range s.Columns {
				var name string
				if col.Alias != "" {
					name = col.Alias
				} else if col.Function != nil {
					name = col.Function.Name
				} else {
					name = col.ColumnName
				}
				fields = append(fields, pgproto3.FieldDescription{
					Name:                 []byte(name),
					TableOID:             0,
					TableAttributeNumber: 0,
					DataTypeOID:          25, // TEXT
					DataTypeSize:         -1,
					TypeModifier:         -1,
					Format:               0,
				})
			}
		} else if system.IsSystemTable(tableName) {
			// For system tables, determine columns from stmt.Columns
			for _, col := range s.Columns {
				var name string
				var oid uint32 = 25 // Default TEXT
				if col.ColumnName == "*" {
					// For SELECT *, assume ssl column for pg_stat_ssl
					if tableName == "pg_stat_ssl" {
						name = "ssl"
						oid = 16 // BOOL
					} else {
						name = "column1"
					}
				} else {
					name = col.ColumnName
					// Set OID based on column name if known
					if name == "ssl" {
						oid = 16 // BOOL
					}
				}
				fields = append(fields, pgproto3.FieldDescription{
					Name:                 []byte(name),
					TableOID:             0,
					TableAttributeNumber: 0,
					DataTypeOID:          oid,
					DataTypeSize:         -1,
					TypeModifier:         -1,
					Format:               0,
				})
			}
		} else {
			// For normal tables, need to get from metadata
			// Simplified
			fields = []pgproto3.FieldDescription{}
		}
		return fields
	default:
		return []pgproto3.FieldDescription{}
	}
}

func (ws *WireServer) handleExecute(backend *pgproto3.Backend, msg *pgproto3.Execute, session *Session) {
	logger.Debugf("Execute portal: '%s'", msg.Portal)
	ps, ok := session.portals[msg.Portal]
	if !ok {
		logger.Errorf("Portal %s does not exist", msg.Portal)
		backend.Send(&pgproto3.ErrorResponse{
			Severity: "ERROR",
			Code:     "26000",
			Message:  "portal does not exist",
		})
		backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		return
	}
	logger.Debugf("Executing prepared statement: %s", ps.Query)
	logger.Debugf("Statement AST: %+v", ps.Stmt)
	logger.Debugf("Parameters: %+v", ps.Params)

	// If there are parameters, we need to re-parse with substituted values
	var stmt statements.SQLStatement
	if len(ps.Params) > 0 {
		// Substitute parameters in query
		substitutedQuery := ps.Query
		for i, param := range ps.Params {
			placeholder := fmt.Sprintf("$%d", i+1)
			var replacement string
			if param == nil {
				replacement = "NULL"
			} else {
				// Check if it's a string parameter - only quote strings
				paramStr, isString := param.(string)
				if isString {
					// Try to parse as number first
					if _, err := strconv.Atoi(paramStr); err == nil {
						// It's a numeric string, don't quote
						replacement = paramStr
					} else {
						// It's a text string, quote it
						replacement = fmt.Sprintf("'%v'", param)
					}
				} else {
					// It's already a non-string type (int, float, etc), don't quote
					replacement = fmt.Sprintf("%v", param)
				}
			}
			substitutedQuery = strings.ReplaceAll(substitutedQuery, placeholder, replacement)
		}
		logger.Debugf("Substituted query: %s", substitutedQuery)

		// Re-parse with substituted values
		var err error
		stmt, err = visitor.ParseSQL(substitutedQuery)
		if err != nil {
			logger.Errorf("Parse error after substitution: %v", err)
			backend.Send(&pgproto3.ErrorResponse{
				Severity: "ERROR",
				Code:     "42601",
				Message:  err.Error(),
			})
			backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
			return
		}
	} else {
		stmt = ps.Stmt
	}

	// Execute the statement
	result, err := executor.ExecuteStatement(stmt, ws.storage, session.params)
	if err != nil {
		logger.Errorf("Execute error: %v", err)
		backend.Send(&pgproto3.ErrorResponse{
			Severity: "ERROR",
			Code:     "XX000",
			Message:  err.Error(),
		})
		backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		return
	}

	logger.Debugf("Execute result: %+v", result)

	// Send result based on statement type (similar to handleQuery)
	switch s := ps.Stmt.(type) {
	case *statements.EmptyStatement:
		backend.Send(&pgproto3.EmptyQueryResponse{})
	case *statements.CreateTableStatement:
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("CREATE TABLE")})
	case *statements.InsertStatement:
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("INSERT")})
	case *statements.DropTableStatement:
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("DROP TABLE")})
	case *statements.SelectStatement:
		// Force text format for system tables since we can't reliably encode them in binary
		formatCodes := ps.ResultFormatCodes
		// TODO избавиться от постоянных проверок IsSystemTable, вынести в стратегию выполнения
		// if s.From != nil && system.IsSystemTable(s.From.TableName) {
		// 	// Override to text format for system tables
		// 	formatCodes = []int16{0}
		// }
		ws.sendSelectResultWithFormats(backend, s, result, true, formatCodes)
	case *statements.SetStatement:
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("SET")})
		// case *statements.ShowStatement:
		// 	ws.sendDescribeResult(backend, result)
	}
}

func (ws *WireServer) sendSelectResultWithFormats(backend *pgproto3.Backend, stmt *statements.SelectStatement, result *table.ExecuteResult, sendRowDesc bool, formatCodes []int16) {
	logger.Debugf("DEBUG sendSelectResultWithFormats: formatCodes=%v", formatCodes)
	logger.Debugf("DEBUG: Select result: %+v\n", result)
	// var rows []map[string]interface{}
	// var columns []string
	// var columnTypes []table.ColType

	// if resultMap, ok := result.(map[string]interface{}); ok {
	// 	// New result format
	// 	columns = resultMap["columns"].([]string)
	// 	columnTypes = resultMap["columnTypes"].([]table.ColType)
	// 	rows = resultMap["rows"].([]map[string]interface{})
	// } else {
	// 	backend.Send(&pgproto3.ErrorResponse{
	// 		Severity: "ERROR",
	// 		Code:     "XX000",
	// 		Message:  "invalid result type",
	// 	})
	// 	return
	// }

	if len(result.Rows) == 0 {
		// No data
		if sendRowDesc {
			rowDesc := &pgproto3.RowDescription{}
			for i, column := range result.Columns {
				oid := uint32(25) // Default TEXT
				switch column.Type {
				case table.ColTypeInt:
					oid = 23 // INT4
				case table.ColTypeBigInt:
					oid = 20 // INT8
				case table.ColTypeFloat:
					oid = 701 // FLOAT8
				case table.ColTypeBool:
					oid = 16 // BOOL
				case table.ColTypeString:
					oid = 25 // TEXT
				}

				// Determine format for this column
				format := int16(0) // default text
				if len(formatCodes) == 1 {
					// Single format applies to all columns
					format = formatCodes[0]
				} else if i < len(formatCodes) {
					// Per-column format
					format = formatCodes[i]
				}

				rowDesc.Fields = append(rowDesc.Fields, pgproto3.FieldDescription{
					Name:                 []byte(column.Name),
					TableOID:             0,
					TableAttributeNumber: 0,
					DataTypeOID:          oid,
					DataTypeSize:         -1,
					TypeModifier:         -1,
					Format:               format,
				})
			}
			backend.Send(rowDesc)
		}
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte(fmt.Sprintf("SELECT %d", 0))})
		return
	}

	// RowDescription
	if sendRowDesc {
		rowDesc := &pgproto3.RowDescription{}
		for i, column := range result.Columns {
			oid := uint32(25) // Default TEXT
			switch column.Type {
			case table.ColTypeInt:
				oid = 23 // INT4
			case table.ColTypeBigInt:
				oid = 20 // INT8
			case table.ColTypeFloat:
				oid = 701 // FLOAT8
			case table.ColTypeBool:
				oid = 16 // BOOL
			case table.ColTypeString:
				oid = 25 // TEXT
			}

			// Determine format for this column
			format := int16(0) // default text
			if len(formatCodes) == 1 {
				// Single format applies to all columns
				format = formatCodes[0]
			} else if i < len(formatCodes) {
				// Per-column format
				format = formatCodes[i]
			}

			rowDesc.Fields = append(rowDesc.Fields, pgproto3.FieldDescription{
				Name:                 []byte(column.Name),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          oid,
				DataTypeSize:         -1,
				TypeModifier:         -1,
				Format:               format,
			})
		}
		backend.Send(rowDesc)
	}

	// Data rows
	for _, row := range result.Rows {
		dataRow := &pgproto3.DataRow{}
		// Determine format for this column
		// format := int16(0) // default text
		// if len(formatCodes) == 1 {
		// 	// Single format applies to all columns
		// 	format = formatCodes[0]
		// } else if i < len(formatCodes) {
		// 	// Per-column format
		// 	format = formatCodes[i]
		// }

		// if format == 1 {
		// 	// Binary format
		// 	dataRow.Values = append(dataRow.Values, encodeBinary(v, i, columnTypes))
		// } else {
		// Text format
		for _, value := range row {
			dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", value)))
		}
		// }
		backend.Send(dataRow)
	}

	backend.Send(&pgproto3.CommandComplete{CommandTag: []byte(fmt.Sprintf("SELECT %d", len(result.Rows)))})
}

// sendSelectResult is a wrapper that calls sendSelectResultWithFormats with text format (0)
func (ws *WireServer) sendSelectResult(backend *pgproto3.Backend, stmt *statements.SelectStatement, result *table.ExecuteResult, sendRowDesc bool) {
	ws.sendSelectResultWithFormats(backend, stmt, result, sendRowDesc, []int16{0})
}

// encodeBinary encodes a value in PostgreSQL binary format
func encodeBinary(v interface{}, colIdx int, columnTypes []table.ColType) []byte {
	if colIdx >= len(columnTypes) {
		// Fallback to text
		return []byte(fmt.Sprintf("%v", v))
	}

	logger.Debugf("DEBUG encodeBinary: colIdx=%d, type=%v, value=%v (Go type %T)", colIdx, columnTypes[colIdx], v, v)

	switch columnTypes[colIdx] {
	case table.ColTypeInt:
		// INT4 - 4 bytes, big-endian
		var val int32
		switch v := v.(type) {
		case int:
			val = int32(v)
		case int32:
			val = v
		case int64:
			val = int32(v)
		case float64:
			val = int32(v)
		case float32:
			val = int32(v)
		default:
			return []byte(fmt.Sprintf("%v", v))
		}
		buf := make([]byte, 4)
		buf[0] = byte(val >> 24)
		buf[1] = byte(val >> 16)
		buf[2] = byte(val >> 8)
		buf[3] = byte(val)
		return buf
	case table.ColTypeString:
		// TEXT binary format is just the string bytes
		if s, ok := v.(string); ok {
			return []byte(s)
		}
		return []byte(fmt.Sprintf("%v", v))
	default:
		// For other types, convert to string
		return []byte(fmt.Sprintf("%v", v))
	}
}

// func (ws *WireServer) sendDescribeResult(backend *pgproto3.Backend, result *table.ExecuteResult) {
// 	rows := result.Rows
// 	if rows == nil {
// 		backend.Send(&pgproto3.ErrorResponse{
// 			Severity: "ERROR",
// 			Code:     "XX000",
// 			Message:  "invalid result type for describe",
// 		})
// 		return
// 	}

// 	if len(rows) == 0 {
// 		// No data
// 		backend.Send(&pgproto3.RowDescription{Fields: []pgproto3.FieldDescription{}})
// 		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("SHOW 0")})
// 		return
// 	}

// 	// Get columns from first row
// 	var columns []string
// 	for k := range rows[0] {
// 		columns = append(columns, k)
// 	}
// 	sort.Strings(columns)

// 	// RowDescription
// 	rowDesc := &pgproto3.RowDescription{}
// 	for _, colName := range columns {
// 		rowDesc.Fields = append(rowDesc.Fields, pgproto3.FieldDescription{
// 			Name:                 []byte(colName),
// 			TableOID:             0,
// 			TableAttributeNumber: 0,
// 			DataTypeOID:          25, // TEXT
// 			DataTypeSize:         -1,
// 			TypeModifier:         -1,
// 			Format:               0,
// 		})
// 	}
// 	backend.Send(rowDesc)

// 	// Send DataRows
// 	for _, row := range rows {
// 		dataRow := &pgproto3.DataRow{}
// 		for _, colName := range columns {
// 			if v, exists := row[colName]; exists && v != nil {
// 				dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", v)))
// 			} else {
// 				dataRow.Values = append(dataRow.Values, nil)
// 			}
// 		}
// 		backend.Send(dataRow)
// 	}

// 	backend.Send(&pgproto3.CommandComplete{CommandTag: []byte(fmt.Sprintf("SHOW %d", len(rows)))})
// }
