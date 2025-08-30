// Package identifier provides the classification of SQL expressions
// into broad syntactic categories such as subqueries, functions,
// aggregates, literals, and plain identifiers.
//
// # Overview
//
// The identifier.Type enum is a low-level building block used by
// SQL tokens (Field, Table, …) and higher-level builders
// (SelectBuilder, …). It allows consistent parsing and validation
// of input expressions without introducing cyclic dependencies.
//
// Classification is purely syntactic, not semantic. For example,
// SUM(qty) is classified as an Aggregate even if it appears in an
// invalid position in the query.
//
// # Categories
//
//   - Invalid:    could not classify
//   - Subquery:   "(SELECT ...)"
//   - Computed:   other parenthesized expressions, e.g. "(a+b)"
//   - Aggregate:  SUM, COUNT, MAX, MIN, AVG
//   - Function:   other calls with parentheses, e.g. JSON_EXTRACT(data)
//   - Literal:    quoted string or numeric constant
//   - Identifier: plain table or column name (default fallback)
//
// # Philosophy
//
//   - Never panic: always return a Type, with Invalid or Unknown
//     as safe fallbacks.
//   - Auditability: preserve the original classification for
//     debugging and logs.
//   - Strict enforcement: higher-level resolvers must reject inputs
//     that do not classify correctly.
//
// Example usage is provided in example_test.go and the package README.
package identifier
