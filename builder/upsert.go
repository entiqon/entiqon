package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/driver"
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
	BaseBuilder
	// insert handles the INSERT INTO portion of the query.
	insert *InsertBuilder
	// conflictColumns are the columns used to detect conflicts.
	conflictColumns []string
	// updateSet holds the assignments applied when a conflict occurs.
	updateSet []Assignment
	returning []string
}

// NewUpsert returns a new instance of UpsertBuilder.
func NewUpsert() *UpsertBuilder {
	dialect := driver.NewGenericDialect()
	return &UpsertBuilder{
		BaseBuilder: BaseBuilder{dialect: dialect},
		insert:      NewInsert().UseDialect(dialect.Name()),
		updateSet:   make([]Assignment, 0),
		returning:   make([]string, 0),
	}
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
	fmt.Println(!b.dialect.SupportsReturning() && len(b.returning) > 0)
	if !b.dialect.SupportsReturning() && len(b.returning) > 0 {
		b.AddStageError("RETURNING", fmt.Errorf("UPSERT: RETURNING is not supported for dialect: %s", b.dialect.Name()))
	} else {
		b.returning = append(b.returning, columns...)
	}
	return b
}

// DoUpdateSet defines how to update columns if a conflict is found.
func (b *UpsertBuilder) DoUpdateSet(assignments ...Assignment) *UpsertBuilder {
	b.updateSet = append(b.updateSet, assignments...)
	return b
}

// UseDialect resolves and applies the dialect_engine.md by name (e.g., "postgres").
// It replaces any previously set dialect on the builder.
func (b *UpsertBuilder) UseDialect(name string) *UpsertBuilder {
	b.insert.dialect = driver.ResolveDialect(name)
	b.dialect = b.insert.dialect
	return b
}

// Build compiles the UPSERT SQL statement and returns the query and arguments.
func (b *UpsertBuilder) Build() (string, []any, error) {
	dialect := b.GetDialect()
	if b.insert == nil {
		b.insert = NewInsert().UseDialect(dialect.Name())
	}

	if b.HasErrors() {
		return "", nil, fmt.Errorf("UPSERT: %d invalid  condition(s)", len(b.GetErrors()))
	}

	insertSQL, args, err := b.insert.BuildInsertOnly()
	if err != nil {
		return "", nil, fmt.Errorf("UPSERT: %w", err)
	}

	tokens := []string{insertSQL}

	// ───────────────────────────────────────────────
	// ON CONFLICT (columns)
	// ───────────────────────────────────────────────
	if len(b.conflictColumns) > 0 {
		var quoted []string
		for _, col := range b.conflictColumns {
			if col == "" {
				return "", nil, fmt.Errorf("UPSERT: empty conflict column name")
			}
			quoted = append(quoted, dialect.QuoteIdentifier(col))
		}
		tokens = append(tokens, fmt.Sprintf("ON CONFLICT (%s)", strings.Join(quoted, ", ")))
	}

	// ───────────────────────────────────────────────
	// DO UPDATE SET or DO NOTHING
	// ───────────────────────────────────────────────
	if len(b.updateSet) == 0 {
		tokens = append(tokens, "DO NOTHING")
	} else {
		var assignments []string
		for _, assign := range b.updateSet {
			if assign.Column == "" || assign.Expr == "" {
				return "", nil, fmt.Errorf("UPSERT: column or expression is empty")
			}
			col := b.GetDialect().QuoteIdentifier(assign.Column)
			assignments = append(assignments, fmt.Sprintf("%s = %s", col, assign.Expr))
		}
		if len(assignments) > 0 {
			tokens = append(tokens, "DO UPDATE SET", strings.Join(assignments, ", "))
		}
	}

	// ───────────────────────────────────────────────
	// RETURNING (dialect-aware)
	// ───────────────────────────────────────────────
	if len(b.returning) > 0 {
		if dialect.SupportsReturning() {
			var returnCols []string
			for _, col := range b.returning {
				returnCols = append(returnCols, dialect.QuoteIdentifier(col))
			}
			tokens = append(tokens, "RETURNING", strings.Join(returnCols, ", "))
		} else {
			return "", nil, fmt.Errorf("UPSERT: RETURNING not supported in dialect: %s", dialect.Name())
		}
	}

	return strings.Join(tokens, " "), args, nil
}
