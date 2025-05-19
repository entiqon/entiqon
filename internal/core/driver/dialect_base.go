package driver

import (
	"fmt"
	"strconv"
)

// BaseDialect provides a foundational implementation of the Dialect interface.
// It can be embedded and selectively overridden by specific dialect structs.
type BaseDialect struct {
	// DialectName holds the unique identifier of the dialect (e.g., "postgres").
	DialectName string
}

// Name returns the identifier of the dialect (e.g., "postgres", "mysql").
func (b *BaseDialect) Name() string {
	return b.DialectName
}

// QuoteLiteral returns the given value wrapped appropriately for SQL output.
//
// ⚠️ This method is NOT SQL-injection safe and MUST NOT be used in production query construction.
// It is intended only for debugging, logging, or test output purposes.
func (b *BaseDialect) QuoteLiteral(value any) string {
	switch v := value.(type) {
	case string:
		return "'" + v + "'"
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// QuoteIdentifier returns the given identifier quoted for PostgreSQL.
//
// This default behavior matches ANSI SQL and is compatible with Postgres.
func (b *BaseDialect) QuoteIdentifier(identifier string) string {
	return `"` + identifier + `"`
}

// SupportsUpsert returns false by default.
// Override in dialects that support upsert syntax (e.g., Postgres).
func (b *BaseDialect) SupportsUpsert() bool {
	return false
}

// SupportsReturning returns false by default.
// Dialects that support RETURNING clauses (like PostgreSQL)
// must override this method to return true.
func (b *BaseDialect) SupportsReturning() bool {
	return false
}

// BuildLimitOffset returns a LIMIT/OFFSET clause based on the provided values.
// If both limit and offset are non-negative, both are included.
// If only one is set, the appropriate clause is emitted.
// If neither is valid, an empty string is returned.
func (b *BaseDialect) BuildLimitOffset(limit, offset int) string {
	switch {
	case limit >= 0 && offset >= 0:
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	case limit >= 0:
		return fmt.Sprintf("LIMIT %d", limit)
	case offset >= 0:
		return fmt.Sprintf("OFFSET %d", offset)
	default:
		return ""
	}
}
