package builder

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/dialect"
)

// SelectBuilder builds simple SELECT queries.
type SelectBuilder struct {
	dialect dialect.Dialect
	columns []string
	table   string
	limit   int
	offset  int
}

// NewSelect creates a new SelectBuilder with the provided dialect.
// If nil is passed, BaseDialect is used by default.
func NewSelect(d dialect.Dialect) *SelectBuilder {
	if d == nil {
		d = &dialect.BaseDialect{}
	}
	return &SelectBuilder{
		dialect: d,
		columns: []string{},
	}
}

// Columns adds columns to the select clause.
func (b *SelectBuilder) Columns(cols ...string) *SelectBuilder {
	b.columns = append(b.columns, cols...)
	return b
}

// Source sets the table name.
func (b *SelectBuilder) Source(table string) *SelectBuilder {
	b.table = table
	return b
}

// Limit sets the limit clause.
func (b *SelectBuilder) Limit(limit int) *SelectBuilder {
	b.limit = limit
	return b
}

// Offset sets the offset clause.
func (b *SelectBuilder) Offset(offset int) *SelectBuilder {
	b.offset = offset
	return b
}

// Build builds the SQL query string.
func (b *SelectBuilder) Build() (string, error) {
	if len(b.columns) == 0 {
		b.columns = append(b.columns, "*")
	}
	if b.table == "" {
		return "", fmt.Errorf("no table specified")
	}

	var sb strings.Builder

	sb.WriteString("SELECT ")

	// Quote columns
	var quotedCols []string
	for _, col := range b.columns {
		if col == "*" {
			quotedCols = append(quotedCols, col) // leave * unquoted
		} else {
			quotedCols = append(quotedCols, b.dialect.QuoteIdentifier(col))
		}
	}
	sb.WriteString(strings.Join(quotedCols, ", "))

	sb.WriteString(" FROM ")
	sb.WriteString(b.dialect.QuoteIdentifier(b.table))

	// Add pagination if present
	if b.limit > 0 || b.offset > 0 {
		sb.WriteString(" ")
		sb.WriteString(b.dialect.PaginationSyntax(b.limit, b.offset))
	}

	return sb.String(), nil
}

func (b *SelectBuilder) String() string {
	sql, err := b.Build()
	if err != nil {
		return fmt.Sprintf("Status ❌: Error building SQL: %v", err)
	}

	// Assuming no parameters yet; if you add params later, update accordingly
	params := []interface{}{}

	paramsStr := ""
	if len(params) > 0 {
		paramsStr = fmt.Sprintf("%v", params)
	}

	return fmt.Sprintf("Status ✅: SQL=%s, Params=%s", sql, paramsStr)
}
