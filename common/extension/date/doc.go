// Package date provides utilities to parse and format time.Time values.
//
// # Overview
//
// The package exposes helpers to parse arbitrary values into Go time.Time.
// It includes normalization, common layouts, and formatting functions.
//
// Supported input types:
//   - string: normalized with Cleaning and parsed against defaults
//   - []byte: interpreted as string
//   - time.Time: returned directly
//   - unsupported (including nil): return error
//
// Functions:
//   - ParseFrom(value any) (time.Time, error)
//   - ParseAndFormat(value any, layout string) (string, error)
//   - Cleaning(raw string) string
//
// Example:
//
//	t, err := date.ParseFrom("2025-08-21")
//	// t = 2025-08-21 00:00:00 +0000 UTC, err = nil
//
//	formatted, err := date.ParseAndFormat("21/08/2025", "2006-01-02")
//	// formatted = "2025-08-21", err = nil
package date
