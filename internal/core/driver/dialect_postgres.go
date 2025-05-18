package driver

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
