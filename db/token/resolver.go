package token

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ClassifyExpression determines the ExpressionKind from a raw expression.
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
func ClassifyExpression(expr string) ExpressionKind {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return -1 // invalid
	}

	// Subquery → starts with "(" and specifically "(SELECT ..."
	if strings.HasPrefix(expr, "(") {
		if strings.HasPrefix(strings.ToUpper(expr), "(SELECT ") {
			return Subquery
		}
		// Any other parenthesized expression → Computed
		return Computed
	}

	// Function or Aggregate → must contain "(" and end with ")"
	if strings.Contains(expr, "(") && strings.Contains(expr, ")") {
		upper := strings.ToUpper(expr)
		if strings.HasPrefix(upper, "SUM(") ||
			strings.HasPrefix(upper, "COUNT(") ||
			strings.HasPrefix(upper, "MAX(") ||
			strings.HasPrefix(upper, "MIN(") ||
			strings.HasPrefix(upper, "AVG(") {
			return Aggregate
		}
		return Function
	}

	// Literal → numeric or quoted string
	if strings.HasPrefix(expr, "'") ||
		strings.HasPrefix(expr, "\"") ||
		isNumeric(expr) {
		return Literal
	}

	// Fallback → treat as plain identifier
	return Identifier
}

func isNumeric(expr string) bool {
	_, err := strconv.ParseFloat(expr, 64)
	return err == nil
}

// ResolveExpr splits an input into expr + alias, then classifies the expr.
// If allowAlias=false, alias parts are rejected.
func ResolveExpr(input string, allowAlias bool) (kind ExpressionKind, expr, alias string, err error) {
	in := strings.TrimSpace(input)
	if in == "" {
		return -1, "", "", errors.New("empty input")
	}

	// ✅ Subquery detection: treat whole input as one expression
	// even if it contains spaces.
	if strings.HasPrefix(in, "(") && strings.HasSuffix(in, ")") {
		return Subquery, in, "", nil
	}

	// If explicit "AS" exists → split at it
	up := strings.ToUpper(in)
	if strings.Contains(up, " AS ") {
		parts := strings.SplitN(in, " AS ", 2)
		expr = strings.TrimSpace(parts[0])
		alias = strings.TrimSpace(parts[1])
		if !allowAlias {
			return -1, "", "", errors.New("alias not allowed: " + input)
		}
		if !IsValidAlias(alias) {
			return -1, "", "", errors.New("invalid alias: " + alias)
		}
		return ClassifyExpression(expr), expr, alias, nil
	}

	// Otherwise, split by whitespace
	parts := strings.Fields(in)
	switch len(parts) {
	case 1:
		expr = parts[0]
	case 2:
		if !allowAlias {
			return -1, "", "", errors.New("alias not allowed: " + input)
		}
		expr, alias = parts[0], parts[1]
		if !IsValidAlias(alias) {
			return -1, "", "", errors.New("invalid alias: " + alias)
		}
	default:
		// Handle cases like "(SELECT * FROM users) u"
		if HasTrailingAliasWithoutAS(in) {
			alias = parts[len(parts)-1]
			expr = strings.Join(parts[:len(parts)-1], " ")
			if !allowAlias {
				return -1, "", "", errors.New("alias not allowed: " + input)
			}
			if !IsValidAlias(alias) {
				return -1, "", "", errors.New("invalid alias: " + alias)
			}
		} else {
			return -1, "", "", errors.New("invalid input: " + input)
		}
	}

	kind = ClassifyExpression(expr)

	// Strict check: Identifiers must be a single token
	if kind == Identifier {
		coreParts := strings.Fields(expr)
		if len(coreParts) != 1 {
			return -1, expr, alias, errors.New("invalid format " + expr)
		}
	}

	return kind, expr, alias, nil
}

// HasTrailingAliasWithoutAS checks if the last space-separated token is an alias candidate
func HasTrailingAliasWithoutAS(expr string) bool {
	up := strings.ToUpper(expr)
	if strings.Contains(up, " AS ") {
		return false // explicit AS → fine
	}

	tokens := strings.Fields(expr)
	if len(tokens) <= 1 {
		return false // single token can't have alias
	}

	last := tokens[len(tokens)-1]
	penultimate := tokens[len(tokens)-2]

	// If the token before last is an operator, this "last" is part of the expression, not alias
	operators := map[string]bool{"||": true, "+": true, "-": true, "*": true, "/": true}
	if operators[penultimate] {
		return false
	}

	// Otherwise, if it looks like an identifier → treat as alias
	if regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(last) {
		return true
	}
	return false
}

// IsValidAlias reports whether s is a valid SQL identifier alias.
// Must match identifier syntax AND not be a reserved keyword.
func IsValidAlias(s string) bool {
	if !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(s) {
		return false
	}
	_, found := reserved[strings.ToUpper(s)]
	return !found
}
