# UpsertBuilder

The `UpsertBuilder` allows you to construct an SQL `INSERT ... ON CONFLICT DO UPDATE` (UPSERT) statement programmatically, with support for parameter binding and dialect-specific identifier escaping.

---

## 🛠️ Constructor
```go
builder.NewUpsert()
```

---

## 🔧 Method Reference
| Method        | Description                                    |
|---------------|------------------------------------------------|
| `WithDialect` | Sets dialect for escaping identifiers          |
| `Into`        | Sets the target table                          |
| `Columns`     | Declares insert columns                        |
| `Values`      | Adds value rows for insertion                  |
| `OnConflict`  | Defines conflict detection columns             |
| `DoUpdateSet` | Assigns updates to apply on conflict           |
| `Returning`   | Specifies which columns to return              |
| `Build()`     | Compiles the full UPSERT SQL with placeholders |

---

## ✍️ Example Usage
```go
q := builder.NewUpsert().
	WithDialect(&dialect.PostgresEngine{}).
	Into("users").
	Columns("id", "email").
	Values(1, "dev@entiqon.dev").
	OnConflict("id").
	DoUpdateSet(
		builder.Assignment{Column: "email", Expr: "EXCLUDED.email"},
	).
	Returning("id", "email")

sql, args, err := q.Build()
```

---

## 🔐 Clause Ordering
The builder guarantees correct SQL ordering:
1. `INSERT INTO ...`
2. `VALUES (...)`
3. `ON CONFLICT (...)`
4. `DO UPDATE SET ...` or `DO NOTHING`
5. `RETURNING ...`

---

## ⚠️ Validation Rules
`Build()` will return an error (not panic) in the following cases:
- Missing `Into(...)`
- Missing `Columns(...)`
- Missing `Values(...)`
- Value count doesn't match column count

---

## 🔄 Dialect Fallback Behavior
If no dialect is set via `WithDialect`, the builder:
- Skips escaping — raw identifiers are used
- All features still work
- Internal methods defensively avoid nil dereferences

---

## 🧪 Full Test Coverage
- Grouped by method sections
- 100% tested: clause behavior, validation errors, and fallback scenarios
