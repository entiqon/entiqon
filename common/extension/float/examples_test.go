package float_test

import (
	"fmt"

	"github.com/entiqon/entiqon/common/extension/float"
)

func ExampleParseFrom_string() {
	f, _ := float.ParseFrom("123.456")
	fmt.Println(f)
	// Output: 123.456
}

func ExampleParseFrom_int() {
	f, _ := float.ParseFrom(42)
	fmt.Println(f)
	// Output: 42
}

func ExampleParseFrom_bool() {
	f, _ := float.ParseFrom(true)
	fmt.Println(f)
	// Output: 1
}

func ExampleParseFrom_pointer() {
	i := 7
	f, _ := float.ParseFrom(&i)
	fmt.Println(f)
	// Output: 7
}
