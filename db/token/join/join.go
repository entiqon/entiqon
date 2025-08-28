package join

import (
	"errors"
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/token/table"
)

// join is the unexported implementation of the Token interface.
//
// It preserves:
//   - kind      → the join type (INNER, LEFT, RIGHT, FULL)
//   - left/right → table operands
//   - condition → the ON clause
//   - err       → error state, if any
type join struct {
	kind      Kind
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
//   - Kind (e.g. join.LeftJoin)
//   - string ("LEFT", "LEFT JOIN", case-insensitive)
//
// Example:
//
//	j1 := join.New(join.InnerJoin, "users", "orders", "u.id = o.user_id")
//	j2 := join.New("LEFT", "users", "orders", "u.id = o.user_id")
//
// If kind is invalid, New returns an errored join immediately.
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
	return newWithKind(Inner, left, right, condition)
}

// NewLeft constructs an explicit LEFT JOIN.
func NewLeft(left, right any, condition string) Token {
	return newWithKind(Left, left, right, condition)
}

// NewRight constructs an explicit RIGHT JOIN.
func NewRight(left, right any, condition string) Token {
	return newWithKind(Right, left, right, condition)
}

// NewFull constructs an explicit FULL JOIN.
func NewFull(left, right any, condition string) Token {
	return newWithKind(Full, left, right, condition)
}

// Clone returns a deep copy of the join token.
func (j *join) Clone() Token {
	return &join{
		kind:      j.kind,
		left:      j.left.Clone(),
		right:     j.right.Clone(),
		condition: j.condition,
		err:       j.err,
	}
}

// Kind returns the type of the join (INNER, LEFT, RIGHT, FULL).
func (j *join) Kind() Kind {
	return j.kind
}

// Left returns the left table operand.
func (j *join) Left() table.Token {
	return j.left
}

// Right returns the right table operand.
func (j *join) Right() table.Token {
	return j.right
}

// Condition returns the ON condition string.
func (j *join) Condition() string {
	return j.condition
}

// Debug returns an auditable representation of the join.
// Includes validity state and error if applicable.
func (j *join) Debug() string {
	valid := j.IsValid()
	leftExpr := "<nil>"
	rightExpr := "<nil>"

	if j.left != nil {
		leftExpr = j.left.Raw()
	}
	if j.right != nil {
		rightExpr = j.right.Raw()
	}

	if !valid {
		return fmt.Sprintf(
			"Join{Kind:%q, Left:%q, Right:%q, Condition:%q, Valid:false, Err:%v}",
			j.kind,
			leftExpr,
			rightExpr,
			strings.TrimSpace(j.condition),
			j.err,
		)
	}

	return fmt.Sprintf(
		"Join{Kind:%q, Left:%q, Right:%q, Condition:%q, Valid:true}",
		j.kind,
		leftExpr,
		rightExpr,
		strings.TrimSpace(j.condition),
	)
}

// Error returns the current error on the join, if any.
func (j *join) Error() error {
	return j.err
}

// IsErrored reports whether the join is in an errored state.
func (j *join) IsErrored() bool {
	return j.err != nil
}

// SetError sets the error and returns the receiver as Token.
func (j *join) SetError(err error) Token {
	j.err = err
	return j
}

// IsRaw always returns false for joins.
// A Join is a SQL-generic construct, not a raw fragment.
func (j *join) IsRaw() bool {
	return false
}

// Raw returns the raw SQL fragment of the join clause.
// If the join is errored, Raw returns an empty string.
func (j *join) Raw() string {
	if j.err != nil {
		return ""
	}
	return fmt.Sprintf("%s %s ON %s",
		j.kind,
		j.right.Raw(),
		strings.TrimSpace(j.condition),
	)
}

// Render produces the canonical SQL fragment for the join.
// It is dialect-agnostic and delegates to Raw().
func (j *join) Render() string {
	return j.Raw()
}

// String returns a concise, loggable representation of the join.
// Valid joins are marked with ✅, invalid ones with ⛔.
func (j *join) String() string {
	base := fmt.Sprintf("%s %s ON %s",
		j.kind,
		func() string {
			if j.right != nil {
				return j.right.Raw()
			}
			return ""
		}(),
		strings.TrimSpace(j.condition),
	)

	if !j.IsValid() {
		return fmt.Sprintf("⛔ join(%q): %v", base, j.err)
	}

	return fmt.Sprintf("✅ join(%q)", base)
}

// IsValid reports whether the join is valid (not errored).
func (j *join) IsValid() bool {
	return !j.IsErrored()
}

// newWithKind is the internal constructor.
// It enforces early exit for invalid kinds and validates operands and condition.
func newWithKind(kind any, left, right any, condition string) Token {
	jk := normalizeKind(kind)
	if !jk.IsValid() {
		// EARLY EXIT: unsupported or invalid join kind
		return (&join{}).SetError(errors.New(jk.String()))
	}

	lt := normalizeTable(left, "left")
	rt := normalizeTable(right, "right")

	j := &join{kind: jk, left: lt, right: rt, condition: condition}

	if lt == nil || rt == nil {
		return j.SetError(fmt.Errorf("join requires both left and right tables"))
	}
	if lt.IsErrored() || rt.IsErrored() {
		var errs []string
		if lt.IsErrored() {
			errs = append(errs, fmt.Sprintf("left: %v", lt.Error()))
		}
		if rt.IsErrored() {
			errs = append(errs, fmt.Sprintf("right: %v", rt.Error()))
		}
		return j.SetError(fmt.Errorf("join invalid: %s", strings.Join(errs, "; ")))
	}
	if condition == "" {
		return j.SetError(fmt.Errorf("join condition is empty"))
	}

	return j
}

// normalizeKind resolves arbitrary input into a Kind.
// Accepts join.Kind or string. Unsupported types normalize to invalid.
func normalizeKind(arg any) Kind {
	if k, ok := arg.(Kind); ok {
		return k
	}
	if s, ok := arg.(string); ok {
		return ParseJoinKindFrom(s)
	}
	return Kind(-1)
}

// normalizeTable resolves a join operand into a table.Token.
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
