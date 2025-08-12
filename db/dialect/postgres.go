// File: db/dialect/postgres.go

package dialect

import (
	"fmt"
	"strings"
)

// PostgresDialect implements Dialect interface for PostgreSQL.
type PostgresDialect struct {
	BaseDialect
}

// Name returns the name of the dialect.
func (d *PostgresDialect) Name() string {
	return "postgres"
}

// QuoteIdentifier quotes an identifier with double quotes,
// and escapes embedded double quotes by doubling them.
func (d *PostgresDialect) QuoteIdentifier(name string) string {
	escaped := strings.ReplaceAll(name, `"`, `""`)
	return `"` + escaped + `"`
}

// Placeholder returns the PostgreSQL-style positional parameter placeholder, e.g., $1, $2, ...
func (d *PostgresDialect) Placeholder(n int) string {
	return fmt.Sprintf("$%d", n)
}

// SupportsReturning returns true because PostgreSQL supports RETURNING clause.
func (d *PostgresDialect) SupportsReturning() bool {
	return true
}

// PaginationSyntax returns the PostgreSQL LIMIT/OFFSET syntax.
func (d *PostgresDialect) PaginationSyntax(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	} else if offset > 0 {
		return fmt.Sprintf("OFFSET %d", offset)
	}
	return ""
}
