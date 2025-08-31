package condition_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/token/types/condition"
)

func ExampleType_isValid() {
	fmt.Println(condition.Single.IsValid())
	fmt.Println(condition.Type(99).IsValid())
	// Output:
	// true
	// false
}

func ExampleType_string() {
	fmt.Println(condition.Invalid.String())
	fmt.Println(condition.Single.String())
	fmt.Println(condition.And.String())
	fmt.Println(condition.Or.String())
	// Output:
	// Invalid
	//
	// AND
	// OR
}

func ExampleType_parseFrom() {
	// From int
	t := condition.ParseFrom(1)
	fmt.Println(t)

	// From valid string
	t = condition.ParseFrom("and")
	fmt.Println(t)

	// From invalid string
	t = condition.ParseFrom("foobar")
	fmt.Println(t)

	// Output:
	//
	// AND
	// Invalid
}
