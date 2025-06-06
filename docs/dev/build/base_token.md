# Entiqon Developer Guide: BaseToken

`BaseToken` is the foundational building block for SQL‐like token structures (e.g., `Column`, `Table`, `Join`,
`Condition`). It normalizes raw input expressions into a parsed `Name` and optional `Alias`, tracks the original `input`
for diagnostics, enforces basic validation, and records a `Kind` classification. Higher‐level tokens embed `BaseToken`
to reuse this logic, leaving context‐specific resolution (e.g., table‐qualification or ownership) to their own
constructors.

> **Note:** `BaseToken` is *not* intended to be used standalone for producing SQL. It’s strictly an internal abstraction
> to simplify composition of more complex tokens.

---

## 🔍 Purpose & Scope

- **Parse & normalize** raw SQL‐like input (e.g., `"users.id AS uid"`) into separate fields:
    - `input` (the original string)
    - `Name` (e.g., `"users.id"`)
    - `Alias` (e.g., `"uid"`)
    - `Error` (if validation fails)

- **Provide safe accessors** to these fields (`GetInput`, `GetName`, `AliasOr`, `GetError`, etc.), all nil‐safe.

- **Implement core interfaces**:
    - `contract.Renderable` (methods: `Raw`, `RenderName`, `RenderAlias`, `String`)
    - `contract.Kindable` (`GetKind`, `SetKind`)
    - `contract.Errorable` (`GetError`, `IsErrored`, `SetError`)

- **Emit diagnostic output** via `String()`, showing token kind, name, alias state, and any error.

Higher‐level tokens (e.g., `Column`) embed `BaseToken` and add fields like `Qualified string` or other semantics; they
then call `Raw`, `RenderName`, or `RenderAlias` on the embedded `BaseToken` as needed.

---

## 🧱 Struct Definition

```go
// BaseToken provides normalized parsing and error handling for a raw SQL token.
// It is designed to be embedded in higher‐level token types (e.g., Column, Table).
//
// Fields:
//   input string        // the original raw input (e.g., "users.id AS uid")
//   Name  string        // the parsed identifier (e.g., "users.id")
//   Alias string        // optional alias (e.g., "uid")
//   Error error         // non‐nil if parsing/validation failed
//   kind  contract.Kind // classification (e.g., ColumnKind, TableKind, ConditionKind)
//
// BaseToken implements contract.Renderable, contract.Kindable, and contract.Errorable.
type BaseToken struct {
input string
Name  string
Alias string
Error error
kind  contract.Kind
}
```

> **Do not confuse** `input` (the original raw string) with `Name` or `Alias`. The field `input` is strictly for
> diagnostics; it does *not* affect rendering.

---

## 🚧 Constructor

### `func NewBaseToken(input string, aliasOverride ...string) *BaseToken`

1. **Trim & validate** `input` is not empty (after trimming).
2. **Reject** any comma (`,`) in `input` (aliases cannot be comma‐separated).
3. **Reject** if `input` (after trimming) starts with `"AS "` or resolves to `"AS"`.
4. **Parse inline alias** via `ParseAlias(input)`; derive:
    - `base` = everything before `AS`
    - `parsedAlias` = token after `AS` (if present)
5. **Apply explicit override** if `aliasOverride[0] != ""`:
    - If `parsedAlias != ""` and `aliasOverride[0] != parsedAlias`, then record an error:
      ```
      "alias conflict: explicit alias "%s" does not match inline alias "%s""
      ```
      But overwrite `Alias = aliasOverride[0]`.
    - Otherwise, set `Alias = aliasOverride[0]`.
6. **Populate** a new `BaseToken{ input: input, Name: base, Alias: parsedAlias (or override) }`.
7. **Return** the pointer; if any step failed, `Error` is non‐nil.

```go
b := NewBaseToken("users.id AS uid")
// → b.input == "users.id AS uid"
// → b.Name  == "users.id"
// → b.Alias == "uid"
// → b.Error == nil

b = NewBaseToken("users.id", "uid")
// → b.input == "users.id"
// → b.Name  == "users.id"
// → b.Alias == "uid"
// → b.Error == nil

b = NewBaseToken("users.id AS user_id", "uid")
// → b.input == "users.id AS user_id"
// → b.Name  == "users.id"
// → b.Alias == "uid"
// → b.Error == fmt.Errorf("alias conflict: explicit alias "uid" does not match inline alias "user_id"")

b = NewBaseToken("", "")
// → b.input == ""
// → b.Name  == ""
// → b.Alias == ""
// → b.Error == fmt.Errorf("invalid input expression: expression is empty")
```

---

## ✅ Methods Overview

| Method         | Signature                                                              | Description                                                                      |
|----------------|------------------------------------------------------------------------|----------------------------------------------------------------------------------|
| `AliasOr`      | `func (b *BaseToken) AliasOr() string`                                 | Returns `Alias` if non‐empty, otherwise returns `Name`. Nil‐safe (returns `""`). |
| `GetError`     | `func (b *BaseToken) GetError() error`                                 | Returns the underlying `Error` (or `nil`). Nil‐safe (returns `nil`).             |
| `GetInput`     | `func (b *BaseToken) GetInput() string`                                | Returns the original raw `input`. Nil‐safe (returns `""`).                       |
| `GetName`      | `func (b *BaseToken) GetName() string`                                 | Returns the parsed `Name`. Nil‐safe (returns `""`).                              |
| `HasError`     | `func (b *BaseToken) HasError() bool`                                  | (*alias*: calls `IsErrored()`). Reports whether `Error != nil`. Nil‐safe.        |
| `IsErrored`    | `func (b *BaseToken) IsErrored() bool`                                 | Returns `true` if `Error != nil`. Nil‐safe (returns `false`).                    |
| `IsAliased`    | `func (b *BaseToken) IsAliased() bool`                                 | Returns `true` if `Alias != ""`. Nil‐safe (returns `false`).                     |
| `IsValid`      | `func (b *BaseToken) IsValid() bool`                                   | Returns `true` if `Name != "" && Error == nil`. Nil‐safe (returns `false`).      |
| `Raw`          | `func (b *BaseToken) Raw() string`                                     | Returns raw SQL form: `"Name"` or `"Name AS Alias"`. Nil‐safe (returns `""`).    |
| `RenderAlias`  | `func (b *BaseToken) RenderAlias(q Quoter, qualified string) string`   |                                                                                  |
| `RenderName`   | `func (b *BaseToken) RenderName(q Quoter) string`                      | See detailed descriptions below.                                                 |
| `SetError`     | `func (b *BaseToken) SetError(input string, err error)`                | Assigns `b.Error = err`; sets `b.input = input`. Nil‐safe (no‐op).               |
| `SetErrorWith` | `func (b *BaseToken) SetErrorWith(input string, err error) *BaseToken` | Alias for `SetError`, allows chaining. Nil‐safe (no‐op).                         |
| `SetKind`      | `func (b *BaseToken) SetKind(k contract.Kind)`                         | Sets internal `kind`. Nil‐safe (no‐op).                                          |
| `GetKind`      | `func (b *BaseToken) GetKind() contract.Kind`                          | Returns `b.kind` or `UnknownKind` if nil/unset.                                  |
| `String`       | `func (b *BaseToken) String() string`                                  | Returns diagnostic string. Nil‐safe (returns `"nil"`).                           |

---

## 📖 Detailed Method Descriptions

### `AliasOr() string`

```go
func (b *BaseToken) AliasOr() string {
if b == nil {
return ""
}
if b.Alias != "" {
return b.Alias
}
return b.Name
}
```

- **Behavior**
    - If `b == nil`, returns `""`.
    - If `Alias != ""`, returns `Alias`.
    - Otherwise, returns `Name`.

- **Usage**
    - Useful when rendering SELECT clauses: prefer `Alias`, fallback to `Name`.

- **Example**
  ```go
  var b *BaseToken
  fmt.Println(b.AliasOr()) // → ""

  b = NewBaseToken("users.email AS em")
  fmt.Println(b.AliasOr()) // → "em"
  ```

---

### `GetError() error`

```go
func (b *BaseToken) GetError() error {
if b == nil {
return nil
}
return b.Error
}
```

- **Behavior**
    - If `b == nil`, returns `nil`.
    - Otherwise returns `b.Error`.

- **Usage**
    - Inspect parsing or alias conflict errors safely.

- **Example**
  ```go
  var b *BaseToken
  fmt.Println(b.GetError()) // → <nil>

  b = NewBaseToken("AS uid")
  fmt.Println(b.GetError()) // → invalid input expression: cannot start with 'AS'
  ```

---

### `GetInput() string`

```go
func (b *BaseToken) GetInput() string {
if b == nil {
return ""
}
return b.input
}
```

- **Behavior**
    - If `b == nil`, returns `""`.
    - Otherwise returns the original raw input.

- **Usage**
    - For logging or error diagnostics.

- **Example**
  ```go
  b := NewBaseToken("users.id AS uid")
  fmt.Println(b.GetInput()) // → "users.id AS uid"

  b = NewBaseToken("")
  fmt.Println(b.GetInput()) // → ""
  ```

---

### `GetName() string`

```go
func (b *BaseToken) GetName() string {
if b == nil {
return ""
}
return b.Name
}
```

- **Behavior**
    - If `b == nil`, returns `""`.
    - Otherwise returns `b.Name`.

- **Usage**
    - Retrieve the parsed identifier.

- **Example**
  ```go
  var b *BaseToken
  fmt.Println(b.GetName()) // → ""

  b = NewBaseToken("users.id")
  fmt.Println(b.GetName()) // → "users.id"
  ```

---

### `HasError() bool`

```go
func (b *BaseToken) HasError() bool {
return b.IsErrored()
}
```

- **Behavior**
    - Alias for `IsErrored()`.
    - Returns `false` if `b == nil` or `b.Error == nil`.
    - Otherwise `true`.

- **Usage**
    - Backward compatibility; prefer `IsErrored()` in new code.

- **Example**
  ```go
  b := NewBaseToken("users.id")
  fmt.Println(b.HasError()) // → false

  b = NewBaseToken("users.id AS user_id", "x")
  fmt.Println(b.HasError()) // → true
  ```

---

### `IsErrored() bool`

```go
func (b *BaseToken) IsErrored() bool {
return b != nil && b.Error != nil
}
```

- **Behavior**
    - Returns `false` if `b == nil` or `b.Error == nil`.
    - Otherwise `true`.

- **Usage**
    - Check if parsing or alias resolution failed.

- **Example**
  ```go
  b := NewBaseToken("users.id")
  fmt.Println(b.IsErrored()) // → false

  b = NewBaseToken("users.id AS user_id", "x")
  fmt.Println(b.IsErrored()) // → true
  ```

---

### `IsAliased() bool`

```go
func (b *BaseToken) IsAliased() bool {
return b != nil && b.Alias != ""
}
```

- **Behavior**
    - Returns `false` if `b == nil` or `b.Alias == ""`.
    - Otherwise `true`.

- **Usage**
    - Determine if there’s an `AS <alias>` clause.

- **Example**
  ```go
  b := NewBaseToken("users.id")
  fmt.Println(b.IsAliased()) // → false

  b = NewBaseToken("users.id AS uid")
  fmt.Println(b.IsAliased()) // → true
  ```

---

### `IsValid() bool`

```go
func (b *BaseToken) IsValid() bool {
return b != nil && b.Error == nil && strings.TrimSpace(b.Name) != ""
}
```

- **Behavior**
    - Returns `false` if `b == nil`, `b.Error != nil`, or `b.Name` is empty.
    - Otherwise `true`.

- **Usage**
    - Check before including in SQL generation.

- **Example**
  ```go
  b := NewBaseToken("users.id")
  fmt.Println(b.IsValid()) // → true

  b = NewBaseToken("")
  fmt.Println(b.IsValid()) // → false
  ```

---

### `Raw() string`

```go
func (b *BaseToken) Raw() string {
if b == nil {
return ""
}
if b.Alias != "" {
return fmt.Sprintf("%s AS %s", b.Name, b.Alias)
}
return b.Name
}
```

- **Behavior**
    - If `b == nil`, returns `""`.
    - If `Alias != ""`, returns `"Name AS Alias"`.
    - Otherwise returns `"Name"`.

- **Usage**
    - Embed in SQL when no quoting is needed.

- **Example**
  ```go
  b := NewBaseToken("users.id")
  fmt.Println(b.Raw()) // → "users.id"

  b = NewBaseToken("users.id AS uid")
  fmt.Println(b.Raw()) // → "users.id AS uid"
  ```

---

### `RenderAlias(q Quoter, qualified string) string`

```go
func (b *BaseToken) RenderAlias(q contract.Quoter, qualified string) string {
if b == nil || qualified == "" {
return qualified
}
if b.Alias == "" {
return qualified
}
if q == nil {
return fmt.Sprintf("%s AS %s", qualified, b.Alias)
}
return fmt.Sprintf("%s AS %s", qualified, q.QuoteIdentifier(b.Alias))
}
```

- **Behavior**
    1. If `b == nil` or `qualified == ""`, return `qualified`.
    2. If `b.Alias == ""`, return `qualified`.
    3. If `q == nil`, return `qualified + " AS " + Alias`.
    4. Otherwise, return `qualified + " AS " + q.QuoteIdentifier(Alias)`.

- **Usage**
    - Construct SELECT or JOIN clauses with proper quoting.

- **Example**
  ```go
  pg := driver.NewPostgresDialect()

  b := &BaseToken{Name: "id"}
  fmt.Println(b.RenderAlias(pg, `"users"."id"`)) // → `"users"."id"`

  b = &BaseToken{Name: "id", Alias: "uid"}
  fmt.Println(b.RenderAlias(pg, `"users"."id"`)) // → `"users"."id" AS "uid"`

  fmt.Println(b.RenderAlias(nil, `"users"."id"`)) // → `"users"."id" AS uid"`

  var bnil *BaseToken
  fmt.Println(bnil.RenderAlias(pg, `"users"."id"`)) // → `"users"."id"`
  ```

---

### `RenderName(q Quoter) string`

```go
func (b *BaseToken) RenderName(q contract.Quoter) string {
if b == nil || b.Name == "" {
return ""
}
if q == nil {
return b.Name
}
return q.QuoteIdentifier(b.Name)
}
```

- **Behavior**
    1. If `b == nil` or `b.Name == ""`, return `""`.
    2. If `q == nil`, return `b.Name`.
    3. Otherwise, return `q.QuoteIdentifier(b.Name)`.

- **Usage**
    - Render a quoted (or unquoted) identifier.

- **Example**
  ```go
  pg := driver.NewPostgresDialect()

  b := &BaseToken{Name: "id"}
  fmt.Println(b.RenderName(pg)) // → `"id"`
  fmt.Println(b.RenderName(nil)) // → "id"`

  var bnil *BaseToken
  fmt.Println(bnil.RenderName(pg)) // → ""
  ```

---

### `SetError(source string, err error)`

```go
func (b *BaseToken) SetError(source string, err error) {
if b == nil {
return
}
b.Error = err
if b.input != source {
b.input = source
}
}
```

- **Behavior**
    - If `b == nil`, do nothing.
    - Otherwise, set `b.Error = err` and update `b.input`.

- **Usage**
    - Called when higher‐level tokens detect semantic errors.

- **Example**
  ```go
  b := NewBaseToken("users.id")
  b.SetError("users.id", fmt.Errorf("permission denied"))
  fmt.Println(b.IsErrored()) // → true
  fmt.Println(b.GetError())  // → permission denied
  fmt.Println(b.GetInput())  // → "users.id"
  ```

---

### `SetErrorWith(source string, err error) *BaseToken`

```go
func (b *BaseToken) SetErrorWith(source string, err error) *BaseToken {
b.SetError(source, err)
return b
}
```

- **Behavior**
    - Same as `SetError`, returns `b` for chaining.

- **Usage**
    - Backward compatibility; use `SetError` otherwise.

- **Example**
  ```go
  b := NewBaseToken("users.id")
  b.SetErrorWith("users.id", fmt.Errorf("missing column")).SetKind(ColumnKind)
  fmt.Println(b.IsErrored()) // → true
  ```

---

### `SetKind(k contract.Kind)`

```go
func (b *BaseToken) SetKind(k contract.Kind) {
if b == nil {
return
}
b.kind = k
}
```

- **Behavior**
    - If `b == nil`, no‐op.
    - Otherwise set `b.kind = k`.

- **Usage**
    - Label the token type (Column, Table, Condition).

- **Example**
  ```go
  b := NewBaseToken("users")
  b.SetKind(contract.TableKind)
  fmt.Println(b.GetKind()) // → TableKind
  ```

---

### `GetKind() contract.Kind`

```go
func (b *BaseToken) GetKind() contract.Kind {
if b == nil {
return contract.UnknownKind
}
return b.kind
}
```

- **Behavior**
    - Returns `contract.UnknownKind` if `b == nil` or unset.
    - Otherwise returns `b.kind`.

- **Usage**
    - Inspect token classification.

- **Example**
  ```go
  b := &BaseToken{Name: "id"}
  fmt.Println(b.GetKind()) // → UnknownKind

  b.SetKind(contract.ColumnKind)
  fmt.Println(b.GetKind()) // → ColumnKind
  ```

---

### `String() string`

```go
func (b *BaseToken) String() string {
if b == nil {
return "nil"
}

kind := contract.UnknownKind
if b.kind != contract.UnknownKind {
kind = b.kind
}

suffix := fmt.Sprintf("[aliased: %t, errored: %t]", b.IsAliased(), b.IsErrored())
if b.IsErrored() {
suffix += fmt.Sprintf(", error: %s", b.GetError())
}

if kind != contract.UnknownKind {
return fmt.Sprintf(`%s("%s") %s`, kind.String(), b.GetName(), suffix)
}
return fmt.Sprintf(`Unknown("%s") %s`, b.GetName(), suffix)
}
```

- **Behavior**
    1. If `b == nil`, returns `"nil"`.
    2. Determine `kind`: `"Unknown"` or `kind.String()`.
    3. Build `suffix`: `"[aliased:…, errored:…]"` and `", error:…"` if errored.
    4. Return formatted string.

- **Usage**
    - Debugging or logging; not SQL‑valid.

- **Examples**
  ```go
  b := NewBaseToken("id")
  b.SetKind(contract.ColumnKind)
  fmt.Println(b.String())
  // → Column("id") [aliased: false, errored: false]

  b = NewBaseToken("users.id AS uid")
  b.SetKind(contract.ColumnKind)
  fmt.Println(b.String())
  // → Column("id") [aliased: true, errored: false]

  b = NewBaseToken("id AS uid", "wrong_alias")
  b.SetKind(contract.ColumnKind)
  fmt.Println(b.String())
  // → Column("id") [aliased: true, errored: true, error: alias conflict: explicit alias "wrong_alias" does not match inline alias "uid"]

  var bnil *BaseToken
  fmt.Println(bnil.String())
  // → "nil"
  ```

---

## 🧪 Unit Tests

Covered in `internal/build/token/base_test.go`:

- **RenderAlias & RenderName** with `contract.Quoter` (e.g., `driver.NewPostgresDialect()`) and `nil`.
- **Raw()** with alias/no‑alias and nil receiver.
- **String()** for all permutations.
- **SetError / GetError** behavior.
- **SetKind / GetKind** behavior.
- **AliasOr / GetInput / GetName / IsAliased / IsErrored / IsValid** edge cases.

---

## 📚 Quick Reference

| Use‑Case                                   | Method Call                                       | Return / Side Effect                                                                                                  |
|--------------------------------------------|---------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------|
| Parse input with inline alias              | `b := NewBaseToken("users.id AS uid")`            | `b.Name == "users.id"`, `b.Alias == "uid"`, `b.Error == nil`                                                          |
| Override inline alias explicitly           | `b := NewBaseToken("users.id", "uid")`            | `b.Alias == "uid"`, `b.Error == nil`                                                                                  |
| Detect alias conflict                      | `b := NewBaseToken("users.id AS user_id", "uid")` | `b.Alias == "uid"`, `b.Error.Error() == "alias conflict: explicit alias "uid" does not match inline alias "user_id""` |
| Retrieve raw SQL form                      | `b.Raw()`                                         | `"users.id AS uid"` (or `"users.id"`)                                                                                 |
| Get quoted identifier                      | `b.RenderName(pgDialect)`                         | `fmt.Sprintf(""%s"", b.Name)`                                                                                         |
| Get quoted alias from qualified expression | `b.RenderAlias(pgDialect, ""users"."id"")`        | `"users"."id" AS "uid"`                                                                                               |
| Check parsing errors                       | `b.IsErrored()` / `b.HasError()`                  | `true` if `b.Error != nil`                                                                                            |
| Get diagnostic string                      | `b.String()`                                      | `Column("id") [aliased: true, errored: false]`                                                                        |
| Access original input                      | `b.GetInput()`                                    | raw input (or `""` if nil)                                                                                            |
| Access parsed name                         | `b.GetName()`                                     | `b.Name` (or `""` if nil)                                                                                             |
| Fallback alias or name                     | `b.AliasOr()`                                     | `Alias` if set, else `Name` (or `""` if nil)                                                                          |
| Assign token kind                          | `b.SetKind(contract.ColumnKind)`                  | no panic if `b == nil`                                                                                                |
| Retrieve token kind                        | `b.GetKind()`                                     | `contract.UnknownKind` if unset/nil, else the assigned kind                                                           |

---

## 🚀 How to Generate Documentation

- **`godoc` CLI**:
  ```shell
  godoc -http=:6060
  ```  
- **`go doc`**:
  ```shell
  go doc github.com/entiqon/entiqon/internal/build/token.BaseToken
  ```  
- **IDE Plugins**: hover to see updated GoDoc.

Re-run after changes to keep docs in sync.

---

2025 — **© Entiqon Project**