<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true" align="left" height="96" width="96"> Builder
</h1>
<h6 align="left">Part of the <a href="https://github.com/entiqon/entiqon">Entiqon</a>::<span>Database</span> toolkit.</h6>

# 🌱 Overview

The `builder` package provides a fluent API to construct SQL `SELECT` queries.  
It is designed to be simple, safe, and dialect-aware.

## 🧩 Features

- **SelectBuilder** for building `SELECT` queries.
- Support for:
  - Fields (`Fields`, `AddFields`)
  - Source (`Source`)
  - Pagination (`Limit`, `Offset`)
  - SQL rendering (`Build`, `String`)
- Default fallback to `SELECT *` if no fields are specified.
- Aggregated error reporting when invalid fields are provided.

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
// Output: SELECT id, name FROM "users" LIMIT 10 OFFSET 20
```

### Fields Handling

- **Fields**: resets the field collection and replaces existing fields.
- **AddFields**: appends to the existing field collection.

Supported inputs:
- `string`: single expression, comma-separated list, inline alias with `AS` or space
- `token.Field` / `*token.Field`
- `expr, alias`
- `expr, alias, isRaw`

Invalid inputs are collected and reported at build time.

### Error Cases

- `Build()` on a nil receiver:
  ```
  ❌ [Build] - Wrong initialization. Cannot build on receiver nil
  ```
- No source specified:
  ```
  ❌ [Build] - No source specified
  ```
- Invalid fields:
  ```
  ❌ [Build] - Invalid fields:
      ⛔️ Field("true"): unsupported type
      ⛔️ Field("123"): unsupported type
  ```

## Status

Currently, supports:
- Field selection and aliasing
- Single source
- Limit and offset
- Error reporting for invalid fields

Planned extensions include joins, ordering, grouping, and parameters.

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project