<img src="https://raw.githubusercontent.com/ialopezg/entiqon/main/assets/entiqon_black.png" alt="Entiqon Logo" style="width: 384px; display: block; margin: auto;" />

<p style="text-align: center; width: 384px; display: block; margin: auto;">
  <a href="https://pkg.go.dev/github.com/ialopezg/entiqon">
    <img src="https://pkg.go.dev/badge/github.com/ialopezg/entiqon.svg" alt="Go Reference" />
  </a>
  <a href="https://goreportcard.com/report/github.com/ialopezg/entiqon">
    <img src="https://goreportcard.com/badge/github.com/ialopezg/entiqon" alt="Go Report Card" />
  </a>
  <a href="https://github.com/ialopezg/entiqon/actions/workflows/ci.yml">
    <img src="https://github.com/ialopezg/entiqon/actions/workflows/ci.yml/badge.svg" alt="Build Status" />
  </a>
  <a href="https://codecov.io/gh/ialopezg/entiqon">
    <img src="https://codecov.io/gh/ialopezg/entiqon/branch/main/graph/badge.svg" alt="Code Coverage" />
  </a>
  <a href="https://github.com/ialopezg/entiqon/releases">
    <img src="https://img.shields.io/github/v/release/ialopezg/entiqon" alt="Latest Release" />
  </a>
  <a href="https://github.com/ialopezg/entiqon/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/ialopezg/entiqon" alt="License" />
  </a>
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

* [`SelectBuilder`](./builder/select.md) — Fluent SELECT with support for aliasing, ordering, and pagination
* [`InsertBuilder`](./builder/insert.md) — Multi-row inserts and RETURNING support
* [`UpdateBuilder`](./builder/update.md) — Strict value assignment and no-alias validation
* [`DeleteBuilder`](./builder/delete.md) — DELETE with optional RETURNING support
* [`UpsertBuilder`](./builder/upsert.md) — PostgreSQL-style UPSERT with conflict resolution

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
