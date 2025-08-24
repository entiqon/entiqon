/**
 * @Author: Isidro Lopez isidro.lopezg@live.com
 * @Date: 2025-08-24 05:42:00
 * @LastEditors: Isidro Lopez isidro.lopezg@live.com
 * @LastEditTime: 2025-08-24 05:42:04
 * @FilePath: db/contract/errorable.go
 * @Description: 这是默认设置,可以在设置》工具》File Description中进行配置
 */
// File: db/contract/errorable.go
//
// Errorable defines error state inspection for tokens and builders.
// See package-level documentation in doc.go for an overview of all
// contracts and their distinct purposes.

package contract

// Errorable defines the contract for objects that may carry an
// internal error state from construction or evaluation.
//
// Errorable provides a consistent way to check whether a token or
// builder is usable (IsErrored) and to retrieve the underlying
// error (Error) if construction or validation failed.
//
// Contrast with:
//   - Renderable: produces canonical SQL output (machine-facing).
//   - Rawable: produces generic SQL fragments (dialect-agnostic).
//   - Stringable: produces audit/log output (human-facing).
//   - Debuggable: produces developer diagnostics.
//   - Clonable[T]: produces semantic copies for safe mutation.
//
// Example:
//
//	u := table.New("") // invalid construction
//	if u.IsErrored() {
//	    fmt.Println("Error:", u.Error())
//	    // Output: Error: table.New: empty table name
//	}
type Errorable interface {
	// IsErrored reports whether the object was constructed
	// with an error or became invalid during initialization.
	IsErrored() bool

	// Error returns the underlying error if any.
	// Returns nil if the object is valid.
	Error() error
}
