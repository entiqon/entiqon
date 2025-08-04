// File: db/driver/styling/alias_style.go

package styling

import "fmt"

// AliasStyle defines how SQL aliases are rendered in dialect-specific implementations.
type AliasStyle int

const (
	// AliasNone disables aliasing entirely.
	AliasNone AliasStyle = iota

	// AliasWithoutKeyword renders alias without the AS keyword (e.g., "users u").
	AliasWithoutKeyword

	// AliasWithKeyword renders alias using the AS keyword (e.g., "users AS u").
	AliasWithKeyword
)

// Format returns the SQL alias representation for a given base identifier and alias,
// based on the AliasStyle.
//
// If alias is empty or the style is AliasNone, only the base is returned.
// Example:
//
//	AliasWithKeyword.Format("users", "u") → "users AS u"
//	AliasWithoutKeyword.Format("users", "u") → "users u"
func (a AliasStyle) Format(base, alias string) string {
	if alias == "" || a == AliasNone {
		return base
	}
	switch a {
	case AliasWithKeyword:
		return fmt.Sprintf("%s AS %s", base, alias)
	case AliasWithoutKeyword:
		return fmt.Sprintf("%s %s", base, alias)
	default:
		return base
	}
}

// FormatWith renders the alias expression using the provided IdentifierQuoter.
// It applies the alias style by quoting both the base and alias using q.QuoteIdentifier(),
// and formatting them according to the AliasStyle.
//
// If alias is empty or AliasStyle is AliasNone, only the quoted base is returned.
// If q is nil, identifiers are returned unquoted.
//
// Example:
//
//	var style AliasStyle = AliasWithKeyword
//	var q IdentifierQuoter = BaseDialect{
//	    QuoteStyle: QuoteDouble,
//	}
//
//	style.FormatWith(q, "users", "u")
//	// Output: `"users" AS "u"`
//
//	style = AliasWithoutKeyword
//	style.FormatWith(q, "users", "u")
//	// Output: `"users" "u"`
//
//	style = AliasNone
//	style.FormatWith(q, "users", "u")
//	// Output: `"users"`
//
// Since: v1.5.0
func (a AliasStyle) FormatWith(q IdentifierQuoter, base, alias string) string {
	quote := func(s string) string {
		if q != nil {
			return q.QuoteIdentifier(s)
		}
		return s
	}

	if alias == "" || a == AliasNone {
		return quote(base)
	}

	switch a {
	case AliasWithKeyword:
		return fmt.Sprintf("%s AS %s", quote(base), quote(alias))
	case AliasWithoutKeyword:
		return fmt.Sprintf("%s %s", quote(base), quote(alias))
	default:
		return quote(base)
	}
}

// IsValid returns true if the AliasStyle is a recognized alias rendering option.
// Valid values include:
//   - AliasNone:              no aliasing
//   - AliasWithoutKeyword:    renders as "table alias"
//   - AliasWithKeyword:       renders as "table AS alias"
//
// Used in BaseDialect.Validate to ensure table/column aliasing is properly configured.
func (a AliasStyle) IsValid() bool {
	return a >= AliasNone && a <= AliasWithKeyword
}
