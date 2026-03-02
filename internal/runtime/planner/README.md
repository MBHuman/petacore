# Query Planner Module

This directory contains the PetaCore Query Planner - responsible for transforming SQL SELECT statements into optimized logical execution plans.

## Overview

The planner implements a two-phase architecture:

1. **Planner Phase** - Converts parsed SQL (AST) into a tree of logical plan nodes
2. **Executor Phase** - Executes the plan tree against the storage engine

This separation provides modularity, testability, and extensibility comparable to modern database systems like DuckDB.

## Module Contents

### Core Files

- **`planner.go`** - Main planner that builds logical plans from SELECT statements
- **`plan.go`** - Definition of all logical plan node types  
- **`executor.go`** - Executor engine that traverses and executes plan trees
- **`set_operations.go`** - Handles UNION, INTERSECT, EXCEPT operations

## Quick Start

### Creating and Executing a Query Plan

```go
package main

import (
    "petacore/internal/runtime/planner"
    "petacore/internal/runtime/rsql/visitor"
)

// 1. Parse SQL to AST
stmt, err := visitor.ParseSQL("SELECT a, b FROM table1 WHERE a > 10")
if err != nil {
    panic(err)
}

// 2. Create logical plan
plannerCtx := planner.PlannerContext{
    Database: "mydb",
    Schema:   "public",
}
plan, err := planner.CreateQueryPlan(
    stmt.(*statements.SelectStatement), 
    plannerCtx,
)
if err != nil {
    panic(err)
}

// 3. Execute plan
executorCtx := planner.ExecutorContext{
    Database: "mydb",
    Schema:   "public",
    Storage:  storageEngine,
    Allocator: memoryAllocator,
}
result, err := planner.ExecutePlan(plan, executorCtx)
if err != nil {
    panic(err)
}

// Use result.Rows and result.Schema
```

## Architecture

### Planning Process

```
SQL Text
  ↓
Parser → AST
  ↓
Visitor → Semantic Objects
  ↓
Planner → LogicalPlan (tree of PlanNodes)
```

### Plan Node Types

The planner generates a tree of logical operators:

| Node Type | Purpose | Semantics |
|-----------|---------|-----------|
| **ScanPlanNode** | Read table | Emit all tuples from table |
| **ProjectPlanNode** | Column selection | Apply SELECT list to each row |
| **FilterPlanNode** | Row filtering | Evaluate WHERE condition |
| **JoinPlanNode** | Table joining | Combine rows from two sources |
| **AggregatePlanNode** | Grouping & aggregation | GROUP BY + aggregate functions |
| **SortPlanNode** | Sorting | ORDER BY + sort direction |
| **LimitPlanNode** | Row limiting | LIMIT + OFFSET |
| **UnionPlanNode** | Set union | Combine sets of rows |
| **IntersectPlanNode** | Set intersection | Rows in both sets |
| **ExceptPlanNode** | Set difference | Rows in left but not right |
| **ValuesPlanNode** | Constant rows | SELECT without FROM |
| **SubqueryPlanNode** | Subqueries | Treat subquery as table |

### Example Plans

**Simple Query**:
```sql
SELECT a, b FROM t1 WHERE x > 10 ORDER BY a LIMIT 5
```

Plan tree:
```
Limit(5)
  │
  └─ Sort(ORDER BY a)
      │
      └─ Filter(WHERE x > 10)
          │
          └─ Project(SELECT a, b)
              │
              └─ Scan(t1)
```

**Union Query**:
```sql
SELECT a FROM t1 UNION SELECT b FROM t2
```

Plan tree:
```
Union(ALL=false)
  │
  ├─ Project(SELECT a)
  │   │
  │   └─ Scan(t1)
  │
  └─ Project(SELECT b)
      │
      └─ Scan(t2)
```

## Key Features

### Modular Design
Each SQL operation is a distinct plan node type, allowing:
- Independent unit testing
- Code reuse
- Easy extension with new operators

### Query Compatibility
Supports:
- Simple SELECT queries
- Complex JOINs (INNER, LEFT, RIGHT, FULL, CROSS)
- Set operations (UNION, INTERSECT, EXCEPT)
- Subqueries in FROM and WHERE clauses
- Aggregation with GROUP BY
- Sorting and limiting
- Expressions and functions

### Performance Considerations
- Short-circuit evaluation for boolean expressions
- Lazy subquery evaluation
- Stream processing where possible
- Memory-efficient MVCC transaction handling

## Integration Points

### With Parser
- Receives `statements.SelectStatement` from SQL parser
- Parser performs syntactic analysis
- Planner performs semantic analysis

### With Executor  
- Generates execution plans for `ExecutePlan()`
- Plan nodes define execution interface
- Executor implements per-node strategies

### With Storage
- Uses `DistributedStorageVClock` for data access
- Works within transaction boundaries
- Respects MVCC snapshot isolation

## Design Advantages

1. **Correctness**: Clear separation of parsing, planning, optimization
2. **Debuggability**: Print plan for analysis, trace execution per node
3. **Testability**: Unit test each node type independently
4. **Extensibility**: Add new operators without touching existing code
5. **Clarity**: Business logic explicit in plan structure

## Future Enhancements

### Optimization Passes
- Cost-based plan selection
- Join ordering optimization
- Predicate push-down
- Constant expression folding

### Extended SQL Support
- Common Table Expressions (CTEs / WITH clause)
- Window functions (OVER clause)
- Advanced aggregation (GROUPING SETS, CUBE, ROLLUP)

### Performance Features
- Parallel plan execution
- Adaptive plan selection
- Incremental result streaming

## Testing

Run tests with:

```bash
go test ./internal/runtime/planner/...
```

Test files cover:
- Plan node creation
- Expression evaluation
- Join execution
- Set operations
- Edge cases (NULL handling, type coercion)

## Related Documentation

See the `docs/` directory for detailed documentation:

- **[Architecture.md](../docs/Architecture.md)** - Overall runtime architecture
- **[Planner.md](../docs/Planner.md)** - Detailed planner design
- **[Executor.md](../docs/Executor.md)** - Plan execution details
- **[TypeSystem.md](../docs/TypeSystem.md)** - SQL type handling
- **[Functions.md](../docs/Functions.md)** - SQL functions

## References

The planner design is inspired by modern database systems:
- DuckDB's modular execution engine
- PostgreSQL's planning and execution separation
- Calcite's logical plan representation

- [ ] Реализовать полноценную агрегацию в `processGroupByAndAggregates`
- [ ] Добавить поддержку ON условий в JOIN
- [ ] Реализовать RIGHT и FULL JOIN
- [ ] Добавить оптимизатор планов
- [ ] Добавить статистику выполнения
- [ ] Добавить EXPLAIN для просмотра планов
