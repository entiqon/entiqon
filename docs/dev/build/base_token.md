# Entiqon Developer Guide: BaseToken

`BaseToken` is the foundational building block for SQL-like token structures (e.g., `Column`, `Table`, `Join`, `Condition`). It normalizes raw input expressions into a parsed `name` and optional `alias`, tracks the original `input` for diagnostics, enforces basic validation, and records a `kind` classification. Higher-level tokens embed `BaseToken` to reuse this logic, leaving context-specific resolution (e.g., table-qualification or ownership) to their own constructors.

> **Note:** `BaseToken` is *not* intended to be used standalone for producing SQL. It’s strictly an internal abstraction to simplify composition of more complex tokens.

---

## 🔍 Purpose & Scope

* **Parse & normalize** raw SQL-like input (e.g., `"users.id AS uid"`) into separate fields:

  * `input` (the original string)
  * `name`  (parsed identifier, unquoted)
  * `alias` (parsed alias, unquoted)
  * `Error` (if validation fails)

* **Provide nil-safe getters** via `contract.Rawable`:

  * `GetInput()` → raw input
  * `GetName()`  → parsed identifier
  * `GetAlias()` → parsed alias
  * `GetRaw()`   → raw SQL fragment (`name` or `name AS alias`)

* **Implement core interfaces**:

  * `contract.Rawable`    (GetInput, GetName, GetAlias, GetRaw)
  * `contract.Renderable` (Raw, RenderName, RenderAlias, String)
  * `contract.Kindable`   (GetKind, SetKind)
  * `contract.Errorable`  (GetError, IsErrored, SetError)

* **Emit diagnostic output** via `String()`, showing token `kind`, `name`, alias state, and any error.

---

## 🧱 Struct Definition

```go
// BaseToken provides parsing and error handling for a raw SQL token.
// It is designed to be embedded in higher-level token types.
//
// Fields:
//   input string        // original raw input (e.g. "users.id AS uid")
//   name  string        // parsed identifier (unquoted)
//   alias string        // parsed alias (unquoted)
//   Error error         // non-nil if parsing/validation failed
//   kind  contract.Kind // classification (e.g. ColumnKind)
//
// DEPRECATED FIELDS (for compatibility):
//   Name  string        // use GetName() instead
//   Alias string        // use GetAlias() instead
//
// BaseToken implements Rawable, Renderable, Kindable, and Errorable.
type BaseToken struct {
    input string
    name  string
    alias string

    // DEPRECATED: use GetName()
    Name  string
    // DEPRECATED: use GetAlias()
    Alias string

    Error error
    kind  contract.Kind
}
```

---

## 🚧 Constructor

```go
func NewBaseToken(raw string, explicitAlias ...string) *BaseToken
```

1. Trim & validate `raw` is non-empty, not comma-separated, doesn’t start with `AS `, and is not literally `AS`.
2. Record `input = raw`.
3. Parse inline alias via `ParseAlias(raw)`, populating `name` and `alias` (and writing `Name`/`Alias` for backward compatibility).
4. Apply explicit override if provided, recording an error on conflict.
5. Return `*BaseToken` with fields set and `Error` if any.

```go
b := NewBaseToken("users.id AS uid")
// b.GetInput() == "users.id AS uid"
// b.GetName()    == "users.id"
// b.GetAlias()   == "uid"
// b.Name         == "users.id"        // for compatibility
// b.Alias        == "uid"             // for compatibility
// b.GetRaw()     == "users.id AS uid"
```

---

## ✅ Methods Overview

| Method        | Signature                                           | Description                                              |
| ------------- | --------------------------------------------------- | -------------------------------------------------------- |
| `GetInput`    | `func (b *BaseToken) GetInput() string`             | Original raw input.                                      |
| `GetName`     | `func (b *BaseToken) GetName() string`              | Parsed identifier (unquoted).                            |
| `GetAlias`    | `func (b *BaseToken) GetAlias() string`             | Parsed alias (unquoted).                                 |
| `GetRaw`      | `func (b *BaseToken) GetRaw() string`               | `name` or `name AS alias`.                               |
| `Raw`         | `func (b *BaseToken) Raw() string`                  | Same as `GetRaw()`.                                      |
| `RenderName`  | `func (b *BaseToken) RenderName(q Quoter) string`   | Quoted name or alias via dialect.                        |
| `RenderAlias` | `func (b *BaseToken) RenderAlias(q Quoter, string)` | `qualified AS alias` with quoting.                       |
| `String`      | `func (b *BaseToken) String() string`               | Diagnostic representation with `kind`, alias, and error. |
| `GetError`    | `func (b *BaseToken) GetError() error`              | Underlying error (if any).                               |
| `IsErrored`   | `func (b *BaseToken) IsErrored() bool`              | Reports whether `Error != nil`.                          |
| `SetError`    | `func (b *BaseToken) SetError(input string, err)`   | Assigns `Error` and updates `input`.                     |
| `SetKind`     | `func (b *BaseToken) SetKind(k Kind)`               | Sets token classification.                               |
| `GetKind`     | `func (b *BaseToken) GetKind() Kind`                | Retrieves token classification.                          |
| `AliasOr`     | `func (b *BaseToken) AliasOr() string`              | `alias` if set, else `name`.                             |
| `IsAliased`   | `func (b *BaseToken) IsAliased() bool`              | `true` if `alias` non-empty.                             |
| `IsValid`     | `func (b *BaseToken) IsValid() bool`                | `true` if `name` non-empty and no `Error`.               |

---

## 📖 Examples

### Rawable usage

```go
var r contract.Rawable = NewBaseToken("users.id AS u")

fmt.Println(r.GetInput()) // users.id AS u
fmt.Println(r.GetName())  // users.id
fmt.Println(r.GetAlias()) // u
fmt.Println(r.GetRaw())   // users.id AS u
```

### Renderable usage

```go
var s contract.Renderable = NewBaseToken("orders.qty")

dialect := driver.NewPostgresDialect()

fmt.Println(s.Raw())               // orders.qty
fmt.Println(s.RenderName(dialect)) // "orders.qty"
fmt.Println(s.String())            // Unknown("orders.qty") [aliased: false, errored: false]
```

---

## 🚀 Generating Documentation

* **godoc CLI**: `godoc -http=:6060`
* **go doc**: `go doc github.com/entiqon/entiqon/internal/build/token.BaseToken`
* **IDE**: Hover over methods for inline GoDoc.

*Re-run after changes to refresh docs.*
