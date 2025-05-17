package token

// ConditionType defines how a condition is logically connected in a WHERE clause.
type ConditionType string

const (
	// ConditionSimple is used for initial WHERE conditions.
	ConditionSimple ConditionType = "SIMPLE"

	// ConditionAnd adds an AND between conditions.
	ConditionAnd ConditionType = "AND"

	// ConditionOr adds an OR between conditions.
	ConditionOr ConditionType = "OR"
)

// Condition represents a conditional expression used in a WHERE clause.
type Condition struct {
	// Type specifies how this condition is logically joined (SIMPLE, AND, OR).
	Type ConditionType

	// Key is the SQL condition expression, e.g., "id = ?".
	Key string

	Alias string // optional label or usage tag

	// Params contains the arguments to be bound to the placeholders in Key.
	Params []any

	// Raw contains the formatted representation, optionally for debug or display.
	Raw string
}
