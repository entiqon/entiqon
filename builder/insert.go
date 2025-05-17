package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/dialect"
)

// InsertBuilder builds a SQL INSERT statement.
//
// It supports inserting data into a table with specified columns and values,
// and can optionally append a RETURNING clause (PostgreSQL).
type InsertBuilder struct {
	dialect   dialect.Engine
	table     string
	columns   []string
	values    [][]any
	returning []string
}

// NewInsert returns a new instance of InsertBuilder.
func NewInsert() *InsertBuilder {
	return &InsertBuilder{
		columns:   make([]string, 0),
		values:    make([][]any, 0),
		returning: make([]string, 0),
	}
}

// WithDialect sets the SQL dialect for identifier escaping.
func (b *InsertBuilder) WithDialect(d dialect.Engine) *InsertBuilder {
	b.dialect = d
	return b
}

// Into sets the target table for the INSERT operation.
func (b *InsertBuilder) Into(table string) *InsertBuilder {
	b.table = table
	return b
}

// Columns defines the column names for the INSERT statement.
func (b *InsertBuilder) Columns(cols ...string) *InsertBuilder {
	b.columns = append(b.columns, cols...)
	return b
}

// Values appends a single row of values to be inserted.
func (b *InsertBuilder) Values(values ...any) *InsertBuilder {
	if len(values) != len(b.columns) {
		panic("values count must match columns count")
	}
	b.values = append(b.values, values)
	return b
}

// Returning specifies columns to return after the insert (e.g., PostgreSQL RETURNING clause).
func (b *InsertBuilder) Returning(columns ...string) *InsertBuilder {
	b.returning = append(b.returning, columns...)
	return b
}

// Build generates the final SQL INSERT statement and returns it
// along with the ordered list of arguments.
func (b *InsertBuilder) Build() (string, []any, error) {
	if b.table == "" {
		return "", nil, fmt.Errorf("INSERT must specify table name")
	}
	if len(b.columns) == 0 {
		return "", nil, fmt.Errorf("INSERT must specify columns")
	}
	if len(b.values) == 0 {
		return "", nil, fmt.Errorf("INSERT must provide at least one row of values")
	}

	var sql strings.Builder
	args := make([]any, 0)

	table := b.table
	if b.dialect != nil {
		table = b.dialect.EscapeIdentifier(table)
	}

	columns := b.columns
	if b.dialect != nil {
		for i, col := range columns {
			columns[i] = b.dialect.EscapeIdentifier(col)
		}
	}

	sql.WriteString("INSERT INTO ")
	sql.WriteString(table)
	sql.WriteString(" (")
	sql.WriteString(strings.Join(columns, ", "))
	sql.WriteString(") VALUES ")

	placeholders := make([]string, len(b.columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	valueBlocks := make([]string, 0, len(b.values))
	for _, row := range b.values {
		if len(row) != len(b.columns) {
			return "", nil, fmt.Errorf("value count does not match column count")
		}
		valueBlocks = append(valueBlocks, "("+strings.Join(placeholders, ", ")+")")
		args = append(args, row...)
	}

	sql.WriteString(strings.Join(valueBlocks, ", "))

	if len(b.returning) > 0 {
		returning := b.returning
		if b.dialect != nil {
			for i, col := range returning {
				returning[i] = b.dialect.EscapeIdentifier(col)
			}
		}
		sql.WriteString(" RETURNING ")
		sql.WriteString(strings.Join(returning, ", "))
	}

	return sql.String(), args, nil
}
