package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// DeleteBuilder builds a SQL DELETE statement.
//
// It supports WHERE clauses and dialect-aware identifier escaping.
type DeleteBuilder struct {
	// dialect defines the SQL dialect used for identifier escaping
	// (e.g., PostgreSQL, MySQL). It is optional but recommended for safety.
	dialect driver.Dialect

	// table holds the name of the target table from which rows will be deleted.
	table string

	// conditions holds the WHERE clause conditions used to filter which rows
	// will be deleted. If empty, no WHERE clause will be added, which may result
	// in deleting all rows.
	conditions []token.Condition
}

// NewDelete creates and returns a new DeleteBuilder.
func NewDelete() *DeleteBuilder {
	return &DeleteBuilder{
		conditions: []token.Condition{},
	}
}

// From sets the target table to delete from.
func (b *DeleteBuilder) From(table string) *DeleteBuilder {
	b.table = table
	return b
}

// Where starts the WHERE clause by replacing any previous conditions.
func (b *DeleteBuilder) Where(condition string, params ...any) *DeleteBuilder {
	b.conditions = []token.Condition{
		token.NewCondition(token.ConditionSimple, condition, params...),
	}
	return b
}

// AndWhere appends an AND condition to the WHERE clause.
func (b *DeleteBuilder) AndWhere(condition string, params ...any) *DeleteBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionAnd, condition, params...),
	)
	return b
}

// OrWhere appends an OR condition to the WHERE clause.
func (b *DeleteBuilder) OrWhere(condition string, params ...any) *DeleteBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionOr, condition, params...),
	)
	return b
}

// UseDialect resolves and sets the SQL dialect engine using its registered name.
// This replaces any previously set dialect on the builder.
func (b *DeleteBuilder) UseDialect(name string) *DeleteBuilder {
	b.dialect = driver.ResolveDialect(name)
	return b
}

// WithDialect sets the SQL dialect engine directly.
//
// Deprecated: Use UseDialect(name string) instead for consistent resolution and future-proofing.
// This method will be removed in v1.4.0.
func (b *DeleteBuilder) WithDialect(name string) *DeleteBuilder {
	b.dialect = driver.ResolveDialect(name)
	return b
}

// Build compiles the final DELETE SQL string and returns it
// along with the bound arguments to be used with database/sql.
//
// Returns an error if the table name is not set or if the condition format is invalid.
func (b *DeleteBuilder) Build() (string, []any, error) {
	if strings.TrimSpace(b.table) == "" {
		return "", nil, fmt.Errorf("DELETE requires a target table")
	}

	table := b.table
	if b.dialect != nil {
		table = b.dialect.QuoteIdentifier(table)
	}

	sql := fmt.Sprintf("DELETE FROM %s", table)
	var args []any

	if len(b.conditions) > 0 {
		sql += " WHERE "
		condSQL, condArgs, err := token.FormatConditions(b.dialect, b.conditions)
		if err != nil {
			return "", nil, err
		}
		sql += condSQL
		args = append(args, condArgs...)
	}

	return sql, args, nil
}
