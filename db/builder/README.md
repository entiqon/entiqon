<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true" align="left" height="96" width="96"> SelectBuilder
</h1>
<h6 align="left">Part of the <a href="https://github.com/entiqon/entiqon">Entiqon</a>::<span>Database</span> toolkit.</h6>

# 🌱 Overview

The `builder` package provides a fluent API to construct SQL `SELECT` queries.  
It is designed to be **simple**, **strict**, and **dialect-aware**.

## 🧩 Features

- **SelectBuilder** for building `SELECT` queries.
- Support for:
  - Fields (`Fields`, `AddFields`) with strict rules
  - Source (`Source`)
  - Conditions (`Where`, `And`, `Or`)
  - Grouping (`GroupBy`, `ThenGroupBy`)
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

Rules:
- `Where` resets conditions (like `Fields`).
- `And` appends with `AND`.
- `Or` appends with `OR`.
- Multiple conditions in one `Where(...)` are normalized with `AND`.

---

### Ordering

You can build `ORDER BY` clauses using `OrderBy` and `ThenOrderBy`:

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

Rules:
- `OrderBy` resets ordering (like `Fields`).
- `ThenOrderBy` appends additional fields.
- Empty or whitespace values are ignored.

---

### Grouping

You can build `GROUP BY` clauses using `GroupBy` and `ThenGroupBy`:

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

Rules:
- `GroupBy` resets grouping (like `Fields`).
- `ThenGroupBy` appends additional grouping fields.
- Empty or whitespace values are ignored.

---

### Debugging Fields

Use `String()` and `Debug()` to understand how a field was parsed:

```go
f := token.NewField("COUNT(*) AS total")

fmt.Println(f.String())
// ✅ Field("COUNT(*) AS total")

fmt.Println(f.Debug())
// ✅ Field("COUNT(*) AS total"): [raw: true, aliased: true, errored: false]

f2 := token.NewField(true)

fmt.Println(f2.String())
// ⛔️ Field("true"): input type unsupported: bool

fmt.Println(f2.Debug())
// ⛔️ Field("true"): [raw: false, aliased: false, errored: true] – input type unsupported: bool
```

---

### Error Cases

- `Build()` on a nil receiver:
  ```
  ❌ [Build] - Wrong initialization. Cannot build on receiver nil
  ```

- No source specified:
  ```
  ❌ [Build] - No source specified
  ```

- Invalid fields (detailed diagnostics):
  ```
  ❌ [Build] - Invalid fields:
      ⛔️ Field("true"): input type unsupported: bool
      ⛔️ Field("false"): input type unsupported: bool
      ⛔️ Field("123"): input type unsupported: int
  ```

- Raw expression with alias but without explicit `AS`:
  ```
  ⛔️ Field("(field1 || field2) alias"): [raw: true, aliased: false, errored: true] – raw expressions must use explicit AS for alias
  ```

---

## Status

Currently, supports:
- Field selection and aliasing (strict rules enforced)
- Single source
- WHERE conditions with AND/OR composition
- GROUP BY with multiple fields
- ORDER BY with multiple fields
- Limit and offset
- Error reporting for invalid fields with ✅/⛔️ diagnostics

Planned extensions include:
- Joins
- Parameter binding

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project
