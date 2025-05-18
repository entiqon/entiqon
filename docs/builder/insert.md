# InsertBuilder

The `InsertBuilder` constructs safe, parameterized SQL `INSERT INTO ... VALUES ...` statements with optional dialect-specific escaping.

---

## 🛠️ Constructor

```go
builder.NewInsert()
```

---

## 🔧 Method Reference

| Method              | Description                                           |
|---------------------|-------------------------------------------------------|
| `WithDialect`       | Sets SQL dialect for identifier escaping              |
| `Into`              | Sets the target table name                            |
| `Columns`           | Declares column names to insert into                  |
| `Values`            | Adds value row(s) to insert                           |
| `Returning`         | Specifies which columns to return (Postgres)          |
| `Build()`           | Builds the full SQL string with optional RETURNING    |
| `BuildInsertOnly()` | Builds just the INSERT ... VALUES part (no RETURNING) |

---

## ✍️ Example Usage

```go
q := builder.NewInsert().
	WithDialect(&dialect.PostgresEngine{}).
	Into("users").
	Columns("id", "email").
	Values(1, "dev@entiqon.dev").
	Returning("id")

sql, args, err := q.Build()
```

---

## 🔐 Clause Ordering

* `INSERT INTO ...`
* `VALUES (...)`
* (optional) `RETURNING ...`

---

## ⚠️ Validation Rules

`Build()` and `BuildInsertOnly()` will return an error when:

* No `Into(...)` is set
* No `Columns(...)` are set
* No `Values(...)` are added
* A row of values has a different length than the number of columns

---

## 🔄 Dialect Fallback Behavior

If no dialect is set:

* Raw identifiers are used as-is
* Placeholder (`?`) remains the same
* All validations still apply

---

## 🧪 Full Test Coverage

* Clause construction, including escaping
* Input validation errors
* Edge cases with multiple rows and types
