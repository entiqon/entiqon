// Package helpers provides classification and resolution utilities for
// SQL token parsing. It defines common helpers to validate identifiers,
// wildcards, and to classify or resolve expressions into their kind and alias.
//
// # Current Rules
//
// The classifier ResolveExpressionType inspects an input string and
// returns its high-level type:
//
//   - Identifier → plain tokens (e.g. "field")
//   - Subquery   → inputs wrapped in parentheses starting with SELECT
//   - Computed   → parenthesized expressions (e.g. "(a+b)")
//   - Aggregate  → SUM(...), COUNT(...), MIN(...), MAX(...), AVG(...)
//   - Function   → any other FUNC(...) form
//   - Literal    → numeric or quoted strings (e.g. "123", "'abc'")
//   - Invalid    → empty or malformed inputs
//
// # Expression Resolution
//
// Beyond classification, the ResolveExpression function provides
// a higher-level resolver that splits an input string into its
// core expression and optional alias.
//
// ResolveExpression is classifier-driven:
//
//   - Identifiers → "id", "id alias", "id AS alias"
//   - Subqueries → "(SELECT ...)", "(SELECT ...) alias", "(SELECT ...) AS alias"
//   - Computed   → "(a+b)", "(a+b) alias", "(a+b) AS alias"
//   - Aggregates → "COUNT(...)", "SUM(...)", with alias variants
//   - Functions  → "FUNC(...)", with alias variants
//   - Literals   → "'text'", "123", with alias variants
//
// Each branch enforces strict rules:
//   - Subqueries and Computed expressions must be parenthesized.
//   - Aggregates and Functions must include parentheses.
//   - Aliases are allowed only if explicitly permitted via allowAlias.
//   - Aliases may use either the space form ("expr alias") or
//     the explicit AS form ("expr AS alias").
//   - Invalid alias formats or reserved keywords are rejected.
//
// This ensures consistency: classification determines the kind,
// resolution enforces alias rules, and all branches are covered
// explicitly without fallthrough defaults.
//
// # Future Dialect-Specific Rules
//
// Dialects may extend classification or resolution with additional
// rules for functions, operators, or literals. For example, PostgreSQL
// introduces JSON operators and type casts, while MySQL introduces
// special string functions. These extensions should be handled by
// layering dialect-specific checks on top of the generic helpers.
package helpers
