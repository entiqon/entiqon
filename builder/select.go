package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/builder"
	"github.com/ialopezg/entiqon/internal/core/builder/bind"
	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// SelectBuilder builds a SQL SELECT query using fluent method chaining.
//
// It supports basic querying with WHERE conditions, ordering, and pagination.
type SelectBuilder struct {
	BaseBuilder
	columns    []token.FieldToken
	from       string
	conditions []token.Condition
	sorting    []string
	take       *int
	skip       *int
}

// NewSelect creates and returns a new instance of SelectBuilder.
func NewSelect() *SelectBuilder {
	return &SelectBuilder{
		BaseBuilder: BaseBuilder{dialect: driver.NewGenericDialect()},
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
func (b *SelectBuilder) From(table string) *SelectBuilder {
	if table == "" {
		b.AddStageError("FROM", fmt.Errorf("table is empty"))
	} else {
		b.from = table
	}
	return b
}

// Where sets the base condition(s) for the WHERE clause.
// It resets any previously added conditions.
func (b *SelectBuilder) Where(condition string, values ...any) *SelectBuilder {
	c := token.NewCondition(token.ConditionSimple, condition, values...)
	if !c.IsValid() {
		b.AddStageError("WHERE", c.Error)
	}
	b.conditions = []token.Condition{c}
	return b
}

// AndWhere adds an AND condition to the WHERE clause.
func (b *SelectBuilder) AndWhere(condition string, values ...any) *SelectBuilder {
	c := token.NewCondition(token.ConditionAnd, condition, values...)
	if !c.IsValid() {
		b.AddStageError("WHERE", c.Error)
	}
	b.conditions = append(b.conditions, c)
	return b
}

// OrWhere adds an OR condition to the WHERE clause.
func (b *SelectBuilder) OrWhere(condition string, values ...any) *SelectBuilder {
	c := token.NewCondition(token.ConditionOr, condition, values...)
	if !c.IsValid() {
		b.AddStageError("WHERE", c.Error)
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
	dialect := b.GetDialect()

	if b.HasErrors() {
		return "", nil, fmt.Errorf("FROM: %d invalid condition(s)", len(b.GetErrors()))
	}

	if b.from == "" {
		return "", nil, fmt.Errorf("FROM: requires a target table")
	}

	// ─────────────────────────────────────────────────────────────
	// Render SELECT columns
	// ─────────────────────────────────────────────────────────────
	columns := "*"
	if len(b.columns) > 0 {
		var rendered []string
		for _, col := range b.columns {
			name := col.Name
			if b.dialect != nil && !col.IsRaw {
				name = dialect.QuoteIdentifier(col.Name)
			}
			if col.Alias != "" {
				name = fmt.Sprintf("%s AS %s", name, col.Alias)
			}
			rendered = append(rendered, name)
		}
		columns = strings.Join(rendered, ", ")
	}

	// ─────────────────────────────────────────────────────────────
	// Render FROM clause (quoted if dialect provided)
	// ─────────────────────────────────────────────────────────────
	from := b.from
	if b.dialect != nil {
		from = b.dialect.QuoteIdentifier(from)
	}

	tokens := []string{
		fmt.Sprintf("SELECT %s", columns),
		fmt.Sprintf("FROM %s", from),
	}

	// ─────────────────────────────────────────────────────────────
	// Render WHERE conditions
	// ─────────────────────────────────────────────────────────────
	var args []any
	if len(b.conditions) > 0 {
		binder := bind.NewParamBinderWithPosition(dialect, len(args)+1)
		whereClause, clauseArgs, err := builder.RenderConditionsWithBinder(dialect, b.conditions, binder)
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
		tokens = append(tokens, dialect.BuildLimitOffset(limit, offset))
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
