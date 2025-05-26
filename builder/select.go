package builder

import (
	"fmt"
	"strings"

	driver2 "github.com/ialopezg/entiqon/driver"
	"github.com/ialopezg/entiqon/internal/build/render"
	"github.com/ialopezg/entiqon/internal/build/token"
	"github.com/ialopezg/entiqon/internal/build/util"
	"github.com/ialopezg/entiqon/internal/core/builder"
	"github.com/ialopezg/entiqon/internal/core/builder/bind"
	core "github.com/ialopezg/entiqon/internal/core/error"
	internal "github.com/ialopezg/entiqon/internal/core/token"
)

// SelectBuilder builds a SQL SELECT query using fluent method chaining.
//
// It supports basic querying with WHERE conditions, ordering, and pagination.
type SelectBuilder struct {
	BaseBuilder
	columns    []token.Column
	sources    []internal.Table
	conditions []internal.Condition
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
func NewSelect(dialect driver2.Dialect) *SelectBuilder {
	base := NewBaseBuilder("select", dialect)

	return &SelectBuilder{
		BaseBuilder: base,
		conditions:  make([]internal.Condition, 0),
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
	b.columns = []token.Column{}
	b.appendColumns(util.ParseColumns(columns...))
	return b
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
	if b.columns == nil {
		b.columns = make([]token.Column, 0)
	}
	b.appendColumns(util.ParseColumns(columns...))
	return b
}

// From sets the target table for the SELECT statement.
//
// Since: v0.0.1
// Updated: v1.5.0
func (b *SelectBuilder) From(table string, alias ...string) *SelectBuilder {
	if table == "" {
		b.Validator.AddStageError(core.StageFrom, fmt.Errorf("table is empty"))
	}
	if len(alias) > 0 {
		b.sources = append(b.sources, internal.NewTableWithAlias(table, alias[0]))
	} else {
		b.sources = append(b.sources, internal.NewTable(table))
	}
	return b
}

// Where sets the base condition(s) for the WHERE clause.
// It resets any previously added conditions.
func (b *SelectBuilder) Where(condition string, values ...any) *SelectBuilder {
	c := internal.NewCondition(internal.ConditionSimple, condition, values...)
	if !c.IsValid() {
		b.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = []internal.Condition{c}
	return b
}

// AndWhere adds an AND condition to the WHERE clause.
func (b *SelectBuilder) AndWhere(condition string, values ...any) *SelectBuilder {
	c := internal.NewCondition(internal.ConditionAnd, condition, values...)
	if !c.IsValid() {
		b.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// OrWhere adds an OR condition to the WHERE clause.
func (b *SelectBuilder) OrWhere(condition string, values ...any) *SelectBuilder {
	c := internal.NewCondition(internal.ConditionOr, condition, values...)
	if !c.IsValid() {
		b.AddStageError(core.StageWhere, c.Error)
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
	b.BaseBuilder.dialect = driver2.ResolveDialect(name)
	return b
}

// Build compiles the SELECT statement and returns it as a string and argument list.
// If the FROM clause is missing, an error is returned.
// Dialect rules (quoting, pagination) are applied if configured.
func (b *SelectBuilder) Build() (string, []any, error) {
	if len(b.sources) == 0 {
		b.AddStageError(core.StageFrom, fmt.Errorf("missing source; expected at least one table source"))
	}

	var tokens []string

	// prepare columns
	columns := "*"
	if len(b.columns) > 0 {
		var rendered []string
		for _, column := range b.columns {
			if out := render.Column(b.Dialect, column); out != "" {
				rendered = append(rendered, out)
			}
		}
		columns = strings.Join(rendered, ", ")
	}

	tokens = append([]string{}, fmt.Sprintf("SELECT %s", columns))

	// prepare sources
	if len(b.sources) > 0 {
		var fromParts []string
		for _, tbl := range b.sources {
			if tbl.IsValid() {
				fromParts = append(fromParts, b.Dialect.RenderFrom(tbl.Name, tbl.Alias))
			}
		}
		tokens = append(tokens, "FROM "+strings.Join(fromParts, ", "))
	}

	// prepare conditions
	var args []any
	if len(b.conditions) > 0 {
		binder := bind.NewParamBinderWithPosition(b.Dialect, len(args)+1)
		whereClause, clauseArgs, err := builder.RenderConditionsWithBinder(b.Dialect, b.conditions, binder)
		if err != nil {
			return "", nil, fmt.Errorf("SELECT: %w", err)
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

func (b *SelectBuilder) appendColumns(cols []token.Column) {
	for _, col := range cols {
		if !col.IsValid() {
			b.Validator.AddStageError(core.StageSelect,
				fmt.Errorf("invalid column: %s — %v", col.String(), col.Error))
		}
		b.columns = append(b.columns, col)
	}
}
