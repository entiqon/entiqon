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

func ExampleSelectBuilder_where() {
	sql, _ := builder.NewSelect(nil).
		Fields("id, name").
		Source("users").
		Where("age > 18", "status = 'active'").
		Build()

	fmt.Println(sql)
	// Output:
	// SELECT id, name FROM users WHERE age > 18 AND status = 'active'
}

func ExampleSelectBuilder_andOr() {
	sql, _ := builder.NewSelect(nil).
		Fields("id").
		Source("users").
		Where("age > 18").
		And("status = 'active'").
		Or("role = 'admin'").
		And("country = 'US'").
		Build()

	fmt.Println(sql)
	// Output:
	// SELECT id FROM users WHERE age > 18 AND status = 'active' OR role = 'admin' AND country = 'US'
}

func ExampleSelectBuilder_whereReset() {
	sql, _ := builder.NewSelect(nil).
		Fields("id").
		Source("users").
		Where("age > 18").
		And("status = 'active'").
		Where("role = 'admin'"). // resets previous conditions
		Build()

	fmt.Println(sql)
	// Output:
	// SELECT id FROM users WHERE role = 'admin'
}

// ExampleSelectBuilder_ordering demonstrates how to use OrderBy and ThenOrderBy
// to build an ORDER BY clause in a SELECT statement.
func ExampleSelectBuilder_ordering() {
	sql, _ := builder.NewSelect(nil).
		Fields("id, name").
		Source("users").
		OrderBy("created_at DESC").
		ThenOrderBy("id ASC").
		Build()

	fmt.Println(sql)
	// Output: SELECT id, name FROM users ORDER BY created_at DESC, id ASC
}

// ExampleSelectBuilder_grouping demonstrates how to use GroupBy and ThenGroupBy
// to build a GROUP BY clause in a SELECT statement.
func ExampleSelectBuilder_grouping() {
	sql, _ := builder.NewSelect(nil).
		Fields("department, COUNT(*) AS total").
		Source("users").
		GroupBy("department").
		ThenGroupBy("role").
		Build()

	fmt.Println(sql)
	// Output: SELECT department, COUNT(*) AS total FROM users GROUP BY department, role
}

// ExampleSelectBuilder_having demonstrates how to use Having, AndHaving,
// and OrHaving to build a HAVING clause in a SELECT statement.
func ExampleSelectBuilder_having() {
	sql, _ := builder.NewSelect(nil).
		Fields("department, COUNT(*) AS total").
		Source("users").
		GroupBy("department").
		Having("COUNT(*) > 5").
		AndHaving("AVG(age) > 30").
		OrHaving("SUM(salary) > 100000").
		Build()

	fmt.Println(sql)
	// Output:
	// SELECT department, COUNT(*) AS total FROM users GROUP BY department HAVING COUNT(*) > 5 AND AVG(age) > 30 OR SUM(salary) > 100000
}
