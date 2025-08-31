package helpers

import (
	stdErr "errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/entiqon/entiqon/common/extension"
	"github.com/entiqon/entiqon/db/contract"
	"github.com/entiqon/entiqon/db/errors"
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

// ResolveExpression classifies and normalizes a SQL expression.
//
// It enforces type validation via ValidateType, and only accepts raw strings.
// Passing an existing token (Field, Table, etc.) is rejected with a clear
// message to use Clone() instead.
//
// Classification covers identifiers, subqueries, computed expressions,
// aggregates, functions, and literals. Inline or explicit aliases are
// supported depending on allowAlias.
func ResolveExpression(
	input any,
	allowAlias bool,
) (identifier.Type, string, string, error) {
	in := strings.TrimSpace(input.(string))
	kind := ResolveExpressionType(in)

	switch kind {
	case identifier.Identifier:
		parts := strings.Fields(in)
		switch len(parts) {
		case 1:
			return identifier.Identifier, parts[0], "", nil
		case 2:
			if !allowAlias {
				return identifier.Invalid, "", "", stdErr.New("alias not allowed: " + in)
			}
			if !IsValidAlias(parts[1]) {
				return identifier.Invalid, parts[0], parts[1],
					stdErr.New("invalid alias: " + parts[1])
			}
			return identifier.Identifier, parts[0], parts[1], nil
		case 3:
			if strings.EqualFold(parts[1], "AS") {
				if !allowAlias {
					return identifier.Invalid, "", "", stdErr.New("alias not allowed: " + in)
				}
				if !IsValidAlias(parts[2]) {
					return identifier.Invalid, parts[0], parts[2],
						stdErr.New("invalid alias: " + parts[2])
				}
				return identifier.Identifier, parts[0], parts[2], nil
			}
			return identifier.Invalid, "", "", stdErr.New("invalid identifier: " + in)
		default:
			return identifier.Invalid, "", "", stdErr.New("invalid identifier: " + in)
		}

	case identifier.Subquery:
		closeIdx := strings.LastIndex(in, ")")
		expr := strings.TrimSpace(in[:closeIdx+1])
		rest := strings.TrimSpace(in[closeIdx+1:])
		alias, err := resolveAlias(rest, in, allowAlias)
		if err != nil {
			return identifier.Invalid, expr, alias, err
		}
		return identifier.Subquery, expr, alias, nil

	case identifier.Computed, identifier.Aggregate, identifier.Function:
		closeIdx := strings.LastIndex(in, ")")
		expr := strings.TrimSpace(in[:closeIdx+1])
		rest := strings.TrimSpace(in[closeIdx+1:])
		alias, err := resolveAlias(rest, in, allowAlias)
		if err != nil {
			return identifier.Invalid, expr, alias, err
		}
		return kind, expr, alias, nil

	case identifier.Literal:
		parts := strings.Fields(in)
		expr := parts[0]
		rest := ""
		if len(parts) > 1 {
			rest = strings.Join(parts[1:], " ")
		}
		alias, err := resolveAlias(rest, in, allowAlias)
		if err != nil {
			return identifier.Invalid, expr, alias, err
		}
		return identifier.Literal, expr, alias, nil

	default:
		if in == "" {
			return identifier.Invalid, "", "", fmt.Errorf("empty identifier is not allowed: %q", in)
		}
		return identifier.Invalid, "", "", stdErr.New("invalid because unknown identifier: " + in)
	}
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

	upper := strings.ToUpper(expr)

	// ✅ Subquery
	// Case 1: Parenthesized SELECT
	if strings.HasPrefix(expr, "(") && strings.HasPrefix(upper, "(SELECT ") {
		return identifier.Subquery
	}
	// Case 2: Bare SELECT
	if strings.HasPrefix(upper, "SELECT ") {
		return identifier.Subquery
	}

	// Computed
	if strings.HasPrefix(expr, "(") && strings.Contains(expr, ")") {
		return identifier.Computed
	}

	// Aggregate
	if strings.HasPrefix(upper, "SUM(") ||
		strings.HasPrefix(upper, "COUNT(") ||
		strings.HasPrefix(upper, "MAX(") ||
		strings.HasPrefix(upper, "MIN(") ||
		strings.HasPrefix(upper, "AVG(") {
		return identifier.Aggregate
	}

	// Function
	if strings.Contains(expr, "(") && strings.Contains(expr, ")") {
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

// ValidateType ensures input type is allowed for token constructors.
//
// Rules:
//   - string → OK
//   - any Validable (Field, Table, etc.) → ErrUnsupportedType
//   - everything else → invalid format with type name
func ValidateType(v any) error {
	switch v.(type) {
	case string:
		return nil

	case contract.Validable:
		// Already a token (Field, Table, etc.)
		return errors.UnsupportedTypeError

	default:
		// Everything else → invalid format
		return fmt.Errorf("invalid format (type %T)", v)
	}
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

// resolveAlias parses "expr" (text after the main expression) into an alias.
// Supports:
//   - "" → no alias
//   - "alias" → simple alias
//   - "AS alias" → explicit alias
//
// Applies allowAlias and IsValidAlias checks.
func resolveAlias(expr, in string, allowAlias bool) (string, error) {
	if expr == "" {
		return "", nil
	}

	parts := strings.Fields(expr)

	// "AS alias"
	if len(parts) == 2 && strings.EqualFold(parts[0], "AS") {
		if !allowAlias {
			return "", fmt.Errorf("alias not allowed: %s", in)
		}
		if !IsValidAlias(parts[1]) {
			return "", fmt.Errorf("invalid alias: %s", parts[1])
		}
		return parts[1], nil
	}

	// simple alias
	if len(parts) == 1 {
		if !allowAlias {
			return "", fmt.Errorf("alias not allowed: %s", in)
		}
		if !IsValidAlias(parts[0]) {
			return "", fmt.Errorf("invalid alias: %s", parts[0])
		}
		return parts[0], nil
	}

	return "", fmt.Errorf("invalid alias format: %s", expr)
}
