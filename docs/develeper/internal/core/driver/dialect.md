# 🧱 Developer Guide: SQL Dialect Engine

This document explains the core infrastructure behind SQL dialects in Entiqon. It describes how dialects are defined, extended, resolved, and used across SQL builders like `SelectBuilder`, `InsertBuilder`, and others.

---

## ✅ Purpose

Dialect resolution allows builders to remain engine-agnostic (e.g., Postgres, MySQL, SQLite) while supporting:

* Quoting of table and column names
* Escape logic for diagnostic output
* Dialect-specific syntax (e.g., `ON CONFLICT`, `LIMIT`, `OFFSET`)

---

## 🧩 Dialect Interface

Each dialect must implement the following interface:

```go
type Dialect interface {
  Name() string
  Quote(identifier string) string
  Escape(value any) string
  SupportsUpsert() bool
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

Each builder (e.g., `SelectBuilder`) should:

* Support `UseDialect(name)` to apply dialect behavior
* Quote identifiers using `dialect.Quote(...)`
* Format pagination using `dialect.BuildLimitOffset(...)`
* Detect upsert support with `dialect.SupportsUpsert()` if applicable

> Legacy `WithDialect(driver.Dialect)` is deprecated in favor of `UseDialect(string)`

---

## 📌 Summary

* The dialect engine allows Entiqon to abstract over SQL syntax differences
* `BaseDialect` offers safe defaults
* `ResolveDialect(...)` centralizes instantiation
* Builders are responsible for invoking dialect methods during query generation

For more information on how builders apply dialect behavior, see the individual builder guides such as `SelectBuilder Developer Guide`.
