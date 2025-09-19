package selects

import (
	"github.com/entiqon/entiqon/db/contract"
	"github.com/entiqon/entiqon/db/token/condition"
	"github.com/entiqon/entiqon/db/token/field"
	"github.com/entiqon/entiqon/db/token/join"
	"github.com/entiqon/entiqon/db/token/table"
)

// SelectBuilder defines the contract for constructing SQL SELECT queries.
//
// It provides methods for defining fields, source tables, joins, conditions,
// grouping, sorting, HAVING clauses, pagination, and building the final query.
// Each mutator returns the builder for chaining; accessors return the current state.
//
// Methods:
//   - Fields / AppendFields / GetFields: define and retrieve the SELECT list
//   - From / Table: set or get the source table
//   - InnerJoin / LeftJoin / RightJoin / FullJoin / CrossJoin / NaturalJoin / Joins: manage JOIN clauses
//   - Where / AndWhere / OrWhere / Conditions: manage WHERE conditions
//   - GroupBy / ThenGroupBy / Groupings: manage GROUP BY expressions
//   - OrderBy / ThenOrderBy / Sorting: manage ORDER BY expressions
//   - Having / AndHaving / OrHaving / HavingConditions: manage HAVING conditions
//   - Take / Limit / Skip / Offset / Pagination: manage LIMIT and OFFSET
//   - Build: construct the final SQL string
//   - Debug / String: return diagnostic or human-readable views
type SelectBuilder interface {
	contract.Debuggable
	contract.Stringable

	// Fields sets the SELECT list, replacing existing fields.
	//
	// Notes:
	//   • Accepts strings, field.Token, or *field.Token.
	//   • Comma-separated strings are split into multiple fields.
	//   • Aliases may be declared with "AS" or a space.
	Fields(fields ...any) SelectBuilder

	// AppendFields adds fields to the SELECT list without clearing existing ones.
	AppendFields(fields ...any) SelectBuilder

	// GetFields returns the current SELECT fields.
	//
	// Notes:
	//   • Returns nil if no fields are defined.
	GetFields() []field.Token

	// From sets the source table.
	//
	// Notes:
	//   • Accepts strings or table.Token.
	//   • Returns an errored token if invalid.
	From(args ...any) SelectBuilder

	// Table returns the current source table token.
	//
	// Notes:
	//   • May be nil if not set.
	//   • May be errored if invalid.
	Table() table.Token

	// InnerJoin adds an INNER JOIN clause.
	InnerJoin(base any, related any, condition string) SelectBuilder

	// LeftJoin adds a LEFT JOIN clause.
	LeftJoin(base any, related any, condition string) SelectBuilder

	// RightJoin adds a RIGHT JOIN clause.
	RightJoin(base any, related any, condition string) SelectBuilder

	// FullJoin adds a FULL JOIN clause.
	FullJoin(base any, related any, condition string) SelectBuilder

	// CrossJoin adds a CROSS JOIN clause (no condition).
	CrossJoin(related any) SelectBuilder

	// NaturalJoin adds a NATURAL JOIN clause (implicit condition).
	NaturalJoin(related any) SelectBuilder

	// Joins returns all JOIN clauses.
	//
	// Notes:
	//   • Returns nil if no joins are defined.
	Joins() []join.Token

	// Where sets the WHERE conditions, replacing existing ones.
	//
	// Notes:
	//   • Accepts condition.Token, *condition.Token, or raw expressions.
	//   • Tokens preserve their declared type.
	//   • Raw expressions adopt the operator of the method used.
	Where(args ...any) SelectBuilder

	// AndWhere appends conditions combined with AND.
	AndWhere(args ...any) SelectBuilder

	// OrWhere appends conditions combined with OR.
	OrWhere(args ...any) SelectBuilder

	// Conditions returns all WHERE conditions.
	//
	// Notes:
	//   • Returns nil if none defined.
	Conditions() []condition.Token

	// GroupBy replaces the GROUP BY clause.
	//
	// Notes:
	//   • Preserves order of fields.
	//   • Passing no arguments clears groupings.
	GroupBy(fields ...string) SelectBuilder

	// ThenGroupBy appends additional GROUP BY fields.
	ThenGroupBy(fields ...string) SelectBuilder

	// Groupings returns all GROUP BY fields.
	Groupings() []string

	// OrderBy replaces the ORDER BY clause.
	//
	// Notes:
	//   • Passing no arguments clears sorting.
	//   • Use ThenOrderBy to append instead.
	OrderBy(fields ...string) SelectBuilder

	// ThenOrderBy appends additional ORDER BY fields.
	ThenOrderBy(fields ...string) SelectBuilder

	// Sorting returns all ORDER BY expressions.
	Sorting() []string

	// Having sets the HAVING clause, replacing existing conditions.
	//
	// Notes:
	//   • Accepts raw strings only.
	//   • Use AndHaving / OrHaving to append.
	Having(conditions ...string) SelectBuilder

	// AndHaving appends HAVING conditions with AND.
	AndHaving(conditions ...string) SelectBuilder

	// OrHaving appends HAVING conditions with OR.
	OrHaving(conditions ...string) SelectBuilder

	// HavingConditions returns all HAVING conditions.
	HavingConditions() []string

	// Take sets LIMIT.
	//
	// Notes:
	//   • Negative values are invalid.
	Take(value int) SelectBuilder

	// Limit returns the LIMIT value.
	//
	// Notes:
	//   • Returns 0 if unset.
	Limit() int

	// Skip sets OFFSET.
	//
	// Notes:
	//   • Negative values are invalid.
	Skip(value int) SelectBuilder

	// Offset returns the OFFSET value.
	Offset() int

	// Pagination returns LIMIT and OFFSET values.
	Pagination() (int, int)

	// Build constructs the final SQL string.
	//
	// Returns:
	//   • SQL string
	//   • Bound values
	//   • Error if invalid
	Build() (string, []interface{}, error)
}

var _ SelectBuilder = (*selectBuilder)(nil)
