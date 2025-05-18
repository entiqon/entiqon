<p align="center">
  <img src="./assets/entiqon_black.png" alt="Entiqon Logo" width="384"/>
</p>

<p align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/ialopezg/entiqon.svg)](https://pkg.go.dev/github.com/ialopezg/entiqon)
[![Latest Release](https://img.shields.io/github/v/release/ialopezg/entiqon)](https://github.com/ialopezg/entiqon/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/ialopezg/entiqon)](https://goreportcard.com/report/github.com/ialopezg/entiqon)
[![Build](https://github.com/ialopezg/entiqon/actions/workflows/ci.yml/badge.svg)](https://github.com/ialopezg/entiqon/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/ialopezg/entiqon/branch/main/graph/badge.svg)](https://codecov.io/gh/ialopezg/entiqon)
[![License](https://img.shields.io/github/license/ialopezg/entiqon)](https://github.com/ialopezg/entiqon/blob/main/LICENSE)

</p>

> ⚙️ A structured, intelligent foundation for building queryable, entity-aware Go systems in Go.

---

## 🌱 Overview

Entiqon is a modular query engine designed to:

* 🧱 Enable composable and type-safe SQL query construction
* 🔄 Provide dialect abstraction with pluggable escaping logic
* 🔍 Support strict validation and zero-tolerance safety
* 🧪 Ensure 100% test coverage with method-based test grouping

---

## ✅ Supported Builders

Each builder has full documentation and example usage:

* [`SelectBuilder`](./docs/Select%20Builder.md) — Fluent SELECT with support for aliasing, ordering, and pagination
* [`InsertBuilder`](./docs/Insert%20Builder.md) — Multi-row inserts and RETURNING support
* [`UpdateBuilder`](./docs/Update%20Builder.md) — Strict value assignment and no-alias validation
* [`DeleteBuilder`](./docs/Delete%20Builder.md) — DELETE with optional RETURNING support
* [`UpsertBuilder`](./docs/Upsert%20Builder%20Test.md) — PostgreSQL-style UPSERT with conflict resolution

---

## 🚀 Quick Start

### ↘️ Installation

```bash
go get github.com/ialopezg/entiqon
```

---

## 🧪 Examples by Builder

Every builder supports Go-style method chaining and returns the compiled SQL and argument slice:

```go
sql, args, err := builder.NewSelect().
	From("users").
	Where("email = ?", "test@entiqon.dev").
	Build()
```

For full examples, visit the documentation linked above.

---

## 🧩 Design Principles

* 📐 **Predictable structure**: every builder follows the same pattern: chain, validate, compile
* 🔐 **Strict validation**: no silent fallbacks, every mistake is caught early
* ⚙️ **Composable**: fields, clauses, and assignments are reusable and abstractable
* 🔄 **Dialects**: support for PostgreSQL (others pluggable)
* 📂 **Method grouping**: test files use visually grouped sections for clarity

---

## 📄 License

[MIT](LICENSE) — © Isidro Lopez / Entiqon Project
