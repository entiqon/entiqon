// File: db/builder/select.go

// Package builder provides a fluent API to construct SQL SELECT queries
// with support for expressions, aliases, pagination, and dialect-aware
// quoting and syntax. It leverages token structures and dialect implementations
// to generate safe, database-specific SQL statements.
//
// SelectBuilder is the primary type for building SELECT queries, allowing
// incremental addition of fields, table source, limits, and offsets.
//
// This package is designed to be extensible and easy to integrate with
// various SQL dialects via the dialect package.
package builder

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/common/extension/collection"
	"github.com/entiqon/entiqon/db/dialect"
	"github.com/entiqon/entiqon/db/token"
)

// SelectBuilder builds simple SELECT queries.
type SelectBuilder struct {
	dialect dialect.Dialect
	fields  *collection.Collection[token.Field]
	table   string
	limit   int
	offset  int
}

// NewSelect creates a new SelectBuilder with the provided dialect.
// If nil is passed, BaseDialect is used by default.
func NewSelect(d dialect.Dialect) *SelectBuilder {
	if d == nil {
		d = &dialect.BaseDialect{}
	}
	return &SelectBuilder{
		dialect: d,
		fields:  collection.New[token.Field](),
	}
}

// Fields adds fields to the select clause.
// Fields initializes the collection of fields, replacing any existing ones.
//
// Supports input arguments of type string or token.Field.
// For strings:
// - comma-separated list is parsed into multiple fields
// - parses alias with AS or space
func (b *SelectBuilder) Fields(cols ...interface{}) *SelectBuilder {
	if b.fields == nil {
		b.fields = collection.New[token.Field]()
	} else {
		b.fields.Clear()
	}
	return b.appendFields(cols...)
}

// AddFields appends to the collection, preserving existing fields.
// If the collection is empty, it initializes it like Fields would.
func (b *SelectBuilder) AddFields(cols ...interface{}) *SelectBuilder {
	return b.appendFields(cols...)
}

// GetFields returns the current list of fields in the builder.
// It returns them by value to avoid external modification of the internal slice.
func (b *SelectBuilder) GetFields() token.FieldCollection {
	return b.fields.Items()
}

// Source sets the table name.
func (b *SelectBuilder) Source(table string) *SelectBuilder {
	b.table = table
	return b
}

// Limit sets the LIMIT clause.
func (b *SelectBuilder) Limit(limit int) *SelectBuilder {
	b.limit = limit
	return b
}

// Offset sets the OFFSET clause.
func (b *SelectBuilder) Offset(offset int) *SelectBuilder {
	b.offset = offset
	return b
}

// String implements fmt.Stringer for convenient status output.
func (b *SelectBuilder) String() string {
	sql, err := b.Build()
	if err != nil {
		return fmt.Sprintf("Status ❌: Error building SQL: %v", err)
	}

	// TODO: update with actual parameters once supported
	//params := []interface{}{}
	//
	//paramsStr := ""
	//if len(params) > 0 {
	//	paramsStr = fmt.Sprintf("%v", params)
	//}

	//return fmt.Sprintf("Status ✅: SQL=%s, Params=%s", sql /*, paramsStr*/)
	return fmt.Sprintf("Status ✅: SQL=%s, Params=%s", sql, "")
}

// Build constructs the SQL query string.
func (b *SelectBuilder) Build() (string, error) {
	if b == nil {
		return "", fmt.Errorf("❌ [Build] - Wrong initialization. Cannot build on receiver nil")
	}

	if b.fields.Length() == 0 {
		b.fields.Add(token.Field{
			Expr:  "*",
			IsRaw: true,
		})
	}

	parts := make([]string, 0, b.fields.Length())
	var bad []string

	for _, field := range b.fields.Items() {
		if field.IsErrored() {
			bad = append(bad,
				fmt.Sprintf("⛔️ Field(%q): %v", field.Input, field.Error))
			continue
		}
		parts = append(parts, field.Render())
	}

	if len(bad) > 0 {
		return "", fmt.Errorf("❌ [Build] - Invalid fields:\n\t%s", strings.Join(bad, "\n\t"))
	}

	if b.table == "" {
		return "", fmt.Errorf("❌ [Build] - No source specified")
	}

	tokens := []string{
		"SELECT",
		strings.Join(parts, ", "),
		"FROM",
		b.table,
	}

	sql := strings.Join(tokens, " ")

	if b.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", b.limit)
	}
	if b.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", b.offset)
	}

	return sql, nil
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// appendFields contains the shared logic for parsing/adding fields.
func (b *SelectBuilder) appendFields(cols ...interface{}) *SelectBuilder {
	switch len(cols) {
	case 0:
		return b
	case 1:
		switch v := cols[0].(type) {
		case string:
			if strings.TrimSpace(v) == "" {
				return b
			}
			if strings.Contains(v, ",") {
				parts := splitAndTrim(v, ",")
				for _, part := range parts {
					b.fields.Add(*token.NewField(part))
				}
			} else {
				b.fields.Add(*token.NewField(v))
			}
		case token.Field:
			b.fields.Add(v)
		case *token.Field:
			if v != nil {
				b.fields.Add(*v)
			}
		default:
			b.fields.Add(*token.NewField(v))
		}
	case 2, 3:
		b.fields.Add(*token.NewField(cols...))
	}
	return b
}
