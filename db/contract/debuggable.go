/**
 * @Author: Isidro Lopez isidro.lopezg@live.com
 * @Date: 2025-08-24 05:42:00
 * @LastEditors: Isidro Lopez isidro.lopezg@live.com
 * @LastEditTime: 2025-08-24 05:42:04
 * @FilePath: db/contract/debuggable.go
 * @Description: 这是默认设置,可以在设置》工具》File Description中进行配置
 */
// File: db/contract/debuggable.go
//
// Debuggable defines developer-facing diagnostic output. See package-level
// documentation in doc.go for an overview of all contracts and
// their distinct purposes.

package contract

// Debuggable defines the contract for objects that can expose their
// internal state for diagnostic or developer-facing output.
//
// Debug() is intended to be *more verbose* and *implementation-oriented*
// than String(), and is typically used when troubleshooting, inspecting
// runtime state, or writing developer-oriented logs.
//
// Key points:
//   - Debug() is not guaranteed to be stable across versions.
//   - Unlike Render(), the output is not meant for machine consumption.
//   - Unlike String(), the output is not designed for audits or end-user logs.
//
// Implementations should include enough internal detail to aid diagnosis,
// but are free to choose the most helpful format (often JSON-like or
// Go-syntax inspired).
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
//	// ExampleToken demonstrates Debuggable.
//	type ExampleToken struct {
//	    Name  string
//	    Alias string
//	    Err   error
//	}
//
//	func (t ExampleToken) Debug() string {
//	    return fmt.Sprintf("ExampleToken{Name:%q, Alias:%q, Err:%v}", t.Name, t.Alias, t.Err)
//	}
//
//	func main() {
//	    var d contract.Debuggable = ExampleToken{Name: "users", Alias: "u"}
//	    fmt.Println(d.Debug())
//	    // Output: ExampleToken{Name:"users", Alias:"u", Err:<nil>}
//	}
type Debuggable interface {
	// Debug returns a developer-facing representation of the object.
	// The output should include enough detail for diagnostics but is
	// not intended for machine parsing or audit logs.
	Debug() string
}
