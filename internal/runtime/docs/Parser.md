# SQL Parser

## Overview

The SQL Parser is responsible for converting SQL text into an Abstract Syntax Tree (AST) that can be processed by the semantic layer and query planner. It is generated from an ANTLR4 grammar and provides both lexical analysis and syntactic parsing.

## Architecture

### Parser Generation

The parser is generated from the SQL grammar defined in `sql.g4`:

```
sql.g4 (ANTLR4 Grammar)
    ↓
antlr4 compiler
    ↓
sqlLexer.go    (Lexical analysis)
sqlParser.go   (Syntactic analysis)
sql_listener.go (Visitor interface)
```

## Core Components

### Lexer (sqlLexer.go)

**Responsibility**: Convert SQL text into a stream of tokens.

**Tokens Generated**:
- Keywords: SELECT, FROM, WHERE, JOIN, etc.
- Identifiers: Table names, column names
- Literals: Numbers, strings
- Operators: +, -, *, /, <, >, =, etc.
- Punctuation: (, ), ;, ,

**Token Stream Example**:
```
SELECT name, age FROM users WHERE age > 18
    ↓
[SELECT, IDENTIFIER(name), COMMA, IDENTIFIER(age), FROM, IDENTIFIER(users), 
 WHERE, IDENTIFIER(age), GT, NUMBER(18)]
```

### Parser (sqlParser.go)

**Responsibility**: Build Abstract Syntax Tree from token stream.

**Parsing Rules** (from sql.g4):
- `statement` → SELECT, INSERT, UPDATE, DELETE, CREATE TABLE, DROP TABLE, etc.
- `selectStatement` → SELECT clause + FROM + WHERE + GROUP BY + ORDER BY + LIMIT
- `expression` → Arithmetic, comparison, logical, function calls
- `functionCall` → Function name + argument list
- `joinClause` → JOIN type + table + ON condition

### Listener Pattern (sql_listener.go)

The parser uses a listener pattern to traverse the AST:

```
Parser
  ↓
AST
  ↓
TreeListener
  ↓
Visitor Events (enterContext, exitContext)
```

## Parsing Flow

### Parse Function

```go
func ParseSQL(sql string) (interface{}, error)
```

**Steps**:
1. Create input stream from SQL string
2. Create lexer from input stream
3. Create token stream from lexer
4. Create parser from token stream
5. Call `statement()` rule on parser
6. Return root ParseTree node

### Error Handling

Parser errors include:
- **Lexical Errors**: Unrecognized tokens
- **Syntactic Errors**: Invalid grammar structure
- **Recovery**: Attempts to continue parsing

**Error Details**:
- Line number
- Column number  
- Expected tokens
- Actual token

## Statement Types

The parser supports all major SQL statements:

### Query Statements

**SELECT**:
```sql
SELECT [DISTINCT] select_list 
FROM table_sources
[WHERE condition]
[GROUP BY group_by_list]
[HAVING condition]
[ORDER BY order_by_list]
[LIMIT count [OFFSET offset]]
```

**UNION/INTERSECT/EXCEPT**:
```sql
select_query {UNION | INTERSECT | EXCEPT} [ALL] select_query
```

### Data Modification Statements

**INSERT**:
```sql
INSERT INTO table (column_list) VALUES (value_list)
```

**UPDATE**:
```sql
UPDATE table SET column = value [WHERE condition]
```

**DELETE**:
```sql
DELETE FROM table [WHERE condition]
```

### Schema Statements

**CREATE TABLE**:
```sql
CREATE TABLE table (
    column_name column_type [constraints],
    ...
)
```

**DROP TABLE**:
```sql
DROP TABLE [IF EXISTS] table_name
```

**TRUNCATE TABLE**:
```sql
TRUNCATE TABLE table_name
```

## Expression Parsing

### Expression Types

1. **Literals**:
   - Integer: `123`, `-456`
   - Float: `3.14`, `1e-3`
   - String: `'hello'`, `'it''s'` (escaping)
   - Boolean: `TRUE`, `FALSE`
   - NULL: `NULL`

2. **Identifiers**:
   - Simple: `column_name`
   - Qualified: `table_name.column_name`
   - Quoted: `"Column Name"`

3. **Functions**:
   ```sql
   function_name(arg1, arg2, ...)
   aggregate_function(*) OVER (PARTITION BY ... ORDER BY ...)
   ```

4. **Operators**:
   - Arithmetic: `+`, `-`, `*`, `/`, `%`
   - Comparison: `=`, `<>`, `<`, `>`, `<=`, `>=`
   - Logical: `AND`, `OR`, `NOT`
   - String: `||` (concatenation)
   - LIKE, BETWEEN, IN, EXISTS

5. **Complex Expressions**:
   - Subqueries: `(SELECT ...)`
   - CASE: `CASE WHEN ... THEN ... ELSE ... END`
   - Casts: `CAST(expr AS type)`
   - Aggregates with filters: `SUM(x) FILTER (WHERE ...)`

### Operator Precedence

Parsed in standard SQL precedence order:
1. Function calls, parentheses
2. Unary +, -
3. *, /, %
4. Binary +, -
5. Comparison operators
6. NOT
7. AND
8. OR

### Type Parsing

Supported SQL types:
- Numeric: BIGINT, INT, SMALLINT, FLOAT, DOUBLE, DECIMAL/NUMERIC
- String: VARCHAR, TEXT, CHAR
- Temporal: DATE, TIME, TIMESTAMP, TIME WITH TIME ZONE
- Binary: BYTEA
- Boolean: BOOLEAN
- Collections: ARRAY

## Parse Tree Structure

The parser builds a hierarchical tree structure:

```
StatementContext
├── SelectStatementContext
│   ├── SelectSpecContext
│   │   ├── SelectItemContext "name"
│   │   └── SelectItemContext "age"
│   ├── FromClauseContext
│   │   └── TableSourceContext
│   │       └── TableNameContext "users"
│   └── WhereClauseContext
│       └── ExpressionContext
│           ├── ComparisonExprContext
│           ├── ColumnRefContext "age"
│           ├── Operator ">"
│           └── LiteralContext "18"
```

## Visitor Integration

The semantic layer uses a visitor pattern to convert parse trees to semantic structures:

### Visitor Pattern

```go
type Visitor interface {
    VisitStatement(ctx *StatementContext) interface{}
    VisitSelectStatement(ctx *SelectStatementContext) interface{}
    VisitTableSource(ctx *TableSourceContext) interface{}
    // ... more visit methods
}
```

### Conversion Flow

```
ParseTree
    ↓
Visitor.visit(tree)
    ↓
Semantic Objects (statements.SelectStatement)
    ↓
Planner
```

## Semantic Actions

The parser doesn't perform semantic analysis; instead:
1. Parser builds syntactically correct tree
2. Visitor converts to semantic objects
3. Semantic layer validates names, types, etc.

This separation allows:
- Clean parser/semantic boundary
- Reusable parse trees
- Better error reporting

## Grammar Highlights

### Key Grammar Rules (Simplified)

```antlr
statement
    : selectStatement
    | insertStatement
    | updateStatement
    | deleteStatement
    | createTableStatement
    | dropTableStatement
    ;

selectStatement
    : selectPrimary (setOp selectPrimary)*
    ;

selectPrimary
    : SELECT [DISTINCT] selectItemList
      FROM tableSourceList
      [WHERE expression]
      [GROUP BY expressionList]
      [HAVING expression]
      [ORDER BY orderByList]
      [LIMIT number [OFFSET number]]
    ;

expression
    : expression op expression
    | '(' expression ')'
    | functionCall
    | columnRef
    | literal
    | caseExpression
    ;
```

## Error Recovery

The parser attempts to recover from errors and continue:

1. **Synchronization**: Skip tokens until recovery point
2. **Suggested Fixes**: Missing semicolon, mismatched parentheses
3. **Error Collection**: Report multiple errors instead of failing on first

Example:
```
SELECT FROM users  -- Missing select list
    ↓
ERROR: Missing select item list
[Recovery: Skip to FROM, continue parsing]
    ↓
Partial tree with error nodes
```

## Performance Characteristics

- **Linear in SQL size**: O(n) where n = SQL length
- **Memory**: O(tree size) for AST
- **No optimization**: Pure syntax analysis

## Parse Tree Caching

For repeated queries:
1. Parse once
2. Cache parse tree
3. Reuse for different parameters

Requires type-stable parameters.

## Related Components

- **Grammar**: `sql.g4`
- **Generated Lexer**: `sqlLexer.go`
- **Generated Parser**: `sqlParser.go`
- **Listener Base**: `sql_base_listener.go`
- **Visitor**: `rsql/visitor/`
- **Statements**: `rsql/statements/`

## DuckDB Comparison

### Similarities

- ANTLR4 grammar for parser generation
- Visitor pattern for conversion
- Separation of parsing and semantic analysis
- Recursive expression evaluation

### Differences

- DuckDB has more optimization passes
- DuckDB has more comprehensive type system
- DuckDB has window function support (can be added)
- PetaCore focuses on simplicity and clarity

## Future Enhancements

1. **Parse Tree Caching**: Cache for repeated queries
2. **Better Error Messages**: Context-aware suggestions
3. **Extended SQL Support**:
   - Window functions
   - CTEs (WITH clause)
   - MORE complex joins
4. **Parameter Support**: Named and positional parameters
5. **Parser Hints**: Optimization guidance in comments
