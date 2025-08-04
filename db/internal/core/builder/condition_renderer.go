// File: db/internal/core/builder/condition_renderer.go

package builder

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/driver"
	"github.com/entiqon/entiqon/db/internal/core/builder/bind"
	token2 "github.com/entiqon/entiqon/db/internal/core/token"
)

// AppendCondition appends a valid condition to a slice and returns the result.
//
// Update: v1.4.0
func AppendCondition(existing []token2.Condition, newCond token2.Condition) []token2.Condition {
	if newCond.IsValid() {
		return append(existing, newCond)
	}
	if len(existing) == 0 && newCond.Type != token2.ConditionSimple {
		newCond.Type = token2.ConditionSimple
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
func RenderConditions(dialect driver.Dialect, conditions []token2.Condition) (string, []any, error) {
	binder := bind.NewParamBinder(dialect)
	return RenderConditionsWithBinder(dialect, conditions, binder)
}

// RenderConditionsWithBinder builds a SQL WHERE clause using the given conditions
// and a provided ParamBinder, which allows control over placeholder indexing.
//
// This is useful in builders like UPDATE or UPSERT where the argument list
// already contains values for SET assignments, and WHERE placeholders must
// continue indexing without collisions.
//
// Example:
//
//	binder := bind.NewParamBinderWithPosition(dialect, len(args))
//	clause, args, err := RenderConditionsWithBinder(dialect, conditions, binder)
//
// Since: v1.4.0
func RenderConditionsWithBinder(dialect driver.Dialect, conditions []token2.Condition, binder *bind.ParamBinder) (string, []any, error) {
	if len(conditions) == 0 {
		return "", nil, nil
	}

	var parts []string

	for _, c := range conditions {
		if !c.IsValid() {
			return "", nil, fmt.Errorf("invalid condition: %v", c.Error)
		}

		placeholders := binder.BindMany(c.Values...)
		placeholderExpr := strings.Join(placeholders, ", ")

		expr := fmt.Sprintf("%s %s %s", dialect.QuoteIdentifier(c.Key), c.Operator, placeholderExpr)
		switch c.Type {
		case token2.ConditionSimple:
			parts = append(parts, expr)
		case token2.ConditionAnd:
			parts = append(parts, "AND "+expr)
		case token2.ConditionOr:
			parts = append(parts, "OR "+expr)
		default:
			return "", nil, fmt.Errorf("unsupported condition type: %v", c.Type)
		}
	}

	return strings.Join(parts, " "), binder.Args(), nil
}
