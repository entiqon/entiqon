package builder

import (
	"fmt"
	"strings"

	driver2 "github.com/ialopezg/entiqon/driver"
	"github.com/ialopezg/entiqon/internal/core/builder"
	"github.com/ialopezg/entiqon/internal/core/builder/bind"
	core "github.com/ialopezg/entiqon/internal/core/error"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// SelectBuilder builds a SQL SELECT query using fluent method chaining.
//
// It supports basic querying with WHERE conditions, ordering, and pagination.
type SelectBuilder struct {
	BaseBuilder
	columns    []token.FieldToken
	from       []token.Table
	conditions []token.Condition
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
		columns:     make([]token.FieldToken, 0),
		conditions:  make([]token.Condition, 0),
		sorting:     make([]string, 0),
	}
}

// Select adds raw column strings (can include aliases like "id", "name AS n").
// Select sets columns using FieldsFromExpr(...) and resets previous entries.
func (b *SelectBuilder) Select(columns ...string) *SelectBuilder {
	b.columns = []token.FieldToken{}
	for _, expr := range columns {
		b.columns = append(b.columns, token.FieldsFromExpr(expr)...)
	}
	return b
}

// AddSelect appends more columns using FieldsFromExpr(...) without resetting.
func (b *SelectBuilder) AddSelect(columns ...string) *SelectBuilder {
	for _, expr := range columns {
		b.columns = append(b.columns, token.FieldsFromExpr(expr)...)
	}
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
		b.from = append(b.from, token.NewTableWithAlias(table, alias[0]))
	} else {
		b.from = append(b.from, token.NewTable(table))
	}
	return b
}

// Where sets the base condition(s) for the WHERE clause.
// It resets any previously added conditions.
func (b *SelectBuilder) Where(condition string, values ...any) *SelectBuilder {
	c := token.NewCondition(token.ConditionSimple, condition, values...)
	if !c.IsValid() {
		b.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = []token.Condition{c}
	return b
}

// AndWhere adds an AND condition to the WHERE clause.
func (b *SelectBuilder) AndWhere(condition string, values ...any) *SelectBuilder {
	c := token.NewCondition(token.ConditionAnd, condition, values...)
	if !c.IsValid() {
		b.AddStageError(core.StageWhere, c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// OrWhere adds an OR condition to the WHERE clause.
func (b *SelectBuilder) OrWhere(condition string, values ...any) *SelectBuilder {
	c := token.NewCondition(token.ConditionOr, condition, values...)
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
	if len(b.from) == 0 {
		b.Validator.AddStageError(core.StageFrom, fmt.Errorf("requires a target table"))
	}

	if err := b.Validate(); err != nil {
		return "", nil, err
	}

	var tokens []string
	// ─────────────────────────────────────────────────────────────
	// Render SELECT columns
	// ─────────────────────────────────────────────────────────────
	columns := "*"
	if len(b.columns) > 0 {
		var rendered []string
		for _, col := range b.columns {
			name := col.Name
			if b.Dialect != nil && !col.IsRaw {
				name = b.Dialect.QuoteIdentifier(col.Name)
			}
			if col.Alias != "" {
				name = fmt.Sprintf("%s AS %s", name, col.Alias)
			}
			rendered = append(rendered, name)
		}
		columns = strings.Join(rendered, ", ")
	}

	tokens = append([]string{}, fmt.Sprintf("SELECT %s", columns))

	// ─────────────────────────────────────────────────────────────
	// Render FROM clause (quoted if dialect provided)
	// ─────────────────────────────────────────────────────────────
	if len(b.from) > 0 {
		var fromParts []string
		for _, tbl := range b.from {
			if tbl.IsValid() {
				fromParts = append(fromParts, b.Dialect.RenderFrom(tbl.Name, tbl.Alias))
			}
		}
		tokens = append(tokens, "FROM "+strings.Join(fromParts, ", "))
	}

	// ─────────────────────────────────────────────────────────────
	// Render WHERE conditions
	// ─────────────────────────────────────────────────────────────
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

	// ─────────────────────────────────────────────────────────────
	// ORDER BY
	// ─────────────────────────────────────────────────────────────
	if len(b.sorting) > 0 {
		tokens = append(tokens, "ORDER BY "+strings.Join(b.sorting, ", "))
	}

	// ─────────────────────────────────────────────────────────────
	// LIMIT/OFFSET via dialect
	// ─────────────────────────────────────────────────────────────
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

	return strings.Join(tokens, " "), args, nil
}
