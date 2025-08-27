// File: db/contract/base_token.go

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
}
