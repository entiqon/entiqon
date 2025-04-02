package builder

import (
	"entiqon/dialect"
	"fmt"
	"strings"
)

type SelectBuilder struct {
	dialect    dialect.Dialect
	columns    []string
	from       string
	conditions []string
	args       []interface{}
	orderBy    []string
	limit      int64
}

func NewSelect(d dialect.Dialect) *SelectBuilder {
	return &SelectBuilder{dialect: d}
}

func (sb *SelectBuilder) Columns(columns ...string) *SelectBuilder {
	sb.columns = columns
	return sb
}

func (sb *SelectBuilder) From(from string) *SelectBuilder {
	sb.from = from
	return sb
}

func (sb *SelectBuilder) Where(conditions string, args ...interface{}) *SelectBuilder {
	sb.conditions = append(sb.conditions, conditions)
	sb.args = append(sb.args, args...)
	return sb
}

func (sb *SelectBuilder) OrderBy(orderBy ...string) *SelectBuilder {
	sb.orderBy = append(sb.orderBy, orderBy...)
	return sb
}

func (sb *SelectBuilder) Limit(limit int64) *SelectBuilder {
	sb.limit = limit
	return sb
}

func (sb *SelectBuilder) Build() (string, []interface{}) {
	binder := NewParamBinder(sb.dialect)

	escapedCols := make([]string, len(sb.columns))
	for i, col := range sb.columns {
		escapedCols[i] = sb.dialect.EscapeIdentifier(col)
	}

	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(escapedCols, ", "), sb.dialect.EscapeIdentifier(sb.from))
	if len(sb.conditions) > 0 {
		// Replace all placeholders with dialect-specific ones
		processed := make([]string, len(sb.conditions))
		argIndex := 0
		for i, condition := range sb.conditions {
			// Replace `?` placeholders manually with dialect placeholders
			count := strings.Count(condition, "?")
			bound := condition
			for j := 0; j < count; j++ {
				ph := binder.Bind(binder.args[argIndex])
				bound = strings.Replace(bound, "?", ph, 1)
				argIndex++
			}
			processed[i] = bound
		}
		query += fmt.Sprintf(" WHERE %s", strings.Join(processed, " AND "))
	}
	if len(sb.orderBy) > 0 {
		query += fmt.Sprintf(" ORDER BY %s", strings.Join(sb.orderBy, " "))
	}
	if sb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", sb.limit)
	}
	return query, binder.Args()
}
