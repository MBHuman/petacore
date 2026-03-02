// TODO пересмотреть wire сервер, разбить на файлы, слишком много строк кода
package wire

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgproto3/v2"
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

// Session stores prepared statements, portals and session params
type Session struct {
	preparedStatements map[string]*PreparedStatement
	portals            map[string]*PreparedStatement
	params             map[string]string
}

// NewSession creates a new session with initialized maps
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
	// Создаем allocator для этого запроса
	allocator := pmem.NewArena(1024 * 1024) // 1MB arena
	defer allocator.Close()

	var stmt statements.SQLStatement
	if strings.TrimSpace(msg.Query) == "" {
		logger.Debug("Empty query in Parse, creating EmptyStatement")
		stmt = &statements.EmptyStatement{}
	} else {
		var err error
		stmt, err = visitor.ParseSQL(allocator, msg.Query)
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
	// Build placeholder FieldDescriptions for SELECT prepared statements when possible.
	// If the SELECT contains a '*' (IsSelectAll) we cannot know column count yet, so leave Columns nil.
	inferredCols := ws.getFieldDescriptions(stmt)
	if sel, ok := stmt.(*statements.SelectStatement); ok {
		hasSelectAll := false
		if sel.Primary != nil {
			for _, c := range sel.Primary.Columns {
				if c.IsSelectAll {
					hasSelectAll = true
					break
				}
			}
		}

		if hasSelectAll {
			// Defer RowDescription until Execute when we know actual columns
			session.preparedStatements[msg.Name].Columns = nil
		} else if sel.Primary != nil && len(inferredCols) == len(sel.Primary.Columns) && len(inferredCols) > 0 {
			// We were able to infer exact column descriptions (e.g., SELECT expressions without table)
			session.preparedStatements[msg.Name].Columns = inferredCols
		} else if sel.Primary != nil {
			// Create placeholder descriptions from AST (names only, default to TEXT)
			var fields []pgproto3.FieldDescription
			for i, col := range sel.Primary.Columns {
				var name string
				if col.Alias != "" {
					name = col.Alias
				} else if col.Function != nil {
					name = col.Function.Name
				} else if col.ColumnName != "" {
					name = col.ColumnName
				} else {
					name = fmt.Sprintf("column%d", i+1)
				}
				oid := uint32(25)
				if col.Function != nil {
					if fn, ok := functions.GetRegisteredFunction(col.Function.Name); ok {
						oid = uint32(fn.GetFunction().ProRetType)
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
			session.preparedStatements[msg.Name].Columns = fields
		} else {
			session.preparedStatements[msg.Name].Columns = nil
		}
	} else {
		session.preparedStatements[msg.Name].Columns = inferredCols
	}
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

func (ws *WireServer) handleQuery(backend *pgproto3.Backend, query string, session *Session) {
	// Handle keep-alive queries and empty queries
	trimmedQuery := strings.TrimSpace(query)
	if trimmedQuery == "" || trimmedQuery == ";" || strings.ToLower(trimmedQuery) == "select 1" || strings.ToLower(trimmedQuery) == "keep alive" {
		// Send empty result for keep-alive queries
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("SELECT 1")})
		backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		return
	}

	// Создаем allocator для этого запроса
	allocator := pmem.NewArena(26 * 1024 * 1024) // 26MB arena

	// Parse and execute query
	stmt, err := visitor.ParseSQL(allocator, query)
	if err != nil {
		backend.Send(&pgproto3.ErrorResponse{
			Severity: "ERROR",
			Code:     "42601",
			Message:  err.Error(),
		})
		backend.Send(&pgproto3.ReadyForQuery{TxStatus: 'I'})
		return
	}

	logger.Debugf("Parsed statement: ", zap.Any("stmt", stmt))

	result, err := executor.ExecuteStatement(allocator, stmt, ws.storage, session.params)
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
	switch stmt.(type) {
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
		ws.sendSelectResult(backend, result, true)
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
		// Do not execute the statement during Describe to infer columns - this
		// may observe a different projection/state than will be produced on
		// Execute and introduce mismatches. Only send RowDescription when we
		// already have `ps.Columns` cached (set earlier by an explicit Prepare
		// flow). This keeps Describe idempotent and avoids side-effects.
		logger.Debugf("Describe prepared statement '%s': ps.Columns=%d (not sending RowDescription on Describe)", msg.Name, len(ps.Columns))
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
		// Do not execute the statement during Describe to infer columns for the
		// same reasons as for prepared statements above. Only send a
		// RowDescription when portal already carries column metadata.
		logger.Debugf("Describe portal '%s': ps.Columns=%d (not sending RowDescription on Describe)", msg.Name, len(ps.Columns))
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

		// With new structure, SelectStatement has Primary or Combined
		// For now, only handle Primary SELECT for field descriptions
		if s.Primary == nil {
			// Combined statement or no primary - return empty for now
			return fields
		}

		primary := s.Primary
		if primary.From == nil {
			// When there's no FROM (SELECT expressions without table), infer column types
			for _, col := range primary.Columns {
				var name string
				var oid uint32 = 25 // default TEXT
				if col.Alias != "" {
					name = col.Alias
				} else if col.Function != nil {
					name = col.Function.Name
					// Try to infer type from registered functions
					if fn, ok := functions.GetRegisteredFunction(col.Function.Name); ok {
						oid = uint32(fn.GetFunction().ProRetType)
					} else {
						// Fallback: if aggregate or known, keep TEXT
						oid = 25
					}
				} else if col.ColumnName != "" {
					name = col.ColumnName
				} else {
					name = "?column?"
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
			return fields
		}

		tableName := primary.From.TableName

		if tableName == "" {
			// System functions
			for _, col := range primary.Columns {
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
			for _, col := range primary.Columns {
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

	// Создаем allocator для этой операции
	allocator := pmem.NewArena(26 * 1024 * 1024) // 26MB arena
	defer allocator.Close()

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
		stmt, err = visitor.ParseSQL(allocator, substitutedQuery)
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
	result, err := executor.ExecuteStatement(allocator, stmt, ws.storage, session.params)
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
	switch ps.Stmt.(type) {
	case *statements.EmptyStatement:
		backend.Send(&pgproto3.EmptyQueryResponse{})
	case *statements.CreateTableStatement:
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("CREATE TABLE")})
	case *statements.InsertStatement:
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("INSERT")})
	case *statements.DropTableStatement:
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("DROP TABLE")})
	case *statements.SelectStatement:
		// For portals, if a RowDescription was already provided during Describe, send that
		formatCodes := ps.ResultFormatCodes
		if len(ps.Columns) > 0 {
			// Send RowDescription based on prepared/portal columns, applying requested format codes
			rowDesc := &pgproto3.RowDescription{}
			for i, fd := range ps.Columns {
				// Determine OID: prefer actual executed result column type when available
				var oid uint32 = uint32(fd.DataTypeOID)
				if result.Schema != nil && i < len(result.Schema.Fields) {
					oid = uint32(result.Schema.Fields[i].OID)
				}
				// Force text format for data rows for now
				format := int16(0)

				rowDesc.Fields = append(rowDesc.Fields, pgproto3.FieldDescription{
					Name:                 fd.Name,
					TableOID:             fd.TableOID,
					TableAttributeNumber: fd.TableAttributeNumber,
					DataTypeOID:          oid,
					DataTypeSize:         fd.DataTypeSize,
					TypeModifier:         fd.TypeModifier,
					Format:               format,
				})
			}
			logger.Debugf("Execute (portal) - sending RowDescription: ps.Columns=%d, result.Schema.Fields=%d", len(ps.Columns), func() int {
				if result.Schema != nil {
					return len(result.Schema.Fields)
				}
				return 0
			}())
			backend.Send(rowDesc)
			// Send only data rows; do not resend RowDescription. Use portal's column count as expected.
			ws.sendSelectResultWithFormats(backend, result, false, formatCodes, len(ps.Columns))
		} else {
			// No prior description; let sendSelectResultWithFormats send RowDescription
			colCount := 0
			if result.Schema != nil {
				colCount = len(result.Schema.Fields)
			}
			logger.Debugf("Execute - no prior description, will send RowDescription from result.Schema: result.Schema.Fields=%d", colCount)
			ws.sendSelectResultWithFormats(backend, result, true, formatCodes, colCount)
		}
	case *statements.SetStatement:
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte("SET")})
		// case *statements.ShowStatement:
		// 	ws.sendDescribeResult(backend, result)
	}
}

func (ws *WireServer) sendSelectResultWithFormats(backend *pgproto3.Backend, result *table.ExecuteResult, sendRowDesc bool, formatCodes []int16, expectedCols int) {
	schemaFieldCount := 0
	if result.Schema != nil {
		schemaFieldCount = len(result.Schema.Fields)
	}
	logger.Debugf("sendSelectResultWithFormats: sendRowDesc=%v, result.Schema.Fields=%d, expectedCols=%d, rows=%d, formatCodes=%v", sendRowDesc, schemaFieldCount, expectedCols, len(result.Rows), formatCodes)

	if len(result.Rows) == 0 {
		// No data
		if sendRowDesc && result.Schema != nil {
			rowDesc := &pgproto3.RowDescription{}
			for _, field := range result.Schema.Fields {
				oid := field.OID
				format := int16(0)
				rowDesc.Fields = append(rowDesc.Fields, pgproto3.FieldDescription{
					Name:                 []byte(field.Name),
					TableOID:             0,
					TableAttributeNumber: 0,
					DataTypeOID:          uint32(oid),
					DataTypeSize:         -1,
					TypeModifier:         -1,
					Format:               format,
				})
			}
			logger.Debug("Sending RowDescription", zap.Int("fields", len(rowDesc.Fields)), zap.Any("fieldsDesc", rowDesc.Fields))
			backend.Send(rowDesc)
		}
		backend.Send(&pgproto3.CommandComplete{CommandTag: []byte(fmt.Sprintf("SELECT %d", 0))})
		return
	}

	// RowDescription
	if sendRowDesc && result.Schema != nil {
		rowDesc := &pgproto3.RowDescription{}
		logger.Debug("sendSelectResultWithFormats: columns", zap.Any("result.Schema.Fields", result.Schema.Fields))
		for _, field := range result.Schema.Fields {
			oid := field.OID
			format := int16(0)
			rowDesc.Fields = append(rowDesc.Fields, pgproto3.FieldDescription{
				Name:                 []byte(field.Name),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          uint32(oid),
				DataTypeSize:         -1,
				TypeModifier:         -1,
				Format:               format,
			})
		}
		logger.Debug("Sending RowDescription", zap.Int("fields", len(rowDesc.Fields)), zap.Any("fieldsDesc", rowDesc.Fields))
		backend.Send(rowDesc)
	}

	// Data rows
	for rIdx, row := range result.Rows {
		if rIdx < 5 && result.Schema != nil {
			logger.Debugf("DataRow %d: columns=%d", rIdx, len(result.Schema.Fields))
		}
		logger.Debug("DEBUG: Sending row:", zap.Any("row", row))

		// Determine how many columns we must emit for this DataRow. If we
		// just sent a RowDescription in this call, use the result.Schema.Fields
		// length; otherwise use the expectedCols provided by the caller
		// (which should match the RowDescription already sent earlier).
		colCount := expectedCols
		if sendRowDesc && result.Schema != nil {
			colCount = len(result.Schema.Fields)
		}

		dataRow := &pgproto3.DataRow{}
		for i := 0; i < colCount; i++ {
			// Extract field from binary row
			var value interface{}
			if result.Schema != nil && i < len(result.Schema.Fields) {
				buf, oid, err := result.Schema.GetField(row, i)
				if err != nil {
					value = nil
				} else if buf == nil {
					value = nil
				} else {
					// Deserialize the value and extract native Go type
					if bt, err := serializers.DeserializeGeneric(buf, oid); err == nil && bt != nil {
						value = bt.IntoGo()
					} else {
						value = nil
					}
				}
			} else {
				value = nil
			}

			// Determine format for this column
			format := int16(0)
			if len(formatCodes) == 1 {
				format = formatCodes[0]
			} else if i < len(formatCodes) {
				format = formatCodes[i]
			}

			if rIdx < 3 {
				// Log per-column format info for first few rows to help debug
				colOID := ptypes.PTypeText
				if result.Schema != nil && i < len(result.Schema.Fields) {
					colOID = result.Schema.Fields[i].OID
				}
				logger.Debugf("Row %d col %d: format=%d, colOID=%d, valueType=%T", rIdx, i, format, colOID, value)
			}

			if value == nil {
				dataRow.Values = append(dataRow.Values, nil)
				continue
			}

			if format == 1 {
				// Binary format encoding based on column OID. If the
				// result doesn't have type information for this index,
				// default to string behavior.
				colOID := ptypes.PTypeText
				if result.Schema != nil && i < len(result.Schema.Fields) {
					colOID = result.Schema.Fields[i].OID
				}
				switch colOID {
				case ptypes.PTypeBool:
					if b, ok := value.(bool); ok {
						if b {
							dataRow.Values = append(dataRow.Values, []byte{1})
						} else {
							dataRow.Values = append(dataRow.Values, []byte{0})
						}
					} else {
						dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", value)))
					}
				case ptypes.PTypeInt2, ptypes.PTypeInt4:
					var v32 int32
					appended := false
					switch vv := value.(type) {
					case int:
						v32 = int32(vv)
					case int32:
						v32 = vv
					case int64:
						v32 = int32(vv)
					case string:
						if parsed, err := strconv.ParseInt(vv, 10, 32); err == nil {
							v32 = int32(parsed)
						} else {
							dataRow.Values = append(dataRow.Values, []byte(vv))
							appended = true
						}
					default:
						dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", vv)))
						appended = true
					}
					if appended {
						break
					}
					buf := make([]byte, 4)
					binary.BigEndian.PutUint32(buf, uint32(v32))
					dataRow.Values = append(dataRow.Values, buf)
				case ptypes.PTypeInt8:
					var v64 int64
					appended := false
					switch vv := value.(type) {
					case int:
						v64 = int64(vv)
					case int32:
						v64 = int64(vv)
					case int64:
						v64 = vv
					case string:
						if parsed, err := strconv.ParseInt(vv, 10, 64); err == nil {
							v64 = parsed
						} else {
							dataRow.Values = append(dataRow.Values, []byte(vv))
							appended = true
						}
					default:
						dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", vv)))
						appended = true
					}
					if appended {
						break
					}
					buf8 := make([]byte, 8)
					binary.BigEndian.PutUint64(buf8, uint64(v64))
					dataRow.Values = append(dataRow.Values, buf8)
				case ptypes.PTypeTimestamp, ptypes.PTypeTimestampz:
					logger.Debug("Handling timestamp column in binary format", zap.Any("value", value), zap.Any("valueType", fmt.Sprintf("%T", value)))
					// Return timestamps as int64 microseconds since epoch in binary format
					var v64 int64
					switch vv := value.(type) {
					case int64:
						v64 = vv
					case int:
						v64 = int64(vv)
					case string:
						if parsed, err := strconv.ParseInt(vv, 10, 64); err == nil {
							v64 = parsed
						} else {
							dataRow.Values = append(dataRow.Values, []byte(vv))
							continue
						}
					default:
						dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", vv)))
						continue
					}
					timestampValue := time.UnixMicro(v64)
					dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", timestampValue)))
				case ptypes.PTypeFloat8:
					var f64 float64
					appended := false
					switch vv := value.(type) {
					case float32:
						f64 = float64(vv)
					case float64:
						f64 = vv
					case string:
						if parsed, err := strconv.ParseFloat(vv, 64); err == nil {
							f64 = parsed
						} else {
							dataRow.Values = append(dataRow.Values, []byte(vv))
							appended = true
						}
					default:
						dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", vv)))
						appended = true
					}
					if appended {
						break
					}
					buff := make([]byte, 8)
					binary.BigEndian.PutUint64(buff, math.Float64bits(f64))
					dataRow.Values = append(dataRow.Values, buff)
				case ptypes.PTypeText, ptypes.PTypeVarchar:
					switch vv := value.(type) {
					case string:
						dataRow.Values = append(dataRow.Values, []byte(vv))
					default:
						dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", vv)))
					}
				// Array types (binary format) - using text representation
				case 1007, 1016: // _int4[], _int8[] and other array types
					// For arrays in binary format, we'll use text representation for simplicity
					dataRow.Values = append(dataRow.Values, []byte(formatArrayAsText(value)))
				default:
					dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", value)))
				}
			} else {
				// Text format
				// Check if we have column type info to format timestamps specially
				colOID := ptypes.PTypeText
				if result.Schema != nil && i < len(result.Schema.Fields) {
					colOID = result.Schema.Fields[i].OID
				}

				// Special handling for timestamp types
				if colOID == ptypes.PTypeTimestamp || colOID == ptypes.PTypeTimestampz {
					var formatted string
					handled := false
					switch ts := value.(type) {
					case int64:
						t := time.UnixMicro(ts)
						if colOID == ptypes.PTypeTimestampz {
							formatted = t.Format("2006-01-02 15:04:05.999999-07")
						} else {
							formatted = t.UTC().Format("2006-01-02 15:04:05.999999")
						}
						handled = true
					case *time.Time:
						if ts != nil {
							if colOID == ptypes.PTypeTimestampz {
								formatted = ts.Format("2006-01-02 15:04:05.999999-07")
							} else {
								formatted = ts.UTC().Format("2006-01-02 15:04:05.999999")
							}
							handled = true
						}
					}
					if handled {
						dataRow.Values = append(dataRow.Values, []byte(formatted))
						continue
					}
				}

				switch v := value.(type) {
				case bool:
					if v {
						dataRow.Values = append(dataRow.Values, []byte("t"))
					} else {
						dataRow.Values = append(dataRow.Values, []byte("f"))
					}
				case int, int8, int16, int32, int64:
					dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", v)))
				case uint, uint8, uint16, uint32, uint64:
					dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", v)))
				case float32:
					dataRow.Values = append(dataRow.Values, []byte(strconv.FormatFloat(float64(v), 'f', -1, 32)))
				case float64:
					dataRow.Values = append(dataRow.Values, []byte(strconv.FormatFloat(v, 'f', -1, 64)))
				case string:
					dataRow.Values = append(dataRow.Values, []byte(v))
				case []string, []int, []int32, []int64, []float32, []float64, []bool:
					// Array types - format as PostgreSQL array text: {val1,val2,val3}
					dataRow.Values = append(dataRow.Values, []byte(formatArrayAsText(v)))
				default:
					dataRow.Values = append(dataRow.Values, []byte(fmt.Sprintf("%v", v)))
				}
			}
		}

		// Sanity log: actual values count must match the RowDescription/colCount
		logger.Debugf("Sending DataRow: values=%d, colCount=%d, sourceRowLen=%d", len(dataRow.Values), colCount, row.FieldCount())
		backend.Send(dataRow)
	}

	backend.Send(&pgproto3.CommandComplete{CommandTag: []byte(fmt.Sprintf("SELECT %d", len(result.Rows)))})
}

// sendSelectResult is a wrapper that calls sendSelectResultWithFormats with text format (0)
func (ws *WireServer) sendSelectResult(backend *pgproto3.Backend, result *table.ExecuteResult, sendRowDesc bool) {
	ws.sendSelectResultWithFormats(backend, result, sendRowDesc, []int16{0}, len(result.Schema.Fields))
}

// formatArrayAsText formats a Go slice as a PostgreSQL array text representation: {val1,val2,val3}
func formatArrayAsText(value interface{}) string {
	switch arr := value.(type) {
	case []string:
		if len(arr) == 0 {
			return "{}"
		}
		var sb strings.Builder
		sb.WriteString("{")
		for i, v := range arr {
			if i > 0 {
				sb.WriteString(",")
			}
			// Quote strings and escape special characters if needed
			if strings.ContainsAny(v, ",{}\" ") || v == "" {
				sb.WriteString(`"`)
				sb.WriteString(strings.ReplaceAll(strings.ReplaceAll(v, `\`, `\\`), `"`, `\"`))
				sb.WriteString(`"`)
			} else {
				sb.WriteString(v)
			}
		}
		sb.WriteString("}")
		return sb.String()
	case []int:
		if len(arr) == 0 {
			return "{}"
		}
		var sb strings.Builder
		sb.WriteString("{")
		for i, v := range arr {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteString("}")
		return sb.String()
	case []int32:
		if len(arr) == 0 {
			return "{}"
		}
		var sb strings.Builder
		sb.WriteString("{")
		for i, v := range arr {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(strconv.FormatInt(int64(v), 10))
		}
		sb.WriteString("}")
		return sb.String()
	case []int64:
		if len(arr) == 0 {
			return "{}"
		}
		var sb strings.Builder
		sb.WriteString("{")
		for i, v := range arr {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteString("}")
		return sb.String()
	case []float32:
		if len(arr) == 0 {
			return "{}"
		}
		var sb strings.Builder
		sb.WriteString("{")
		for i, v := range arr {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
		}
		sb.WriteString("}")
		return sb.String()
	case []float64:
		if len(arr) == 0 {
			return "{}"
		}
		var sb strings.Builder
		sb.WriteString("{")
		for i, v := range arr {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
		}
		sb.WriteString("}")
		return sb.String()
	case []bool:
		if len(arr) == 0 {
			return "{}"
		}
		var sb strings.Builder
		sb.WriteString("{")
		for i, v := range arr {
			if i > 0 {
				sb.WriteString(",")
			}
			if v {
				sb.WriteString("t")
			} else {
				sb.WriteString("f")
			}
		}
		sb.WriteString("}")
		return sb.String()
	default:
		// Fallback to standard string representation
		return fmt.Sprintf("%v", value)
	}
}
