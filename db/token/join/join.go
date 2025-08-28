package join

import (
	"errors"
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/token/table"
)

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
// The kind argument may be provided as:
//
//   - token.Kind (e.g. token.Left)
//   - string ("LEFT", "LEFT JOIN", case-insensitive)
//
// Example:
//
//	j1 := join.New(token.Inner, "users", "orders", "u.id = o.user_id")
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

func (j *join) Kind() Kind {
	return j.kind
}

func (j *join) Left() table.Token {
	return j.left
}

func (j *join) Right() table.Token {
	return j.right
}

func (j *join) Condition() string {
	return j.condition
}

// Debug returns an auditable representation of the join.
// If either side is invalid, the join is reported as invalid.
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

// Error returns the current error on the Join, if any.
func (j *join) Error() error {
	return j.err
}

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

// Raw returns the full join clause using the raw representation
// of its components. This is mainly useful when the right side
// is a raw table (e.g. a subquery).
func (j *join) Raw() string {
	if j.err != nil {
		return ""
	}
	return fmt.Sprintf("%s %s ON %s",
		j.kind,
		j.right.Raw(), // if right is raw (subquery), this shows it
		strings.TrimSpace(j.condition),
	)
}

// Render produces the canonical SQL join fragment
// in a dialect-agnostic way.
func (j *join) Render() string {
	return j.Raw()
}

// String returns a loggable representation of the join.
// Always includes the clause shape; marks invalid joins with ⛔.
func (j *join) String() string {
	// always build the base clause
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

func (j *join) IsValid() bool {
	return !j.IsErrored()
}

// newWithKind is the internal constructor.
// Enforces early exit if kind is invalid.
func newWithKind(kind any, left, right any, condition string) Token {
	jk := normalizeKind(kind)
	if !jk.IsValid() {
		// EARLY EXIT: unsupported or invalid join kind
		return (&join{}).SetError(errors.New(jk.String()))
	}

	lt := normalizeTable(left, "left")
	rt := normalizeTable(right, "right")

	j := &join{kind: jk, left: lt, right: rt, condition: condition}

	// further validations
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

// normalizeKind resolves any input into a Kind.
// Supports token.Kind and string. Any unsupported
// input normalizes into an invalid Kind (-1).
func normalizeKind(arg any) Kind {
	if k, ok := arg.(Kind); ok {
		// already a Kind → return directly
		return k
	}
	if s, ok := arg.(string); ok {
		// string → parse into Kind
		return ParseJoinKindFrom(s)
	}
	// unsupported type → invalid
	return Kind(-1)
}

// normalizeTable resolves left/right args into table.Token.
// Accepts string (delegates to table.New), table.Token,
// or nil. Unsupported types produce an errored table.Token.
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
