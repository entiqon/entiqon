// File: common/math/decimal/parser.go

// Package decimal provides utilities for parsing and rounding decimal numbers.
package decimal

import (
	"math"

	"github.com/entiqon/entiqon/common/extension/float"
)

// ParseFrom converts the input value into a float64, rounding to the specified decimal precision.
//
// Supported input types include basic numeric types, strings representing numeric values, and booleans:
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
//   - string (parsed as float64)
//   - bool (true returns 1.0, false returns 0.0)
//
// The parsed float64 value is rounded to `precision` decimal places using standard rounding rules.
//
// Parameters:
//   - value: the input value to parse and convert.
//   - precision: the number of decimal places to round the output value.
//     If precision is negative, no rounding is applied and the raw parsed float64 is returned.
//
// Returns:
//   - float64: the parsed and rounded floating-point value.
//   - error: non-nil if parsing fails or if the input type is unsupported.
//
// Example usage:
//
//	val, err := ParseFrom("123.456789", 2)
//	// val == 123.46, err == nil
//
//	val, err := ParseFrom(42, 0)
//	// val == 42.0, err == nil
//
//	val, err := ParseFrom(true, 3)
//	// val == 1.0, err == nil
//
//	val, err := ParseFrom("invalid", 2)
//	// val == 0.0, err != nil
func ParseFrom(value interface{}, precision int) (float64, error) {
	f, err := float.ParseFrom(value)
	if err != nil {
		return 0, err
	}
	if precision < 0 {
		return f, nil
	}
	factor := math.Pow(10, float64(precision))
	return math.Round(f*factor) / factor, nil
}
