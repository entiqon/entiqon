// Package contract defines small, reusable behavioral contracts (interfaces)
// that core tokens (Field, Table, Condition, etc.) and builders implement
// to enable polymorphic behavior without tight coupling between packages.
//
// Each contract describes one narrow aspect of behavior, allowing
// composable capabilities. By design, contracts are intentionally
// minimalistic and orthogonal:
//
//   - BaseToken: core identity (input, expression, alias, validation).
//   - Renderable: canonical SQL output, dialect-aware.
//   - Rawable: generic SQL fragment, dialect-agnostic.
//   - Stringable: human-facing representation for logs and audits.
//   - Debuggable: developer-facing diagnostic view of internal state.
//   - Clonable[T]: semantic clone for safe mutation.
//
// Separation of concerns:
//
//   - BaseToken → identity & validation
//   - Renderable → query generation
//   - Rawable → generic fragments
//   - Stringable → logs/audit
//   - Debuggable → developer diagnostics
//   - Clonable[T] → safe duplication
//
// Example:
//
//	package main
//
//	import (
//	    "fmt"
//	    "entiqon/db/contract"
//	)
//
//	// ExampleToken implements several contracts.
//	type ExampleToken struct {
//	    Name  string
//	    Alias string
//	}
//
//	func (e ExampleToken) Render() string   { return fmt.Sprintf("%s AS %s", e.Name, e.Alias) }
//	func (e ExampleToken) Raw() string      { return fmt.Sprintf("%s AS %s", e.Name, e.Alias) }
//	func (e ExampleToken) String() string   { return fmt.Sprintf("[Token name=%q alias=%q]", e.Name, e.Alias) }
//	func (e ExampleToken) Debug() string    { return fmt.Sprintf("ExampleToken{Name:%q, Alias:%q}", e.Name, e.Alias) }
//	func (e ExampleToken) Clone() *ExampleToken {
//	    return &ExampleToken{Name: e.Name, Alias: e.Alias}
//	}
//
//	func main() {
//	    var bt contract.BaseToken = ExampleToken{Name: "users", Alias: "u"}
//	    var r contract.Renderable = ExampleToken{Name: "users", Alias: "u"}
//	    var w contract.Rawable    = ExampleToken{Name: "users", Alias: "u"}
//	    var s contract.Stringable = ExampleToken{Name: "users", Alias: "u"}
//	    var d contract.Debuggable = ExampleToken{Name: "users", Alias: "u"}
//	    var c contract.Clonable[*ExampleToken] = ExampleToken{Name: "users", Alias: "u"}
//
//	    fmt.Println(bt.Input(), bt.Expr(), bt.Alias(), bt.IsAliased(), bt.IsValid())
//	    fmt.Println(r.Render()) // dialect-aware SQL
//	    fmt.Println(w.Raw())    // generic SQL
//	    fmt.Println(s.String()) // audit/log
//	    fmt.Println(d.Debug())  // developer diagnostic
//	    fmt.Println(c.Clone())  // safe copy
//	}
package contract
