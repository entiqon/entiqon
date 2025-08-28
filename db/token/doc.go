// Package token provides low-level primitives to represent SQL query
// fragments in a dialect-agnostic way.
//
// # Overview
//
// The token package defines atomic building blocks such as Field,
// Table, and Join. These tokens are consumed by higher-level builders
// (e.g. SelectBuilder) to assemble safe, expressive, and auditable
// SQL statements.
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
//   - join: represents JOIN clauses (INNER, LEFT, etc.) with
//     validation of join kind and conditions.
//
// # Supporting modules
//
//   - resolver: centralizes input type validation and expression
//     resolution (expr, alias, kind) with strict rules and
//     subquery detection.
//
//   - expression_kind: classifies raw expressions into categories
//     (Identifier, Computed, Literal, Subquery, Function, Aggregate)
//     and validates identifier aliases, including reserved keyword
//     checks.
//
// # Roadmap
//
// Future tokens will include:
//   - conditions (WHERE / HAVING)
//   - functions (aggregates, JSON, custom expressions)
//
// Contracts will progressively enforce stricter auditability across
// all tokens.
package token
