package builder

import (
	"fmt"
	"strings"
)

// DeleteBuilder builds a SQL DELETE statement.
//
// It supports WHERE clauses and optional RETURNING fields (e.g., for PostgreSQL).
type DeleteBuilder struct {
	// from defines the target table to delete from.
	from string

	// where holds WHERE clause expressions (joined with AND).
	where []string

	// args holds placeholder arguments for WHERE conditions.
	args []any

	// returning lists the columns to return (PostgreSQL only).
	returning []string
}

// NewDelete returns a new DeleteBuilder instance.
func NewDelete() *DeleteBuilder {
	return &DeleteBuilder{
		where:     make([]string, 0),
		args:      make([]any, 0),
		returning: make([]string, 0),
	}
}

// From sets the table to delete from.
func (b *DeleteBuilder) From(table string) *DeleteBuilder {
	b.from = table
	return b
}

// Where adds a WHERE clause with optional placeholder arguments.
//
// Example: .Where("id = ?", 42)
func (b *DeleteBuilder) Where(condition string, args ...any) *DeleteBuilder {
	b.where = append(b.where, condition)
	b.args = append(b.args, args...)
	return b
}

// Returning specifies which columns to return after deletion (PostgreSQL only).
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
