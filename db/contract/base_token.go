// File: db/contract/base_token.go

package contract

import "github.com/entiqon/entiqon/db/token/types/identifier"

// BaseToken is implemented by all SQL tokens that carry an expression
// and optional alias. It defines common behaviors for identity and validation.
type BaseToken interface {
	// ExpressionKind returns the semantic classification of the
	// underlying expression (Identifier, Literal, Subquery, etc.).
	ExpressionKind() identifier.Type

	Identifiable

	// Alias returns the alias name, if present.
	Alias() string

	// IsAliased reports whether the token has an alias.
	IsAliased() bool
}
