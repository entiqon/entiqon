// Package selects provides a builder for SQL SELECT statements.
//
// # Overview
//
// SelectBuilder constructs SELECT queries with support for:
//
//   - Fields (columns, expressions, aliases)
//   - Source tables (FROM)
//   - Joins (INNER, LEFT, RIGHT, FULL, CROSS, NATURAL)
//   - Conditions (WHERE)
//   - Grouping (GROUP BY)
//   - Filtering (HAVING)
//   - Sorting (ORDER BY)
//   - Pagination (LIMIT and OFFSET)
//
// # Example
//
//	sb := selects.New(nil).
//	    Fields("id", "name AS username").
//	    From("users u").
//	    Where("u.active = true").
//	    OrderBy("created_at DESC").
//	    Take(10)
//
//	sql, args, err := sb.Build()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(sql, args)
//
// # Notes
//
//   - Mutators return the builder for chaining.
//   - Accessors expose the current state (fields, joins, etc.).
//   - Invalid tokens are carried forward and surfaced at Build.
//   - Passing nil as dialect defaults to BaseDialect.
package selects
