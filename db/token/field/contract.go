// File: db/token/field/contract.go

package field

import "github.com/entiqon/entiqon/db/contract"

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
	// BaseToken provides access to the original input, normalized
	// expression, alias, and validation state.
	contract.BaseToken

	// Clonable allows producing a semantic copy safe for independent mutation.
	contract.Clonable[Token]

	// Debuggable returns a detailed diagnostic view of the field,
	// including markers for aliased/raw/errored state.
	contract.Debuggable

	// Errorable reports whether the field is invalid and provides its error.
	contract.Errorable[Token]

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

	// HasOwner reports whether the field is qualified by a table name or alias.
	//
	// It returns true if an owner (e.g. table alias "u" in "u.id") has been
	// assigned, and false otherwise.
	HasOwner() bool

	// Owner returns the owning table name (or alias) if one is set.
	//
	// Must be called only when HasOwner() returns true; otherwise the result
	// is undefined.
	Owner() *string

	// SetOwner assigns or clears the owning table name (or alias).
	//
	// Passing nil clears the owner. Implementations should ensure that calling
	// SetOwner does not violate immutability guarantees (typically by applying
	// the change on a clone rather than mutating the original instance).
	SetOwner(owner *string)
}

// Ensure *Field implements Token at compile time.
var _ Token = (*field)(nil)
