# 📘 Developer Guide: SelectBuilder Dialect Support

This document explains how `SelectBuilder` in Entiqon uses SQL dialects to generate engine-specific queries, such as PostgreSQL-compatible SELECT statements. It is builder-focused and assumes dialect resolution is already understood. For dialect internals, see the \[SQL Dialect Engine Guide].

---

## ✅ Overview

`SelectBuilder` now supports dialect-based query generation. This enables:

* Proper quoting of identifiers (`"id"`, `"users"`)
* Dialect-specific pagination (`LIMIT`, `OFFSET`)
* Accurate formatting of `WHERE` conditions

---

## 🧱 Dialect Injection

Use the `UseDialect(name string)` method to apply a dialect by name.

> ℹ️ The older `WithDialect(driver.Dialect)` method is now **deprecated** in favor of `UseDialect(...)`, which is more readable and consistent with fluent API design.

```go
query, args, err := NewSelect().
  Select("id", "created_at").
  From("users").
  Where("is_two_factor_enabled = true").
  UseDialect("postgres").
  Build()
```

> Internally uses `driver.ResolveDialect(...)` to apply quoting, escaping, and pagination behavior.

---

## 🛠️ Quoting Behavior

`SelectBuilder` applies dialect quoting rules automatically to:

* Column names in the SELECT clause
* The table name in the FROM clause
* The left-hand side of simple `WHERE` conditions

### Example

```sql
SELECT "id", "created_at" FROM "users" WHERE "is_two_factor_enabled" = true
```

> Aliases are **not quoted** by default.

---

## 📊 Limit/Offset Handling

Pagination is applied through the dialect’s `BuildLimitOffset(...)` method:

```go
LIMIT 10 OFFSET 5
```

The builder will fallback to manual LIMIT/OFFSET if no dialect is configured.

---

## 🔃 Migration Notes

* `UseDialect("postgres")` is preferred
* `WithDialect(driver.Dialect)` remains temporarily available but is deprecated
* Field names and tables should follow PostgreSQL naming conventions: lowercase and `snake_case`

---

## 📌 Summary

* `SelectBuilder` now fully supports dialect-aware query building
* Integration with quoting, condition parsing, and pagination is complete
* Backward compatibility is maintained with legacy method marked deprecated
* For underlying dialect implementation, see the **SQL Dialect Engine Guide**

---

For questions or contributions, refer to the Entiqon core SQL builder package.

\[SQL Dialect Engine Guide]: Sql Dialect Engine Guide
