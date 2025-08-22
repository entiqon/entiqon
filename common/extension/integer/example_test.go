package integer_test

import (
	"fmt"

	"github.com/entiqon/entiqon/common/extension/integer"
)

// ExampleParseFrom demonstrates using integer.ParseFrom
// with different types of inputs.
func ExampleParseFrom() {
	// From int
	i1, _ := integer.ParseFrom(42)
	fmt.Println(i1)

	// From float (truncated)
	i2, _ := integer.ParseFrom(3.99)
	fmt.Println(i2)

	// From string
	i3, _ := integer.ParseFrom("123")
	fmt.Println(i3)

	// From bool
	i4, _ := integer.ParseFrom(true)
	fmt.Println(i4)

	// Output:
	// 42
	// 3
	// 123
	// 1
}
