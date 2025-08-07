// File: common/number/parser_test.go

package number_test

import (
	"testing"

	"github.com/entiqon/entiqon/common/number"
)

func TestParseFrom(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		want      int
		expectErr bool
	}{
		{"int", int(42), 42, false},
		{"int8", int8(-12), -12, false},
		{"int16", int16(1000), 1000, false},
		{"int32", int32(123456), 123456, false},
		{"int64", int64(-987654321), -987654321, false},

		{"uint", uint(42), 42, false},
		{"uint8", uint8(255), 255, false},
		{"uint16", uint16(65535), 65535, false},
		{"uint32", uint32(4294967295), 4294967295, false},

		{"uint64Ranged", uint64(1 << 30), 1 << 30, false},
		{"uint64Overflow", uint64(1 << 63), 0, true},

		{"float32", float32(123.456), 123, false},
		{"float64", -9876.54321, -9876, false},

		{"StringValid", "12345", 12345, false},
		{"StringNegative", "-42", -42, false},
		{"StringInvalid", "12abc", 0, true},

		{"BoolTrue", true, 1, false},
		{"BoolFalse", false, 0, false},

		{"UnsupportedTypeStruct", struct{}{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := number.ParseFrom(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("ParseFrom(%v) expected error but got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseFrom(%v) unexpected error: %v", tt.input, err)
					return
				}
				if got != tt.want {
					t.Errorf("ParseFrom(%v) = %d; want %d", tt.input, got, tt.want)
				}
			}
		})
	}
}
