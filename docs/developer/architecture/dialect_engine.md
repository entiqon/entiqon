# 🧱 Developer Guide: SQL Dialect Engine

This document explains the core infrastructure behind SQL dialects in Entiqon. It describes how dialects are defined, extended, resolved, and used across SQL builders like `SelectBuilder`, `InsertBuilder`, and others.

---

## ✅ Purpose

Dialect resolution allows builders to remain engine-agnostic (e.g., Postgres, MySQL, SQLite) while supporting:

* Quoting of table and column names
* Escape logic for diagnostic output
* Dialect-specific syntax (e.g., `ON CONFLICT`, `LIMIT`, `OFFSET`, `RETURNING`)

---

## 🧩 Dialect Interface

The `Dialect` interface defines how SQL builders interact with engine-specific behaviors. Each method provides a specific capability that a dialect may override.

```go
// Dialect represents SQL dialect-specific behaviors for quoting,
// escaping, pagination, and feature support.
type Dialect interface {
// Name returns the name of the dialect (e.g., "postgres").
Name() string

// Quote returns a quoted SQL identifier (e.g., column/table name)
// according to the dialect rules.
Quote(identifier string) string

// Escape returns a debug-safe string representation of a value.
// This is not meant for actual query building — use placeholders instead.
Escape(value any) string

// SupportsUpsert returns true if the dialect supports native UPSERT syntax.
SupportsUpsert() bool

// SupportsReturning returns true if the dialect supports RETURNING clauses.
SupportsReturning() bool

// BuildLimitOffset returns the dialect-specific LIMIT/OFFSET clause.
BuildLimitOffset(limit, offset int) string
}
```

---

## 🪜 BaseDialect Implementation

To simplify extension, Entiqon provides a `BaseDialect` struct that can be embedded:

```go
type BaseDialect struct {
DialectName string
}
```

This provides default behaviors for ANSI-style quoting and fallback escape/limit behavior.

---

## 🛠️ Resolving Dialects

Use the global `ResolveDialect(name string)` function to obtain the proper dialect implementation:

```go
switch name {
case "postgres":
return NewPostgresDialect()
case "mysql":
return NewMySQLDialect()
default:
return &BaseDialect{DialectName: "generic"}
}
```

---

## 📐 Quoting Behavior

### Rules:

* Uses **double quotes** for identifiers (ANSI standard, Postgres default)
* Applies quoting to: column names, table names, and optional expressions
* **Aliases are not quoted** by default

### Example:

```sql
SELECT "id", "created_at" FROM "users"
```

---

## ⛔ Escaping

* The `Escape(value any)` method is used only for diagnostics/debugging.
* Real query construction must use parameter placeholders (`?`, `$1`, etc.).

---

## 🧱 Dialect Extensions

To define a new dialect:

1. Create a struct embedding `BaseDialect`
2. Override only what differs
3. Register in `ResolveDialect()`

Example:

```go
type MySQLDialect struct {
BaseDialect
}

func (d *MySQLDialect) BuildLimitOffset(limit, offset int) string {
return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}
```

---

## 🔁 Dialect Usage in Builders

### ✅ Supported RETURNING Behavior by Dialect

| Dialect           | Method                | Returns |
|-------------------|-----------------------|---------|
| `BaseDialect`     | `SupportsReturning()` | `false` |
| `PostgresDialect` | `SupportsReturning()` | `true`  |
| Others (future)   | override as needed    |         |

Each builder (e.g., `SelectBuilder`, `InsertBuilder`) should:

* Support `UseDialect(name)` to apply dialect behavior
* Quote identifiers using `dialect.Quote(...)`
* Format pagination using `dialect.BuildLimitOffset(...)`
* Detect upsert support with `dialect.SupportsUpsert()`
* Detect `RETURNING` clause support via `dialect.SupportsReturning()`

> Legacy `WithDialect(driver.Dialect)` is deprecated in favor of `UseDialect(string)`

---

## 📌 Summary

* The dialect engine allows Entiqon to abstract over SQL syntax differences
* `BaseDialect` offers safe defaults
* `ResolveDialect(...)` centralizes instantiation
* Builders are responsible for invoking dialect methods during query generation
* Dialect capabilities like `RETURNING` are now validated per engine

For more information on how builders apply dialect behavior, see the individual builder guides such as `SelectBuilder Developer Guide`.
