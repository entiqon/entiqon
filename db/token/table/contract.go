package table

import "github.com/entiqon/entiqon/db/contract"

// Token defines the contract implemented by all table tokens.
// A table token represents a SQL source (table or subquery)
// with optional alias, and supports rendering, cloning,
// validation, and debugging.
//
// Responsibilities:
//   - Enforce strict construction and validation of table sources.
//   - Provide canonical SQL output (Render) and raw fragments (Raw).
//   - Surface validation errors explicitly (Errorable).
//   - Expose identity and structural details (BaseToken, Name).
//
// This contract is consumed by higher-level builders (e.g. SelectBuilder)
// to safely assemble SQL FROM and JOIN clauses.
//
// Implementations:
//   - *table (default constructor: table.New(...))
type Token interface {
	// BaseToken provides identity information:
	//   - Input()   → original input string
	//   - Expr()    → resolved core expression (e.g. table name)
	//   - Alias()   → alias if present
	//   - ExpressionKind()    → classified expression kind
	//   - IsAliased() → true if alias is present
	//   - IsValid()   → true if structurally valid
	contract.BaseToken

	// Clonable provides safe duplication of a table token,
	// preserving its state without sharing underlying pointers.
	contract.Clonable[Token]

	// Debuggable produces verbose developer-oriented diagnostics
	// including raw/alias/errored flags.
	contract.Debuggable

	// Errorable surfaces explicit error states and never panics.
	// SetError(err) preserves input for auditability.
	contract.Errorable[Token]

	// Rawable returns a dialect-agnostic raw fragment
	// suitable for embedding directly in SQL.
	contract.Rawable

	// Renderable produces canonical SQL output (expr + alias).
	contract.Renderable

	// Stringable produces a concise log/debug string
	// prefixed with validity markers (✅/❌).
	contract.Stringable

	// Validable exposes structural validation via IsValid().
	contract.Validable

	// Name returns the normalized table identifier
	// (base name without alias). For subqueries, this is
	// the unaliased expression string.
	Name() string
}

// Ensure *table implements the Token contract at compile time.
var _ Token = (*table)(nil)
