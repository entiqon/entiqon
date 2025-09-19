// File: db/builder/select.go

package builder

import (
	"fmt"
	"strings"

	"github.com/entiqon/db/driver"
	bt "github.com/entiqon/db/internal/build/token"
	"github.com/entiqon/db/internal/core/builder"
	"github.com/entiqon/db/internal/core/builder/bind"
	"github.com/entiqon/db/internal/core/errors"
	ct "github.com/entiqon/db/internal/core/token"
)

// SelectBuilder builds a SQL SELECT query using fluent method chaining.
//
// It supports basic querying with WHERE conditions, ordering, and pagination.
type SelectBuilder struct {
	BaseBuilder
	columns    []*bt.Column
	sources    []*bt.Table
	conditions []ct.Condition
	sorting    []string
	take       *int
	skip       *int
}

// NewSelect creates a new SelectBuilder using the given SQL dialect.
//
// If the provided dialect is nil, it defaults to driver.NewGenericDialect().
// The builder name is automatically set to "select".
//
// Since: v1.4.0
func NewSelect(dialect driver.Dialect) *SelectBuilder {
	if dialect == nil {
		dialect = driver.NewGenericDialect()
	}
	base := NewBaseBuilder("select", dialect)

	return &SelectBuilder{
		BaseBuilder: base,
		columns:     make([]*bt.Column, 0),
		sources:     make([]*bt.Table, 0),
		conditions:  make([]ct.Condition, 0),
		sorting:     make([]string, 0),
	}
}

// Select sets the list of columns to be used in the SELECT clause.
//
// Columns are parsed from strings using util.ParseColumns(...), and support
// inline aliasing (e.g., "id AS uid", "user_id uid"). Any invalid column
// is stored and its error tracked in the builder validator.
//
// # Examples
//
//	Select("id", "name AS customer")
//	  → SELECT id, name AS customer
//
//	Select("users.id AS uid", "users.name")
//	  → SELECT users.id AS uid, users.name
//
//	Select("id, name") // multiple in one string
//	  → SELECT id, name
func (b *SelectBuilder) Select(columns ...string) *SelectBuilder {
	return b.ClearSelect().addColumns(columns...)
}

// AddSelect appends one or more columns to the SELECT clause of the query,
// preserving any previously defined columns.
//
// This is complementary to Select(...), which replaces the column list entirely.
// AddSelect should ideally be called after Select(...) to preserve the semantic
// flow of query construction.
//
// Note: If AddSelect is called before Select, the internal column list will be
// automatically initialized. However, this is not considered best practice.
// It is recommended to follow a hierarchical flow — define the base columns with
// Select(...) first, then extend with AddSelect(...) as needed.
//
// Each input string may contain a single column or a comma-separated list.
// The method internally uses util.ParseColumns to handle splitting, trimming,
// and inline aliasing (e.g., "id AS user_id").
//
// Invalid columns are included with their Error field populated and may be
// filtered or logged by downstream logic.
//
// Examples:
//
//	Select("id")
//	AddSelect("name")
//	  → SELECT id, name
//
//	AddSelect("id, name")
//	  → SELECT id, name
//
//	AddSelect("id", "name AS customer")
//	  → SELECT id, name AS customer
func (b *SelectBuilder) AddSelect(columns ...string) *SelectBuilder {
	return b.addColumns(columns...)
}

// From sets the target table for the SELECT statement.
//
// Since: v0.0.1
// Updated: v1.5.0
func (b *SelectBuilder) From(table string, alias ...string) *SelectBuilder {
	source := bt.NewTable(table, alias...)
	if !source.IsValid() {
		b.Validator.AddStageError(errors.StageFrom,
			fmt.Errorf("invalid column: %s — %v", source.String(), source.GetError()))
	}
	b.sources = append(b.sources, bt.NewTable(table, alias...))
	return b
}

// Where sets the base condition(s) for the WHERE clause.
// It resets any previously added conditions.
func (b *SelectBuilder) Where(condition string, values ...any) *SelectBuilder {
	c := ct.NewCondition(ct.ConditionSimple, condition, values...)
	if !c.IsValid() {
		b.AddStageError(errors.StageWhere, c.Error)
	}
	b.conditions = []ct.Condition{c}
	return b
}

// AndWhere adds an AND condition to the WHERE clause.
func (b *SelectBuilder) AndWhere(condition string, values ...any) *SelectBuilder {
	c := ct.NewCondition(ct.ConditionAnd, condition, values...)
	if !c.IsValid() {
		b.AddStageError(errors.StageWhere, c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// OrWhere adds an OR condition to the WHERE clause.
func (b *SelectBuilder) OrWhere(condition string, values ...any) *SelectBuilder {
	c := ct.NewCondition(ct.ConditionOr, condition, values...)
	if !c.IsValid() {
		b.AddStageError(errors.StageWhere, c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// OrderBy appends a column or expression to the ORDER BY clause.
func (b *SelectBuilder) OrderBy(column string) *SelectBuilder {
	b.sorting = append(b.sorting, column)
	return b
}

// Take limits the number of rows returned by the query (engine-agnostic equivalent).
func (b *SelectBuilder) Take(value int) *SelectBuilder {
	b.take = &value
	return b
}

// Skip offsets the rows returned by the query (engine-agnostic equivalent).
func (b *SelectBuilder) Skip(value int) *SelectBuilder {
	b.skip = &value
	return b
}

// UseDialect resolves and applies the dialect by name (e.g., "postgres").
// It replaces any previously set dialect on the builder.
func (b *SelectBuilder) UseDialect(name string) *SelectBuilder {
	b.BaseBuilder.dialect = driver.ResolveDialect(name)
	return b
}

// Build compiles the SELECT statement and returns it as a string and argument list.
// If the FROM clause is missing, an error is returned.
// Dialect rules (quoting, pagination) are applied if configured.
func (b *SelectBuilder) Build() (string, []any, error) {
	if b.Dialect == nil {
		b.Dialect = driver.NewGenericDialect()
	}
	if len(b.sources) == 0 {
		return "", nil, fmt.Errorf("%s: missing source; expected at least one table source", errors.StageFrom)
	}

	var tokens []string

	// prepare columns
	columns := "*"
	if len(b.columns) > 0 {
		single := len(b.sources) == 1
		var defaultAlias string
		if single && b.sources[0].IsAliased() {
			// Only pick up an explicit alias
			defaultAlias = b.sources[0].GetAlias()
		}

		var parts []string
		for _, col := range b.columns {
			rendered := col.Render(b.Dialect)
			if single {
				col = col.WithTable(b.sources[0])
				part := col.RenderAlias(b.Dialect, col.GetName())
				// Single‐source: prefix only if we have an explicit alias
				if defaultAlias != "" {
					parts = append(parts, fmt.Sprintf("%s.%s", defaultAlias, part))
				} else {
					parts = append(parts, part)
				}
			} else {
				// Multiple sources: must have a Table and it must be registered
				if col.Table == nil {
					return "", nil, fmt.Errorf(
						"%s: column %q must have an owner table in multi-source SELECT",
						errors.StageSelect, col.GetName(),
					)
				}
				found := false
				for _, src := range b.sources {
					if src == col.Table {
						found = true
						break
					}
				}
				if !found {
					return "", nil, fmt.Errorf(
						"%s: column %q refers to table %q which is not in builder sources",
						errors.StageSelect, col.GetName(), col.Table.GetName(),
					)
				}
				parts = append(parts, fmt.Sprintf("%s.%s", col.Table.AliasOr(), rendered))
			}
		}

		columns = strings.Join(parts, ", ")
	}

	tokens = append([]string{}, fmt.Sprintf("SELECT %s", columns))

	// prepare sources
	if len(b.sources) > 0 {
		var fromParts []string
		for _, tbl := range b.sources {
			if out := tbl.Render(b.Dialect); out != "" {
				fromParts = append(fromParts, out)
			}
		}
		if len(fromParts) > 0 {
			tokens = append(tokens, "FROM", strings.Join(fromParts, ", "))
		}
	}

	// prepare conditions
	var args []any
	if len(b.conditions) > 0 {
		binder := bind.NewParamBinderWithPosition(b.Dialect, len(args)+1)
		whereClause, clauseArgs, err := builder.RenderConditionsWithBinder(b.Dialect, b.conditions, binder)
		if err != nil {
			return "", nil, fmt.Errorf("%s: %w", errors.StageWhere, err)
		}
		tokens = append(tokens, "WHERE", whereClause)
		args = append(args, clauseArgs...)
	}

	// prepare sorting
	if len(b.sorting) > 0 {
		tokens = append(tokens, "ORDER BY "+strings.Join(b.sorting, ", "))
	}

	// prepare pagination, if applicable
	limit, offset := -1, -1
	if b.take != nil {
		limit = *b.take
	}
	if b.skip != nil {
		offset = *b.skip
	}

	if limit > 0 && offset > 0 {
		tokens = append(tokens, b.Dialect.BuildLimitOffset(limit, offset))
	} else {
		if limit > 0 {
			tokens = append(tokens, fmt.Sprintf("LIMIT %d", limit))
		}
		if offset > 0 {
			tokens = append(tokens, fmt.Sprintf("OFFSET %d", offset))
		}
	}

	// validate for errors
	if err := b.Validate(); err != nil {
		return "", nil, err
	}

	// render
	return strings.Join(tokens, " "), args, nil
}

// ClearSelect removes all previously selected columns.
//
// This is used internally by Select(...) to ensure the SELECT clause
// reflects only the explicitly provided columns.
func (b *SelectBuilder) ClearSelect() *SelectBuilder {
	b.columns = make([]*bt.Column, 0)
	return b
}

// addColumns parses and appends the given column expressions,
// optionally associating them with a table token if there is
// exactly one valid source defined.
//
// This is a shared internal method between Select(...) and AddSelect(...).
// It delegates token creation to util.ParseColumns(...) and
// passes any applicable source token for column qualification.
func (b *SelectBuilder) addColumns(columns ...string) *SelectBuilder {
	var table *bt.Table
	if len(b.sources) == 1 && b.sources[0].IsValid() {
		table = b.sources[0]
	}
	b.appendColumns(bt.NewColumnsFrom(columns...), table)
	return b
}

// appendColumns stores a list of Column tokens into the builder,
// applying optional source resolution from a provided table token.
//
// Each column is validated for structural correctness. If a column
// is invalid, a StageSelect error is recorded in the builder's validator.
// If a table token is provided and valid, it is assigned to each column
// via WithTable, allowing later rendering to reference the table alias.
//
// This method is used internally by addColumns(...), which handles
// source detection and expression parsing.
//
// # Example
//
//	users := token.NewTable("users AS u")
//	cols := util.ParseColumns("id", "email")
//
//	b.appendColumns(cols, &users)
//
//	// Rendered: SELECT u.id, u.email FROM users AS u
func (b *SelectBuilder) appendColumns(cols []*bt.Column,
	table *bt.Table) {
	for _, col := range cols {
		if col.IsErrored() {
			b.Validator.AddStageError(errors.StageSelect, fmt.Errorf("invalid column: %s", col.String()))
		}

		// Assign table for qualification and rendering
		if table != nil && table.IsValid() {
			col.WithTable(table)
		}

		b.columns = append(b.columns, col)
	}
}
