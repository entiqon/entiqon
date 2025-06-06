// Package token contains the internal type classification used to identify different token structures
// such as Column, Table, or Condition, in a type-safe and efficient way.
package token

// Kind represents the internal classification of a token type.
// It is used to distinguish between different token implementations (e.g., Column, Table, Condition)
// without relying on interface assertions or type reflection.
//
// Kind enables semantic inspection and dispatching logic across builder components.
//
// # Example
//
//	var k GetKind = ColumnKind
//	fmt.Println(k) // → 1
type Kind int

const (
	// UnknownKind is the default value assigned to tokens that do not declare a specific kind.
	// This allows for safe default behavior in case the token is malformed or unrecognized.
	UnknownKind Kind = iota

	// ColumnKind identifies tokens that represent columns in a SQL query.
	ColumnKind

	// TableKind identifies tokens that represent tables or data sources in the FROM or JOIN clause.
	TableKind

	// ConditionKind identifies tokens that express a conditional clause (e.g., WHERE, ON, HAVING).
	ConditionKind
)

// String returns the printable form of the Kind for logging and diagnostics.
func (k Kind) String() string {
	switch k {
	case ColumnKind:
		return "Column"
	case TableKind:
		return "Table"
	case ConditionKind:
		return "Condition"
	default:
		return "Unknown"
	}
}

// Kinded defines an interface for tokens that expose a classification Kind,
// which identifies their role in a SQL query (e.g., Column, Table, Condition).
//
// Implementers are expected to provide both getter and setter methods
// for internal token classification, useful for rendering and validation logic.
//
// Typical implementers include: BaseToken, Column, Table, Condition
//
// # Example
//
//	var k Kinded = NewBaseToken("users.id")
//	k.SetKind(ColumnKind)
//	fmt.Println(k.GetKind()) // → ColumnKind
//
// # Output:
//
//	ColumnKind
type Kinded interface {
	GetKind() Kind  // Returns the internal classification type
	SetKind(k Kind) // Assigns the classification kind to the token
}
