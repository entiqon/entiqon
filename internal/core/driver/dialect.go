package driver

// Dialect defines the behavior required to generate SQL syntax
// that conforms to a specific database engine (e.g., Postgres, MySQL).
type Dialect interface {
	// Name returns the name of the dialect.md (e.g., "postgres", "mysql").
	Name() string

	// Quote wraps an SQL identifier (e.g., column or table name)
	// using the appropriate quotation syntax for the target database.
	Quote(identifier string) string

	// Escape formats a value for safe inclusion in SQL strings.
	// This is intended for debugging and diagnostics only — NOT for query injection.
	Escape(value any) string

	// SupportsUpsert indicates whether the target dialect.md supports
	// native upsert syntax (e.g., INSERT ... ON CONFLICT).
	SupportsUpsert() bool

	// BuildLimitOffset constructs the SQL clause for pagination,
	// using LIMIT and OFFSET keywords as appropriate for the dialect.md.
	BuildLimitOffset(limit, offset int) string
}
