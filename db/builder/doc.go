// Package builder provides a fluent API to construct SQL SELECT queries.
//
// The SelectBuilder type allows incremental composition of SELECT statements
// with support for fields, sources, limits, and offsets. It is intended to be
// simple, extensible, and dialect-aware.
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
//   - Fields(...interface{}): reset and set fields for the select clause.
//   - AddFields(...interface{}): append fields without resetting.
//   - Source(string): set the source table.
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
//		Limit(10).
//		Offset(20).
//		Build()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(sql) // SELECT id, name FROM users LIMIT 10 OFFSET 20
//
// With no fields specified, the builder defaults to:
//
//	SELECT * FROM "users"
package builder
