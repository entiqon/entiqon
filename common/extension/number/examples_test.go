package number_test

import (
	"fmt"

	"github.com/entiqon/entiqon/common/extension/number"
)

func ExampleParseFrom_round() {
	n, _ := number.ParseFrom("123.6", true)
	fmt.Println(n)
	// Output: 124
}

func ExampleParseFrom_strict() {
	_, err := number.ParseFrom(123.4, false)
	fmt.Println(err != nil)
	// Output: true
}

func ExampleParseFrom_bool() {
	n, _ := number.ParseFrom(true, false)
	fmt.Println(n)
	// Output: 1
}

func ExampleParseFrom_int() {
	n, _ := number.ParseFrom(42, false)
	fmt.Println(n)
	// Output: 42
}
