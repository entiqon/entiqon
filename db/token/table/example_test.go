package table_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/token/table"
)

// ExampleNew_plain demonstrates creating a plain table.
func ExampleNew_plain() {
	t := table.New("users")
	fmt.Println(t.Render())
	// Output: users
}

// ExampleNew_inlineAlias demonstrates inline aliasing.
func ExampleNew_inlineAlias() {
	t := table.New("users u")
	fmt.Println(t.Render())
	// Output: users AS u
}

// ExampleNew_explicitAlias demonstrates explicit aliasing with two arguments.
func ExampleNew_explicitAlias() {
	t := table.New("users", "u")
	fmt.Println(t.Render())
	// Output: users AS u
}

// ExampleNew_subquery demonstrates creating a subquery table with alias.
func ExampleNew_subquery() {
	t := table.New("(SELECT COUNT(id) FROM users)", "t")
	fmt.Println(t.Render())
	// Output: (SELECT COUNT(id) FROM users) AS t
}

// ExampleTable_Raw demonstrates the Raw() method.
func ExampleTable_Raw() {
	t := table.New("users u")
	fmt.Println(t.Raw())
	// Output: users AS u
}

// ExampleTable_IsRaw demonstrates the IsRaw() method.
func ExampleTable_IsRaw() {
	t1 := table.New("users u")
	fmt.Println(t1.IsRaw())

	t2 := table.New("(SELECT COUNT(id) FROM users)", "t")
	fmt.Println(t2.IsRaw())
	// Output:
	// false
	// true
}

// ExampleTable_IsAliased demonstrates the IsAliased() method.
func ExampleTable_IsAliased() {
	t := table.New("users u")
	fmt.Println(t.IsAliased())
	// Output: true
}

// ExampleTable_IsValid demonstrates the IsValid() method.
func ExampleTable_IsValid() {
	t := table.New("users u")
	fmt.Println(t.IsValid())

	bad := table.New("users AS") // invalid
	fmt.Println(bad.IsValid())
	// Output:
	// true
	// false
}

// ExampleTable_String demonstrates String() output for logging.
func ExampleTable_String() {
	t := table.New("users", "u")
	fmt.Println(t.String())
	// Output: ✅ Table(users AS u)
}

// ExampleTable_Debug demonstrates Debug() output for diagnostics.
func ExampleTable_Debug() {
	t := table.New("users u")
	fmt.Println(t.Debug())
	// Output: ✅ Table("users u"): [raw:false, aliased:true, errored:false]
}

// ExampleTable_Error demonstrates handling of invalid input.
func ExampleTable_Error() {
	t := table.New("users AS") // invalid alias
	fmt.Println(t.String())
	fmt.Println(t.IsErrored())
	fmt.Println(t.Error())
	// Output:
	// ❌ Table("users AS"): invalid format "users AS"
	// true
	// invalid format "users AS"
}

// ExampleTable_Clone demonstrates the Clone() method.
func ExampleTable_Clone() {
	t := table.New("users u")
	clone := t.Clone()
	fmt.Println(clone.Render())
	// Output: users AS u
}
