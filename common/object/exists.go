// File: common/object/exists.go

package object

// Exists checks if the key exists in the object.
func Exists(object map[string]any, key string) bool {
	_, exists := object[key]
	return exists
}
