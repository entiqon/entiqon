// File: db/token/field.go

package field

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/entiqon/entiqon/db/token"
)

// field represents a column/field or expression in a SELECT clause.
//
// It holds the original user input, a parsed SQL expression (without alias),
// an optional alias, a flag indicating whether the expression should be
// treated as raw SQL (and therefore not quoted by any dialect), and a
// possible error captured during parsing/validation.
type field struct {
	kind token.ExpressionKind

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
// A field represents a column, identifier, or computed expression in a
// SELECT clause. New enforces strict validation of its inputs:
//
//   - Always returns a non-nil *field. Even on failure, the *field carries
//     an error that can be checked with IsValid() or IsErrored().
//
//   - The first argument must be a string expression. Passing another field
//     or *field results in an error, with guidance to use Clone() instead.
//
//   - Empty input is rejected:
//     New()                → error "empty input is not allowed"
//     New("")              → error "empty expression is not allowed"
//
//   - Two arguments (expr, alias):
//     New("id", "user_id") → valid
//     New("id", "")        → error "empty alias is not allowed"
//
//   - Three arguments (expr, alias, isRaw):
//     New("id", "alias", true)   → valid
//     New("id", "alias", "true") → error "isRaw: invalid type string, expected bool"
//     New("id", "", true)        → error "empty alias is not allowed"
//
//   - One argument (expr):
//
//   - Identifiers: treated as plain fields (e.g., "id", "users.id").
//
//   - Computed expressions (functions, arithmetic, subqueries) are supported.
//     Aliases may be specified either with explicit "AS" or with shorthand:
//
//     New("SUM(price) AS total")     → valid, alias="total"
//     New("SUM(price) total")        → valid (shorthand), alias="total"
//     New("(SELECT 1) AS one")       → valid subquery, alias="one"
//     New("(SELECT 1) one")          → valid shorthand subquery, alias="one"
//
//     If a computed expression has no alias, New auto-generates one:
//
//     New("SUM(price)") → alias="raw_expr_ab12cd" (deterministic hash)
//
//   - Invalid arity (more than 3 args) is rejected with:
//     error "invalid field constructor signature"
//
// Example usage:
//
//	f1 := field.New("id")
//	f2 := field.New("id", "user_id")
//	f3 := field.New("COUNT(*) AS total")
//	f4 := field.New("name", "username", false)
//	f5 := field.New("(SELECT MAX(price) FROM sales) total_sales")
//
// Each call produces a *field that preserves the original input for auditing,
// enforces strict validation, and guarantees non-nil returns for safe chaining.
//
// Errors never cause panics; instead, they are stored inside the *field
// and exposed via IsErrored(), Error(), String(), or Debug().
func New(input ...any) Token {
	f := &field{
		owner: nil,
		input: fmt.Sprint(input...),
	}
	if len(input) == 0 {
		f.SetError(errors.New("empty input is not allowed"))
		return f
	}

	// Validate expr/input type
	if err := validateType(input[0]); err != nil {
		if err.Error() == "unsupported type: field" {
			f.SetError(fmt.Errorf("%s; if you want to create a copy, use Clone() instead", err.Error()))
		} else {
			f.SetError(fmt.Errorf("expr has %v", err))
		}
		return f
	}

	f.expr = strings.TrimSpace(input[0].(string))
	if f.expr == "" {
		f.SetError(errors.New("empty expression is not allowed"))
		return f
	}

	// --- special handling for "*" ---
	if f.expr == "*" {
		if len(input) == 1 {
			f.alias = ""
			f.isRaw = false
			return f
		}
		f.SetError(errors.New("'*' cannot be aliased or raw"))
		return f
	}

	switch len(input) {
	case 3:
		f.input = fmt.Sprint(input[0], " ", input[1], " ", input[2])
		// expr, alias, isRaw
		if err := validateType(input[1]); err != nil {
			f.SetError(fmt.Errorf("alias has %v", err))
			return f
		}
		f.alias = strings.TrimSpace(input[1].(string))
		if f.alias == "" {
			f.SetError(errors.New("empty alias is not allowed"))
			return f
		}

		isRaw, ok := input[2].(bool)
		if !ok {
			f.SetError(fmt.Errorf("isRaw has invalid type: %T, expected bool", input[2]))
		}
		f.isRaw = isRaw
		return f

	case 2:
		f.input = fmt.Sprint(input[0], " ", input[1])
		// expr, alias
		if err := validateType(input[1]); err != nil {
			f.SetError(fmt.Errorf("alias has %v", err))
			return f
		}
		f.alias = strings.TrimSpace(input[1].(string))
		if f.alias == "" {
			f.SetError(errors.New("empty alias is not allowed"))
			return f
		}
		f.isRaw = isRawExpr(f.expr)
		return f

	case 1:
		f.isRaw = isRawExpr(f.expr)
		parsedExpr, parsedAlias, _ := parseAlias(f.expr)
		f.expr = parsedExpr
		f.alias = parsedAlias
		if f.isRaw {
			// Adapted behavior: accept trailing alias without AS
			if token.HasTrailingAliasWithoutAS(f.input) {
				parsedExpr, parsedAlias, _ := parseAlias(f.input)
				f.expr = parsedExpr
				f.alias = parsedAlias
				return f
			}

			if parsedAlias == "" {
				// No alias at all → auto-generate one
				f.alias = autoAlias(f.expr)
			}
			return f
		}
		return f
	}
	f.SetError(errors.New("invalid field constructor signature"))
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
	// Reuse the standard constructor.
	f := New(input...)

	if owner == "" {
		f.SetError(errors.New("owner is empty"))
		return f
	}

	// Attach the table.
	f.SetOwner(&owner)
	return f
}

func (f *field) ExpressionKind() token.ExpressionKind {
	return f.kind
}

// HasOwner returns the owning table name or alias if set.
func (f *field) HasOwner() bool { return f.owner != nil && *f.owner != "" }

// Owner returns the owning table name or alias if set.
// If none is set, returns a pointer to the empty string.
func (f *field) Owner() *string {
	return f.owner
}

// SetOwner assigns or clears the owning table name or alias.
// Passing nil resets the owner to an empty string.
func (f *field) SetOwner(owner *string) {
	f.owner = owner
}

// Input returns the original raw input string.
func (f *field) Input() string { return f.input }

// Expr returns the parsed SQL expression without alias.
func (f *field) Expr() string { return f.expr }

// Alias returns the optional alias.
func (f *field) Alias() string { return f.alias }

// IsAliased reports whether the field has a non-empty alias.
func (f *field) IsAliased() bool {
	return strings.TrimSpace(f.alias) != ""
}

// IsValid reports whether the field is considered valid.
//
// A field is invalid if it carries an error, or Expr() is empty/whitespace,
func (f *field) IsValid() bool {
	if f.err != nil {
		return false
	}
	return true
}

// Clone returns a semantic copy of the field.
// Errors and all state are preserved.
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
//	⛔️ field("false"): [raw: false, aliased: false, errored: true] – input type unsupported: bool
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

// IsErrored reports whether the field carries a non-nil error.
func (f *field) IsErrored() bool {
	return f.err != nil
}

// SetError assigns an error to the field. Intended for use during
// construction/parsing to capture validation failures.
func (f *field) SetError(err error) Token {
	f.err = err
	return f
}

// IsRaw reports whether the field was explicitly constructed as raw
// (via the two-argument form or as a subquery).
func (f *field) IsRaw() bool { return f.isRaw }

// Raw returns the raw SQL snippet for the field without dialect quoting.
//
// If an alias is present, it returns "Expr AS Alias"; otherwise it returns Expr.
// Errors do not suppress output — callers should check IsErrored() explicitly.
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

// Render implements contract.Renderable by returning the raw SQL snippet
// for the field (no dialect quoting). If an alias is present, it renders
// "Expr AS Alias"; otherwise it returns Expr.
//
// This method is stable and machine-facing, suitable for builders.
// Dialect-specific quoting may later override Raw vs Render.
func (f *field) Render() string {
	return f.Raw()
}

// String implements fmt.Stringer and returns a concise,
// user-friendly representation suitable for UI/UX display.
//
// Example outputs:
//
//	✅ field("users.id AS user_id")
//	✅ field("COUNT(*) AS total")
//	⛔ field("SUM()"): empty expression is not allowed
//
// For developer diagnostics (flags, error state), use Debug().
// For SQL output, use Render() or Raw().
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

func autoAlias(expr string) string {
	const prefix = "expr_alias_"
	h := sha1.New()
	h.Write([]byte(expr))
	sum := hex.EncodeToString(h.Sum(nil))
	return prefix + sum[:6] // 6 hex chars => ~16.7 million combinations
}

// isRawExpr detects whether the expression contains raw SQL indicators.
func isRawExpr(expr string) bool {
	expr = strings.TrimSpace(expr)

	// Parentheses usually indicate a raw expression (function call, grouping)
	if strings.ContainsAny(expr, "()") {
		return true
	}

	// Common SQL arithmetic or concatenation operators
	operators := []string{"||", "+", "-", "*", "/"}
	for _, op := range operators {
		if strings.Contains(expr, op) {
			return true
		}
	}

	return false
}

// parseAlias attempts to split an expression into [expr, alias].
//
// Priority:
//  1. Explicit "AS" (case-insensitive).
//  2. Subquery shorthand: "(SELECT ... ) alias" → expr="(SELECT ...)", alias="alias".
//  3. Function-style shorthand: FUNC(...), may allow alias after the closing ")".
//  4. Operator chains with no wrapper (col1+col2, col1||col2, etc.) are treated as
//     full expressions with no alias unless explicit AS is used.
//  5. Fallback: last token is alias, everything before it is expr.
func parseAlias(s string) (expr, alias string, fromAS bool) {
	s = strings.TrimSpace(s)

	// 1. Explicit AS first (case-insensitive).
	re := regexp.MustCompile(`(?i)\s+AS\s+`)
	parts := re.Split(s, 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
	}

	// 2. Subquery shorthand: "(... ) alias"
	if strings.HasPrefix(s, "(") {
		closing := strings.LastIndex(s, ")")
		if closing > 0 && closing < len(s)-1 {
			expr := strings.TrimSpace(s[:closing+1])
			alias := strings.TrimSpace(s[closing+1:])
			return expr, alias, false
		}
	}

	// 3. Function-style shorthand (e.g. SUM(...))
	if strings.Contains(s, "(") && strings.HasSuffix(s, ")") {
		// no alias → return whole string as expr
		return s, "", false
	}

	// 4. Operator-heavy expressions.
	// If starts with a function name, allow alias parsing.
	if strings.ContainsAny(s, "+-*/") || strings.Contains(s, "||") {
		// crude detection: function call starts with a word + "("
		if matched, _ := regexp.MatchString(`^[A-Za-z_][A-Za-z0-9_]*\s*\(`, s); !matched {
			// not a function call → treat as pure operator expr, no alias
			return s, "", false
		}
	}

	// 5. Fallback: last token = alias, everything before = expr.
	tokens := strings.Fields(s)
	if len(tokens) >= 2 {
		expr := strings.Join(tokens[:len(tokens)-1], " ")
		alias := tokens[len(tokens)-1]
		return expr, alias, false
	}

	// Default: no alias
	return s, "", false
}

func validateType(input any) error {
	switch v := input.(type) {
	case field, *field:
		return errors.New("unsupported type: field")
	case string:
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}
