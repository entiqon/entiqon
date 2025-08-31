package contract_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/contract"
	"github.com/entiqon/entiqon/db/token/field"
	"github.com/entiqon/entiqon/db/token/table"
)

// ExampleBaseToken demonstrates using a Field as a BaseToken.
func ExampleIdentifiable() {
	t := field.New("users", "u")
	var bt contract.Identifiable = t
	fmt.Println(fmt.Sprintf(
		"Input=%q, Expr=%q", bt.Input(), bt.Expr()))
	// Output: Input="users u", Expr="users"
}

// ExampleBaseToken demonstrates using a Field as a BaseToken.
func ExampleBaseToken() {
	t := field.New("users", "u")
	var bt contract.BaseToken = t
	fmt.Println(fmt.Sprintf(
		"Input=%q, Expr=%q, Alias=%q, Aliased=%t",
		bt.Input(), bt.Expr(), bt.Alias(), bt.IsAliased()))
	// Output: Input="users u", Expr="users", Alias="u", Aliased=true
}

// ExampleClonable demonstrates using a Table as a Clonable.
func ExampleClonable() {
	t := table.New("users", "u")
	var c contract.Clonable[table.Token] = t
	clone := c.Clone()
	fmt.Println(clone.Render())
	// Output: users AS u
}

// ExampleDebuggable demonstrates using a Table as a Debuggable.
func ExampleDebuggable() {
	t := table.New("users", "u")
	var d contract.Debuggable = t
	fmt.Println(d.Debug())
	// Output:
	// ✅ Table("users u"): [raw:false, aliased:true, errored:false]
}

// ExampleErrorable demonstrates using a Table as an Errorable.
func ExampleErrorable() {
	// invalid construction
	t := table.New("users AS")

	var e contract.Errorable[table.Token] = t
	fmt.Println(e.IsErrored())
	fmt.Println(e.Error())

	// manually mark an otherwise valid table as errored
	valid := table.New("products")
	valid.SetError(fmt.Errorf("manual mark as errored"))
	fmt.Println(valid.IsErrored())
	fmt.Println(valid.Error())

	// Output:
	// true
	// invalid alias: AS
	// true
	// manual mark as errored
}

// ExampleRawable demonstrates using a Table as a Rawable.
func ExampleRawable() {
	t := table.New("users", "u")
	var r contract.Rawable = t
	fmt.Println(r.Raw())
	// Output: users AS u
}

// ExampleRenderable demonstrates using a Table as a Renderable.
func ExampleRenderable() {
	t := table.New("users", "u")
	var r contract.Renderable = t
	fmt.Println(r.Render())
	// Output: users AS u
}

// ExampleStringable demonstrates using a Table as a Stringable.
func ExampleStringable() {
	t := table.New("users", "u")
	var s contract.Stringable = t
	fmt.Println(s.String())
	// Output: ✅ Table(users AS u)
}

// ExampleValidable demonstrates using a Table as a Validable.
func ExampleValidable() {
	t := table.New("users", "u")
	var v contract.Validable = t
	fmt.Println(v.IsValid())
	// Output: true
}
