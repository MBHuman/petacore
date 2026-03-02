# Query Planner Architecture

## Overview

The Query Planner is responsible for converting SQL SELECT statements into a tree-based logical execution plan. It performs semantic analysis, query optimization, and generates a sequence of operations that implement the query.

## Key Components

### Planner Interface

```go
func CreateQueryPlan(stmt *statements.SelectStatement, ctx PlannerContext) (*QueryPlan, error)
```

**Input**: 
- A parsed SELECT statement (AST)
- Planning context (database name, schema, etc.)

**Output**: 
- A `QueryPlan` object containing the logical plan tree

### PlannerContext

```go
type PlannerContext struct {
    Database string  // Current database name
    Schema   string  // Current schema name
}
```

Provides context for plan creation, including database and schema information.

## Plan Nodes

### Node Interface

All plan nodes implement:

```go
type PlanNode interface {
    NodeType() string               // Returns node type name
    OutputColumns() []table.TableColumn  // Returns output schema
}
```

### Node Types

#### 1. **ScanPlanNode**
Table scan operation.

```go
type ScanPlanNode struct {
    Database      string  // Database name
    Schema        string  // Schema name
    TableName     string  // Table to scan
    Alias         string  // Table alias (if any)
    IsSystemTable bool    // Whether this is a system table
}
```

**Semantics**: Reads all tuples from a table and emits them as rows.

#### 2. **ProjectPlanNode**
Column projection (SELECT list).

```go
type ProjectPlanNode struct {
    Input   PlanNode           // Child node
    Columns []items.SelectItem // Selected columns/expressions
}
```

**Semantics**: Applies the SELECT list to each input row, computing expressions and selecting columns.

#### 3. **FilterPlanNode**
WHERE clause filtering.

```go
type FilterPlanNode struct {
    Input     PlanNode             // Child node
    Condition *items.WhereClause   // WHERE condition
}
```

**Semantics**: Evaluates the WHERE predicate for each input row and emits only rows where the condition is true.

#### 4. **JoinPlanNode**
Table join operation.

```go
type JoinPlanNode struct {
    Left      PlanNode              // Left input
    Right     PlanNode              // Right input
    JoinType  string                // INNER, LEFT, RIGHT, FULL
    Condition parser.IExpressionContext // ON condition
}
```

**Semantics**: Joins rows from left and right inputs according to join type and ON condition.

**Supported Join Types**:
- INNER JOIN
- LEFT OUTER JOIN
- RIGHT OUTER JOIN
- FULL OUTER JOIN
- CROSS JOIN (no condition)

#### 5. **AggregatePlanNode**
GROUP BY aggregation.

```go
type AggregatePlanNode struct {
    Input      PlanNode          // Child node
    GroupBy    []items.SelectItem // GROUP BY expressions
    Aggregates []items.SelectItem // Aggregate functions
}
```

**Semantics**: Groups input rows by GroupBy expressions and computes aggregate functions for each group.

**Supported Aggregates**:
- COUNT
- SUM
- AVG
- MIN
- MAX
- ARRAY_AGG

#### 6. **SortPlanNode**
ORDER BY sorting.

```go
type SortPlanNode struct {
    Input   PlanNode            // Child node
    OrderBy []items.OrderByItem // ORDER BY specifications
}
```

**Semantics**: Sorts input rows according to ORDER BY expressions and sort directions.

#### 7. **LimitPlanNode**
LIMIT/OFFSET.

```go
type LimitPlanNode struct {
    Input  PlanNode // Child node
    Limit  int      // Number of rows to return
    Offset int      // Number of rows to skip
}
```

**Semantics**: Returns up to `Limit` rows, skipping `Offset` rows.

#### 8. **UnionPlanNode**
UNION operation.

```go
type UnionPlanNode struct {
    Left  PlanNode // Left input
    Right PlanNode // Right input
    All   bool     // true for UNION ALL, false for UNION
}
```

**Semantics**: Concatenates rows from left and right inputs. If `All` is false, removes duplicates.

#### 9. **IntersectPlanNode**
INTERSECT operation.

```go
type IntersectPlanNode struct {
    Left  PlanNode // Left input
    Right PlanNode // Right input
    All   bool     // true for INTERSECT ALL, false for INTERSECT
}
```

**Semantics**: Returns rows that appear in both left and right inputs. If `All` is false, removes duplicates.

#### 10. **ExceptPlanNode**
EXCEPT operation (set difference).

```go
type ExceptPlanNode struct {
    Left  PlanNode // Left input
    Right PlanNode // Right input
    All   bool     // true for EXCEPT ALL, false for EXCEPT
}
```

**Semantics**: Returns rows from left that don't appear in right. If `All` is false, removes duplicates.

#### 11. **ValuesPlanNode**
SELECT without FROM (constant values).

```go
type ValuesPlanNode struct {
    Values [][]parser.IExpressionContext // Literal value rows
}
```

**Semantics**: Returns constant rows without reading from storage.

#### 12. **SubqueryPlanNode**
Subqueries in FROM clause.

```go
type SubqueryPlanNode struct {
    Query *QueryPlan // Subquery plan
    Alias string     // Subquery alias
}
```

**Semantics**: Executes subquery and treats result as a table.

## Planning Algorithm

### Main Flow

```
buildPlanFromSelectStatement(stmt)
├─ if stmt.IsCombined():
│  └─ buildCombinedPlan() → UnionPlanNode, IntersectPlanNode, or ExceptPlanNode
└─ if stmt.IsPrimary():
   └─ buildPrimaryPlan()
      ├─ if no FROM: ValuesPlanNode
      └─ if has FROM:
         └─ buildFromClause() → ScanPlanNode or JoinPlanNode
            ├─ if WHERE: wrap in FilterPlanNode
            ├─ if GROUP BY: wrap in AggregatePlanNode
            ├─ if ORDER BY: wrap in SortPlanNode
            └─ if LIMIT: wrap in LimitPlanNode
```

### Key Planning Steps

1. **SELECT Type Detection**
   - Primary SELECT (simple query)
   - Combined SELECT (UNION/INTERSECT/EXCEPT)

2. **FROM Clause Processing**
   - Single table → ScanPlanNode
   - Multiple tables → Joined ScanPlanNode
   - Subquery → SubqueryPlanNode

3. **WHERE Filtering**
   - Wrap child in FilterPlanNode

4. **GROUP BY Aggregation**
   - Wrap child in AggregatePlanNode
   - Compute aggregate expressions

5. **ORDER BY Sorting**
   - Wrap child in SortPlanNode

6. **LIMIT/OFFSET**
   - Wrap child in LimitPlanNode

## Example Plans

### Simple Query
```sql
SELECT a, b FROM t1 WHERE x > 10 ORDER BY a LIMIT 5
```

Plan Tree:
```
Limit(5)
  Sort(a ASC)
    Filter(x > 10)
      Project(a, b)
        Scan(t1)
```

### Join Query
```sql
SELECT a, c FROM t1 JOIN t2 ON t1.id = t2.id
```

Plan Tree:
```
Project(a, c)
  Join(INNER, t1.id = t2.id)
    Scan(t1)
    Scan(t2)
```

### UNION Query
```sql
SELECT a FROM t1 UNION SELECT b FROM t2
```

Plan Tree:
```
Union(ALL=false)
  Project(a)
    Scan(t1)
  Project(b)
    Scan(t2)
```

### GROUP BY Query
```sql
SELECT dept, COUNT(*) FROM employees GROUP BY dept
```

Plan Tree:
```
Aggregate(GROUP BY [dept], COUNT(*))
  Scan(employees)
```

### Subquery
```sql
SELECT * FROM (SELECT * FROM t1 WHERE x > 5) sub
```

Plan Tree:
```
SubqueryPlanNode
  └─ Scan with Filter
```

## QueryPlan Structure

```go
type QueryPlan struct {
    Root              PlanNode                       // Root of plan tree
    Statement         *statements.SelectStatement    // Original SQL statement
    IsReadOnly        bool                           // true if SELECT
    SubqueryExecutor  subquery.SubqueryExecutor      // Subquery executor function
}
```

## Planning Strategies

### Join Ordering
- Currently uses left-to-right ordering from SQL
- Future: Cost-based optimizer can reorder joins

### Push-Down Optimization
- Filters pushed down to table scans (implicit)
- Projection pushed down (explicit in plan)

### Subquery Handling
- Correlated subqueries evaluated per row
- Uncorrelated subqueries can be cached

## Error Handling

The planner validates:
- **Column References**: Columns exist in tables
- **Type Compatibility**: Operations are valid for types
- **Query Structure**: Required clauses are present
- **Ambiguities**: Duplicate column names, etc.

Errors return descriptive messages indicating the semantic issue.

## Extensions

### Adding New Plan Nodes

1. Define struct implementing `PlanNode` interface
2. Implement `NodeType()` method
3. Implement `OutputColumns()` method for schema tracking
4. Add execution logic in `executor.go`
5. Update builder functions to create node type

### Adding Set Operations

1. Define PlanNode struct in `plan.go`
2. Add builder in `buildCombinedPlan()`
3. Add execution in `executeSet()` function

## Related Files

- **Planner**: `planner/planner.go`
- **Plan Nodes**: `planner/plan.go`
- **Executor**: `planner/executor.go`
- **Set Operations**: `planner/set_operations.go`
- **Statements**: `rsql/statements/`
- **Items**: `rsql/items/`

## Future Enhancements

1. **Cost-Based Optimization**: Estimate cardinality and choose optimal plans
2. **Parallel Execution**: Multi-threaded plan execution
3. **Incremental Planning**: Incremental compilation of subqueries
4. **Query Hints**: User guidance for plan selection
5. **Materialization Strategy**: Decide when to materialize intermediate results
