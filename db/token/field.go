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
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/entiqon/entiqon/db/contract"
)

// Ensure Field implements contract.Renderable at compile time.
var _ contract.Renderable = (*Field)(nil)

// Ensure Field implements contract.Cloanable[*Field] at compile time.
var _ contract.Cloanable[*Field] = (*Field)(nil)

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
func NewField(inputs ...any) *Field {
	if len(inputs) == 0 {
		return &Field{}
	}

	// Validate expr/input type
	if err := validateType(inputs[0]); err != nil {
		fd := &Field{
			Input: fmt.Sprint(inputs[0]),
		}
		if err.Error() == "input is a Field" {
			fd.setError(fmt.Errorf("%s; if you want to create a copy, use Clone() instead", err.Error()))
		} else {
			fd.setError(err)
		}
		return fd
	}
	expr := strings.TrimSpace(inputs[0].(string))

	switch len(inputs) {
	case 3:
		// input, alias, isRawExpr
		if err := validateType(inputs[1]); err != nil {
			return &Field{
				Input: expr,
				Error: fmt.Errorf("%s: %s", err.Error(), "alias must be a string"),
			}
		}
		alias := strings.TrimSpace(inputs[1].(string))

		isRaw, ok := inputs[2].(bool)
		if !ok {
			fd := &Field{
				Input: expr,
			}
			fd.setError(errors.New("isRaw must be a bool"))
			return fd
		}

		return &Field{
			Input: expr,
			Expr:  expr,
			Alias: alias,
			IsRaw: isRaw,
		}

	case 2:
		// expr, alias
		if err := validateType(inputs[1]); err != nil {
			return &Field{
				Input: expr,
				Error: err,
			}
		}
		alias := strings.TrimSpace(inputs[1].(string))
		return &Field{
			Input: expr,
			Expr:  expr,
			Alias: alias,
			IsRaw: isRawExpr(expr),
		}

	case 1:
		isRaw := isRawExpr(expr)

		if isRaw {
			parsedExpr := expr // always preserve full expression

			// Explicit AS → parse alias
			if strings.Contains(strings.ToUpper(expr), " AS ") {
				_, parsedAlias, _ := parseAlias(expr)
				return &Field{
					Input: expr,
					Expr:  parsedExpr,
					Alias: parsedAlias,
					IsRaw: true,
				}
			}

			// Has space but no AS → error (alias without AS is invalid for raw)
			if HasTrailingAliasWithoutAS(expr) {
				fd := &Field{Input: parsedExpr}
				fd.setError(errors.New("raw expressions must use explicit AS for alias"))
				return fd
			}

			// No alias at all → auto-generate
			auto := autoAlias(parsedExpr)
			return &Field{
				Input: fmt.Sprintf("%s AS %s", parsedExpr, auto),
				Expr:  parsedExpr,
				Alias: auto,
				IsRaw: true,
			}
		}

		// Not raw — safe to parse normally
		parsedExpr, parsedAlias, _ := parseAlias(expr)
		return &Field{
			Input: expr,
			Expr:  parsedExpr,
			Alias: parsedAlias,
			IsRaw: false,
		}
	}

	fd := &Field{}
	fd.setError(errors.New("invalid NewField signature"))
	return fd
}

// Clone returns a semantic copy of the Field. A nil receiver yields nil.
// This nil-preserving behavior makes Clone safe to call unconditionally.
func (f *Field) Clone() *Field {
	if f == nil {
		return nil
	}
	cp := *f
	return &cp
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
// If Alias is set, it returns the lower-cased alias. Otherwise, it derives
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

	base := fmt.Sprintf(
		"%s Field(%q): [raw: %t, aliased: %t, errored: %t]",
		icon,
		f.Name(),
		f.IsRaw,
		aliased,
		errored,
	)
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

func autoAlias(expr string) string {
	const prefix = "raw_expr_"
	h := sha1.New()
	h.Write([]byte(expr))
	sum := hex.EncodeToString(h.Sum(nil))
	return prefix + sum[:6] // 6 hex chars => ~16.7 million combinations
}

// HasTrailingAliasWithoutAS checks if the last space-separated token is an alias candidate
func HasTrailingAliasWithoutAS(expr string) bool {
	up := strings.ToUpper(expr)
	if strings.Contains(up, " AS ") {
		return false // explicit AS → fine
	}

	tokens := strings.Fields(expr)
	if len(tokens) <= 1 {
		return false // single token can't have alias
	}

	last := tokens[len(tokens)-1]
	penultimate := tokens[len(tokens)-2]

	// If the token before last is an operator, this "last" is part of the expression, not alias
	operators := map[string]bool{"||": true, "+": true, "-": true, "*": true, "/": true}
	if operators[penultimate] {
		return false
	}

	// Otherwise, if it looks like an identifier → treat as alias
	if regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(last) {
		return true
	}
	return false
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

func parseAlias(s string) (expr, alias string, fromAS bool) {
	s = strings.TrimSpace(s)

	// Try explicit AS first (case-insensitive)
	re := regexp.MustCompile(`(?i)\s+AS\s+`)
	parts := re.Split(s, 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
	}

	// Fallback: split by first space
	parts = strings.Fields(s)
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), false
	}

	return s, "", false
}

func validateType(input any) error {
	switch v := input.(type) {
	case Field, *Field:
		return errors.New("input is a Field")
	case string:
		return nil
	default:
		return fmt.Errorf("input type unsupported: %T", v)
	}
}
