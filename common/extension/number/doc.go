// Package number provides utilities for parsing numeric values into int with optional rounding.
//
// # Overview
//
// The main function ParseFrom converts supported input types into int.
// Floats can be interpreted strictly (must be integer within tolerance)
// or leniently (rounded to nearest integer).
//
// Supported input types:
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64 (with overflow check)
//   - float32, float64 (strict or rounded depending on flag)
//   - string (parsed as integer or float)
//   - bool (true = 1, false = 0)
//
// Unsupported inputs (including nil, slices, maps, structs, complex, functions)
// return an error.
//
// Example:
//
//	val, err := number.ParseFrom("123.6", true)
//	// val = 124, err = nil
//
//	val, err := number.ParseFrom(123.4, false)
//	// val = 0, err != nil
package number
