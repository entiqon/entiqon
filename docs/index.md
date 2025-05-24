# 📚 Entiqon Builder Documentation Index

Welcome to the Entiqon SQL Builder documentation. This suite includes validated, dialect-aware builders with full test coverage and modular design.

---

## 📘 Core Builder Guides

- [SelectBuilder](./developer/builder/select_builder.md) – Fluent SELECT queries with pagination, WHERE logic, and dialect quoting.
- [InsertBuilder](./developer/builder/insert_builder.md) – Secure INSERT queries with multi-row support and RETURNING clause.
- [UpdateBuilder](./developer/builder/update_builder.md) – UPDATE statements with SET chaining and conflict validation.
- [DeleteBuilder](./developer/builder/delete_builder.md) – DELETE queries with conditional filters and dialect injection.
- [UpsertBuilder](./developer/builder/upsert_builder.md) – INSERT ON CONFLICT resolution for PostgreSQL-compatible dialects.

---

## ⚙️ Architecture & Shared Concepts

- [Dialect Exposure Guide](./developer/core/driver/dialect.md) *(WIP)* – Describes how custom dialects integrate with builders.
- [StageToken Usage](./developer/builder/builder_guide.md#stagetoken) – Explains clause tagging and error traceability.
- [ParamBinder Flow](./developer/builder/builder_guide.md#parambinder) – Covers how parameters are handled by dialect.

---

## 📦 Releases

- [CHANGELOG](./CHANGELOG.md)
- [Release Notes v1.5.0](./releases/release-notes-v1.5.0.md)

---

All builders follow the same principles:
- ✅ 100% coverage or near
- 🔐 Validation-safe
- 🧠 Stage-aware error tagging
- 🧩 Compatible with external dialect extensions

Onward, builder.
