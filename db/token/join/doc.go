// Package token provides primitives to represent SQL JOIN clauses
// in a dialect-agnostic way. A Join token binds two tables with
// a specific token kind (INNER, LEFT, RIGHT, FULL) and an ON condition.
//
// # Philosophy
//
// Join construction exposes two layers of API:
//
//   - Safe constructors — NewInner, NewLeft, NewRight, NewFull
//     These are the recommended entrypoints for most code. They provide
//     compile-time safety: the token kind is fixed and guaranteed valid.
//     Using them keeps queries intention-revealing and free from runtime
//     surprises.
//
//   - Flexible constructor — New(kind any, left, right, condition)
//     This variant is for advanced scenarios where the token kind must
//     be decided dynamically (e.g. driven by configuration, DSLs, or
//     user input). It accepts either a token.Kind or a free-form
//     string ("LEFT", "LEFT JOIN", case-insensitive). If the kind is
//     not valid, New returns an errored token immediately.
//
// # Guidance
//
// Prefer the safe constructors in application code. Reserve New(kind,…)
// for frameworks or tooling layers where extensibility is critical.
// This separation ensures day-to-day usage remains safe while still
// allowing power users to integrate joins dynamically when necessary.
//
// # Example
//
//	// Safe usage
//	j1 := token.NewInner("users", "orders", "users.id = orders.user_id")
//	j2 := token.NewLeft("products", "categories", "products.cat_id = categories.id")
//
//	// Flexible usage
//	j3 := token.New("RIGHT", "employees", "departments", "employees.dep_id = departments.id")
//	if !j3.IsValid() {
//		log.Fatalf("invalid token: %v", j3.Error())
//	}
package join
