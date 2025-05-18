package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/dialect"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// UpdateBuilder builds a SQL UPDATE query with fluent syntax and dialect support.
type UpdateBuilder struct {
	dialect     dialect.Engine
	table       string
	assignments []token.FieldToken
	conditions  []token.Condition
}

// NewUpdate creates a new UpdateBuilder.
func NewUpdate() *UpdateBuilder {
	return &UpdateBuilder{
		assignments: []token.FieldToken{},
		conditions:  []token.Condition{},
	}
}

// Table sets the table to update.
func (b *UpdateBuilder) Table(name string) *UpdateBuilder {
	b.table = name
	return b
}

// Set adds a column assignment using ordered Field.
// The value is stored temporarily in .Alias for consistency.
func (b *UpdateBuilder) Set(column string, value any) *UpdateBuilder {
	b.assignments = append(b.assignments, token.Field(column).WithValue(value))
	return b
}

// Where sets the base WHERE clause.
func (b *UpdateBuilder) Where(condition string, params ...any) *UpdateBuilder {
	b.conditions = []token.Condition{
		token.NewCondition(token.ConditionSimple, condition, params...),
	}
	return b
}

// AndWhere adds an AND clause.
func (b *UpdateBuilder) AndWhere(condition string, params ...any) *UpdateBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionAnd, condition, params...),
	)
	return b
}

// OrWhere adds an OR clause.
func (b *UpdateBuilder) OrWhere(condition string, params ...any) *UpdateBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionOr, condition, params...),
	)
	return b
}

// Build renders the UPDATE SQL query and returns the query + args.
func (b *UpdateBuilder) Build() (string, []any, error) {
	if b.table == "" {
		return "", nil, fmt.Errorf("UPDATE requires a target table")
	}
	if len(b.assignments) == 0 {
		return "", nil, fmt.Errorf("UPDATE must define at least one column assignment")
	}

	var sets []string
	var args []any

	for _, f := range b.assignments {
		if f.Alias != "" {
			return "", nil, fmt.Errorf("UPDATE does not support column aliasing: '%s AS %s'", f.Name, f.Alias)
		}

		name := f.Name
		if b.dialect != nil && !f.IsRaw {
			name = b.dialect.EscapeIdentifier(name)
		}

		sets = append(sets, fmt.Sprintf("%s = ?", name))
		args = append(args, f.Value)
	}

	sql := fmt.Sprintf("UPDATE %s SET %s", b.table, strings.Join(sets, ", "))

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

// WithDialect sets the SQL dialect for escaping.
func (b *UpdateBuilder) WithDialect(d dialect.Engine) *UpdateBuilder {
	b.dialect = d
	return b
}

func collectConditionArgs(conditions []token.Condition) []any {
	var args []any
	for _, c := range conditions {
		args = append(args, c.Params...)
	}
	return args
}
