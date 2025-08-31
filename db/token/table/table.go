package table

import (
	stdErr "errors"
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/errors"
	"github.com/entiqon/entiqon/db/token/helpers"
	"github.com/entiqon/entiqon/db/token/types/identifier"
)

// table represents a SQL table token.
//
// A table encapsulates the user-provided input, the normalized
// base table name, and an optional alias. It provides multiple
// forms for different audiences:
//
//   - Identity: Kind(), Input(), Expr(), Name(), Alias(), IsAliased()
//   - Lifecycle: Clone(), IsErrored(), Error(), SetError(), IsValid()
//   - Representations: Raw(), Render(), String(), Debug()
//
// Internal fields:
//   - kind: identifies the token category (reserved for future use).
//   - input: the exact user-provided string, preserved verbatim.
//   - name: the normalized base table name, derived from input.
//   - alias: the optional alias (parsed or explicitly set).
//   - err: any construction or validation error.
//   - isRaw: whether the table was constructed via the two-argument
//     form or as a subquery.
type table struct {
	kind identifier.Type

	// input is the exact user-provided string, never modified.
	input string

	// name is the normalized base table name (derived from input).
	name string

	// alias is the optional alias (parsed or explicitly set).
	alias string

	// err holds any construction or validation error.
	err error

	// isRaw reports whether the table was constructed via the
	// explicit two-argument form or is a subquery.
	isRaw bool
}

// New constructs a new Table from user input.
//
// Accepted forms:
//   - table.New("users")          → name="users"
//   - table.New("users u")        → name="users", alias="u"
//   - table.New("users AS u")     → name="users", alias="u"
//   - table.New("users", "u")     → name="users", alias="u", isRaw=true
//   - table.New("(SELECT ...)", "t")
//     → name="(SELECT ...)", alias="t", isRaw=true
//   - table.New("(SELECT ...) AS t")
//     → name="(SELECT ...)", alias="t", isRaw=true
//
// The first argument is always preserved verbatim in input.
// If construction fails, the returned table is errored but
// still carries the original input for diagnostics.
func New(input ...any) Token {
	t := &table{
		kind:  identifier.Invalid,
		input: fmt.Sprint(input...),
	}

	if len(input) == 0 {
		return t.SetError(stdErr.New("empty input is not allowed"))
	}

	if len(input) > 2 {
		return t.SetError(fmt.Errorf("invalid table constructor signature: %d args", len(input)))
	}

	// Type validation (string only)
	if err := helpers.ValidateType(input[0]); err != nil {
		if stdErr.Is(err, errors.UnsupportedTypeError) {
			return t.SetError(fmt.Errorf(
				"%w; if you want to create a copy, use Clone() instead", err,
			))
		}
		return t.SetError(fmt.Errorf("expr has %w", err))
	}

	t.input = strings.Join(helpers.Stringify(input), " ") // keep audit trail

	// always resolve the first part
	kind, expr, alias, err := helpers.ResolveExpression(fmt.Sprint(input[0]), true)
	if err != nil {
		return t.SetError(err)
	}
	t.kind, t.name, t.alias = kind, expr, alias

	if len(input) == 2 {
		a, ok := input[1].(string)
		if !ok {
			return t.SetError(fmt.Errorf("alias must be a string, got %T", input[1]))
		}
		if err := helpers.ValidateAlias(a); err != nil {
			return t.SetError(err)
		}
		t.alias = a
	}

	// ✅ one place only: context rule
	if t.kind == identifier.Literal || t.kind == identifier.Aggregate {
		return t.SetError(fmt.Errorf(
			"%s %q cannot be used as a table source",
			strings.ToLower(t.kind.String()), t.name,
		))
	}

	return t
}

// ExpressionKind returns the kind of token (always table).
func (t *table) ExpressionKind() identifier.Type { return t.kind }

// Input returns the exact user-provided input string.
func (t *table) Input() string { return t.input }

// Expr returns the normalized table name (without alias).
func (t *table) Expr() string { return t.name }

// Name returns the normalized base table name.
func (t *table) Name() string { return t.name }

// Alias returns the alias of the table, if any.
func (t *table) Alias() string { return t.alias }

// IsAliased reports whether the table has an alias.
func (t *table) IsAliased() bool { return t.alias != "" }

// Clone returns a semantic copy of the table.
func (t *table) Clone() Token {
	return &table{
		input: t.input,
		name:  t.name,
		alias: t.alias,
		err:   t.err,
		isRaw: t.isRaw,
	}
}

// Debug returns a developer-facing representation of the table.
//
// The output is verbose and intended for diagnostics, showing the
// original input and internal flags (raw, aliased, errored).
// If the table is errored, Debug also appends the error message.
func (t *table) Debug() string {
	flags := fmt.Sprintf(
		"[raw:%v, aliased:%v, errored:%v]",
		t.IsRaw(),
		t.IsAliased(),
		t.IsErrored(),
	)
	if !t.IsValid() {
		return fmt.Sprintf("❌ Table(%q): %s {err=%v}", t.input, flags, t.err)
	}
	return fmt.Sprintf("✅ Table(%q): %s", t.input, flags)
}

// IsErrored reports whether the table was constructed with an error.
func (t *table) IsErrored() bool { return t.err != nil }

// Error returns the underlying construction error, if any.
func (t *table) Error() error { return t.err }

// SetError assigns an error to the table and returns it.
//
// This is primarily used during parsing to capture validation failures.
func (t *table) SetError(err error) Token {
	t.err = err
	return t
}

// IsRaw reports whether the table represents a raw expression
// (subquery, explicit 2-arg form, or anything that is not a plain identifier).
func (t *table) IsRaw() bool {
	switch t.kind {
	case identifier.Subquery, identifier.Computed, identifier.Function, identifier.Aggregate:
		return true
	default:
		return false
	}
}

// Raw returns the generic SQL fragment of the table, including alias if present.
//
// Unlike Render(), Raw() does not apply dialect-specific quoting or rewriting.
// It simply reflects the normalized fragment as a plain SQL string.
func (t *table) Raw() string {
	if !t.IsValid() {
		return ""
	}
	if t.alias != "" {
		return fmt.Sprintf("%s AS %s", t.name, t.alias)
	}
	return t.name
}

// Render returns the canonical SQL representation of the table.
//
// Render() differs from Raw() in that it represents the resolved form
// the builder will actually use. It may later be extended to apply
// dialect-specific quoting or rewriting.
//
// If the table is invalid or errored, Render() returns an empty string.
func (t *table) Render() string {
	if !t.IsValid() {
		return ""
	}
	if t.alias != "" {
		return fmt.Sprintf("%s AS %s", t.name, t.alias)
	}
	return t.name
}

// String returns the human-facing representation of the table.
//
// Unlike Render(), which is used in query building, String() is
// intended for logs and audits. It produces a concise summary of
// the table state.
//
// If the table is invalid or errored, String() reports the user input
// and the associated error. If valid, it reports the normalized form
// with alias if present, prefixed with a ✅ marker.
func (t *table) String() string {
	if !t.IsValid() {
		// Always show original input and error
		return fmt.Sprintf("❌ Table(%q): %v", t.input, t.err)
	}
	if t.alias != "" {
		return fmt.Sprintf("✅ Table(%s AS %s)", t.name, t.alias)
	}
	return fmt.Sprintf("✅ Table(%s)", t.name)
}

// IsValid reports whether the table has a non-empty name and no error.
func (t *table) IsValid() bool { return !t.IsErrored() && t.name != "" }
