<h1><img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true" align="center" height="64" width="64"> Common Module</h1>

## 🌱 Overview

Entiqon is a modular SQL query engine for Go focused on:

* 🧱 Composable and type-safe SQL builders
* 🔄 Dialect abstraction and pluggable formatting logic
* 🔍 Strict validation with tagged error context
* 🧪 Full method-based test coverage

---

## 🚀 Quick Start

```bash
go get github.com/entiqon/db/v2
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

- [Builder Architecture](../docs/dev/builder/builder_guide.md) — Dialects, StageToken, ParamBinder integration

### Builders

- [SelectBuilder](../docs/dev/builder/select_builder.md)
- [InsertBuilder](../docs/dev/builder/insert_builder.md)
- [UpdateBuilder](../docs/dev/builder/update_builder.md)
- [DeleteBuilder](../docs/dev/builder/delete_builder.md)
- [UpsertBuilder](../docs/dev/builder/upsert_builder.md)

### Extensions

- [Custom Driver Guide](../docs/dev/core/driver/custom_driver_guide.md)

---

## 📄 License

[MIT](../LICENSE) — © Isidro Lopez / Entiqon Project