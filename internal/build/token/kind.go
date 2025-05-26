package token

// Kind represents the category or role of a token
// within the SQL generation or query-building pipeline.
//
// This type is used to label tokens like Column, Table, Join, etc.
// for logging, diagnostics, and builder-specific rendering.
//
// The Kind is passed to BaseToken.String(...) for structured output.
//
// Since: v1.6.0
type Kind string

const (
	// KindColumn represents a SQL column token.
	KindColumn Kind = "Column"

	// KindTable represents a SQL table or table alias.
	KindTable Kind = "Table"

	// KindJoin represents a join clause or joinable entity.
	KindJoin Kind = "Join"

	// KindExpression represents a general SQL expression or subquery.
	KindExpression Kind = "Expression"
)
