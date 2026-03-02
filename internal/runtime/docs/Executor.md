# Query Executor

## Overview

The Query Executor is responsible for executing a logical query plan against the storage engine. It recursively traverses the plan tree, processes tuples at each step, and produces a result set.

## Architecture

### ExecutorContext

Provides the execution environment:

```go
type ExecutorContext struct {
    Database         string                          // Current database
    Schema           string                          // Current schema
    Storage          *storage.DistributedStorageVClock // Storage engine
    SubqueryExecutor subquery.SubqueryExecutor       // For nested queries
    Allocator        pmem.Allocator                  // Memory allocator
    GoCtx            context.Context                 // Go context
}
```

### Main Execution Function

```go
func ExecutePlan(plan *QueryPlan, ctx ExecutorContext) (*table.ExecuteResult, error)
```

**Execution Model**:
1. Begins a transaction if not in one
2. Recursively executes plan nodes
3. Returns final result set
4. Commits transaction on success

## Execution Pipeline

### Tuple Processing Flow

```
Input Tuples (from child node)
    ↓
Node-Specific Processing
(filtering, transformation, aggregation)
    ↓
Output Tuples (to parent node)
```

Each plan node type has corresponding execution logic:

### Node Execution Functions

#### executeScan()
Retrieves tuples from a physical table.

```go
func executeScan(
    node *ScanPlanNode,
    plan *QueryPlan,
    ctx ExecutorContext,
    tx *storage.DistributedTransactionVClock,
    runtimeParams map[int]interface{},
) (*table.ExecuteResult, error)
```

**Process**:
1. Resolve table metadata from storage
2. Retrieve all tuples from table
3. Return as ExecuteResult with schema

#### executeProject()
Applies the SELECT list to each input row.

```go
func executeProject(
    node *ProjectPlanNode,
    plan *QueryPlan,
    ctx ExecutorContext,
    tx *storage.DistributedTransactionVClock,
    runtimeParams map[int]interface{},
) (*table.ExecuteResult, error)
```

**Process**:
1. Execute child node to get input rows
2. For each input row:
   - Evaluate SELECT expressions
   - Extract selected columns
   - Construct output row with new schema
3. Return result with projected schema

**Expression Evaluation**:
- Column references: Extract from input schema
- Literals: Return constant value
- Functions: Call registered function
- Complex expressions: Recursive evaluation

#### executeFilter()
Evaluates WHERE clause conditions.

```go
func executeFilter(
    node *FilterPlanNode,
    plan *QueryPlan,
    ctx ExecutorContext,
    tx *storage.DistributedTransactionVClock,
    runtimeParams map[int]interface{},
) (*table.ExecuteResult, error)
```

**Process**:
1. Execute child node
2. For each row:
   - Evaluate WHERE condition
   - Include row only if condition is true
3. Return filtered rows with same schema

#### executeJoin()
Joins rows from two inputs.

```go
func executeJoin(
    node *JoinPlanNode,
    plan *QueryPlan,
    ctx ExecutorContext,
    tx *storage.DistributedTransactionVClock,
    runtimeParams map[int]interface{},
) (*table.ExecuteResult, error)
```

**Process**:

1. **Execute both children** to get left and right inputs

2. **For each join type**:

   **INNER JOIN**:
   - For each left row:
     - For each right row:
       - Evaluate ON condition
       - If true, create combined row

   **LEFT OUTER JOIN**:
   - Same as INNER
   - Plus: Include left row with NULLs if no right match

   **RIGHT OUTER JOIN**:
   - Same as INNER
   - Plus: Include right row with NULLs if no left match

   **FULL OUTER JOIN**:
   - Include all matched pairs
   - Include unmatched left rows with right NULLs
   - Include unmatched right rows with left NULLs

   **CROSS JOIN** (no ON clause):
   - Cartesian product of inputs

#### executeAggregate()
Implements GROUP BY with aggregation functions.

```go
func executeAggregate(
    node *AggregatePlanNode,
    plan *QueryPlan,
    ctx ExecutorContext,
    tx *storage.DistributedTransactionVClock,
    runtimeParams map[int]interface{},
) (*table.ExecuteResult, error)
```

**Process**:
1. Execute child node
2. Group rows by GROUP BY expressions:
   - Compute group key for each row
   - Collect rows into groups
3. For each group:
   - Evaluate GROUP BY expressions → group key
   - Compute aggregate functions:
     - COUNT: Count rows in group
     - SUM: Sum values
     - AVG: Average values
     - MIN/MAX: Minimum/maximum value
     - ARRAY_AGG: Collect values into array
4. Return one row per group

#### executeSort()
Sorts rows.

```go
func executeSort(
    node *SortPlanNode,
    plan *QueryPlan,
    ctx ExecutorContext,
    tx *storage.DistributedTransactionVClock,
    runtimeParams map[int]interface{},
) (*table.ExecuteResult, error)
```

**Process**:
1. Execute child node
2. Materialize all rows into memory
3. Sort by ORDER BY expressions:
   - Primary sort key: First ORDER BY item
   - Secondary sort keys: Subsequent items
   - Each can be ASC or DESC
4. Return sorted rows

#### executeLimit()
Applies LIMIT and OFFSET.

```go
func executeLimit(
    node *LimitPlanNode,
    plan *QueryPlan,
    ctx ExecutorContext,
    tx *storage.DistributedTransactionVClock,
    runtimeParams map[int]interface{},
) (*table.ExecuteResult, error)
```

**Process**:
1. Execute child node
2. Skip first OFFSET rows
3. Take up to LIMIT rows
4. Return result

#### executeUnion()
UNION operation.

```go
func executeSetOperation(
    left, right *table.ExecuteResult,
    operation string, // "UNION", "INTERSECT", "EXCEPT"
    all bool,
) (*table.ExecuteResult, error)
```

**UNION Semantics**:
- `UNION`: Concatenate, remove duplicates
- `UNION ALL`: Concatenate all rows

**Process**:
1. Execute left and right children
2. Concatenate result sets
3. If not ALL: Remove duplicate rows
4. Return combined result

#### executeIntersect()
INTERSECT operation (rows in both sets).

**Process**:
1. Execute left and right children  
2. For each left row:
   - Check if exists in right result set
   - If yes and not ALL: mark for output
   - If yes and ALL: output for each right occurrence
3. Return intersected rows

#### executeExcept()
EXCEPT operation (rows in left but not right).

**Process**:
1. Execute left and right children
2. For each left row:
   - Check if exists in right result set
   - If no: include in output
   - If yes and ALL: include once per net occurrence
3. Return difference

### Expression Evaluation

The executor evaluates expressions by dispatching to type-specific evaluators:

```go
func evaluateExpression(expr interface{}, row *table.Row, schema *table.Schema) (any, error)
```

**Expression Types**:
- **Literals**: Integer, float, string, boolean, NULL
- **Column References**: Extract from row using schema
- **Functions**: Dispatch to function registry
- **Operators**: Binary operators (+, -, *, /, <, >, etc.)
- **Subqueries**: Execute and fetch single value
- **CASE Expressions**: Conditional evaluation
- **Type Casts**: Convert between types

## Result Format

### ExecuteResult

```go
type ExecuteResult struct {
    Schema  *Schema      // Output schema
    Rows    [][]byte     // Serialized tuples
    RowCount int64       // Number of rows
}
```

### Schema

```go
type Schema struct {
    Fields []Field // Column definitions
}

type Field struct {
    Name  string     // Column name
    Type  oid.Oid    // SQL type OID
    // Additional metadata
}
```

## Performance Considerations

### Streaming vs. Materialization

- **Streaming**: Project, Filter pass through rows
- **Materialization**: Sort, Group By must buffer results

### Memory Usage

- Large intermediate results can consume significant memory
- Allocator manages memory across operations
- MVCC maintains version visibility

### Index Usage

Currently: Full table scans
Future: Index support for optimized Scan operations

## Error Handling

Executor handles errors at each level:
- **Type errors**: Incompatible operands
- **NULL handling**: Three-valued logic
- **Integer overflow**: Out of range values
- **Division by zero**: Special handling
- **Missing columns**: Reference errors

## Transaction Integration

Execution occurs within a transaction:

```go
err = ctx.Storage.RunTransactionWithAllocator(
    ctx.Allocator,
    func(tx *storage.DistributedTransactionVClock) error {
        // All plan execution happens here
        result, err = executePlanNode(plan.Root, ...)
        return err
    },
)
```

**Guarantees**:
- Snapshot isolation
- MVCC consistency
- Automatic rollback on error
- Atomic execution

## Special Operations

### CREATE TABLE

Executed outside query plan:

```go
ExecuteCreateTable(stmt, ctx) → (*ExecuteResult, error)
```

- Validates table definition
- Creates schema
- Creates system table entries

### DROP TABLE / TRUNCATE

Similar DDL handling:
- Removes table metadata
- Deallocates storage
- Updates system tables

## Subquery Evaluation

Subqueries can appear in:
- FROM clause → SubqueryPlanNode
- WHERE clause → Correlation evaluation
- SELECT list → Scalar subqueries

**Execution**:
1. Detect subquery in expression
2. Execute subquery plan
3. Return scalar value or cache result set
4. Use in parent expression evaluation

## Type Coercion

The executor automatically coerces types:

```
Integer ← String (parse)
Float ← Integer (promote)  
String ← Any (stringify)
Boolean ← Comparison results
```

## Related Files

- **Executor**: `planner/executor.go`
- **Result Types**: `rsql/table/`
- **Functions**: `functions/function.go`
- **Storage**: `../storage/`
- **Helpers**: `rhelpers/`
