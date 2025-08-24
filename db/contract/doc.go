/**
 * @Author: Isidro Lopez isidro.lopezg@live.com
 * @Date: 2025-08-24 05:42:00
 * @LastEditors: Isidro Lopez isidro.lopezg@live.com
 * @LastEditTime: 2025-08-24 05:42:04
 * @FilePath: db/contract/doc.go
 * @Description: 这是默认设置,可以在设置》工具》File Description中进行配置
 */
// Package contract defines small, reusable behavioral contracts (interfaces)
// that core tokens (Field, Table, Condition, etc.) and builders implement
// to enable polymorphic behavior without tight coupling between packages.
//
// Each contract describes one narrow aspect of behavior, allowing
// composable capabilities. By design, contracts are intentionally
// minimalistic and orthogonal:
//
//   - Renderable: canonical SQL output, dialect-aware.
//   - Rawable: generic SQL fragment, dialect-agnostic.
//   - Stringable: human-facing representation for logs and audits.
//   - Debuggable: developer-facing diagnostic view of internal state.
//   - Clonable[T]: semantic clone contract for safe mutation.
//
// These contracts separate concerns:
//
//   - Renderable is for machine-facing query generation.
//   - Rawable is for generic SQL fragments used as the basis for dialects.
//   - Stringable is for audit/logging, human-facing but concise.
//   - Debuggable is for diagnostics, verbose, not stable across versions.
//   - Clonable[T] is for safe duplication of tokens or builders.
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
//	    var r contract.Renderable = ExampleToken{Name: "users", Alias: "u"}
//	    var w contract.Rawable    = ExampleToken{Name: "users", Alias: "u"}
//	    var s contract.Stringable = ExampleToken{Name: "users", Alias: "u"}
//	    var d contract.Debuggable = ExampleToken{Name: "users", Alias: "u"}
//	    var c contract.Clonable[*ExampleToken] = ExampleToken{Name: "users", Alias: "u"}
//
//	    fmt.Println(r.Render()) // dialect-aware SQL
//	    fmt.Println(w.Raw())    // generic SQL
//	    fmt.Println(s.String()) // audit/log
//	    fmt.Println(d.Debug())  // developer diagnostic
//	    fmt.Println(c.Clone())  // safe copy
//	}
package contract
