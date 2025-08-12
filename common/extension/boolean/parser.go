// File: common/boolean.go

package boolean

import (
	"fmt"
	"reflect"
	"strings"
)

// ParseFrom attempts to parse the given value into a boolean.
//
// Supported input types:
// - bool: returns the value directly.
// - string: accepts "true", "false", "1", "0", "yes", "no" (case-insensitive).
// - integer types: zero is false, non-zero is true.
// - unsigned integer types: zero is false, non-zero is true.
// - floating point types: zero is false, non-zero is true.
//
// Returns an error if the value cannot be interpreted as a boolean.
//
// Examples:
//
//	boolVal, err := ParseFrom(true)
//	boolVal, err := ParseFrom("yes")
//	boolVal, err := ParseFrom(1)
//	boolVal, err := ParseFrom(0.0)
//	boolVal, err := ParseFrom("invalid") // error returned
func ParseFrom(value any) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		s := strings.ToLower(strings.TrimSpace(v))
		switch s {
		case "true", "1", "yes":
			return true, nil
		case "false", "0", "no":
			return false, nil
		default:
			return false, fmt.Errorf("invalid boolean string: %q", v)
		}
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() != 0, nil
	case float32, float64:
		return reflect.ValueOf(v).Float() != 0, nil
	default:
		return false, fmt.Errorf("unsupported type %T for boolean parsing", v)
	}
}
