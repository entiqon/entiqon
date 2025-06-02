<p align="center">
    <img src="https://raw.githubusercontent.com/entiqon/entiqon/main/assets/entiqon_black.png" alt="Entiqon Logo" width="200"/>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/entiqon/entiqon"><img src="https://pkg.go.dev/badge/github.com/entiqon/entiqon.svg" alt="Go Reference" /></a>
  <a href="https://goreportcard.com/report/github.com/entiqon/entiqon"><img src="https://goreportcard.com/badge/github.com/entiqon/entiqon" alt="Go Report Card" /></a>
  <a href="https://github.com/entiqon/entiqon/actions/workflows/ci.yml"><img src="https://github.com/entiqon/entiqon/actions/workflows/ci.yml/badge.svg" alt="Build Status" /></a>
  <a href="https://codecov.io/gh/entiqon/entiqon"><img src="https://codecov.io/gh/entiqon/entiqon/branch/main/graph/badge.svg" alt="Code Coverage" /></a>
  <a href="https://github.com/entiqon/entiqon/releases"><img src="https://img.shields.io/github/v/release/entiqon/entiqon" alt="Latest Release" /></a>
  <a href="https://entiqon.github.io/entiqon/"><img src="https://img.shields.io/badge/docs-online-blue?logo=github" alt="Documentation" /></a>
  <a href="https://github.com/entiqon/entiqon/blob/main/LICENSE"><img src="https://img.shields.io/github/license/entiqon/entiqon" alt="License" /></a>
</p>

> ⚙️ A structured, intelligent foundation for building queryable, entity-aware Go systems.

---

💡 Originally created by [Isidro Lopez](https://github.com/ialopezg)  
🏢 Maintained by the [Entiqon Organization](https://github.com/entiqon)

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
go get github.com/entiqon/entiqon
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

- [SelectBuilder](docs/dev/builder/select_builder.md)
- [InsertBuilder](docs/dev/builder/insert_builder.md)
- [UpdateBuilder](docs/dev/builder/update_builder.md)
- [DeleteBuilder](docs/dev/builder/delete_builder.md)
- [UpsertBuilder](docs/dev/builder/upsert_builder.md)

### Extensions

- [Custom Driver Guide](docs/dev/core/driver/custom_driver_guide.md)

---

## 📦 Releases

- [v1.6.0 - Keystone](./releases/release-notes-v1.6.0.md)
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