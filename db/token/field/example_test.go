package field_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/token/field"
)

//
// BaseToken contract
//

// ExampleField_Input shows that Input() preserves the original input string.
func ExampleToken_input() {
	f := field.New("SUM(qty) total")
	fmt.Println(f.Input())
	// Output: SUM(qty) total
}

// ExampleField_Expr shows Expr() returns the parsed expression without alias.
func ExampleToken_expr() {
	f := field.New("SUM(qty) total")
	fmt.Println(f.Expr())
	// Output: SUM(qty)
}

// ExampleField_Alias shows Alias() returns the parsed alias.
func ExampleToken_alias() {
	f := field.New("SUM(qty) total")
	fmt.Println(f.Alias())
	// Output: total
}

// ExampleField_IsAliased shows IsAliased() is true if alias is set.
func ExampleToken_isAliased() {
	f := field.New("SUM(qty) total")
	fmt.Println(f.IsAliased())
	// Output: true
}

// ExampleField_IsRaw shows IsRaw() is true for computed expressions.
func ExampleToken_isRaw() {
	f := field.New("SUM(qty) total")
	fmt.Println(f.IsRaw())
	// Output: true
}

//
// Errorable contract
//

// ExampleField_Error shows Error() returns nil for valid fields and non-nil for invalid ones.
func ExampleToken_error() {
	f := field.New("")
	fmt.Println(f.Error())
	// Output: empty expression is not allowed
}

// ExampleField_IsErrored shows IsErrored() is true when the field has an error.
func ExampleToken_isErrored() {
	f := field.New("")
	fmt.Println(f.IsErrored())
	// Output: true
}

//
// Token (ownership) contract
//

// ExampleField_HasOwner shows HasOwner() reports if a field is bound to a table.
func ExampleToken_HasOwner() {
	f := field.New("id")
	fmt.Println(f.HasOwner())
	// Output: false
}

// ExampleField_Owner shows Owner() returns the bound owner if set.
func ExampleToken_owner() {
	f := field.NewWithTable("users", "id")
	fmt.Println(*f.Owner())
	// Output: users
}

// ExampleField_SetOwner shows SetOwner() attaches a table owner.
func ExampleToken_setOwner() {
	f := field.New("id")
	owner := "orders"
	f.SetOwner(&owner)
	fmt.Println(*f.Owner())
	// Output: orders
}

//
// Clonable contract
//

// ExampleField_Clone shows Clone() returns a deep copy of the field.
func ExampleToken_clone() {
	orig := field.New("id user_id")
	cl := orig.Clone()
	fmt.Println(cl.Expr(), cl.Alias())
	// Output: id user_id
}

//
// Debugging and logging
//

// ExampleField_Debug shows Debug() returns a detailed diagnostic view.
func ExampleToken_debug() {
	f := field.New("COUNT(*) AS total")
	fmt.Println(f.Debug())
	// Output: ✅ field("COUNT(*) AS total"): [raw: true, aliased: true, errored: false]
}

// ExampleField_String shows String() returns concise logging output.
func ExampleToken_string() {
	f := field.New("")
	fmt.Println(f.String())
	// Output: ⛔ field(""): empty expression is not allowed
}

//
// Rawable and Renderable contracts
//

// ExampleField_Raw shows Raw() returns SQL-generic rendering.
func ExampleToken_raw() {
	f := field.New("SUM(qty)", "total", true)
	fmt.Println(f.Raw())
	// Output: SUM(qty) AS total
}

// ExampleField_Render shows Render() produces dialect-agnostic SQL.
func ExampleToken_render() {
	f := field.NewWithTable("users", "id", "user_id")
	fmt.Println(f.Render())
	// Output: users.id AS user_id
}

// ExampleField_IsValid shows IsValid() reports validity.
func ExampleToken_isValid() {
	f := field.New("id")
	fmt.Println(f.IsValid())
	// Output: true
}
