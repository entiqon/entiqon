package builder

import (
	"fmt"
	"strings"
)

// InsertBuilder builds a SQL INSERT statement.
//
// It supports inserting data into a table with specified columns and values,
// and can optionally append a RETURNING clause (PostgreSQL).
type InsertBuilder struct {
	intoTable string
	columns   []string
	rows      [][]any
	returning []string
}

// NewInsert returns a new instance of InsertBuilder.
func NewInsert() *InsertBuilder {
	return &InsertBuilder{
		columns:   make([]string, 0),
		rows:      make([][]any, 0),
		returning: make([]string, 0),
	}
}

// Into sets the target table for the INSERT operation.
func (b *InsertBuilder) Into(table string) *InsertBuilder {
	b.intoTable = table
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
	b.rows = append(b.rows, values)
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
	if b.intoTable == "" {
		return "", nil, fmt.Errorf("no target table specified")
	}
	if len(b.columns) == 0 {
		return "", nil, fmt.Errorf("no columns specified")
	}
	if len(b.rows) == 0 {
		return "", nil, fmt.Errorf("no values provided")
	}

	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(b.intoTable)
	sb.WriteString(" (")
	sb.WriteString(strings.Join(b.columns, ", "))
	sb.WriteString(") VALUES ")

	placeholderRow := "(" + strings.TrimRight(strings.Repeat("?, ", len(b.columns)), ", ") + ")"
	valuesPlaceholders := make([]string, len(b.rows))
	args := make([]any, 0, len(b.rows)*len(b.columns))

	for i, row := range b.rows {
		valuesPlaceholders[i] = placeholderRow
		args = append(args, row...)
	}

	sb.WriteString(strings.Join(valuesPlaceholders, ", "))

	if len(b.returning) > 0 {
		sb.WriteString(" RETURNING ")
		sb.WriteString(strings.Join(b.returning, ", "))
	}

	return sb.String(), args, nil
}
