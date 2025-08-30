package helpers_test

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/token/helpers"
	"github.com/entiqon/entiqon/db/token/types/identifier"
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

// ExampleIsValidIdentifier demonstrates quick true/false checks.
func ExampleIsValidAlias() {
	fmt.Println(helpers.IsValidAlias("user_id"))
	fmt.Println(helpers.IsValidAlias("123abc"))

	// Output:
	// true
	// false
}

// ExampleValidateTrailingAlias demonstrates extracting a trailing alias
// when no explicit AS is present.
func ExampleValidateTrailingAlias() {
	// Valid trailing alias
	alias, err := helpers.ValidateTrailingAlias("(SELECT * FROM users) u")
	fmt.Println(alias, err == nil)

	// Invalid: reserved keyword as alias
	alias, err = helpers.ValidateTrailingAlias("(SELECT * FROM users) SELECT")
	fmt.Println(alias, err)

	// Invalid: explicit AS → not considered trailing alias
	alias, err = helpers.ValidateTrailingAlias("(SELECT * FROM users) AS u")
	fmt.Println(alias, err)

	// Output:
	// u true
	//  invalid trailing alias "SELECT": alias is a reserved keyword: "SELECT"
	//  explicit AS found, not a trailing alias
}

// ExampleGenerateAlias demonstrates generating a deterministic alias
// from a prefix and expression string.
func ExampleGenerateAlias() {
	// Function expression with "fn" prefix
	got := helpers.GenerateAlias(identifier.Function.Alias(), "SUM(price)")
	fmt.Println(fmt.Sprintf(
		"Contains(fn)=%t, Length=%d",
		strings.Contains(got, "fn"),
		len(got),
	))

	// Subquery expression with "sq" prefix
	got = helpers.GenerateAlias(identifier.Subquery.Alias(), "(SELECT * FROM users)")
	fmt.Println(fmt.Sprintf(
		"Contains(sq)=%t, Length=%d",
		strings.Contains(got, "fn"),
		len(got),
	))

	// Output:
	// Contains(fn)=true, Length=9
	// Contains(sq)=false, Length=9
}

// ExampleValidateWildcard demonstrates validating use of the "*"
// wildcard in field expressions.
func ExampleValidateWildcard() {
	// Valid: bare "*" without alias
	fmt.Println(helpers.ValidateWildcard("*", ""))

	// Invalid: "*" aliased → not allowed
	fmt.Println(helpers.ValidateWildcard("*", "total"))

	// Not a wildcard → ignored by this helper
	fmt.Println(helpers.ValidateWildcard("id", "alias"))

	// Output:
	// <nil>
	//'*' cannot be aliased or raw
	//<nil>
}

// ExampleClassifyExpression demonstrates classifying raw SQL expressions
// into identifier.Type categories.
func ExampleResolveExpressionType() {
	fmt.Println(helpers.ResolveExpressionType("(SELECT * FROM users)"))      // Subquery
	fmt.Println(helpers.ResolveExpressionType("(a+b)"))                      // Computed
	fmt.Println(helpers.ResolveExpressionType("SUM(price)"))                 // Aggregate
	fmt.Println(helpers.ResolveExpressionType("JSON_EXTRACT(data, '$.id')")) // Function
	fmt.Println(helpers.ResolveExpressionType("'abc'"))                      // Literal
	fmt.Println(helpers.ResolveExpressionType("users"))                      // Identifier
	fmt.Println(helpers.ResolveExpressionType(""))                           // Invalid

	// Output:
	// Subquery
	// Computed
	// Aggregate
	// Function
	// Literal
	// Identifier
	// Invalid
}
