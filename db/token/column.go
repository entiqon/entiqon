// File: db/token/column.go

// Package token defines SQL token types used by the query builder,
// such as columns, tables, conditions, joins, etc.
//
// Tokens encapsulate parsed components of SQL statements,
// holding both raw user input and normalized representations.
// This enables structured, safe, and extensible SQL generation.
//
// The Column token represents a column or expression in a SELECT clause,
// including expression, alias, and raw SQL flag.
package token

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/contract"
)

// Column represents a column or expression in a SELECT clause.
//
// It holds the original user input, parsed SQL expression without alias,
// optional alias, a flag indicating if the expression is raw SQL, and
// an error field indicating validation or parsing errors.
type Column struct {
	// Input is the raw user input string that generated this token.
	// It is retained for debugging, error reporting, or regeneration.
	Input string

	// Expr is the SQL expression or column name without the alias.
	Expr string

	// Alias is an optional alias for the column.
	// If set, the generated SQL will include an "AS Alias" clause.
	Alias string

	// IsRaw indicates whether Expr should be treated as raw SQL
	// and thus not quoted by the dialect.
	IsRaw bool

	// Error holds any validation or parsing error encountered.
	// It is nil if the column is considered valid.
	Error error
}

// IsAliased returns true if the column has a non-empty alias.
func (c *Column) IsAliased() bool {
	return strings.TrimSpace(c.Alias) != ""
}

// IsErrored returns true if the column has a non-nil error.
func (c *Column) IsErrored() bool {
	return c.Error != nil
}

// IsValid returns true if the column is considered valid.
//
// A column is invalid if it has an error, if its Expr is empty or whitespace,
// or if the derived Name is empty.
func (c *Column) IsValid() bool {
	if c.IsErrored() {
		return false
	}
	if strings.TrimSpace(c.Expr) == "" {
		return false
	}
	if c.Name() == "" {
		return false
	}
	return true
}

// Name returns the unique identifier for the column.
//
// It returns the lowercase Alias if set, otherwise derives a normalized
// name from Expr by removing non-alphanumeric characters and lowercasing.
func (c *Column) Name() string {
	if c.Alias != "" {
		return strings.ToLower(c.Alias)
	}
	return deriveNameFromExpr(c.Expr)
}

// Raw returns the raw SQL representation of the column.
//
// If an alias is present, it returns "Expr AS Alias".
// Otherwise, it returns Expr as is.
func (c *Column) Raw() string {
	if c.Alias != "" {
		return fmt.Sprintf("%s AS %s", c.Expr, c.Alias)
	}
	return c.Expr
}

// Render returns the raw SQL snippet for the column.
// It does not apply any dialect quoting. If an alias is present,
// it renders "Expr AS Alias"; otherwise it returns Expr as-is.
func (c *Column) Render() string {
	if c.IsAliased() {
		return fmt.Sprintf("%s AS %s", c.Expr, strings.TrimSpace(c.Alias))
	}
	return c.Expr
}

// String returns a human-readable string representation of the Column.
//
// It includes a status icon (✅ for valid, ⛔️ for errored), the Name,
// alias presence, and error state. If errored, the error message is appended.
func (c *Column) String() string {
	aliased := c.IsAliased()
	errored := !c.IsValid()

	icon := "✅"
	if errored {
		icon = "⛔️"
	}

	base := fmt.Sprintf("%s Column(%q) [alias: %t, errored: %t]", icon, c.Name(), aliased, errored)
	if errored && c.Error != nil {
		return fmt.Sprintf("%s – %s", base, c.Error.Error())
	}
	return base
}

// deriveNameFromExpr derives a normalized identifier from an expression.
//
// It removes all non-alphanumeric characters and converts to lowercase.
// This provides a stable name when no alias is provided.
func deriveNameFromExpr(expr string) string {
	var b strings.Builder
	for _, r := range expr {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return strings.ToLower(b.String())
}

// Ensure Column implements contract.Renderable at compile time.
var _ contract.Renderable = (*Column)(nil)
