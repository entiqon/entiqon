package token

import (
	"fmt"
	"strings"
	"time"

	"github.com/ialopezg/entiqon/internal/core/driver"
)

// AppendCondition appends a valid condition to a slice and returns the result.
//
// Update: v1.4.0
func AppendCondition(existing []Condition, newCond Condition) []Condition {
	if newCond.IsValid() {
		return append(existing, newCond)
	}
	return existing
}

// FormatConditions builds SQL condition fragments from a list of conditions and returns the SQL string and bound args.
//
// It uses ParamBinder to safely generate dialect-aware placeholders.
//
// Updated: v1.4.0
func FormatConditions(dialect driver.Dialect, conditions []Condition) (string, []any, error) {
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
		case ConditionSimple:
			parts = append(parts, expr)
		case ConditionAnd:
			parts = append(parts, "AND "+expr)
		case ConditionOr:
			parts = append(parts, "OR "+expr)
		default:
			return "", nil, fmt.Errorf("unsupported condition type: %v", c.Type)
		}
	}

	return strings.Join(parts, " "), binder.Args(), nil
}

// AreCompatibleTypes checks whether all provided values belong to the same compatible type group.
//
// Since: v1.4.0
func AreCompatibleTypes(values ...any) bool {
	if len(values) < 2 {
		return false
	}

	switch values[0].(type) {
	case string:
		for _, v := range values[1:] {
			if _, ok := v.(string); !ok {
				return false
			}
		}
	case int, int64, float32, float64:
		for _, v := range values[1:] {
			switch v.(type) {
			case int, int64, float32, float64:
				// ok
			default:
				return false
			}
		}
	case time.Time:
		for _, v := range values[1:] {
			if _, ok := v.(time.Time); !ok {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func extractConditionParts(input string) (field string, operator string, value string, ok bool) {
	supportedOps := []string{
		"NOT IN", "IN", "BETWEEN", "<>", "!=", ">=", "<=", "LIKE", "=", ">", "<",
	}
	input = strings.TrimSpace(input)
	upper := strings.ToUpper(input)

	for _, op := range supportedOps {
		if idx := strings.Index(upper, op); idx > 0 {
			field = strings.TrimSpace(input[:idx])
			value = strings.TrimSpace(input[idx+len(op):])
			return field, op, value, true
		}
	}
	return "", "", "", false
}
