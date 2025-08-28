package join

import (
	"fmt"
	"strings"
)

// Kind enumerates the supported SQL JOIN types.
//
// Kind is defined as an integer enum (via iota) to provide type safety,
// while rendering uses String() for canonical SQL keywords.
//
// Typical usage:
//
//	j := join.NewLeft("users", "orders", "users.id = orders.user_id")
//	if j.Kind() == token.Left {
//		fmt.Println("left join detected")
//	}
type Kind int

const (
	// Inner represents the canonical INNER JOIN between two tables.
	Inner Kind = iota

	// Left represents the canonical LEFT JOIN (a.k.a. LEFT OUTER JOIN).
	Left

	// Right represents the canonical RIGHT JOIN (a.k.a. RIGHT OUTER JOIN).
	Right

	// Full represents the canonical FULL JOIN (a.k.a. FULL OUTER JOIN).
	Full
)

// String returns the SQL keyword for the Kind.
//
// The output is dialect-agnostic, using standard SQL92 keywords.
// Any unrecognized value renders as "INVALID JOIN".
func (k Kind) String() string {
	switch k {
	case Inner:
		return "INNER JOIN"
	case Left:
		return "LEFT JOIN"
	case Right:
		return "RIGHT JOIN"
	case Full:
		return "FULL JOIN"
	default:
		return fmt.Sprintf("invalid join type (%d)", int(k))
	}
}

// IsValid reports whether the Kind is one of the recognized types
// (INNER, LEFT, RIGHT, FULL). InvalidJoin and any unrecognized value
// return false.
func (k Kind) IsValid() bool {
	return k >= Inner && k <= Full
}

// ParseJoinKindFrom converts a free-form string into a Kind.
//
// It accepts common variants such as "INNER", "INNER JOIN",
// "LEFT", "LEFT JOIN", etc., case-insensitively.
// If the input does not match a known kind, it returns -1,
// which is not a valid Kind.
//
// Example:
//
//	jk := token.ParseJoinKindFrom("left join")
//	if !jk.IsValid() {
//		log.Fatal("invalid join kind")
//	}
//	fmt.Println(jk.String()) // LEFT JOIN
func ParseJoinKindFrom(s string) Kind {
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
	default:
		return -1 // invalid, IsValid() will return false
	}
}
