// File: db/driver/styling/contracts.go

// Package styling provides utilities and configuration types for rendering SQL tokens,
// such as quoting identifiers or formatting table/column aliases in dialect-specific ways.

package styling

// IdentifierQuoter defines the minimal quoting capability required for alias or token formatting.
//
// This interface is used to decouple the styling logic (e.g., FormatWith) from the full Dialect interface,
// avoiding circular dependencies between styling and dialect packages.
//
// Any type (such as BaseDialect) that implements:
//
//	func QuoteIdentifier(name string) string
//
// can satisfy this interface and be passed to FormatWith methods in styling utilities.
//
// Example:
//
//	func (a AliasStyle) FormatWith(q IdentifierQuoter, base, alias string) string {
//	    return fmt.Sprintf("%s AS %s", q.QuoteIdentifier(base), q.QuoteIdentifier(alias))
//	}
type IdentifierQuoter interface {
	// QuoteIdentifier returns the dialect-specific representation of an identifier,
	// applying the appropriate quoting strategy (e.g., "name", `name`, [name]).
	QuoteIdentifier(name string) string
}
