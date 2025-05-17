package token

import (
	"fmt"
	"strings"
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
	var c Condition
	c.Set(conditionType, condition, params...)
	return c
}

// Set assigns the condition's internal structure, resolving raw formatting.
func (c *Condition) Set(conditionType ConditionType, condition string, params ...any) {
	c.Type = conditionType
	c.Key = condition
	c.Params = params

	raw := condition
	for _, val := range params {
		raw = fmt.Sprintf("(%s)", strings.Replace(raw, "?", fmt.Sprintf("'%v'", val), 1))
	}
	c.Raw = raw
}

// IsValid checks if the condition has a non-empty key.
func (c *Condition) IsValid() bool {
	return strings.TrimSpace(c.Key) != ""
}
