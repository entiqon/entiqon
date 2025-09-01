// Package types defines low-level, dependency-free classifications used by SQL tokens.
//
// These enums and registries provide the common building blocks for higher-level
// tokens (Field, Table, Join, Condition, …) and builders (SelectBuilder, …),
// ensuring consistent semantics and avoiding cyclic dependencies.
//
// ## Purpose
//
// Types capture "what kind of thing is this?" without rendering or ownership.
// For example, whether an expression is an identifier, a function, or a literal.
// Tokens (Field, Table, Condition, …) then embed these classifications.
//
// ## Subpackages
//
//   - identifier: classification of SQL expressions (Identifier, Subquery, Literal,
//     Aggregate, Function, Computed)
//   - join:       type-safe join kinds (Inner, Left, Right, Full, Cross, Natural)
//   - condition:  condition clause kinds (Single, And, Or)
//   - operator:   comparison and predicate operators (=, !=, IN, BETWEEN, IS NULL, …)
//
// ## Notes
//
// * All types are designed to be dependency-free: no reliance on token internals.
// * Each type provides:
//   - Strongly typed enum or registry
//   - `String()` canonical form
//   - `ParseFrom(any)` coercion from string/int/etc.
//   - `IsValid()` validation
//   - Aliases for concise references (e.g. `gte`, `nin`, `isnull`)
//
// * Operator.Type is registry-backed with deterministic ordering and O(1) parsing.
// * Tests in each subpackage reach 100% coverage and include runnable examples.
//
// By separating classification logic into this package, Entiqon ensures that
// tokens remain focused on ownership, rendering, and composition, while all
// semantic categories remain centralized and reusable.
package types
