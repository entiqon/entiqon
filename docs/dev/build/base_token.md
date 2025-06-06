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

| Method             | Description                                                      |
|--------------------|------------------------------------------------------------------|
| `GetName()`        | Returns the token name safely (even if nil)                      |
| `GetSource()`      | Returns the original source string safely                        |
| `AliasOr()`        | Returns alias or name fallback                                   |
| `IsAliased()`      | Returns true if alias is defined                                 |
| `IsErrored()`      | Returns true if error is set                                     |
| `Validate()`       | (planned) Assigns error if malformed or incomplete               |
| `SetErrorWith()`   | Manually set error and associate with source                     |
| `Raw()`            | Returns raw SQL representation (`name` or `name AS alias`)       |
| `RenderName()`     | Applies dialect-safe quoting to name                             |
| `RenderAlias()`    | Applies dialect-safe aliasing to a qualified expression          |
| `String(kind)`     | Returns diagnostic view for logs and test outputs                |

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