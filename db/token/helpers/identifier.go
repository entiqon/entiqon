package helpers

import (
	"fmt"
	"regexp"
)

// identifierPattern defines the allowed structure of SQL identifiers.
//   - First character: letter (A–Z, a–z) or underscore (_)
//   - Remaining characters: letters, digits (0–9), or underscores
var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidateIdentifier checks if s is a valid SQL identifier.
// Returns nil if valid, or an error describing why it is invalid.
//
// Reasons for invalid identifiers:
//   - Empty string
//   - Starts with a digit or invalid character
//   - Contains disallowed characters (e.g. space, dash, punctuation)
func ValidateIdentifier(s string) error {
	if s == "" {
		return fmt.Errorf("identifier cannot be empty")
	}
	if !identifierPattern.MatchString(s) {
		first := s[0]
		if first >= '0' && first <= '9' {
			return fmt.Errorf("identifier cannot start with digit: %q", s)
		}
		return fmt.Errorf("invalid identifier syntax: %q", s)
	}
	return nil
}

// IsValidIdentifier is a convenience wrapper that returns true if the identifier
// is valid, false otherwise. Prefer ValidateIdentifier when you need the reason.
func IsValidIdentifier(s string) bool {
	return ValidateIdentifier(s) == nil
}

// ValidateWildcard checks whether a wildcard expression ("*")
// is used in a valid context.
//
// Rules:
//   - Bare "*" is allowed only without an alias.
//   - If "*" is aliased or marked as raw, it is invalid.
//
// Example:
//
//	ValidateWildcard("*", "")       → ok
//	ValidateWildcard("*", "total")  → error ("* cannot be aliased or raw")
//
// Future: dialect packages may extend this to support qualified
// wildcards (e.g. "table.*") and additional rules.
func ValidateWildcard(expr, alias string) error {
	if expr == "*" && alias != "" {
		return fmt.Errorf("'*' cannot be aliased or raw")
	}
	return nil
}
