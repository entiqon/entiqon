// File: db/token/field/contract.go

package condition

import (
	"github.com/entiqon/db/contract"
	"github.com/entiqon/db/token/types/condition"
	"github.com/entiqon/db/token/types/operator"
)

// Token is the contract implemented by *Field.
//
// A Token represents a SQL field (column or expression) that may be
// aliased and owned by a table. It combines core identity/validation
// from BaseToken with shared token contracts (Renderable, Rawable, etc.),
// and provides access to the owning table when relevant.
//
// Implementations must guarantee immutability of construction, and
// produce safe, independent clones for mutation via Clonable.
//
// Typical use cases:
//   - SelectBuilder fields (SELECT id, name AS username)
//   - Expressions (COUNT(*), JSON_EXTRACT(data, '$.field'))
//   - Qualified fields with table alias (u.id AS user_id)
//
// Example:
//
//	f := field.New("u.id", "user_id")
//	var tok field.Token = f
//	fmt.Println(tok.Render()) // "u.id AS user_id"
type Token interface {
	// Kindable returns the semantic classification of the
	// underlying expression (Identifier, Literal, Subquery, etc.).
	contract.Kindable[condition.Type]

	// Identifiable is implemented by tokens that expose their core identity
	// without aliasing.
	contract.Identifiable

	// Errorable reports whether the field is invalid and provides its error.
	contract.Errorable[Token]

	// Debuggable returns a detailed diagnostic view of the field,
	// including markers for aliased/raw/errored state.
	contract.Debuggable

	// Rawable returns the original input string(s), useful for auditing/logging.
	contract.Rawable

	// Renderable returns the SQL-safe representation of the field.
	contract.Renderable

	// Stringable provides a human-readable summary of the field,
	// typically used for debugging.
	contract.Stringable

	// Validable reports whether the table token is structurally valid.
	// A token is valid if it has no error and its essential invariants
	// (e.g., non-empty name) are satisfied.
	contract.Validable

	// Name returns the binding key for the parameter if a value
	Name() string

	Operator() operator.Type

	// Value returns the bound value associated to the condition, if any.
	// It is intended for parameter binding in prepared statements.
	Value() any
}

// Ensure *Field implements Token at compile time.
var _ Token = (*token)(nil)
