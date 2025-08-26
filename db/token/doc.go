// Package token provides low-level primitives to represent SQL query
// fragments in a dialect-agnostic way.
//
// # Overview
//
// The token package defines building blocks such as Field and Table
// that are consumed by higher-level builders (e.g. SelectBuilder) to
// assemble safe, expressive, and auditable SQL statements.
//
// Key principles enforced across all tokens:
//   - Immutability: tokens are never mutated after construction;
//     cloning is explicit.
//   - Auditability: identity aspects (input, expression, alias, owner,
//     validation) are separated into contracts.
//   - Consistency: all tokens share common contracts like BaseToken,
//     Renderable, and Errorable.
//
// Subpackages
//
//   - field: represents a column or expression in a SELECT clause,
//     with support for aliasing, raw expressions, validation, and
//     diagnostics.
//
//   - table: represents a SQL source (table or view) used in FROM /
//     JOIN clauses, with aliasing and validation support.
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
