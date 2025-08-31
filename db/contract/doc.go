// Package contract defines small, reusable behavioral contracts (interfaces)
// that core tokens (Field, Table, Join, Condition, etc.) and builders implement
// to enable polymorphic behavior without tight coupling between packages.
//
// Each contract describes one narrow aspect of behavior, allowing
// composable capabilities. By design, contracts are intentionally
// minimalistic and orthogonal:
//
//   - Identifiable: raw input and normalized expression identity
//   - Aliasable: alias surface (Alias, IsAliased)
//   - BaseToken: core identity (input, expression, alias)
//   - Clonable[T]: semantic clone for safe mutation
//   - Debuggable: developer-facing diagnostic view of internal state
//   - Errorable[T]: error state management and reporting
//   - Rawable: generic SQL fragment, dialect-agnostic
//   - Renderable: canonical SQL output, dialect-aware
//   - Stringable: human-facing representation for logs and audits
//   - Validable: structural validation
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
//	    err   error
//	}
//
//	func (e ExampleToken) Clone() *ExampleToken {
//	    return &ExampleToken{Name: e.Name, Alias: e.Alias}
//	}
//	func (e ExampleToken) Debug() string                     { return fmt.Sprintf("ExampleToken{Name:%q, Alias:%q}", e.Name, e.Alias) }
//	func (e ExampleToken) IsErrored() bool                   { return e.err != nil }
//	func (e ExampleToken) Error() error                      { return e.err }
//	func (e *ExampleToken) SetError(err error) *ExampleToken { e.err = err; return e }
//	func (e ExampleToken) Raw() string                       { return fmt.Sprintf("%s AS %s", e.Name, e.Alias) }
//	func (e ExampleToken) Render() string                    { return fmt.Sprintf("%s AS %s", e.Name, e.Alias) }
//	func (e ExampleToken) String() string                    { return fmt.Sprintf("[Token name=%q alias=%q]", e.Name, e.Alias) }
//	func (e ExampleToken) IsValid() bool                     { return e.err == nil }
//	func (e ExampleToken) Input() string                     { return e.Name }
//	func (e ExampleToken) Expr() string                      { return e.Name }
//
//	func main() {
//	    var id contract.Identifiable = ExampleToken{Name: "users", Alias: "u"}
//	    var al contract.Aliasable   = ExampleToken{Name: "users", Alias: "u"}
//	    var bt contract.BaseToken    = ExampleToken{Name: "users", Alias: "u"}
//	    var c  contract.Clonable[*ExampleToken] = ExampleToken{Name: "users", Alias: "u"}
//	    var d  contract.Debuggable   = ExampleToken{Name: "users", Alias: "u"}
//	    var e  contract.Errorable[*ExampleToken] = &ExampleToken{Name: "users", Alias: "u"}
//	    var w  contract.Rawable      = ExampleToken{Name: "users", Alias: "u"}
//	    var r  contract.Renderable   = ExampleToken{Name: "users", Alias: "u"}
//	    var s  contract.Stringable   = ExampleToken{Name: "users", Alias: "u"}
//	    var v  contract.Validable    = ExampleToken{Name: "users", Alias: "u"}
//
//	    fmt.Println(id.Input(), id.Expr())   // identity only
//	    fmt.Println(al.Alias(), al.IsAliased()) // alias surface
//	    fmt.Println(bt.Input(), bt.Expr(), bt.Alias(), bt.IsAliased())
//	    fmt.Println(c.Clone())                         // safe copy
//	    fmt.Println(d.Debug())                         // developer diagnostic
//	    fmt.Println(e.IsErrored(), e.Error())          // error state
//	    fmt.Println(w.Raw())                           // generic SQL
//	    fmt.Println(r.Render())                        // dialect-aware SQL
//	    fmt.Println(s.String())                        // audit/log
//	    fmt.Println(v.IsValid())                       // structural validation
//	}
package contract
