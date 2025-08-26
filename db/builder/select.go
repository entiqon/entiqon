// File: db/builder/select.go

package builder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/common/extension/collection"
	"github.com/entiqon/entiqon/db/dialect"
	"github.com/entiqon/entiqon/db/token/field"
	"github.com/entiqon/entiqon/db/token/table"
)

// SelectBuilder builds simple SELECT queries.
type SelectBuilder struct {
	dialect    dialect.Dialect
	fields     *collection.Collection[field.Field]
	source     *table.Table
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
		fields:     collection.New[field.Field](),
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
		b.fields = collection.New[field.Field]()
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
func (b *SelectBuilder) GetFields() *collection.Collection[field.Field] {
	return b.fields
}

// Source sets the source table for the SELECT statement.
// Accepts the same argument patterns as table.New:
//   - Source("users")
//   - Source("users u")
//   - Source("users AS u")
//   - Source("users", "u")
//   - Source("(SELECT ...)", "t")
//   - Source(table.New("orders", "o"))
func (b *SelectBuilder) Source(args ...any) *SelectBuilder {
	switch v := args[0].(type) {
	case string:
		strArgs := make([]string, 0, len(args))
		for _, a := range args {
			strArgs = append(strArgs, fmt.Sprintf("%v", a))
		}
		b.source = table.New(strArgs...)
	case *table.Table:
		if len(args) == 1 {
			b.source = v
		} else {
			// multiple args but first is *table.Table → invalid
			b.source = table.New("too many arguments")
		}
	default:
		// anything else → fallback to errored table
		b.source = table.New(fmt.Sprintf("%v", v))
	}

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

// Debug returns a developer-facing representation of the SelectBuilder.
//
// The output is verbose and intended for diagnostics, showing the
// internal state including counts of fields, clauses, and the current
// source. It does not attempt to render SQL, unlike Build().
//
// Unlike String(), which is intended for concise human-facing audit logs,
// Debug() exposes detailed state primarily for developers.
//
// Example output:
//
//	✅ SelectBuilder{fields:2, source:✅ Table(users AS u), where:0, groupBy:0, having:0, orderBy:0}
//	❌ SelectBuilder{fields:0, source:<nil>, where:0, groupBy:0, having:0, orderBy:0}
//
// Example usage:
//
//	sb := builder.NewSelect().
//	    Fields("id", "name").
//	    Source("users u")
//
//	fmt.Println(sb.Debug())
//	// Output: ✅ SelectBuilder{fields:2, source:✅ Table(users AS u), where:0, groupBy:0, having:0, orderBy:0}
func (b *SelectBuilder) Debug() string {
	if b == nil {
		return "❌ SelectBuilder(nil)"
	}

	status := "✅"
	if b.source == nil || b.source.IsErrored() {
		status = "❌"
	}

	src := "source:<nil>"
	if b.source != nil {
		src = fmt.Sprintf("source: %s", b.source.String())
	}

	fieldsLen, whereLen, groupLen, havingLen, orderLen := 0, 0, 0, 0, 0
	if b.fields != nil {
		fieldsLen = b.fields.Length()
	}
	if b.conditions != nil {
		whereLen = b.conditions.Length()
	}
	if b.groupings != nil {
		groupLen = b.groupings.Length()
	}
	if b.having != nil {
		havingLen = b.having.Length()
	}
	if b.sorting != nil {
		orderLen = b.sorting.Length()
	}

	return fmt.Sprintf(
		"%s SelectBuilder{fields:%d, %s, where:%d, groupBy:%d, having:%d, orderBy:%d}",
		status,
		fieldsLen,
		src,
		whereLen,
		groupLen,
		havingLen,
		orderLen,
	)
}

// String returns the human-facing representation of the SelectBuilder.
//
// It provides a concise status log of the builder, suitable for audits
// and logs. Unlike Build(), it does not attempt to render SQL.
//
// The conditions count is the sum of WHERE and HAVING clauses.
//
// Example output:
//
//	✅ SelectBuilder: status: ready, fields=3, no conditions, grouped=false, sorted=true
//	❌ SelectBuilder: status: invalid, no selected fields, no conditions, grouped=false, sorted=false
//	✅ SelectBuilder: status: ready, fields=3, conditions=2, grouped=true, sorted=false
func (b *SelectBuilder) String() string {
	if b == nil {
		return "❌ SelectBuilder: status=<nil> – wrong initialization"
	}

	status, icon := "ready", "✅"
	if b.source == nil || !b.source.IsValid() {
		return fmt.Sprintf("❌ SelectBuilder: status=invalid – no source specified")
	}

	fieldsStr := "no selected fields"
	if b.fields != nil && b.fields.Length() > 0 {
		fieldsStr = fmt.Sprintf("fields=%d", b.fields.Length())
	}

	// total conditions = WHERE + HAVING
	totalConditions := 0
	if b.conditions != nil {
		totalConditions += b.conditions.Length()
	}
	if b.having != nil {
		totalConditions += b.having.Length()
	}

	conditionsStr := "no conditions"
	if totalConditions > 0 {
		conditionsStr = fmt.Sprintf("conditions=%d", totalConditions)
	}

	grouped := b.groupings != nil && b.groupings.Length() > 0
	sorted := b.sorting != nil && b.sorting.Length() > 0

	return fmt.Sprintf(
		"%s SelectBuilder: status: %s, %s, %s, grouped=%v, sorted=%v",
		icon, status, fieldsStr, conditionsStr, grouped, sorted,
	)
}

// Build constructs the SQL query string.
func (b *SelectBuilder) Build() (string, error) {
	if b == nil {
		return "", fmt.Errorf("❌ [Build] - Wrong initialization. Cannot build on receiver nil")
	}

	if b.source == nil || !b.source.IsValid() {
		return "", errors.New(b.String())
	}

	if b.fields.Length() == 0 {
		b.fields.Add(*field.New("*"))
	}

	parts := make([]string, 0, b.fields.Length())
	var bad []string
	for _, f := range b.fields.Items() {
		if f.IsErrored() {
			bad = append(bad,
				fmt.Sprintf("⛔️ Field(%q): %v", f.Input(), f.Error()))
			continue
		}
		parts = append(parts, f.Render())
	}
	if len(bad) > 0 {
		return "", fmt.Errorf("❌ [Build] - Invalid fields:\n\t%s", strings.Join(bad, "\n\t"))
	}

	tokens := []string{
		"SELECT",
		strings.Join(parts, ", "),
		"FROM",
		b.source.Render(),
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
					b.fields.Add(*field.New(part))
				}
			} else {
				b.fields.Add(*field.New(v))
			}
		case field.Field:
			b.fields.Add(v)
		case *field.Field:
			if v != nil {
				b.fields.Add(*v)
			}
		default:
			b.fields.Add(*field.New(v))
		}
	case 2, 3:
		b.fields.Add(*field.New(cols...))
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
