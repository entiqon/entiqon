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
}

// NewDelete creates and returns a new DeleteBuilder.
func NewDelete() *DeleteBuilder {
	return &DeleteBuilder{
		conditions: []token.Condition{},
	}
}

// From sets the target table to delete table.
func (b *DeleteBuilder) From(table string) *DeleteBuilder {
	b.table = table
	return b
}

// Where starts the WHERE clause.
func (b *DeleteBuilder) Where(condition string, params ...any) *DeleteBuilder {
	b.conditions = []token.Condition{
		token.NewCondition(token.ConditionSimple, condition, params...),
	}
	return b
}

// AndWhere appends an AND condition.
func (b *DeleteBuilder) AndWhere(condition string, params ...any) *DeleteBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionAnd, condition, params...),
	)
	return b
}

// OrWhere appends an OR condition.
func (b *DeleteBuilder) OrWhere(condition string, params ...any) *DeleteBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionOr, condition, params...),
	)
	return b
}

// WithDialect sets the dialect engine for identifier escaping.
func (b *DeleteBuilder) WithDialect(e dialect.Engine) *DeleteBuilder {
	b.dialect = e
	return b
}

// Build compiles the DELETE SQL query and returns it with parameter args.
func (b *DeleteBuilder) Build() (string, []any, error) {
	if strings.TrimSpace(b.table) == "" {
		return "", nil, fmt.Errorf("DELETE requires a target table")
	}

	table := b.table
	if b.dialect != nil {
		table = b.dialect.EscapeIdentifier(table)
	}

	sql := fmt.Sprintf("DELETE FROM %s", table)
	var args []any

	if len(b.conditions) > 0 {
		var parts []string
		for _, c := range b.conditions {
			switch c.Type {
			case token.ConditionSimple:
				parts = append(parts, c.Key)
			case token.ConditionAnd, token.ConditionOr:
				parts = append(parts, fmt.Sprintf("%s %s", c.Type, c.Key))
			default:
				return "", nil, fmt.Errorf("invalid condition type: %s", c.Type)
			}
		}
		sql += " WHERE " + strings.Join(parts, " ")
		args = append(args, collectConditionArgs(b.conditions)...)
	}

	return sql, args, nil
}
