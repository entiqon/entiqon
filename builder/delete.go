package builder

import (
	"fmt"
	"strings"
)

type DeleteBuilder struct {
	from      string
	where     []string
	args      []any
	returning []string
}

func NewDelete() *DeleteBuilder {
	return &DeleteBuilder{
		where:     make([]string, 0),
		args:      make([]any, 0),
		returning: make([]string, 0),
	}
}

func (b *DeleteBuilder) From(table string) *DeleteBuilder {
	b.from = table
	return b
}

func (b *DeleteBuilder) Where(condition string, args ...any) *DeleteBuilder {
	b.where = append(b.where, condition)
	b.args = append(b.args, args...)
	return b
}

func (b *DeleteBuilder) Returning(columns ...string) *DeleteBuilder {
	b.returning = append(b.returning, columns...)
	return b
}

// Build compiles the DELETE SQL query and returns the statement and argument list.
func (b *DeleteBuilder) Build() (string, []any, error) {
	if b.from == "" {
		return "", nil, fmt.Errorf("no FROM table specified")
	}

	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	sb.WriteString(b.from)

	if len(b.where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.where, " AND "))
	}

	if len(b.returning) > 0 {
		sb.WriteString(" RETURNING ")
		sb.WriteString(strings.Join(b.returning, ", "))
	}

	return sb.String(), b.args, nil
}
