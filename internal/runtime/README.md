# PetaCore SQL Runtime

The PetaCore runtime is a production-grade SQL query engine implementing PostgreSQL compatibility with Snapshot Isolation (MVCC) concurrency control. It provides a clean separation between parsing, planning, and execution phases.

## Quick Overview

```
SQL Input → Parser → Planner → Executor → Results
```

The runtime consists of several integrated modules that together implement a complete SQL database engine.

## Core Modules

### 📋 [Parser](./docs/Parser.md)
Converts SQL text into an Abstract Syntax Tree (AST).

**Location**: `parser/`  
**Generation**: ANTLR4 from `sql.g4` grammar  
**Responsibility**: Lexical and syntactic analysis

### 🎯 [Query Planner](./planner/)
Transforms parsed SQL into a tree of logical execution operators.

**Location**: `planner/`  
**Key Functions**: `CreateQueryPlan()`, `ExecutePlan()`  
**Responsibility**: Semantic analysis and plan generation

### ⚡ [Executor](./docs/Executor.md)
Executes logical plans against the storage engine.

**Location**: `planner/executor.go`, `executor/`  
**Responsibility**: Tuple processing and result generation

### 🔧 [Semantic Layer](./rsql/)
Intermediate representation providing structure for SQL constructs.

**Location**: `rsql/`  
**Components**:
- `statements/` - SQL statements
- `items/` - Expressions and clauses
- `table/` - Result sets
- `visitor/` - AST to semantic conversion

### 🎨 [Functions System](./docs/Functions.md)
SQL functions, operators, and type coercion.

**Location**: `functions/`  
**Features**:
- Scalar functions (string, numeric, date/time)
- Aggregate functions (COUNT, SUM, AVG, etc.)
- Type casting and coercion

### 📦 [Type System](./docs/TypeSystem.md)
PostgreSQL-compatible SQL types.

**Location**: Referenced through SDK  
**Supported Types**:
- Numeric: SMALLINT, INTEGER, BIGINT, FLOAT, DOUBLE, NUMERIC
- String: TEXT, VARCHAR, CHAR
- Temporal: DATE, TIME, TIMESTAMP
- Special: BOOLEAN, BYTEA, ARRAY

### 🔌 [Wire Protocol](./docs/Wire.md)
PostgreSQL Frontend/Backend Protocol v3.0 implementation.

**Location**: `wire/`  
**Purpose**: Enable compatibility with PostgreSQL clients

### 🗂️ [System Tables](./system/)
PostgreSQL-compatible system catalog implementation.

**Tables**: pg_class, pg_attribute, pg_namespace, pg_type, etc.

### 🧦 [Runtime Helpers](./rhelpers/)
Utilities for expression evaluation and type conversion.

**Location**: `rhelpers/`  
**Components**:
- `rmodels/` - Runtime expression models
- `subquery/` - Subquery execution
- `revaluate/` - Expression evaluation
- `rops/` - Runtime operations

## Architecture Overview

See [Architecture.md](./docs/Architecture.md) for detailed system design.

### Query Execution Pipeline

```
1. SQL Text
   ↓
2. Lexer/Parser → Abstract Syntax Tree
   ↓
3. Visitor Pattern → Semantic Layer
   ↓
4. Planner → Logical Plan (tree of operators)
   ↓
5. Executor → Result Set
   ↓
6. Wire Protocol → Client Response
```

### Plan Node Types

The planner generates a tree of logical operators:

- **ScanPlanNode** - Table scan
- **ProjectPlanNode** - SELECT list projection
- **FilterPlanNode** - WHERE filtering
- **JoinPlanNode** - Table joins
- **AggregatePlanNode** - GROUP BY aggregation
- **SortPlanNode** - ORDER BY sorting
- **LimitPlanNode** - LIMIT/OFFSET
- **UnionPlanNode** - UNION operation
- **IntersectPlanNode** - INTERSECT operation
- **ExceptPlanNode** - EXCEPT (set difference)
- **SubqueryPlanNode** - Subqueries
- **ValuesPlanNode** - SELECT without FROM

## Key Features

### SQL Support
- ✅ SELECT queries with all clauses
- ✅ JOINs (INNER, LEFT, RIGHT, FULL, CROSS)
- ✅ Set operations (UNION, INTERSECT, EXCEPT)
- ✅ Subqueries in FROM and WHERE
- ✅ Aggregation with GROUP BY/HAVING
- ✅ Expression evaluation
- ✅ SQL functions and operators
- ✅ Type casting and coercion

### Data Definition Language
- ✅ CREATE TABLE
- ✅ DROP TABLE
- ✅ TRUNCATE TABLE
- ✅ System table querying

### Concurrency & Isolation
- ✅ MVCC (Multi-Version Concurrency Control)
- ✅ Snapshot Isolation
- ✅ Automatic transaction management
- ✅ Conflict-free read operations

### PostgreSQL Compatibility
- ✅ Wire protocol v3.0
- ✅ Prepared statements
- ✅ System catalog tables
- ✅ Type system (OIDs)
- ✅ Standard SQL functions

## Documentation Structure

```
docs/
├── Architecture.md      # System design overview
├── Planner.md          # Query planner details
├── Executor.md         # Plan execution engine
├── Parser.md           # SQL parsing
├── Functions.md        # SQL functions
├── TypeSystem.md       # Type handling
└── Wire.md             # Protocol implementation

planner/
├── README.md           # Planner module guide
├── planner.go          # Main planner
├── plan.go             # Plan node definitions
├── executor.go         # Plan execution
└── set_operations.go   # Set operations
```

## Usage Example

### Complete Query Execution

```go
package main

import (
    "petacore/internal/runtime/planner"
    "petacore/internal/runtime/rsql/visitor"
)

func main() {
    // 1. Parse SQL
    sql := "SELECT name, age FROM users WHERE age > 18"
    stmt, err := visitor.ParseSQL(sql)
    if err != nil {
        panic(err)
    }

    // 2. Create plan
    planCtx := planner.PlannerContext{
        Database: "testdb",
        Schema:   "public",
    }
    plan, err := planner.CreateQueryPlan(
        stmt.(*statements.SelectStatement),
        planCtx,
    )
    if err != nil {
        panic(err)
    }

    // 3. Execute plan
    execCtx := planner.ExecutorContext{
        Database:  "testdb",
        Schema:    "public",
        Storage:   storageEngine,
        Allocator: memAllocator,
    }
    result, err := planner.ExecutePlan(plan, execCtx)
    if err != nil {
        panic(err)
    }

    // 4. Process results
    for _, row := range result.Rows {
        // Handle row data
    }
}
```

## Testing

Run all runtime tests:

```bash
go test ./internal/runtime/...
```

Run specific module tests:

```bash
go test ./internal/runtime/planner/...    # Planner tests
go test ./internal/runtime/executor/...   # Executor tests
```

## Design Principles

1. **Separation of Concerns**: Parser, planner, executor are distinct
2. **Tree-Based Execution**: Plans and results are tree structures
3. **Composability**: Nest operators for complex queries
4. **Extensibility**: Add new operators without modifying existing code
5. **Type Safety**: Leverage Go's type system
6. **Simplicity**: Clear, maintainable code over optimization

## Performance Characteristics

- **Parsing**: O(n) where n = SQL length
- **Planning**: O(n) where n = AST size
- **Execution**: Depends on data size and operations
- **Memory**: Buffering determined by operator (streaming where possible)

## Transaction Model

- **Isolation Level**: Snapshot Isolation
- **Concurrency Control**: MVCC
- **Consistency**: ACID (with storage layer support)
- **Transactions**: Automatic or explicit control

## Storage Integration

Abstracts storage through clean interface:

```go
interface DistributedStorageVClock {
    Read()          // Tuple retrieval
    Write()         // Tuple insertion
    Update()        // Tuple modification
    Delete()        // Tuple deletion
    CreateTable()   // Schema creation
    DropTable()     // Schema removal
}
```

## Extension Points

### Add New Operators

1. Define `PlanNode` in `planner/plan.go`
2. Add builder in `planner/planner.go`
3. Implement executor in `planner/executor.go`
4. Add tests

### Add Functions

1. Register in `functions/function.go`
2. Implement evaluation logic
3. Add type resolution
4. Test with queries

### Add SQL Features

1. Update `sql.g4` grammar
2. Extend visitor in `rsql/visitor/`
3. Update planner
4. Add executor support

## Known Limitations & Future Work

### Current Limitations
- No optimizer (uses default plan)
- Single-threaded execution
- Limited aggregate functions
- No window functions
- No CTEs (WITH clause)

## Contributing

When extending the runtime:

1. Follow existing code style
2. Add comprehensive tests
3. Update relevant documentation
4. Ensure type safety
5. Consider performance implications

## Files Organization

```
internal/runtime/
├── parser/              # ANTLR4 generated parser
├── planner/             # Query planner
│   ├── README.md        # Planner guide
│   ├── planner.go       # Main planner
│   ├── plan.go          # Node definitions
│   ├── executor.go      # Executor
│   └── set_operations.go
├── executor/            # DDL executor
├── functions/           # SQL functions
├── rsql/                # Semantic layer
├── rhelpers/            # Runtime helpers
├── wire/                # PostgreSQL protocol
├── system/              # System tables
├── docs/                # Documentation
└── README.md            # This file
```

## Related Components

- **Storage**: `internal/storage/` - Data persistence with MVCC
- **SDK**: `sdk/` - Type system and serializers
- **Logger**: `internal/logger/` - Logging infrastructure

## References & Resources

### Project Structure
- See [Architecture.md](./docs/Architecture.md) for detailed design

### Module Details
- [Planner Documentation](./docs/Planner.md)
- [Executor Documentation](./docs/Executor.md)
- [Parser Documentation](./docs/Parser.md)
- [Functions Documentation](./docs/Functions.md)
- [Type System Documentation](./docs/TypeSystem.md)
- [Wire Protocol Documentation](./docs/Wire.md)

### Inspiration
- DuckDB - Modular in-process OLAP database
- PostgreSQL - Mature SQL database
- Calcite - Query optimization framework

## License

Same as PetaCore main project
