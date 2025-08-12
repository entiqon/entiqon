// File: common/math/decimal/parser_test.go

package decimal_test

import (
	"math"
	"testing"

	"github.com/entiqon/entiqon/common/extension/decimal"
)

func TestDecimalParseFrom(t *testing.T) {
	const epsilon = 1e-6

	tests := []struct {
		name      string
		input     interface{}
		precision int
		want      float64
		expectErr bool
	}{
		{"NoRoundingStringFloat", "1.23456789", -1, 1.23456789, false},
		{"Round2DecimalsString", "1.23456789", 2, 1.23, false},
		{"Round0DecimalsFloat", 123.456, 0, 123, false},
		{"IntInput", 42, 3, 42, false},
		{"BoolTrue", true, 2, 1, false},
		{"BoolFalse", false, 2, 0, false},
		{"InvalidString", "abc", 2, 0, true},
		{"UnsupportedType", struct{}{}, 2, 0, true},
		{"Float32", float32(2.71828), 5, 2.71828, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decimal.ParseFrom(tt.input, tt.precision)
			if tt.expectErr {
				if err == nil {
					t.Errorf("ParseFrom(%v, %d) expected error but got nil", tt.input, tt.precision)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseFrom(%v, %d) unexpected error: %v", tt.input, tt.precision, err)
				return
			}
			if math.Abs(got-tt.want) > epsilon {
				t.Errorf("ParseFrom(%v, %d) = %f; want %f", tt.input, tt.precision, got, tt.want)
			}
		})
	}
}
