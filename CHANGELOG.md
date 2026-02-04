# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.1] - 2026-02-05

### Added
- 🎉 Initial release of PetaCore
- MVCC (Multi-Version Concurrency Control) implementation with HLC and VClock hybrid
- Distributed storage with plug and play storage engine support (For now only ETCD and InMemory supported)
- In-memory storage backend for testing
- PostgreSQL wire protocol compatibility
- Basic SQL support:
  - DDL: CREATE TABLE, DROP TABLE, TRUNCATE TABLE
  - DML: INSERT, SELECT
  - JOIN operations
  - GROUP BY support
  - System tables (information_schema, pg_catalog)
- Snapshot Isolation transaction level
- Vector clock-based conflict resolution
- Command-line flags for flexible configuration
- Comprehensive test suite:
  - Unit tests with coverage reporting
  - Integration tests via run-tests.sh
- GitHub Actions CI/CD pipeline
- Docker Compose configurations for ETCD, PostgreSQL, Redis
- REST API server
- REPL interface for interactive queries

### Features
- Store type selection: `etcd` or `inmemory`
- Configurable ETCD endpoints and key prefixes
- Automatic system tables initialization
- Wire protocol server on port 5432
- Graceful shutdown handling

### Documentation
- README with quickstart guide
- Coverage reporting setup
- Project structure documentation
- Command-line flags reference

[0.0.1]: https://github.com/MBHuman/petacore/releases/tag/v0.0.1
