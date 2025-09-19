package dialect

// Options defines the shared capabilities and behaviors
// across all SQL dialect implementations.
type Options struct {
	// Name of the dialect ("generic", "postgres", "mysql", "sqlite", etc.)
	Name string

	// QuoteStyle defines how identifiers are quoted.
	// Examples: `"` (ANSI/SQLDialect, Postgres), "`" (MySQL), "[" (SQL Server).
	QuoteStyle string

	// PlaceholderStyle defines how parameters are rendered.
	// Examples:
	//   "?"    (SQLDialect, MySQL, SQLite),
	//   "$%d"  (Postgres),
	//   "@p%d" (SQL Server).
	PlaceholderStyle string

	// AllowMerge indicates support for MERGE statements.
	AllowMerge bool

	// AllowUpsert indicates support for INSERT ... ON CONFLICT / DUPLICATE KEY UPDATE.
	AllowUpsert bool

	// ForcedAliasing defines aliasing rules.
	//   - true  → requires explicit "AS" for all aliases.
	//   - false → allows implicit aliases (e.g., "table t").
	ForcedAliasing bool

	// EnableReturning indicates whether INSERT/UPDATE/DELETE
	// can use RETURNING clauses.
	EnableReturning bool

	// SupportsCTE indicates whether WITH (Common Table Expressions) are supported.
	SupportsCTE bool

	// SupportsWindowFunctions indicates support for OVER() / windowed expressions.
	SupportsWindowFunctions bool

	// MaxPlaceholderIndex defines the maximum index supported for placeholders
	// (useful for dialects like Oracle where :1..:N may be limited).
	// Zero or negative → no limit.
	MaxPlaceholderIndex int
}
