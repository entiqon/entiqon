package token_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/token"
)

// ExampleNewField demonstrates the basic instantiation forms of Field.
func ExampleNewField() {
	// Single expression
	f1 := token.NewField("id")

	// Expression with alias (space-separated)
	f2 := token.NewField("id user_id")

	// Expression with alias (AS keyword)
	f3 := token.NewField("COUNT(*) AS total")

	// Expression and alias separately
	f4 := token.NewField("name", "username")

	// Expression and alias with raw flag
	f5 := token.NewField("JSON_EXTRACT(data, '$.name')", "extracted", true)

	fmt.Println(f1.Render())
	fmt.Println(f2.Render())
	fmt.Println(f3.Render())
	fmt.Println(f4.Render())
	fmt.Println(f5.Render())

	// Output:
	// id
	// id AS user_id
	// COUNT(*) AS total
	// name AS username
	// JSON_EXTRACT(data, '$.name') AS extracted
}

// ExampleField_clone demonstrates the Clone method.
func ExampleField_clone() {
	original := token.NewField("id", "user_id")
	clone := original.Clone()

	fmt.Println(original.Render())
	fmt.Println(clone.Render())
	fmt.Println(original == clone) // should be false, deep copy

	// Output:
	// id AS user_id
	// id AS user_id
	// false
}

// ExampleField_invalid demonstrates invalid field creation.
func ExampleField_invalid() {
	// Passing an unsupported type (e.g. bool) results in an errored field.
	f := token.NewField(true)

	if f.IsErrored() {
		fmt.Println(f.String()) // concise log view
		fmt.Println(f.Debug())  // detailed diagnostic view
	}

	// Output:
	// ⛔️ Field("true"): input type unsupported: bool
	// ⛔️ Field("true"): [raw: false, aliased: false, errored: true] – input type unsupported: bool
}

// ExampleField_debug demonstrates valid fields with Debug output.
func ExampleField_debug() {
	f := token.NewField("COUNT(*) AS total")
	fmt.Println(f.String())
	fmt.Println(f.Debug())

	// Output:
	// ✅ Field("COUNT(*) AS total")
	// ✅ Field("COUNT(*) AS total"): [raw: true, aliased: true, errored: false]
}
