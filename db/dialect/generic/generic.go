package generic

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/entiqon/db/dialect"
)

//
// Generic Dialect
//

// dialectImpl provides the ANSI/Generic implementation of the dialect.SQLDialect
// interface. It is unexported to prevent direct instantiation; consumers should
// always use the New() constructor.
//
// The Generic dialect adheres to ANSI SQL as closely as possible without
// introducing vendor-specific extensions. It serves as a safe default when
// generating SQL queries for databases with unknown or mixed compatibility.
//
// Key behaviors:
//   - Identifiers are quoted using ANSI double quotes (").
//   - Placeholders are always "?" and do not support indexing.
//   - LIMIT and OFFSET are rendered in the standard form.
//   - RETURNING, MERGE, and UPSERT clauses are not supported.
type dialectImpl struct {
	opts dialect.Options
}

// Compile-time check: ensure dialectImpl implements SQLDialect
var _ dialect.SQLDialect = (*dialectImpl)(nil)

//
// Constructor
//

// New returns a new ANSI-compliant generic dialect. The returned value
// implements the dialect.SQLDialect interface.
//
// Example:
//
//	d := generic.New()
//	sql := fmt.Sprintf("SELECT %s FROM %s%s",
//	    d.QuoteIdentifier("id"),
//	    d.QuoteIdentifier("users"),
//	    d.PaginationSyntax(10, 0))
//
// Produces:
//
//	SELECT "id" FROM "users" LIMIT 10
func New() dialect.SQLDialect {
	return &dialectImpl{
		opts: dialect.Options{
			Name:                    "generic",
			QuoteStyle:              `"`,
			PlaceholderStyle:        "?",
			AllowMerge:              false,
			AllowUpsert:             false,
			ForcedAliasing:          false,
			EnableReturning:         false,
			SupportsCTE:             true,
			SupportsWindowFunctions: true,
			MaxPlaceholderIndex:     0,
		},
	}
}

//
// Interface Implementation
//

// Name returns the identifier of this dialect, which is always "generic".
func (d *dialectImpl) Name() string {
	return d.opts.Name
}

// Options returns the feature capability matrix for the Generic dialect.
// Consumers can inspect the returned Options to determine whether a feature
// such as RETURNING or UPSERT is supported before generating SQL that uses it.
func (d *dialectImpl) Options() dialect.Options {
	return d.opts
}

// QuoteIdentifier returns an ANSI-quoted SQL identifier using the configured
// quote style. Identifiers are quoted only when necessary (mixed case, spaces,
// or symbols); lowercase unqualified names are returned as-is.
//
// Example:
//
//	d.QuoteIdentifier("users")    // → users
//	d.QuoteIdentifier("UserData") // → "UserData"
func (d *dialectImpl) QuoteIdentifier(name string) string {
	if name == "" {
		return ""
	}
	if strings.ToLower(name) == name && !strings.ContainsAny(name, " -") {
		return name
	}
	return fmt.Sprintf("%s%s%s", d.opts.QuoteStyle, name, d.opts.QuoteStyle)
}

// QuoteLiteral safely quotes a literal value for inline use in SQL statements.
// This method is provided for scenarios where parameter binding is not used,
// such as debugging or migration scripts. In production, bound parameters are
// strongly preferred for safety and performance.
//
// Supported types:
//   - nil     → "NULL"
//   - string  → escaped and wrapped in single quotes
//   - bool    → "TRUE" or "FALSE"
//   - numbers → rendered in decimal form
//   - time.Time → formatted as UTC 'YYYY-MM-DD HH:MM:SS'
//
// Example:
//
//	d.QuoteLiteral("O'Reilly")   // → 'O''Reilly'
//	d.QuoteLiteral(true)         // → TRUE
//	d.QuoteLiteral(42)           // → 42
//	d.QuoteLiteral(time.Now())   // → '2025-09-19 07:32:00'
func (d *dialectImpl) QuoteLiteral(literal any) string {
	switch v := literal.(type) {
	case nil:
		return "NULL"
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return strconv.FormatFloat(toFloat64(v), 'f', -1, 64)
	case time.Time:
		return fmt.Sprintf("'%s'", v.UTC().Format("2006-01-02 15:04:05"))
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// PaginationSyntax generates an ANSI-compliant LIMIT/OFFSET clause.
// If both limit and offset are non-positive, an empty string is returned.
//
// Example:
//
//	d.PaginationSyntax(10, 0)   // → " LIMIT 10"
//	d.PaginationSyntax(10, 20)  // → " LIMIT 10 OFFSET 20"
//	d.PaginationSyntax(0, 0)    // → ""
func (d *dialectImpl) PaginationSyntax(limit, offset int) string {
	if limit <= 0 && offset <= 0 {
		return ""
	}
	var sb strings.Builder
	if limit > 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	}
	if offset > 0 {
		sb.WriteString(fmt.Sprintf(" OFFSET %d", offset))
	}
	return sb.String()
}

// Placeholder always returns "?" for the Generic dialect. The index parameter
// is ignored since ANSI placeholders are not positional.
//
// Example:
//
//	d.Placeholder(1)  // → "?"
//	d.Placeholder(99) // → "?"
func (d *dialectImpl) Placeholder(_ int) string {
	return d.opts.PlaceholderStyle
}

// toFloat64 normalizes float32/float64 values into float64.
// If the input is not a float, it safely returns 0.
func toFloat64(v any) float64 {
	var value float64
	switch n := v.(type) {
	case float32:
		value = float64(n)
	case float64:
		value = n
	}
	return value
}
