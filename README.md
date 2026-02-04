# PetaCore

[![Tests](https://github.com/MBHuman/petacore/actions/workflows/tests.yml/badge.svg)](https://github.com/MBHuman/petacore/actions/workflows/tests.yml)
[![codecov](https://codecov.io/gh/MBHuman/petacore/branch/main/graph/badge.svg)](https://codecov.io/gh/MBHuman/petacore)
[![Go Report Card](https://goreportcard.com/badge/github.com/MBHuman/petacore)](https://goreportcard.com/report/github.com/MBHuman/petacore)

Distributed MVCC database с поддержкой SQL

## Возможности

- 🔄 **MVCC (Multi-Version Concurrency Control)**
- 📊 **SQL поддержка** через PostgreSQL wire protocol
- 🌐 **Распределенное хранилище** на базе ETCD или in-memory для тестов
- 🔒 **Snapshot Isolation** уровень изоляции транзакций