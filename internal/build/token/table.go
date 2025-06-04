// filename: internal/build/token/table.go

package token

import (
	"fmt"
	"strings"
)

// Table represents a SQL table reference, typically used in FROM, INTO, or JOIN clauses.
//
// It embeds BaseToken to track the name, optional alias, and semantic error state.
// A Table token can be qualified with an alias (e.g., "users AS u") but does not support
// column-level detail — it represents the table object only.
//
// Since: v1.6.0
type Table struct {
	*BaseToken
}

// NewErroredTable creates a Table token with an attached error.
//
// This is used when a table reference fails parsing, validation, or resolution.
// The resulting Table will have its Name and Alias unset, and the provided
// error will be stored in its BaseToken for tracking and diagnostics.
//
// # Example
//
//	tbl := token.NewErroredTable(fmt.Errorf("empty table reference"))
//	fmt.Println(tbl.String())
//
//	// Output:
//	Table("") [aliased: false, errored: true, error: empty table reference]
//
// This constructor allows builders and parsers to retain structurally invalid
// tokens in the build stream while still reporting errors through validators.
func NewErroredTable(err error) *Table {
	return &Table{BaseToken: NewErroredToken(err)}
}

// NewTable creates a Table token from an expression and optional alias.
//
// Supports aliasing via:
//   - "AS" keyword: "users AS u"
//   - space-delimited: "users u"
//
// If an inline alias and explicit alias are both present but conflict,
// the returned Table includes an error.
//
// Comma-separated expressions are not allowed and result in an errored Table.
//
// # Examples:
//
//	NewTable("users")                    → Table{Name: "users"}
//	NewTable("users AS u")               → Table{Name: "users", Alias: "u"}
//	NewTable("users", "u")               → Table{Name: "users", Alias: "u"}
//	NewTable("users AS x", "y")          → Table{Error: alias mismatch}
//	NewTable("users, orders")            → Table{Error: comma-separated input not allowed}
func NewTable(expr string, alias ...string) *Table {
	trimmed := strings.TrimSpace(expr)
	if trimmed == "" {
		return (&Table{BaseToken: &BaseToken{}}).SetErrorWith(expr, fmt.Errorf("table expression is empty"))
	}

	if strings.Contains(expr, ",") {
		return (&Table{BaseToken: &BaseToken{}}).
			SetErrorWith(expr, fmt.Errorf("invalid table expression: aliases must not be comma-separated"))
	}

	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "AS ") {
		return (&Table{BaseToken: &BaseToken{}}).
			SetErrorWith(expr, fmt.Errorf("invalid table expression: cannot start with 'AS'"))
	}

	base, parsedAlias := ParseAlias(expr)
	if base == "" {
		return (&Table{BaseToken: &BaseToken{}}).
			SetErrorWith(expr, fmt.Errorf("source table is empty"))
	}

	source := &Table{
		BaseToken: &BaseToken{
			Source: expr,
			Name:   base,
		},
	}

	if len(alias) > 0 && alias[0] != "" {
		source.Alias = alias[0]
		if parsedAlias != "" && alias[0] != parsedAlias {
			return source.SetErrorWith(expr, fmt.Errorf(
				"alias conflict: explicit alias %q does not match inline alias %q",
				alias[0],
				parsedAlias,
			))
		}
	} else {
		source.Alias = parsedAlias
	}

	return source
}

// Raw returns the SQL-safe table reference, including alias if present.
//
// This output is not dialect-quoted. Quoting should be handled during rendering.
//
// # Examples
//
//	Table{Name: "users"} → "users"
//	Table{Name: "users", Alias: "u"} → "users AS u"
func (t *Table) Raw() string {
	if t.IsAliased() {
		return fmt.Sprintf("%s AS %s", t.Name, t.Alias)
	}
	return t.Name
}

// String returns a debug-friendly view of the Table token.
//
// The output includes the table name, alias status, and error state if applicable.
// This is intended for logging and diagnostics, not SQL rendering.
//
// # Example Output
//
//	Table("users") [aliased: false, errored: false]
//	Table("users") [aliased: true, errored: true, error: alias mismatch]
func (t *Table) String() string {
	s := fmt.Sprintf("Table(%q) [aliased: %v, errored: %v", t.Name, t.IsAliased(), t.HasError())
	if t.HasError() {
		s += fmt.Sprintf(", error: %s", t.Error.Error())
	}
	s += "]"
	return s
}

// SetErrorWith records an error and source expression on the table token.
//
// This method delegates to BaseToken.SetErrorWith and returns the updated *Table,
// allowing fluent chaining in constructors or resolution logic.
//
// # Example
//
//	t := NewTable("users AS u")
//	t.SetErrorWith("users AS u", fmt.Errorf("alias not allowed"))
//	fmt.Println(t.String())
//
//	// Output:
//	Table("users") [aliased: true, errored: true, error: alias not allowed]
func (t *Table) SetErrorWith(expr string, err error) *Table {
	t.BaseToken.SetErrorWith(expr, err)
	return t
}

// Ensure Table satisfies the GenericToken interface.
var _ GenericToken = &Table{}
