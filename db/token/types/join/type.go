package join

import (
	"strings"
)

// Type enumerates the supported SQL JOIN types.
//
// Type is defined as an integer enum (via iota) to provide type safety,
// while rendering uses String() for canonical SQL keywords.
//
// Typical usage:
//
//	j := join.NewLeft("users", "orders", "users.id = orders.user_id")
//	if j.Type() == join.Left {
//		fmt.Println("left join detected")
//	}
type Type int

const (
	// Invalid represents a not classified value
	Invalid Type = iota

	// Inner represents the canonical INNER JOIN between two tables.
	Inner

	// Left represents the canonical LEFT JOIN (a.k.a. LEFT OUTER JOIN).
	Left

	// Right represents the canonical RIGHT JOIN (a.k.a. RIGHT OUTER JOIN).
	Right

	// Full represents the canonical FULL JOIN (a.k.a. FULL OUTER JOIN).
	Full

	// Cross represents the canonical CROSS JOIN (Cartesian product).
	Cross

	// Natural represents the canonical NATURAL JOIN (implicit column match).
	Natural
)

// String returns the SQL keyword for the Type.
//
// The output is dialect-agnostic, using standard SQL92 keywords.
// Any unrecognized value renders as "INVALID JOIN".
func (k Type) String() string {
	switch k {
	case Inner:
		return "INNER JOIN"
	case Left:
		return "LEFT JOIN"
	case Right:
		return "RIGHT JOIN"
	case Full:
		return "FULL JOIN"
	case Cross:
		return "CROSS JOIN"
	case Natural:
		return "NATURAL JOIN"
	default:
		return "INVALID"
	}
}

// IsValid reports whether the Type is one of the recognized types.
// Invalid and any unrecognized values return false.
func (k Type) IsValid() bool {
	return k >= Inner && k <= Natural
}

// ParseFrom converts a free-form string into a Type.
//
// It accepts common variants such as "INNER", "INNER JOIN",
// "LEFT", "LEFT JOIN", etc., case-insensitively.
// If the input does not match a known kind, it returns -1,
// which is not a valid Type.
//
// Example:
//
//	jk := join.ParseJoinKindFrom("left join")
//	if !jk.IsValid() {
//		log.Fatal("invalid join kind")
//	}
//	fmt.Println(jk.String()) // LEFT JOIN
func ParseFrom(s string) Type {
	normalized := strings.ToUpper(strings.TrimSpace(s))

	switch normalized {
	case "INNER", "INNER JOIN":
		return Inner
	case "LEFT", "LEFT JOIN":
		return Left
	case "RIGHT", "RIGHT JOIN":
		return Right
	case "FULL", "FULL JOIN":
		return Full
	case "CROSS", "CROSS JOIN":
		return Cross
	case "NATURAL", "NATURAL JOIN":
		return Natural
	default:
		return Invalid
	}
}
