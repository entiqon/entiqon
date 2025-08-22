// Package decimal provides utilities to parse values into float64 with optional rounding.
//
// # Overview
//
// The main function ParseFrom converts supported input types into float64 and
// applies decimal rounding using math.Round. A negative precision disables rounding.
//
// Supported input types:
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
//   - string (parsed as float64)
//   - bool (true = 1.0, false = 0.0)
//
// Unsupported inputs (including nil, slices, maps, structs, complex, functions)
// return an error.
//
// Example:
//
//	val, err := decimal.ParseFrom("123.456789", 2)
//	// val = 123.46, err = nil
//
//	val, err := decimal.ParseFrom(42, 0)
//	// val = 42.0, err = nil
//
//	val, err := decimal.ParseFrom(true, 3)
//	// val = 1.0, err = nil
//
//	val, err := decimal.ParseFrom("invalid", 2)
//	// val = 0.0, err != nil
package decimal
