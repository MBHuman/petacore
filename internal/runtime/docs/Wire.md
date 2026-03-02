# PostgreSQL Wire Protocol

## Overview

The Wire Protocol module implements PostgreSQL Frontend/Backend Protocol v3.0, enabling compatibility with PostgreSQL clients and tools. This allows any tool that speaks the PostgreSQL protocol to connect to PetaCore.

## Protocol Architecture

### Interaction Model

```
Client                      Server
  │                          │
  ├─── Startup Message ────→ │
  │                          │
  │ ← Authentication Request │
  │                          │
  ├─── Authentication ────→ │
  │                          │
  │ ← Authentication OK      │
  │                          │
  │ ← Parameter Status       │
  │                          │
  │ ← Ready for Query        │
  │                          │
  ├─── Query ────────────→ │
  │                          │
  │ ← Row Description        │
  │                          │
  │ ← Data Row (multiple)    │
  │                          │
  │ ← Command Complete       │
  │                          │
  │ ← Ready for Query        │
  │                          │
  ├─── Terminate ────────→ │
```

## Message Types

### Client Messages

#### Startup Message

Initiates connection with protocol version and parameters.

```go
type StartupMessage struct {
    ProtocolVersion uint32
    Parameters      map[string]string // user, database, etc.
}
```

**Parameters**:
- `user`: Username
- `database`: Database to connect
- `application_name`: Application identifier
- `client_encoding`: Client character encoding

#### Query Message

Sends SQL query to server.

```go
type QueryMessage struct {
    Query string // SQL text
}
```

#### Parse Message

Prepares parameterized query.

```go
type ParseMessage struct {
    Name       string   // Prepared statement name
    Query      string   // SQL with $1, $2, etc.
    ParameterTypes []uint32 // OIDs of parameter types (optional)
}
```

#### Bind Message

Associates parameters with prepared statement.

```go
type BindMessage struct {
    DestPortalName   string        // Portal name
    PreparedName     string        // Prepared statement name
    ParameterFormats []int16       // 0=text, 1=binary
    ParameterValues  [][]byte      // Serialized values
    ResultFormats    []int16       // 0=text, 1=binary
}
```

#### Execute Message

Executes bound portal.

```go
type ExecuteMessage struct {
    PortalName string // Portal name
    MaxRows    int32  // 0=all rows
}
```

#### Describe Message

Requests info about prepared statement or portal.

```go
type DescribeMessage struct {
    Type char        // 'S'=statement, 'P'=portal
    Name string      // Name of statement/portal
}
```

#### Sync Message

Forces backend to flush and close implicit transaction.

```go
type SyncMessage struct {}
```

#### Terminate Message

Closes client connection.

```go
type TerminateMessage struct {}
```

### Server Messages

#### Authentication

Messages related to authentication:

```go
// AuthenticationOk: Connection confirmed
type AuthenticationOk struct {}

// AuthenticationMD5Password: MD5 password requested
type AuthenticationMD5Password struct {
    Salt [4]byte
}

// AuthenticationCleartextPassword: Plaintext password
type AuthenticationCleartextPassword struct {}
```

#### ParameterStatus

Notifies of connection parameters.

```go
type ParameterStatus struct {
    Name  string // Parameter name
    Value string // Parameter value
}
```

**Common Parameters**:
- `server_version`: PostgreSQL version string
- `server_encoding`: Server character encoding
- `database_encoding`: Database character encoding
- `client_encoding`: Client encoding
- `application_name`: Application name

#### ReadyForQuery

Indicates backend ready for new query.

```go
type ReadyForQuery struct {
    Status byte // 'I'=idle, 'T'=transaction, 'E'=error
}
```

#### RowDescription

Describes structure of query result.

```go
type RowDescription struct {
    Fields []FieldDescription
}

type FieldDescription struct {
    Name        string // Column name
    TableOID    uint32 // Table OID (0 if not from table)
    AttrNum     int16  // Attribute number
    TypeOID     uint32 // Data type OID
    TypeLen     int16  // Type size (-1 if variable)
    TypeMod     int32  // Type modifier
    Format      int16  // 0=text, 1=binary
}
```

#### DataRow

Sends one result row.

```go
type DataRow struct {
    Values [][]byte // Serialized column values
}
```

Column values are NULL-terminated if text format, or raw bytes if binary.

#### CommandComplete

Indicates successful command completion.

```go
type CommandComplete struct {
    CommandTag string // "SELECT n", "INSERT 0 n", etc.
}
```

#### Error Response

Sends error information.

```go
type ErrorResponse struct {
    Fields map[byte]string // Error codes and messages
}
```

**Field Types**:
- `'S'`: Severity (ERROR, WARNING, NOTICE, etc.)
- `'C'`: SQLSTATE code (5-char code)
- `'M'`: Human-readable message
- `'D'`: Detail
- `'H'`: Hint
- `'P'`: Position in query

#### Notice Response

Sends non-error messages.

```go
type NoticeResponse struct {
    Fields map[byte]string // Similar to ErrorResponse
}
```

## Wire Format

### Message Structure

```
[MESSAGE_TYPE (1 byte)][LENGTH (4 bytes)][PAYLOAD]
```

**Length**: Total message length including length field (4 bytes), but NOT the message type byte.

### Example Messages

#### Query
```
'Q'          -- Message type
00 00 00 0E  -- Length (14 bytes)
53 45 4C 45  -- "SELECT 1"
43 54 20 31
00           -- Null terminator
```

#### DataRow
```
'D'          -- Message type
00 00 00 0F  -- Length
00 01        -- Number of columns
00 00 00 01  -- Column 1 length (1 byte)
31           -- '1'
```

## Session State Management

### Transaction State

```
Ready
  ↓
Implicit transaction started
  ├→ Query successful  → Implicit commit
  ├→ Query error      → Implicit rollback (later)
  └→ COMMIT/ROLLBACK  → Explicit control
```

**Transaction Modes**:
- Autocommit: Each query commits implicitly
- Transaction: Multi-query transaction

### Session Parameters

Maintained per-connection:

```go
type SessionState struct {
    User            string
    Database        string
    Schema          string
    ClientEncoding  string
    DateStyle       string
    TimeZone        string
    SearchPath      []string
    // Custom parameters
    CustomParams    map[string]string
}
```

## Error Handling

### Error Response Format

```
'E'          -- Error response
[FIELD_TYPE][VALUE]...
'\0'         -- Terminator

Field Types:
'S' → Severity
'C' → SQLSTATE code
'M' → Message
'D' → Detail
```

**SQLSTATE Codes** (PostgreSQL-compatible):
- `00000` - Successful completion
- `3D000` - Invalid database name
- `42P01` - Undefined table
- `42703` - Undefined column
- `42802` - Wrong number of parameters
- `23000` - Integrity constraint violation

## Prepared Statements

### Parse / Bind / Execute Flow

Benefits:
- Query compilation happens once
- Parameters sent separately (security)
- Type information preserved

### Example

```
1. Parse: "SELECT * FROM users WHERE id = $1"
2. Bind: [oid=23, value=123]
3. Execute
4. Result rows
```

## Type Handling

### Type OIDs

Column types identified by PostgreSQL OIDs:

```
23       → INT4
20       → INT8
701      → FLOAT8
25       → TEXT
1114     → TIMESTAMP
1082     → DATE
16       → BOOLEAN
```

### Format Specification

Each column can be text or binary:

- **Text (0)**: ASCII representation
- **Binary (1)**: Native PostgreSQL binary format

## Performance Considerations

### Buffering

- Messages buffered for efficiency
- Flush after each command
- Optional pipelining for multiple commands

### Memory Usage

- Large result sets streamed
- Native binary format for efficiency
- Pool reuse for message objects

## Example: SELECT Query

```
Client → Server
'Q'
00 00 00 0E
SELECT 1\0

Server → Client
'T'  [RowDescription]
  int8 columns
  Column: name="?column?", type=23

'D'  [DataRow]
  value = "1"

'C'  [CommandComplete]
  tag = "SELECT 1"

'Z'  [ReadyForQuery]
  status = 'I'
```

## Extended Query Protocol

### Parse - Analyze Phase

Sends query template without parameters.

**Advantages**:
- Early error detection
- Type resolution
- Prepared statement caching

### Bind - Bind Parameters

Associates actual parameter values.

**Security**:
- Type-checked parameter binding
- SQL injection prevention
- Automatic serialization

### Execute - Run Query

Sends execute request with row limit.

**Streaming**:
- Multiple data rows returned
- Client can fetch more with another Execute
- Efficient for large results

## Streaming Result Protocol

For large result sets:

```
1. Send RowDescription
2. Send DataRow (multiple)
3. Can send suspend with row limit
4. Resume with Execute message
5. Send more DataRow messages
6. Send CommandComplete
```

## Copy Protocol (Future)

For bulk data loading:

```
COPY table FROM STDIN
[each line is a data row]
\.
```

## Notification Messages (Future)

Asynchronous notifications for LISTEN/NOTIFY:

```
'A'  [NotificationResponse]
  channel = "mychannel"
  payload = "data"
```

## Connection Security

### Authentication Methods

1. **Trust**: No password
2. **MD5**: MD5 hash of password
3. **SCRAM-SHA-256**: Salted password hash
4. **Plaintext** (not recommended)

### SSL/TLS

Protocol supports upgrading to SSL:

```
Client → StartupMessage (ssl=1)
Server → 'S' (SSL supported) or 'N' (not supported)
Client → Upgrade to TLS
```

## Related Files

- **Protocol Server**: `internal/runtime/wire/`
- **Message Types**: `internal/runtime/wire/` (message.go, etc.)
- **Connection Handler**: `internal/runtime/wire/server.go`
- **Prepared Statements**: `internal/runtime/wire/prepared.go`

## Testing

Protocol compliance tested against:
- psql (PostgreSQL client)
- pgAdmin
- JDBC drivers
- Python psycopg2
- Go pq driver
