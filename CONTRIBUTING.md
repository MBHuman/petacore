# Contributing to PetaCore

Thank you for your interest in contributing to PetaCore! This document provides guidelines and instructions for participating in the project.

## Code of Conduct

We are committed to providing a welcoming and inspiring community for all. Please read and follow our [Code of Conduct](./CODE_OF_CONDUCT.md).

## Getting Started

### Prerequisites

Before you begin, ensure you have:

- **Go 1.16 or later** - Download from [golang.org](https://golang.org/doc/install)
- **Git** - For version control
- **ANTLR4** - Required for SQL parser (optional if not modifying grammar)
- **PostgreSQL client** - For testing (`psql` or similar)
- **Docker** (optional) - For running ETCD and PostgreSQL in containers

### Development Setup

1. **Fork the repository**
   ```bash
   # Go to https://github.com/MBHuman/petacore and click "Fork"
   ```

2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/petacore.git
   cd petacore
   ```

3. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/MBHuman/petacore.git
   ```

4. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-number
   ```

5. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

6. **Run tests to verify setup**
   ```bash
   ./scripts/run-tests.sh
   ```

## Areas for Contribution

### High-Impact Areas

- **Query Optimization** - Implement cost-based optimizer
- **Index Support** - Add B-tree and hash indexes
- **Window Functions** - Implement OVER clause
- **Performance** - Profile and optimize execution
- **Documentation** - Improve guides and examples
- **Testing** - Add comprehensive test coverage

### Good First Issues

- Adding built-in SQL functions
- Improving error messages
- Writing documentation
- Adding test cases
- Bug fixes in non-critical paths

### Architecture Areas

**Parser** (`internal/runtime/parser/`)
- SQL grammar enhancements
- Better error reporting
- Performance improvements

**Planner** (`internal/runtime/planner/`)
- Query optimization rules
- Cost estimation
- Plan rewriting

**Executor** (`internal/runtime/planner/executor.go`)
- New operators
- Performance tuning
- Memory optimization

**Storage** (`internal/storage/`)
- Index implementations
- Compression
- Persistence improvements

**Type System** (`sdk/types/`)
- New type support
- Serialization
- Coercion rules

## Workflow

### 1. Identify the Work

**For new features:**
- Check [GitHub Issues](https://github.com/MBHuman/petacore/issues) for existing discussions
- Open a discussion before starting major work
- Describe the feature and implementation approach

**For bugs:**
- Search existing issues first
- Include reproduction steps and expected behavior
- Attach error logs if relevant

**For improvements:**
- Discuss in an issue first for non-trivial changes
- Document the rationale for the change

### 2. Make Your Changes

#### Code Style

Write clear, idiomatic Go following the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments):

```go
// Good - clear variable names
func calculateQueryCost(nodes []PlanNode) float64 {
    totalCost := 0.0
    for _, node := range nodes {
        totalCost += node.EstimatedCost()
    }
    return totalCost
}

// Avoid - unclear names
func calc(n []N) float64 {
    t := 0.0
    for _, x := range n {
        t += x.C()
    }
    return t
}
```

**Conventions:**
- Use descriptive variable names
- Keep functions focused and small (< 50 lines preferred)
- Comment exported functions and types
- Add comments for complex logic
- Use standard library types and interfaces

#### Documentation

Update documentation for:
- **Public types** - Add godoc comments
- **Functions** - Explain parameters and return values
- **Packages** - Document purpose and key concepts
- **Complex logic** - Explain the "why", not just the "what"

Example:
```go
// CalculateSelectivity estimates the selectivity of a filter predicate.
// Returns a value between 0.0 (no rows) and 1.0 (all rows).
// Uses statistics if available, otherwise uses heuristics.
func (p *Planner) CalculateSelectivity(pred Filter) float64 {
```

#### Testing

All code changes must include tests:

```go
// Test function naming: Test<FunctionName>
func TestCalculateQueryCost(t *testing.T) {
    // Arrange
    nodes := []PlanNode{
        &ScanPlanNode{Table: "users"},
        &FilterPlanNode{Condition: ...},
    }

    // Act
    cost := calculateQueryCost(nodes)

    // Assert
    if cost <= 0 {
        t.Errorf("expected positive cost, got %f", cost)
    }
}
```

**Test Coverage:**
- Write tests for new functionality
- Aim for >80% coverage in core modules
- Include edge cases and error conditions
- Add integration tests for user-facing features

Run tests:
```bash
# All tests
go test ./...

# Specific package
go test ./internal/runtime/planner/...

# With coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./internal/benchmark/...
```

### 3. Commit Your Changes

#### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code style (formatting, semicolons, etc.)
- `refactor` - Code refactoring without feature changes
- `perf` - Performance improvements
- `test` - Test additions or modifications
- `chore` - Build, dependency, or tooling changes

**Scope:** One of `parser`, `planner`, `executor`, `storage`, `types`, `functions`, `wire`, `docs`

**Examples:**

```bash
git commit -m "feat(planner): implement join ordering optimization

- Use cardinality estimation to determine optimal join order
- Reduces query execution time by 30% on multi-table queries
- Includes comprehensive tests

Fixes #123"
```

```bash
git commit -m "fix(executor): handle NULL correctly in aggregate functions

COUNT(*) now correctly counts NULL values as 0.

Fixes #456"
```

```bash
git commit -m "docs(README): update quick start section"
```

#### Commit Guidelines

- Make logical, atomic commits
- One feature or fix per commit
- Don't mix concerns
- Write descriptive commit messages
- Reference related issues

### 4. Push and Create Pull Request

```bash
# Keep your branch up to date
git fetch upstream
git rebase upstream/main

# Push to your fork
git push origin feature/your-feature-name
```

#### Pull Request Description

```markdown
## Description
Brief description of changes and why they're needed.

## Related Issues
Fixes #123
Related to #456

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Added unit tests
- [ ] Added integration tests
- [ ] All tests pass
- [ ] Manual testing completed

## Documentation
- [ ] Updated README
- [ ] Updated relevant docs
- [ ] Added code comments for complex logic

## Performance Impact
- [ ] No performance changes
- [ ] Performance improvement (describe)
- [ ] Performance regression (explain mitigation)

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-reviewed the code
- [ ] Commented complex logic
- [ ] Updated documentation
- [ ] No new compiler warnings
- [ ] Tests pass locally
```

### 5. Code Review

**What to expect:**
- Constructive feedback on code quality
- Questions about design decisions
- Suggestions for improvements
- Requests for additional tests

**How to respond:**
- Address all comments
- Ask for clarification if needed
- Push additional commits for review changes
- Thank reviewers for feedback

## Testing Requirements

### Unit Tests

Test individual functions and classes:

```go
func TestFilterPlanNode_Execute(t *testing.T) {
    node := &FilterPlanNode{
        Condition: &WhereClause{...},
    }
    
    result, err := executeFilter(node, plan, ctx, tx, params)
    
    require.NoError(t, err)
    require.NotNil(t, result)
}
```

### Integration Tests

Test interactions between components:

```go
func TestSelectQuery_WithJoin(t *testing.T) {
    sql := "SELECT * FROM users JOIN orders ON users.id = orders.user_id"
    
    result, err := executeQuery(sql, testDB)
    
    require.NoError(t, err)
    require.Equal(t, 50, len(result.Rows))
}
```

### Regression Tests

Add tests that would catch the bug you're fixing:

```go
func TestAggregate_CountWithNULL(t *testing.T) {
    // This test would have caught the NULL bug
    rows := [][]interface{}{
        {1},
        {nil},
        {3},
    }
    
    count := aggregateCount(rows)
    assert.Equal(t, int64(2), count)  // NULL shouldn't be counted
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestName ./package

# Run with coverage report
go test -cover ./...

# Generate coverage HTML
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detector
go test -race ./...
```

## Coding Standards

### Go Best Practices

1. **Error Handling**
   ```go
   // Good
   if err != nil {
       return nil, fmt.Errorf("failed to parse query: %w", err)
   }

   // Avoid
   if err != nil {
       return nil, err
   }
   ```

2. **Interface Usage**
   ```go
   // Accept interfaces, return concrete types
   func Process(reader io.Reader) (*Result, error) {
       // ...
   }
   ```

3. **Naming**
   ```go
   // Receiver names should be short (1-2 chars)
   func (p *Planner) CreatePlan(stmt Statement) (*Plan, error) {

   // Interface names: -er suffix
   type Reader interface { Read() }

   // Avoid stuttering
   // Good: storage.New()
   // Bad: storage.NewStorage()
   ```

4. **Concurrency**
   ```go
   // Use mutexes for shared state
   type Cache struct {
       mu    sync.RWMutex
       data  map[string]interface{}
   }

   func (c *Cache) Get(key string) interface{} {
       c.mu.RLock()
       defer c.mu.RUnlock()
       return c.data[key]
   }
   ```

5. **Documentation**
   ```go
   // ExecutePlan executes a query plan and returns the result set.
   //
   // The execution happens within a transaction context and respects
   // snapshot isolation semantics.
   //
   // Parameters:
   //   plan - The logical execution plan to run
   //   ctx  - The execution context with storage and allocator
   //
   // Returns:
   //   A result set containing rows and schema, or an error if execution fails.
   func ExecutePlan(plan *QueryPlan, ctx ExecutorContext) (*table.ExecuteResult, error) {
   ```

### Performance Considerations

When implementing features:

1. **Avoid unnecessary allocations**
   ```go
   // Good - pre-allocate if you know size
   rows := make([]Row, 0, estimatedSize)

   // Avoid - grows repeatedly
   var rows []Row
   for ... {
       rows = append(rows, ...)
   }
   ```

2. **Use appropriate data structures**
   ```go
   // For membership testing use map
   seen := make(map[string]bool)

   // For ordered iteration use slice
   items := make([]Item, 0)
   ```

3. **Profile before optimizing**
   ```bash
   go test -cpuprofile=cpu.prof -bench=. ./internal/benchmark/...
   go tool pprof cpu.prof
   ```

## Documentation

### README Updates

When adding features, update relevant sections:
- Feature lists
- Quick start examples
- Architecture diagrams
- Configuration options

### Doc Files

Add detailed documentation in `internal/runtime/docs/`:
- Architecture changes
- New operators
- Type system extensions
- Protocol changes

### Examples

Add example queries to `examples/`:
- Complex JOINs
- Aggregations with GROUP BY
- Subqueries
- Set operations

### Inline Comments

Comment complex logic:

```go
// Estimate selectivity using column cardinality statistics if available,
// otherwise use heuristic based on operator type.
// This is a simplified model; production systems should use more sophisticated histograms.
selectivity := estimateSelectivityHeuristic(filter)
```

## Getting Help

### Communication Channels

- **GitHub Issues** - Report bugs and request features
- **GitHub Discussions** - Ask questions and discuss ideas
- **Pull Request Comments** - Technical discussions during review

### Useful Resources

- [PetaCore Documentation](./internal/runtime/docs/)
- [Go Documentation](https://golang.org/doc/)
- [SQL Standards](https://en.wikipedia.org/wiki/SQL)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## Contributor Recognition

Contributors are recognized in:
- [CONTRIBUTORS.md](./CONTRIBUTORS.md)
- Release notes
- GitHub contributors page

Thank you for contributions!

## Legal

By contributing to PetaCore, you agree that:

1. Your contribution is your own original work
2. You have the right to grant the license
3. Your contribution is licensed under the same license as the project

## Further Questions?

Open an issue with the `question` label or start a discussion on GitHub.

---

Happy coding! 🚀
