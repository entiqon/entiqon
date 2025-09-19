# SelectBuilder

> Part of [Entiqon](../../../) / [Database](../../) / [Builder](../)

The `SelectBuilder` constructs SQL `SELECT` queries in Go with a **fluent, safe, and dialect-aware API**.  
It is now part of the `builder/selects` subpackage.  


---

## ✨ Features

- Define fields, source tables, joins, conditions, grouping, ordering, having, and pagination.
- **Strict rules for field parsing**:
  - Single string → one or many fields (comma-split).
  - `"id AS alias"` or `"id alias"` → field with alias.
  - Two args (`expr, alias`) → field with alias.
  - Three args (`expr, alias, raw`) → raw expression with alias.
  - Passing an existing `Field` requires `.Clone()`, not `NewField`.
- Chainable mutators (`Fields`, `AddFields`, `Source`, `Join`, `Where`, …).
- Safe by design: invalid tokens are carried and surfaced at `Build()`.
- Default fallback to `SELECT *` if no fields are specified.
- Diagnostics: `String()` for concise view, `Debug()` for detailed state dump.

---

## 🚀 Quick Example

```go
import "github.com/entiqon/entiqon/db/builder/selects"

sb := selects.New(nil).
    Fields("id", "name").
    Source("users u").
    Where("u.active = true").
    OrderBy("created_at DESC").
    Limit(10).
    Offset(20)

sql, err := sb.Build()
if err != nil {
    log.Fatal(err)
}
fmt.Println(sql)
```

Output:

```sql
SELECT id, name
FROM users AS u
WHERE u.active = true
ORDER BY created_at DESC
LIMIT 10 OFFSET 20
```

---

## 🔍 Usage

### Fields

```go
sb := selects.New(nil).
    Fields("id, email AS user_email")
// SELECT id, email AS user_email
```

### Source

```go
sb := selects.New(nil).
    Fields("id").
    Source("users u")
// SELECT id FROM users AS u
```

### Joins

```go
sb := selects.New(nil).
    Fields("u.id", "o.id").
    Source("users u").
    InnerJoin("users u", "orders o", "u.id = o.user_id").
    LeftJoin("orders o", "payments p", "o.id = p.order_id").
    CrossJoin("currencies c").
    NaturalJoin("states s")
// SELECT u.id, o.id FROM users AS u
// INNER JOIN orders o ON u.id = o.user_id
// LEFT JOIN payments p ON o.id = p.order_id
// CROSS JOIN currencies c
// NATURAL JOIN states s
```

### Where

```go
sb := selects.New(nil).
    From("users").
    Where("active = true").
    And("country = 'USA'").
    Or("role = 'admin'")
// SELECT * FROM users WHERE active = true AND country = 'USA' OR role = 'admin'
```

### Group By / Having

```go
sb := selects.New(nil).
    Fields("department, COUNT(*) AS total").
    Source("users").
    GroupBy("department").
    Having("COUNT(*) > 5").
    AndHaving("AVG(age) > 30")
// SELECT department, COUNT(*) AS total
// FROM users
// GROUP BY department
// HAVING COUNT(*) > 5 AND AVG(age) > 30
```

### Order By

```go
sb := selects.New(nil).
    Fields("id, name").
    Source("users").
    OrderBy("created_at DESC").
    ThenOrderBy("id ASC")
// SELECT id, name FROM users ORDER BY created_at DESC, id ASC
```

### Pagination

```go
sb := selects.New(nil).
    Source("users").
    Limit(10).
    Offset(20)
// SELECT * FROM users LIMIT 10 OFFSET 20
```

---

## 🛠 Diagnostics

- `String()` → concise human-readable status  
- `Debug()` → verbose internal state with ✅/❌ markers

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Project
