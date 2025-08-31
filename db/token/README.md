> Part of [Entiqon](https://github.com/entiqon/entiqon) / [Database](../)

# Token

---

## 🎯 Purpose

The **token** module provides low-level primitives that represent SQL
fragments in a dialect-agnostic way.  
These tokens are consumed by higher-level builders (e.g. `SelectBuilder`)
to assemble safe, expressive, and auditable SQL statements.

---

## 📦 Sub-packages

| Package                | Purpose                                                                                                                                                                         |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [`field`](./field)     | Represents a column, identifier, or computed expression (with aliasing, raw expressions, validation).                                                                           |
| [`table`](./table)     | Represents a SQL source (table or view) used in `FROM` / `JOIN` clauses with aliasing and validation.                                                                           |
| [`join`](./join)       | Represents JOIN clauses (`INNER`, `LEFT`, `RIGHT`, `FULL`, `CROSS`, `NATURAL`) using **join.Type** for strict validation and safe construction.                                 |
| [`types`](./types)     | Groups type enums (`identifier`, `join`, `condition`) that classify SQL expressions, joins, and conditions for consistent validation and rendering.                             |
| [`helpers`](./helpers) | Provides reusable validation utilities for identifiers, aliases, trailing aliases, reserved keywords, wildcards, deterministic alias generation, and expression classification. |

---

## 🚧 Roadmap

Planned tokens:
- **conditions** (WHERE / HAVING)
- **functions** (aggregates, JSON, custom expressions)

Contracts will progressively enforce stricter auditability across all tokens.

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project
