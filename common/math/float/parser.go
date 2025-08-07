// File: common/math/float/parser.go

// Package float provides utilities to parse various input types into float64 values.
package float

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ParseFrom converts a variety of input types into a float64 without rounding.
//
// Supported input types:
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64, uintptr
//   - float32, float64
//   - string (parsed as float64)
//   - bool (true returns 1.0, false returns 0.0)
//   - pointers and interfaces wrapping any of the above
//
// Returns an error if the input type is unsupported or parsing fails.
//
// Example:
//
//	f, err := float.ParseFrom("123.456")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(f) // 123.456
func ParseFrom(value interface{}) (float64, error) {
	if value == nil {
		return 0, fmt.Errorf("nil value")
	}

	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 0, fmt.Errorf("nil pointer encountered")
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
	case reflect.String:
		str := strings.TrimSpace(v.String())
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse string to float64: %w", err)
		}
		return f, nil
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type %T", value)
	}
}
