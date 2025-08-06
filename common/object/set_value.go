// File: common/object/set_value.go

package object

import "reflect"

// SetValue assigns value into object[key] only if the key is missing
// or the existing value differs (deep equality). Returns the (possibly initialized) map.
func SetValue[T any](object map[string]any, key string, value T) map[string]any {
	if object == nil {
		object = make(map[string]any)
	}

	existing, ok := object[key]
	if ok {
		if reflect.DeepEqual(existing, value) {
			// Values deeply equal, no update
			return object
		}
	}

	object[key] = value
	return object
}
