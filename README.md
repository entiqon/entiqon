<p align="center">
    <img src="https://raw.githubusercontent.com/ialopezg/entiqon/main/assets/entiqon_black.png" alt="Entiqon Logo" width="200"/>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/ialopezg/entiqon"><img src="https://pkg.go.dev/badge/github.com/ialopezg/entiqon.svg" alt="Go Reference" /></a>
  <a href="https://goreportcard.com/report/github.com/ialopezg/entiqon"><img src="https://goreportcard.com/badge/github.com/ialopezg/entiqon" alt="Go Report Card" /></a>
  <a href="https://github.com/ialopezg/entiqon/actions/workflows/ci.yml"><img src="https://github.com/ialopezg/entiqon/actions/workflows/ci.yml/badge.svg" alt="Build Status" /></a>
  <a href="https://codecov.io/gh/ialopezg/entiqon"><img src="https://codecov.io/gh/ialopezg/entiqon/branch/main/graph/badge.svg" alt="Code Coverage" /></a>
  <a href="https://github.com/ialopezg/entiqon/releases"><img src="https://img.shields.io/github/v/release/ialopezg/entiqon" alt="Latest Release" /></a>
  <a href="https://ialopezg.github.io/entiqon/"><img src="https://img.shields.io/badge/docs-online-blue?logo=github" alt="Documentation" /></a>
  <a href="https://github.com/ialopezg/entiqon/blob/main/LICENSE"><img src="https://img.shields.io/github/license/ialopezg/entiqon" alt="License" /></a>
</p>

> ⚙️ A structured, intelligent foundation for building queryable, entity-aware Go systems.

---

## 🌱 Overview

Entiqon is a modular SQL query engine for Go focused on:

* 🧱 Composable and type-safe SQL builders
* 🔄 Dialect abstraction and pluggable formatting logic
* 🔍 Strict validation with tagged error context
* 🧪 Full method-based test coverage

---

## 🚀 Quick Start

```bash
go get github.com/ialopezg/entiqon
```

```go
sql, args, err := builder.NewSelect().
From("users").
Where("email = ?", "test@entiqon.dev").
Build()
```

---

## 📘 Developer Guides

### Architecture & Concepts

- [Builder Architecture](./builder_guide_updates.md) — Dialects, StageToken, ParamBinder integration

### Builders

- [SelectBuilder](./docs/developer/builder/select_builder.md) — Fluent SELECT with WHERE, AND, OR, LIMIT, and ordering
- [InsertBuilder](./docs/developer/builder/insert_builder.md) — Multi-row inserts and dialect-aware RETURNING
- [UpdateBuilder](./docs/developer/builder/update_builder.md) — No-alias column assignment and safe clause enforcement
- [DeleteBuilder](./docs/developer/builder/delete_builder.md) — DELETE with WHERE and optional RETURNING
- [UpsertBuilder](./docs/developer/builder/upsert_builder.md) — Full INSERT ... ON CONFLICT DO UPDATE/NOTHING support

### Extensions

- [Custom Driver Guide](./docs/developer/core/driver/custom_driver_guide.md) — How to implement dialects and understand
  quoting policies

---

## 📦 Releases

- [v1.5.0 - Without Sin](./releases/release_notes_v1.5.0.md)
- [CHANGELOG](./CHANGELOG.md)

---

## 📏 Principles & Best Practices

* 🧼 Clarity over brevity — use explicit method names
* 🚫 Deprecations are tested and clearly marked
* 🔐 Validate every path — no silent failures
* 🧩 Always quote identifiers through the dialect

---

## 🧩 Design Philosophy

* 📐 Chain → Validate → Compile
* 🧠 Tag errors with `StageToken`
* ⚙️ Compose with safe abstractions
* 📂 Group test methods visually

---

## 📄 License

[MIT](./LICENSE) — © Isidro Lopez / Entiqon Project