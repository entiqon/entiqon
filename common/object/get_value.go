// File: common/object/get_value.go

package object

// GetValue returns the value for key cast to T.
// If the key is missing or the type assertion fails, returns defaultVal.
func GetValue[T any](object map[string]any, key string, defaultVal T) T {
	raw, exists := object[key]
	if !exists {
		return defaultVal
	}
	casted, ok := raw.(T)
	if !ok {
		return defaultVal
	}
	return casted
}
