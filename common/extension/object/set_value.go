// File: common/object/set_value.go

package object

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// SetValue sets the property or key named `key` on `object` to `value`.
// Supports:
// - map[string]any: sets key/value pair.
// - pointer to struct: sets exported field via reflection.
// Returns updated object of same type or error on failure.
//
// Example usage:
//
//	m := map[string]any{"foo": 1}
//	m, _ = SetValue(m, "foo", 2) // map updated
//
//	type S struct { Foo int }
//	s := &S{Foo: 1}
//	s, _ = SetValue(s, "Foo", 2) // struct field updated
func SetValue[O any, T any](object O, key string, value T) (O, error) {
	// Use reflect.Value on object
	v := reflect.ValueOf(object)
	if !v.IsValid() {
		return object, errors.New("invalid object")
	}

	// Handle nil pointer
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return object, errors.New("nil pointer passed")
	}

	switch v.Kind() {
	case reflect.Map:
		// Must be map[string]any
		if v.Type().Key().Kind() != reflect.String {
			return object, fmt.Errorf("map key type is %v, expected string", v.Type().Key().Kind())
		}

		// Initialize map if nil
		if v.IsNil() {
			newMap := reflect.MakeMap(v.Type())
			reflect.ValueOf(&object).Elem().Set(newMap)
			v = reflect.ValueOf(object)
		}

		keyVal := reflect.ValueOf(key)
		valVal := reflect.ValueOf(value)
		v.SetMapIndex(keyVal, valVal)

		return object, nil

	case reflect.Ptr:
		// Dereference pointer
		elem := v.Elem()
		if elem.Kind() != reflect.Struct {
			return object, fmt.Errorf("pointer does not point to struct")
		}
		return setStructField(object, elem, key, value)

	default:
		return object, fmt.Errorf("unsupported kind %v, must be map or pointer to struct", v.Kind())
	}
}

func setStructField[O any, T any](object O, v reflect.Value, key string, value T) (O, error) {
	typ := v.Type()

	// Find field case-insensitive
	var field reflect.StructField
	var found bool
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if strings.EqualFold(f.Name, key) {
			field = f
			found = true
			break
		}
	}

	if !found {
		return object, fmt.Errorf("field %q not found in struct", key)
	}

	fv := v.FieldByName(field.Name)
	if !fv.IsValid() || !fv.CanSet() {
		return object, fmt.Errorf("field %q cannot be set (unexported or non-addressable)", key)
	}

	valVal := reflect.ValueOf(value)

	if valVal.Type().AssignableTo(fv.Type()) {
		fv.Set(valVal)
		return object, nil
	}

	if valVal.Type().ConvertibleTo(fv.Type()) {
		fv.Set(valVal.Convert(fv.Type()))
		return object, nil
	}

	return object, fmt.Errorf("cannot assign value of type %v to field %q of type %v", valVal.Type(), key, fv.Type())
}
