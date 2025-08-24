package contract_test

//import (
//	"fmt"
//
//	"github.com/entiqon/entiqon/db/contract"
//	"github.com/entiqon/entiqon/db/token/table"
//)
//
//// ExampleRenderable demonstrates using a Table as a Renderable.
//func ExampleRenderable() {
//	t := table.New("users", "u")
//	var r contract.Renderable = t
//	fmt.Println(r.Render())
//	// Output: users AS u
//}
//
//// ExampleRawable demonstrates using a Table as a Rawable.
//func ExampleRawable() {
//	t := table.New("users", "u")
//	var r contract.Rawable = t
//	fmt.Println(r.Raw())
//	// Output: users AS u
//}
//
//// ExampleStringable demonstrates using a Table as a Stringable.
//func ExampleStringable() {
//	t := table.New("users", "u")
//	var s contract.Stringable = t
//	fmt.Println(s.String())
//	// Output: ✅ Table(users AS u)
//}
//
//// ExampleDebuggable demonstrates using a Table as a Debuggable.
//func ExampleDebuggable() {
//	t := table.New("users", "u")
//	var d contract.Debuggable = t
//	fmt.Println(d.Debug())
//	// Output: ✅ Table("users AS u"): [raw:true, aliased:true, errored:false]
//}
//
//// ExampleClonable demonstrates using a Table as a Clonable.
//func ExampleClonable() {
//	t := table.New("users", "u")
//	var c contract.Clonable[*table.Table] = t
//	clone := c.Clone()
//	fmt.Println(clone.Render())
//	// Output: users AS u
//}
//
//// ExampleErrorable demonstrates using a Table as an Errorable.
//func ExampleErrorable() {
//	// invalid construction
//	t := table.New("users AS")
//
//	var e contract.Errorable = t
//	fmt.Println(e.IsErrored())
//	fmt.Println(e.Error())
//	// Output:
//	// true
//	// invalid format "users AS"
//}
