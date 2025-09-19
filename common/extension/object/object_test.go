// File: common/object/object_test.go

package object_test

import (
	"testing"
	"time"

	"github.com/entiqon/common/extension/object"
)

type TestStruct struct {
	LineNo       int
	SKU          string
	privateField string // unexported
	Date         time.Time
}

func (ts TestStruct) SKUValue() string {
	return ts.SKU
}

func TestObject(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		t.Run("NilObject", func(t *testing.T) {
			if object.Exists(nil, "anykey") {
				t.Errorf("Exists(nil) = true; want false")
			}
		})

		t.Run("NonStringKeyMap", func(t *testing.T) {
			m := map[int]any{1: "one"}
			if object.Exists(m, "anykey") {
				t.Errorf("Exists should be false for map with non-string keys")
			}
		})

		t.Run("MapStringAny", func(t *testing.T) {
			m := map[string]any{"Foo": 123}
			if !object.Exists(m, "foo") {
				t.Error("Exists should be true for map key 'foo' (case-insensitive)")
			}
			if object.Exists(m, "bar") {
				t.Error("Exists should be false for missing key 'bar'")
			}
		})

		t.Run("Struct", func(t *testing.T) {
			type S struct {
				Foo int
				bar int // unexported
			}
			s := S{Foo: 10}
			if !object.Exists(s, "foo") {
				t.Error("Exists should be true for struct field 'foo' (case-insensitive)")
			}
			if object.Exists(s, "bar") {
				t.Error("Exists should be false for unexported field 'bar'")
			}
			if object.Exists(s, "baz") {
				t.Error("Exists should be false for missing field 'baz'")
			}
		})

		t.Run("PointerToStruct", func(t *testing.T) {
			type S struct {
				Foo int
			}
			s := &S{Foo: 20}
			if !object.Exists(s, "foo") {
				t.Error("Exists should be true for pointer to struct field 'foo'")
			}
		})

		t.Run("NilPointer", func(t *testing.T) {
			type S struct {
				Foo int
			}
			var s *S = nil
			if object.Exists(s, "foo") {
				t.Error("Exists should be false for nil pointer")
			}
		})

		t.Run("UnsupportedKind", func(t *testing.T) {
			if object.Exists(123, "foo") {
				t.Error("Exists should be false for unsupported kind")
			}
		})
	})

	t.Run("GetValue", func(t *testing.T) {
		t.Run("MapValues", func(t *testing.T) {
			m := map[string]any{
				"Int":    42,
				"String": "hello",
				"Float":  3.14,
			}

			tests := []struct {
				key        string
				def        any
				want       any
				typ        string
				shouldFind bool
			}{
				{"Int", 0, 42, "int", true},
				{"String", "default", "hello", "string", true},
				{"Float", 0.0, 3.14, "float64", true},
				{"Missing", "default", "default", "string", false},
				{"WrongType", 100, 100, "int", false},
			}

			for _, tt := range tests {
				t.Run(tt.key, func(t *testing.T) {
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
				})
			}
		})

		t.Run("WithDefaults", func(t *testing.T) {
			m := map[string]any{}

			gotInt := object.GetValue[int](m, "missing_int", 123)
			if gotInt != 123 {
				t.Errorf("GetValue[int] missing key: got %v, want default %v", gotInt, 123)
			}

			gotStr := object.GetValue[string](m, "missing_str", "default")
			if gotStr != "default" {
				t.Errorf("GetValue[string] missing key: got %q, want default %q", gotStr, "default")
			}

			m["val"] = "not_int"
			gotInt2 := object.GetValue[int](m, "val", 456)
			if gotInt2 != 456 {
				t.Errorf("GetValue[int] wrong type: got %v, want default %v", gotInt2, 456)
			}
		})

		t.Run("StructAndPointerValues", func(t *testing.T) {
			rawStruct := TestStruct{
				LineNo: 99,
				SKU:    "XYZ789",
				Date:   time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			}

			rawStructPtr := &rawStruct

			tests := []struct {
				name     string
				input    any
				key      string
				def      any
				expected any
				typ      string
			}{
				{"StructIntField", rawStruct, "LineNo", 0, 99, "int"},
				{"StructStringField", rawStruct, "SKU", "", "XYZ789", "string"},
				{"StructMethod", rawStruct, "SKUValue", "", "XYZ789", "string"},
				{"StructMissingKey", rawStruct, "Missing", "default", "default", "string"},
				{"PointerToStructInt", rawStructPtr, "LineNo", 0, 99, "int"},
				{"PointerToStructMethod", rawStructPtr, "SKUValue", "", "XYZ789", "string"},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					switch def := tt.def.(type) {
					case int:
						got := object.GetValue[int](tt.input, tt.key, def)
						if got != tt.expected {
							t.Errorf("GetValue[int](%q) = %v; want %v", tt.key, got, tt.expected)
						}
					case string:
						got := object.GetValue[string](tt.input, tt.key, def)
						if got != tt.expected {
							t.Errorf("GetValue[string](%q) = %v; want %v", tt.key, got, tt.expected)
						}
					default:
						t.Fatalf("Unsupported default value type %T in test", def)
					}
				})
			}
		})

		t.Run("NilObject", func(t *testing.T) {
			got := object.GetValue[int](nil, "foo", 123)
			if got != 123 {
				t.Errorf("GetValue(nil) = %v; want 123", got)
			}
		})

		t.Run("UnsupportedKind_Int", func(t *testing.T) {
			got := object.GetValue[int](123, "foo", 456)
			if got != 456 {
				t.Errorf("GetValue[int](int) = %v; want 456", got)
			}
		})

		t.Run("UnsupportedKind_Slice", func(t *testing.T) {
			got := object.GetValue[string]([]string{"a"}, "foo", "default")
			if got != "default" {
				t.Errorf("GetValue[string](slice) = %v; want default", got)
			}
		})

		t.Run("MapKeyCaseInsensitive", func(t *testing.T) {
			m := map[string]any{"FOO": 10}
			got := object.GetValue[int](m, "foo", 0)
			if got != 10 {
				t.Errorf("GetValue case insensitive map key = %v; want 10", got)
			}
		})

		t.Run("StructUnexportedField", func(t *testing.T) {
			type S struct {
				Public  int
				private int
			}
			s := S{Public: 1, private: 2}
			got := object.GetValue[int](s, "private", 42)
			if got != 42 {
				t.Errorf("GetValue private field should return default, got %v", got)
			}
		})

		t.Run("FieldCaseInsensitive", func(t *testing.T) {
			type S struct {
				MyField int
			}
			s := S{MyField: 42}
			got := object.GetValue[int](s, "myfield", 0) // clave en minúsculas, campo con mayúscula
			if got != 42 {
				t.Errorf("GetValue case-insensitive field = %v; want 42", got)
			}
		})

		t.Run("NilPointer", func(t *testing.T) {
			type S struct {
				Field int
			}
			var ptr *S = nil
			got := object.GetValue[int](ptr, "Field", 7)
			if got != 7 {
				t.Errorf("GetValue(nil pointer) = %v; want 7", got)
			}
		})
	})

	t.Run("SetValue", func(t *testing.T) {
		t.Run("Map", func(t *testing.T) {
			m := map[string]any{}

			m2, err := object.SetValue(m, "LineNo", 100)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v, ok := m2["LineNo"].(int); !ok || v != 100 {
				t.Errorf("expected map[LineNo]=100, got %v", m2["LineNo"])
			}

			m2, err = object.SetValue(m2, "SKU", "ABC")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v, ok := m2["SKU"].(string); !ok || v != "ABC" {
				t.Errorf("expected map[SKU]=ABC, got %v", m2["SKU"])
			}
		})

		t.Run("NilMap", func(t *testing.T) {
			var m map[string]any // nil map

			m2, err := object.SetValue(m, "NewKey", 123)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if m2 == nil {
				t.Fatalf("SetValue did not initialize nil map")
			}
			if v, ok := m2["NewKey"].(int); !ok || v != 123 {
				t.Errorf("SetValue failed to set 'NewKey'")
			}
		})

		t.Run("StructPointer", func(t *testing.T) {
			s := &TestStruct{}

			s2, err := object.SetValue(s, "LineNo", 555)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s2.LineNo != 555 {
				t.Errorf("expected struct LineNo=555, got %v", s2.LineNo)
			}

			s2, err = object.SetValue(s2, "SKU", "ZZZ")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s2.SKU != "ZZZ" {
				t.Errorf("expected struct SKU=ZZZ, got %v", s2.SKU)
			}
		})

		t.Run("StructFieldNotFound", func(t *testing.T) {
			s := &TestStruct{}
			_, err := object.SetValue(s, "NonExistentField", 10)
			if err == nil || err.Error() != `field "NonExistentField" not found in struct` {
				t.Errorf("SetValue with non-existent field: expected error, got %v", err)
			}
		})

		t.Run("StructPointerUnexportedField", func(t *testing.T) {
			s := &TestStruct{}

			_, err := object.SetValue(s, "privateField", "secret")
			if err == nil {
				t.Errorf("expected error for unexported field, got nil")
			}
		})

		t.Run("InvalidType", func(t *testing.T) {
			var i int = 0

			_, err := object.SetValue(i, "LineNo", 100)
			if err == nil {
				t.Errorf("expected error for unsupported type, got nil")
			}
		})

		t.Run("InvalidObject", func(t *testing.T) {
			var nilIface interface{} = nil
			_, err := object.SetValue(nilIface, "foo", 123)
			if err == nil || err.Error() != "invalid object" {
				t.Errorf("SetValue with invalid object: expected error 'invalid object', got %v", err)
			}
		})

		t.Run("MapWithNonStringKey", func(t *testing.T) {
			type badMap map[int]any
			bm := badMap{1: "one"}
			_, err := object.SetValue(bm, "1", "uno")
			if err == nil {
				t.Errorf("SetValue with map[int]any should return error")
			}
		})

		t.Run("NilPointer", func(t *testing.T) {
			var s *TestStruct = nil
			_, err := object.SetValue(s, "LineNo", 10)
			if err == nil {
				t.Errorf("SetValue with nil pointer should return error")
			}
		})

		t.Run("PointerToNonStruct", func(t *testing.T) {
			i := 5
			iptr := &i
			_, err := object.SetValue(iptr, "foo", "bar")
			if err == nil {
				t.Errorf("SetValue pointer to non-struct should error")
			}
		})

		t.Run("ValueConvertible", func(t *testing.T) {
			s := &TestStruct{}
			var val int32 = 123
			s2, err := object.SetValue(s, "LineNo", val) // LineNo es int, val es int32
			if err != nil {
				t.Fatalf("unexpected error in convertible value: %v", err)
			}
			if s2.LineNo != int(val) {
				t.Errorf("expected LineNo = %d, got %d", val, s2.LineNo)
			}
		})

		t.Run("ValueNotAssignableOrConvertible", func(t *testing.T) {
			s := &TestStruct{}
			_, err := object.SetValue(s, "LineNo", "not an int")
			if err == nil {
				t.Errorf("SetValue with wrong value type should error")
			}
		})

		t.Run("UnsupportedKind", func(t *testing.T) {
			_, err := object.SetValue(5, "LineNo", 10)
			if err == nil {
				t.Errorf("SetValue with unsupported kind should error")
			}
		})
	})
}
