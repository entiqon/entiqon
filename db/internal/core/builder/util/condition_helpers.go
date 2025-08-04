// File: db/internal/core/builder/util/condition_helpers.go

package util

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var dollarPlaceholderRegex = regexp.MustCompile(`\$\d+`)

// AllSameType checks whether all values in the slice have the same Go type.
//
// This is used in condition resolution to validate that operators like IN and
// BETWEEN receive semantically consistent input (e.g., all strings, all ints).
//
// It accepts a []any and returns true if all values share the same reflect.Type.
//
// Example:
//
//	allSameType([]any{"a", "b", "c"}) → true
//	allSameType([]any{1, "b"})        → false
//
// Since: v1.4.0
func AllSameType(values []any) bool {
	if len(values) < 2 {
		return true
	}
	first := reflect.TypeOf(values[0])
	for _, v := range values[1:] {
		if reflect.TypeOf(v) != first {
			return false
		}
	}
	return true
}

// ContainsUnboundPlaceholder checks if the input condition string contains
// any placeholder token: ?, :name, or $1-style positional.
//
// Used to reject expressions like "status = ?" with no bound value.
//
// Since: v1.4.0
func ContainsUnboundPlaceholder(input string) bool {
	input = strings.TrimSpace(input)

	return strings.Contains(input, "?") ||
		strings.Contains(input, ":") ||
		dollarPlaceholderRegex.MatchString(input)
}

// InferLiteralType tries to cast a single string into a Go scalar type.
// Supports bool, int, float, and fallback to string.
//
// Example:
//
//	"true"   → bool(true)
//	"42"     → int(42)
//	"3.14"   → float64(3.14)
//	"hello"  → string("hello")
//
// Since: v1.4.0
func InferLiteralType(input string) any {
	input = strings.TrimSpace(input)

	if i, err := strconv.Atoi(input); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(input, 64); err == nil {
		return f
	}
	switch strings.ToLower(input) {
	case "true":
		return true
	case "false":
		return false
	}
	return input
}

// IsPlaceholderExpression returns true if name contains ? without known field/operator structure
func IsPlaceholderExpression(input string) bool {
	return strings.Contains(input, "?") && !strings.ContainsAny(input, "=<>!IN")
}

func ParsePlaceholderPattern(input string) (field string, operator string, ok bool) {
	upper := strings.ToUpper(input)

	switch {
	case strings.Contains(upper, " BETWEEN ? AND ?"):
		return strings.TrimSpace(strings.Split(upper, "BETWEEN")[0]), "BETWEEN", true
	case strings.Contains(upper, " NOT LIKE ?"):
		return strings.TrimSpace(strings.Split(upper, "NOT LIKE")[0]), "NOT LIKE", true
	case strings.Contains(upper, " LIKE ?"):
		return strings.TrimSpace(strings.Split(upper, "LIKE")[0]), "LIKE", true
	case strings.Contains(upper, " NOT IN ?"):
		return strings.TrimSpace(strings.Split(upper, "NOT IN")[0]), "NOT IN", true
	case strings.Contains(upper, " IN ?"):
		return strings.TrimSpace(strings.Split(upper, "IN")[0]), "IN", true
	case strings.Contains(upper, " >= ?"):
		return strings.TrimSpace(strings.Split(upper, ">=")[0]), ">=", true
	case strings.Contains(upper, " <= ?"):
		return strings.TrimSpace(strings.Split(upper, "<=")[0]), "<=", true
	case strings.Contains(upper, " <> ?"):
		return strings.TrimSpace(strings.Split(upper, "<>")[0]), "<>", true
	case strings.Contains(upper, " != ?"):
		return strings.TrimSpace(strings.Split(upper, "!=")[0]), "!=", true
	case strings.Contains(upper, " > ?"):
		return strings.TrimSpace(strings.Split(upper, ">")[0]), ">", true
	case strings.Contains(upper, " < ?"):
		return strings.TrimSpace(strings.Split(upper, "<")[0]), "<", true
	case strings.Contains(upper, " = ?"):
		return strings.TrimSpace(strings.Split(upper, "=")[0]), "=", true
	}
	return "", "", false
}
