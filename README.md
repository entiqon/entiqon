
<p align="center" style="text-align: center; width: 256px; display: block; margin: auto;">
    <img src="https://raw.githubusercontent.com/ialopezg/entiqon/main/assets/entiqon_black.png" align="center" alt="Entiqon Logo" style="width: 200px; display: block; margin: auto;" />
</p>
<br/>

<p align="center" style="text-align: center; width: 384px; display: block; margin: auto;">
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
  <a href="https://ialopezg.github.io/entiqon/">
    <img src="https://img.shields.io/badge/docs-online-blue?logo=github" alt="Documentation" />
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



## ✅ Developer Guides

Each builder is fully documented and independently tested:

* [SelectBuilder](./docs/developer/builder/select_builder.md) — Fluent SELECT with WHERE, AND, OR, LIMIT, and ordering
* [InsertBuilder](./docs/developer/builder/insert_builder.md) — Multi-row inserts and dialect-aware RETURNING
* [UpdateBuilder](./docs/developer/builder/update_builder_full_guide.md) — No-alias column assignment and safe clause enforcement
* [DeleteBuilder](./docs/developer/builder/delete_builder.md) — DELETE with WHERE and optional RETURNING
* [UpsertBuilder](./docs/developer/builder/upsert_builder_full_guide.md) — Full INSERT ... ON CONFLICT DO UPDATE/NOTHING support
* [Dialect Guide](./docs/developer/architecture/dialect_guide.md) — How to implement dialects and understand quoting policies

---

## 📏 Principles & Best Practices

* 🧼 Prefer clarity over brevity: use full method names (e.g., `QuoteIdentifier` instead of `QuoteIdent`)
* 🚫 Deprecated methods should be tested until removal and marked with a clear version timeline
* 📦 Builders should remain fluent and composable — every call must return the builder
* 🧪 Every public method must be test-covered (≥100%) — including deprecations
* 🧩 Avoid hardcoded identifiers; always route through dialect-safe quoting


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
