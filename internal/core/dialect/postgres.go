package dialect

import (
	"fmt"
	"strings"
)

// Ensure PostgresEngine implements Engine
var _ Engine = (*PostgresEngine)(nil)

// Escape characters used by PostgreSQL
const (
	PostgresIdentifierQuote = `"` // used for identifiers
	PostgresStringQuote     = `'` // used for values
)

// PostgresEngine implements the Engine interface for PostgreSQL.
type PostgresEngine struct{}

// EscapeParam safely formats values for raw SQL preview.
func (e *PostgresEngine) EscapeParam(value any) string {
	switch v := value.(type) {
	case string:
		return PostgresStringQuote + strings.ReplaceAll(v, "'", "''") + PostgresStringQuote
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// EscapeIdentifier escapes table/column names using double quotes.
func (e *PostgresEngine) EscapeIdentifier(name string) string {
	return PostgresIdentifierQuote + strings.ReplaceAll(name, `"`, `""`) + PostgresIdentifierQuote
}

// DialectName returns the name of the dialect.
func (e *PostgresEngine) DialectName() string {
	return "postgres"
}
