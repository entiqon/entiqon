// Package helpers provides classification and resolution utilities for
// SQL token parsing. It defines common helpers to validate identifiers,
// wildcards, classify/resolve expressions, and parse condition operators.
//
// # Expression Classification
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
// Beyond classification, ResolveExpression splits an input string into
// its core expression and optional alias. Aliases may use either:
//
//   - Trailing form → "expr alias"
//   - AS form       → "expr AS alias"
//
// Alias validation enforces reserved keywords and identifier rules,
// and can be toggled via the allowAlias flag.
//
// # Condition Resolution
//
// The ResolveCondition function parses SQL-like condition inputs into:
//
//   - field → the identifier being compared (e.g. "id")
//   - op    → the detected operator (e.g. =, IN, BETWEEN, IS NULL)
//   - value → the right-hand side, normalized as:
//   - scalar (any Go type)
//   - slice for IN/NOT IN and BETWEEN
//   - nil for IS NULL / IS NOT NULL
//
// Examples:
//
//	"id = 1"                   → field="id", op="=", value=1
//	"price BETWEEN 1 AND 10"   → field="price", op="BETWEEN", value=[1,10]
//	"lastname IN ('a','b')"    → field="lastname", op="IN", value=["a","b"]
//	"deleted_at IS NULL"       → field="deleted_at", op="IS NULL", value=nil
//
// If the operator is missing, a bare identifier defaults to "=".
// For example, "id" is resolved as field="id", op="=", value=nil.
//
// # Condition Validation
//
// The IsValidSlice helper ensures operator/value consistency:
//
//   - IN / NOT IN → require a non-empty slice
//   - BETWEEN     → requires exactly 2 values
//
// Invalid operator/value pairs are rejected during condition construction.
//
// # Supporting Utilities
//
// Additional helpers provide low-level parsing and normalization:
//
//   - parseBetween → splits "x AND y" into [x,y]
//   - parseList    → parses CSV or parenthesized lists into []any
//   - coerceScalar → converts string tokens into int, float64, nil, or string
//   - ToParamKey   → converts identifiers like "users.id" into safe keys ("users_id")
//   - splitCSVRespectingQuotes → splits lists while preserving quoted commas
//
// # Future Dialect-Specific Rules
//
// Dialects may extend classification, resolution, or condition parsing
// with additional rules for functions, operators, or literals.
// For example, PostgreSQL introduces JSON operators and type casts,
// while MySQL introduces special string functions.
package helpers
