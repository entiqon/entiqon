// File: db/driver/styling/placeholder_style.go

package styling

import "fmt"

// PlaceholderStyle defines how SQL parameter placeholders are rendered
// in parameterized queries. Each style corresponds to a different SQL dialect.
//
// It determines how values are substituted during query preparation,
// e.g., using "?", "$1", or ":name" depending on the engine.
//
// The zero value is PlaceholderUnset, which indicates that the style
// is not configured and must be explicitly set.
//
// # Example:
//
//	switch dialect.Placeholder {
//	case PlaceholderQuestion:
//	    return "?"
//
//	case PlaceholderDollar:
//	    return fmt.Sprintf("$%d", index+1)
//
//	case PlaceholderNamed:
//	    return fmt.Sprintf(":%s", name)
//
//	case PlaceholderAt:
//	    return fmt.Sprintf("@%s", name)
//	}
type PlaceholderStyle int

const (
	// PlaceholderUnset indicates the placeholder style is not configured.
	// This is the default zero value and is invalid during dialect validation.
	PlaceholderUnset PlaceholderStyle = iota

	// PlaceholderQuestion uses "?" — unnumbered positional.
	// Common in MySQL, SQLite, and other lightweight engines.
	PlaceholderQuestion

	// PlaceholderDollar uses "$1", "$2", etc. — numbered positional.
	// Used in PostgreSQL.
	PlaceholderDollar

	// PlaceholderNamed uses ":name" — named parameters.
	// Typical for Oracle, older DB2, and various ORM interfaces.
	PlaceholderNamed

	// PlaceholderAt uses "@name" — alternate named parameter style.
	// Used in SQL Server, Sybase, and some ADO-based engines.
	PlaceholderAt
)

// Format returns a placeholder string based on the given positional index.
// This applies only to positional styles (Question, Dollar).
//
// Example:
//
//	PlaceholderQuestion.Format(1) → "?"
//	PlaceholderDollar.Format(3)   → "$3"
func (p PlaceholderStyle) Format(index int) string {
	switch p {
	case PlaceholderQuestion:
		return "?"
	case PlaceholderDollar:
		return fmt.Sprintf("$%d", index)
	default:
		return "?"
	}
}

// FormatNamed returns a placeholder for named styles like ":name" or "@name".
// Positional styles will default to "?".
//
// Example:
//
//	PlaceholderNamed.FormatNamed("id") → ":id"
//	PlaceholderAt.FormatNamed("uid")   → "@uid"
func (p PlaceholderStyle) FormatNamed(name string) string {
	switch p {
	case PlaceholderNamed:
		return fmt.Sprintf(":%s", name)
	case PlaceholderAt:
		return fmt.Sprintf("@%s", name)
	default:
		return "?"
	}
}

// IsValid returns true if the PlaceholderStyle is a defined value.
// Valid values include:
//   - PlaceholderQuestion: "?"
//   - PlaceholderDollar:   "$1"
//   - PlaceholderNamed:    ":name"
//   - PlaceholderAt:       "@name"
//
// Used in BaseDialect.Validate to ensure placeholder formatting is configured.
func (p PlaceholderStyle) IsValid() bool {
	return p > PlaceholderUnset && p <= PlaceholderAt
}
