package helpers_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/token/helpers"
)

// ExampleIsValidIdentifier demonstrates quick true/false checks.
func ExampleIsValidIdentifier() {
	fmt.Println(helpers.IsValidIdentifier("user_id"))
	fmt.Println(helpers.IsValidIdentifier("123abc"))

	// Output:
	// true
	// false
}

// ExampleValidateIdentifier demonstrates detailed validation errors.
func ExampleValidateIdentifier() {
	// Valid identifier → nil error
	fmt.Println(helpers.ValidateIdentifier("user_id"))

	// Empty identifier
	fmt.Println(helpers.ValidateIdentifier(""))

	// Starts with digit
	fmt.Println(helpers.ValidateIdentifier("123abc"))

	// Invalid syntax (dash)
	fmt.Println(helpers.ValidateIdentifier("user-name"))

	// Non-ASCII identifiers (strict mode rejects them)
	fmt.Println(helpers.ValidateIdentifier("café"))
	fmt.Println(helpers.ValidateIdentifier("mañana"))
	fmt.Println(helpers.ValidateIdentifier("niño"))

	// Output:
	// <nil>
	// identifier cannot be empty
	// identifier cannot start with digit: "123abc"
	// invalid identifier syntax: "user-name"
	// invalid identifier syntax: "café"
	// invalid identifier syntax: "mañana"
	// invalid identifier syntax: "niño"
}
