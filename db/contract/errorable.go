// File: db/contract/errorable.go

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
type Errorable[T any] interface {
	// IsErrored reports whether the object was constructed
	// with an error or became invalid during initialization.
	IsErrored() bool

	// Error returns the underlying error if any.
	// Returns nil if the object is valid.
	Error() error

	// SetError marks the token as errored with the given error.
	// Implementations must store the error for reporting in Error().
	SetError(err error) T
}
