package decimal_test

import (
	"fmt"

	"github.com/entiqon/entiqon/common/extension/decimal"
)

func ExampleParseFrom_rounding() {
	v, _ := decimal.ParseFrom("123.456789", 2)
	fmt.Println(v)
	// Output: 123.46
}

func ExampleParseFrom_integer() {
	v, _ := decimal.ParseFrom(42, 0)
	fmt.Println(v)
	// Output: 42
}

func ExampleParseFrom_bool() {
	v, _ := decimal.ParseFrom(true, 3)
	fmt.Println(v)
	// Output: 1
}

func ExampleParseFrom_invalid() {
	_, err := decimal.ParseFrom("invalid", 2)
	fmt.Println(err != nil)
	// Output: true
}
