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

	"github.com/entiqon/entiqon/db/dialect"
	"github.com/entiqon/entiqon/db/token"
)

// SelectBuilder builds simple SELECT queries.
type SelectBuilder struct {
	dialect dialect.Dialect
	fields  token.FieldCollection
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
		fields:  token.FieldCollection{},
	}
}

// Fields adds fields to the select clause.
// Supports input arguments of type string or token.Field.
// For strings:
// - comma-separated list is parsed into multiple fields
// - parses alias with AS or space
func (b *SelectBuilder) Fields(cols ...interface{}) *SelectBuilder {
	switch len(cols) {
	case 1:
		switch v := cols[0].(type) {
		case string:
			// Ignore empty strings — no field should be added for blank input.
			if strings.TrimSpace(v) == "" {
				return b
			}

			if strings.Contains(v, ",") {
				// Comma-separated list → split and add each
				parts := splitAndTrim(v, ",")
				for _, part := range parts {
					b.fields.Add(*token.NewField(part))
				}
			} else {
				// Single field or inline alias
				b.fields.Add(*token.NewField(v))
			}
		case token.Field:
			// Append Field directly
			b.fields = append(b.fields, v)

		case *token.Field:
			// Append dereferenced pointer
			if v != nil {
				b.fields = append(b.fields, *v)
			}
		default:
			// Either ignore silently or record as an errored field
			b.fields.Add(*token.NewField(v))
		}
		return b
	case 2, 3:
		b.fields.Add(*token.NewField(cols...))
		return b
	}

	return b
}

// GetFields returns the current list of fields in the builder.
// It returns them by value to avoid external modification of the internal slice.
func (b *SelectBuilder) GetFields() token.FieldCollection {
	return b.fields
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

// Build constructs the SQL query string.
func (b *SelectBuilder) Build() (string, error) {
	if b.fields.Length() == 0 {
		b.fields.Add(token.Field{
			Expr:  "*",
			IsRaw: true,
		})
	}
	if b.table == "" {
		return "", fmt.Errorf("no table specified")
	}

	parts := make([]string, 0, b.fields.Length())
	for _, field := range b.fields {
		parts = append(parts, field.Render())
	}

	tokens := []string{
		"SELECT",
		strings.Join(parts, ", "),
		"FROM",
		b.dialect.QuoteIdentifier(b.table),
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

// isRaw detects whether the expression contains raw SQL indicators.
func isRaw(expr string) bool {
	if strings.ContainsAny(expr, "()") {
		return true
	}

	operators := []string{"||", "+", "-", "*", "/"}
	for _, op := range operators {
		if strings.Contains(expr, op) {
			return true
		}
	}

	return false
}

// parseInput extracts expression and optional alias from the input string.
func parseInput(s string) (expr, alias string) {
	upper := strings.ToUpper(s)
	if idx := strings.Index(upper, " AS "); idx >= 0 {
		expr = strings.TrimSpace(s[:idx])
		alias = strings.TrimSpace(s[idx+4:])
		return
	}

	lastSpace := strings.LastIndex(s, " ")
	if lastSpace > 0 {
		expr = strings.TrimSpace(s[:lastSpace])
		alias = strings.TrimSpace(s[lastSpace+1:])
		return
	}

	expr = strings.TrimSpace(s)
	alias = ""
	return
}

func splitAndTrim(s, sep string) []string {
	if len(s) == 0 {
		return nil
	}
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
