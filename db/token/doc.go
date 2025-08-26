// Package token provides low-level primitives to represent SQL query
// fragments in a dialect-agnostic way.
//
// # Overview
//
// The token package defines atomic building blocks such as Field and
// Table. These tokens are consumed by higher-level builders (e.g.
// SelectBuilder) to assemble safe, expressive, and auditable SQL
// statements.
//
// # Doctrine
//
//   - Never panic: constructors always return a non-nil token,
//     even if errored.
//   - Auditability: preserve original input for logs and debugging.
//   - Strict validation: invalid inputs are rejected immediately
//     with explicit errors.
//   - Delegation: parsing rules live in tokens, not in builders.
//   - Clarity: responsibilities are split into explicit contracts.
//
// # Contracts
//
// All tokens implement a shared set of contracts:
//
//   - BaseToken   — identity (input, expression, alias, validity)
//   - Errorable   — explicit error state, never panic
//   - Clonable    — safe duplication with preserved state
//   - Rawable     — SQL-generic rendering (expr, alias, owner)
//   - Renderable  — dialect-agnostic String() output
//   - Stringable  — concise diagnostic/logging string
//   - Ownerable   — ownership binding (HasOwner, Owner, SetOwner)
//
// # Subpackages
//
//   - field: represents a column, identifier, or computed expression
//     with aliasing, validation, and diagnostics.
//
//   - table: represents a SQL source (table or view) used in FROM /
//     JOIN clauses with aliasing and validation support.
//
// # Roadmap
//
// Future tokens will include:
//   - conditions (WHERE / HAVING)
//   - joins (INNER, LEFT, etc.)
//   - functions (aggregates, JSON, custom expressions)
//
// Contracts will progressively enforce stricter auditability across
// all tokens.
package token
