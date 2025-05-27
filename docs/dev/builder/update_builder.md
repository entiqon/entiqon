
**Fluent UPDATE builder with validation, error tagging and dialect-safe formatting.**

# 🛠️ UpdateBuilder Developer Guide

`UpdateBuilder` constructs SQL `UPDATE` statements using a fluent API that ensures clarity, correctness, and dialect-aware formatting.

---

## 📆 Since

**Available since:** `v1.2.0`

---


## 🔧 Method Reference
| Method        | Description                                         |
|---------------|-----------------------------------------------------|
| `WithDialect` | Sets SQL dialect for escaping                       |
| `Table`       | Sets the table to update                            |
| `Set`         | Assigns values to columns using `Field.WithValue()` |
| `Where`       | Sets the WHERE clause                               |
| `Returning`   | Adds a RETURNING clause                             |
| `Build()`     | Compiles the full SQL statement                     |

---

## ✍️ Example Usage

## ✅ Supported Features

- Target table declaration via `.Table(...)`
- Column assignment with `.Set(...)` (no aliasing allowed)
- Condition building using `.Where(...)`, `.AndWhere(...)`, `.OrWhere(...)`
- Dialect support for identifier quoting
- Deprecation-safe: `.WithDialect(...)` still works and is covered

---

## 🧱 Fluent API Example

```go
sql, args, err := builder.NewUpdate().
    Table("users").
    Set("status", "active").
    Where("id = ?", 42).
    UseDialect("postgres").
    Build()
```

Produces:

```sql
UPDATE "users" SET "status" = ? WHERE id = ?
```

Args:

```go
[]any{"active", 42}
```

---

## 🔐 Validation

- Requires `.Table(...)` and at least one `.Set(...)` call
- Rejects aliases: calling `.Set("email AS contact", ...)` returns an error
- Safe fallback: `.WithDialect(...)` is marked deprecated but still functional

---

## 🧪 Test Coverage

✅ **100% tested**, grouped by method and edge cases.

| Area             | Covered           |
|------------------|-------------------|
| `.Table(...)`    | ✅                 |
| `.Set(...)`      | ✅                 |
| `.Where/And/Or`  | ✅                 |
| `.Build()`       | ✅                 |
| `.UseDialect()`  | ✅                 |
| `.WithDialect()` | ✅                 |
| Validation       | ✅                 |
| Aliased field    | ✅ error triggered |

> 🧪 Even if necessary, **tests will be tested.**
> Because coverage isn't just a number — it's a philosophy.

---

## 💡 Best Practices

- ✅ Use `UseDialect(...)` for proper identifier quoting
- 🚫 Do not alias fields in `.Set(...)`
- 📏 Chain conditions instead of nesting or formatting manually

---

## 📚 Related

* [InsertBuilder Guide](./insert_builder.md)
* [SelectBuilder Guide](./select_builder.md)
* [DeleteBuilder Guide](./delete_builder.md)
* [UpsertBuilder Guide](./upsert_builder.md)
* [Dialect Guide](../driver/dialect.md)