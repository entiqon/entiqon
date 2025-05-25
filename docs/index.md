# 📚 Entiqon Builder Documentation Index

Welcome to the Entiqon SQL Builder documentation. This suite includes validated, dialect-aware builders with full test
coverage and modular design.

---

## 📘 Core Builder Guides

- [SelectBuilder](dev/builder/select_builder.md) – Fluent SELECT queries with pagination, WHERE logic, and dialect
  quoting.
- [InsertBuilder](dev/builder/insert_builder.md) – Secure INSERT queries with multi-row support and RETURNING clause.
- [UpdateBuilder](dev/builder/update_builder.md) – UPDATE statements with SET chaining and conflict validation.
- [DeleteBuilder](dev/builder/delete_builder.md) – DELETE queries with conditional filters and dialect injection.
- [UpsertBuilder](dev/builder/upsert_builder.md) – INSERT ON CONFLICT resolution for PostgreSQL-compatible dialects.

---

## ⚙️ Architecture & Shared Concepts

- [Dialect Exposure Guide](dev/core/driver/dialect.md) *(WIP)* – Describes how custom dialects integrate with builders.
- [Token System Guide](dev/build/token.md) – Covers how tokens like `Column` are parsed, validated, and consumed by
  builders.
- [StageToken Usage](dev/builder/builder_guide.md#stagetoken) – Explains clause tagging and error traceability.
- [ParamBinder Flow](dev/builder/builder_guide.md#parambinder) – Covers how parameters are handled by dialect.

---

## 📦 Releases
- [Overview](./releases/index.md)
- [v1.6.0 - Keystone](./releases/release-notes-v1.6.0.md)
- [Full Changelog](./CHANGELOG.md)

---

All builders follow the same principles:

- ✅ 100% coverage or near
- 🔐 Validation-safe
- 🧠 Stage-aware error tagging
- 🧩 Compatible with external dialect extensions

Onward, builder.
