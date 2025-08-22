// Package boolean provides utilities for parsing and formatting boolean values.
//
// # Overview
//
// This package exposes ParseFrom to parse arbitrary values into bools, and
// BoolToString to format a boolean into custom string representations.
//
// Supported parsing inputs for ParseFrom:
//   - bool: returned directly
//   - string: accepts "true", "false", "1", "0", "yes", "no",
//     "on", "off", "y", "n", "t", "f" (case-insensitive, trimmed)
//   - integers and unsigned integers: zero = false, non-zero = true
//   - floating points: zero = false, non-zero = true
//
// Functions:
//   - ParseFrom(value any) (bool, error)
//   - BoolToString(b bool, trueStr, falseStr string) string
//
// Deprecated:
//   - BoolToStr: alias of BoolToString, use BoolToString instead.
//
// Example:
//
//	val, err := boolean.ParseFrom("yes")					-> // val = true, err = nil
//
//	s := boolean.BoolToString(false, "enabled", "disabled")	-> // s = "disabled"
package boolean
