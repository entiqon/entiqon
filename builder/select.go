package builder

import (
	"fmt"
	"strings"
)

type SelectBuilder struct {
	columns    []string
	from       string
	conditions []string
	args       []interface{}
	orderBy    []string
	limit      int64
}

func NewSelect() *SelectBuilder {
	return &SelectBuilder{}
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
	query := fmt.Sprintf("SELECT %s FROM %s", sb.columns, sb.from)
	if len(sb.conditions) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(sb.conditions, " AND "))
	}
	if len(sb.orderBy) > 0 {
		query += fmt.Sprintf(" ORDER BY %s", strings.Join(sb.orderBy, " "))
	}
	if sb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", sb.limit)
	}
	return query, sb.args
}
