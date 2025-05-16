package builder

import (
	"fmt"
	"strings"
)

// UpsertBuilder builds a SQL UPSERT (INSERT ... ON CONFLICT DO UPDATE) statement.
//
// It will eventually combine InsertBuilder and UpdateBuilder logic internally.
type UpsertBuilder struct {
	insert          *InsertBuilder
	conflictColumns []string
	updateSet       map[string]string
}

func NewUpsert() *UpsertBuilder {
	return &UpsertBuilder{
		insert: NewInsert(),
	}
}

func (b *UpsertBuilder) Into(table string) *UpsertBuilder {
	b.insert.Into(table)
	return b
}

func (b *UpsertBuilder) Columns(cols ...string) *UpsertBuilder {
	b.insert.Columns(cols...)
	return b
}

func (b *UpsertBuilder) Values(values ...any) *UpsertBuilder {
	b.insert.Values(values...)
	return b
}

func (b *UpsertBuilder) OnConflict(columns ...string) *UpsertBuilder {
	b.conflictColumns = append(b.conflictColumns, columns...)
	return b
}

func (b *UpsertBuilder) DoUpdateSet(set map[string]string) *UpsertBuilder {
	if b.updateSet == nil {
		b.updateSet = make(map[string]string)
	}
	for col, expr := range set {
		b.updateSet[col] = expr
	}
	return b
}

// Build compiles the UPSERT SQL query (PostgreSQL style) with ON CONFLICT clause.
func (b *UpsertBuilder) Build() (string, []any, error) {
	sql, args, err := b.insert.Build()
	if err != nil {
		return "", nil, err
	}

	if len(b.conflictColumns) == 0 {
		return "", nil, fmt.Errorf("ON CONFLICT requires at least one column")
	}
	if len(b.updateSet) == 0 {
		return "", nil, fmt.Errorf("DO UPDATE SET requires at least one assignment")
	}

	var sb strings.Builder
	sb.WriteString(sql)

	sb.WriteString(" ON CONFLICT (")
	sb.WriteString(strings.Join(b.conflictColumns, ", "))
	sb.WriteString(") DO UPDATE SET ")

	assignments := make([]string, 0, len(b.updateSet))
	for col, expr := range b.updateSet {
		assignments = append(assignments, fmt.Sprintf("%s = %s", col, expr))
	}
	sb.WriteString(strings.Join(assignments, ", "))

	return sb.String(), args, nil
}
