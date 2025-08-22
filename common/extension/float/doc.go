// Package float provides utilities to parse values into float64.
//
// # Overview
//
// The package exposes ParseFrom to convert supported inputs into float64
// without rounding. Pointers and interface wrappers are dereferenced before parsing.
//
// Supported input types:
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64, uintptr
//   - float32, float64
//   - string (parsed as float64)
//   - bool (true = 1.0, false = 0.0)
//   - pointers and interfaces wrapping any of the above
//
// Unsupported inputs (including nil, structs, slices, maps, complex, functions)
// return an error.
//
// Example:
//
//	val, err := float.ParseFrom("123.456")
//	// val = 123.456, err = nil
//
//	val, err := float.ParseFrom(42)
//	// val = 42.0, err = nil
//
//	val, err := float.ParseFrom(true)
//	// val = 1.0, err = nil
package float
