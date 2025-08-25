<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true" align="left" height="82" width="82" alt="entiqon"> Core Contract
</h1>
<h6 align="left">Part of the <a href="../../README.md">Entiqon</a> / <a href="../README.md">Database</a> toolkit.</h6>

## 🧩 Overview

The `contract` package defines small, reusable behavioral contracts (interfaces)
that core tokens (`Field`, `Table`, `Condition`, etc.) and builders implement to
enable polymorphic behavior without tight coupling between packages.

Contracts are intentionally minimalistic and orthogonal. Each one describes a
narrow capability that can be composed with others.

---

## Available Contracts

### [BaseToken](./base_token.go)
- **Purpose**: Core identity and validation for tokens.
- **Methods**:
  - `Input() string`
  - `Expr() string`
  - `Alias() string`
  - `IsAliased() bool`
  - `IsValid() bool`
- **Usage**: Ensures tokens like `Field` and `Table` consistently expose their
  raw input, normalized expression, alias, and validity.

### [Renderable](./renderable.go)
- **Purpose**: Canonical, dialect-aware SQL output (machine-facing).
- **Method**: `Render() string`

### [Rawable](./rawable.go)
- **Purpose**: Generic SQL fragments, dialect-agnostic.
- **Methods**:
  - `Raw() string`
  - `IsRaw() bool`

### [Stringable](./stringable.go)
- **Purpose**: Human-facing audit/log output.
- **Method**: `String() string`

### [Debuggable](./debuggable.go)
- **Purpose**: Developer-facing diagnostic output.
- **Method**: `Debug() string`

### [Clonable](./clonable.go)
- **Purpose**: Semantic cloning for safe mutation.
- **Method**: `Clone() T`

### [Errorable](./errorable.go)
- **Purpose**: Error state inspection and propagation for tokens/builders.
- **Methods**:
  - `IsErrored() bool`
  - `Error() error`
  - `SetError(err error)`

---

## Examples

See [example_test.go](./example_test.go) for runnable examples of all contracts:

```go
t := table.New("users", "u")

var bt contract.BaseToken = t
fmt.Println(bt.Input(), bt.Expr(), bt.Alias(), bt.IsAliased(), bt.IsValid())
// Output: users users u true true

var r contract.Renderable = t
fmt.Println(r.Render())
// Output: users AS u

var e contract.Errorable = table.New("users AS")
fmt.Println(e.IsErrored()) // true
fmt.Println(e.Error())     // invalid format "users AS"

// Example of SetError usage
tok := table.New("products")
et := tok.(*table.Table)
et.SetError(fmt.Errorf("manual mark as errored"))
fmt.Println(et.IsErrored()) // true
```

---

## Philosophy

- **Never panic**: Constructors always return a token, even if errored.
- **Auditability**: `Input()` is always preserved for logs.
- **Consistency**: All tokens share a common identity contract (`BaseToken`).
- **Separation of concerns**:
  - BaseToken → identity & validation
  - Renderable → SQL generation
  - Rawable → generic fragments
  - Stringable → logs/audit
  - Debuggable → developer diagnostics
  - Clonable → safe duplication
  - Errorable → error handling

---

This package underpins the entire query builder layer. Contracts ensure tokens
like `Field` and `Table` behave consistently and predictably across the system.

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project
