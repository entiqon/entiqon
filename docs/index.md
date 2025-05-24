# 📚 Entiqon Builder Documentation Index

Welcome to the Entiqon SQL Builder documentation. This suite includes validated, dialect-aware builders with full test coverage and modular design.

---

## 📘 Core Builder Guides

- [SelectBuilder](./select_builder_updated.md) – Fluent SELECT queries with pagination, WHERE logic, and dialect quoting.
- [InsertBuilder](./insert_builder_updated.md) – Secure INSERT queries with multi-row support and RETURNING clause.
- [UpdateBuilder](./update_builder_updated.md) – UPDATE statements with SET chaining and conflict validation.
- [DeleteBuilder](./delete_builder_updated.md) – DELETE queries with conditional filters and dialect injection.
- [UpsertBuilder](./upsert_builder_updated.md) – INSERT ON CONFLICT resolution for PostgreSQL-compatible dialects.

---

## ⚙️ Architecture & Shared Concepts

- [Dialect Exposure Guide](./dialect_engine.md) *(WIP)* – Describes how custom dialects integrate with builders.
- [StageToken Usage](./builder_guide_updates.md) – Explains clause tagging and error traceability.
- [ParamBinder Flow](./builder_guide_updates.md) – Covers how parameters are handled by dialect.

---

## 📦 Releases

- [CHANGELOG](./CHANGELOG_v1.5.0.md)
- [Release Notes v1.5.0](./release_notes_v1.5.0.md)

---

All builders follow the same principles:
- ✅ 100% coverage or near
- 🔐 Validation-safe
- 🧠 Stage-aware error tagging
- 🧩 Compatible with external dialect extensions

Onward, builder.
