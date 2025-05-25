# 📘 InsertBuilder Developer Guide
**Safe and validated INSERT builder with dialect injection and RETURNING support.**


This document describes the behavior, constraints, and dialect integration strategy for Entiqon's `InsertBuilder`. It ensures insert operations are consistent, secure, and aligned with SQL engine-specific capabilities.

---

## ✅ Overview

`InsertBuilder` constructs safe, parameterized SQL `INSERT` statements. It supports:

* Fluent configuration of table, columns, and values
* Multiple row insertion
* Dialect-specific identifier quoting
* Optional `RETURNING` clause (PostgreSQL, etc.)

---

## 🧱 Method Summary

| Method              | Description                                     |
| ------------------- | ----------------------------------------------- |
| `NewInsert()`       | Creates a new insert builder                    |
| `Into()`            | Sets the table name                             |
| `Columns()`         | Defines the column names (must match values)    |
| `Values()`          | Adds a row of values                            |
| `Returning()`       | Sets fields for RETURNING clause (if supported) |
| `UseDialect()`      | Injects a dialect by name (e.g., "postgres")    |
| `Build()`           | Renders full query + args                       |
| `BuildInsertOnly()` | Renders INSERT statement only (no RETURNING)    |

---

## 🛑 Column Aliases Are Not Allowed

Column aliasing (e.g., `email AS contact`) is rejected in `InsertBuilder`. If a column contains an alias, `Build()` returns an error:

```go
Columns("email AS contact") // ❌ invalid
```

> Aliases are meant for SELECT clauses and are not valid in INSERT column lists.

---

## 🧩 Dialect Integration

Use `UseDialect("postgres")` to apply dialect behavior:

* Quotes table and column names (e.g., `"users"`, `"email"`)
* Enables support for dialect-specific features like `RETURNING`
* Disallows `RETURNING` unless the dialect supports it

Example:

```go
Insert().
  Into("users").
  Columns("email", "role").
  Values("admin@example.com", "admin").
  Returning("id").
  UseDialect("postgres").
  Build()
```

Yields:

```sql
INSERT INTO "users" ("email", "role") VALUES (?, ?) RETURNING "id"
```

---

## 🧪 Test Coverage

Every path in `InsertBuilder` is covered:

| Function              | Coverage  | Status |
| --------------------- | --------- | ------ |
| Public API methods    | 100%      | ✅      |
| Internal builders     | 100%      | ✅      |
| Dialect integration   | Verified  | ✅      |
| RETURNING enforcement | Covered   | ✅      |
| Alias rejection       | Validated | ✅      |

* Full support for multiple `Values(...)` rows
* Guardrails for mismatched row/column counts
* Safe error when using `RETURNING` without a supported dialect
* Rejection of column aliases
* `BuildInsertOnly()` logic also dialect-aware

---

## ✅ RETURNING Clause Behavior

* Only allowed when `SupportsReturning()` returns `true`
* Currently supported by `PostgresDialect`
* Builder fails gracefully with an error for unsupported dialects

---

## 📌 Summary

* `InsertBuilder` is dialect-aware, secure, and strictly validated
* Compatible with PostgreSQL and fallback-safe with generic dialects
* Covers `Build()` and `BuildInsertOnly()` paths fully
* All behaviors tested to 100% coverage