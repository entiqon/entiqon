// File: common/math/float/parser_test.go

package float_test

import (
	"math"
	"testing"

	"github.com/entiqon/common/extension/float"
)

func TestParseFrom(t *testing.T) {
	const epsilon = 1e-6

	tests := []struct {
		name      string
		input     interface{}
		want      float64
		expectErr bool
	}{
		{"int", 42, 42.0, false},
		{"int8", int8(-12), -12.0, false},
		{"uint64", uint64(1 << 30), float64(1 << 30), false},
		{"float64", 3.14159, 3.14159, false},
		{"float32", float32(2.71828), float64(2.71828), false},

		{"boolTrue", true, 1.0, false},
		{"boolFalse", false, 0.0, false},

		{"boolPointerTrue", func() interface{} { b := true; return &b }(), 1.0, false},
		{"boolPointerFalse", func() interface{} { b := false; return &b }(), 0.0, false},

		{"stringInt", "123", 123.0, false},
		{"stringFloat", "123.456", 123.456, false},
		{"stringInvalid", "abc", 0, true},

		{"interfaceInt", func() interface{} { var i interface{} = 123; return i }(), 123.0, false},
		{"interfaceBool", func() interface{} { var b interface{} = true; return b }(), 1.0, false},

		{"nilValue", nil, 0, true},
		{"pointerToInt", func() interface{} { i := 7; return &i }(), 7.0, false},
		{"unsupportedStruct", struct{}{}, 0, true},
		{"nilPointer", func() interface{} { var p *int = nil; return p }(), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := float.ParseFrom(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("ParseFrom(%v) expected error but got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseFrom(%v) unexpected error: %v", tt.input, err)
				return
			}
			if math.Abs(got-tt.want) > epsilon {
				t.Errorf("ParseFrom(%v) = %f; want %f", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseFrom_InterfaceDirect(t *testing.T) {
	var val interface{} = 123
	got, err := float.ParseFrom(val)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 123 {
		t.Errorf("got %v, want 123", got)
	}
}
