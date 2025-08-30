package identifier

import "strings"

// Type represents the syntactic classification of a SQL expression.
//
// Resolution categories:
//   - Invalid:    could not classify
//   - Subquery:   "(SELECT ...)"
//   - Computed:   other parenthesized expressions, e.g. "(a+b)"
//   - Aggregate:  SUM, COUNT, MAX, MIN, AVG
//   - Function:   other calls with parentheses, e.g. JSON_EXTRACT(data)
//   - Literal:    quoted string or numeric constant
//   - Identifier: plain name (default fallback)
type Type int

const (
	Invalid Type = iota
	Subquery
	Computed
	Aggregate
	Function
	Literal
	Identifier
)

// Alias returns the short two-letter code used when generating
// automatic aliases for this expression kind.
//
//	Identifier → "id"
//	Literal    → "lt"
//	Function   → "fn"
//	Aggregate  → "ag"
//	Computed   → "cp"
//	Subquery   → "sq"
//	Default    → "ex"
func (k Type) Alias() string {
	switch k {
	case Identifier:
		return "id"
	case Literal:
		return "lt"
	case Function:
		return "fn"
	case Aggregate:
		return "ag"
	case Computed:
		return "cp"
	case Subquery:
		return "sq"
	default:
		return "ex"
	}
}

// IsValid reports whether the Type is recognized.
func (k Type) IsValid() bool {
	switch k {
	case Identifier, Literal, Function, Aggregate, Computed, Subquery:
		return true
	default:
		return false
	}
}

// ParseFrom attempts to coerce an arbitrary value into a Type.
// Strings are matched against String() labels (case-insensitive).
// Integers are cast directly if valid. Otherwise returns Invalid.
func (k Type) ParseFrom(value any) Type {
	switch v := value.(type) {
	case Type:
		return v
	case int:
		t := Type(v)
		if t.IsValid() {
			return t
		}
		return Invalid
	case string:
		switch s := normalize(v); s {
		case "identifier":
			return Identifier
		case "literal":
			return Literal
		case "function":
			return Function
		case "aggregate":
			return Aggregate
		case "computed":
			return Computed
		case "subquery":
			return Subquery
		default:
			return Invalid
		}
	default:
		return Invalid
	}
}

// String returns a human-readable label for the Type.
func (k Type) String() string {
	switch k {
	case Invalid:
		return "Invalid"
	case Identifier:
		return "Identifier"
	case Computed:
		return "Computed"
	case Literal:
		return "Literal"
	case Subquery:
		return "Subquery"
	case Function:
		return "Function"
	case Aggregate:
		return "Aggregate"
	default:
		return "Unknown"
	}
}

// normalize lowercases and trims a string for matching.
func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
