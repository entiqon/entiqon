package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/builder"
	"github.com/ialopezg/entiqon/internal/core/builder/bind"
	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// UpdateBuilder builds a SQL UPDATE query with fluent syntax and dialect_engine.md support.
type UpdateBuilder struct {
	BaseBuilder
	table       string
	assignments []token.FieldToken
	conditions  []token.Condition
}

// NewUpdate creates a new UpdateBuilder.
// NewUpdate creates a new instance of UpdateBuilder.
// This is the entry point for building a SQL UPDATE query fluently.
func NewUpdate() *UpdateBuilder {
	return &UpdateBuilder{
		BaseBuilder: BaseBuilder{dialect: driver.NewGenericDialect()},
		assignments: []token.FieldToken{},
		conditions:  []token.Condition{},
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
	b.assignments = append(b.assignments, token.Field(column).WithValue(value))
	return b
}

// Where sets the base WHERE clause.
// Where sets the base WHERE clause as a simple condition.
// This replaces any existing conditions.
func (b *UpdateBuilder) Where(condition string, values ...any) *UpdateBuilder {
	c := token.NewCondition(token.ConditionSimple, condition, values)
	if !c.IsValid() {
		b.errors = append(b.errors, builder.Error{
			Token:  "WHERE",
			Errors: []error{c.Error},
		})
	}
	b.conditions = append([]token.Condition{}, c)
	return b
}

// AndWhere adds an AND clause.
// AndWhere appends a condition with an AND operator to the current WHERE clause.
func (b *UpdateBuilder) AndWhere(condition string, values ...any) *UpdateBuilder {
	c := token.NewCondition(token.ConditionAnd, condition, values)
	if !c.IsValid() {
		b.errors = append(b.errors, builder.Error{
			Token:  "WHERE",
			Errors: []error{c.Error},
		})
	}
	b.conditions = append(b.conditions, c)
	return b
}

// OrWhere adds an OR clause.
// OrWhere appends a condition with an OR operator to the current WHERE clause.
func (b *UpdateBuilder) OrWhere(condition string, values ...any) *UpdateBuilder {
	c := token.NewCondition(token.ConditionOr, condition, values)
	if !c.IsValid() {
		b.errors = append(b.errors, builder.Error{
			Token:  "WHERE",
			Errors: []error{c.Error},
		})
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
	if !b.HasDialect() {
		_ = b.GetDialect()
	}

	if b.HasErrors() {
		return "", nil, fmt.Errorf("UPDATE: %d invalid condition(s)", len(b.GetErrors()))
	}

	if b.table == "" {
		return "", nil, fmt.Errorf("UPDATE: requires a target table")
	}
	if len(b.assignments) == 0 {
		return "", nil, fmt.Errorf("UPDATE: must define at least one column assignment")
	}

	dialect := b.GetDialect()
	var tokens []string
	var sets []string
	var args []any

	for _, f := range b.assignments {
		if f.Alias != "" {
			return "", nil, fmt.Errorf("UPDATE: column aliasing is not supported: '%s AS %s'", f.Name, f.Alias)
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

	return strings.Join(tokens, " "), args, nil
}
