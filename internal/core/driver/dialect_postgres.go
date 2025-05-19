package driver

import "fmt"

// PostgresDialect defines SQL behavior specific to PostgreSQL.
type PostgresDialect struct {
	BaseDialect
}

// NewPostgresDialect returns a new PostgresDialect instance.
func NewPostgresDialect() *PostgresDialect {
	return &PostgresDialect{
		BaseDialect: BaseDialect{DialectName: "postgres"},
	}
}

// SupportsUpsert returns true for PostgreSQL, which supports
// INSERT ... ON CONFLICT for native upsert operations.
func (d *PostgresDialect) SupportsUpsert() bool {
	return true
}

// SupportsReturning returns true because PostgreSQL supports
// RETURNING clauses on INSERT, UPDATE, and DELETE statements.
func (d *PostgresDialect) SupportsReturning() bool {
	return true
}

// Placeholder returns PostgreSQL-style placeholders using 1-based indexing.
// Example: "$1", "$2", etc.
func (d *PostgresDialect) Placeholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

// QuoteIdentifier returns the given identifier quoted for PostgreSQL.
func (d *PostgresDialect) QuoteIdentifier(identifier string) string {
	return `"` + identifier + `"`
}
