package builder

import (
	"fmt"
	"strings"
)

type UpdateBuilder struct {
	table string
	set   map[string]any
	where []string
	args  []any
}

// NewUpdate returns a new UpdateBuilder instance.
func NewUpdate() *UpdateBuilder {
	return &UpdateBuilder{
		set:   make(map[string]any),
		where: make([]string, 0),
		args:  make([]any, 0),
	}
}

// Table sets the table to be updated.
func (b *UpdateBuilder) Table(name string) *UpdateBuilder {
	b.table = name
	return b
}

// Set defines a field and value to update.
func (b *UpdateBuilder) Set(field string, value any) *UpdateBuilder {
	b.set[field] = value
	return b
}

// Where adds a WHERE clause with arguments.
func (b *UpdateBuilder) Where(condition string, args ...any) *UpdateBuilder {
	b.where = append(b.where, condition)
	b.args = append(b.args, args...)
	return b
}

// Build compiles the UPDATE statement and returns the SQL string and arguments.
func (b *UpdateBuilder) Build() (string, []any, error) {
	if b.table == "" {
		return "", nil, fmt.Errorf("no table specified")
	}
	if len(b.set) == 0 {
		return "", nil, fmt.Errorf("no fields specified in SET clause")
	}

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(b.table)
	sb.WriteString(" SET ")

	setClauses := make([]string, 0, len(b.set))
	args := make([]any, 0, len(b.set)+len(b.args))

	for col, val := range b.set {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}

	sb.WriteString(strings.Join(setClauses, ", "))

	if len(b.where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.where, " AND "))
		args = append(args, b.args...)
	}

	return sb.String(), args, nil
}
