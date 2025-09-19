// File: db/builder/selects/example_test.go

package selects_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/builder/selects"
	"github.com/entiqon/entiqon/db/token/types/operator"
)

func ExampleSelectBuilder_fields() {
	sb := selects.New(nil).
		Fields("id").
		AppendFields("name", "username"). // expr + alias (2-args)
		From("users")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT id, name AS username FROM users
}

func ExampleSelectBuilder_appendFields() {
	sb := selects.New(nil).
		From("users").
		Fields("id").
		AppendFields("created_at")

	sql, _, _ := sb.Build()

	fmt.Println(sql)
	// Output: SELECT id, created_at FROM users
}

func ExampleSelectBuilder_from() {
	sb := selects.New(nil).
		Fields("id").
		From("users u")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT id FROM users AS u
}

func ExampleSelectBuilder_innerJoin() {
	sb := selects.New(nil).
		Fields("u.id, o.id").
		From("users u").
		InnerJoin("users u", "orders o", "u.id = o.user_id")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT u.id, o.id FROM users AS u INNER JOIN orders AS o ON u.id = o.user_id
}

func ExampleSelectBuilder_leftJoin() {
	sb := selects.New(nil).
		From("users u").
		LeftJoin("users u", "profiles p", "u.id = p.user_id")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM users AS u LEFT JOIN profiles AS p ON u.id = p.user_id
}

func ExampleSelectBuilder_rightJoin() {
	sb := selects.New(nil).
		From("orders o").
		RightJoin("orders o", "payments p", "o.id = p.order_id")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM orders AS o RIGHT JOIN payments AS p ON o.id = p.order_id
}

func ExampleSelectBuilder_fullJoin() {
	sb := selects.New(nil).
		From("customers c").
		FullJoin("customers c", "subscriptions s", "c.id = s.customer_id")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM customers AS c FULL JOIN subscriptions AS s ON c.id = s.customer_id
}

func ExampleSelectBuilder_crossJoin() {
	sb := selects.New(nil).
		From("users u").
		CrossJoin("roles r")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM users AS u CROSS JOIN roles AS r
}

func ExampleSelectBuilder_naturalJoin() {
	sb := selects.New(nil).
		From("employees e").
		NaturalJoin("departments d")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM employees AS e NATURAL JOIN departments AS d
}

func ExampleSelectBuilder_where() {
	sb := selects.New(nil).
		From("users").
		Where("age", operator.GreaterThan, 18)

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM users WHERE age > :age
}

func ExampleSelectBuilder_andWhere() {
	sb := selects.New(nil).
		From("users").
		Where("active = true").
		AndWhere("country = 'USA'")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM users WHERE active = :active AND country = :country
}

func ExampleSelectBuilder_orWhere() {
	sb := selects.New(nil).
		From("users").
		Where("active = true").
		OrWhere("country = 'USA'")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM users WHERE active = :active OR country = :country
}

func ExampleSelectBuilder_groupBy() {
	sb := selects.New(nil).
		Fields("COUNT(id)", "department").
		From("users").
		GroupBy("department")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT COUNT(id) AS department FROM users GROUP BY department
}

func ExampleSelectBuilder_thenGroupBy() {
	sb := selects.New(nil).
		Fields("department, role, COUNT(id) AS collaborators").
		From("users").
		GroupBy("department").
		ThenGroupBy("role")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT department, role, COUNT(id) AS collaborators FROM users GROUP BY department, role
}

func ExampleSelectBuilder_orderBy() {
	sb := selects.New(nil).
		Fields("id, name").
		From("users").
		OrderBy("created_at DESC")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT id, name FROM users ORDER BY created_at DESC
}

func ExampleSelectBuilder_thenOrderBy() {
	sb := selects.New(nil).
		Fields("id, name").
		From("users").
		OrderBy("created_at DESC").
		ThenOrderBy("id")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT id, name FROM users ORDER BY created_at DESC, id
}

func ExampleSelectBuilder_having() {
	sb := selects.New(nil).
		Fields("department_id, COUNT(id)").
		From("users").
		GroupBy("department_id").
		Having("COUNT(id) > 5")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT department_id, COUNT(id) FROM users GROUP BY department_id HAVING COUNT(id) > 5
}

func ExampleSelectBuilder_andHaving() {
	sb := selects.New(nil).
		Fields("department_id, COUNT(id)").
		From("users").
		GroupBy("department_id").
		Having("COUNT(id) > 5").
		AndHaving("COUNT(id) < 100")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT department_id, COUNT(id) FROM users GROUP BY department_id HAVING COUNT(id) > 5 AND COUNT(id) < 100
}

func ExampleSelectBuilder_orHaving() {
	sb := selects.New(nil).
		Fields("department_id, COUNT(id)").
		From("users").
		GroupBy("department_id").
		Having("COUNT(id) > 5").
		OrHaving("COUNT(id) = 1")

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT department_id, COUNT(id) FROM users GROUP BY department_id HAVING COUNT(id) > 5 OR COUNT(id) = 1
}

func ExampleSelectBuilder_take() {
	sb := selects.New(nil).
		From("users").
		Take(10)

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM users LIMIT 10
}

func ExampleSelectBuilder_skip() {
	sb := selects.New(nil).
		From("users").
		Skip(20)

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM users OFFSET 20
}

func ExampleSelectBuilder_pagination() {
	sb := selects.New(nil).
		From("users").
		Take(10).
		Skip(20)

	sql, _, _ := sb.Build()
	fmt.Println(sql)
	// Output: SELECT * FROM users LIMIT 10 OFFSET 20
}
