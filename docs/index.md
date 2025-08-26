<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_logo.png?raw=true" align="left" height="82" width="82" alt="entiqon"> Entiqon Core
</h1>
<p align="left">A structured, intelligent foundation for building queryable, entity-aware Go systems.</p>

---

Welcome to the **Entiqon Core** documentation. This suite includes validated, dialect-aware builders with full test
coverage and modular design.

---

## 📦 Packages

| Icon                                                                                                                                          | Package                                   | Description                           |                                                Guides                                                 |
|-----------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------|---------------------------------------|:-----------------------------------------------------------------------------------------------------:|
| <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_sharicon.png?raw=true.png" height="32" width="32" alt="Common Icon" />  | [**Common**](packages/common/overview.md) | Shared utilities and helper functions | <img src="https://img.icons8.com/ios-glyphs/24/000000/checked-checkbox.png" width="20" height="20" /> |
| <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true.png" height="32" width="32" alt="Database Icon" /> | [**Database**](packages/database.md)      | Modular SQL Query builder             | <img src="https://img.icons8.com/ios-glyphs/24/000000/checked-checkbox.png" width="20" height="20" /> |

---

## 📘 Developer Guides

### 📚 Common

- [ProcessStage](packages/common/guides/ProcessStage_Developer_Guide.md) – How to use and extend the ProcessStage type
  for
  stage-aware error handling and workflow tracking.

### 📚 Database

- [SelectBuilder](dev/builder/select_builder.md) – Fluent SELECT queries with pagination, WHERE logic, and dialect
  quoting.
- [InsertBuilder](dev/builder/insert_builder.md) – Secure INSERT queries with multi-row support and RETURNING clause.
- [UpdateBuilder](dev/builder/update_builder.md) – UPDATE statements with SET chaining and conflict validation.
- [DeleteBuilder](dev/builder/delete_builder.md) – DELETE queries with conditional filters and dialect injection.
- [UpsertBuilder](dev/builder/upsert_builder.md) – INSERT ON CONFLICT resolution for PostgreSQL-compatible dialects.

---

## ⚙️ Architecture & Shared Concepts

- [Dialect Exposure Guide](dev/driver/dialect.md) *(WIP)* – Describes how custom dialects integrate with builders.
- [Styling Guide](dev/driver/styling.md) – Details `QuoteStyle`, `AliasStyle`, and `PlaceholderStyle` configuration and
  behavior.
- [Token System Guide](dev/build/token.md) – Covers how tokens like `Column` are parsed, validated, and consumed by
  builders.
- [StageToken Usage](dev/builder/builder_guide.md#stagetoken) – Explains clause tagging and error traceability.
- [ParamBinder Flow](dev/builder/builder_guide.md#parambinder) – Covers how parameters are handled by dialect.

---

## 📦 Releases

- [Overview](./releases/index.md)
- [v1.11.0 - Forge](./releases/release-notes-v1.11.0.md)
- [Full Changelog](./CHANGELOG.md)

---

All builders follow the same principles:

- ✅ 100% coverage or near
- 🔐 Validation-safe
- 🧠 Stage-aware error tagging
- 🧩 Compatible with external dialect extensions

Onward, builder.

---

## 📄 License

[MIT](../LICENSE) — © Isidro Lopez / Entiqon Project