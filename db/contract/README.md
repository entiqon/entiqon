# Contract Package

The `contract` package defines small, reusable behavioral contracts (interfaces)
that core tokens (`Field`, `Table`, `Condition`, etc.) and builders implement to
enable polymorphic behavior without tight coupling between packages.

Contracts are intentionally minimalistic and orthogonal. Each one describes a
narrow capability that can be composed with others.

---

## Available Contracts

### [Renderable](./renderable.go)
- **Purpose**: Canonical, dialect-aware SQL output (machine-facing).
- **Method**: `Render() string`
- **Usage**: Used by builders to generate final SQL fragments.

### [Rawable](./rawable.go)
- **Purpose**: Generic SQL fragments, dialect-agnostic.
- **Methods**:  
  - `Raw() string`  
  - `IsRaw() bool`
- **Usage**: Useful for introspection, debugging, and as a base for dialect rewriting.

### [Stringable](./stringable.go)
- **Purpose**: Human-facing audit/log output.
- **Method**: `String() string`
- **Usage**: Concise log/audit representation of tokens.

### [Debuggable](./debuggable.go)
- **Purpose**: Developer-facing diagnostic output.
- **Method**: `Debug() string`
- **Usage**: Verbose inspection with internal state flags and error state.

### [Clonable](./clonable.go)
- **Purpose**: Semantic cloning for safe mutation.
- **Method**: `Clone() T`
- **Usage**: Duplicate tokens/builders without sharing state.

### [Errorable](./errorable.go)
- **Purpose**: Error state inspection for tokens/builders.
- **Methods**:  
  - `IsErrored() bool`  
  - `Error() error`
- **Usage**: Detect invalid tokens and retrieve construction errors.

---

## Examples

See [example_test.go](./example_test.go) for runnable examples of all contracts:

```go
t := table.New("users", "u")

var r contract.Renderable = t
fmt.Println(r.Render())
// Output: users AS u

var e contract.Errorable = table.New("users AS")
fmt.Println(e.IsErrored()) // true
fmt.Println(e.Error())     // invalid format "users AS"
```

---

## Philosophy

- **Never panic**: Constructors always return a token, even if errored.
- **Auditability**: `input` is always preserved for logs.
- **Separation of concerns**:
  - Render → SQL generation
  - Raw → Generic fragments
  - String → Audit logs
  - Debug → Developer diagnostics
  - Clonable → Safe duplication
  - Errorable → Error handling

---

This package underpins the entire query builder layer. Contracts ensure tokens
like `Field` and `Table` behave consistently and predictably across the system.
