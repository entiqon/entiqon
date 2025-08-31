> Part of [Entiqon](https://github.com/entiqon/entiqon) / [Database](../)

# Contract

## 🧩 Overview

The `contract` package defines small, reusable behavioral contracts (interfaces)  
that core tokens (`Field`, `Table`, `Join`, `Condition`, etc.) and builders implement  
to enable polymorphic behavior without tight coupling between packages.

Contracts are intentionally **minimalistic, composable, and auditable**.  
They focus on clear identity, safe mutation, and consistent rendering across tokens,  
while leaving implementation details to each package.

---

## Available Contracts (strict order)

| Contract                          | Purpose                                                    | Methods                                                                       |
|-----------------------------------|------------------------------------------------------------|-------------------------------------------------------------------------------|
| [Identifiable](./identifiable.go) | Raw input and normalized expression identity (alias-free). | `Input() string`<br>`Expr() string`                                           |
| [Aliasable](./aliasable.go)       | Alias surface for tokens.                                  | `Alias() string`<br>`IsAliased() bool`                                        |
| [BaseToken](./base_token.go)      | Core identity for tokens, including alias.                 | `Input() string`<br>`Expr() string`<br>`Alias() string`<br>`IsAliased() bool` |
| [Clonable](./clonable.go)         | Semantic cloning for safe mutation.                        | `Clone() T`                                                                   |
| [Debuggable](./debuggable.go)     | Developer-facing diagnostic output.                        | `Debug() string`                                                              |
| [Errorable](./errorable.go)       | Error state inspection and propagation.                    | `IsErrored() bool`<br>`Error() error`<br>`SetError(err error) T`              |
| [Rawable](./rawable.go)           | Generic SQL fragments, dialect-agnostic.                   | `Raw() string`<br>`IsRaw() bool`                                              |
| [Renderable](./renderable.go)     | Canonical, dialect-aware SQL output.                       | `Render() string`                                                             |
| [Stringable](./stringable.go)     | Human-facing audit/log output.                             | `String() string`                                                             |
| [Validable](./validable.go)       | Structural validation.                                     | `IsValid() bool`                                                              |

---

## Examples

See [example_test.go](./example_test.go) for runnable examples of all contracts:

```go
t := table.New("users", "u")

var k contract.Kindable[myKind] = &myToken{}
k.SetKind(Custom)
fmt.Println("Kind=", k.Kind()) // classification

var id contract.Identifiable = t
fmt.Println(id.Input(), id.Expr()) // identity only

var al contract.Aliasable = t
fmt.Println(al.Alias(), al.IsAliased()) // alias surface

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

_**Note**:_  
- Use `Kindable` for classification enums (e.g., `condition.Type`, `identifier.Type`).  
- Use `Identifiable` when aliasing must be excluded (e.g., in `Condition` tokens).


---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project