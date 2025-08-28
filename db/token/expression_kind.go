package token

import (
	"regexp"
	"strconv"
	"strings"
)

var reserved = map[string]struct{}{
	"AS": {}, "SELECT": {}, "FROM": {}, "WHERE": {}, "JOIN": {},
	"ON": {}, "GROUP": {}, "ORDER": {}, "BY": {}, "LIMIT": {},
	"INSERT": {}, "UPDATE": {}, "DELETE": {}, "INTO": {}, "VALUES": {},
	"CREATE": {}, "ALTER": {}, "DROP": {}, "TABLE": {}, "INDEX": {},
	// extend with the keywords you care about
}

// ExpressionKind classifies the semantic type of Field or Table.
//
// This unifies classification for tokens that represent values
// (Field) or sources (Table). Join tokens keep their own JoinKind.
//
//   - Identifier — plain column reference (e.g. users.id)
//   - Computed   — computed expression (e.g. COUNT(*))
//   - Literal    — constant value (e.g. 'abc', 42)
//   - Subquery   — nested SELECT used as a field or a table
//   - Function   — function call used as a field or a table
type ExpressionKind int

const (
	Identifier ExpressionKind = iota // column reference
	Computed                         // computed expression
	Literal                          // constant value
	Subquery                         // subquery as field or table
	Function                         // function call as field or table
	Aggregate                        // aggregate functions like SUM(), COUNT(), etc.
)

// String returns a human-readable label for the ExpressionKind.
func (k ExpressionKind) String() string {
	switch k {
	case Identifier:
		return "IDENTIFIER"
	case Computed:
		return "COMPUTED"
	case Literal:
		return "LITERAL"
	case Subquery:
		return "SUBQUERY"
	case Function:
		return "FUNCTION"
	case Aggregate:
		return "AGGREGATE"
	default:
		return "INVALID"
	}
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

// IsValid reports whether the ExpressionKind is recognized.
func (k ExpressionKind) IsValid() bool {
	return k >= Identifier && k <= Aggregate
}

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

func isNumeric(expr string) bool {
	_, err := strconv.ParseFloat(expr, 64)
	return err == nil
}
