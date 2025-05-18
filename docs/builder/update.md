# UpdateBuilder

The `UpdateBuilder` constructs SQL `UPDATE` statements with support for safe value assignments, WHERE clauses, and dialect-specific escaping.

---

## 🛠️ Constructor

```go
builder.NewUpdate()
```

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

```go
q := builder.NewUpdate().
	WithDialect(&dialect.PostgresEngine{}).
	Table("users").
	Set(
		builder.Field{Name: "email"}.WithValue("new@domain.com"),
	).
	Where("id = ?").
	Returning("id")

sql, args, err := q.Build()
```

---

## 🔐 Clause Ordering

1. `UPDATE ...`
2. `SET ...`
3. `WHERE ...`
4. `RETURNING ...`

---

## ⚠️ Validation Rules

`Build()` returns error if:

* `Table(...)` not set
* No fields passed to `Set(...)`
* Any field has an alias (`Name AS alias`) — **strictly disallowed**

---

## 🔄 Dialect Fallback Behavior

* Raw column names are used when no dialect is defined
* Internal safety checks protect against nil dereference

---

## 🧪 Full Test Coverage

* Strict alias validation (rejected)
* Dialect escaping with and without RETURNING
* Multi-field and single-field SET usage
