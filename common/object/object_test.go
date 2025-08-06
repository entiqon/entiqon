// File: common/object/object_test.go

package object_test

import (
	"reflect"
	"testing"

	"github.com/entiqon/entiqon/common/object"
)

func TestObject(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		m := map[string]any{"foo": 123}

		if !object.Exists(m, "foo") {
			t.Errorf("Exists: expected true for key 'foo'")
		}

		if object.Exists(m, "bar") {
			t.Errorf("Exists: expected false for missing key 'bar'")
		}
	})

	t.Run("GetValue", func(t *testing.T) {
		t.Run("GetValue", func(t *testing.T) {
			m := map[string]any{
				"int":    42,
				"string": "hello",
				"float":  3.14,
			}

			tests := []struct {
				key        string
				def        any
				want       any
				typ        string
				shouldFind bool
			}{
				{"int", 0, 42, "int", true},
				{"string", "default", "hello", "string", true},
				{"float", 0.0, 3.14, "float64", true},
				{"missing", "default", "default", "string", false},
				{"wrongtype", 100, 100, "int", false},
			}

			for _, tt := range tests {
				switch tt.typ {
				case "int":
					got := object.GetValue[int](m, tt.key, tt.def.(int))
					if got != tt.want {
						t.Errorf("GetValue[int](%q) = %v; want %v", tt.key, got, tt.want)
					}
				case "string":
					got := object.GetValue[string](m, tt.key, tt.def.(string))
					if got != tt.want {
						t.Errorf("GetValue[string](%q) = %v; want %v", tt.key, got, tt.want)
					}
				case "float64":
					got := object.GetValue[float64](m, tt.key, tt.def.(float64))
					if got != tt.want {
						t.Errorf("GetValue[float64](%q) = %v; want %v", tt.key, got, tt.want)
					}
				}
			}
		})

		t.Run("WithDefaults", func(t *testing.T) {
			m := map[string]any{}

			// Missing key returns default
			gotInt := object.GetValue[int](m, "missing_int", 123)
			if gotInt != 123 {
				t.Errorf("GetValue[int] missing key: got %v, want default %v", gotInt, 123)
			}

			gotStr := object.GetValue[string](m, "missing_str", "default")
			if gotStr != "default" {
				t.Errorf("GetValue[string] missing key: got %q, want default %q", gotStr, "default")
			}

			// Key present but wrong type returns default
			m["val"] = "not_int"
			gotInt2 := object.GetValue[int](m, "val", 456)
			if gotInt2 != 456 {
				t.Errorf("GetValue[int] wrong type: got %v, want default %v", gotInt2, 456)
			}
		})
	})

	t.Run("SetValue", func(t *testing.T) {
		var m map[string]any

		// Test setting nil map initializes it
		m = object.SetValue(m, "a", 1)
		if m == nil {
			t.Fatalf("SetValue did not initialize nil map")
		}
		if v, ok := m["a"].(int); !ok || v != 1 {
			t.Errorf("SetValue did not set key 'a' to 1")
		}

		// Test updating with different value
		m = object.SetValue(m, "a", 2)
		if v := m["a"].(int); v != 2 {
			t.Errorf("SetValue did not update key 'a' to 2")
		}

		// Test skipping update for equal value
		m = object.SetValue(m, "a", 2)
		if v := m["a"].(int); v != 2 {
			t.Errorf("SetValue changed value unexpectedly")
		}

		// Test with complex types (slices)
		slice1 := []int{1, 2, 3}
		slice2 := []int{1, 2, 3}
		slice3 := []int{4, 5, 6}

		m = object.SetValue(m, "slice", slice1)
		if !reflect.DeepEqual(m["slice"], slice1) {
			t.Errorf("SetValue failed to set slice1")
		}

		m = object.SetValue(m, "slice", slice2) // equal slice, should not update
		if !reflect.DeepEqual(m["slice"], slice1) {
			t.Errorf("SetValue updated slice despite deep equality")
		}

		m = object.SetValue(m, "slice", slice3) // different slice, should update
		if !reflect.DeepEqual(m["slice"], slice3) {
			t.Errorf("SetValue failed to update slice to slice3")
		}
	})
}
