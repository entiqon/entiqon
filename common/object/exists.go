// File: common/object/exists.go

// Package object provides utilities for dynamic access and manipulation
// of flexible objects represented as maps or structs.
//
// It supports case-insensitive property/key lookups, generic typed retrieval,
// safe setting of struct fields or map entries, and handles pointers and methods.
//
// This package simplifies working with heterogeneous data structures where
// both map-based and struct-based objects need to be accessed uniformly.
package object

import (
	"reflect"
	"strings"
)

// Exists checks whether the given key/property exists in the object.
// Supports map[string]any and structs (including pointer to struct).
// Case-insensitive key/field matching.
func Exists(object any, key string) bool {
	if object == nil {
		return false
	}

	v := reflect.ValueOf(object)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return false
		}
		for _, mapKey := range v.MapKeys() {
			if mapKey.Kind() == reflect.String && strings.EqualFold(mapKey.String(), key) {
				return true
			}
		}
		return false

	case reflect.Struct:
		typ := v.Type()
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if strings.EqualFold(f.Name, key) && f.PkgPath == "" { // only exported fields
				return true
			}
		}
		return false

	default:
		return false
	}
}
