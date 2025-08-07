// File: common/math/number/parser.go

// Package number provides utilities for parsing numeric values to integers with optional rounding.
package number

import (
	"fmt"
	"math"
	"strconv"
)

// ParseFrom attempts to parse an integer value from a generic input.
//
// Supported input types include:
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64 (with overflow check)
//   - float32, float64 (fractional parts are either rounded or strictly validated)
//   - string (parsed as integer or float with rounding behavior)
//   - bool (false = 0, true = 1)
//
// The round parameter controls float handling:
//   - round == true: floats are rounded to nearest integer.
//   - round == false: floats must be within 1e-9 of an integer, otherwise an error is returned.
//
// Returns the parsed integer or an error if parsing/conversion fails.
//
// Example usage:
//
//	val, err := number.ParseFrom("123.6", true)  // val = 124, err = nil
//	val, err := number.ParseFrom(123.4, false)   // err returned, not close to integer
func ParseFrom(value interface{}, round bool) (int, error) {
	const epsilon = 1e-9

	checkFloat := func(f float64) (int, error) {
		if !round {
			if diff := math.Abs(f - math.Round(f)); diff > epsilon {
				return 0, fmt.Errorf("float value %v is not integer within tolerance", f)
			}
		}
		return int(math.Round(f)), nil
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		maxInt := uint64(^uint(0) >> 1)
		if v > maxInt {
			return 0, fmt.Errorf("uint64 value %v too large for int", v)
		}
		return int(v), nil
	case float32:
		return checkFloat(float64(v))
	case float64:
		return checkFloat(v)
	case string:
		i, err := strconv.Atoi(v)
		if err == nil {
			return i, nil
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse string to number: %w", err)
		}
		return checkFloat(f)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type %T for ParseFrom", value)
	}
}
