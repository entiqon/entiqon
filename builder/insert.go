package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/dialect"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// InsertBuilder builds a SQL INSERT statement.
//
// It supports inserting data into a table with specified columns and values,
// and can optionally append a RETURNING clause (PostgreSQL).
type InsertBuilder struct {
	dialect   dialect.Engine
	table     string
	columns   []token.FieldToken
	values    [][]any
	returning []string
}

// NewInsert returns a new instance of InsertBuilder.
func NewInsert() *InsertBuilder {
	return &InsertBuilder{
		columns:   []token.FieldToken{},
		values:    [][]any{},
		returning: []string{},
	}
}

// Into sets the target table for the INSERT operation.
func (b *InsertBuilder) Into(table string) *InsertBuilder {
	b.table = table
	return b
}

// Columns sets the column definitions using FieldFrom(...) and resets existing ones.
func (b *InsertBuilder) Columns(names ...string) *InsertBuilder {
	b.columns = []token.FieldToken{}
	for _, name := range names {
		b.columns = append(b.columns, token.Field(name))
	}
	return b
}

// Values appends a row of values using a map of column name to value.
// Each call adds a new row. The map must contain every column defined in Columns().
func (b *InsertBuilder) Values(row ...any) *InsertBuilder {
	b.values = append(b.values, row)
	return b
}

// Returning adds columns for RETURNING clause support (e.g., "id", "created_at").
func (b *InsertBuilder) Returning(fields ...string) *InsertBuilder {
	b.returning = append(b.returning, fields...)
	return b
}

// Build compiles the full INSERT SQL statement along with arguments.
// Returns an error if the structure is invalid or values are missing.
func (b *InsertBuilder) Build() (string, []any, error) {
	if b.table == "" {
		return "", nil, fmt.Errorf("INSERT requires a target table")
	}
	if len(b.columns) == 0 {
		return "", nil, fmt.Errorf("INSERT must define at least one column")
	}
	if len(b.values) == 0 {
		return "", nil, fmt.Errorf("INSERT must contain at least one row of values")
	}

	colNames := make([]string, len(b.columns))
	for i, col := range b.columns {
		name := col.Name
		if b.dialect != nil && !col.IsRaw {
			name = b.dialect.EscapeIdentifier(name)
		}
		if col.Alias != "" {
			name = fmt.Sprintf("%s AS %s", name, col.Alias)
		}
		colNames[i] = name
	}

	var placeholders []string
	var args []any

	for i, row := range b.values {
		if len(row) != len(b.columns) {
			return "", nil, fmt.Errorf("row %d has %d values; expected %d", i+1, len(row), len(b.columns))
		}

		marks := make([]string, len(row))
		for j := range row {
			marks[j] = "?"
		}

		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(marks, ", ")))
		args = append(args, row...)
	}

	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		b.table,
		strings.Join(colNames, ", "),
		strings.Join(placeholders, ", "),
	)

	if len(b.returning) > 0 {
		stmt += " RETURNING " + strings.Join(b.returning, ", ")
	}

	return stmt, args, nil
}

// WithDialect sets the SQL dialect engine (e.g., PostgresEngine) for identifier escaping.
func (b *InsertBuilder) WithDialect(e dialect.Engine) *InsertBuilder {
	b.dialect = e
	return b
}
