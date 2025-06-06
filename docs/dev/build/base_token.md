# Entiqon Developer Guide: BaseToken

`BaseToken` is the foundational building block for SQL token structures such as `Column`, `Table`, or `Join`. It provides normalized parsing of raw SQL expressions, manages name and alias resolution, and standardizes error handling across tokens.

---

## 🔍 Purpose

`BaseToken` is designed for internal composition of higher-level tokens and is **not meant to be used on its own** in the rendering pipeline.

It helps:
- Parse and normalize input expressions
- Extract alias and name
- Store raw source input for traceability
- Handle validation errors early

---

## 🧱 Structure

```go
type BaseToken struct {
    Source string // Original raw input
    Name   string // Resolved name
    Alias  string // Optional alias
    Error  error  // Semantic or structural conflict
}
```

---

## 🛠 Constructor

### `func NewBaseToken(input string, alias ...string) *BaseToken`

Performs:
- Input trimming and validation
- Inline alias resolution (e.g., `AS`)
- Optional alias override
- Error assignment if malformed

#### Examples:

```go
NewBaseToken("users.id")               // ✅ Valid
NewBaseToken("users.id AS uid")       // ✅ Valid inline alias
NewBaseToken("users.id", "uid")       // ✅ Valid explicit alias
NewBaseToken("users.id AS uid", "x")  // ⚠️ Alias conflict: x != uid
NewBaseToken("AS uid")                // ❌ Invalid: cannot start with 'AS'
NewBaseToken("")                      // ❌ Invalid: input is empty
NewBaseToken("id, name")              // ❌ Invalid: comma-separated
```

---


## ✅ Methods Overview

| Method           | Description                                                |
|------------------|------------------------------------------------------------|
| `GetName()`      | Returns the token name safely (even if nil)                |
| `GetSource()`    | Returns the original source string safely                  |
| `AliasOr()`      | Returns alias or name fallback                             |
| `IsAliased()`    | Returns true if alias is defined                           |
| `IsErrored()`    | Returns true if error is set                               |
| `Validate()`     | (planned) Assigns error if malformed or incomplete         |
| `GetKind()`      | Returns the Kind classification of the token               |
| `SetKind()`      | Assigns a Kind to the token                                |
| `SetErrorWith()` | Manually set error and associate with source               |
| `Raw()`          | Returns raw SQL representation (`name` or `name AS alias`) |
| `RenderName()`   | Applies dialect-safe quoting to name                       |
| `RenderAlias()`  | Applies dialect-safe aliasing to a qualified expression    |
| `String()`       | Returns diagnostic view for logs and test outputs          |


---

## 🧠 Design Highlights

- **Centralized validation**: errors are caught during token construction
- **Aliasing control**: supports both inline and explicit `AS`
- **Safe rendering**: always checks for nil and empty input
- **Future-proof**: planned to include `TokenKind` metadata

---

## 🧪 Typical Usage (in Column)

```go
type Column struct {
    *BaseToken
    Qualified string
}

col := &Column{BaseToken: NewBaseToken("users.id AS uid")}
fmt.Println(col.Raw())           // → "users.id AS uid"
fmt.Println(col.IsAliased())     // → true
fmt.Println(col.RenderName(nil)) // → "users.id"
```

---

## 🧷 Notes

- **Conflicts** between inline alias and explicit alias are retained as errors.
- **IsValid** may be refactored into `Validate()` + `IsErrored()`.

---

## 📌 Status

Included in Entiqon v1.6.0 and continuously updated during the normalization plan.
Supports testing and rendering across all dialects through composable tokens.

---

### 📌 `String()` Method

Returns a human-readable diagnostic representation of the token. This is not intended for SQL rendering, but rather for test debugging and developer feedback.

If the token is `nil`, it returns the literal string `"nil"` — which **should not** be used in test assertions or string comparisons.

```go
func (b *BaseToken) String() string {
    if b == nil {
        return "nil" // defensive fallback — not assertable in tests
    }

    suffix := fmt.Sprintf("aliased: %v, errored: %v", b.IsAliased(), b.HasError())
    if b.HasError() {
        suffix += fmt.Sprintf(", error: %s", b.Error.Error())
    }

    return fmt.Sprintf("%s("%s") [%s]", b.kind.String(), b.Name, suffix)
}
```

#### Example

```go
b := NewBaseToken("users.id AS uid")
b.SetKind(ColumnKind)
fmt.Println(b.String())
// Output:
// Column("id") [aliased: true, errored: false]
```


---

### 📌 `GetKind()` Method

Returns the token kind (type classification) associated with the token instance.

```go
func (b *BaseToken) GetKind() Kind {
    if b == nil {
        return UnknownKind
    }
    return b.kind
}
```

#### Example

```go
b := NewBaseToken("id")
b.SetKind(ColumnKind)
fmt.Println(b.GetKind()) // Output: ColumnKind
```

---

### 📌 `SetKind()` Method

Sets the token kind, allowing the token to declare its identity in type-safe rendering or diagnostics.

```go
func (b *BaseToken) SetKind(k Kind) {
    b.kind = k
}
```

#### Example

```go
b := NewBaseToken("id")
b.SetKind(TableKind)
```
---

### 📌 `GetName()` Method

Returns the token name in a nil-safe way.

```go
func (b *BaseToken) GetName() string {
    if b == nil {
        return ""
    }
    return b.Name
}
```

---

### 📌 `GetSource()` Method

Returns the original expression used to create the token.

```go
func (b *BaseToken) GetSource() string {
    if b == nil {
        return ""
    }
    return b.Source
}
```

---

### 📌 `AliasOr()` Method

Returns the alias if available, or the name as fallback.

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

---

### 📌 `IsAliased()` Method

Indicates whether the token has a defined alias.

```go
func (b *BaseToken) IsAliased() bool {
    return b != nil && b.Alias != ""
}
```

---

### 📌 `IsErrored()` Method

Indicates whether the token contains a parsing or semantic error.

```go
func (b *BaseToken) IsErrored() bool {
    return b != nil && b.Error != nil
}
```

---

### 📌 `Validate()` Method (Planned)

Validates the token, resolving defaults and identifying semantic problems. Will be used in future automatic builders.

```go
func (b *BaseToken) Validate() *BaseToken
```

> Note: Implementation pending.

---

### 📌 `SetErrorWith()` Method

Assigns an error to the token and sets its `Source`.

```go
func (b *BaseToken) SetErrorWith(source string, err error) *BaseToken {
    b.Error = err
    if b.Source == "" {
        b.Source = source
    }
    return b
}
```

---

### 📌 `Raw()` Method

Returns a raw string representation for SQL output.

```go
func (b *BaseToken) Raw() string
```

Returns `name` or `name AS alias`, depending on aliasing.

---

### 📌 `RenderName()` Method

Quotes the name according to dialect.

```go
func (b *BaseToken) RenderName(d Dialect) string
```

---

### 📌 `RenderAlias()` Method

Quotes and renders alias from a qualified expression.

```go
func (b *BaseToken) RenderAlias(d Dialect, qualified string) string
```

---

### 📌 `String()` Method

Returns a diagnostic representation.

```go
func (b *BaseToken) String() string {
    if b == nil {
        return "nil"
    }

    suffix := fmt.Sprintf("aliased: %v, errored: %v", b.IsAliased(), b.HasError())
    if b.HasError() {
        suffix += fmt.Sprintf(", error: %s", b.Error.Error())
    }

    return fmt.Sprintf("%s("%s") [%s]", b.kind.String(), b.Name, suffix)
}
```

---

### 📌 `GetKind()` Method

Returns the kind (token type).

```go
func (b *BaseToken) GetKind() Kind {
    if b == nil {
        return UnknownKind
    }
    return b.kind
}
```

---

### 📌 `SetKind()` Method

Assigns the token kind.

```go
func (b *BaseToken) SetKind(k Kind) {
    b.kind = k
}
```
