// File: db/token/field.go

// Package token defines SQL token types used by the query builder,
// such as fields (columns), tables, conditions, joins, and more.
//
// Tokens encapsulate parsed components of SQL statements, preserving both
// raw user input and normalized representations to enable structured,
// safe, and extensible SQL generation.
//
// Field represents a column/field or expression within a SELECT list,
// including its expression, optional alias, and flags relevant to rendering.
package token

import (
	"errors"
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/contract"
)

// Ensure Field implements contract.Renderable at compile time.
var _ contract.Renderable = (*Field)(nil)

// Field represents a column/field or expression in a SELECT clause.
//
// It holds the original user input, a parsed SQL expression (without alias),
// an optional alias, a flag indicating whether the expression should be
// treated as raw SQL (and therefore not quoted by any dialect), and a
// possible error captured during parsing/validation.
type Field struct {
	// Input is the raw user input string that produced this token.
	// It is retained for debugging, error reporting, or regeneration.
	Input string

	// Expr is the SQL expression or field name without the alias.
	Expr string

	// Alias is an optional alias for the field. If set, generated SQL
	// should include an "AS Alias" clause.
	Alias string

	// IsRaw indicates whether Expr should be treated as raw SQL
	// and thus not quoted by a dialect.
	IsRaw bool

	// Error holds any validation or parsing error encountered.
	// It is nil when the field is considered valid.
	Error error
}

// NewField constructs a Field from the given components, trimming the alias
// and initializing a validation error when the expression is empty or when
// a derived name would be empty.
func NewField(input string, expr string, alias string, isRaw bool) *Field {
	fd := &Field{
		Input: input,
		Expr:  expr,
		Alias: strings.TrimSpace(alias),
		IsRaw: isRaw,
	}
	if strings.TrimSpace(expr) == "" {
		fd.setError(errors.New("expression cannot be empty"))
	} else if fd.Name() == "" {
		fd.setError(errors.New("derived field name cannot be empty"))
	}
	return fd
}

// IsAliased reports whether the field has a non-empty alias.
func (f *Field) IsAliased() bool {
	return strings.TrimSpace(f.Alias) != ""
}

// IsErrored reports whether the field carries a non-nil error.
func (f *Field) IsErrored() bool {
	return f.Error != nil
}

// IsValid reports whether the field is considered valid.
//
// A field is invalid if it carries an error, if Expr is empty/whitespace,
// or if the derived Name would be empty.
func (f *Field) IsValid() bool {
	if f.IsErrored() {
		return false
	}
	if strings.TrimSpace(f.Expr) == "" {
		return false
	}
	if f.Name() == "" {
		return false
	}
	return true
}

// Name returns a stable identifier for the field.
//
// If Alias is set, it returns the lower-cased alias. Otherwise it derives
// a name from Expr by removing non-alphanumeric characters and lower-casing.
func (f *Field) Name() string {
	if f.Alias != "" {
		return strings.ToLower(f.Alias)
	}
	return deriveNameFromExpr(f.Expr)
}

// Raw returns the raw SQL snippet for the field without dialect quoting.
//
// If an alias is present it returns "Expr AS Alias"; otherwise it returns Expr.
func (f *Field) Raw() string {
	if f.IsAliased() {
		return fmt.Sprintf("%s AS %s", f.Expr, strings.TrimSpace(f.Alias))
	}
	return f.Expr
}

// Render implements contract.Renderable by returning the raw SQL snippet
// for the field (no dialect quoting). If an alias is present, it renders
// "Expr AS Alias"; otherwise it returns Expr.
func (f *Field) Render() string {
	return f.Raw()
}

// String returns a human-readable representation of the field, including a
// status icon (✅ when valid, ⛔️ when invalid), the resolved Name, whether
// it is aliased, and error state. If errored, the error message is appended.
func (f *Field) String() string {
	aliased := f.IsAliased()
	errored := !f.IsValid()

	icon := "✅"
	if errored {
		icon = "⛔️"
	}

	base := fmt.Sprintf("%s Field(%q) [alias: %t, errored: %t]", icon, f.Name(), aliased, errored)
	if errored && f.Error != nil {
		return base + " – " + f.Error.Error()
	}
	return base
}

// deriveNameFromExpr derives a normalized identifier from an expression by
// removing all non-alphanumeric characters and converting to lower case.
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

// setError assigns an error to the field. Intended for use during
// construction/parsing to capture validation failures.
func (f *Field) setError(err error) { f.Error = err }
