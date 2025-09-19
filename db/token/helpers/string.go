package helpers

import "fmt"

// Stringify converts a slice of arbitrary values (`[]any`) into a slice of
// their string representations.
//
// Each element is converted using `fmt.Sprint`, preserving Go’s default
// formatting for numbers, booleans, strings, and custom types with
// `String()` methods.
//
// Example:
//
//	values := []any{42, true, "hello"}
//	result := Stringify(values)
//	// result == []string{"42", "true", "hello"}
//
// This utility is useful when you need a uniform `[]string` for logging,
// SQL parameter binding, or serialization.
func Stringify(input []any) []string {
	parts := make([]string, len(input))
	for i, v := range input {
		parts[i] = fmt.Sprint(v)
	}
	return parts
}
