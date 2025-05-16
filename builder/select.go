package builder

import (
	"fmt"
	"strings"
)

// ConditionType defines how a condition is logically connected in a WHERE clause.
type ConditionType string

const (
	// ConditionSimple is used for initial WHERE conditions.
	ConditionSimple ConditionType = "SIMPLE"

	// ConditionAnd adds an AND between conditions.
	ConditionAnd ConditionType = "AND"

	// ConditionOr adds an OR between conditions.
	ConditionOr ConditionType = "OR"
)

// ConditionToken represents a conditional expression used in a WHERE clause.
type ConditionToken struct {
	// Type specifies how this condition is logically joined (SIMPLE, AND, OR).
	Type ConditionType

	// Condition is the SQL condition string (e.g., "id = 1").
	Condition string
}

// SelectBuilder builds a SQL SELECT query using fluent method chaining.
//
// It supports basic querying with WHERE conditions, ordering, and pagination.
type SelectBuilder struct {
	columns    []string
	from       string
	conditions []ConditionToken
	sorting    []string
	take       *int
	skip       *int
}

// NewSelect creates and returns a new instance of SelectBuilder.
func NewSelect() *SelectBuilder {
	return &SelectBuilder{
		columns:    make([]string, 0),
		conditions: make([]ConditionToken, 0),
		sorting:    make([]string, 0),
	}
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
func (sb *SelectBuilder) Where(conditions ...string) *SelectBuilder {
	sb.conditions = []ConditionToken{}
	sb.addCondition(ConditionSimple, conditions...)
	return sb
}

// AndWhere adds an AND condition to the WHERE clause.
func (sb *SelectBuilder) AndWhere(conditions ...string) *SelectBuilder {
	sb.addCondition(ConditionAnd, conditions...)
	return sb
}

// OrWhere adds an OR condition to the WHERE clause.
func (sb *SelectBuilder) OrWhere(conditions ...string) *SelectBuilder {
	sb.addCondition(ConditionOr, conditions...)
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
func (sb *SelectBuilder) Build() (string, error) {
	if sb.from == "" {
		return "", fmt.Errorf("FROM clause is required")
	}

	columns := "*"
	if len(sb.columns) > 0 {
		columns = strings.Join(sb.columns, ", ")
	}

	tokens := []string{
		fmt.Sprintf("SELECT %s", columns),
		fmt.Sprintf("FROM %s", sb.from),
	}

	if len(sb.conditions) > 0 {
		var parts []string
		for _, condition := range sb.conditions {
			switch condition.Type {
			case ConditionSimple:
				parts = append(parts, condition.Condition)
			case ConditionAnd, ConditionOr:
				parts = append(parts, fmt.Sprintf("%s %s", condition.Type, condition.Condition))
			default:
				return "", fmt.Errorf("invalid condition type: %s", condition.Type)
			}
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

	return strings.Join(tokens, " "), nil
}

// addCondition adds a logical condition to the WHERE clause.
func (sb *SelectBuilder) addCondition(conditionType ConditionType, conditions ...string) {
	if len(conditions) == 0 {
		return
	}
	if len(conditions) == 1 {
		sb.conditions = append(sb.conditions, ConditionToken{
			Type:      conditionType,
			Condition: conditions[0],
		})
		return
	}
	// Join and wrap in parentheses
	joiner := string(conditionType)
	group := "(" + strings.Join(conditions, fmt.Sprintf(" %s ", joiner)) + ")"
	sb.conditions = append(sb.conditions, ConditionToken{
		Type:      conditionType,
		Condition: group,
	})
}
