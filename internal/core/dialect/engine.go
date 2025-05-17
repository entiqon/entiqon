package dialect

// Engine defines behavior for dialect-specific escaping and formatting.
type Engine interface {
	// EscapeParam returns a formatted SQL-safe string for the given value.
	// Used for logging, previewing, or inline SQL when applicable.
	EscapeParam(value any) string

	// EscapeIdentifier returns the dialect-safe name for columns or tables.
	// e.g., "user" → \"user\" (Postgres), `user` (MySQL)
	EscapeIdentifier(name string) string

	// DialectName returns the string identifier for this dialect (e.g., "postgres", "mysql").
	DialectName() string
}
