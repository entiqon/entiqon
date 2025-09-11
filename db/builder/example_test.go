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

	// Output: ❌ [Build] - Invalid fields:
	//	⛔️ field(""): expr has invalid format (type bool)
	//	⛔️ field(""): expr has invalid format (type bool)
	//	⛔️ field(""): expr has invalid format (type int)
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

func ExampleSelectBuilder_limit() {
	sql, _ := builder.NewSelect(nil).
		Source("table").
		Where("field = 1").
		Limit(1).
		Build()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM table WHERE field = 1 LIMIT 1
}

// --- JOIN Examples ---

func ExampleSelectBuilder_innerJoin() {
	sql, _ := builder.NewSelect(nil).
		Fields("u.id, o.total").
		Source("users u").
		InnerJoin("users", "orders o", "u.id = o.user_id").
		Build()

	fmt.Println(sql)
	// Output: SELECT u.id, o.total FROM users AS u INNER JOIN orders AS o ON u.id = o.user_id
}

func ExampleSelectBuilder_leftJoin() {
	sql, _ := builder.NewSelect(nil).
		Fields("u.id, o.total").
		Source("users u").
		LeftJoin("users", "orders o", "u.id = o.user_id").
		Build()

	fmt.Println(sql)
	// Output: SELECT u.id, o.total FROM users AS u LEFT JOIN orders AS o ON u.id = o.user_id
}

func ExampleSelectBuilder_rightJoin() {
	sql, _ := builder.NewSelect(nil).
		Fields("u.id, o.total").
		Source("users u").
		RightJoin("users", "orders o", "u.id = o.user_id").
		Build()

	fmt.Println(sql)
	// Output: SELECT u.id, o.total FROM users AS u RIGHT JOIN orders AS o ON u.id = o.user_id
}

func ExampleSelectBuilder_fullJoin() {
	sql, _ := builder.NewSelect(nil).
		Fields("u.id, o.total").
		Source("users u").
		FullJoin("users", "orders o", "u.id = o.user_id").
		Build()

	fmt.Println(sql)
	// Output: SELECT u.id, o.total FROM users AS u FULL JOIN orders AS o ON u.id = o.user_id
}

func ExampleSelectBuilder_crossJoin() {
	sql, _ := builder.NewSelect(nil).
		Fields("u.id, r.role").
		Source("users u").
		CrossJoin("users u", "roles r"). // condition ignored
		Build()

	fmt.Println(sql)
	// Output: SELECT u.id, r.role FROM users AS u CROSS JOIN roles AS r
}

func ExampleSelectBuilder_naturalJoin() {
	sql, _ := builder.NewSelect(nil).
		Fields("e.id, d.name").
		Source("employees e").
		NaturalJoin("departments d").
		Build()

	fmt.Println(sql)
	// Output: SELECT e.id, d.name FROM employees AS e NATURAL JOIN departments AS d
}

// ExampleSelectBuilder_getFields demonstrates how to access fields
// already added to the builder using GetFields().
func ExampleSelectBuilder_getFields() {
	sb := builder.NewSelect(nil).
		Fields("id").
		AddFields("name").
		Source("users")

	// Get the fields as a collection
	fields := sb.GetFields()
	for _, f := range fields.Items() {
		fmt.Print(f.Render(), " ")
	}

	// Output: id name
}

// ExampleSelectBuilder_debug demonstrates the developer-facing Debug()
// output, which includes internal state such as counts of clauses.
func ExampleSelectBuilder_debug() {
	sb := builder.NewSelect(nil).
		Fields("id").
		Source("users").
		Where("age > 18").
		OrderBy("id DESC")

	fmt.Println(sb.Debug())
	// Example output (counts may vary):
	// ✅ SelectBuilder{fields:1, source:users, where:1, groupBy:0, having:0, orderBy:1}
}

// ExampleSelectBuilder_string demonstrates the concise, human-facing
// status output of String(), suitable for logs or audits.
func ExampleSelectBuilder_string() {
	sb := builder.NewSelect(nil).
		Fields("id").
		Source("users").
		Where("age > 18")

	fmt.Println(sb.String())
	// Example output (status may vary):
	// ✅ SelectBuilder: status: ready, fields=1, conditions=1, grouped=false, sorted=false
}
