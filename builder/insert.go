package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// InsertBuilder builds a SQL INSERT statement.
//
// It supports inserting data into a table with specified columns and values,
// and can optionally append a RETURNING clause (PostgreSQL).
type InsertBuilder struct {
	dialect   driver.Dialect
	table     string
	columns   []token.FieldToken
	values    [][]any
	returning []token.FieldToken
}

// NewInsert returns a new instance of InsertBuilder.
func NewInsert() *InsertBuilder {
	return &InsertBuilder{
		columns:   []token.FieldToken{},
		values:    [][]any{},
		returning: []token.FieldToken{},
	}
}

// Into sets the target table for the INSERT operation.
func (ib *InsertBuilder) Into(table string) *InsertBuilder {
	ib.table = table
	return ib
}

// Columns sets the column definitions using FieldFrom(...) and resets existing ones.
func (ib *InsertBuilder) Columns(names ...string) *InsertBuilder {
	ib.columns = []token.FieldToken{}
	for _, name := range names {
		ib.columns = append(ib.columns, token.Field(name))
	}
	return ib
}

// Values appends a row of values using a map of column name to value.
// Each call adds a new row. The map must contain every column defined in Columns().
func (ib *InsertBuilder) Values(row ...any) *InsertBuilder {
	ib.values = append(ib.values, row)
	return ib
}

// Returning adds one or more column names to the RETURNING clause.
// It parses string expressions into FieldTokens.
// If called multiple times, it appends to the existing list.
func (ib *InsertBuilder) Returning(fields ...string) *InsertBuilder {
	for _, f := range fields {
		ib.returning = append(ib.returning, token.FieldsFromExpr(f)...)
	}
	return ib
}

// UseDialect resolves and applies the SQL dialect by name (e.g., "postgres").
// This method configures how identifiers (tables, columns) are quoted
// and how engine-specific syntax is emitted.
func (ib *InsertBuilder) UseDialect(name string) *InsertBuilder {
	ib.dialect = driver.ResolveDialect(name)
	return ib
}

// WithDialect sets the SQL dialect engine used for quoting identifiers.
// It may be removed in a future version in favor of the string-based resolver.
//
// Deprecated: Use UseDialect(name string) instead for consistent resolution and future-proofing.
func (ib *InsertBuilder) WithDialect(name string) *InsertBuilder {
	ib.dialect = driver.ResolveDialect(name)
	return ib
}

// Build compiles the full INSERT SQL statement along with arguments.
// Returns an error if the structure is invalid or values are missing.
// Build compiles the full INSERT SQL statement along with arguments.
// Returns an error if the structure is invalid or values are missing.
func (ib *InsertBuilder) Build() (string, []any, error) {
	if ib.table == "" {
		return "", nil, fmt.Errorf("INSERT requires a target table")
	}
	if len(ib.columns) == 0 {
		return "", nil, fmt.Errorf("INSERT requires at least one field")
	}
	if len(ib.values) == 0 {
		return "", nil, fmt.Errorf("INSERT requires at least one row of values")
	}

	// ─────────────────────────────────────
	// Quote table name
	// ─────────────────────────────────────
	table := ib.table
	if ib.dialect != nil {
		table = ib.dialect.Quote(ib.table)
	}

	// ─────────────────────────────────────
	// Quote field names
	// ─────────────────────────────────────
	var quotedCols []string
	for _, column := range ib.columns {
		if column.Alias != "" {
			return "", nil, fmt.Errorf("column aliasing is not supported in INSERT statements")
		}

		name := column.Name
		if ib.dialect != nil && !column.IsRaw {
			name = ib.dialect.Quote(name)
		}
		quotedCols = append(quotedCols, name)
	}
	columns := strings.Join(quotedCols, ", ")

	// ─────────────────────────────────────
	// Construct value placeholders and args
	// ─────────────────────────────────────
	var valuesSQL []string
	var args []any

	for i, row := range ib.values {
		if len(row) != len(ib.columns) {
			return "", nil, fmt.Errorf("row %d has %d values, expected %d", i+1, len(row), len(ib.columns))
		}

		placeholders := make([]string, len(row))
		for i, val := range row {
			placeholders[i] = "?"
			args = append(args, val)
		}
		valuesSQL = append(valuesSQL, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
	}

	// ─────────────────────────────────────
	// Compose base INSERT statement
	// ─────────────────────────────────────
	tokens := []string{
		fmt.Sprintf("INSERT INTO %s (%s)", table, columns),
		"VALUES",
		strings.Join(valuesSQL, ", "),
	}

	// ─────────────────────────────────────
	// Append RETURNING clause if allowed
	// ─────────────────────────────────────
	if len(ib.returning) > 0 {
		if ib.dialect == nil || !ib.dialect.SupportsReturning() {
			name := "generic"
			if ib.dialect != nil {
				name = ib.dialect.Name()
			}
			return "", nil, fmt.Errorf("RETURNING is not supported by the active dialect: %q", name)
		}

		var quotedRet []string
		for _, field := range ib.returning {
			name := field.Name
			if !field.IsRaw && ib.dialect != nil {
				name = ib.dialect.Quote(name)
			}
			quotedRet = append(quotedRet, name)
		}
		tokens = append(tokens, "RETURNING "+strings.Join(quotedRet, ", "))
	}

	return strings.Join(tokens, " "), args, nil
}

func (ib *InsertBuilder) BuildInsertOnly() (string, []any, error) {
	if ib.table == "" {
		return "", nil, fmt.Errorf("table name is required")
	}
	if len(ib.columns) == 0 {
		return "", nil, fmt.Errorf("at least one column is required")
	}
	if len(ib.values) == 0 {
		return "", nil, fmt.Errorf("at least one set of values is required")
	}
	for i, row := range ib.values {
		if len(row) != len(ib.columns) {
			return "", nil, fmt.Errorf("row %d: expected %d values, got %d", i+1, len(ib.columns), len(row))
		}
	}

	tableName := ib.table
	if ib.dialect != nil {
		tableName = ib.dialect.Quote(tableName)
	}

	quotedCols := make([]string, len(ib.columns))
	for i, col := range ib.columns {
		quotedCols[i] = col.Name
		if ib.dialect != nil {
			quotedCols[i] = ib.dialect.Quote(col.Name)
		}
	}

	placeholders := make([]string, len(ib.values))
	args := make([]any, 0, len(ib.values)*len(ib.columns))
	for i, row := range ib.values {
		args = append(args, row...)
		ph := make([]string, len(row))
		for j := range row {
			ph[j] = "?"
		}
		placeholders[i] = "(" + strings.Join(ph, ", ") + ")"
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(quotedCols, ", "),
		strings.Join(placeholders, ", "))

	return sql, args, nil
}
