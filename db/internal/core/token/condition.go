// File: db/internal/core/token/condition.go

package token

// Condition represents a conditional expression used in a WHERE clause.
//
// Since: v0.0.1
// Updated; v1.4.0
type Condition struct {
	// Type specifies how this condition is logically joined (SIMPLE, AND, OR).
	Type ConditionType

	// Key is the SQL condition expression, e.g., "id = ?".
	Key string

	// Operator is the SQL operator to use (e.g., =, IN, LIKE, BETWEEN).
	Operator string

	// Values contains the arguments to be bound to the placeholders in Key.
	Values []any

	// Alias is optional label or usage tag
	Alias string

	// Raw contains the formatted representation, optionally for debug or display.
	Raw string

	// Error if presents means the condition is invalid
	Error error
}
