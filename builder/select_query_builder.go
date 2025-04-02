package builder

import (
	"fmt"
	"strings"
)

type SelectQueryBuilder struct {
	columns    []string
	from       string
	conditions []string
	sorting    []string
	take       *int
	skip       *int
}

// Select sets the columns to retrieve
func (sb *SelectQueryBuilder) Select(columns ...string) *SelectQueryBuilder {
	sb.columns = columns
	return sb
}

// From sets the table to select from
func (sb *SelectQueryBuilder) From(from string) *SelectQueryBuilder {
	sb.from = from
	return sb
}

// Where initializes the WHERE conditions (resets any existing conditions)
func (sb *SelectQueryBuilder) Where(condition string) *SelectQueryBuilder {
	sb.conditions = []string{condition}
	return sb
}

// AndWhere adds an AND condition
func (sb *SelectQueryBuilder) AndWhere(condition string) *SelectQueryBuilder {
	sb.conditions = append(sb.conditions, fmt.Sprintf("AND %s", condition))
	return sb
}

// OrWhere adds an OR condition
func (sb *SelectQueryBuilder) OrWhere(condition string) *SelectQueryBuilder {
	sb.conditions = append(sb.conditions, fmt.Sprintf("OR %s", condition))
	return sb
}

// OrderBy adds an ORDER BY clause
func (sb *SelectQueryBuilder) OrderBy(column string) *SelectQueryBuilder {
	sb.sorting = append(sb.sorting, column)
	return sb
}

// Take sets the LIMIT (engine-agnostic equivalent)
func (sb *SelectQueryBuilder) Take(value int) *SelectQueryBuilder {
	sb.take = &value
	return sb
}

// Skip sets the OFFSET (engine-agnostic equivalent)
func (sb *SelectQueryBuilder) Skip(value int) *SelectQueryBuilder {
	sb.skip = &value
	return sb
}

// Build builds the SQL string
func (sb *SelectQueryBuilder) Build() (string, error) {
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
		tokens = append(tokens, "WHERE "+strings.Join(sb.conditions, " "))
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
