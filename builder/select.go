package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// SelectBuilder builds a SQL SELECT query using fluent method chaining.
//
// It supports basic querying with WHERE conditions, ordering, and pagination.
type SelectBuilder struct {
	dialect    driver.Dialect
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
		columns:    make([]token.FieldToken, 0),
		conditions: make([]token.Condition, 0),
		sorting:    make([]string, 0),
	}
}

// Select adds raw column strings (can include aliases like "id", "name AS n").
// Select sets columns using FieldsFromExpr(...) and resets previous entries.
func (sb *SelectBuilder) Select(columns ...string) *SelectBuilder {
	sb.columns = []token.FieldToken{}
	for _, expr := range columns {
		sb.columns = append(sb.columns, token.FieldsFromExpr(expr)...)
	}
	return sb
}

// AddSelect appends more columns using FieldsFromExpr(...) without resetting.
func (sb *SelectBuilder) AddSelect(columns ...string) *SelectBuilder {
	for _, expr := range columns {
		sb.columns = append(sb.columns, token.FieldsFromExpr(expr)...)
	}
	return sb
}

// From sets the target table for the SELECT statement.
func (sb *SelectBuilder) From(from string) *SelectBuilder {
	sb.from = from
	return sb
}

// Where sets the base condition(s) for the WHERE clause.
// It resets any previously added conditions.
func (sb *SelectBuilder) Where(condition string, params ...any) *SelectBuilder {
	sb.conditions = token.AppendCondition(
		[]token.Condition{},
		token.NewCondition(token.ConditionSimple, condition, params...),
	)
	return sb
}

// AndWhere adds an AND condition to the WHERE clause.
func (sb *SelectBuilder) AndWhere(condition string, params ...any) *SelectBuilder {
	sb.conditions = token.AppendCondition(sb.conditions, token.NewCondition(token.ConditionAnd, condition, params...))
	return sb
}

// OrWhere adds an OR condition to the WHERE clause.
func (sb *SelectBuilder) OrWhere(condition string, params ...any) *SelectBuilder {
	sb.conditions = token.AppendCondition(sb.conditions, token.NewCondition(token.ConditionOr, condition, params...))
	return sb
}

// OrderBy appends a column or expression to the ORDER BY clause.
func (sb *SelectBuilder) OrderBy(column string) *SelectBuilder {
	sb.sorting = append(sb.sorting, column)
	return sb
}

// Take limits the number of rows returned by the query (engine-agnostic equivalent).
func (sb *SelectBuilder) Take(value int) *SelectBuilder {
	sb.take = &value
	return sb
}

// Skip offsets the rows returned by the query (engine-agnostic equivalent).
func (sb *SelectBuilder) Skip(value int) *SelectBuilder {
	sb.skip = &value
	return sb
}

// UseDialect resolves and applies the dialect.md by name (e.g., "postgres").
// It replaces any previously set dialect.md on the builder.
func (sb *SelectBuilder) UseDialect(name string) *SelectBuilder {
	sb.dialect = driver.ResolveDialect(name)
	return sb
}

// WithDialect sets the dialect.md engine used for escaping and quoting.
//
// Deprecated: use UseDialect(name string) instead for improved clarity
// and consistency with the fluent builder style.
func (sb *SelectBuilder) WithDialect(name string) *SelectBuilder {
	sb.dialect = driver.ResolveDialect(name)
	return sb
}

// Build compiles the SELECT statement and returns it as a string and argument list.
// If the FROM clause is missing, an error is returned.
// Dialect rules (quoting, pagination) are applied if configured.
func (sb *SelectBuilder) Build() (string, []any, error) {
	if sb.from == "" {
		return "", nil, fmt.Errorf("FROM clause is required")
	}

	// ─────────────────────────────────────────────────────────────
	// Render SELECT columns
	// ─────────────────────────────────────────────────────────────
	columns := "*"
	if len(sb.columns) > 0 {
		var rendered []string
		for _, col := range sb.columns {
			name := col.Name
			if sb.dialect != nil && !col.IsRaw {
				name = sb.dialect.Quote(col.Name)
			}
			if col.Alias != "" {
				name = fmt.Sprintf("%s AS %s", name, col.Alias)
			}
			rendered = append(rendered, name)
		}
		columns = strings.Join(rendered, ", ")
	}

	// ─────────────────────────────────────────────────────────────
	// Render FROM clause (quoted if dialect.md provided)
	// ─────────────────────────────────────────────────────────────
	from := sb.from
	if sb.dialect != nil {
		from = sb.dialect.Quote(from)
	}

	tokens := []string{
		fmt.Sprintf("SELECT %s", columns),
		fmt.Sprintf("FROM %s", from),
	}

	// ─────────────────────────────────────────────────────────────
	// Render WHERE conditions
	// ─────────────────────────────────────────────────────────────
	if len(sb.conditions) > 0 {
		var parts []string
		for _, condition := range sb.conditions {
			rendered := condition.Key

			if sb.dialect != nil {
				if parsed := strings.SplitN(condition.Key, "=", 2); len(parsed) == 2 {
					field := strings.TrimSpace(parsed[0])
					right := strings.TrimSpace(parsed[1])
					rendered = fmt.Sprintf("%s = %s", sb.dialect.Quote(field), right)
				}
			}

			switch condition.Type {
			case token.ConditionSimple:
				parts = append(parts, rendered)
			case token.ConditionAnd, token.ConditionOr:
				parts = append(parts, fmt.Sprintf("%s %s", condition.Type, rendered))
			default:
				return "", nil, fmt.Errorf("invalid condition type: %s", condition.Type)
			}
		}
		tokens = append(tokens, fmt.Sprintf("WHERE %s", strings.Join(parts, " ")))
	}

	// ─────────────────────────────────────────────────────────────
	// ORDER BY
	// ─────────────────────────────────────────────────────────────
	if len(sb.sorting) > 0 {
		tokens = append(tokens, "ORDER BY "+strings.Join(sb.sorting, ", "))
	}

	// ─────────────────────────────────────────────────────────────
	// LIMIT/OFFSET via dialect.md
	// ─────────────────────────────────────────────────────────────
	limit, offset := -1, -1
	if sb.take != nil {
		limit = *sb.take
	}
	if sb.skip != nil {
		offset = *sb.skip
	}

	if sb.dialect != nil && (limit >= 0 || offset >= 0) {
		tokens = append(tokens, sb.dialect.BuildLimitOffset(limit, offset))
	} else {
		if limit >= 0 {
			tokens = append(tokens, fmt.Sprintf("LIMIT %d", limit))
		}
		if offset >= 0 {
			tokens = append(tokens, fmt.Sprintf("OFFSET %d", offset))
		}
	}

	return strings.Join(tokens, " "), sb.collectArgs(), nil
}

// collectArgs gathers all condition parameters.
func (sb *SelectBuilder) collectArgs() []any {
	var args []any
	for _, c := range sb.conditions {
		args = append(args, c.Params...)
	}
	return args
}
