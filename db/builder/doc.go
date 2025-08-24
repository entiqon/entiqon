// Package builder provides a fluent API to construct SQL SELECT queries.
//
// The SelectBuilder type allows incremental composition of SELECT statements
// with support for fields, sources, conditions, ordering, limits, and offsets.
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
//   - Where(...string): reset and set conditions for the WHERE clause.
//   - And(...string): append conditions with AND.
//   - Or(...string): append conditions with OR.
//   - OrderBy(...string): reset and set fields for the ORDER BY clause.
//   - ThenOrderBy(...string): append additional ORDER BY fields.
//   - Limit(int): apply LIMIT.
//   - Offset(int): apply OFFSET.
//   - Build(): construct the SQL string or return an error.
//   - String(): status output with the built SQL or error.
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
// # Conditions
//
// Are expressed as raw strings. They are combined as follows:
//
//   - Where() clears existing conditions and sets new ones.
//   - And() appends conditions joined by AND.
//   - Or() appends conditions joined by OR.
//   - Multiple conditions in one call to Where() are normalized with AND.
//
// For example:
//
//	sb := builder.NewSelect(nil).
//		Fields("id, name").
//		Source("users").
//		Where("age > 18", "status = 'active'").
//		Or("role = 'admin'").
//		And("country = 'US'")
//
//	sql, _ := sb.Build()
//	// SELECT id, name FROM users WHERE age > 18 AND status = 'active' OR role = 'admin' AND country = 'US'
//
// # Ordering
//
// Ordering is expressed as raw strings (e.g. "created_at DESC"). They are combined as follows:
//
//   - OrderBy() clears existing ordering and sets new ones.
//   - ThenOrderBy() appends additional ORDER BY fields.
//   - Empty or whitespace-only strings are ignored.
//
// For example:
//
//	sb := builder.NewSelect(nil).
//		Fields("id, name").
//		Source("users").
//		OrderBy("created_at DESC").
//		ThenOrderBy("id ASC")
//
//	sql, _ := sb.Build()
//	// SELECT id, name FROM users ORDER BY created_at DESC, id ASC
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
//		Fields("id, name").
//		Source("users").
//		Where("age > 18").
//		OrderBy("created_at DESC").
//		Limit(10).
//		Offset(20).
//		Build()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(sql) // SELECT id, name FROM users WHERE age > 18 ORDER BY created_at DESC LIMIT 10 OFFSET 20
//
// With no fields specified, the builder defaults to:
//
//	SELECT * FROM "users"
package builder
