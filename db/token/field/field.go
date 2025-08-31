// File: db/token/field.go

package field

import (
	stdErr "errors"
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/errors"
	"github.com/entiqon/entiqon/db/token/helpers"
	"github.com/entiqon/entiqon/db/token/types/identifier"
)

// field represents a column or expression in a SELECT clause.
//
// It holds the original user input, a parsed SQL expression (without alias),
// an optional alias, a flag indicating whether the expression should be
// treated as raw SQL (and therefore not quoted by any dialect), and a
// possible error captured during parsing/validation.
type field struct {
	kind identifier.Type

	// owner returns the optional table associated with this field.
	// If no table was set, it returns nil.
	owner *string

	// input is the raw user input string that produced this token.
	// It is retained for debugging, error reporting, or regeneration.
	input string

	// expr is the SQL expression or field name without the alias.
	expr string

	// alias is an optional alias for the field. If set, generated SQL
	// should include an "AS Alias" clause.
	alias string

	// isRaw indicates whether Expr should be treated as raw SQL
	// and thus not quoted by a dialect.
	isRaw bool

	// err holds any validation or parsing error encountered.
	// It is nil when the field is considered valid.
	err error
}

// New constructs a *field token from the given arguments.
//
// A field represents a column, identifier, literal, or computed expression
// in a SELECT clause. New enforces strict validation of its inputs:
//
//   - Always returns a non-nil *field. Even on failure, the *field carries
//     an error that can be checked with IsValid() or IsErrored().
//
//   - The first argument must be a string expression. Passing another field
//     or *field results in an error, with guidance to use Clone() instead.
//
//   - Empty input is rejected:
//     New()                → error "empty input is not allowed"
//     New("")              → error "empty identifier: \"\""
//
//   - One argument (expr):
//     New("id")                → Identifier, expr="id"
//     New("id user_id")        → Identifier, expr="id", alias="user_id"
//     New("id AS user_id")     → Identifier, expr="id", alias="user_id"
//     New("SUM(price) AS t")   → Aggregate, expr="SUM(price)", alias="t"
//     New("(SELECT 1) one")    → Subquery, expr="(SELECT 1)", alias="one"
//     New("'a' msg")           → Literal, expr="'a'", alias="msg"
//
//   - Two arguments (expr, alias):
//     New("id", "user_id")     → Identifier with explicit alias
//     New("*", "alias")        → error "'*' cannot be aliased or raw"
//     New("id", 123)           → error "alias must be a string, got int"
//     New("id", "123alias")    → error "invalid alias identifier cannot start with digit"
//
//   - Passing an existing Field:
//     New(field.New("id"))     → error "unsupported type; use Clone() instead"
//
// Validation is layered:
//  1. Type validation — only strings are accepted; other tokens must be
//     cloned, and non-strings are rejected with descriptive errors.
//  2. Classification — expressions are categorized via ResolveExpressionType.
//  3. Resolution — expr and alias are extracted according to category.
//
// Wildcards ("*") are only valid without an alias; attempts to alias them
// result in errors.
//
// Example usage:
//
//	f1 := field.New("id")
//	f2 := field.New("id", "user_id")
//	f3 := field.New("COUNT(*) AS total")
//	f4 := field.New("(SELECT MAX(price) FROM sales) total_sales")
//
// Each call produces a *field that preserves the original input for auditing,
// enforces strict validation, and guarantees non-nil returns for safe chaining.
//
// Errors never cause panics; instead, they are stored inside the *field
// and exposed via IsErrored(), Error(), String(), or Debug().
func New(input ...any) Token {
	f := &field{
		kind:  identifier.Invalid,
		owner: nil,
		input: "",
	}

	if len(input) == 0 {
		return f.SetError(stdErr.New("empty input is not allowed"))
	}

	if len(input) > 2 {
		return f.SetError(fmt.Errorf("invalid field constructor signature: %d args", len(input)))
	}

	// Type validation (string only)
	if err := helpers.ValidateType(input[0]); err != nil {
		if stdErr.Is(err, errors.UnsupportedTypeError) {
			return f.SetError(fmt.Errorf(
				"%w; if you want to create a copy, use Clone() instead", err,
			))
		}
		return f.SetError(fmt.Errorf("expr has %w", err))
	}

	f.input = strings.Join(helpers.Stringify(input), " ") // keep audit trail

	// always resolve the first part
	kind, expr, alias, err := helpers.ResolveExpression(fmt.Sprint(input[0]), true)
	if err != nil {
		return f.SetError(err)
	}
	f.kind, f.expr, f.alias = kind, expr, alias

	if len(input) == 2 {
		a, ok := input[1].(string)
		if !ok {
			return f.SetError(fmt.Errorf("alias must be a string, got %T", input[1]))
		}
		if err := helpers.ValidateAlias(a); err != nil {
			return f.SetError(err)
		}
		f.alias = a
	}

	if err := helpers.ValidateWildcard(f.expr, f.alias); err != nil {
		return f.SetError(err)
	}

	return f
}

// NewWithTable constructs a field bound to a specific table.
//
// Example:
//
//	u := table.New("users")
//	f := field.NewWithTable(u, "id", "user_id")
//	// Renders as: users.id AS user_id
func NewWithTable(owner string, input ...any) Token {
	f := New(input...)

	if owner == "" {
		f.SetError(stdErr.New("owner is empty"))
		return f
	}

	f.SetOwner(&owner)
	return f
}

// ExpressionKind returns the identifier.Type classification of the field.
func (f *field) ExpressionKind() identifier.Type { return f.kind }

// HasOwner reports whether the field has a table owner assigned.
func (f *field) HasOwner() bool { return f.owner != nil && *f.owner != "" }

// Owner returns the owning table name or alias if set.
func (f *field) Owner() *string { return f.owner }

// SetOwner assigns or clears the owning table name or alias.
// Passing nil resets the owner.
func (f *field) SetOwner(owner *string) { f.owner = owner }

// Input returns the original raw input string provided to the constructor.
func (f *field) Input() string { return f.input }

// Expr returns the resolved SQL expression without alias.
func (f *field) Expr() string { return f.expr }

// Alias returns the optional alias.
func (f *field) Alias() string { return f.alias }

// IsAliased reports whether the field has a non-empty alias.
func (f *field) IsAliased() bool { return strings.TrimSpace(f.alias) != "" }

// Clone returns a semantic copy of the field, preserving all state and errors.
func (f *field) Clone() Token {
	cp := f
	if f.owner != nil {
		owner := *f.owner
		cp.owner = &owner
	}
	return cp
}

// Debug returns a compact diagnostic view of the field.
//
// Example (valid):
//
//	✅ field("COUNT(*) AS total"): [raw: true, aliased: true, errored: false]
//
// Example (invalid):
//
//	⛔ field("false"): [raw: false, aliased: false, errored: true] – expr has invalid format
func (f *field) Debug() string {
	flags := fmt.Sprintf(
		"[raw: %v, aliased: %v, errored: %v]",
		f.isRaw,
		f.alias != "",
		f.err != nil,
	)

	if f.err != nil {
		return fmt.Sprintf("⛔️ field(%q): %s – %v", f.input, flags, f.err)
	}
	return fmt.Sprintf("✅ field(%q): %s", f.input, flags)
}

// Error returns the underlying construction error, if any.
func (f *field) Error() error { return f.err }

// IsErrored reports whether the field is invalid (kind=Invalid or err non-nil).
func (f *field) IsErrored() bool {
	return f.ExpressionKind() == identifier.Invalid || f.err != nil
}

// SetError assigns an error to the field and returns itself.
func (f *field) SetError(err error) Token {
	f.err = err
	return f
}

// IsRaw reports whether the field should be rendered as raw SQL
// (subquery, computed, function, aggregate).
func (f *field) IsRaw() bool {
	switch f.kind {
	case identifier.Subquery, identifier.Computed, identifier.Function, identifier.Aggregate:
		return true
	default:
		return false
	}
}

// Raw returns the SQL snippet for the field without dialect quoting.
// If an alias is present, it returns "Expr AS Alias".
// Errors do not suppress output — callers should check IsErrored().
func (f *field) Raw() string {
	base := f.expr
	if f.alias != "" {
		base = fmt.Sprintf("%s AS %s", base, f.alias)
	}
	if f.owner != nil && !f.isRaw {
		base = fmt.Sprintf("%s.%s", *f.owner, base)
	}
	return base
}

// Render implements contract.Renderable by returning Raw().
func (f *field) Render() string { return f.Raw() }

// String implements fmt.Stringer, returning a user-friendly view of the field.
// For developer diagnostics, use Debug(). For SQL, use Render()/Raw().
func (f *field) String() string {
	base := f.expr
	if f.alias != "" {
		base = fmt.Sprintf("%s AS %s", base, f.alias)
	}
	if f.owner != nil && !f.isRaw {
		base = fmt.Sprintf("%s.%s", *f.owner, base)
	}

	if f.err != nil {
		return fmt.Sprintf("⛔ field(%q): %v", base, f.err)
	}
	return fmt.Sprintf("✅ field(%q)", base)
}

// IsValid reports whether the field is valid (no errors).
func (f *field) IsValid() bool { return !f.IsErrored() }
