# Token

---

## 🎯 Purpose

The **token** module provides low-level primitives that represent SQL
fragments in a dialect-agnostic way.  
These tokens are consumed by higher-level builders (e.g. `SelectBuilder`)
to assemble safe, expressive, and auditable SQL statements.

---

## 📦 Sub-packages

| Package            | Purpose                                                                                               |
|--------------------|-------------------------------------------------------------------------------------------------------|
| [`field`](./field) | Represents a column, identifier, or computed expression (with aliasing, raw expressions, validation). |
| [`table`](./table) | Represents a SQL source (table or view) used in `FROM` / `JOIN` clauses with aliasing and validation. |
| [`join`](./join)   | Represents JOIN clauses (`INNER`, `LEFT`, etc.) with strict validation of join kind and conditions.   |

---

## 🔧 Module Helpers

| Module                                   | Purpose                                                                                                                                           |
|------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| [`resolver`](./resolver.go)              | Centralizes input type validation and expression resolution (expr, alias, kind), with subqueries.                                                 |
| [`ExpressionKind`](./expression_kind.go) | Classifies raw expressions (Identifier, Computed, Literal, Subquery, Function, Aggregate) and validates aliases, rejecting reserved SQL keywords. |

---

## 🚧 Roadmap

Planned tokens:
- **conditions** (WHERE / HAVING)
- **functions** (aggregates, JSON, custom expressions)

Contracts will progressively enforce stricter auditability across all tokens.

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project

