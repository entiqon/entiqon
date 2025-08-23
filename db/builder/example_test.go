// File: db/builder/example_test.go

package builder_test

import (
	"fmt"
	"log"

	"github.com/entiqon/entiqon/db/builder"
)

func ExampleSelectBuilder_basic() {
	sql, err := builder.NewSelect(nil).
		Fields("id, name").
		Source("users").
		Limit(10).
		Offset(20).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sql)
	// Output:
	// SELECT id, name FROM users LIMIT 10 OFFSET 20
}

func ExampleSelectBuilder_addFields() {
	sql, err := builder.NewSelect(nil).
		Fields("id").
		AddFields("name").
		Source("users").
		Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sql)
	// Output:
	// SELECT id, name FROM users
}

func ExampleSelectBuilder_emptyFields() {
	sql, err := builder.NewSelect(nil).
		Source("users").
		Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sql)
	// Output:
	// SELECT * FROM users
}

func ExampleSelectBuilder_invalidFields() {
	_, err := builder.NewSelect(nil).
		Fields(true).     // rejected
		AddFields(false). // rejected
		AddFields(123).   // rejected
		Source("users").
		Build()

	// Print the aggregated error
	if err != nil {
		fmt.Println(err.Error())
	}
	// Output:
	// ❌ [Build] - Invalid fields:
	//	⛔️ Field("true"): input type unsupported: bool
	//	⛔️ Field("false"): input type unsupported: bool
	//	⛔️ Field("123"): input type unsupported: int
}
