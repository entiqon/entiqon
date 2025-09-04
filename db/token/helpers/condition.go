package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/entiqon/entiqon/db/token/types/operator"
)

// ResolveCondition parses a SQL-like condition and returns:
//
//	expr  -> always normalized as "field = :field" (the '=' is a canonical placeholder)
//	op    -> the real operator detected (e.g., BETWEEN, IN, =, IS NULL, ...)
//	value -> RHS normalized: scalar (any), []any for multivalued ops, or nil for IS NULL
//
// Examples:
//
//	"id = 1"                      -> expr="id = :id", op="=", value=1
//	"price BETWEEN 1 AND 3.5"     -> expr="price = :price", op="BETWEEN", value=[1 3.5]
//	"lastname IN ('a','b','c')"   -> expr="lastname = :lastname", op="IN", value=["a","b","c"]
//	"deleted_at IS NULL"          -> expr="deleted_at = :deleted_at", op="IS NULL", value=nil
func ResolveCondition(input string) (field string, op operator.Type, value any, err error) {
	in := strings.TrimSpace(input)

	// Detect operator token
	left, oper, right, hasOperator := resolveExpression(in)
	field = left
	if !hasOperator {
		// Default case: bare identifier, fallback to "="
		field = strings.TrimSpace(in)
		if field == "" {
			return "", operator.Invalid, nil, fmt.Errorf("empty condition")
		}
		if !IsValidIdentifier(field) {
			return "", operator.Invalid, nil, fmt.Errorf("invalid condition expression: %q", input)
		}
		op = operator.Equal
		value = nil
		return
	}

	op = operator.ParseFrom(oper)
	// Parse RHS by operator
	switch op {
	case operator.Between:
		values, err := parseBetween(right)
		if err != nil {
			return "", operator.Invalid, nil, err
		}
		value = values
	case operator.In, operator.NotIn:
		values, err := parseList(right)
		if err != nil {
			return "", operator.Invalid, nil, err
		}
		value = values
	case operator.IsNull, operator.IsNotNull:
		value = nil
	default:
		value = coerceScalar(strings.TrimSpace(right))
	}

	return
}

// IsValidSlice validates that the provided value slice matches
// the requirements of the given operator.
//
// - IN / NOT IN require a non-empty slice
// - BETWEEN requires a slice of exactly 2 elements
// Returns true if valid, false otherwise.
func IsValidSlice(op operator.Type, value any) bool {
	if value == nil {
		return false
	}

	// Normalize slice length regardless of element type
	var length int
	switch vv := value.(type) {
	case []any:
		length = len(vv)
	case []int:
		length = len(vv)
	case []string:
		length = len(vv)
	case []int64:
		length = len(vv)
	case []float64:
		length = len(vv)
	default:
		return false
	}

	switch op {
	case operator.In, operator.NotIn:
		return length > 0
	case operator.Between:
		return length == 2
	default:
		return false
	}
}

// resolveExpression returns (field, opToken, rhs, hasOperator).
// Uses a lowercase copy for search but slices from the original to preserve casing.
func resolveExpression(input string) (field string, op string, value string, hasOperator bool) {
	orig := strings.TrimSpace(input)
	in := strings.ToLower(orig)

	ops := operator.GetKnownOperators()
	for _, o := range ops {
		i := indexOf(in, strings.ToLower(o))
		if i < 0 {
			continue
		}

		j := i + len(o)
		field = strings.TrimSpace(orig[:i])
		op = strings.TrimSpace(orig[i:j]) // preserves original casing/spacing
		value = strings.TrimSpace(orig[j:])
		hasOperator = true
		return
	}

	return
}

// indexOf locates op in s. For word operators it enforces simple boundaries:
// start|space before and space|end|paren after. For symbol operators it uses raw Index.
func indexOf(s, op string) int {
	// symbol-ish if op contains non-letters (e.g., "!=", "<=", "<>")
	isWord := true
	for _, r := range op {
		if (r < 'a' || r > 'z') && r != ' ' {
			isWord = false
			break
		}
	}
	if !isWord {
		return strings.Index(s, op)
	}

	// word operator: scan and enforce loose word boundaries
	for i := 0; ; {
		j := strings.Index(s[i:], op)
		if j < 0 {
			return -1
		}
		j += i
		// check left boundary
		if j == 0 || isBoundary(s[j-1]) {
			// check right boundary
			r := j + len(op)
			if r >= len(s) || isBoundary(s[r]) {
				return j
			}
		}
		i = j + 1
	}
}

func isBoundary(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '(', ')', ',', ';':
		return true
	default:
		return false
	}
}

func parseBetween(rhs string) ([]any, error) {
	// Expect shapes like: "1 AND 3", "'a' and 'z'"
	low := strings.ToLower(rhs)
	parts := strings.SplitN(low, " and ", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid BETWEEN clause: %q", rhs)
	}
	// Use original rhs to extract accurate slices
	i := strings.Index(strings.ToLower(rhs), " and ")
	left := strings.TrimSpace(rhs[:i])
	right := strings.TrimSpace(rhs[i+5:])
	return []any{
		coerceScalar(trimParensOrQuotes(left)),
		coerceScalar(trimParensOrQuotes(right)),
	}, nil
}

func parseList(rhs string) ([]any, error) {
	s := strings.TrimSpace(rhs)
	// Allow optional parentheses: "(1, 2, 3)" or "1, 2, 3"
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		s = strings.TrimSpace(s[1 : len(s)-1])
	}
	if s == "" {
		return []any{}, errors.New("empty list")
	}
	raw := splitCSVRespectingQuotes(s)
	out := make([]any, 0, len(raw))
	for _, item := range raw {
		item = trimParensOrQuotes(strings.TrimSpace(item))
		out = append(out, coerceScalar(item))
	}
	return out, nil
}

func coerceScalar(s string) any {
	// Try int, then float, then NULL, else unquote string
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f // keep float64, don't cast to int
	}
	if strings.EqualFold(strings.TrimSpace(s), "null") {
		return nil
	}
	return trimQuotes(s)
}

func trimParensOrQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && ((s[0] == '(' && s[len(s)-1] == ')') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return trimQuotes(s)
}

func trimQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && ((s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"')) {
		return s[1 : len(s)-1]
	}
	return s
}

// splitCSVRespectingQuotes splits on commas not inside single quotes.
func splitCSVRespectingQuotes(s string) []string {
	var out []string
	start := 0
	inQuote := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'':
			inQuote = !inQuote
		case ',':
			if !inQuote {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	out = append(out, s[start:])
	return out
}

// ToParamKey converts a field like `users.id` or `u."last name"` into a safe key:
//
//	users.id        -> users_id
//	"last name"     -> last_name
//	u."last name"   -> u_last_name
func ToParamKey(field string) string {
	// strip quotes/backticks and turn non-alphanum to underscore
	var b strings.Builder
	for _, r := range field {
		switch {
		case r == '.' || r == '-':
			b.WriteByte('_')
		case r == '`' || r == '"' || r == '\'':
			// skip
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_':
			b.WriteRune(r)
		case unicode.IsSpace(r):
			b.WriteByte('_')
		default:
			b.WriteByte('_')
		}
	}
	key := b.String()
	key = strings.Trim(key, "_")
	if key == "" {
		return "field"
	}
	return strings.ToLower(key)
}
