package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/dialect"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// DeleteBuilder builds a SQL DELETE statement.
//
// It supports WHERE clauses and optional RETURNING fields (e.g., for PostgreSQL).
type DeleteBuilder struct {
	dialect dialect.Engine

	// table defines the target table to delete table.
	table string

	// conditions holds WHERE clause expressions (joined with AND).
	conditions []token.Condition

	// args holds placeholder arguments for WHERE conditions.
	args []any

	// returning lists the columns to return (PostgreSQL only).
	returning []string
}

// NewDelete returns a new DeleteBuilder instance.
func NewDelete() *DeleteBuilder {
	return &DeleteBuilder{
		conditions: make([]token.Condition, 0),
		args:       make([]any, 0),
		returning:  make([]string, 0),
	}
}

// WithDialect sets the SQL dialect for escaping identifiers.
func (b *DeleteBuilder) WithDialect(d dialect.Engine) *DeleteBuilder {
	b.dialect = d
	return b
}

// From sets the table to delete table.
func (b *DeleteBuilder) From(table string) *DeleteBuilder {
	b.table = table
	return b
}

// Where adds a WHERE clause with optional placeholder arguments.
//
// Example: .Where("id = ?", 42)
func (b *DeleteBuilder) Where(condition string, params ...any) *DeleteBuilder {
	b.conditions = token.AppendCondition(
		[]token.Condition{},
		token.NewCondition(token.ConditionSimple, condition, params...),
	)
	return b
}

// AndWhere adds an AND condition.
func (b *DeleteBuilder) AndWhere(condition string, params ...any) *DeleteBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionAnd, condition, params...),
	)
	return b
}

// OrWhere adds an OR condition.
func (b *DeleteBuilder) OrWhere(condition string, params ...any) *DeleteBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionOr, condition, params...),
	)
	return b
}

// Returning specifies which columns to return after deletion (PostgreSQL only).
func (b *DeleteBuilder) Returning(columns ...string) *DeleteBuilder {
	b.returning = append(b.returning, columns...)
	return b
}

// Build compiles the DELETE SQL query and returns the statement and argument list.
func (b *DeleteBuilder) Build() (string, []any, error) {
	if b.table == "" {
		return "", nil, fmt.Errorf("no FROM table specified")
	}

	var sql strings.Builder
	var args []any

	table := b.table
	if b.dialect != nil {
		table = b.dialect.EscapeIdentifier(table)
	}

	sql.WriteString("DELETE FROM ")
	sql.WriteString(table)

	if len(b.conditions) > 0 {
		var parts []string
		for _, cond := range b.conditions {
			switch cond.Type {
			case token.ConditionSimple:
				parts = append(parts, cond.Key)
			case token.ConditionAnd, token.ConditionOr:
				parts = append(parts, fmt.Sprintf("%s %s", cond.Type, cond.Key))
			default:
				return "", nil, fmt.Errorf("invalid condition type: %s", cond.Type)
			}
			args = append(args, cond.Params...)
		}
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(parts, " "))
	}

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
