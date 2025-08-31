// Package join provides low-level primitives to represent SQL JOIN clauses
// in a dialect-agnostic way.
//
// # Overview
//
// The join package defines the enumeration Type, which classifies the supported
// SQL JOIN operations. Each Type is safe to use in builders and renders to the
// canonical SQL92 keyword form.
//
// Supported join types include:
//   - INNER JOIN
//   - LEFT JOIN
//   - RIGHT JOIN
//   - FULL JOIN
//   - CROSS JOIN
//   - NATURAL JOIN
//
// # Usage
//
// A Type can be created directly or parsed from user input:
//
//	j := join.Inner
//	fmt.Println(j.String()) // "INNER JOIN"
//
//	j2 := join.ParseFrom("cross join")
//	if !j2.IsValid() {
//		log.Fatal("invalid join type")
//	}
//	fmt.Println(j2) // "CROSS JOIN"
//
// # Design
//
//   - Type is defined as an integer enum to provide type safety.
//   - Rendering uses String() for canonical keywords.
//   - Invalid or unrecognized values are represented by Type Invalid.
//
// The join package is consumed by higher-level builders such as SelectBuilder,
// which combine Table, Field, Condition, and Join tokens into complete SQL
// statements.
package join
