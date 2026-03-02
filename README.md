# PetaCore

[![Tests](https://github.com/MBHuman/petacore/actions/workflows/tests.yml/badge.svg)](https://github.com/MBHuman/petacore/actions/workflows/tests.yml)
[![codecov](https://codecov.io/gh/MBHuman/petacore/branch/main/graph/badge.svg)](https://codecov.io/gh/MBHuman/petacore)
[![Go Report Card](https://goreportcard.com/badge/github.com/MBHuman/petacore)](https://goreportcard.com/report/github.com/MBHuman/petacore)

Distributed MVCC database with SQL support, implementing PostgreSQL compatibility and Snapshot Isolation for transactional consistency.

## Features

- 🔄 **MVCC (Multi-Version Concurrency Control)** - Non-blocking concurrent access with snapshot consistency
- 📊 **PostgreSQL-Compatible SQL** - Full SQL query support via PostgreSQL wire protocol (v3.0)
- 🌐 **Distributed Storage** - Built on ETCD for production deployments or in-memory for testing
- 🔒 **Snapshot Isolation** - Strong transactional guarantees without locking
- ⚡ **High Performance** - Modular query engine with clean separation of parsing, planning, and execution
- 🎯 **Standards Conformant** - Compatible with standard PostgreSQL tools and libraries
- 🧪 **Well-Tested** - Comprehensive test suite with CI/CD integration

## Quick Start

### Installation

```bash
git clone https://github.com/MBHuman/petacore.git
cd petacore
go build -o pcore ./cmd/api
```

### Running the Server

```bash
# Start API server on port 5432
go run cmd/pcore/main.go -store inmemory -log-file test.log

# Or with custom configuration
PCORE_PORT=5433 PCORE_DB=testdb go run cmd/pcore/main.go -store inmemory -log-file test.log
```

### Connecting with PostgreSQL Client

```bash
psql -h localhost -p 5432 -U postgres -d testdb
```

### Creating Tables and Running Queries

```sql
CREATE TABLE users (
    id BIGINT,
    name VARCHAR,
    age INTEGER
);

INSERT INTO users VALUES (1, 'Alice', 30);
INSERT INTO users VALUES (2, 'Bob', 25);

SELECT * FROM users WHERE age > 26;
```

## Architecture

PetaCore consists of several integrated subsystems:

**SQL Query Engine** → Parses SQL, builds logical execution plans, and executes them efficiently
**Transaction Layer** → Implements MVCC snapshot isolation with PostgreSQL compatibility
**Distributed Storage** → Uses ETCD backend for data persistence and distributed consensus
**Wire Protocol Server** → Implements PostgreSQL Frontend/Backend Protocol v3.0

See [internal/runtime/docs/Architecture.md](./internal/runtime/docs/Architecture.md) for detailed design documentation.

## Core Modules

### Runtime
The SQL execution engine with modular design:
- **Parser** - ANTLR4-based SQL parser generating AST
- **Planner** - Converts SQL to logical query plans
- **Executor** - Executes plans with stream processing
- **Functions** - Rich set of built-in SQL functions
- **Type System** - PostgreSQL-compatible types with casting

See [internal/runtime/README.md](./internal/runtime/README.md) for complete details.

### Storage
Multi-version concurrency control with distributed support:
- **Distributed Storage** - ETCD-backed persistent storage
- **MVCC** - Multi-version snapshots for read consistency
- **Transaction Manager** - Snapshot isolation implementation
- **In-Memory** - Optional in-memory backend for testing

### Wire Protocol
PostgreSQL Frontend/Backend Protocol v3.0:
- **Client Connection** - Startup, authentication, ready states
- **Query Execution** - Simple and extended protocols
- **Prepared Statements** - Type-safe parameterized queries
- **Result Streaming** - Efficient row-by-row transmission

## SQL Support

### Queries
- ✅ SELECT with all standard clauses
- ✅ JOINs (INNER, LEFT, RIGHT, FULL, CROSS)
- ✅ Set operations (UNION, INTERSECT, EXCEPT)
- ✅ Subqueries
- ✅ Aggregation (GROUP B)

### Data Definition
- ✅ CREATE TABLE
- ✅ DROP TABLE  
- ✅ TRUNCATE TABLE
- ✅ System tables (pg_class, pg_attribute, etc.)

### Data Manipulation
- ✅ INSERT
- ✅ DELETE

### Functions
- ✅ String functions (UPPER, LOWER, SUBSTR, etc.)
- ✅ Numeric functions (ABS, ROUND, SQRT, etc.)
- ✅ Aggregate functions (COUNT, SUM, AVG, MIN, MAX)
- ✅ Date/Time functions (NOW, DATE, EXTRACT, etc.)
- ✅ Type casting and coercion

## Concurrency & Transactions

PetaCore uses Snapshot Isolation for transactional consistency:


Features:
- **ACID Guarantees** - Atomicity, Consistency, Isolation, Durability
- **MVCC** - Readers never block writers, writers never block readers
- **Snapshot Isolation** - Each transaction sees consistent snapshot
- **Conflict Detection** - Detects and prevents write-write conflicts

## Performance

### Characteristics
- O(n) parsing and planning complexity
- Streaming execution where possible
- Lazy subquery evaluation
- Efficient memory usage through MVCC

### Optimization
- (Current) Basic plan generation
- (Planned) Cost-based optimization
- (Planned) Index support
- (Planned) Parallel execution

## Configuration

### Environment Variables

```bash
PCORE_PORT          # Server port (default: 5432)
PCORE_HOST          # Server host (default: localhost)
PCORE_DB            # Default database (default: postgres)
PCORE_LOG_LEVEL     # Logging level (default: info)
PCORE_STORAGE       # Storage backend: etcd or memory (default: memory)
PCORE_ETCD_ENDPOINTS # ETCD endpoints for distributed mode
```


## Building

### Prerequisites
- Go 1.16 or later
- ANTLR4 (for parser generation)

### Build Commands

```bash
# Build API server
go build -o pcore ./cmd/api

# Build REPL client
go build -o pcore-repl ./cmd/repl

# Run tests
go test ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Testing

```bash
# Run all tests
./scripts/run-tests.sh

# Run specific package tests
go test ./internal/runtime/...
go test ./internal/storage/...

# Run benchmarks
go test -bench=. ./internal/benchmark/...
```

## Documentation

- **[Architecture](./internal/runtime/docs/Architecture.md)** - System design overview
- **[Planner Documentation](./internal/runtime/docs/Planner.md)** - Query planning details
- **[Executor Documentation](./internal/runtime/docs/Executor.md)** - Plan execution engine
- **[Type System](./internal/runtime/docs/TypeSystem.md)** - SQL type handling
- **[Wire Protocol](./internal/runtime/docs/Wire.md)** - PostgreSQL protocol details
- **[Functions Reference](./internal/runtime/docs/Functions.md)** - Built-in functions
- **[Runtime Module](./internal/runtime/README.md)** - SQL runtime guide

## Project Structure

```
petacore/
├── cmd/                      # Command-line tools
│   ├── api/                  # PostgreSQL wire API server
│   ├── pcore/                # PetaCore CLI tool
│   ├── repl/                 # Interactive REPL
│   └── repl_vclock/          # REPL with vector clock support
├── internal/                 # Internal packages
│   ├── runtime/              # SQL query engine
│   │   ├── docs/             # Detailed documentation
│   │   ├── parser/           # SQL parser (ANTLR4)
│   │   ├── planner/          # Query planner
│   │   ├── executor/         # DDL executor
│   │   ├── functions/        # SQL functions
│   │   ├── rsql/             # Semantic layer
│   │   ├── rhelpers/         # Runtime helpers
│   │   ├── wire/             # PostgreSQL protocol
│   │   └── system/           # System tables
│   ├── storage/              # Data storage layer
│   │   ├── distributed_storage.go
│   │   ├── distributed_storage_vclock.go
│   │   └── ...
│   ├── core/                 # Core data structures
│   └── logger/               # Logging
├── sdk/                      # SDK and type system
│   ├── types/                # SQL type definitions
│   ├── serializers/          # Type serialization
│   └── pmem/                 # Memory allocation
├── plugins/                  # Extension plugins
├── tests/                    # Integration tests
├── examples/                 # Example queries
└── README.md                 # This file
```

## Design Philosophy

PetaCore follows principles inspired by DuckDB:

1. **Modularity** - Separate parser, planner, executor phases
2. **Clarity** - Clean code over clever optimization
3. **Extensibility** - Easy to add operators, functions, types
4. **Correctness** - Comprehensive testing and validation
5. **Standards Compliance** - PostgreSQL compatibility

## Limitations & Future Work

### Current Limitations
- Single-node by default (distributed via ETCD in progress)
- No query optimizer (default left-to-right plans)
- Limited aggregate functions
- No window functions yet
- No CTEs (WITH clause)
- No full-text search

### Planned Features
- ✨ Cost-based query optimization
- ✨ Parallel query execution
- ✨ Index support (B-tree, Hash)
- ✨ Window functions
- ✨ CTEs and recursive queries
- ✨ JSON/JSONB types
- ✨ Full-text search
- ✨ Columnar storage option (DuckDB-style)

## Contributing

We welcome contributions! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

See our [Contributing Guidelines](./CONTRIBUTING.md) for details.

## License

PetaCore is licensed under the MIT License. See [LICENSE](./LICENSE) for details.

## Community

- **Issues** - Report bugs and request features on GitHub
- **Discussions** - Ask questions and discuss ideas
- **Documentation** - Contribute to improving docs

## Related Projects

- [DuckDB](https://duckdb.org/) - Embedded OLAP database (design inspiration)
- [PostgreSQL](https://www.postgresql.org/) - World's most advanced database (compatibility target)
- [ETCD](https://etcd.io/) - Distributed key-value store (storage backend)

## Authors

PetaCore is developed and maintained by the community. See [CONTRIBUTORS.md](./CONTRIBUTORS.md) for a list of contributors.

---

**Status**: Early development. Use for evaluation and testing only in production-like scenarios.

For questions or support, please open an issue on GitHub.