// filename: /builder/delete.go

package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/builder"
	"github.com/ialopezg/entiqon/internal/core/builder/bind"
	"github.com/ialopezg/entiqon/internal/core/driver"
	core "github.com/ialopezg/entiqon/internal/core/error"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// DeleteBuilder builds DELETE SQL queries with optional WHERE and LIMIT clauses.
type DeleteBuilder struct {
	BaseBuilder
	binder     bind.ParamBinder
	table      string
	conditions []token.Condition
	limit      int
}

// NewDelete creates a new DeleteBuilder and resolves its base dialect.
// Updated: v1.4.0
func NewDelete() *DeleteBuilder {
	dialect := driver.NewGenericDialect()

	return &DeleteBuilder{
		BaseBuilder: BaseBuilder{
			dialect: dialect,
		},
		binder: *bind.NewParamBinder(dialect),
		limit:  -1,
	}
}

// UseDialect overrides the dialect for the delete builder.
//
// Updated: v1.4.0
func (b *DeleteBuilder) UseDialect(name string) *DeleteBuilder {
	b.BaseBuilder.UseDialect(name)
	return b
}

// From sets the table to delete from.
func (b *DeleteBuilder) From(table string) *DeleteBuilder {
	if table == "" {
		b.AddStageError("FROM", fmt.Errorf("table is empty"))
	} else {
		b.table = table
	}
	return b
}

// Where sets the initial WHERE condition and resets previous ones.
func (b *DeleteBuilder) Where(condition string, values ...any) *DeleteBuilder {
	c := token.NewCondition(token.ConditionSimple, condition, values...)
	if !c.IsValid() {
		b.errors.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = append([]token.Condition{}, c)
	return b
}

// AndWhere adds an AND condition.
func (b *DeleteBuilder) AndWhere(condition string, values ...any) *DeleteBuilder {
	c := token.NewCondition(token.ConditionAnd, condition, values...)
	if !c.IsValid() {
		b.errors.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// OrWhere adds OR condition.
func (b *DeleteBuilder) OrWhere(condition string, values ...any) *DeleteBuilder {
	c := token.NewCondition(token.ConditionOr, condition, values...)
	if !c.IsValid() {
		b.errors.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// Limit sets a row limit on the DELETE operation.
func (b *DeleteBuilder) Limit(n int) *DeleteBuilder {
	b.limit = n
	return b
}

// Build assembles the DELETE SQL query with placeholders.
// Updated: v1.4.0
func (b *DeleteBuilder) Build() (string, []any, error) {
	var dialect = b.dialect
	if !b.HasDialect() {
		dialect = b.GetDialect()
	}

	if b.HasErrors() {
		return "", nil, fmt.Errorf("DELETE: %d invalid condition(s)", len(b.GetErrors()))
	}
	if b.table == "" {
		return "", nil, fmt.Errorf("DELETE: requires a target table")
	}

	tokens := []string{"DELETE FROM", dialect.QuoteIdentifier(b.table)}
	var args []any

	if len(b.conditions) > 0 {
		binder := bind.NewParamBinderWithPosition(dialect, len(args)+1)
		whereClause, condArgs, err := builder.RenderConditionsWithBinder(dialect, b.conditions, binder)
		if err != nil {
			return "", nil, fmt.Errorf("UPDATE: %w", err)
		}
		tokens = append(tokens, "WHERE", whereClause)
		args = append(args, condArgs...)
	}

	if b.limit >= 0 {
		limitClause := dialect.BuildLimitOffset(b.limit, -1)
		if limitClause != "" {
			tokens = append(tokens, limitClause)
		}
	}

	return strings.Join(tokens, " "), args, nil
}
