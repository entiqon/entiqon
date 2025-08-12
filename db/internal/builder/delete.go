// File: db/builder/delete.go
// Description: Provides DeleteBuilder for constructing DELETE SQL statements.
// Since: v1.4.0

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

// DeleteBuilder builds a SQL DELETE statement with optional WHERE and LIMIT clauses.
type DeleteBuilder struct {
	BaseBuilder

	binder     bind.ParamBinder
	table      string
	alias      string
	conditions []token2.Condition
	limit      int
}

// NewDelete creates a new DeleteBuilder using the given SQL dialect.
//
// If the provided dialect is nil, it defaults to driver.NewGenericDialect().
// The builder name is automatically set to "delete".
//
// Since: v1.4.0
func NewDelete(dialect driver.Dialect) *DeleteBuilder {
	base := NewBaseBuilder("delete", dialect)

	return &DeleteBuilder{
		BaseBuilder: base,
		binder:      *bind.NewParamBinder(base.Dialect),
		limit:       -1,
	}
}

// From sets the target table for the DELETE operation.
func (b *DeleteBuilder) From(table string, alias ...string) *DeleteBuilder {
	if table == "" {
		b.Validator.AddStageError(core.StageFrom, fmt.Errorf("table is empty"))
	}
	b.table = table
	if len(alias) > 0 {
		b.alias = strings.TrimSpace(alias[0])
	}
	return b
}

// Where sets the initial WHERE clause for the DELETE statement,
// replacing any previously defined conditions.
func (b *DeleteBuilder) Where(condition string, values ...any) *DeleteBuilder {
	c := token2.NewCondition(token2.ConditionSimple, condition, values...)
	if !c.IsValid() {
		b.Validator.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = []token2.Condition{c}
	return b
}

// AndWhere appends a condition to the WHERE clause using logical AND.
func (b *DeleteBuilder) AndWhere(condition string, values ...any) *DeleteBuilder {
	c := token2.NewCondition(token2.ConditionAnd, condition, values...)
	if !c.IsValid() {
		b.Validator.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// OrWhere appends a condition to the WHERE clause using logical OR.
func (b *DeleteBuilder) OrWhere(condition string, values ...any) *DeleteBuilder {
	c := token2.NewCondition(token2.ConditionOr, condition, values...)
	if !c.IsValid() {
		b.Validator.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// Limit sets the maximum number of rows to delete.
func (b *DeleteBuilder) Limit(n int) *DeleteBuilder {
	b.limit = n
	return b
}

// Build assembles the DELETE SQL query and returns the final SQL string and arguments.
//
// If validation fails, an error is returned describing any missing elements or invalid conditions.
//
// Returns:
//   - SQL string (e.g., DELETE FROM users WHERE id = $1)
//   - Arguments for parameterized execution
//   - Error if the builder state is invalid
//
// Updated: v1.4.0
func (b *DeleteBuilder) Build() (string, []any, error) {
	if b.table == "" {
		b.Validator.AddStageError(core.StageFrom, fmt.Errorf("table is empty"))
	}

	var whereClause string
	var args []any

	if len(b.conditions) > 0 {
		binder := bind.NewParamBinderWithPosition(b.Dialect, 1)
		var condErr error
		whereClause, args, condErr = builder.RenderConditionsWithBinder(b.Dialect, b.conditions, binder)
		if condErr != nil {
			b.Validator.AddStageError(core.StageWhere, condErr)
		}
	}

	if err := b.Validate(); err != nil {
		return "", nil, err
	}

	tokens := []string{"DELETE FROM", b.Dialect.QuoteIdentifier(b.table)}

	if whereClause != "" {
		tokens = append(tokens, "WHERE", whereClause)
	}

	if b.limit >= 0 {
		if limitClause := b.Dialect.BuildLimitOffset(b.limit, -1); limitClause != "" {
			tokens = append(tokens, limitClause)
		}
	}

	return strings.Join(tokens, " "), args, nil
}
