// Package condition provides primitives to represent SQL conditional
// expressions (WHERE, HAVING, ON).
//
// # Overview
//
// A condition expresses a boolean predicate used to filter or join rows.
// Unlike fields or tables, conditions combine expressions with values
// and logical operators (AND, OR, etc.).
//
// The package defines a Type enum and associated helpers to classify
// condition tokens:
//
//   - Single: a single expression, such as "id = 1" or "age > 30".
//   - And:    a composite condition joined with logical AND.
//   - Or:     a composite condition joined with logical OR.
//   - Invalid: an unrecognized or unsupported condition type.
//
// # Example
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/entiqon/db/token/condition"
//	)
//
//	func main() {
//	    c := condition.Single
//	    fmt.Println(c, c.IsValid()) // Output: "" true
//
//	    c2 := condition.Type(99)
//	    fmt.Println(c2, c2.IsValid()) // Output: Invalid false
//
//	    c3 := condition.And
//	    fmt.Println(c3.String()) // Output: AND
//	}
//
// # Integration
//
// Condition types are consumed by higher-level builders (e.g. SelectBuilder)
// to render WHERE and HAVING clauses, or by join tokens to render ON clauses.
// Conditions will later be extended to support structured construction,
// parameter binding, and logical composition (e.g. condition.And(c1, c2)).
//
// # Contracts
//
// Condition tokens implement the Kindable[Type] contract, which allows builders
// to query and mutate the classification in a type-safe manner.
package condition
