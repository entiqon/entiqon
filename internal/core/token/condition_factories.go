// filename: /internal/core/token/condition_factories.go

package token

import "fmt"

// NewCondition creates a condition from either a raw string expression like
// "status != active" or structured field/value input.
//
// - If one param is passed: parses inline condition (e.g., "score >= 100")
// - If two params are passed: interprets as field + value (defaults to "=")
// - If more than two, returns an error.
//
// Since: v0.0.1
// Updated: v1.4.0
func NewCondition(conditionType ConditionType, condition string, params ...any) Condition {
	switch len(params) {
	case 0:
		// Inline syntax like: "status = active"
		f, o, v, ok := extractConditionParts(condition)
		if !ok {
			return Condition{Error: fmt.Errorf("%s: unable to parse condition: %q", conditionType, condition)}
		}
		if v == "" {
			return Condition{Error: fmt.Errorf("%s: empty value on condition: %q", conditionType, condition)}
		}
		return NewConditionWithOperator(conditionType, f, o, v)

	case 1:
		// Explicit syntax like: NewCondition(type, "status", "active")
		return NewConditionWithOperator(conditionType, condition, "=", params...)

	default:
		return Condition{Error: fmt.Errorf("%s: invalid number of parameters on condition", conditionType)}
	}
}

// NewConditionBetween creates a BETWEEN condition with exactly two values.
//
// Since: v1.4.0
func NewConditionBetween(conditionType ConditionType, field string, start any, end any) Condition {
	if field == "" || start == nil || end == nil {
		return Condition{Error: fmt.Errorf("%s: BETWEEN requires a field and two non-nil values", conditionType)}
	}
	if s, ok := start.(string); ok && s == "" {
		return Condition{Error: fmt.Errorf("%s: BETWEEN start value cannot be empty string", conditionType)}
	}
	if e, ok := end.(string); ok && e == "" {
		return Condition{Error: fmt.Errorf("%s: BETWEEN end value cannot be empty string", conditionType)}
	}
	if !AreCompatibleTypes(start, end) {
		return Condition{Error: fmt.Errorf("%s: BETWEEN values must be of compatible types: got %T and %T", conditionType, start, end)}
	}
	return NewConditionWithOperator(conditionType, field, "BETWEEN", start, end)
}

// NewConditionIn creates an IN condition with multiple values.
//
// Since: v1.4.0
func NewConditionIn(conditionType ConditionType, field string, values ...any) Condition {
	if !AreCompatibleTypes(values...) {
		return Condition{Error: fmt.Errorf("%s: IN values must be of compatible types on: %q", conditionType, field)}
	}
	return NewConditionWithOperator(conditionType, field, "IN", values...)
}

// NewConditionNotIn creates a NOT IN condition with multiple values.
//
// Since: v1.4.0
func NewConditionNotIn(conditionType ConditionType, field string, values ...any) Condition {
	if !AreCompatibleTypes(values...) {
		return Condition{Error: fmt.Errorf("%s: NOT IN values must be of compatible types", conditionType)}
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
	if len(values) == 1 {
		if inner, ok := values[0].([]any); ok && (operator == "IN" || operator == "NOT IN") {
			values = inner
		}
	}

	c := Condition{
		Type:     conditionType,
		Key:      field,
		Operator: operator,
		Values:   values,
	}

	if field == "" || operator == "" || len(values) == 0 {
		c.Error = fmt.Errorf("%s: invalid condition parameters: field='%s', operator='%s', values=%d", conditionType, field, operator, len(values))
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
