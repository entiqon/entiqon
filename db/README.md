<p align="center">
    <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true.png" align="left" height="128" width="128">
</p>

---

# Entiqon DB Module

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
go get github.com/entiqon/db
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