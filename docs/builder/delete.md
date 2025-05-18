# DeleteBuilder

The `DeleteBuilder` helps you build SQL `DELETE FROM ... WHERE ...` statements with parameterization and optional dialect escaping.

---

## 🛠️ Constructor

```go
builder.NewDelete()
```

---

## 🔧 Method Reference

| Method        | Description                               |
|---------------|-------------------------------------------|
| `WithDialect` | Sets SQL dialect for escaping identifiers |
| `From`        | Sets the target table for deletion        |
| `Where`       | Sets the WHERE clause                     |
| `Returning`   | Adds RETURNING clause (e.g. Postgres)     |
| `Build()`     | Compiles the full SQL DELETE query        |

---

## ✍️ Example Usage

```go
q := builder.NewDelete().
	WithDialect(&dialect.PostgresEngine{}).
	From("users").
	Where("id = ?").
	Returning("id")

sql, args, err := q.Build()
```

---

## 🔐 Clause Ordering

1. `DELETE FROM ...`
2. `WHERE ...` (optional)
3. `RETURNING ...` (optional)

---

## ⚠️ Validation Rules

`Build()` returns an error if:

* `From(...)` is missing

---

## 🔄 Dialect Fallback Behavior

* If no dialect is set:

    * Raw identifiers are used
    * Placeholders remain `?`
    * All validations still apply

---

## 🧪 Test Coverage

* Table name validation
* WHERE clause behavior
* Dialect escaping
* RETURNING clause handling
