package field_test

import (
	"fmt"

	"github.com/entiqon/db/token/field"
)

// ExampleField_Input shows that Input() preserves the original input string.
func ExampleToken_input() {
	f := field.New("SUM(qty) total")
	fmt.Println(f.Input())

	// Output: SUM(qty) total
}

//
//// ExampleField_Expr shows Expr() returns the parsed expression without alias.
//func ExampleToken_expr() {
//	f := field.New("SUM(qty) total")
//	fmt.Println(f.Expr())
//	// Output: SUM(qty)
//}
//
//// ExampleField_Alias shows Alias() returns the parsed alias.
//func ExampleToken_alias() {
//	f := field.New("SUM(qty) total")
//	fmt.Println(f.Alias())
//	// Output: total
//}
//
//// ExampleField_IsAliased shows IsAliased() is true if alias is set.
//func ExampleToken_isAliased() {
//	f := field.New("SUM(qty) total")
//	fmt.Println(f.IsAliased())
//	// Output: true
//}
//
//// ExampleField_IsRaw shows IsRaw() is true for computed expressions.
//func ExampleToken_isRaw() {
//	f := field.New("SUM(qty) total")
//	fmt.Println(f.IsRaw())
//	// Output: true
//}
//
////
//// Errorable contract
////
//
//// ExampleField_Error shows Error() returns nil for valid fields and non-nil for invalid ones.
//func ExampleToken_error() {
//	f := field.New("")
//	fmt.Println(f.Error())
//	// Output: empty expression is not allowed
//}
//
//// ExampleField_IsErrored shows IsErrored() is true when the field has an error.
//func ExampleToken_isErrored() {
//	f := field.New("")
//	fmt.Println(f.IsErrored())
//	// Output: true
//}
//
////
//// Token (ownership) contract
////
//
//// ExampleField_HasOwner shows HasOwner() reports if a field is bound to a table.
//func ExampleToken_HasOwner() {
//	f := field.New("id")
//	fmt.Println(f.HasOwner())
//	// Output: false
//}
//
//// ExampleField_Owner shows Owner() returns the bound owner if set.
//func ExampleToken_owner() {
//	f := field.NewWithTable("users", "id")
//	fmt.Println(*f.Owner())
//	// Output: users
//}
//
//// ExampleField_SetOwner shows SetOwner() attaches a table owner.
//func ExampleToken_setOwner() {
//	f := field.New("id")
//	owner := "orders"
//	f.SetOwner(&owner)
//	fmt.Println(*f.Owner())
//	// Output: orders
//}
//
////
//// Clonable contract
////
//
//// ExampleField_Clone shows Clone() returns a deep copy of the field.
//func ExampleToken_clone() {
//	orig := field.New("id user_id")
//	cl := orig.Clone()
//	fmt.Println(cl.Expr(), cl.Alias())
//	// Output: id user_id
//}
//
////
//// Debugging and logging
////
//
//// ExampleField_Debug shows Debug() returns a detailed diagnostic view.
//func ExampleToken_debug() {
//	f := field.New("COUNT(*) AS total")
//	fmt.Println(f.Debug())
//	// Output: ✅ field("COUNT(*) AS total"): [raw: false, aliased: true, errored: false]
//}
//
//// ExampleField_String shows String() returns concise logging output.
//func ExampleToken_string() {
//	f := field.New("")
//	fmt.Println(f.String())
//	// Output: ⛔ field(""): empty expression is not allowed
//}
//
////
//// Rawable and Renderable contracts
////
//
//// ExampleField_Raw shows Raw() returns SQL-generic rendering.
//func ExampleToken_raw() {
//	f := field.New("SUM(qty)", "total", true)
//	fmt.Println(f.Raw())
//	// Output: SUM(qty) AS total
//}
//
//// ExampleField_Render shows Render() produces dialect-agnostic SQL.
//func ExampleToken_render() {
//	f := field.NewWithTable("users", "id", "user_id")
//	fmt.Println(f.Render())
//	// Output: users.id AS user_id
//}
//
//// ExampleField_IsValid shows IsValid() reports validity.
//func ExampleToken_isValid() {
//	f := field.New("id")
//	fmt.Println(f.IsValid())
//	// Output: true
//}

func ExampleToken_new() {
	f := field.New()
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New("field alias extra")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New(field.New("field"))
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New(123456)
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New("")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New("*")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New("field")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New("field alias")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New("field AS alias")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New("field", "alias")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New("*", "alias")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New("field", 123456)
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	f = field.New("field", "123alias")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	// Aggregate function
	f = field.New("SUM(price) AS total")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	// Subquery
	f = field.New("(SELECT id FROM users) user_id")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	// Literal
	f = field.New("'hello' msg")
	fmt.Println(f.ExpressionKind(), f.Input(), f.Expr(), f.Alias(), f.IsAliased(), f.Error())

	// Output:
	// Invalid    false empty input is not allowed
	// Invalid field alias extra   false invalid identifier: field alias extra
	// Invalid Field("field")   false unsupported type; if you want to create a copy, use Clone() instead
	// Invalid 123456   false expr has invalid format (type int)
	// Invalid    false empty identifier is not allowed: ""
	// Identifier * *  false <nil>
	// Identifier field field  false <nil>
	// Identifier field alias field alias true <nil>
	// Identifier field AS alias field alias true <nil>
	// Identifier field alias field alias true <nil>
	// Identifier * alias * alias true '*' cannot be aliased or raw
	// Identifier field 123456 field  false alias must be a string, got int
	// Identifier field 123alias field  false invalid alias identifier cannot start with digit: "123alias"
	// Aggregate SUM(price) AS total SUM(price) total true <nil>
	// Subquery (SELECT id FROM users) user_id (SELECT id FROM users) user_id true <nil>
	// Literal 'hello' msg 'hello' msg true <nil>
}
