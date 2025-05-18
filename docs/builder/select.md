# SelectBuilder

The `SelectBuilder` builds SQL `SELECT` queries with support for projections, filtering, ordering, and dialect-aware identifier escaping.

---

## 🛠️ Constructor

```go
builder.NewSelect()
```

---

## 🔧 Method Reference

| Method        | Description                                    |
|---------------|------------------------------------------------|
| `WithDialect` | Sets SQL dialect for identifier escaping       |
| `From`        | Sets the table or subquery to select from      |
| `Select`      | Adds raw column names or expressions           |
| `SelectAs`    | Adds `ColumnAs` structs with enforced aliasing |
| `Where`       | Sets the WHERE clause expression               |
| `OrderBy`     | Sets ORDER BY clause (raw string)              |
| `Limit`       | Sets LIMIT count                               |
| `Offset`      | Sets OFFSET value                              |
| `Build()`     | Builds the full SQL SELECT query               |

---

## ✍️ Example Usage

```go
q := builder.NewSelect().
	WithDialect(&dialect.PostgresEngine{}).
	Select("id", "email").
	From("users").
	Where("id = ?").
	Limit(10).
	Offset(5)

sql, args, err := q.Build()
```

---

## 🧩 Features

* Supports aliasing with `SelectAs` and `ColumnAs`
* Supports raw or parameterized WHERE clause
* Dialect-aware escaping for all identifiers

---

## ⚠️ Validation Rules

* `From(...)` is required; `Build()` returns error if missing
* `Select(...)` is optional (defaults to `*` if omitted)
* Supports both aliased and non-aliased column declaration

---

## 🔄 Dialect Fallback Behavior

* If no dialect is set, column/table names are used as-is
* Placeholders use `?` regardless of dialect

---

## 🧪 Test Coverage

* Select column combinations
* Aliased vs raw selections
* Dialect escaping edge cases
* Missing FROM clause validation
