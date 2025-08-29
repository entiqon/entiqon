package helpers_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/token/helpers"
)

// ExampleIsValidIdentifier shows how to validate strings as SQL identifiers.
func ExampleIsValidIdentifier() {
	// Valid identifiers
	fmt.Println(helpers.IsValidIdentifier("id"))
	fmt.Println(helpers.IsValidIdentifier("user_id"))
	fmt.Println(helpers.IsValidIdentifier("_col123"))

	// Invalid identifiers
	fmt.Println(helpers.IsValidIdentifier("123abc"))
	fmt.Println(helpers.IsValidIdentifier("first-name"))
	fmt.Println(helpers.IsValidIdentifier(""))

	// Output:
	// true
	// true
	// true
	// false
	// false
	// false
}
