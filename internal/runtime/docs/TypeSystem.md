# Type System

## Overview

PetaCore implements a SQL type system compatible with PostgreSQL, providing a rich set of scalar and collection types. The type system handles automatic type coercion, comparison, and serialization.

## SQL Types

### Numeric Types

#### Integer Types

| Type | Size | Range | Example |
|------|------|-------|---------|
| `SMALLINT`/`INT2` | 2 bytes | -32,768 to 32,767 | `123` |
| `INTEGER`/`INT`/`INT4` | 4 bytes | -2.1B to 2.1B | `123456` |
| `BIGINT`/`INT8` | 8 bytes | -9.2E18 to 9.2E18 | `123456789` |

#### Floating Point Types

| Type | Size | Precision | Example |
|------|------|-----------|---------|
| `REAL`/`FLOAT4` | 4 bytes | ~6 digits | `3.14::float4` |
| `DOUBLE PRECISION`/`FLOAT8`/`DOUBLE` | 8 bytes | ~15 digits | `3.141592654` |

#### Exact Numeric Type

| Type | Precision | Scale | Example |
|------|-----------|-------|---------|
| `NUMERIC`/`DECIMAL` | Configurable | Configurable | `123.45` |

### String Types

| Type | Characteristics | Example |
|------|-----------------|---------|
| `VARCHAR(n)` | Variable length, max n | `'hello'` |
| `CHAR(n)` | Fixed length, padded | `'hello   '` |
| `TEXT` | Variable length, unlimited | `'hello'` |

**String Literals**:
- Single quotes: `'hello'`
- Escape single quote: `'it''s'`
- E-string syntax: `E'hello\nworld'`

### Date/Time Types

| Type | Range | Format | Example |
|------|-------|--------|---------|
| `DATE` | 4713 BC - 5874897 AD | YYYY-MM-DD | `2024-03-02` |
| `TIME` | 00:00:00 - 23:59:59 | HH:MM:SS | `14:30:00` |
| `TIMESTAMP` | Like DATE + TIME | YYYY-MM-DD HH:MM:SS | `2024-03-02 14:30:00` |
| `TIMESTAMP WITH TIME ZONE` | Like TIMESTAMP + TZ offset | With UTC offset | `2024-03-02 14:30:00+01:00` |
| `INTERVAL` | Duration | Signed time interval | `'1 day 2 hours'` |

### Boolean Type

| Type | Values | Example |
|------|--------|---------|
| `BOOLEAN`/`BOOL` | TRUE, FALSE, NULL | `TRUE`, `FALSE` |

### Binary Type

| Type | Usage | Example |
|------|-------|---------|
| `BYTEA` | Binary data | `'\x012345'` |

### Collection Types

| Type | Description | Example |
|------|-------------|---------|
| `ARRAY` | Ordered collection | `ARRAY[1, 2, 3]` |

### Special Types

| Type | Purpose | Example |
|------|---------|---------|
| `NULL` | Unknown/missing value | `NULL` |
| `VOID` | No value | Return type |

## Type Operators and Control Structures

### Type Casting

#### Implicit Casting

Automatic type promotion in operations:

```sql
SELECT 1 + 2.5        -- 1 (INT) promoted to FLOAT → 3.5
SELECT 'hello' || 123 -- 123 converted to TEXT → 'hello123'
```

**Promotion Order**:
```
SMALLINT → INTEGER → BIGINT → NUMERIC → FLOAT → DOUBLE
```

#### Explicit Casting

Using CAST operator:

```sql
CAST(value AS type)
-- Example:
CAST('123' AS INTEGER)      -- Parse string to int
CAST(123.45 AS INTEGER)     -- Truncate
CAST(123 AS VARCHAR)        -- Stringify
```

### Type Comparison

Comparison rules:

1. **Same Type**: Direct comparison
2. **Numeric Types**: Value comparison with promotion
3. **String Types**: Lexicographic comparison
4. **Date/Time**: Temporal comparison
5. **Different Categories**: Type error or conversion

### NULL Handling

Three-valued logic:

```
NULL = NULL           → NULL
value = NULL          → NULL
NULL IS NULL          → TRUE
value IS NOT NULL     → TRUE (if value != NULL)
COALESCE(a, b, c)     → First non-NULL value
```

## Type Categories

### Numeric Category

- SMALLINT, INTEGER, BIGINT
- NUMERIC, DECIMAL
- REAL, DOUBLE PRECISION

Operations:
- Arithmetic: `+`, `-`, `*`, `/`, `%`
- Comparison: `<`, `>`, `=`, etc.
- Aggregate: `SUM`, `AVG`, `MIN`, `MAX`

### String Category

- TEXT, VARCHAR, CHAR

Operations:
- Concatenation: `||`
- Pattern matching: `LIKE`
- String functions: `UPPER`, `LOWER`, `LENGTH`, etc.

### Temporal Category

- DATE, TIME, TIMESTAMP, TIMESTAMP WITH TIME ZONE

Operations:
- Comparison: `<`, `>`, `=`, etc.
- Arithmetic: `date + interval`, etc.
- Extraction: `EXTRACT`, `DATE_PART`

### Boolean Category

Operations:
- Logical: `AND`, `OR`, `NOT`
- Comparison: `=`, `<>`

## Type Conversions

### String to Numeric

```
'123' → 123
'3.14' → 3.14
'123abc' → Error (invalid format)
```

### Numeric to String

```
123 → '123'
3.14 → '3.14'
```

### Date/Time Parsing

```
'2024-03-02' (DATE format) → DATE
'14:30:00' (TIME format) → TIME
'2024-03-02 14:30:00' → TIMESTAMP
```

### Boolean Conversion

```
TRUE → 1 (as int)
FALSE → 0 (as int)
'true', 't', 'yes', 'y', '1' → TRUE
'false', 'f', 'no', 'n', '0' → FALSE
```

## OID (Object Identifier)

PostgreSQL compatibility - types identified by OID:

```
INT2 = 21
INT4 = 23
INT8 = 20
FLOAT4 = 700
FLOAT8 = 701
TEXT = 25
VARCHAR = 1043
BOOLEAN = 16
TIMESTAMP = 1114
DATE = 1082
TIME = 1083
BYTEA = 17
```

## Type Checking

The type system validates:

1. **Function Arguments**: Type compatibility
2. **Operator Operands**: Valid combinations
3. **Assignment**: Target type compatibility
4. **Collation**: String comparison rules

## Serialization and Deserialization

Types are serialized for storage and network transmission:

### Binary Format

Each type has:
- **Serializer**: value → bytes
- **Deserializer**: bytes → value

Example: INT4
```
123 → [00 00 00 7B] (big-endian)
```

### Protocol

Each value includes:
- OID (4 bytes) - Type identifier
- Length (4 bytes) - Data length
- Data - Serialized value

## Null Semantics

### NULL Propagation

Operations with NULL produce NULL:

```
SELECT 1 + NULL         → NULL
SELECT NULL || 'text'   → NULL
SELECT UPPER(NULL)      → NULL
SELECT COALESCE(NULL, 'default')  → 'default'
```

### Exception: IS NULL

```
SELECT NULL IS NULL     → TRUE
SELECT 1 IS NULL        → FALSE
SELECT NULL IS NOT NULL → FALSE
```

## Type Registry

Internal type configuration:

```go
type Type struct {
    OID        uint32        // PostgreSQL OID
    Name       string        // Type name
    Category   TypeCategory  // Numeric, String, etc.
    Size       int32         // Fixed size (-1 if variable)
    Serializer Serializer    // Conversion to bytes
    Comparable bool          // Can be compared
    Orderable  bool          // Can be ordered
}
```

## Array Types

### Array Literals

```sql
ARRAY[1, 2, 3]
ARRAY['a', 'b', 'c']
ARRAY[1, 2, 3]::INT4[]
```

### Array Operations

```
arr[1]        -- Access element (1-indexed)
arr[1:3]      -- Slice
array_length(arr)
array_append(arr, value)
array_concat(arr1, arr2)
```

## Future Type Support

1. **JSON/JSONB**: JSON data type
2. **UUID**: Unique identifiers
3. **ENUM**: User-defined enumerations
4. **RANGE**: Range types (INT4RANGE, DATERANGE, etc.)
5. **Composite**: User-defined types
6. **CITEXT**: Case-insensitive text

## Type Inference

For expressions without explicit types:

1. If all operands have same type → use that type
2. If numeric types mixed → promote to highest
3. If string and numeric → convert numeric to string
4. Otherwise → error

## Performance

- **Type checks**: O(1) lookup from OID
- **Serialization**: ~O(data size)
- **Comparison**: Type-specific, usually O(data size)

## Related Files

- **Type Definitions**: `sdk/types/`
- **Serializers**: `sdk/serializers/`
- **Type Coercion**: `internal/runtime/functions/`
- **Operators**: `internal/runtime/rhelpers/`
