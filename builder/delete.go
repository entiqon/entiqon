// filename: /builder/delete.go

package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/builder"
	"github.com/ialopezg/entiqon/internal/core/builder/bind"
	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// DeleteBuilder builds DELETE SQL queries with optional WHERE and LIMIT clauses.
type DeleteBuilder struct {
	BaseBuilder

	dialect    driver.Dialect
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
		dialect: dialect,
		binder:  *bind.NewParamBinder(dialect),
		limit:   -1,
	}
}

// UseDialect overrides the dialect for the delete builder.
// Updated: v1.4.0
func (db *DeleteBuilder) UseDialect(d driver.Dialect) *DeleteBuilder {
	db.BaseBuilder.UseDialect(d)
	return db
}

// From sets the table to delete from.
func (db *DeleteBuilder) From(table string) *DeleteBuilder {
	db.table = table
	return db
}

// Where sets the initial WHERE condition and resets previous ones.
func (db *DeleteBuilder) Where(condition string, values ...any) *DeleteBuilder {
	c := token.NewCondition(token.ConditionSimple, condition, values...)
	if !c.IsValid() {
		db.errors = append(db.errors, builder.Error{
			Token:  "WHERE",
			Errors: []error{c.Error},
		})
	}
	db.conditions = append([]token.Condition{}, c)
	return db
}

// AndWhere adds an AND condition.
func (db *DeleteBuilder) AndWhere(condition string, values ...any) *DeleteBuilder {
	c := token.NewCondition(token.ConditionAnd, condition, values...)
	if !c.IsValid() {
		db.errors = append(db.errors, builder.Error{
			Token:  "WHERE",
			Errors: []error{c.Error},
		})
	}
	db.conditions = append(db.conditions, c)
	return db
}

// OrWhere adds OR condition.
func (db *DeleteBuilder) OrWhere(condition string, values ...any) *DeleteBuilder {
	c := token.NewCondition(token.ConditionOr, condition, values...)
	if !c.IsValid() {
		db.errors = append(db.errors, builder.Error{
			Token:  "WHERE",
			Errors: []error{c.Error},
		})
	}
	db.conditions = append(db.conditions, c)
	return db
}

// Limit sets a row limit on the DELETE operation.
func (db *DeleteBuilder) Limit(n int) *DeleteBuilder {
	db.limit = n
	return db
}

// Build assembles the DELETE SQL query with placeholders.
// Updated: v1.4.0
func (db *DeleteBuilder) Build() (string, []any, error) {
	if db.HasErrors() {
		return "", nil, fmt.Errorf("delete builder: %d invalid condition(s)", len(db.GetErrors()))
	}

	if db.table == "" {
		return "", nil, fmt.Errorf("DELETE requires a target table")
	}

	dialect := db.resolveDialect()

	tokens := []string{"DELETE FROM", dialect.QuoteIdentifier(db.table)}

	var args []any
	if len(db.conditions) > 0 {
		whereClause, clauseArgs, err := builder.RenderConditions(dialect, db.conditions)
		if err != nil {
			return "", nil, fmt.Errorf("delete builder: %w", err)
		}
		tokens = append(tokens, "WHERE", whereClause)
		args = append(args, clauseArgs...)
	}

	if db.limit >= 0 {
		limitClause := dialect.BuildLimitOffset(db.limit, -1)
		if limitClause != "" {
			tokens = append(tokens, limitClause)
		}
	}

	return strings.Join(tokens, " "), args, nil
}
