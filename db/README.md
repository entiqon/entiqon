# Database
Part of the [Entiqon Core](https://github.com/entiqon/entiqon).

## 🌱 Overview
Entiqon/db is a modular SQL query engine for Go, designed for composable, type-safe, and auditable SQL generation.

---

## 🧭 Doctrine
- **Never panic** — always return a token or builder, errors are embedded not thrown.
- **Auditability** — preserve user input for logs and error context.
- **Strict validation** — invalid expressions rejected early.
- **Delegation** — tokens own parsing/validation, builders compose them.

---

## 🛠 Capabilities
## 🛠 Capabilities

| Module                 | Feature                | Purpose                                                                    | Status      |
|------------------------|------------------------|----------------------------------------------------------------------------|-------------|
| [builder](./builder)   | [insert](./builder)    | High-level SQL builder for INSERT statements                               | 📝 Planned  |
|                        | [select](./builder)    | High-level SQL builder for SELECT statements (stable and production-ready) | ✅ Stable    |
|                        | [update](./builder)    | High-level SQL builder for UPDATE statements                               | 📝 Planned  |
|                        | [delete](./builder)    | High-level SQL builder for DELETE statements                               | 📝 Planned  |
|                        | [upsert](./builder)    | High-level SQL builder for UPSERT / MERGE statements                       | 📝 Planned  |
| [token](./token)       | [field](./token/field) | Dialect-agnostic representation of SQL fields/expressions                  | ✅ Stable    |
|                        | [table](./token/table) | Dialect-agnostic representation of SQL tables/sources                      | 🚧 On Going |
| [contract](./contract) | BaseToken              | Common base for tokens (shared identity, ownership, validity checks)       | 📝 Planned  |
|                        | Clonable               | Ensures semantic cloning of tokens without mutation                        | ✅ Stable    |
|                        | Debuggable             | Provides developer-facing diagnostics (`Debug()`)                          | ✅ Stable    |
|                        | Rawable                | Provides generic SQL output for logging (`Raw()`)                          | ✅ Stable    |
|                        | Renderable             | Provides dialect-aware SQL rendering (`Render()`)                          | ✅ Stable    |
|                        | Stringable             | Provides UX-facing, human-friendly string representations (`String()`)     | ✅ Stable    |

---

## 📄 License
[MIT](../LICENSE) — © Entiqon Project
