// filename: db/internal/build/token/generic.go

package token

// GenericToken defines the minimal contract for any structured token
// used in SQL query construction and validation.
//
// This interface is implemented by types like Column, Table, Join, and Condition.
// It allows a query builder to inspect and resolve token properties in a uniform way.
//
// Implementations are expected to carry at least a name (or identifier),
// an optional alias, and a mechanism for error detection.
//
// Tokens must be usable in SQL generation after passing IsValid() checks.
// If an internal Error exists, HasError() will return true, and the builder
// may choose to skip, replace, or report the token accordingly.
//
// The Raw method should return the SQL-compatible identifier (e.g., "users.id")
// without quoting or dialect-specific formatting.
//
// The String method is used strictly for debugging/logging and should not
// be used in actual SQL output.
//
// # Examples
//
// The following logic may be used in a SQL SELECT builder:
//
//	for _, tok := range selectedTokens {
//	    if !tok.IsValid() {
//	        log.Printf("invalid token: %s", tok.String())
//	        continue
//	    }
//	    clause.Write(tok.Raw())
//	}
//
// A valid Column token might produce:
//
//	Raw():     "users.name"
//	String():  Column("users.name") [aliased: false, qualified: true]
type GenericToken interface {
	// HasError reports whether the token has encountered a validation,
	// parsing, or resolution error.
	//
	// This method does not imply structural validity; use IsValid() instead
	// to determine whether the token should be used in query output.
	HasError() bool

	// IsValid reports whether the token is safe to use in query generation.
	//
	// A token is considered valid if it has a resolvable identifier
	// (e.g., non-empty Name) and no internal Error state.
	IsValid() bool
}
