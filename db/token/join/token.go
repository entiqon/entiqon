package join

import (
	"errors"
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/token/table"
	"github.com/entiqon/entiqon/db/token/types/join"
)

// token is the unexported implementation of the Token interface.
//
// It preserves:
//   - kind      → the token type (INNER, LEFT, RIGHT, FULL)
//   - left/right → table operands
//   - condition → the ON clause
//   - err       → error state, if any
type token struct {
	kind      join.Type
	left      table.Token
	condition string
	right     table.Token
	err       error
}

// New constructs a Join with an explicit kind.
//
// This constructor is the most flexible and is intended for advanced usage.
// The kind argument may be:
//
//   - Type (e.g. token.LeftJoin)
//   - string ("LEFT", "LEFT JOIN", case-insensitive)
//
// Example:
//
//	j1 := token.New(token.InnerJoin, "users", "orders", "u.id = o.user_id")
//	j2 := token.New("LEFT", "users", "orders", "u.id = o.user_id")
//
// If kind is invalid, New returns an errored token immediately.
//
// # Guidance
//
// Use New when building dynamic or extensible query logic, such as when
// parsing configuration, building a DSL, or integrating external inputs.
// For most application code, prefer the safe wrappers:
// NewInner, NewLeft, NewRight, NewFull — these are intention-revealing
// and avoid the possibility of invalid kinds.
func New(kind any, left, right any, condition string) Token {
	return newWithKind(kind, left, right, condition)
}

// NewInner constructs an explicit INNER JOIN.
func NewInner(left, right any, condition string) Token {
	return newWithKind(join.Inner, left, right, condition)
}

// NewLeft constructs an explicit LEFT JOIN.
func NewLeft(left, right any, condition string) Token {
	return newWithKind(join.Left, left, right, condition)
}

// NewRight constructs an explicit RIGHT JOIN.
func NewRight(left, right any, condition string) Token {
	return newWithKind(join.Right, left, right, condition)
}

// NewFull constructs an explicit FULL JOIN.
func NewFull(left, right any, condition string) Token {
	return newWithKind(join.Full, left, right, condition)
}

// NewCross constructs an explicit CROSS JOIN.
func NewCross(left, right any) Token {
	return newWithKind(join.Cross, left, right, "")
}

// NewNatural constructs an explicit NATURAL JOIN.
func NewNatural(left, right any) Token {
	return newWithKind(join.Natural, left, right, "")
}

// Clone returns a deep copy of the token.
func (t *token) Clone() Token {
	return &token{
		kind:      t.kind,
		left:      t.left.Clone(),
		right:     t.right.Clone(),
		condition: t.condition,
		err:       t.err,
	}
}

// Kind returns the type of the token (INNER, LEFT, RIGHT, FULL).
func (t *token) Kind() join.Type {
	return t.kind
}

// Left returns the left table operand.
func (t *token) Left() table.Token {
	return t.left
}

// Right returns the right table operand.
func (t *token) Right() table.Token {
	return t.right
}

// Condition returns the ON condition string.
func (t *token) Condition() string {
	return t.condition
}

// Debug returns an auditable representation of the token.
// Includes validity state and error if applicable.
func (t *token) Debug() string {
	valid := t.IsValid()
	leftExpr := "<nil>"
	rightExpr := "<nil>"

	if t.left != nil {
		leftExpr = t.left.Raw()
	}
	if t.right != nil {
		rightExpr = t.right.Raw()
	}

	if !valid {
		return fmt.Sprintf(
			"Join{Type:%q, Left:%q, Right:%q, Condition:%q, Valid:false, Err:%v}",
			t.kind,
			leftExpr,
			rightExpr,
			strings.TrimSpace(t.condition),
			t.err,
		)
	}

	return fmt.Sprintf(
		"Join{Type:%q, Left:%q, Right:%q, Condition:%q, Valid:true}",
		t.kind,
		leftExpr,
		rightExpr,
		strings.TrimSpace(t.condition),
	)
}

// Error returns the current error on the token, if any.
func (t *token) Error() error {
	return t.err
}

// IsErrored reports whether the token is in an errored state.
func (t *token) IsErrored() bool {
	return t.err != nil
}

// SetError sets the error and returns the receiver as Token.
func (t *token) SetError(err error) Token {
	t.err = err
	return t
}

// IsRaw always returns false for joins.
// A Join is a SQL-generic construct, not a raw fragment.
func (t *token) IsRaw() bool {
	return false
}

// Raw returns the raw SQL fragment of the token clause.
// If the token is errored, Raw returns an empty string.
func (t *token) Raw() string {
	if t.err != nil {
		return ""
	}

	// 🔑 Special rendering for CROSS / NATURAL joins
	if t.kind == join.Cross || t.kind == join.Natural {
		return fmt.Sprintf("%s %s", t.kind, t.right.Raw())
	}

	return fmt.Sprintf("%s %s ON %s",
		t.kind,
		t.right.Raw(),
		strings.TrimSpace(t.condition),
	)
}

// Render produces the canonical SQL fragment for the token.
// It is dialect-agnostic and delegates to Raw().
func (t *token) Render() string {
	return t.Raw()
}

// String returns a concise, loggable representation of the token.
// Valid joins are marked with ✅, invalid ones with ⛔.
func (t *token) String() string {
	var base string
	if t.kind == join.Cross || t.kind == join.Natural {
		base = fmt.Sprintf("%s %s", t.kind, t.right.Raw())
	} else {
		base = fmt.Sprintf("%s %s ON %s",
			t.kind,
			func() string {
				if t.right != nil {
					return t.right.Raw()
				}
				return ""
			}(),
			strings.TrimSpace(t.condition),
		)
	}

	if !t.IsValid() {
		return fmt.Sprintf("⛔ token(%q): %v", base, t.err)
	}
	return fmt.Sprintf("✅ token(%q)", base)
}

// IsValid reports whether the token is valid (not errored).
func (t *token) IsValid() bool {
	return !t.IsErrored()
}

func newWithKind(kind any, left, right any, condition string) Token {
	jk := normalizeKind(kind)
	if !jk.IsValid() {
		// EARLY EXIT: unsupported or invalid token kind
		return (&token{}).SetError(errors.New(jk.String()))
	}

	lt := normalizeTable(left, "left")
	rt := normalizeTable(right, "right")

	j := &token{kind: jk, left: lt, right: rt, condition: condition}

	if lt == nil || rt == nil {
		return j.SetError(fmt.Errorf("token requires both left and right tables"))
	}
	if lt.IsErrored() || rt.IsErrored() {
		var errs []string
		if lt.IsErrored() {
			errs = append(errs, fmt.Sprintf("left: %v", lt.Error()))
		}
		if rt.IsErrored() {
			errs = append(errs, fmt.Sprintf("right: %v", rt.Error()))
		}
		return j.SetError(fmt.Errorf("token invalid: %s", strings.Join(errs, "; ")))
	}

	// 🔑 Special case for CROSS / NATURAL: they must NOT have conditions
	if jk == join.Cross || jk == join.Natural {
		j.condition = ""
		return j
	}

	// For all other join kinds: require condition
	if condition == "" {
		return j.SetError(fmt.Errorf("token condition is empty"))
	}

	return j
}

// normalizeKind resolves arbitrary input into a Type.
// Accepts token.Kind or string. Unsupported types normalize to invalid.
func normalizeKind(arg any) join.Type {
	if k, ok := arg.(join.Type); ok {
		return k
	}
	if s, ok := arg.(string); ok {
		return join.ParseFrom(s)
	}
	return join.Type(-1)
}

// normalizeTable resolves a token operand into a table.Token.
// Accepts table.Token, string (delegates to table.New), or nil.
// Unsupported types produce an errored table.Token.
func normalizeTable(arg any, side string) table.Token {
	if arg == nil {
		return nil
	}
	switch v := arg.(type) {
	case table.Token:
		return v
	case string:
		return table.New(v)
	default:
		t := table.New(fmt.Sprintf("invalid_%s", side))
		return t.SetError(fmt.Errorf("unsupported %s type %T", side, v))
	}
}
