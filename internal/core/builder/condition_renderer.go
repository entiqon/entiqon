package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// AppendCondition appends a valid condition to a slice and returns the result.
//
// Update: v1.4.0
func AppendCondition(existing []token.Condition, newCond token.Condition) []token.Condition {
	if newCond.IsValid() {
		return append(existing, newCond)
	}
	return existing
}

// RenderConditions builds SQL condition fragments from a list of conditions and returns the SQL string and bound args.
//
// It uses ParamBinder to safely generate dialect-aware placeholders.
//
// TODO(v1.5+):
// If the condition contains multiple logical expressions (e.g. "a = 1 AND b = 2"),
// or is explicitly marked as grouped, wrap the expression in parentheses
// to ensure correct logical precedence when combined with other conditions.
//
// Example: WHERE (a = 1 AND b = 2) OR (c = 3)
//
// Updated: v1.4.0
func RenderConditions(dialect driver.Dialect, conditions []token.Condition) (string, []any, error) {
	if len(conditions) == 0 {
		return "", nil, nil
	}

	var parts []string
	binder := driver.NewParamBinder(dialect)

	for _, c := range conditions {
		if !c.IsValid() {
			return "", nil, fmt.Errorf("invalid condition: %v", c.Error)
		}

		placeholders := binder.BindMany(c.Values...)
		placeholderExpr := strings.Join(placeholders, ", ")

		expr := fmt.Sprintf("%s %s %s", c.Key, c.Operator, placeholderExpr)
		switch c.Type {
		case token.ConditionSimple:
			parts = append(parts, expr)
		case token.ConditionAnd:
			parts = append(parts, "AND "+expr)
		case token.ConditionOr:
			parts = append(parts, "OR "+expr)
		default:
			return "", nil, fmt.Errorf("unsupported condition type: %v", c.Type)
		}
	}

	return strings.Join(parts, " "), binder.Args(), nil
}
