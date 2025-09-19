// Package builder provides modular, dialect-aware SQL builders.
//
// # Overview
//
// The builder package offers a consistent, chainable API for constructing SQL
// statements across multiple dialects. Each builder is scoped to a specific
// SQL verb, with its own package:
//
//   - selects — SELECT queries
//   - inserts — INSERT queries
//   - updates — UPDATE queries
//   - deletes — DELETE queries
//   - merge   — MERGE queries
//
// Builders enforce consistent parsing, validation, and rendering rules, while
// preserving Go idioms. Each mutator returns the builder for chaining, while
// accessors expose the current state (fields, joins, groupings, etc.).
//
// # Philosophy
//
//   - Fluent: Builders are chainable and self-explanatory.
//   - Safe: Invalid expressions are carried as tokens, surfaced at Build.
//   - Modular: Each SQL verb has its own package and contract.
//   - Dialect-aware: BaseDialect included; dialects can override rendering.
//   - Testable: 100% coverage enforced across all builders.
//
// # Quick Example
//
//	sb := selects.New(nil).
//	    Fields("id", "name").
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
// # See also
//
// For detailed usage and examples, refer to each builder’s documentation:
//   - selects
//   - inserts
//   - updates
//   - deletes
//   - merge
package builder
