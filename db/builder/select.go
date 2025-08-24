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
	dialect    dialect.Dialect
	fields     *collection.Collection[token.Field]
	table      string
	conditions *collection.Collection[string]
	groupings  *collection.Collection[string]
	sorting    *collection.Collection[string]
	having     *collection.Collection[string]
	limit      int
	offset     int
}

// NewSelect creates a new SelectBuilder with the provided dialect.
// If nil is passed, BaseDialect is used by default.
func NewSelect(d dialect.Dialect) *SelectBuilder {
	if d == nil {
		d = &dialect.BaseDialect{}
	}
	return &SelectBuilder{
		dialect:    d,
		fields:     collection.New[token.Field](),
		conditions: collection.New[string](),
		sorting:    collection.New[string](),
		groupings:  collection.New[string](),
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

// Where sets the initial WHERE clause conditions, replacing any existing ones.
// Empty or whitespace-only expressions are ignored.
func (b *SelectBuilder) Where(conditions ...string) *SelectBuilder {
	return b.addConditions("", true, conditions...)
}

// And appends additional conditions to the WHERE clause.
// Empty or whitespace-only expressions are ignored.
func (b *SelectBuilder) And(conditions ...string) *SelectBuilder {
	return b.addConditions("AND", false, conditions...)
}

// Or appends additional conditions to the WHERE clause with OR.
// Empty or whitespace-only expressions are ignored.
// Or appends additional conditions joined with OR.
// If this is the first condition, it behaves like Where().
func (b *SelectBuilder) Or(conditions ...string) *SelectBuilder {
	return b.addConditions("OR", false, conditions...)
}

// GroupBy sets the GROUP BY clause, replacing existing groupings.
func (b *SelectBuilder) GroupBy(fields ...string) *SelectBuilder {
	return b.addGroupBy(true, fields...)
}

// ThenGroupBy appends additional GROUP BY fields, preserving existing ones.
func (b *SelectBuilder) ThenGroupBy(fields ...string) *SelectBuilder {
	return b.addGroupBy(false, fields...)
}

// OrderBy sets the ORDER BY clause, replacing existing order fields.
func (b *SelectBuilder) OrderBy(fields ...string) *SelectBuilder {
	return b.addOrderBy(true, fields...)
}

// Having sets the HAVING clause, replacing any existing conditions.
func (b *SelectBuilder) Having(conditions ...string) *SelectBuilder {
	return b.addHaving("", true, conditions...)
}

// AndHaving appends additional conditions with AND.
func (b *SelectBuilder) AndHaving(conditions ...string) *SelectBuilder {
	return b.addHaving("AND", false, conditions...)
}

// OrHaving appends additional conditions with OR.
func (b *SelectBuilder) OrHaving(conditions ...string) *SelectBuilder {
	return b.addHaving("OR", false, conditions...)
}

// ThenOrderBy appends additional ORDER BY fields, preserving existing ones.
func (b *SelectBuilder) ThenOrderBy(fields ...string) *SelectBuilder {
	return b.addOrderBy(false, fields...)
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

	if b.conditions != nil && b.conditions.Length() > 0 {
		sql += " WHERE " + strings.Join(b.conditions.Items(), " ")
	}

	if b.groupings != nil && b.groupings.Length() > 0 {
		sql += " GROUP BY " + strings.Join(b.groupings.Items(), ", ")
	}

	if b.sorting != nil && b.sorting.Length() > 0 {
		sql += " ORDER BY " + strings.Join(b.sorting.Items(), ", ")
	}

	if b.having != nil && b.having.Length() > 0 {
		sql += " HAVING " + strings.Join(b.having.Items(), " ")
	}

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

// addConditions adds one or more conditions to the builder.
// prefix is "", "AND", or "OR" depending on caller semantics.
// If reset is true, the collection is cleared first.
func (b *SelectBuilder) addConditions(
	prefix string,
	reset bool, conditions ...string,
) *SelectBuilder {
	if b.conditions == nil {
		b.conditions = collection.New[string]()
	} else if reset {
		b.conditions.Clear()
	}

	for _, cond := range conditions {
		trimmed := strings.TrimSpace(cond)
		if trimmed == "" {
			continue
		}

		// First condition → always stored without prefix
		if b.conditions.Length() == 0 {
			b.conditions.Add(trimmed)
			continue
		}

		// Subsequent conditions
		if prefix != "" {
			b.conditions.Add(prefix + " " + trimmed)
		} else {
			// If no prefix (Where with multiple args), default to AND
			b.conditions.Add("AND " + trimmed)
		}
	}
	return b
}

// addGroupBy ensures init and handles reset
func (b *SelectBuilder) addGroupBy(reset bool, fields ...string) *SelectBuilder {
	if b.groupings == nil {
		b.groupings = collection.New[string]()
	} else if reset {
		b.groupings.Clear()
	}

	for _, f := range fields {
		trimmed := strings.TrimSpace(f)
		if trimmed != "" {
			b.groupings.Add(trimmed)
		}
	}
	return b
}

// addHaving adds one or more conditions to the builder.
// prefix is "", "AND", or "OR" depending on caller semantics.
// If reset is true, the collection is cleared first.
func (b *SelectBuilder) addHaving(prefix string, reset bool, conditions ...string) *SelectBuilder {
	if b.having == nil {
		b.having = collection.New[string]()
	} else if reset {
		b.having.Clear()
	}

	for _, cond := range conditions {
		trimmed := strings.TrimSpace(cond)
		if trimmed == "" {
			continue
		}

		if b.having.Length() == 0 {
			b.having.Add(trimmed)
		} else if prefix != "" {
			b.having.Add(prefix + " " + trimmed)
		} else {
			// default joiner for multiple conditions in Having is AND
			b.having.Add("AND " + trimmed)
		}
	}
	return b
}

// addOrderBy ensures init and handles reset
func (b *SelectBuilder) addOrderBy(reset bool, fields ...string) *SelectBuilder {
	if b.sorting == nil {
		b.sorting = collection.New[string]()
	} else if reset {
		b.sorting.Clear()
	}

	for _, f := range fields {
		trimmed := strings.TrimSpace(f)
		if trimmed != "" {
			b.sorting.Add(trimmed)
		}
	}
	return b
}
