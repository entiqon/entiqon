// File: common/number/parser.go

// Package number provides utilities for parsing numeric values from generic input types.
package number

import (
	"fmt"
	"strconv"
)

// ParseFrom attempts to parse an integer value from an input of generic type.
//
// The input value can be any of the following types:
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64 (uint64 values are checked for overflow)
//   - float32, float64 (fractional parts are truncated)
//   - string (parsed using strconv.Atoi)
//   - bool (false = 0, true = 1)
//
// If the input is a string, it will be parsed as a base-10 integer.
// If the input is a numeric type, it will be converted to int.
// If the input is a bool, false returns 0, true returns 1.
// If the input type is unsupported or conversion fails, an error is returned.
//
// Note: Conversion from float will truncate the fractional part without rounding.
//
// Examples:
//
//	ParseFrom("123") => 123, nil
//	ParseFrom(123.9) => 123, nil
//	ParseFrom(true) => 1, nil
//	ParseFrom(false) => 0, nil
//
// Returns the parsed integer or an error if parsing/conversion is not possible.
func ParseFrom(value interface{}) (int, error) {
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
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("failed to parse string to int: %w", err)
		}
		return i, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type %T for ParseFrom", value)
	}
}
