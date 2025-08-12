// File: db/dialect/base.go

// Package dialect provides interfaces and implementations
// for SQL dialects to generate database-specific SQL syntax.
//
// BaseDialect offers a generic SQL dialect with common behaviors
// suitable for many relational databases.
package dialect

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Dialect defines the behavior required to generate SQL syntax
// tailored to a specific database dialect.
type Dialect interface {
	// Name returns the name of the dialect, e.g., "postgres", "mysql".
	Name() string

	// QuoteIdentifier quotes an SQL identifier such as table or column names,
	// ensuring it is escaped according to dialect rules.
	QuoteIdentifier(name string) string

	// QuoteLiteral safely quotes a literal value (string, number, bool, time, etc.)
	// for inclusion in SQL statements, escaping where necessary.
	QuoteLiteral(literal any) string

	// PaginationSyntax generates the SQL syntax for limiting
	// and offsetting query results according to dialect rules.
	PaginationSyntax(limit, offset int) string

	// Placeholder returns the parameter placeholder string for
	// the given parameter index (e.g., "?" or "$1").
	Placeholder(index int) string

	// SupportsReturning indicates if the dialect supports
	// a RETURNING clause on INSERT/UPDATE/DELETE statements.
	SupportsReturning() bool
}

// BaseDialect is a generic SQL dialect implementation providing
// default behaviors common to many relational databases.
type BaseDialect struct{}

// Name returns the dialect name "base" indicating generic SQL.
func (d *BaseDialect) Name() string {
	return "base"
}

// QuoteIdentifier returns the given identifier quoted with double quotes,
// suitable for most SQL databases.
func (d *BaseDialect) QuoteIdentifier(name string) string {
	return `"` + name + `"`
}

// QuoteLiteral quotes and escapes a literal value according to its type.
// Strings are escaped to double single-quotes, numeric types are formatted plainly,
// booleans converted to "true"/"false", nil to NULL, and time.Time formatted as timestamp.
func (d *BaseDialect) QuoteLiteral(value any) string {
	switch v := value.(type) {
	case string:
		escaped := strings.ReplaceAll(v, "'", "''")
		return "'" + escaped + "'"
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return "NULL"
	case time.Time:
		return "'" + v.Format("2006-01-02 15:04:05") + "'"
	default:
		escaped := strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''")
		return "'" + escaped + "'"
	}
}

// Placeholder returns "?" as a generic positional parameter placeholder.
func (d *BaseDialect) Placeholder(n int) string {
	return "?"
}

// SupportsReturning returns false by default, indicating the dialect
// does not support RETURNING clauses.
func (d *BaseDialect) SupportsReturning() bool {
	return false
}

// PaginationSyntax returns a SQL LIMIT/OFFSET clause according to the given limit and offset.
// If both are zero or negative, returns an empty string.
func (d *BaseDialect) PaginationSyntax(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	} else if offset > 0 {
		return fmt.Sprintf("OFFSET %d", offset)
	}
	return ""
}
