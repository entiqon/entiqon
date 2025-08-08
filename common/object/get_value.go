// File: common/object/get_value.go

package object

import (
	"reflect"
	"strings"
)

// GetValue attempts to retrieve a property named `key` from `object`
// and returns it as type T. If not found or type assertion fails,
// returns defaultVal.
// Supports:
// - map[string]any
// - structs (exported fields or methods with zero args)
// - pointer to struct
func GetValue[T any](object any, key string, defaultVal T) T {
	if object == nil {
		return defaultVal
	}

	v := reflect.ValueOf(object)

	// Handle pointer to underlying value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		// Support case-insensitive string keys
		for _, mapKey := range v.MapKeys() {
			if mapKey.Kind() == reflect.String && strings.EqualFold(mapKey.String(), key) {
				val := v.MapIndex(mapKey)
				if val.IsValid() {
					if casted, ok := val.Interface().(T); ok {
						return casted
					}
				}
			}
		}

	case reflect.Struct:
		// Try to find field by name (case-insensitive)
		field := v.FieldByName(key)
		if !field.IsValid() {
			for i := 0; i < v.NumField(); i++ {
				f := v.Type().Field(i)
				if strings.EqualFold(f.Name, key) {
					field = v.Field(i)
					break
				}
			}
		}
		if field.IsValid() && field.CanInterface() {
			if casted, ok := field.Interface().(T); ok {
				return casted
			}
		}

		// Try method with key name (getter)
		method := v.MethodByName(key)
		if method.IsValid() && method.Type().NumIn() == 0 && method.Type().NumOut() == 1 {
			results := method.Call(nil)
			if len(results) == 1 {
				if casted, ok := results[0].Interface().(T); ok {
					return casted
				}
			}
		}

	default:
		// For any other kind, return defaultVal
		return defaultVal
	}

	return defaultVal
}
