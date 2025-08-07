// File: common/number/parser.go

// Package number provides utilities for parsing numeric values from generic input types.
package number

import (
	"fmt"
	"math"
	"strconv"
)

// ParseFrom attempts to parse an integer value from an input of generic type.
//
// The input value can be any of the following types:
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64 (uint64 values are checked for overflow)
//   - float32, float64 (fractional parts are checked to be close to integers within a small tolerance if round=false,
//     or always rounded to nearest int if round=true)
//   - string (parsed using strconv.Atoi, or strconv.ParseFloat with same rounding rules)
//   - bool (false = 0, true = 1)
//
// If the input is a string, it will be parsed as a base-10 integer if possible;
// otherwise parsed as float with rounding behavior controlled by round flag.
// If the input is a numeric type, it will be converted to int.
// If the input is a bool, false returns 0, true returns 1.
// If the input type is unsupported or conversion fails, an error is returned.
//
// Parameters:
//   - value: the input value to parse (any supported type)
//   - round: if true, always round floats to nearest int;
//     if false, reject floats not within 1e-9 tolerance of an integer.
//
// Returns the parsed integer or an error if parsing/conversion is not possible.
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
		if v > uint64(^uint(0)>>1) {
			return 0, fmt.Errorf("uint64 value too large for int")
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
