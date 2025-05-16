package builder

import (
	"fmt"
	"strings"
)

// UpsertBuilder builds a SQL UPSERT (INSERT ... ON CONFLICT DO UPDATE) statement.
//
// It composes an internal InsertBuilder and extends it with conflict resolution clauses.
type UpsertBuilder struct {
	// insert handles the INSERT INTO portion of the UPSERT query.
	insert *InsertBuilder
	// conflictColumns are the column(s) that trigger conflict resolution when duplicates occur.
	conflictColumns []string
	// updateSet contains assignments to apply when a conflict is detected.
	// Example: map["name"] = "EXCLUDED.name"
	updateSet map[string]string
}

// NewUpsert returns a new instance of UpsertBuilder.
func NewUpsert() *UpsertBuilder {
	return &UpsertBuilder{
		insert: NewInsert(),
	}
}

// Into sets the target table for the UPSERT statement.
func (b *UpsertBuilder) Into(table string) *UpsertBuilder {
	b.insert.Into(table)
	return b
}

// Columns defines the column names for the UPSERT statement.
func (b *UpsertBuilder) Columns(cols ...string) *UpsertBuilder {
	b.insert.Columns(cols...)
	return b
}

// Values appends a row of values to insert into the table.
func (b *UpsertBuilder) Values(values ...any) *UpsertBuilder {
	b.insert.Values(values...)
	return b
}

// OnConflict defines the column(s) to check for UPSERT conflict handling.
func (b *UpsertBuilder) OnConflict(columns ...string) *UpsertBuilder {
	b.conflictColumns = append(b.conflictColumns, columns...)
	return b
}

// Returning specifies the RETURNING clause for the UPSERT query.
func (b *UpsertBuilder) Returning(columns ...string) *UpsertBuilder {
	b.insert.Returning(columns...)
	return b
}

// DoUpdateSet defines how to update columns if a conflict is found.
//
// For example: map["name"] = "EXCLUDED.name"
func (b *UpsertBuilder) DoUpdateSet(set map[string]string) *UpsertBuilder {
	if b.updateSet == nil {
		b.updateSet = make(map[string]string)
	}
	for col, expr := range set {
		b.updateSet[col] = expr
	}
	return b
}

// Build compiles the UPSERT SQL query with ON CONFLICT and DO UPDATE clauses.
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
