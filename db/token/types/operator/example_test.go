package operator_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/token/types/operator"
)

func ExampleType_Alias() {
	fmt.Println(fmt.Sprintf("[%s]: alias=%s", operator.Equal, operator.Equal.Alias()))
	op := operator.ParseFrom(false)
	fmt.Println(fmt.Sprintf("[%s]: alias=%s", op, op.Alias()))

	// Output:
	// [=]: alias=eq
	// [Invalid]: alias=invalid
}

func ExampleType_IsValid() {
	fmt.Println(fmt.Sprintf("[%s]: valid=%t", operator.Equal, operator.Equal.IsValid()))

	op := operator.ParseFrom(1)
	fmt.Println(fmt.Sprintf("[%s]: valid=%t", op, op.IsValid()))

	// Output:
	// [=]: valid=true
	// [Invalid]: valid=false
}

func ExampleType_String() {
	fmt.Println(fmt.Sprintf("[%s]: valid=%t", operator.Equal, operator.Equal.IsValid()))
	fmt.Println(fmt.Sprintf("[%s]: valid=%t", operator.Between, operator.Between.IsValid()))

	// Output:
	// [=]: valid=true
	// [BETWEEN]: valid=true
}

func ExampleGetKnownOperators() {
	fmt.Println(len(operator.GetKnownOperators()), "active operators")

	// Output:
	// 15 active operators
}

func ExampleParseFrom() {
	fmt.Println(operator.ParseFrom(1))
	fmt.Println(operator.ParseFrom(operator.Equal))
	fmt.Println(operator.ParseFrom("lte"))
	fmt.Println(operator.ParseFrom("nin"))
	fmt.Println(operator.ParseFrom("is null"))
	fmt.Println(operator.ParseFrom("IN"))
	fmt.Println(operator.ParseFrom([]byte(">=")))

	// Output:
	// Invalid
	// =
	// <=
	// NOT IN
	// IS NULL
	// IN
	// >=
}
