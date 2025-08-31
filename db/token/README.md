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
| [`helpers`](./helpers) | Provides reusable validation utilities for identifiers, aliases, trailing aliases, reserved keywords, wildcards, deterministic alias generation, and expression classification. |

---

## 🔧 Module Helpers

| Module                                   | Purpose                                                                                                                                                    |
|------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [`resolver`](./resolver.go)              | Centralizes input type validation and expression resolution (expr, alias, type), with subqueries and alias handling.                                       |
| [`ExpressionKind`](./expression_kind.go) | Classifies raw expressions (Identifier, Computed, Literal, Subquery, Function, Aggregate, Invalid). Provides short alias codes for deterministic aliasing. |
| [`identifier`](./identifier.go)          | Enum-based type classification for SQL expressions. Supports strict validation, parsing, alias generation, and safe rendering.                             |
| [`alias`](./helpers/alias.go)            | Validation utilities for aliases: `IsValidAlias`, `ValidateAlias`, reserved keyword rejection, trailing alias detection.                                   |
| [`wildcard`](./helpers/wildcard.go)      | Validation utilities for wildcards: `ValidateWildcard` ensures `*` cannot be aliased or raw.                                                               |
| [`trailing`](./helpers/trailing.go)      | Detects and validates trailing aliases when no `AS` keyword is present.                                                                                    |
| [`reserved`](./helpers/reserved.go)      | Provides reserved SQL keywords list and helpers to validate alias safety.                                                                                  |
| [`generator`](./helpers/generator.go)    | Deterministic alias generation using short codes and SHA-1 hashes of expressions.                                                                          |

---

## 🚧 Roadmap

Planned tokens:
- **conditions** (WHERE / HAVING)
- **functions** (aggregates, JSON, custom expressions)

Contracts will progressively enforce stricter auditability across all tokens.

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project

