# Functions System

## Overview

The Functions System provides a mechanism for SQL functions and operators. It includes built-in functions (aggregate, scalar, window), user-defined functions, and operator resolution.

## Architecture

### Function Registry

The function system is organized hierarchically:

```
Function Registry
├── Scalar Functions
│   ├── String Functions
│   ├── Numeric Functions
│   ├── Date/Time Functions
│   └── Type Conversion Functions
├── Aggregate Functions
│   ├── COUNT
│   ├── SUM
│   ├── AVG
│   ├── MIN/MAX
│   └── Custom Aggregates
├── Operators
│   ├── Arithmetic (+, -, *, /, %)
│   ├── Comparison (<, >, =, etc.)
│   └── Logical (AND, OR, NOT)
└── User-Defined Functions (Future)
```

## Function Interface

```go
type Function interface {
    Name() string              // Function name
    Arity() int               // Expected number of args (-1 for variadic)
    ReturnType(args []types.Type) types.Type
    Invoke(args []interface{}) (interface{}, error)
}
```

## Built-in Functions

### Scalar Functions

#### String Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `UPPER` | `UPPER(text) → text` | Convert to uppercase |
| `LOWER` | `LOWER(text) → text` | Convert to lowercase |
| `LENGTH` | `LENGTH(text) → int` | String length |
| `SUBSTRING` | `SUBSTRING(text, pos, len) → text` | Extract substring |
| `TRIM` | `TRIM(text) → text` | Remove leading/trailing spaces |
| `LTRIM` | `LTRIM(text) → text` | Remove leading spaces |
| `RTRIM` | `RTRIM(text) → text` | Remove trailing spaces |
| `REPLACE` | `REPLACE(text, from, to) → text` | Replace substring |
| `POSITION` | `POSITION(substr IN text) → int` | Find substring position |
| `CONCAT` | `CONCAT(text, text, ...) → text` | Concatenate strings |
| `STARTS_WITH` | `STARTS_WITH(text, prefix) → bool` | Check prefix |
| `ENDS_WITH` | `ENDS_WITH(text, suffix) → bool` | Check suffix |

#### Numeric Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `ABS` | `ABS(numeric) → numeric` | Absolute value |
| `ROUND` | `ROUND(numeric, digits) → numeric` | Round to digits |
| `FLOOR` | `FLOOR(numeric) → numeric` | Round down |
| `CEIL` | `CEIL(numeric) → numeric` | Round up |
| `SQRT` | `SQRT(numeric) → numeric` | Square root |
| `POWER` | `POWER(numeric, exponent) → numeric` | Exponentiation |
| `MOD` | `MOD(x, y) → numeric` | Modulo |
| `SIGN` | `SIGN(numeric) → int` | Sign (-1, 0, 1) |
| `GREATEST` | `GREATEST(numeric, ...) → numeric` | Maximum value |
| `LEAST` | `LEAST(numeric, ...) → numeric` | Minimum value |

#### Date/Time Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `NOW` | `NOW() → timestamp` | Current timestamp |
| `CURRENT_DATE` | `CURRENT_DATE() → date` | Current date |
| `CURRENT_TIME` | `CURRENT_TIME() → time` | Current time |
| `DATE` | `DATE(timestamp) → date` | Extract date |
| `EXTRACT` | `EXTRACT(part FROM timestamp) → int` | Extract date part |
| `DATE_PART` | `DATE_PART(part, timestamp) → int` | Extract date part |
| `INTERVAL` | `INTERVAL 'value' unit` | Create interval |
| `DATE_TRUNC` | `DATE_TRUNC(unit, timestamp) → timestamp` | Truncate to unit |

#### Type Conversion

| Function | Signature | Description |
|----------|-----------|-------------|
| `CAST` | `CAST(expr AS type)` | Explicit type conversion |
| `TO_TEXT` | `TO_TEXT(any) → text` | Convert to text |
| `TO_INT` | `TO_INT(text) → int` | Parse integer |
| `TO_FLOAT` | `TO_FLOAT(text) → float` | Parse float |
| `TO_BOOL` | `TO_BOOL(text) → bool` | Parse boolean |

### Aggregate Functions

Aggregates combine multiple rows into a single value.

#### Standard Aggregates

```go
func COUNT() → int64
func COUNT(DISTINCT expr) → int64
func SUM(expr) → numeric
func AVG(expr) → numeric
func MIN(expr) → any
func MAX(expr) → any
func STRING_AGG(expr, separator) → text
func ARRAY_AGG(expr) → array
```

#### Aggregate Evaluation

1. **Initialization**: Create empty accumulator
2. **For each input row**:
   - Evaluate expression
   - Update accumulator
3. **Finalization**: Return accumulated value

Example: SUM aggregate

```
Input rows: [1, 2, 3, 4, 5]
Accumulator: 0
  → 0 + 1 = 1
  → 1 + 2 = 3
  → 3 + 3 = 6
  → 6 + 4 = 10
  → 10 + 5 = 15
Output: 15
```

### Operators

SQL operators are implemented as binary and unary functions.

#### Arithmetic Operators

```
a + b  →  Add(a, b)
a - b  →  Subtract(a, b)
a * b  →  Multiply(a, b)
a / b  →  Divide(a, b)
a % b  →  Modulo(a, b)
-a     →  Negate(a)
+a     →  Plus(a)
```

**Type Rules**:
- Integer + Integer → Integer
- Float + Float → Float
- Integer + Float → Float (promotion)
- String + String → String (concatenation with ||)

#### Comparison Operators

```
a = b   →  Equal
a <> b  →  NotEqual
a < b   →  Less
a > b   →  Greater
a <= b  →  LessOrEqual
a >= b  →  GreaterOrEqual
```

**Type Coercion**:
- Numeric: Direct comparison
- String: Lexicographic comparison
- Date: Temporal comparison
- NULL: Always false or UNKNOWN

#### Logical Operators

```
a AND b  →  LogicalAnd(a, b)  [short-circuit]
a OR b   →  LogicalOr(a, b)   [short-circuit]
NOT a    →  LogicalNot(a)
```

**Three-valued Logic**:
- `TRUE AND TRUE` → TRUE
- `TRUE AND FALSE` → FALSE
- `TRUE AND NULL` → NULL
- `FALSE AND NULL` → FALSE (short-circuit)

#### LIKE Operator

```
text LIKE pattern
```

**Pattern Matching**:
- `%` matches any sequence
- `_` matches single character
- `\` escapes special characters

Example:
```
'hello' LIKE 'h%'      → TRUE
'hello' LIKE 'h_llo'   → TRUE
'hello' LIKE '%ll%'    → TRUE
'hello' LIKE 'h\%'     → FALSE
```

#### IN Operator

```
expr IN (value_list)
expr IN (SELECT ...)
```

Equivalent to:
```
expr = value1 OR expr = value2 OR expr = value3 ...
```

#### BETWEEN Operator

```
expr BETWEEN low AND high
```

Equivalent to:
```
expr >= low AND expr <= high
```

#### Comparison with NULL

```sql
NULL = NULL         → NULL (not TRUE)
expr = NULL         → NULL
expr IS NULL        → boolean (TRUE if NULL, FALSE otherwise)
expr IS NOT NULL    → boolean (opposite of IS NULL)
COALESCE(expr, default) → expr if not NULL, else default
```

## Function Resolution

### Resolution Process

When a function call appears in a query:

1. **Name Lookup**: Find all functions with that name
2. **Arity Check**: Match argument count
3. **Type Check**: Match argument types or find coercion
4. **Resolution**: Select best-matching function
5. **Invocation**: Call with actual arguments

### Type Coercion

The system automatically coerces types when needed:

```
Function: CONCAT(text, text) → text
Call: CONCAT(123, 456)

Coercion:
  123 (int) → '123' (text)
  456 (int) → '456' (text)

Result: CONCAT('123', '456') → '123456'
```

### Function Overloading

Multiple functions can have the same name with different signatures:

```go
func SUM(postgres.INT8 array) → INT8
func SUM(postgres.FLOAT8 array) → FLOAT8
func AVG(postgres.INT8 array) → FLOAT8
func AVG(postgres.FLOAT8 array) → FLOAT8
```

## Statistical Aggregate Functions

### Internal State

Aggregate functions maintain internal state:

```go
type AggregateState interface {
    Add(value interface{}) error      // Add value
    Result() interface{}              // Get result
    Copy() AggregateState             // Copy for DISTINCT
    Merge(other AggregateState) error // Merge parallel states
}
```

### Example: COUNT Implementation

```go
type CountState struct {
    count int64
}

func (s *CountState) Add(value interface{}) error {
    if value != nil {
        s.count++
    }
    return nil
}

func (s *CountState) Result() interface{} {
    return s.count
}
```

## User-Defined Functions (Future)

Framework for custom functions:

```sql
CREATE FUNCTION my_func(x INT, y TEXT) RETURNS INT AS
  'SELECT ...'
LANGUAGE SQL;
```

Allows:
- Custom scalar functions
- Custom aggregate functions
- PL/pgSQL or other languages

## Special Functions

### System Functions

```
version()         → PetaCore version
database()        → Current database
schema()          → Current schema
current_user()    → Current user
```

### Meta Functions

```
typeof(expr)      → SQL type name
oid_type(oid)     → Type name from OID
```

## Operator Precedence and Associativity

```
Level  Operators              Associativity
─────────────────────────────────────────
1      (parentheses, functions)  N/A
2      +, - (unary)           Right
3      *, /, %                Left
4      +, - (binary)          Left  
5      ||                     Left
6      =, <, >, <>, <=, >=, IN, LIKE, BETWEEN  N/A
7      IS NULL, IS NOT NULL   N/A
8      NOT                    Right
9      AND                    Left
10     OR                     Left
```

## Performance Optimization

### Function Dispatch

- **Static dispatch**: Function resolved at plan time
- **Dynamic dispatch**: Function resolved at runtime

### Inlining

Simple functions inlined at execution time:
- ABS(x)
- SIGN(x)
- UPPER(x)

### Vectorization (Future)

Process vectors of values:

```
SUM([1, 2, 3, 4, 5]) → 15  [single vector operation]
```

## Error Handling

Functions handle errors gracefully:

```
SELECT 1 / 0         → ERROR: division by zero
SELECT CAST('abc' AS INT) → ERROR: cannot convert text to integer  
SELECT SQRT(-1)      → NULL or ERROR: domain error
```

## Configuration

### Function Registration

Functions are registered with their metadata:

```go
registry.Register(FunctionSpec{
    Name: "UPPER",
    Arity: 1,
    ReturnType: func(args []Type) Type { return TEXT },
    Invoke: func(args []interface{}) (interface{}, error) {
        text, _ := args[0].(string)
        return strings.ToUpper(text), nil
    },
})
```

## Related Files

- **Main**: `internal/runtime/functions/function.go`
- **Operators**: `internal/runtime/rhelpers/parser.go`
- **Expression Evaluation**: `internal/runtime/rhelpers/revaluate/`
- **Runtime Models**: `internal/runtime/rhelpers/rmodels/`

## Testing

Functions come with comprehensive tests:

- Unit tests per function
- Integration tests with queries
- Type coercion tests
- Edge case tests (NULL, overflow, etc.)

## Future Extensions

1. **Window Functions**: OVER clause support
2. **CTEs**: Common table expressions
3. **Recursive Functions**: Self-referential definitions
4. **JSON Functions**: JSON path operations
5. **More Aggregate Functions**: PERCENTILE, MEDIAN, STDDEV
6. **Analytics Functions**: LAG, LEAD, RANK
