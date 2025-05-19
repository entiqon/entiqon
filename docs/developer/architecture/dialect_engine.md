
# 🧭 Dialect Interface Guide

This guide explains how to implement and extend SQL dialects in Entiqon.

---

## 🔁 Interface Overview

Every dialect must implement the following methods:

```go
type Dialect interface {
    Name() string
    QuoteIdentifier(identifier string) string
    QuoteLiteral(value any) string
    BuildLimitOffset(limit, offset int) string
    SupportsUpsert() bool
    SupportsReturning() bool
}
```

---

## 🆕 Quoting Policy (since v1.2.0)

To improve clarity and safety, the old ambiguous method:

```go
Escape(value any)
```

has been replaced with two explicit alternatives:

| Method             | Purpose                      | Example Usage                     |
|--------------------|------------------------------|-----------------------------------|
| `QuoteIdentifier`  | Escapes table/column names   | `"user_id"`                       |
| `QuoteLiteral`     | Escapes literal values       | `'abc'`, `42`, `true`             |

### ⚠️ Warning
- `QuoteLiteral` is **not SQL-injection safe** and is meant only for logging/debugging.
- Do **not** use it in actual query strings.

---

## 🔧 Dialect Usage by Builder

| Builder         | Uses `QuoteIdentifier` | Uses `QuoteLiteral` | Requires Dialect? |
|------------------|------------------------|----------------------|--------------------|
| `SelectBuilder`  | ✅ Yes                 | ⚠️ Only for debug    | Optional           |
| `InsertBuilder`  | ✅ Yes                 | ⚠️ For logs only     | Optional           |
| `UpdateBuilder`  | ✅ Yes                 | ⚠️ For logs only     | Optional           |
| `DeleteBuilder`  | ✅ Yes                 | ❌ Not used           | Optional           |
| `UpsertBuilder`  | ✅ Yes                 | ⚠️ For logs only     | Optional           |

---

## 🗑️ Deprecated Methods

| Method         | Status        | Notes                          |
|----------------|---------------|--------------------------------|
| `Escape(...)`  | ❌ Removed     | Replaced by `QuoteIdentifier` and `QuoteLiteral` |
| `WithDialect`  | ⚠️ Deprecated | Use `UseDialect(...)` instead. Will be removed in v1.4.0.

---

## ✅ Example: PostgresDialect

```go
type PostgresDialect struct {
    BaseDialect
}

func (d *PostgresDialect) QuoteIdentifier(identifier string) string {
    return `"` + identifier + `"`
}
```

---

## 🔄 Migrating a Custom Dialect

If you've implemented your own dialect, follow these steps:
1. ✅ Rename `Escape(...)` → `QuoteLiteral(...)` (if for values)
2. ✅ Add `QuoteIdentifier(...)` for proper SQL quoting
3. 🔁 Update any usages of `Escape(...)` in builders

> All core builders now rely exclusively on `QuoteIdentifier`.

---

## 📚 Related

* [InsertBuilder Guide](./insert_builder.md)
* [UpdateBuilder Guide](./update_builder.md)
* [DeleteBuilder Guide](./delete_builder.md)
* [SelectBuilder Guide](./select_builder.md)
* [UpsertBuilder Guide](./upsert_builder_full_guide.md)
