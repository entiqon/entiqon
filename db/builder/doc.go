// Package builder provides a fluent API to construct SQL SELECT queries.
//
// The SelectBuilder type allows incremental composition of SELECT statements
// with support for fields, sources, joins, conditions, grouping, having, ordering, limits, and offsets.
// It is simple, extensible, and dialect-aware.
//
// # Overview
//
// A builder is created with NewSelect. If no dialect is provided, BaseDialect
// is used by default.
//
//	sb := builder.NewSelect(nil)
//
// The builder supports two styles of adding fields:
//   - Fields(): initializes the field collection and replaces any existing fields.
//   - AddFields(): appends to the current field collection.
//
// If no fields are provided, Build() will fall back to `SELECT *`.
//
// # Methods
//
//   - Fields(...interface{}): reset and set fields for the SELECT clause.
//   - AddFields(...interface{}): append fields without resetting.
//   - Source(string): set the source table.
//   - InnerJoin(left, right, on): add INNER JOIN.
//   - LeftJoin(left, right, on): add LEFT JOIN.
//   - RightJoin(left, right, on): add RIGHT JOIN.
//   - FullJoin(left, right, on): add FULL JOIN.
//   - CrossJoin(left, right, on): add CROSS JOIN.
//   - NaturalJoin(left, right, on): add NATURAL JOIN.
//   - Where(...string): reset and set conditions for the WHERE clause.
//   - And(...string): append conditions with AND.
//   - Or(...string): append conditions with OR.
//   - GroupBy(...string): reset and set fields for the GROUP BY clause.
//   - ThenGroupBy(...string): append additional GROUP BY fields.
//   - Having(...string): reset and set conditions for the HAVING clause.
//   - AndHaving(...string): append conditions with AND.
//   - OrHaving(...string): append conditions with OR.
//   - OrderBy(...string): reset and set fields for the ORDER BY clause.
//   - ThenOrderBy(...string): append additional ORDER BY fields.
//   - Limit(int): apply LIMIT.
//   - Offset(int): apply OFFSET.
//   - Build(): construct the SQL string or return an error.
//   - String(): status output with a concise summary.
//   - Debug(): developer-facing detailed internal state.
//
// # Field Inputs
//
// Fields support the following forms:
//   - string: single expression, comma-separated list, inline alias via AS or space
//   - token.Field / *token.Field: field objects are added directly
//   - expr, alias
//   - expr, alias, isRaw
//
// Unsupported inputs are rejected and recorded as errored fields.
// Build() will collect all invalid fields and return a descriptive error.
//
// # Joins
//
// Join methods add JOIN clauses between the source table and another table.
//
//   - InnerJoin(left, right, on): INNER JOIN
//   - LeftJoin(left, right, on): LEFT JOIN
//   - RightJoin(left, right, on): RIGHT JOIN
//   - FullJoin(left, right, on): FULL JOIN
//   - CrossJoin(left, right, on): CROSS JOIN (Cartesian product)
//   - NaturalJoin(left, right, on): NATURAL JOIN (implicit column matching)
//
// Example:
//
//	sb := builder.NewSelect(nil).
//	    Fields("u.id").
//	    AddFields("o.id").
//	    AddFields("p.amount").
//	    Source("users u").
//	    InnerJoin("users u", "orders o", "u.id = o.user_id").
//	    LeftJoin("orders o", "payments p", "o.id = p.order_id").
//	    CrossJoin("orders o", "currencies c").
//	    NaturalJoin("departments d", "states s").
//	    Where("u.active = true").
//	    OrderBy("p.amount DESC").
//	    Limit(10).
//	    Offset(20)
//
//	sql, _ := sb.Build()
//	// SELECT u.id, o.id, p.amount
//	// FROM users u
//	// INNER JOIN orders o ON u.id = o.user_id
//	// LEFT JOIN payments p ON o.id = p.order_id
//	// CROSS JOIN currencies c
//	// NATURAL JOIN states s
//	// WHERE u.active = true
//	// ORDER BY p.amount DESC
//	// LIMIT 10
//	// OFFSET 20
//
// # Conditions
//
// Are expressed as raw strings. They are combined as follows:
//
//   - Where() clears existing conditions and sets new ones.
//   - And() appends conditions joined by AND.
//   - Or() appends conditions joined by OR.
//   - Multiple conditions in one call to Where() are normalized with AND.
//
// # Grouping
//
//   - GroupBy() clears existing grouping and sets new ones.
//   - ThenGroupBy() appends additional GROUP BY fields.
//
// # Having
//
//   - Having() clears existing conditions and sets new ones.
//   - AndHaving() appends conditions joined by AND.
//   - OrHaving() appends conditions joined by OR.
//
// # Ordering
//
//   - OrderBy() clears existing ordering and sets new ones.
//   - ThenOrderBy() appends additional ORDER BY fields.
//
// # Diagnostics
//
// - Debug(): verbose internal state, intended for developers.
// - String(): concise human-facing status, suitable for audit logs.
//
// # Error Handling
//
// Build() may return an error in the following cases:
//   - The builder receiver is nil.
//   - No source has been specified.
//   - One or more invalid fields were provided.
//
// # Example
//
//	sql, err := builder.NewSelect(nil).
//	    Fields("id, name").
//	    Source("users").
//	    Where("age > 18").
//	    GroupBy("department").
//	    Having("COUNT(*) > 5").
//	    OrderBy("created_at DESC").
//	    Limit(10).
//	    Offset(20).
//	    Build()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(sql)
//
// With no fields specified, the builder defaults to:
//
//	SELECT * FROM "users"
package builder
