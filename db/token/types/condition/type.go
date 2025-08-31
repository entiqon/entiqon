package condition

import "strings"

// Type classifies the kind of condition expression.
//
// Conditions appear in SQL WHERE, HAVING, and ON clauses, and can be
// either atomic (a single comparison) or composite (chained with logical
// operators).
//
// The Type enumeration allows builders and tokens to represent the
// structure of a condition in a consistent way.
//
// Valid values are:
//   - Invalid: unrecognized or unsupported condition.
//   - Single:  a single expression such as "id = 1".
//   - And:     a composite condition joined with logical AND.
//   - Or:      a composite condition joined with logical OR.
type Type int

const (
	// Invalid indicates an unrecognized or unsupported condition type.
	Invalid Type = iota

	// Single represents an atomic expression, e.g. "id = 1".
	Single

	// And represents a composite condition joined with logical AND.
	// For example: "age > 18 AND active = true".
	And

	// Or represents a composite condition joined with logical OR.
	// For example: "status = 'open' OR status = 'pending'".
	Or
)

// IsValid reports whether the Type is a recognized classification.
// Returns true for Single, And, and Or; false otherwise.
func (k Type) IsValid() bool {
	switch k {
	case Single, And, Or:
		return true
	default:
		return false
	}
}

// String returns a human-readable label for the Type.
//
// Output values:
//   - Invalid → "Invalid"
//   - Single  → "" (empty string, since single expressions render
//     without a prefix operator)
//   - And     → "AND"
//   - Or      → "OR"
func (k Type) String() string {
	switch k {
	case Invalid:
		return "Invalid"
	case Single:
		return ""
	case And:
		return "AND"
	case Or:
		return "OR"
	default:
		return "Invalid"
	}
}

// ParseFrom attempts to coerce an arbitrary value into a Type.
//
// Supported conversions:
//   - Type: returned as-is.
//   - int: cast directly if the value corresponds to a valid Type.
//   - string: matched case-insensitively against known labels
//     ("and", "or"). Anything else yields Invalid.
//
// Examples:
//
//	var t Type
//	t = t.ParseFrom("and") // And
//	t = t.ParseFrom(3)     // Or
//	t = t.ParseFrom("foo") // Invalid
func ParseFrom(value any) Type {
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
		case "and":
			return And
		case "or":
			return Or
		default:
			return Invalid
		}
	default:
		return Invalid
	}
}

// normalize lowercases and trims a string for matching during parsing.
func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
