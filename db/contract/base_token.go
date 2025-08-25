// File: db/contract/base_token.go
//
// BaseToken defines semantic clone behavior for all SQL tokens
// (Field, Table, etc.) to provide consistent access to their
// original input, normalized expression, alias, and validation
// status. See package-level documentation in doc.go for an
// overview of all contracts and their distinct purposes.

package contract

// BaseToken is implemented by all SQL tokens that carry an expression
// and optional alias. It defines common behaviors for identity and validation.
type BaseToken interface {
	// Input returns the raw user-provided string(s)
	// before parsing/normalization. Useful for auditing and logs.
	Input() string

	// Expr returns the parsed/normalized SQL expression
	// (e.g. "u.id" or "COUNT(*)").
	Expr() string

	// Alias returns the alias name, if present.
	Alias() string

	// IsAliased reports whether the token has an alias.
	IsAliased() bool

	// IsValid reports whether the token was successfully parsed
	// and is valid for rendering.
	IsValid() bool
}
