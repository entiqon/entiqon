// Package condition defines SQL condition tokens used by the Entiqon
// SQL builder. A condition token represents a logical predicate
// (e.g. "id = 1", "age > 18", "status IN ('active', 'pending')")
// that can be composed into WHERE or HAVING clauses.
//
// # Purpose
//
// Condition tokens wrap an expression, its operator, and any bound
// values into a strongly-typed construct that implements shared
// Entiqon contracts such as Renderable, Debuggable, and Validable.
// This allows conditions to be validated, inspected, and rendered
// safely into SQL strings across dialects.
//
// # Types
//
// The primary entry points are:
//
//   - Token interface: contract implemented by all conditions.
//   - New(kind, expr, [value]): generic constructor.
//   - NewAnd / NewOr: convenience constructors for logical composition.
//
// # Examples
//
// Basic usage with inline expressions:
//
//	cond := condition.New(ct.Single, "id = 1")
//	if cond.IsValid() {
//	    fmt.Println(cond.Render()) // "id = :id"
//	    fmt.Println(cond.Value())  // 1
//	}
//
// With operator and value:
//
//	cond := condition.New(ct.Single, "age", operator.GreaterThan, 18)
//	fmt.Println(cond.Render()) // "age > :age"
//	fmt.Println(cond.Value())  // 18
//
// Logical composition:
//
//	cond := condition.NewAnd("status = ?", "active")
//	fmt.Println(cond.Kind()) // And
//
// # Integration
//
// Condition tokens are designed to integrate with higher-level
// builders such as SelectBuilder:
//
//	select := builder.NewSelect("users").
//	    Where(condition.New(ct.Single, "age", operator.GreaterThan, 21))
//	fmt.Println(select.Render())
//
// This package is low-level and focused on representing conditions.
// Composition, parameter binding, and SQL dialect handling are
// delegated to higher layers of the Entiqon stack.
package condition
