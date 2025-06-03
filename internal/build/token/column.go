package token

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/driver"
)

// Column represents a SQL column reference within a SELECT, WHERE, or ORDER BY clause.
//
// It supports optional table qualification (e.g., "users.id") and aliasing (e.g., "id AS user_id").
// The Column token embeds a BaseToken to handle naming, aliasing, and error tracking.
//
// It is designed to participate in query builders or token streams where SQL elements
// need to be validated, introspected, and rendered dialect-neutrally.
//
// # Usage
//
//	c := NewColumn("users.id AS uid")
//	fmt.Println(c.Raw())     	// "users.id"
//	fmt.Println(c.IsAliased()) 	// true
//	fmt.Println(c.String())  	// Column("users.id") [aliased: true, qualified: true]
//
// File: internal/build/token/column.go
// Since: v1.6.0
type Column struct {
	*BaseToken

	// Table holds the attached table token for qualification and rendering.
	// It is automatically parsed or explicitly set using WithTable().
	// When non-nil and valid, the column is considered qualified.
	//
	// # Example
	//
	//	c := NewColumn("id", "alias").WithTable(NewTable("users AS u"))
	//	fmt.Println(c.Raw()) // "u.alias"
	Table *Table

	// TableName holds the extracted table prefix (if any) from the original expression.
	// Used for validation or debugging, not for rendering.
	//
	// # Example
	//
	//	c := NewColumn("users.id")
	//	fmt.Println(c.TableName) // "users"
	TableName string
}

// NewColumn constructs a Column token from a raw expression and optional alias.
//
// If no alias is provided, it attempts to extract one inline using the SQL keyword "AS",
// or through space/comma separation. If both an inline alias and an explicit alias
// are provided, and they conflict, an error is stored in the token.
//
// If the expression is qualified (e.g., "users.id"), the table prefix is used
// to create a Table token automatically, which informs downstream resolution.
//
// # Examples
//
//	NewColumn("id")                    → name: "id"
//	NewColumn("id AS uid")             → name: "id", alias: "uid"
//	NewColumn("users.id")              → name: "id", table: "users"
//	NewColumn("id", "alias")           → name: "id", alias: "alias"
//	NewColumn("users.id", "alias")     → name: "id", table: "users", alias: "alias"
func NewColumn(expr string, alias ...string) *Column {
	trimmed := strings.TrimSpace(expr)
	if trimmed == "" {
		return (&Column{BaseToken: &BaseToken{}}).
			SetErrorWith(expr, fmt.Errorf("column expression is empty"))
	}

	if strings.Contains(expr, ",") {
		return (&Column{BaseToken: &BaseToken{}}).
			SetErrorWith(expr, fmt.Errorf("invalid column expression: aliases must not be comma-separated"))
	}

	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "AS ") {
		return (&Column{BaseToken: &BaseToken{}}).
			SetErrorWith(expr, fmt.Errorf("invalid column expression: cannot start with 'AS'"))
	}

	base, parsedAlias := ParseAlias(expr)
	tableName, column := ParseTableColumn(base)
	if column == "" {
		return (&Column{BaseToken: &BaseToken{}}).
			SetErrorWith(expr, fmt.Errorf("column name is required"))
	}

	col := &Column{
		BaseToken: &BaseToken{
			Source: expr,
			Name:   column,
		},
	}

	if len(alias) > 0 && alias[0] != "" {
		col.Alias = alias[0]
		if parsedAlias != "" && alias[0] != parsedAlias {
			return col.SetErrorWith(expr, fmt.Errorf(
				"alias conflict: explicit alias %q does not match inline alias %q", alias[0], parsedAlias),
			)
		}
		parsedAlias = alias[0]
	}
	col.BaseToken.Alias = parsedAlias

	if tableName != "" {
		return col.WithTable(NewTable(tableName))
	}

	return col
}

// IsQualified reports whether the column is considered qualified.
//
// A column is qualified if:
//   - It has a non-nil Table reference (attached or parsed), AND
//   - It is structurally valid (i.e., passes IsValid)
//
// This check is used to determine if the column should render with a table prefix
// during SQL generation. It reflects the effective qualification status, not just
// string-based parsing of the expression.
//
// # Examples
//
//	NewColumn("id")                 → false
//	NewColumn("users.id")           → true
//	NewColumn("id").WithTable(...)  → true, assuming the column is valid
func (c *Column) IsQualified() bool {
	return c.Table != nil && c.Table.IsValid()
}

// Raw returns the SQL-safe, dialect-neutral representation of the column,
// including table qualification and aliasing if applicable.
//
// If the column is qualified (i.e., has a table context), the table prefix will be used.
// Table.Alias is preferred over Table.Name when available.
//
// If the column is aliased, the result will follow SQL's "AS" convention
// (e.g., "users.id AS uid").
//
// # Examples
//
//	NewColumn("id")                             → "id"
//	NewColumn("id", "alias")                    → "id AS alias"
//	NewColumn("users.id")                       → "users.id"
//	NewColumn("users.id", "alias")              → "users.id AS alias"
//	NewColumn("id", "alias").WithTable(t("u"))  → "u.alias"
//
// Raw returns the SQL-safe, dialect-neutral representation of the column.
//
// The rendering logic adapts based on qualification and aliasing status:
//
//   - If the column is qualified and aliased: it renders as `table_alias.column_alias`
//     This is typical when multiple source tables exist and aliasing prevents ambiguity.
//
//   - If the column is aliased only: it renders as `column_name AS alias`
//     Used for simple renaming in single-table contexts.
//
//   - If the column is qualified only: it renders as `table.column`
//     This is needed for joins or multi-source queries.
//
//   - If the column is neither qualified nor aliased: it renders as `column_name`
//     The simplest fallback for one-source queries.
//
// This behavior allows SelectBuilder to switch between minimal and disambiguated output
// depending on the number of sources involved in the query.
//
// # Examples
//
//	NewColumn("id")                                  	→ "id"
//	NewColumn("id", "alias")                         	→ "id AS alias"
//	NewColumn("users.id")                            	→ "users.id"
//	NewColumn("users.id", "alias")                   	→ "users.id AS alias"
//	NewColumn("id", "alias").WithTable(NewTable("u")) 	→ "u.alias"
func (c *Column) Raw() string {
	if c.IsAliased() && c.IsQualified() {
		prefix := c.Table.Alias
		if prefix == "" {
			prefix = c.Table.Name
		}
		return fmt.Sprintf("%s.%s", prefix, c.Alias)
	}
	if c.IsAliased() {
		return fmt.Sprintf("%s AS %s", c.Name, c.Alias)
	}
	if c.IsQualified() {
		prefix := c.Table.Alias
		if prefix == "" {
			prefix = c.Table.Name
		}
		return fmt.Sprintf("%s.%s", prefix, c.Name)
	}
	return c.Name
}

// Render returns the dialect-aware SQL string for the column,
// using proper quoting for table prefixes and aliasing if applicable.
//
// This method replaces Raw() when generating SQL that must be
// dialect-specific, such as for PostgreSQL, MySQL, or SQLite.
//
// It uses the table alias (if available) for qualified output,
// and falls back to the table name if not.
//
// # Examples
//
//	NewColumn("id").Render(postgres)                   		→ `"id"`
//	NewColumn("id AS uid").Render(postgres)            		→ `"id" AS "uid"`
//	NewColumn("users.id").WithTable(t("users")).Render(postgres) 	→ `"users"."id"`
//	NewColumn("id AS uid").WithTable(t("users u")).Render(postgres)	→ `"u"."id" AS "uid"`
func (c *Column) Render(d driver.Dialect) string {
	if c == nil || !c.IsValid() {
		return ""
	}

	prefix := ""
	if c.IsQualified() && c.Table != nil {
		prefix = c.Table.Alias
		if prefix == "" {
			prefix = c.Table.Name
		}
	}

	qualified := c.Name
	if prefix != "" {
		qualified = fmt.Sprintf("%s.%s", d.QuoteIdentifier(prefix), d.QuoteIdentifier(qualified))
	} else {
		qualified = d.QuoteIdentifier(c.Name)
	}

	return c.RenderAlias(d, qualified)
}

// SetErrorWith records an error and source expression on the column token.
//
// This method delegates to BaseToken.SetErrorWith and returns the updated *Column,
// enabling fluent usage during parsing or resolution failures.
//
// # Example
//
//	c := NewColumn("id AS uid", "user_id")
//	c.SetErrorWith("id AS uid", fmt.Errorf("alias conflict"))
//	fmt.Println(c.String())
//
//	// Output:
//	Column("id") [aliased: true, qualified: false, errored: true, error: alias conflict]
func (c *Column) SetErrorWith(expr string, err error) *Column {
	c.BaseToken.SetErrorWith(expr, err)
	return c
}

// String returns a structured diagnostic view of the column token.
//
// This method reports the internal state of the token using its base name,
// and flags for aliasing, qualification, and error presence.
// It is intended for developer inspection and debugging only.
//
// Unlike Raw(), this output is not suitable for SQL rendering.
//
// # Example Output:
//
//	Column("id") [aliased: false, qualified: false, errored: false]
//	Column("id") [aliased: true, qualified: false, errored: true]
//	Column("m3CUNO") [aliased: true, qualified: true, errored: false]
func (c *Column) String() string {
	s := fmt.Sprintf("Column(%q) [aliased: %v, qualified: %v, errored: %v",
		c.Name, c.IsAliased(), c.IsQualified(), c.HasError(),
	)
	if c.HasError() {
		s += fmt.Sprintf(", error: %s", c.Error.Error())
	}
	s += "]"
	return s
}

// WithTable attaches a full table token to the column.
//
// This enables alias-aware rendering by giving the column access
// to both the table name and its alias, if present. It also triggers
// column resolution, ensuring consistency between column qualification
// and table context.
//
// # Example
//
//	users := token.NewTable("users AS u")
//	col := token.NewColumn("id", "user_id").WithTable(users)
//	fmt.Println(col.Raw()) // "u.user_id"
func (c *Column) WithTable(table *Table) *Column {
	if c == nil || !c.IsValid() || c.HasError() {
		return c
	}

	// validate only if the column’s table name does not match EITHER name or alias
	if c.Table != nil && table != nil {
		if c.Table.Name != table.Name && c.Table.Name != table.Alias {
			c.BaseToken.Error = fmt.Errorf("table mismatch: column refers to %q, but table is %q (alias: %q)", c.Table.Name, table.Name, table.Alias)
			return c
		}
	}
	if table != nil {
		c.Table = table
	}

	return c
}

// Ensure Column satisfies the GenericToken interface.
var _ GenericToken = &Column{}
