// File: db/contract/validable.go

package contract

import (
	"errors"
)

var (
	// ErrUnsupportedType is returned by ValidateType when the input
	// is a token or another unsupported value instead of a raw string.
	ErrUnsupportedType = errors.New("unsupported type")
)

// Validable is implemented by tokens that can report whether
// they are structurally valid.
//
// # Semantics
//
//   - A token is considered valid when it has no internal error
//     and satisfies the minimal rules required by its type.
//   - Validation is lightweight and non-intrusive: it does not
//     perform deep semantic analysis or dialect resolution.
//   - For example, a Join is valid if both sides are present
//     and no construction error was set. A Field is valid if
//     its expression and alias were parsed correctly.
//
// # Usage
//
// This contract is intentionally narrow, allowing higher-level
// builders (e.g., SelectBuilder) to quickly determine whether
// a token may participate in query generation without needing
// to inspect its internal state.
//
// Tokens that do not require validation should not implement
// this interface.
//
// Example:
//
//	func validateToken(t contract.Validable) {
//	    if !t.IsValid() {
//	        fmt.Println("Token is not valid")
//	    }
//	}
type Validable interface {
	// IsValid reports whether the token is considered valid.
	//
	// It must return true only when the token has no error
	// and its essential invariants are satisfied.
	IsValid() bool
}
