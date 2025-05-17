package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/dialect"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// UpdateBuilder builds a SQL UPDATE statement.
//
// It supports setting fields and composing WHERE clauses with arguments.
type UpdateBuilder struct {
	dialect     dialect.Engine
	table       string            // target table
	assignments map[string]any    // column-value pairs for SET
	conditions  []token.Condition // raw SQL conditions
	args        []any             // arguments for WHERE placeholders
}

// NewUpdate returns a new UpdateBuilder instance.
func NewUpdate() *UpdateBuilder {
	return &UpdateBuilder{
		assignments: make(map[string]any),
		conditions:  make([]token.Condition, 0),
		args:        make([]any, 0),
	}
}

// WithDialect sets the SQL dialect for escaping identifiers.
func (b *UpdateBuilder) WithDialect(d dialect.Engine) *UpdateBuilder {
	b.dialect = d
	return b
}

// Table sets the table name to update.
func (b *UpdateBuilder) Table(name string) *UpdateBuilder {
	b.table = name
	return b
}

// Set defines a field and value to update.
func (b *UpdateBuilder) Set(field string, value any) *UpdateBuilder {
	b.assignments[field] = value
	return b
}

// Where adds a WHERE clause with placeholders and binds arguments.
//
// Example:
//
//	.Where("status = ? AND created_at > ?", "active", "2023-01-01")
func (b *UpdateBuilder) Where(condition string, params ...any) *UpdateBuilder {
	b.conditions = []token.Condition{}
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionSimple, condition, params...),
	)
	return b
}

// AndWhere adds an AND condition.
func (b *UpdateBuilder) AndWhere(condition string, params ...any) *UpdateBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionAnd, condition, params...),
	)
	return b
}

// OrWhere adds an OR condition.
func (b *UpdateBuilder) OrWhere(condition string, params ...any) *UpdateBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionOr, condition, params...),
	)
	return b
}

// Build compiles the UPDATE statement into a SQL string and a list of arguments.
//
// It returns an error if the table is missing or no SET fields are defined.
func (b *UpdateBuilder) Build() (string, []any, error) {
	if b.table == "" {
		return "", nil, fmt.Errorf("no table specified")
	}
	if len(b.assignments) == 0 {
		return "", nil, fmt.Errorf("no assignments provided")
	}

	var sql strings.Builder
	var args []any

	table := b.table
	if b.dialect != nil {
		table = b.dialect.EscapeIdentifier(table)
	}

	sql.WriteString("UPDATE ")
	sql.WriteString(table)
	sql.WriteString(" SET ")

	assignments := make([]string, 0, len(b.assignments))
	for col, val := range b.assignments {
		colName := col
		if b.dialect != nil {
			colName = b.dialect.EscapeIdentifier(col)
		}
		assignments = append(assignments, fmt.Sprintf("%s = ?", colName))
		args = append(args, val)
	}
	sql.WriteString(strings.Join(assignments, ", "))

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

	return sql.String(), args, nil
}
