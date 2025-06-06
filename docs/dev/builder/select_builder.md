# 📘 SelectBuilder Developer Guide
**Dialect-aware SELECT builder with full clause tagging and error validation.**


This document explains the behavior, features, and dialect integration of Entiqon's `SelectBuilder`. It serves both as an implementation reference for contributors and a usage guide for consumers.

---

## ✅ Overview

`SelectBuilder` allows constructing SQL `SELECT` statements programmatically, with full support for:

* Fluent method chaining
* SQL identifier quoting
* Dialect-based pagination (e.g., `LIMIT/OFFSET`)
* WHERE condition logic

---

## 🧱 Method Reference

| Method          | Description                                |
|-----------------|--------------------------------------------|
| `NewSelect()`   | Creates a new select builder               |
| `Select(...)`   | Adds selected columns                      |
| `From(name)`    | Sets the table to select from              |
| `Where(...)`    | Adds a base WHERE clause                   |
| `AndWhere(...)` | Adds an AND condition                      |
| `OrWhere(...)`  | Adds an OR condition                       |
| `Take(n)`       | Sets LIMIT                                 |
| `Skip(n)`       | Sets OFFSET                                |
| `UseDialect()`  | Applies a named dialect (e.g., "postgres") |
| `Build()`       | Compiles the SQL query and argument list   |

> Deprecated: `WithDialect(driver.Dialect)` — use `UseDialect(name string)` instead.

---

## 🧩 Dialect Integration

When using `.UseDialect("postgres")`, `SelectBuilder` becomes dialect-aware:

* Table and column names are quoted using dialect rules (e.g., `"users"`, `"email"`)
* Pagination clauses (`LIMIT`, `OFFSET`) are inserted using dialect-specific formatting
* WHERE clause fields are also quoted if simple expressions like `field = ?` are used

Example:

```go
NewSelect().
  Select("id", "created_at").
  From("users").
  Where("is_active = true").
  UseDialect("postgres").
  Build()
```

Yields:

```sql
SELECT "id", "created_at" FROM "users" WHERE "is_active" = true
```

---

## 🚫 Aliases Are Not Quoted

Aliases set in column definitions (e.g., `AS something`) are passed through as-is and not quoted by the dialect. This avoids accidental case sensitivity in PostgreSQL and other engines.

---


---

## 🔍 Column–Table Resolution Matrix

When columns are added via `Select(...)` or `AddSelect(...)`, `SelectBuilder` performs automatic resolution when exactly one valid source table is defined using `From(...)` or `FromToken(...)`.

The following rules apply:

| Column Expression | Table Provided              | Outcome                 | Example Output |
|-------------------|-----------------------------|-------------------------|----------------|
| `"id"`            | `nil`                       | ✅ Render as-is          | `"id"`         |
| `"users.id"`      | `Name: "orders"`            | ❌ Error: name mismatch  | —              |
| `"users.id"`      | `Name: "users", Alias: "u"` | ✅ Renders with alias    | `"u.id"`       |
| `"id"`            | `Alias: "u"`                | ✅ Use alias as prefix   | `"u.id"`       |
| `"u.id"`          | `Alias: "u"`                | ✅ Remains as-is         | `"u.id"`       |
| `"u.id"`          | `Alias: "x"`                | ❌ Error: alias mismatch | —              |

> Columns with inline qualification (e.g., `"users.id"` or `"u.id"`) always take priority.
> If their qualification does not match the resolved source, a validation error is registered.

This mechanism enables intelligent column rendering without requiring repeated table prefixes.


## 🧪 Validation Rules

* `.From()` is required — `Build()` returns an error if missing
* `.Select()` can be empty (defaults to `*`)
* `.Where()` and other conditions are optional but support structured chaining
* `Take()` and `Skip()` apply pagination through the dialect if one is set

---

## 🧪 Test Coverage

| Function               | Coverage  | Status |
|------------------------|-----------|--------|
| Public API methods     | 100%      | ✅      |
| Internal builder logic | 100%      | ✅      |
| Dialect integration    | Verified  | ✅      |
| WHERE quoting          | Covered   | ✅      |
| Pagination logic       | Covered   | ✅      |
| Alias behavior         | Validated | ✅      |

---

## 📌 Summary

* `SelectBuilder` is safe, composable, and SQL-dialect aware
* Quoting, pagination, and clause generation adapt to the chosen dialect
* `UseDialect("postgres")` is the preferred way to activate dialect-specific behavior
* Fully tested and documented

For internal dialect engine details, see the SQL Dialect Engine Guide.