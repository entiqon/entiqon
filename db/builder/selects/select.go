// File: db/builder/select.go

package selects

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/common/extension/collection"
	"github.com/entiqon/entiqon/db/dialect"
	"github.com/entiqon/entiqon/db/token/condition"
	"github.com/entiqon/entiqon/db/token/field"
	"github.com/entiqon/entiqon/db/token/join"
	"github.com/entiqon/entiqon/db/token/table"
	ct "github.com/entiqon/entiqon/db/token/types/condition"
	jt "github.com/entiqon/entiqon/db/token/types/join"
)

// SelectBuilder builds simple SELECT queries.
type selectBuilder struct {
	dialect    dialect.Dialect
	fields     *collection.Collection[field.Token]
	table      table.Token
	joins      *collection.Collection[join.Token]
	conditions *collection.Collection[condition.Token]
	groupings  *collection.Collection[string]
	sorting    *collection.Collection[string]
	having     *collection.Collection[string]
	take       int
	skip       int
}

// New creates a new SelectBuilder with the provided dialect.
// If nil is passed, BaseDialect is used by default.
func New(d dialect.Dialect) SelectBuilder {
	if d == nil {
		d = &dialect.BaseDialect{}
	}
	return &selectBuilder{dialect: d}
}

// Fields sets the SELECT list, replacing existing fields.
//
// Usage:
//
//	sb.Fields("id", "name AS username")
//
// Notes:
//   - Accepts strings, field.Token, or *field.Token.
//   - Comma-separated strings are split into multiple fields.
//   - Aliases are parsed with "AS" or space.
func (b *selectBuilder) Fields(fields ...any) SelectBuilder {
	return b.appendFields(true, fields...)
}

// AppendFields adds fields to the existing SELECT list.
//
// Usage:
//
//	sb.AppendFields("created_at")
//
// Notes:
//   - Same argument rules as Fields.
//   - Preserves any previously added fields.
func (b *selectBuilder) AppendFields(fields ...any) SelectBuilder {
	return b.appendFields(false, fields...)
}

// GetFields returns the current list of fields in the builder.
// It returns them by value to avoid external modification of the internal slice.
func (b *selectBuilder) GetFields() []field.Token {
	if b.fields == nil {
		return nil
	}
	return b.fields.Items()
}

// From sets the source table for the query.
//
// Usage:
//
//	sb.From("users u")
//	sb.From(table.New("orders", "o"))
//
// Notes:
//   - Accepts strings or table.Token.
//   - Returns an errored token if invalid.
func (b *selectBuilder) From(args ...any) SelectBuilder {
	if len(args) == 1 {
		switch v := args[0].(type) {
		case table.Token:
			b.table = v
			return b
		case *table.Token:
			if v != nil {
				b.table = *v
				return b
			}
		}
	}
	b.table = table.New(args...)
	return b
}

// Table returns the table token currently set on the builder.
//
// The returned value is always a table.Token. It will never be nil:
//   - If From(...) was called with valid arguments, the token represents that table.
//   - If From(...) was called with invalid or no arguments, the token is still returned
//     but marked as errored (IsErrored() == true).
//
// This allows callers to safely inspect the builder’s table state without
// additional nil checks.
func (b *selectBuilder) Table() table.Token {
	return b.table
}

// InnerJoin adds an INNER JOIN clause to the query.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("users u").
//	    InnerJoin("users u", "orders o", "u.id = o.user_id")
//
// The left parameter can match the current From table (with or without alias).
// If it matches, the existing table token is reused; otherwise, a new table
// token is constructed. The condition must be provided.
func (b *selectBuilder) InnerJoin(base, related any, condition string) SelectBuilder {
	return b.appendJoin(jt.Inner, base, related, condition)
}

// LeftJoin adds a LEFT JOIN clause to the query.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("users u").
//	    LeftJoin("users u", "orders o", "u.id = o.user_id")
//
// The left parameter can match the current From table (with or without alias).
// If it matches, the existing table token is reused; otherwise, a new table
// token is constructed. The condition must be provided.
func (b *selectBuilder) LeftJoin(base, related any, condition string) SelectBuilder {
	return b.appendJoin(jt.Left, base, related, condition)
}

// RightJoin adds a RIGHT JOIN clause to the query.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("users u").
//	    RightJoin("users u", "orders o", "u.id = o.user_id")
//
// The left parameter can match the current From table (with or without alias).
// If it matches, the existing table token is reused; otherwise, a new table
// token is constructed. The condition must be provided.
func (b *selectBuilder) RightJoin(base, related any, condition string) SelectBuilder {
	return b.appendJoin(jt.Right, base, related, condition)
}

// FullJoin adds a FULL JOIN clause to the query.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("users u").
//	    FullJoin("users u", "orders o", "u.id = o.user_id")
//
// The left parameter can match the current From table (with or without alias).
// If it matches, the existing table token is reused; otherwise, a new table
// token is constructed. The condition must be provided.
func (b *selectBuilder) FullJoin(base, related any, condition string) SelectBuilder {
	return b.appendJoin(jt.Full, base, related, condition)
}

// CrossJoin adds a CROSS JOIN (Cartesian product) to the query.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("users u").
//	    CrossJoin("roles r")
//
// CROSS JOIN never accepts a condition; any ON clause is invalid SQL.
// Only the right-hand table is required.
func (b *selectBuilder) CrossJoin(related any) SelectBuilder {
	return b.appendJoin(jt.Cross, b.table, related, "")
}

// NaturalJoin adds a NATURAL JOIN to the query.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("employees e").
//	    NaturalJoin("departments d")
//
// NATURAL JOIN derives the join condition implicitly from columns
// with matching names. No condition is accepted by SQL.
// Only the right-hand table is required.
func (b *selectBuilder) NaturalJoin(related any) SelectBuilder {
	// Force no condition
	return b.appendJoin(jt.Natural, b.table, related, "")
}

// Joins returns the list of join tokens currently attached to the builder.
//
// The returned slice is a snapshot of the internal join collection. Each
// element is a join.Token representing a JOIN clause that was previously
// added via InnerJoin, LeftJoin, RightJoin, FullJoin, CrossJoin, or NaturalJoin.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("users u").
//	    InnerJoin("users u", "orders o", "u.id = o.user_id").
//	    LeftJoin("users u", "profiles p", "u.id = p.user_id")
//
//	joins := sb.Joins()
//	for _, j := range joins {
//	    fmt.Println(j.Kind(), j.Left().Render(), j.Right().Render())
//	}
//
// If no joins were added, the returned slice is empty. It will never be nil.
func (b *selectBuilder) Joins() []join.Token {
	if b.joins == nil {
		return nil
	}
	return b.joins.Items()
}

// Where adds one or more conditions to the builder.
//
// Behavior:
//   - If the argument is a condition.Token (or *condition.Token), it is added directly
//     and keeps its declared type (ct.Single, ct.And, ct.Or, etc.).
//   - If the argument is a raw expression (string, int, operator, etc.), it is wrapped
//     into a condition.New(...) call. For the very first condition, ct.Single is used;
//     subsequent raw expressions inherit the method type (AND for AndWhere, OR for OrWhere).
//   - Mixin is allowed: tokens and raw expressions may be passed together. Tokens are
//     preserved with their declared type, while raw expressions are normalized to the
//     operator of the method being called.
//
// Examples:
//
//	// Pure raw → Single
//	sb := builder.New(nil).From("users").
//	    Where("age", ">", 30)
//
//	// Pure tokens
//	sb := builder.New(nil).From("users").
//	    Where(condition.New(ct.Single, "age >= 45"), condition.New(ct.And, "country = 'USA'"))
//
//	// Mixin of tokens + raw
//	c := condition.New(ct.Or, "state", operator.Like, "New Jersey")
//	sb := builder.New(nil).From("users").
//	    Where(
//	        condition.New(ct.Single, "age >= 45"),
//	        condition.New(ct.And, "age <= 50"),
//	        &c,
//	        "state = 'Texas'",
//	    )
//
// Produces:
//
//	... WHERE age >= 45 AND age <= 50 OR state LIKE 'New Jersey' OR state = 'Texas'
func (b *selectBuilder) Where(args ...any) SelectBuilder {
	return b.addConditions(true, ct.Single, args...)
}

// AndWhere adds conditions combined with logical AND.
//
// Behavior:
//   - Tokens keep their declared type.
//   - Raw arguments are wrapped with ct.And.
//   - Mixed calls are valid: tokens are preserved, raw arguments inherit AND.
//
// Example:
//
//	sb := builder.New(nil).From("users").
//	    AndWhere(condition.New(ct.Or, "country = 'USA'"), "age >= 21")
//
// Produces:
//
//	... WHERE OR country = 'USA' AND age >= 21
func (b *selectBuilder) AndWhere(conditions ...any) SelectBuilder {
	return b.addConditions(false, ct.And, conditions...)
}

// OrWhere adds conditions combined with logical OR.
//
// Behavior:
//   - Tokens keep their declared type.
//   - Raw arguments are wrapped with ct.Or.
//   - Mixed calls are valid: tokens are preserved, raw arguments inherit OR.
//
// Example:
//
//	sb := builder.New(nil).From("users").
//	    OrWhere(condition.New(ct.And, "state = 'NY'"), "state = 'Texas'")
//
// Produces:
//
//	... WHERE AND state = 'NY' OR state = 'Texas'
func (b *selectBuilder) OrWhere(conditions ...any) SelectBuilder {
	return b.addConditions(false, ct.Or, conditions...)
}

// Conditions returns all condition tokens currently attached to the builder.
//
// The returned slice contains each condition added via Where, AndWhere, or OrWhere.
// Elements may be either valid or errored condition.Token values depending on how
// they were constructed. No additional validation is performed here; errors are
// surfaced later during Build().
//
// The slice is empty if no conditions have been added. If the internal collection
// has not yet been initialized, nil is returned.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("users").
//	    Where("age", ">", 30).
//	    AndWhere(condition.New(ct.Or, "country = 'USA'"))
//
//	conditions := sb.Conditions()
//	for _, c := range conditions {
//	    fmt.Println(c.Render())
//	}
//
// Produces:
//
//	age > 30
//	OR country = 'USA'
func (b *selectBuilder) Conditions() []condition.Token {
	if b.conditions == nil {
		return nil
	}
	return b.conditions.Items()
}

// GroupBy replaces any existing GROUP BY expressions with the given fields.
//
// Each argument is treated as a raw SQL expression and added in order to
// the GROUP BY clause. Passing no arguments clears all groupings.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("orders").
//	    GroupBy("customer_id", "status")
//
// // Renders:
// //   SELECT ... FROM orders GROUP BY customer_id, status
func (b *selectBuilder) GroupBy(fields ...string) SelectBuilder {
	return b.appendGroupBy(true, fields...)
}

// ThenGroupBy appends additional GROUP BY expressions to the builder.
//
// Unlike GroupBy, this method preserves any existing groupings and adds
// the new fields after them, in the order provided. Passing no arguments
// is a no-op.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("orders").
//	    GroupBy("customer_id").
//	    ThenGroupBy("status")
//
// // Renders:
// //   SELECT ... FROM orders GROUP BY customer_id, status
func (b *selectBuilder) ThenGroupBy(fields ...string) SelectBuilder {
	return b.appendGroupBy(false, fields...)
}

// Groupings returns the list of grouping expressions currently attached
// to the builder.
//
// Each string in the slice corresponds to a GROUP BY expression added via
// GroupBy() or similar methods. The slice preserves the order in which
// groupings were added.
//
// If no groupings have been added, nil is returned.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("orders").
//	    GroupBy("customer_id", "status")
//
//	groupings := sb.Groupings()
//	for _, g := range groupings {
//	    fmt.Println(g)
//	}
//
// Produces:
//
//	customer_id
//	status
func (b *selectBuilder) Groupings() []string {
	if b.groupings == nil {
		return nil
	}
	return b.groupings.Items()
}

// OrderBy replaces any existing ORDER BY expressions with the given fields.
//
// Each argument is treated as a raw SQL expression and added in order to
// the ORDER BY clause. Passing no arguments clears all sorting fields.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("orders").
//	    OrderBy("created_at DESC", "id")
//
// // Renders:
// //   SELECT ... FROM orders ORDER BY created_at DESC, id
func (b *selectBuilder) OrderBy(fields ...string) SelectBuilder {
	return b.appendOrderBy(true, fields...)
}

// ThenOrderBy appends additional ORDER BY expressions to the builder.
//
// Unlike OrderBy, this method preserves any existing sort fields and adds
// the new ones after them, in the order provided. Passing no arguments
// is a no-op.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("orders").
//	    OrderBy("created_at DESC").
//	    ThenOrderBy("id")
//
// // Renders:
// //   SELECT ... FROM orders ORDER BY created_at DESC, id
func (b *selectBuilder) ThenOrderBy(fields ...string) SelectBuilder {
	return b.appendOrderBy(false, fields...)
}

// Sorting returns the list of ORDER BY expressions currently attached
// to the builder.
//
// The returned slice preserves the order in which fields were added
// through OrderBy and ThenOrderBy. If no sorting has been defined,
// nil is returned.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("orders").
//	    OrderBy("created_at DESC", "id")
//
//	sortFields := sb.Sorting()
//	for _, f := range sortFields {
//	    fmt.Println(f)
//	}
//
// Produces:
//
//	created_at DESC
//	id
func (b *selectBuilder) Sorting() []string {
	if b.sorting == nil {
		return nil
	}
	return b.sorting.Items()
}

// Having sets the HAVING clause, replacing any existing conditions.
func (b *selectBuilder) Having(conditions ...string) SelectBuilder {
	return b.appendHaving("", true, conditions...)
}

// AndHaving appends additional conditions with AND.
func (b *selectBuilder) AndHaving(conditions ...string) SelectBuilder {
	return b.appendHaving("AND", false, conditions...)
}

// OrHaving appends additional conditions with OR.
func (b *selectBuilder) OrHaving(conditions ...string) SelectBuilder {
	return b.appendHaving("OR", false, conditions...)
}

// HavingConditions returns the list of HAVING expressions currently attached
// to the builder.
//
// The returned slice contains each condition added through Having,
// AndHaving, or OrHaving, in the order they were specified.
//
// If no HAVING conditions have been added, nil is returned.
//
// Example:
//
//	sb := builder.New(nil).
//	    From("orders").
//	    GroupBy("customer_id").
//	    Having("COUNT(*) > 5").
//	    OrHaving("SUM(amount) > 1000")
//
//	havingConditions := sb.HavingConditions()
//	for _, h := range havingConditions {
//	    fmt.Println(h)
//	}
//
// Produces:
//
//	COUNT(*) > 5
//	OR SUM(amount) > 1000
func (b *selectBuilder) HavingConditions() []string {
	if b.having == nil {
		return nil
	}
	return b.having.Items()
}

// Take sets LIMIT.
//
// Usage:
//
//	sb.Take(10)
//
// Notes:
//   - Negative values are invalid.
//   - Equivalent to SQL LIMIT.
func (b *selectBuilder) Take(value int) SelectBuilder {
	b.take = value
	return b
}

func (b *selectBuilder) Limit() int {
	return b.take
}

// Skip sets OFFSET.
//
// Usage:
//
//	sb.Skip(20)
//
// Notes:
//   - Negative values are invalid.
//   - Equivalent to SQL OFFSET.
func (b *selectBuilder) Skip(value int) SelectBuilder {
	b.skip = value
	return b
}

func (b *selectBuilder) Offset() int {
	return b.skip
}

func (b *selectBuilder) Pagination() (int, int) {
	return b.take, b.skip
}

// Debug returns a developer-facing representation of the SelectBuilder.
//
// The output is verbose and intended for diagnostics, showing the
// internal state including counts of fields, clauses, and the current
// table. It does not attempt to render SQL, unlike Build().
//
// Unlike String(), which is intended for concise human-facing audit logs,
// Debug() exposes detailed state primarily for developers.
//
// Example output:
//
//	✅ SelectBuilder{fields:2, table:✅ Table(users AS u), where:0, groupBy:0, having:0, orderBy:0}
//	❌ SelectBuilder{fields:0, table:<nil>, where:0, groupBy:0, having:0, orderBy:0}
//
// Example usage:
//
//	sb := builder.New().
//	    Fields("id", "name").
//	    From("users u")
//
//	fmt.Println(sb.Debug())
//	// Output: ✅ SelectBuilder{fields:2, table:✅ Table(users AS u), where:0, groupBy:0, having:0, orderBy:0}
func (b *selectBuilder) Debug() string {
	src := "table:<nil>"
	if b.table != nil {
		src = fmt.Sprintf("table:%s", b.table.String())
	}

	fieldsLen, joinLength, whereLen, groupLen, havingLen, orderLen := 0, 0, 0, 0, 0, 0
	if b.fields != nil {
		fieldsLen = b.fields.Length()
	}
	if b.joins != nil {
		joinLength = b.joins.Length()
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
		"SelectBuilder{%s, fields:%d, join:%d, where:%d, groupBy:%d, having:%d, orderBy:%d, limit:%d, offset:%d}",
		src,
		fieldsLen,
		joinLength,
		whereLen,
		groupLen,
		havingLen,
		orderLen,
		b.take,
		b.skip,
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
//	SelectBuilder: status: ready, fields=3, no conditions, grouped=false, sorted=true
//	SelectBuilder: status: invalid, no selected fields, no conditions, grouped=false, sorted=false
//	SelectBuilder: status: ready, fields=3, conditions=2, grouped=true, sorted=false
func (b *selectBuilder) String() string {
	status := "ready"
	if b.table == nil || !b.table.IsValid() {
		return fmt.Sprintf("SelectBuilder: status=invalid – no table specified")
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

	joined := b.joins != nil && b.joins.Length() > 0
	grouped := b.groupings != nil && b.groupings.Length() > 0
	sorted := b.sorting != nil && b.sorting.Length() > 0

	return fmt.Sprintf("SelectBuilder: status:%s, table:%s, %s, joined=%t, %s, grouped=%t, sorted=%t",
		status, b.table.String(), fieldsStr, joined, conditionsStr, grouped, sorted,
	)
}

// Build constructs the SQL query string.
func (b *selectBuilder) Build() (string, []any, error) {
	if b.table == nil {
		return "", nil, fmt.Errorf(
			"[Select] – Errors:\n  From:\n    no table specified",
		)
	}

	if b.table.IsErrored() {
		return "", nil, fmt.Errorf(
			"[Select] – Errors:\n  From:\n    %v",
			b.table.Error(),
		)
	}

	var fields string
	if b.fields != nil && b.fields.Length() > 0 {
		parts := make([]string, 0, b.fields.Length())
		var bad []string
		for _, f := range b.fields.Items() {
			if f.IsErrored() {
				bad = append(bad,
					fmt.Sprintf("Field(%q): %v", f.Input(), f.Error()))
				continue
			}
			parts = append(parts, f.Render())
		}
		if len(bad) > 0 {
			return "", nil, fmt.Errorf("[Select] - Fields:\n\t%s", strings.Join(bad, "\n\t"))
		}

		fields = strings.Join(parts, ", ")
	} else {
		fields = field.New("*").Render()
	}

	tokens := []string{
		"SELECT",
		fields,
		"FROM",
		b.table.Render(),
	}

	sql := strings.Join(tokens, " ")

	if b.joins != nil && b.joins.Length() > 0 {
		parts := make([]string, 0, b.joins.Length())
		var bad []string
		for _, j := range b.joins.Items() {
			if j.IsErrored() {
				bad = append(bad, fmt.Sprintf("Join(%q): %v", j.Left(), j.Error()))
				continue
			}
			parts = append(parts, j.Render())
		}
		if len(bad) > 0 {
			return "", nil, fmt.Errorf("[Select] - Join:\n\t%s", strings.Join(bad, "\n\t"))
		}

		sql += " " + strings.Join(parts, " ")
	}

	var values []any
	if b.conditions != nil && b.conditions.Length() > 0 {
		parts := make([]string, 0, b.conditions.Length())
		var bad []string

		for _, c := range b.conditions.Items() {
			if c.IsErrored() {
				bad = append(bad,
					fmt.Sprintf("Condition(%q): %v", c.Input(), c.Error()))
				continue
			}
			parts = append(parts, c.Render())
			if v := c.Value(); v != nil {
				values = append(values, v)
			}
		}

		if len(bad) > 0 {
			return "", nil, fmt.Errorf(
				"[Select] - Where:\n\t%s",
				strings.Join(bad, "\n\t"),
			)
		}

		if len(parts) > 0 {
			sql += " WHERE " + strings.Join(parts, " ")
		}
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

	if b.take > 0 {
		sql += fmt.Sprintf(" LIMIT %d", b.take)
	}
	if b.skip > 0 {
		sql += fmt.Sprintf(" OFFSET %d", b.skip)
	}

	return strings.TrimSpace(sql), values, nil
}

// appendFields is the shared logic for parsing/adding fields.
func (b *selectBuilder) appendFields(reset bool, fields ...any) SelectBuilder {
	if b.fields == nil {
		b.fields = collection.New[field.Token]()
	} else if reset {
		b.fields.Clear()
	}

	if len(fields) == 0 {
		return b
	}

	switch len(fields) {
	case 1:
		switch v := fields[0].(type) {
		case string:
			s := strings.TrimSpace(v)
			if strings.Contains(s, ",") {
				for _, part := range splitAndTrim(s, ",") {
					b.fields.Add(field.New(part))
				}
			} else {
				b.fields.Add(field.New(s))
			}

		case field.Token:
			b.fields.Add(v)

		case *field.Token:
			if v != nil {
				b.fields.Add(*v)
			}

		default:
			// fallback to generic constructor
			b.fields.Add(field.New(v))
		}

	default:
		// 2 args → expr + alias
		// >2 args → will be rejected by field.New
		b.fields.Add(field.New(fields...))
	}

	return b
}

// appendJoin is an internal helper that constructs and appends a JOIN clause
// to the builder's join collection. It should not be used directly; instead,
// call one of the public Join methods (InnerJoin, LeftJoin, RightJoin, FullJoin,
// CrossJoin, or NaturalJoin).
//
// Behavior:
//   - If the left argument matches the current table table (by Render() or Name()),
//     the existing table token is reused.
//   - Otherwise, a new table.Token is constructed from the left argument.
//   - The right argument is passed through to join.New* constructors as-is.
//   - The condition is only relevant for INNER, LEFT, RIGHT, and FULL joins.
//     CROSS and NATURAL ignore it.
//
// The join kind is determined by the jt.Type parameter; if a specific constructor
// is not matched, join.New(kind, ...) is used as a fallback.
func (b *selectBuilder) appendJoin(
	kind jt.Type,
	left any,
	right any,
	on string,
) *selectBuilder {
	if b.joins == nil {
		b.joins = collection.New[join.Token]()
	}

	leftTable := resolveBase(left, b.table)

	if kind == jt.Inner {
		b.joins.Add(join.NewInner(leftTable, right, on))
	} else if kind == jt.Left {
		b.joins.Add(join.NewLeft(leftTable, right, on))
	} else if kind == jt.Right {
		b.joins.Add(join.NewRight(leftTable, right, on))
	} else if kind == jt.Full {
		b.joins.Add(join.NewFull(leftTable, right, on))
	} else if kind == jt.Cross {
		b.joins.Add(join.NewCross(leftTable, right))
	} else if kind == jt.Natural {
		b.joins.Add(join.NewNatural(leftTable, right))
	}
	return b
}

func resolveBase(left any, src table.Token) table.Token {
	switch v := left.(type) {
	case table.Token:
		// if semantically equal, reuse src
		if sameFromTable(v, src) {
			return src
		}
		return v
	case *table.Token:
		if sameFromTable(*v, src) {
			return src
		}
		return *v
	case string:
		// compare with source string forms
		if v == src.Input() || v == src.Render() || v == src.Name() || v == src.String() {
			return src
		}
		return table.New(v)
	default:
		return table.New(left)
	}
}

func sameFromTable(a, b table.Token) bool {
	if a == nil || b == nil {
		return false
	}
	return a.Name() == b.Name() && a.Alias() == b.Alias()
}

// addConditions adds new conditions. If reset=true, clears first.
//
// Behavior:
//   - condition.Token → added directly.
//   - *condition.Token → dereferenced and added.
//   - Any other type (string, int, etc.) → accumulated into a slice and passed
//     as arguments to condition.New(kind, rawArgs...). The result is added as a token.
//   - Errors are carried by condition.New and surfaced later in Build().
//
// This allows mixing tokens and raw expressions in one call.
func (b *selectBuilder) addConditions(
	reset bool,
	kind ct.Type,
	args ...any,
) SelectBuilder {
	if b.conditions == nil {
		b.conditions = collection.New[condition.Token]()
	} else if reset {
		b.conditions.Clear()
	}

	if len(args) == 0 {
		return b
	}

	var rawArgs []any
	flushRaw := func() {
		if len(rawArgs) == 0 {
			return
		}
		rawType := kind
		if b.conditions.Length() == 0 {
			rawType = ct.Single
		} else if rawType == ct.Single {
			rawType = ct.And
		}
		b.conditions.Add(condition.New(rawType, rawArgs...))
		rawArgs = nil
	}

	for _, a := range args {
		switch v := a.(type) {
		case condition.Token:
			flushRaw()
			b.conditions.Add(v)
		case *condition.Token:
			flushRaw()
			if v != nil {
				b.conditions.Add(*v)
			}
		default:
			rawArgs = append(rawArgs, v)
		}
	}
	flushRaw()

	return b
}

// appendGroupBy ensures init and handles reset
func (b *selectBuilder) appendGroupBy(reset bool, fields ...string) SelectBuilder {
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

// appendOrderBy ensures init and handles reset
func (b *selectBuilder) appendOrderBy(reset bool, fields ...string) SelectBuilder {
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

// appendHaving adds one or more conditions to the builder.
// prefix is "", "AND", or "OR" depending on caller semantics.
// If reset is true, the collection is cleared first.
func (b *selectBuilder) appendHaving(
	prefix string,
	reset bool,
	conditions ...string,
) *selectBuilder {
	if b.having == nil {
		b.having = collection.New[string]()
	} else if reset {
		b.having.Clear()
	}

	for _, c := range conditions {
		trimmed := strings.TrimSpace(c)
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
