package token

// BaseToken defines the minimal behavior of any SQL token used internally
// during query construction.
//
// It enables the builder pipeline to process and validate raw SQL components
// such as columns, conditions, and values in a uniform manner.
//
// Implementations must fulfill:
//
//   - String(): returns a log/debug-safe, human-readable representation.
//   - IsValid(): performs basic integrity checks.
//   - Raw(): returns the raw SQL-compatible representation (not dialect-aware).
type BaseToken interface {
	// IsValid returns whether the token can be used in a valid SQL statement
	IsValid() bool

	// Raw returns the core SQL expression (unquoted, dialect-neutral)
	Raw() string

	// String returns a representation suitable for debugging/logging
	String() string
}
