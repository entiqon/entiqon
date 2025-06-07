# Entiqon Developer Guide: BaseToken

`BaseToken` is the foundational building block for SQL-like token structures (e.g., `Column`, `Table`, `Join`, `Condition`). It normalizes raw input expressions into a parsed identifier and alias, tracks the original input for diagnostics, enforces validation, and records a `Kind` classification. Higher-level tokens embed `BaseToken` to reuse this logic, leaving context-specific resolution (e.g., table qualification) to their own constructors.

> **Note:** `BaseToken` is an internal abstraction. Use its public constructors and interfaces, not direct field access.

---

## 🔍 Purpose & Scope

- **Parse & normalize** raw input (e.g., `"users.id AS u"`) into:
  - `input` – original string passed to the parser
  - `name` – parsed identifier (e.g., `"users.id"`)
  - `alias` – parsed alias (e.g., `"u"`)
  - `Error` – non-nil if validation or alias conflict occurred
  - `kind` – classification (`contract.Kind`)
- **Provide accessors**:
  - `GetInput()` • original input
  - `GetName()` • unquoted identifier
  - `GetAlias()` • unquoted alias
  - `GetRaw()` • reconstructed SQL fragment: `name` or `name AS alias`
- **Implement core interfaces**:
  - `contract.Rawable` • raw data getters
  - `contract.Renderable` • `RenderName`, `RenderAlias`, `String`
  - `contract.Errorable` • `GetError`, `IsErrored`, `SetError`
  - `contract.Kindable` • `SetKind`, `GetKind`
- **Support mutation**:
  - `SetName()` • override identifier
  - `SetAlias()` • override alias

---

## 🧱 Struct Definition

```go
// BaseToken normalizes a raw SQL token for later rendering or diagnostics.
type BaseToken struct {
    input string        // original raw input
    name  string        // parsed identifier
    alias string        // parsed alias
    Error error         // non-nil if validation failed
    kind  contract.Kind // token classification
}
```

> Fields `name` and `alias` are unexported; use `GetName()` and `GetAlias()`.  
> Direct field access is deprecated and may be removed in future.

---

## 🚧 Constructor

### `func NewBaseToken(input string, explicitAlias ...string) *BaseToken`

1. **Trim & validate** `input`: non-empty, no commas, not starting with `"AS "` or equaling `"AS"`.  
2. **Parse inline alias** via `ParseAlias(input)` → `(parsedName, parsedAlias)`.  
3. **Apply `explicitAlias`** if provided:
   - If both inline and explicit present and differ → record error.
   - Override alias with explicit.  
4. **Store** `b.input`, `b.name`, `b.alias`.  
5. Return `*BaseToken` with `Error` set if any validation failed.

```go
b := NewBaseToken("users.id AS u")
// b.GetInput() → "users.id AS u"
// b.GetName()  → "users.id"
// b.GetAlias() → "u"
// b.GetError() → nil
```

---

## ✅ Methods Overview

| Method                  | Description                                                                                     |
|-------------------------|-------------------------------------------------------------------------------------------------|
| `GetInput()`            | Original input string (nil-safe).                                                               |
| `GetName()`             | Parsed identifier, unquoted (nil-safe).                                                        |
| `GetAlias()`            | Parsed alias, unquoted (nil-safe).                                                             |
| `AliasOr()`             | `GetAlias()` if non-empty, else `GetName()`.                                                   |
| `GetRaw()`              | `name` or `name AS alias` (nil-safe).                                                          |
| **Errorable**           |                                                                                                 |
| `IsErrored()`           | `true` if an error was recorded (nil-safe).                                                     |
| `GetError()`            | Recorded `error` or `nil` (nil-safe).                                                           |
| `SetError(input, err)`  | Record an error and update `input` for diagnostics (nil-safe).                                  |
| **Mutation**            |                                                                                                 |
| `SetName(name)`         | Override parsed name (nil-safe).                                                                |
| `SetAlias(alias)`       | Override parsed alias (nil-safe).                                                              |
| **Kindable**            |                                                                                                 |
| `SetKind(kind)`         | Set token classification (nil-safe).                                                            |
| `GetKind()`             | Retrieve token classification.                                                                  |
| **Renderable**          |                                                                                                 |
| `RenderName(q Quoter)`  | Quote identifier (alias if present) via `Quoter`; unquoted if `q == nil`.                       |
| `RenderAlias(q, qual)`  | Quote and append alias to `qual` if present; fallback if missing.                               |
| `String()`              | Diagnostic string: `Kind("name") [aliased:bool, errored:bool, error:<msg>]`.                   |

---

## 📖 Detailed Method Examples

### GetRaw

```go
func (b *BaseToken) GetRaw() string {
    if b == nil { return "" }
    if b.alias == "" { return b.name }
    return fmt.Sprintf("%s AS %s", b.name, b.alias)
}
```

<aside>
// **Example**

// Raw form without alias:
//   b := NewBaseToken("users.id")
//   fmt.Println(b.GetRaw()) // Output: users.id

// With alias:
//   b := NewBaseToken("users.id AS u")
//   fmt.Println(b.GetRaw()) // Output: users.id AS u
</aside>

### SetName

```go
func (b *BaseToken) SetName(name string) {
    if b == nil { return }
    b.name = name
}
```

<aside>
// **Example**

// Override parsed name:
//   b := NewBaseToken("users.email AS e")
//   b.SetName("email")
//   fmt.Println(b.GetName()) // Output: email
</aside>

### SetAlias

```go
func (b *BaseToken) SetAlias(alias string) {
    if b == nil { return }
    b.alias = alias
}
```

<aside>
// **Example**

// Override parsed alias:
//   b := NewBaseToken("users.id AS u")
//   b.SetAlias("uid")
//   fmt.Println(b.GetAlias()) // Output: uid
</aside>

---

2025 — **© Entiqon Project**
