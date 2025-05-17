package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/dialect"
	"github.com/ialopezg/entiqon/internal/core/token"
)

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
func (b *UpdateBuilder) Table(table string) *UpdateBuilder {
	b.table = table
	return b
}

// Set adds column=value assignments using ordered Field slice.
func (b *UpdateBuilder) Set(column string, value any) *UpdateBuilder {
	b.assignments = append(b.assignments, token.Field(column).As(fmt.Sprintf("%v", value)))
	return b
}

// Where sets the initial condition (replaces any existing).
func (b *UpdateBuilder) Where(condition string, params ...any) *UpdateBuilder {
	b.conditions = []token.Condition{
		token.NewCondition(token.ConditionSimple, condition, params...),
	}
	return b
}

// AndWhere adds AND condition.
func (b *UpdateBuilder) AndWhere(condition string, params ...any) *UpdateBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionAnd, condition, params...),
	)
	return b
}

// OrWhere adds OR condition.
func (b *UpdateBuilder) OrWhere(condition string, params ...any) *UpdateBuilder {
	b.conditions = token.AppendCondition(
		b.conditions,
		token.NewCondition(token.ConditionOr, condition, params...),
	)
	return b
}

// WithDialect sets the SQL dialect.
func (b *UpdateBuilder) WithDialect(d dialect.Engine) *UpdateBuilder {
	b.dialect = d
	return b
}

// Build assembles the SQL UPDATE statement.
func (b *UpdateBuilder) Build() (string, []any, error) {
	if b.table == "" {
		return "", nil, fmt.Errorf("UPDATE requires a target table")
	}
	if len(b.assignments) == 0 {
		return "", nil, fmt.Errorf("UPDATE must include at least one assignment")
	}

	var sets []string
	var args []any
	for _, field := range b.assignments {
		col := field.Name
		if b.dialect != nil && !field.IsRaw {
			col = b.dialect.EscapeIdentifier(col)
		}
		sets = append(sets, fmt.Sprintf("%s = ?", col))
		args = append(args, field.Alias) // value passed via Alias field (temporary)
	}

	tokens := []string{
		fmt.Sprintf("UPDATE %s", b.table),
		fmt.Sprintf("SET %s", strings.Join(sets, ", ")),
	}

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
		tokens = append(tokens, "WHERE "+strings.Join(parts, " "))
		args = append(args, collectConditionArgs(b.conditions)...)
	}

	return strings.Join(tokens, " "), args, nil
}

func collectConditionArgs(conds []token.Condition) []any {
	var args []any
	for _, c := range conds {
		args = append(args, c.Params...)
	}
	return args
}
