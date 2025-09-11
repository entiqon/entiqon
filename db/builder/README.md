<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true" align="left" height="96" width="96"> SelectBuilder
</h1>
<p align="left">Part of the <strong><a href="https://github.com/entiqon/entiqon">Entiqon</a></strong> / <a href="">Database</a> toolkit.</p>

# 🌱 Overview

The `builder` package provides a fluent API to construct SQL `SELECT` queries.  
It is designed to be **simple**, **strict**, and **dialect-aware**.

## 🧩 Features

- **SelectBuilder** for building `SELECT` queries.
- Support for:
  - Fields (`Fields`, `AddFields`) with strict rules
  - Source (`Source`)
  - Joins (`InnerJoin`, `LeftJoin`, `RightJoin`, `FullJoin`, `CrossJoin`, `NaturalJoin`)
  - Conditions (`Where`, `And`, `Or`)
  - Grouping (`GroupBy`, `ThenGroupBy`)
  - Having (`Having`, `AndHaving`, `OrHaving`)
  - Ordering (`OrderBy`, `ThenOrderBy`)
  - Pagination (`Limit`, `Offset`)
  - SQL rendering (`Build`, `String`)
- Default fallback to `SELECT *` if no fields are specified.
- Aggregated error reporting when invalid fields or invalid usage are provided.

## Usage

### Basic Example

```go
sql, err := builder.NewSelect(nil).
    Fields("id, name").
    Source("users").
    Limit(10).
    Offset(20).
    Build()
if err != nil {
    log.Fatal(err)
}
fmt.Println(sql)
// Output: SELECT id, name FROM users LIMIT 10 OFFSET 20
```

---

### Join Example

```go
sql, err := builder.NewSelect(nil).
    Fields("u.id").
    AddFields("o.id").
    AddFields("p.amount").
    Source("users u").
    InnerJoin("users u", "orders o", "u.id = o.user_id").
    LeftJoin("orders o", "payments p", "o.id = p.order_id").
    CrossJoin("orders o", "currencies c").
    NaturalJoin("departments d", "states s").
    Where("u.active = true").
    OrderBy("p.amount DESC").
    Limit(10).
    Offset(20).
    Build()

if err != nil {
    log.Fatal(err)
}
fmt.Println(sql)

// Output:
// SELECT u.id, o.id, p.amount
// FROM users u
// INNER JOIN orders o ON u.id = o.user_id
// LEFT JOIN payments p ON o.id = p.order_id
// CROSS JOIN currencies c
// NATURAL JOIN states s
// WHERE u.active = true
// ORDER BY p.amount DESC
// LIMIT 10
// OFFSET 20
```

---

### Field Rules (Strict)

Fields are always normalized into a `token.Field`. The following rules are enforced:

1. **Single string**
   - `"id"` → one field  
   - `"id, name, email"` → multiple fields (comma split)  
   - `"id user_id"` → field with alias (`id AS user_id`)  
   - `"id AS user_id"` → field with alias (`id AS user_id`)  
   - `"COUNT(id) AS total"` → raw expression with alias  
   ⚠️ If a raw expression has a trailing alias without `AS`, it is rejected:
   ```go
   .Fields("(field1 || field2) alias")
   // Error: raw expressions must use explicit AS for alias
   ```

2. **Two arguments (`string, string`)**
   - `.Fields("id", "user_id")` → `id AS user_id`  
   ⚠️ This does not mean “two fields.” If you want two fields, use:
   ```go
   .Fields("id, user_id")
   ```

3. **Three arguments (`string, string, bool`)**
   - `.Fields("COUNT(*)", "total", true)` → raw expression with alias  

4. **Multiple `Field` objects**
   - `.Fields(*token.NewField("id"), *token.NewField("name"))`  
   ⚠️ Passing a `Field` into `NewField` is rejected; use `.Clone()` instead.

---

### Conditions

You can build `WHERE` clauses using `Where`, `And`, and `Or`:

```go
sql, _ := builder.NewSelect(nil).
    Fields("id, name").
    Source("users").
    Where("age > 18", "status = 'active'"). // normalized with AND
    Or("role = 'admin'").
    And("country = 'US'").
    Build()

fmt.Println(sql)
// Output: SELECT id, name FROM users WHERE age > 18 AND status = 'active' OR role = 'admin' AND country = 'US'
```

---

### Ordering

```go
sql, _ := builder.NewSelect(nil).
    Fields("id, name").
    Source("users").
    OrderBy("created_at DESC").
    ThenOrderBy("id ASC").
    Build()

fmt.Println(sql)
// Output: SELECT id, name FROM users ORDER BY created_at DESC, id ASC
```

---

### Grouping

```go
sql, _ := builder.NewSelect(nil).
    Fields("id, COUNT(*) AS total").
    Source("users").
    GroupBy("department").
    ThenGroupBy("role").
    Build()

fmt.Println(sql)
// Output: SELECT id, COUNT(*) AS total FROM users GROUP BY department, role
```

---

### Having

```go
sql, _ := builder.NewSelect(nil).
    Fields("department, COUNT(*) AS total").
    Source("users").
    GroupBy("department").
    Having("COUNT(*) > 5").
    AndHaving("AVG(age) > 30").
    OrHaving("SUM(salary) > 100000").
    Build()

fmt.Println(sql)
// Output: SELECT department, COUNT(*) AS total FROM users GROUP BY department HAVING COUNT(*) > 5 AND AVG(age) > 30 OR SUM(salary) > 100000
```

---

### Debugging Fields

Use `String()` and `Debug()` to understand how a field was parsed.

---

## Status

Currently, supports:
- Field selection and aliasing (strict rules enforced)
- Single source
- Joins (INNER, LEFT, RIGHT, FULL, CROSS, NATURAL)
- WHERE conditions with AND/OR composition
- GROUP BY with multiple fields
- HAVING with AND/OR composition
- ORDER BY with multiple fields
- Limit and offset
- Error reporting for invalid fields with ✅/⛔️ diagnostics

Planned extensions include:
- Parameter binding
- Subqueries as sources

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project
