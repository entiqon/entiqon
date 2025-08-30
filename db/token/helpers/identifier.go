package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/entiqon/entiqon/common/extension"
	"github.com/entiqon/entiqon/db/token/types/identifier"
)

// identifierPattern defines the allowed structure of SQL identifiers.
//   - First character: letter (A–Z, a–z) or underscore (_)
//   - Remaining characters: letters, digits (0–9), or underscores
var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// IsValidIdentifier is a convenience wrapper that returns true if the identifier
// is valid, false otherwise. Prefer ValidateIdentifier when you need the reason.
func IsValidIdentifier(s string) bool {
	return ValidateIdentifier(s) == nil
}

// ResolveExpressionType determines the ExpressionKind from a raw expression.
//
// Notes:
//   - This function expects the "core" expression without alias tokens.
//     (Alias stripping is handled earlier in table.New.)
//   - Classification is purely syntactic, not semantic. For example,
//     SUM(qty) is classified as Aggregate even if used incorrectly
//     as a table source.
//
// Resolution order:
//  1. Subquery: expression starts with "(" and begins with "(SELECT ...)"
//  2. Computed: any other parenthesized expression, e.g. "(a+b)"
//  3. Aggregate: aggregate functions (SUM, COUNT, MAX, MIN, AVG)
//  4. Function: other calls with parentheses, e.g. JSON_EACH(data)
//  5. Literal: quoted string or numeric constant
//  6. Identifier: default fallback (plain table or column name)
func ResolveExpressionType(expr string) identifier.Type {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return identifier.Invalid // invalid
	}

	// Subquery → starts with "(" and specifically "(SELECT ..."
	if strings.HasPrefix(expr, "(") {
		if strings.HasPrefix(strings.ToUpper(expr), "(SELECT ") {
			return identifier.Subquery
		}
		// Any other parenthesized expression → Computed
		return identifier.Computed
	}

	// Function or Aggregate → must contain "(" and end with ")"
	if strings.Contains(expr, "(") && strings.Contains(expr, ")") {
		upper := strings.ToUpper(expr)
		if strings.HasPrefix(upper, "SUM(") ||
			strings.HasPrefix(upper, "COUNT(") ||
			strings.HasPrefix(upper, "MAX(") ||
			strings.HasPrefix(upper, "MIN(") ||
			strings.HasPrefix(upper, "AVG(") {
			return identifier.Aggregate
		}
		return identifier.Function
	}

	// Literal → numeric or quoted string
	if strings.HasPrefix(expr, "'") ||
		strings.HasPrefix(expr, "\"") ||
		extension.NumberOr(expr, -1) != -1 {
		return identifier.Literal
	}

	// Fallback → treat as plain identifier
	return identifier.Identifier
}

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
