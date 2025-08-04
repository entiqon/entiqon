// filename: db/internal/build/token/table.go

package token

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/driver"
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
		return (&Table{BaseToken: NewBaseToken(expr)}).SetErrorWith(expr, fmt.Errorf("table expression is empty"))
	}

	if strings.Contains(expr, ",") {
		return (&Table{BaseToken: NewBaseToken(expr)}).
			SetErrorWith(expr, fmt.Errorf("invalid table expression: unexpected comma — aliases must not be comma-separated"))
	}

	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "AS ") {
		return (&Table{BaseToken: NewBaseToken(expr)}).
			SetErrorWith(expr, fmt.Errorf("invalid table expression: cannot start with 'AS'"))
	}

	base, parsedAlias := ParseAlias(expr)
	source := &Table{BaseToken: NewBaseToken(expr)}
	source.SetName(base)

	if len(alias) > 0 && alias[0] != "" {
		source.SetAlias(alias[0])
		if parsedAlias != "" && alias[0] != parsedAlias {
			return source.SetErrorWith(expr, fmt.Errorf(
				"alias conflict: explicit alias %q does not match inline alias %q",
				alias[0],
				parsedAlias,
			))
		}
	} else {
		source.SetAlias(parsedAlias)
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
		return fmt.Sprintf("%s AS %s", t.GetName(), t.GetAlias())
	}
	return t.GetName()
}

// Render returns the dialect-quoted table name and alias (if present).
//
// # Example
//
//	tbl := NewTable("users u")
//	fmt.Println(tbl.Render(postgres)) → `"users" AS "u"`
func (t *Table) Render(d driver.Dialect) string {
	if t == nil || t.GetName() == "" {
		return ""
	}
	if t.IsAliased() {
		return fmt.Sprintf("%s AS %s", d.QuoteIdentifier(t.GetName()), d.QuoteIdentifier(t.GetAlias()))
	}
	return d.QuoteIdentifier(t.GetName())
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
	if t == nil {
		return "Table(nil)"
	}

	s := fmt.Sprintf("Table(%q) [aliased: %t, errored: %t]", t.GetName(), t.IsAliased(), t.IsErrored())
	if t.IsErrored() {
		s += fmt.Sprintf(" – %s", t.GetError().Error())
	}
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
