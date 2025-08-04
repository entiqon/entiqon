// File: db/internal/core/token/condition_type.go

package token

// ConditionType defines how a condition is logically connected in a WHERE clause.
//
// Since: v0.0.1
// Update: v1.4.0
type ConditionType string

const (
	// ConditionSimple is used for initial WHERE conditions.
	//
	// Since: v0.0.1
	// Update: v1.4.0
	ConditionSimple ConditionType = "SIMPLE"

	// ConditionAnd adds an AND between conditions.
	//
	// Since: v0.0.1
	// Update: v1.4.0
	ConditionAnd ConditionType = "AND"

	// ConditionOr adds an OR between conditions.
	//
	// Since: v0.0.1
	// Update: v1.4.0
	ConditionOr ConditionType = "OR"
)
