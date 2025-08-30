package identifier_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/token/types/identifier"
)

// ExampleType demonstrates how to use the identifier.Type enum.
func ExampleType() {
	var t identifier.Type

	// Functions render as capitalized names
	t = identifier.Function
	fmt.Println(t)

	// Subqueries render as capitalized names
	t = identifier.Subquery
	fmt.Println(t)

	// Invalid or unknown types render as "Invalid" or "Unknown"
	fmt.Println(identifier.Invalid)
	fmt.Println(identifier.Type(99))

	// Output:
	// Function
	// Subquery
	// Invalid
	// Unknown
}
