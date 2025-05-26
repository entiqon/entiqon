package token

import (
	"fmt"
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
//	fmt.Println(c.Raw())     // "users.id"
//	fmt.Println(c.IsAliased()) // true
//	fmt.Println(c.String())  // Column("users.id") [aliased: true, qualified: true]
//
// File: internal/build/token/column.go
// Since: v1.6.0
type Column struct {
	BaseToken

	// TableName holds the optional table or alias prefix for the column.
	// It is extracted from expressions like "users.id" or may be set later via WithTable().
	TableName string
}

// NewColumn constructs a Column token from a raw expression and optional alias.
//
// If no alias is provided, it attempts to extract one inline using the SQL keyword "AS",
// or through space/comma separation. If both an inline alias and an explicit alias
// are provided and they conflict, an error is stored in the token.
//
// # Examples
//
//	NewColumn("id")                    → name: "id"
//	NewColumn("id AS uid")             → name: "id", alias: "uid"
//	NewColumn("users.id")              → name: "id", table: "users"
//	NewColumn("id", "alias")           → name: "id", alias: "alias"
//	NewColumn("id AS uid", "alias")    → alias mismatch error
func NewColumn(expr string, alias ...string) Column {
	base, parsedAlias := ParseAlias(expr)
	table, column := ParseTableColumn(base)

	var finalAlias string
	var err error
	if len(alias) > 0 {
		finalAlias = alias[0]
		if parsedAlias != "" && parsedAlias != finalAlias {
			err = fmt.Errorf("alias mismatch: inline alias '%s' ≠ provided alias '%s'", parsedAlias, finalAlias)
		}
	} else {
		finalAlias = parsedAlias
	}

	return Column{
		TableName: table,
		BaseToken: BaseToken{
			Name:  column,
			Alias: finalAlias,
			Error: err,
		},
	}
}

// IsQualified reports whether the column has a table prefix.
//
// This is useful for distinguishing between "id" and "users.id" in query formatting.
func (c Column) IsQualified() bool {
	return c.TableName != ""
}

// WithTable assigns or checks the table prefix for the column.
//
// If the column was already parsed with a table name and a different table
// is provided here, an error is recorded and the original table is retained.
//
// # Example
//
//	c := NewColumn("users.id")
//	c = c.WithTable("orders") // sets error due to mismatch
func (c Column) WithTable(name string) Column {
	if c.TableName != "" && c.TableName != name {
		c.Error = fmt.Errorf("table mismatch: cannot override '%s' with '%s'", c.TableName, name)
		return c
	}
	c.TableName = name
	return c
}

// Raw returns the SQL-safe, dialect-neutral representation of the column,
// including table qualification and aliasing if applicable.
//
// If the column has a table prefix, it will be prepended (e.g., "users.id").
// If the column is aliased, the result will follow SQL's "AS" convention
// (e.g., "users.id AS uid").
//
// # Examples
//
//	Column{Name: "id"}                            → "id"
//	Column{Name: "id", Alias: "uid"}              → "id AS uid"
//	Column{TableName: "users", Name: "id"}        → "users.id"
//	Column{TableName: "users", Name: "id", Alias: "uid"} → "users.id AS uid"
func (c Column) Raw() string {
	base := c.Name
	if c.IsQualified() {
		base = c.TableName + "." + c.Name
	}
	if c.IsAliased() {
		return fmt.Sprintf("%s AS %s", base, c.Alias)
	}
	return base
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
func (c Column) String() string {
	s := fmt.Sprintf("Column(%q) [aliased: %v, qualified: %v, errored: %v",
		c.Name, c.IsAliased(), c.IsQualified(), c.HasError(),
	)
	if c.HasError() {
		s += fmt.Sprintf(", error: %s", c.Error.Error())
	}
	s += "]"
	return s
}

// Ensure Column satisfies the GenericToken interface.
var _ GenericToken = Column{}
