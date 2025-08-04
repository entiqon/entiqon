// File: db/internal/core/token/condition_factories.go

package token

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/internal/core/builder/util"
)

// NewCondition builds a Condition by parsing a condition expression and optional value(s).
// It supports semantic validation, placeholder patterns, and operator resolution,
// producing dialect-safe conditions for use in SQL builders.
//
// Examples:
//
//	NewCondition(token.ConditionSimple, "status = active")           // inline
//	NewCondition(token.ConditionSimple, "status = ?", "active")      // placeholder
//	NewCondition(token.ConditionSimple, "status IN ?", []any{"a", "b"})
//	NewCondition(token.ConditionSimple, "price BETWEEN ? AND ?", []any{10, 20})
//
// Supported operators:
//
//	=, !=, <>, >, <, >=, <=
//	IN, NOT IN
//	BETWEEN
//	LIKE, NOT LIKE
//
// Returns a Condition with .Error if:
//   - A placeholder is used with no value
//   - Types are mixed in, IN or BETWEEN
//   - Length is incorrect for BETWEEN
//
// Since: v0.0.1
// Updated: v1.4.0
func NewCondition(conditionType ConditionType, name string, value ...any) Condition {
	if value == nil || len(value) == 0 {
		// Attempt inline literal resolution
		field, operator, literal, ok := extractConditionParts(name)
		if literal == "" {
			return Condition{Error: fmt.Errorf("unable to parse condition")}
		}
		if util.ContainsUnboundPlaceholder(literal) {
			return Condition{Error: fmt.Errorf("placeholder without a value for %s=", name)}
		}
		v := util.InferLiteralType(literal)
		if !ok {
			return Condition{Error: fmt.Errorf("%s: unable to parse inline condition: %q", conditionType, name)}
		}
		return NewConditionWithOperator(conditionType, field, operator, v)
	}

	if len(value) > 2 {
		return NewConditionWithOperator(conditionType, name, "IN", value...)
	}

	field, operator, l, ok := extractConditionParts(name)
	if !ok {
		return NewConditionWithOperator(conditionType, name, "=", value[0])
	}
	// Flatten if single slice is passed
	if len(value) == 1 {
		if inner, ok := value[0].([]any); ok {
			value = inner
		}
	}

	literal := util.InferLiteralType(l)
	if value == nil {
		return NewConditionWithOperator(conditionType, field, operator, literal)
	}

	switch operator {
	case "IN", "NOT IN":
		if len(value) == 1 {
			// Downgrade to "=" or "!=" if only one value
			operator = map[string]string{"IN": "=", "NOT IN": "!="}[operator]
		} else {
			if !util.AllSameType(value) {
				return Condition{Error: fmt.Errorf("%s: values for %s must be of the same type", conditionType, operator)}
			}
		}

	case "BETWEEN":
		if len(value) != 2 {
			return Condition{Error: fmt.Errorf("%s: BETWEEN requires exactly 2 values", conditionType)}
		}
		if !util.AllSameType(value) {
			return Condition{Error: fmt.Errorf("%s: BETWEEN values must be of the same type", conditionType)}
		}
	}

	return NewConditionWithOperator(conditionType, field, operator, value...)
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
		switch v := values[0].(type) {
		case []any:
			values = v
		case []int:
			values = make([]any, len(v))
			for i, x := range v {
				values[i] = x
			}
		case []string:
			values = make([]any, len(v))
			for i, x := range v {
				values[i] = x
			}
		}
	}

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

	if strings.TrimSpace(field) == "" || strings.TrimSpace(operator) == "" || values == nil {
		c.Error = fmt.Errorf("%s: invalid condition parameters: field='%s', operator='%s', values=%d", conditionType, field, operator, len(values))
		return c
	}

	switch operator {
	case "IN", "NOT IN":
		c.Raw = fmt.Sprintf("%s %s (:%s)", field, operator, field)
	case "BETWEEN":
		if len(values) == 2 {
			c.Raw = fmt.Sprintf("%s BETWEEN :%s_start AND :%s_end", field, field, field)
		}
	default:
		c.Raw = fmt.Sprintf("%s %s :%s", field, operator, field)
	}

	return c
}
