// File: db/builder/update.go

package builder

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/driver"
	"github.com/entiqon/entiqon/db/internal/core/builder"
	"github.com/entiqon/entiqon/db/internal/core/builder/bind"
	core "github.com/entiqon/entiqon/db/internal/core/errors"
	token2 "github.com/entiqon/entiqon/db/internal/core/token"
)

// UpdateBuilder builds a SQL UPDATE query with fluent syntax and dialect_engine.md support.
type UpdateBuilder struct {
	BaseBuilder
	table       string
	assignments []token2.FieldToken
	conditions  []token2.Condition
}

// NewUpdate creates a new UpdateBuilder using the given SQL dialect.
//
// If the provided dialect is nil, it defaults to driver.NewGenericDialect().
// The builder name is automatically set to "update".
//
// Since: v1.4.0
func NewUpdate(dialect driver.Dialect) *UpdateBuilder {
	base := NewBaseBuilder("update", dialect)

	return &UpdateBuilder{
		BaseBuilder: base,
		assignments: []token2.FieldToken{},
		conditions:  []token2.Condition{},
	}
}

// Table sets the table to update.
// Table sets the target table for the UPDATE query.
// The table name will be quoted using the selected dialect during Build().
func (b *UpdateBuilder) Table(name string) *UpdateBuilder {
	b.table = name
	return b
}

// Set adds a column assignment using ordered Field.
// The value is stored temporarily in .Alias for consistency.
// Set adds a column=value assignment to the UPDATE statement.
// The value will be stored and interpolated with placeholders.
// NOTE: Alias use is disallowed and will raise an error during Build().
func (b *UpdateBuilder) Set(column string, value any) *UpdateBuilder {
	b.assignments = append(b.assignments, token2.Field(column).WithValue(value))
	return b
}

// Where sets the base WHERE clause.
// Where sets the base WHERE clause as a simple condition.
// This replaces any existing conditions.
func (b *UpdateBuilder) Where(condition string, values ...any) *UpdateBuilder {
	c := token2.NewCondition(token2.ConditionSimple, condition, values)
	if !c.IsValid() {
		b.errors.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = append([]token2.Condition{}, c)
	return b
}

// AndWhere adds an AND clause.
// AndWhere appends a condition with an AND operator to the current WHERE clause.
func (b *UpdateBuilder) AndWhere(condition string, values ...any) *UpdateBuilder {
	c := token2.NewCondition(token2.ConditionAnd, condition, values)
	if !c.IsValid() {
		b.errors.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// OrWhere adds an OR clause.
// OrWhere appends a condition with an OR operator to the current WHERE clause.
func (b *UpdateBuilder) OrWhere(condition string, values ...any) *UpdateBuilder {
	c := token2.NewCondition(token2.ConditionOr, condition, values)
	if !c.IsValid() {
		b.errors.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// UseDialect resolves and applies the dialect_engine.md by name (e.g., "postgres").
// It replaces any previously set dialect_engine.md on the builder.
func (b *UpdateBuilder) UseDialect(name string) *UpdateBuilder {
	b.BaseBuilder.UseDialect(name)
	return b
}

// Build renders the UPDATE SQL query and returns the query + args.
func (b *UpdateBuilder) Build() (string, []any, error) {
	if b.table == "" {
		b.Validator.AddStageError(core.StageFrom, fmt.Errorf("requires a target table"))
	}
	if len(b.assignments) == 0 {
		b.Validator.AddStageError(core.StageSet, fmt.Errorf("must define at least one column assignment"))
	}

	dialect := b.GetDialect()
	var tokens []string
	var sets []string
	var args []any

	for _, f := range b.assignments {
		if f.Alias != "" {
			b.Validator.AddStageError(
				core.StageSet,
				fmt.Errorf("column aliasing is not supported: '%s AS %s'", f.Name, f.Alias),
			)
			continue
		}

		name := f.Name
		if !f.IsRaw {
			name = dialect.QuoteIdentifier(name)
		}

		placeholder := dialect.Placeholder(len(args) + 1)
		sets = append(sets, fmt.Sprintf("%s = %s", name, placeholder))
		args = append(args, f.Value)
	}

	tokens = append(tokens, "UPDATE", dialect.QuoteIdentifier(b.table), "SET", strings.Join(sets, ", "))

	if len(b.conditions) > 0 {
		binder := bind.NewParamBinderWithPosition(dialect, len(args)+1)
		whereClause, condArgs, err := builder.RenderConditionsWithBinder(dialect, b.conditions, binder)
		if err != nil {
			return "", nil, fmt.Errorf("UPDATE: %w", err)
		}
		tokens = append(tokens, "WHERE", whereClause)
		args = append(args, condArgs...)
	}

	if err := b.Validate(); err != nil {
		return "", nil, err
	}

	return strings.Join(tokens, " "), args, nil
}
