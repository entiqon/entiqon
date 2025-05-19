package driver

// Dialect represents SQL dialect-specific behaviors for quoting,
// escaping, pagination, and advanced SQL features like UPSERT and RETURNING.
type Dialect interface {
	// Name returns the name of the dialect (e.g., "postgres", "mysql").
	Name() string

	// Quote wraps a column or table name using dialect-specific syntax.
	// Example: postgres uses double quotes → "users"
	Quote(identifier string) string

	// Escape returns a safely escaped string version of a value,
	// for debugging/logging purposes only. Not used in actual queries.
	Escape(value any) string

	// SupportsUpsert indicates whether the dialect supports native
	// upsert operations such as "INSERT ... ON CONFLICT".
	SupportsUpsert() bool

	// SupportsReturning indicates whether the dialect supports RETURNING
	// clauses, e.g., "INSERT ... RETURNING id". Only PostgreSQL and similar
	// engines support this. Used in InsertBuilder and future builders.
	SupportsReturning() bool

	// BuildLimitOffset generates the correct LIMIT/OFFSET clause for pagination.
	// This varies across dialects.
	BuildLimitOffset(limit, offset int) string
}
