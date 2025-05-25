package styling

import "fmt"

// PlaceholderStyle defines how SQL placeholders are rendered
// in parameterized queries (e.g., ?, $1, :name).
type PlaceholderStyle int

const (
	// PlaceholderQuestion uses "?" — unnumbered positional (e.g., MySQL, SQLite).
	PlaceholderQuestion PlaceholderStyle = iota

	// PlaceholderDollar uses "$1", "$2", ... — numbered positional (e.g., PostgreSQL).
	PlaceholderDollar

	// PlaceholderNamed uses ":name" — named style (e.g., Oracle).
	PlaceholderNamed

	// PlaceholderAt uses "@name" — named style (e.g., SQL Server).
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
	return p >= PlaceholderQuestion && p <= PlaceholderAt
}
