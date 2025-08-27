# 📜 Contract

## 🧩 Overview

The `contract` package defines small, reusable behavioral contracts (interfaces)
that core tokens (`Field`, `Table`, `Join`, `Condition`, etc.) and builders implement
to enable polymorphic behavior without tight coupling between packages.

Contracts are intentionally minimalistic and orthogonal. Each one describes a
narrow capability that can be composed with others.

---

## Available Contracts (strict order)

### [BaseToken](./base_token.go)
- **Purpose**: Core identity for tokens.
- **Methods**:
  - `Input() string`
  - `Expr() string`
  - `Alias() string`
  - `IsAliased() bool`
- **Usage**: Ensures tokens like `Field` and `Table` consistently expose their
  raw input, normalized expression, and alias.

### [Clonable](./clonable.go)
- **Purpose**: Semantic cloning for safe mutation.
- **Method**: `Clone() T`

### [Debuggable](./debuggable.go)
- **Purpose**: Developer-facing diagnostic output.
- **Method**: `Debug() string`

### [Errorable](./errorable.go)
- **Purpose**: Error state inspection and propagation for tokens/builders.
- **Methods**:
  - `IsErrored() bool`
  - `Error() error`
  - `SetError(err error) T`

### [Rawable](./rawable.go)
- **Purpose**: Generic SQL fragments, dialect-agnostic.
- **Methods**:
  - `Raw() string`
  - `IsRaw() bool`

### [Renderable](./renderable.go)
- **Purpose**: Canonical, dialect-aware SQL output (machine-facing).
- **Method**: `Render() string`

### [Stringable](./stringable.go)
- **Purpose**: Human-facing audit/log output.
- **Method**: `String() string`

### [Validable](./validable.go)
- **Purpose**: Structural validation.
- **Method**: `IsValid() bool`

---

## Examples

See [example_test.go](./example_test.go) for runnable examples of all contracts:

```go
t := table.New("users", "u")

var bt contract.BaseToken = t
fmt.Println(bt.Input(), bt.Expr(), bt.Alias(), bt.IsAliased())
// Output: users users u true

var c contract.Clonable[*table.Table] = t
fmt.Println(c.Clone()) // safe copy

var d contract.Debuggable = t
fmt.Println(d.Debug()) // developer diagnostic

var e contract.Errorable[*table.Table] = t
fmt.Println(e.IsErrored(), e.Error()) // error state

var w contract.Rawable = t
fmt.Println(w.Raw()) // generic SQL

var r contract.Renderable = t
fmt.Println(r.Render()) // dialect-aware SQL

var s contract.Stringable = t
fmt.Println(s.String()) // audit/log

var v contract.Validable = t
fmt.Println(v.IsValid()) // structural validation
```

---

## Philosophy

- **Never panic**: Constructors always return a token, even if errored.
- **Auditability**: `Input()` is always preserved for logs.
- **Consistency**: All tokens share a common identity contract (`BaseToken`).
- **Separation of concerns**:
  - BaseToken → identity
  - Clonable → safe duplication
  - Debuggable → developer diagnostics
  - Errorable → error handling
  - Rawable → generic fragments
  - Renderable → SQL generation
  - Stringable → logs/audit
  - Validable → validation

---

This package underpins the entire query builder layer. Contracts ensure tokens
like `Field`, `Table`, and `Join` behave consistently and predictably across the system.

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project
