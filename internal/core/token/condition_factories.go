// filename: /internal/core/token/condition_factories.go

package token

import (
	"fmt"
)

// NewCondition creates a condition from either a raw string expression like
// "status != active" or structured field/value input.
//
// - If one param is passed: parses inline condition (e.g., "score >= 100")
// - If two params are passed: interprets as field + value (defaults to "=")
// - If more than two, returns an error.
//
// Since: v0.0.1
// Updated: v1.4.0
func NewCondition(conditionType ConditionType, params ...any) Condition {
	if len(params) > 2 {
		return Condition{Error: fmt.Errorf("invalid number of parameters")}
	}

	var field string
	var operator string
	var value any

	if len(params) == 1 {
		raw, ok := params[0].(string)
		if !ok {
			return Condition{Error: fmt.Errorf("inline condition must be a string")}
		}

		field, operator, value, ok = extractConditionParts(raw)
		if !ok {
			return Condition{Error: fmt.Errorf("unable to parse inline condition: %q", raw)}
		}
	} else {
		var ok bool
		field, ok = params[0].(string)
		if !ok {
			return Condition{Error: fmt.Errorf("field must be a string")}
		}
		operator = "="
		value = params[1]
	}

	return NewConditionWithOperator(conditionType, field, operator, value)
}

// NewConditionBetween creates a BETWEEN condition with exactly two values.
//
// Since: v1.4.0
func NewConditionBetween(conditionType ConditionType, field string, start any, end any) Condition {
	if field == "" || start == nil || end == nil {
		return Condition{Error: fmt.Errorf("BETWEEN requires a field and two non-nil values")}
	}
	if s, ok := start.(string); ok && s == "" {
		return Condition{Error: fmt.Errorf("BETWEEN start value cannot be empty string")}
	}
	if e, ok := end.(string); ok && e == "" {
		return Condition{Error: fmt.Errorf("BETWEEN end value cannot be empty string")}
	}
	if !AreCompatibleTypes(start, end) {
		return Condition{Error: fmt.Errorf("BETWEEN values must be of compatible types: got %T and %T", start, end)}
	}
	return NewConditionWithOperator(conditionType, field, "BETWEEN", start, end)
}

// NewConditionIn creates an IN condition with multiple values.
//
// Since: v1.4.0
func NewConditionIn(conditionType ConditionType, field string, values ...any) Condition {
	if !AreCompatibleTypes(values...) {
		return Condition{Error: fmt.Errorf("IN values must be of compatible types")}
	}
	return NewConditionWithOperator(conditionType, field, "IN", values...)
}

// NewConditionNotIn creates a NOT IN condition with multiple values.
//
// Since: v1.4.0
func NewConditionNotIn(conditionType ConditionType, field string, values ...any) Condition {
	if !AreCompatibleTypes(values...) {
		return Condition{Error: fmt.Errorf("NOT IN values must be of compatible types")}
	}
	return NewConditionWithOperator(conditionType, field, "NOT IN", values...)
}

// NewConditionGreaterThan creates a > condition with a single value.
//
// Since: v1.4.0
func NewConditionGreaterThan(conditionType ConditionType, field string, value any) Condition {
	return NewConditionWithOperator(conditionType, field, ">", value)
}

// NewConditionGreaterThanOrEqual creates a >= condition with a single value.
//
// Since: v1.4.0
func NewConditionGreaterThanOrEqual(conditionType ConditionType, field string, value any) Condition {
	return NewConditionWithOperator(conditionType, field, ">=", value)
}

// NewConditionLessThan creates a < condition with a single value.
//
// Since: v1.4.0
func NewConditionLessThan(conditionType ConditionType, field string, value any) Condition {
	return NewConditionWithOperator(conditionType, field, "<", value)
}

// NewConditionLessThanOrEqual creates a <= condition with a single value.
//
// Since: v1.4.0
func NewConditionLessThanOrEqual(conditionType ConditionType, field string, value any) Condition {
	return NewConditionWithOperator(conditionType, field, "<=", value)
}

// NewConditionLike creates a LIKE condition with a single pattern value.
//
// Since: v1.4.0
func NewConditionLike(conditionType ConditionType, field string, pattern any) Condition {
	return NewConditionWithOperator(conditionType, field, "LIKE", pattern)
}

// NewConditionNotEqual creates a != condition with a single value.
//
// Since: v1.4.0
func NewConditionNotEqual(conditionType ConditionType, field string, value any) Condition {
	return NewConditionWithOperator(conditionType, field, "!=", value)
}

// NewConditionWithOperator creates a condition with an explicit operator and one or more values.
//
// Since: v1.4.0
func NewConditionWithOperator(conditionType ConditionType, field, operator string, values ...any) Condition {
	c := Condition{
		Type:     conditionType,
		Key:      field,
		Operator: operator,
		Values:   values,
	}

	if field == "" || operator == "" || len(values) == 0 {
		c.Error = fmt.Errorf("invalid condition parameters: field='%s', operator='%s', values=%d", field, operator, len(values))
		return c
	}

	if operator == "IN" || operator == "NOT IN" {
		c.Raw = fmt.Sprintf("%s %s (:%s)", field, operator, field)
	} else if operator == "BETWEEN" && len(values) == 2 {
		c.Raw = fmt.Sprintf("%s BETWEEN :%s_start AND :%s_end", field, field, field)
	} else {
		c.Raw = fmt.Sprintf("%s %s :%s", field, operator, field)
	}

	return c
}
