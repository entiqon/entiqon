# Token System Guide

This guide explains how internal tokens like `Column` are structured, parsed, validated, and consumed by query builders.

## Principles

- Tokens are internal, immutable, and self-validating.
- Each token may carry an `Error`, which must be checked by the builder.
- Tokens do not handle dialect quoting, SQL rendering, or validation logic beyond structural correctness.

## Column

### Construction

Use `NewColumn(expr string, alias ...string)`. Supports:
- `"id"` → name only
- `"user_id AS uid"` → inline alias
- `NewColumn("email", "primary_email")` → explicit alias

### Fields

| Field | Description                     |
|-------|---------------------------------|
| Name  | The column name                 |
| Alias | Optional alias                  |
| Error | Non-nil if the input is invalid |

### Validation

Call `IsValid()` or check `Error != nil`.

```go
col := NewColumn(" AS uid")
if col.Error != nil {
	log.Warnf("Invalid column: %s — %v", col.String(), col.Error)
}
```

## Parser: ParseColumns

Use `ParseColumns(...)` to convert one or more inputs into `Column` tokens.

- Handles comma-separated values
- Keeps invalid tokens in-place with `Error` set
- Ideal for builder-stage filtering

```go
cols := ParseColumns("id, name", "user_id AS uid")
```

## Recommended Usage in Builders

1. Parse using `ParseColumns(...)`
2. Validate with `IsValid()` or `Error`
3. Add valid columns to builder
4. Record or skip invalid ones

## 🧠 Column Resolution Behavior

Columns support both inline qualification (e.g., `"users.id"`) and table-based attachment via `.WithTable()`.

### Qualification Rules

- A column is **qualified** if:
  - It includes an inline table prefix (`"users.id"`)
  - Or a table is attached via `.WithTable()` or `NewColumnWith(...)`

- A column is **invalid** if:
  - Its inline qualifier does not match the attached table’s name or alias
  - It is qualified but has no table context

### Field Overview

| Field     | Description                                                  |
|-----------|--------------------------------------------------------------|
| `Name`    | The column name                                              |
| `Alias`   | Optional alias (used in `AS`)                                |
| `Table`   | Attached `Table` token (used for alias-aware rendering)      |
| `TableName` | Inline qualifier parsed from input (e.g., `"users"` in `"users.id"`) |
| `Error`   | Non-nil if validation fails                                  |

---

### Examples

```go
users := NewTable("users AS u")

// Unqualified column, resolved using table alias
col := NewColumn("id", "user_id").WithTable(users)
fmt.Println(col.Raw()) // u.user_id

// Qualified column with match
col2 := NewColumn("users.id").WithTable(users)
fmt.Println(col2.Raw()) // u.id

// Qualified column with mismatch
bad := NewColumn("orders.id").WithTable(users)
fmt.Println(bad.Error) // column qualifier "orders" does not match table "users" (alias: "u")
```