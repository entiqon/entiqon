package condition

import (
	"errors"
	"fmt"
	"strings"

	"github.com/entiqon/db/token/helpers"
	ct "github.com/entiqon/db/token/types/condition"
	"github.com/entiqon/db/token/types/operator"
)

type token struct {
	kind     ct.Type
	input    string
	name     string
	operator operator.Type
	expr     string
	value    any
	err      error
}

// New creates a new condition token of the given kind.
//
// The expr parameter is the logical expression to evaluate (e.g. "age > 18").
// The value parameter is an optional bound value used for parameterized queries.
// The resulting Token includes both expr and value concatenated as input, and
// the expr is normalized through helpers.ResolveExpression.
//
// Example:
//
//	cond := condition.New(ct.Single, "age > ?", 18)
//	if cond.IsErrored() {
//	    log.Fatal(cond.Error())
//	}
func New(kind ct.Type, input ...any) Token {
	t := &token{
		kind:  kind,
		input: strings.Join(helpers.Stringify(input), " "),
	}

	// Validate kind first
	if !t.kind.IsValid() {
		return t.SetError(errors.New("invalid condition type"))
	}

	if len(input) == 0 {
		return t.SetError(fmt.Errorf("invalid condition input: %q", input))
	}

	// expr must be string
	expr, ok := input[0].(string)
	if !ok {
		return t.SetError(fmt.Errorf("expr must be string, got %T", input[0]))
	}

	// Always resolve expression
	if expr == "" {
		return t.SetError(errors.New("empty expression"))
	}

	// Support (field, operator.Type, value)
	if len(input) == 3 {
		op := operator.ParseFrom(input[1])
		if !op.IsValid() {
			return t.SetError(fmt.Errorf("invalid operator: %v", op))
		}

		t.name = helpers.ToParamKey(expr)
		t.operator = op
		t.value = input[2]
		t.expr = fmt.Sprintf("%s %s :%s", expr, op, t.name)

		// Slice-based operator validation (IN, NOT IN, BETWEEN)
		if op == operator.In || op == operator.NotIn || op == operator.Between {
			if !helpers.IsValidSlice(op, t.value) {
				return t.SetError(fmt.Errorf("invalid value list for operator %s", op))
			}
		}

		return t
	}

	field, op, value, err := helpers.ResolveCondition(expr)
	if err != nil {
		return t.SetError(err)
	}

	t.name = helpers.ToParamKey(field)
	if op != operator.IsNull && op != operator.IsNotNull {
		t.expr = fmt.Sprintf("%s %s :%s", field, op, t.name)
	}
	t.operator = op
	if len(input) > 1 {
		t.value = input[1]
	} else {
		t.value = value
	}

	return t
}

// NewAnd creates a condition token of type And.
//
// This is a convenience constructor equivalent to calling
// New(ct.And, expr, value).
//
// Example:
//
//	cond := condition.NewAnd("country = ?", "US")
func NewAnd(expr string, args ...any) Token {
	return New(ct.And, expr, args)
}

// NewOr creates a condition token of type Or.
//
// This is a convenience constructor equivalent to calling
// New(ct.Or, expr, args).
//
// Example:
//
//	cond := condition.NewOr("status = ?", "active")
func NewOr(expr string, args ...any) Token {
	return New(ct.Or, expr, args)
}

// Kind returns the condition type of the token.
//
// It indicates whether the condition is a Single expression,
// or a logical And / Or composition. If the token was created
// with an invalid type, Kind still returns that value, but the
// token will carry an error accessible via Error().
//
// Example:
//
//	cond := condition.NewAnd("age > ?", 18)
//	fmt.Println(cond.Kind()) // Output: And
func (t *token) Kind() ct.Type {
	return t.kind
}

// SetKind assigns a new condition type to the token.
//
// This mutator is rarely needed in practice, since condition
// tokens are typically constructed with the desired type using
// New, NewAnd, or NewOr. It is provided to satisfy the Kindable
// contract and to allow controlled mutation when necessary.
//
// Example:
//
//	cond := condition.New(ct.Single, "age > ?", 18)
//	cond.SetKind(ct.And)
//	fmt.Println(cond.Kind()) // Output: And
func (t *token) SetKind(value ct.Type) {
	t.kind = value
}

// Input returns the original input string provided when the
// token was constructed.
//
// The input is built by joining the expr and value parameters
// passed to New, NewAnd, or NewOr. It represents the raw
// condition as received before normalization.
//
// Example:
//
//	cond := condition.NewAnd("age > ?", 18)
//	fmt.Println(cond.Input()) // Output: "age > ? 18"
func (t *token) Input() string { return t.input }

// Expr returns the normalized expression of the condition token.
//
// The expression is derived from the expr argument passed to the
// constructor. If expr was empty, it is resolved from the input
// using helpers.ResolveExpression. This ensures the token always
// carries a syntactically classified expression suitable for SQL
// rendering.
//
// Example:
//
//	cond := condition.New(ct.Single, "age > ?", 18)
//	fmt.Println(cond.Expr()) // Output: "age > ?"
func (t *token) Expr() string { return t.expr }

// Name returns the condition’s identifier.
//
// When built from a bare identifier with value, it is used as the
// parameter name. For example:
//
//	cond := condition.New(ct.And, "id", 1)
//	fmt.Println(cond.Name()) // "id"
//
// Otherwise Name is empty.
func (t *token) Name() string { return t.name }

// Operator returns the structured SQL operator type for this condition.
// Returns operator.Invalid if parsing failed or was unsupported.
func (t *token) Operator() operator.Type { return t.operator }

// Value returns the bound value associated with the condition token.
//
// The value corresponds to the parameter placeholder in the expression,
// if any. It can be of any type (string, int, time.Time, etc.) and is
// intended for use with parameterized queries. If no value was provided
// when the token was constructed, Value returns nil.
//
// Example:
//
//	cond := condition.New(ct.Single, "age > ?", 18)
//	fmt.Println(cond.Value()) // Output: 18
func (t *token) Value() any { return t.value }

// Error returns the error carried by the condition token, if any.
//
// A token may contain an error if it was constructed with an
// invalid kind, if the expression was empty, or if expression
// resolution failed. When no error has occurred, Error returns nil.
//
// Example:
//
//	cond := condition.New(ct.Single, "", nil)
//	if cond.Error() != nil {
//	    fmt.Println("construction failed:", cond.Error())
//	}
func (t *token) Error() error { return t.err }

// IsErrored reports whether the condition token carries an error.
//
// It is a convenience check equivalent to (t.Error() != nil).
// This is useful for short-circuiting logic when working with
// collections of tokens.
//
// Example:
//
//	cond := condition.New(ct.Single, "", nil)
//	if cond.IsErrored() {
//	    fmt.Println("token is invalid:", cond.Error())
//	}
func (t *token) IsErrored() bool { return t.err != nil }

// SetError assigns the given error to the condition token and
// returns the token itself.
//
// This method is primarily used internally by constructors to
// propagate validation failures. It can also be used by higher-level
// components to attach contextual errors while preserving the token.
//
// Example:
//
//	cond := condition.New(ct.Single, "age > ?", 18)
//	cond = cond.SetError(errors.New("forced failure"))
//	if cond.IsErrored() {
//	    fmt.Println(cond.Error()) // Output: forced failure
//	}
func (t *token) SetError(err error) Token {
	t.err = err
	return t
}

// Debug returns a developer-friendly string representation of the
// condition token. It is intended for diagnostics and logging only,
// not for SQL rendering.
//
// The output includes the raw input string, the condition type, the
// normalized expression, the bound value, and any error state. This
// helps trace both what was passed in and how it was resolved.
//
// Example:
//
//	cond := condition.New(ct.And, "age > ?", 18)
//	fmt.Println(cond.Debug())
//	// Example output:
//	// Condition{Input="age > ? 18", Type:"And", Expression="age > ?", Value=18, Error=<nil>}
func (t *token) Debug() string {
	return fmt.Sprintf(
		"Condition{Input=%q, Type:%q, Expression=%q, Value=%v, Error=%v}",
		t.input, t.kind, t.expr, t.value, t.err,
	)
}

// IsRaw reports whether the condition token should be treated as a
// raw SQL fragment.
//
// For now this always returns false, since condition tokens are not
// yet integrated with dialect-aware raw expression handling. The
// method is included to satisfy the Rawable contract and will be
// extended in the future when raw conditions are supported.
//
// Example:
//
//	cond := condition.New(ct.Single, "age > ?", 18)
//	fmt.Println(cond.IsRaw()) // Output: false
func (t *token) IsRaw() bool {
	return false
}

// Raw returns the unformatted expression of the condition token.
//
// At present, Raw simply returns the normalized expression stored
// in the token. In future iterations this method will delegate to
// a dialect-aware formatter, so the returned string will reflect
// the quoting and rendering rules of the target SQL dialect.
//
// Example:
//
//	cond := condition.New(ct.Single, "age > ?", 18)
//	fmt.Println(cond.Raw()) // Output: "age > ?"
func (t *token) Raw() string {
	return t.expr
}

// Render produces the SQL output of the condition token to be
// consumed by builders.
//
// At this stage Render returns the generic SQL expression held
// in the token. In future iterations, this method will delegate
// to dialect-specific logic so that expressions are properly
// quoted and rendered according to the target SQL dialect.
//
// Example:
//
//	cond := condition.New(ct.Single, "age > ?", 18)
//	fmt.Println(cond.Render()) // Output: "age > ?"
func (t *token) Render() string {
	if t.kind == ct.Single {
		return t.expr
	}
	return fmt.Sprintf("%s %s", t.kind.String(), t.expr)
}

// String returns a concise, human-readable representation of the
// condition token suitable for logging.
//
// Unlike Render, which produces the SQL expression for builders,
// String enriches the output with type, expression, value, and
// error state.
//
// Example:
//
//	cond := condition.New(ct.And, "age > ?", 18)
//	fmt.Println(cond.String())
//	// Example output:
//	// Condition("And"): expr="age > ?", value=18, errored=false
func (t *token) String() string {
	return fmt.Sprintf(
		"Condition(%q): name=%q, value=%v, errored=%v",
		t.expr, t.name, t.value, t.IsErrored(),
	)
}

// IsValid reports whether the condition token is valid.
//
// A token is considered valid when it does not carry an error.
// This is the inverse of IsErrored() and is typically used to
// guard builder logic from including invalid conditions.
//
// Example:
//
//	cond := condition.New(ct.Single, "age > ?", 18)
//	if cond.IsValid() {
//	    fmt.Println("ready to use:", cond.Render())
//	}
func (t *token) IsValid() bool {
	return t.err == nil
}
