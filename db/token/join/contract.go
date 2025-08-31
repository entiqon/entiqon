package join

import (
	"github.com/entiqon/entiqon/db/contract"
	"github.com/entiqon/entiqon/db/token/table"
	"github.com/entiqon/entiqon/db/token/types/join"
)

// Token is the contract implemented by all Join tokens.
//
// Join tokens expose the standard contracts used across all tokens:
//
//   - Clonable   → safe deep copies
//   - Debuggable → human-readable diagnostics
//   - Errorable  → explicit error state
//   - Rawable    → raw SQL fragment
//   - Renderable → dialect-agnostic rendering
//   - Stringable → concise string form
//   - Validable  → structural validity
//
// In addition, Join tokens expose token-specific accessors:
//
//   - Kind()      → the token kind (INNER, LEFT, RIGHT, FULL)
//   - Left()      → the left table operand
//   - Right()     → the right table operand
//   - Condition() → the ON condition expression
//
// Example:
//
//	j := token.NewInner("users", "orders", "users.id = orders.user_id")
//	fmt.Println(j.Type())      // INNER JOIN
//	fmt.Println(j.Left())      // users
//	fmt.Println(j.Right())     // orders
//	fmt.Println(j.Condition()) // users.id = orders.user_id
type Token interface {
	contract.Clonable[Token]
	contract.Debuggable
	contract.Errorable[Token]
	contract.Rawable
	contract.Renderable
	contract.Stringable
	contract.Validable

	// Kind reports the type of token (INNER, LEFT, RIGHT, FULL).
	Kind() join.Type

	// Left returns the left table operand.
	Left() table.Token

	// Right returns the right table operand.
	Right() table.Token

	// Condition returns the ON condition expression.
	Condition() string
}

// Ensure token implements the Token interface.
var _ Token = (*token)(nil)
