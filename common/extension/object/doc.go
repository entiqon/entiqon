// Package object provides utilities for dynamic access and manipulation
// of heterogeneous objects using reflection.
//
// # Overview
//
// This package allows case-insensitive access to maps and structs,
// providing uniform functions to check, retrieve, and update properties.
//
// Supported object types:
//   - map[string]any (case-insensitive key matching)
//   - struct (exported fields, case-insensitive)
//   - pointer to struct
//
// Functions:
//
//	Exists(object any, key string) bool
//	  Checks whether a property exists in a map or struct.
//
//	GetValue[T any](object any, key string, defaultVal T) T
//	  Retrieves a property by key, returning defaultVal if missing or type mismatch.
//
//	SetValue[O any, T any](object O, key string, value T) (O, error)
//	  Sets a property on a map or struct pointer and returns updated object.
//
// Example:
//
//	m := map[string]any{"Foo": 123}
//	ok := object.Exists(m, "foo") // true
//
//	val := object.GetValue[int](m, "foo", 0) // 123
//
//	m, _ = object.SetValue(m, "Bar", "baz")
//	// m["Bar"] = "baz"
package object
