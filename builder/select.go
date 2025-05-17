package builder

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/core/dialect"
	"github.com/ialopezg/entiqon/internal/core/token"
)

// SelectBuilder builds a SQL SELECT query using fluent method chaining.
//
// It supports basic querying with WHERE conditions, ordering, and pagination.
type SelectBuilder struct {
	dialect    dialect.Engine
	columns    []string
	from       string
	conditions []token.Condition
	sorting    []string
	take       *int
	skip       *int
}

// NewSelect creates and returns a new instance of SelectBuilder.
func NewSelect() *SelectBuilder {
	return &SelectBuilder{
		columns:    make([]string, 0),
		conditions: make([]token.Condition, 0),
		sorting:    make([]string, 0),
	}
}

// WithDialect sets the dialect engine used for escaping.
func (sb *SelectBuilder) WithDialect(d dialect.Engine) *SelectBuilder {
	sb.dialect = d
	return sb
}

// Select specifies the columns to retrieve.
func (sb *SelectBuilder) Select(columns ...string) *SelectBuilder {
	sb.columns = columns
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
	sb.conditions = []token.Condition{}
	sb.addCondition(token.ConditionSimple, condition, params...)
	return sb
}

// AndWhere adds an AND condition to the WHERE clause.
func (sb *SelectBuilder) AndWhere(condition string, params ...any) *SelectBuilder {
	sb.addCondition(token.ConditionAnd, condition, params...)
	return sb
}

// OrWhere adds an OR condition to the WHERE clause.
func (sb *SelectBuilder) OrWhere(condition string, params ...any) *SelectBuilder {
	sb.addCondition(token.ConditionOr, condition, params...)
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

// Build compiles the SELECT statement and returns it as a string.
// It returns an error if essential parts (like the FROM clause) are missing.
func (sb *SelectBuilder) Build() (string, []any, error) {
	if sb.from == "" {
		return "", nil, fmt.Errorf("FROM clause is required")
	}

	from := sb.from
	if sb.dialect != nil {
		from = sb.dialect.EscapeIdentifier(sb.from)
	}

	columns := "*"
	if len(sb.columns) > 0 {
		if sb.dialect != nil {
			for i, column := range sb.columns {
				sb.columns[i] = sb.dialect.EscapeIdentifier(column)
			}
		}
		columns = strings.Join(sb.columns, ", ")
	}

	tokens := []string{
		fmt.Sprintf("SELECT %s", columns),
		fmt.Sprintf("FROM %s", from),
	}

	var args []any
	if len(sb.conditions) > 0 {
		var parts []string
		for _, condition := range sb.conditions {
			switch condition.Type {
			case token.ConditionSimple:
				parts = append(parts, condition.Key)
			case token.ConditionAnd, token.ConditionOr:
				parts = append(parts, fmt.Sprintf("%s %s", condition.Type, condition.Key))
			default:
				return "", nil, fmt.Errorf("invalid condition type: %s", condition.Type)
			}
			args = append(args, condition.Params...)
		}
		tokens = append(tokens, fmt.Sprintf("WHERE %s", strings.Join(parts, " ")))
	}

	if len(sb.sorting) > 0 {
		tokens = append(tokens, "ORDER BY "+strings.Join(sb.sorting, ", "))
	}

	if sb.take != nil {
		tokens = append(tokens, fmt.Sprintf("LIMIT %d", *sb.take))
	}
	if sb.skip != nil {
		tokens = append(tokens, fmt.Sprintf("OFFSET %d", *sb.skip))
	}

	return strings.Join(tokens, " "), args, nil
}

// addCondition adds a logical condition to the WHERE clause.
func (sb *SelectBuilder) addCondition(conditionType token.ConditionType, condition string, params ...any) {
	if condition == "" {
		return
	}

	raw := condition
	for _, val := range params {
		raw = fmt.Sprintf("(%s)", strings.Replace(raw, "?", fmt.Sprintf("'%v'", val), 1))
	}

	sb.conditions = append(sb.conditions, token.Condition{
		Type:   conditionType,
		Key:    condition,
		Params: params,
		Raw:    raw,
	})
}
