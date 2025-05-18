package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/dialect"
)

// Assignment represents a column update assignment like col = expr.
type Assignment struct {
	// Column is the name of the column to be updated.
	Column string
	// Expr is the SQL expression to assign to the column.
	Expr string
}

// UpsertBuilder builds a SQL UPSERT (INSERT ... ON CONFLICT DO UPDATE) statement.
//
// It composes an internal InsertBuilder and extends it with conflict resolution clauses.
type UpsertBuilder struct {
	// insert handles the INSERT INTO portion of the query.
	insert *InsertBuilder
	// conflictColumns are the columns used to detect conflicts.
	conflictColumns []string
	// updateSet holds the assignments applied when a conflict occurs.
	updateSet []Assignment
}

// NewUpsert returns a new instance of UpsertBuilder.
func NewUpsert() *UpsertBuilder {
	return &UpsertBuilder{
		insert: NewInsert(),
	}
}

// WithDialect sets the SQL dialect to use for escaping identifiers.
func (b *UpsertBuilder) WithDialect(d dialect.Engine) *UpsertBuilder {
	b.insert.WithDialect(d)
	return b
}

// Into sets the target table for the UPSERT operation.
func (b *UpsertBuilder) Into(table string) *UpsertBuilder {
	b.insert.Into(table)
	return b
}

// Columns sets the list of columns for the insert portion of the UPSERT.
func (b *UpsertBuilder) Columns(cols ...string) *UpsertBuilder {
	b.insert.Columns(cols...)
	return b
}

// Values adds a new row of values for the insert portion of the UPSERT.
func (b *UpsertBuilder) Values(values ...any) *UpsertBuilder {
	b.insert.Values(values...)
	return b
}

// OnConflict specifies the columns that determine the conflict for the UPSERT.
func (b *UpsertBuilder) OnConflict(columns ...string) *UpsertBuilder {
	b.conflictColumns = append(b.conflictColumns, columns...)
	return b
}

// Returning specifies the RETURNING clause for the UPSERT statement.
func (b *UpsertBuilder) Returning(columns ...string) *UpsertBuilder {
	b.insert.Returning(columns...)
	return b
}

// DoUpdateSet defines how to update columns if a conflict is found.
func (b *UpsertBuilder) DoUpdateSet(assignments ...Assignment) *UpsertBuilder {
	b.updateSet = append(b.updateSet, assignments...)
	return b
}

// Build compiles the UPSERT SQL statement and returns the query and arguments.
func (b *UpsertBuilder) Build() (string, []any, error) {
	sql, args, err := b.insert.BuildInsertOnly()
	if err != nil {
		return "", nil, err
	}

	if len(b.conflictColumns) > 0 {
		sql += " " + b.buildConflictClause()
	}

	if len(b.updateSet) > 0 {
		sql += " " + b.buildUpdateSetClause()
	} else {
		sql += " DO NOTHING"
	}

	if len(b.insert.returning) > 0 {
		sql += " RETURNING " + strings.Join(b.getEscapedIdentifiers(b.insert.returning), ", ")
	}

	return sql, args, nil
}

func (b *UpsertBuilder) buildConflictClause() string {
	escaped := make([]string, len(b.conflictColumns))
	for i, col := range b.conflictColumns {
		if b.insert.dialect != nil {
			escaped[i] = b.insert.dialect.EscapeIdentifier(col)
		} else {
			escaped[i] = col
		}
	}
	return "ON CONFLICT (" + strings.Join(escaped, ", ") + ")"
}

func (b *UpsertBuilder) buildUpdateSetClause() string {
	parts := make([]string, len(b.updateSet))
	for i, assign := range b.updateSet {
		col := assign.Column
		if b.insert.dialect != nil {
			col = b.insert.dialect.EscapeIdentifier(assign.Column)
		}
		parts[i] = fmt.Sprintf("%s = %s", col, assign.Expr)
	}
	return "DO UPDATE SET " + strings.Join(parts, ", ")
}

func (b *UpsertBuilder) getEscapedIdentifiers(cols []string) []string {
	if b.insert.dialect == nil {
		return cols
	}
	escaped := make([]string, len(cols))
	for i, col := range cols {
		escaped[i] = b.insert.dialect.EscapeIdentifier(col)
	}
	return escaped
}
