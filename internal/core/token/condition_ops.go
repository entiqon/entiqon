package token

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/driver"
)

// AppendCondition appends a valid condition to a slice and returns the result.
func AppendCondition(existing []Condition, newCond Condition) []Condition {
	if newCond.IsValid() {
		return append(existing, newCond)
	}
	return existing
}

// NewCondition creates a Condition instance with the given type, key, and parameters.
func NewCondition(conditionType ConditionType, condition string, params ...any) Condition {
	return Condition{}.Set(conditionType, condition, params...)
}

// Set assigns the condition's internal structure, resolving raw formatting.
func (c Condition) Set(conditionType ConditionType, condition string, params ...any) Condition {
	c.Type = conditionType
	c.Key = condition
	c.Params = params

	raw := condition
	for _, val := range params {
		raw = fmt.Sprintf("(%s)", strings.Replace(raw, "?", fmt.Sprintf("'%v'", val), 1))
	}
	c.Raw = raw

	return c
}

// IsValid checks if the condition has a non-empty key.
func (c Condition) IsValid() bool {
	return strings.TrimSpace(c.Key) != ""
}

// FormatConditions formats a slice of Condition tokens into a SQL WHERE clause
// expression and its associated parameter values.
//
// It returns a condition string like "id = ? AND status = ?" and a flat slice of arguments.
// If a dialect is provided, it's available for escaping identifiers,
// but it's not used in the default implementation.
//
// The output is intended for direct use in SQL builders like SelectBuilder or DeleteBuilder.
func FormatConditions(dialect driver.Dialect, conditions []Condition) (string, []any, error) {
	var parts []string
	var args []any

	for _, c := range conditions {
		switch c.Type {
		case ConditionSimple:
			parts = append(parts, c.Key)
		case ConditionAnd, ConditionOr:
			parts = append(parts, fmt.Sprintf("%s %s", c.Type, c.Key))
		default:
			return "", nil, fmt.Errorf("invalid condition type: %s", c.Type)
		}
		args = append(args, c.Params...)
	}

	return strings.Join(parts, " "), args, nil
}
