// File: common/number/parser_test.go

package number_test

import (
	"testing"

	"github.com/entiqon/entiqon/common/extension/number"
)

func TestNumberParseFrom(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		round     bool
		want      int
		expectErr bool
	}{
		// Integer types (signed)
		{"int", int(42), false, 42, false},
		{"int8", int8(-12), false, -12, false},
		{"int16", int16(1000), false, 1000, false},
		{"int32", int32(123456), false, 123456, false},
		{"int64", int64(-987654321), false, -987654321, false},

		// Unsigned integer types
		{"uint", uint(42), false, 42, false},
		{"uint8", uint8(255), false, 255, false},
		{"uint16", uint16(65535), false, 65535, false},
		{"uint32", uint32(4294967295), false, 4294967295, false},
		{"uint64Ranged", uint64(1 << 30), false, 1 << 30, false},

		// Overflow uint64
		{"uint64Overflow", uint64(1 << 63), false, 0, true},

		// Integer types with round=true (should behave the same)
		{"intRoundTrue", int(42), true, 42, false},
		{"int8RoundTrue", int8(-12), true, -12, false},
		{"int16RoundTrue", int16(1000), true, 1000, false},
		{"int32RoundTrue", int32(123456), true, 123456, false},
		{"int64RoundTrue", int64(-987654321), true, -987654321, false},

		{"uintRoundTrue", uint(42), true, 42, false},
		{"uint8RoundTrue", uint8(255), true, 255, false},
		{"uint16RoundTrue", uint16(65535), true, 65535, false},
		{"uint32RoundTrue", uint32(4294967295), true, 4294967295, false},
		{"uint64RangedRoundTrue", uint64(1 << 30), true, 1 << 30, false},

		{"uint64OverflowRoundTrue", uint64(1 << 63), true, 0, true},

		// Bool tests
		{"boolTrue", true, false, 1, false},
		{"boolFalse", false, false, 0, false},
		{"boolTrueRoundTrue", true, true, 1, false},
		{"boolFalseRoundTrue", false, true, 0, false},

		// String integers
		{"stringInt", "12345", false, 12345, false},
		{"stringIntNegative", "-42", false, -42, false},
		{"stringIntRoundTrue", "12345", true, 12345, false},
		{"stringIntNegativeRoundTrue", "-42", true, -42, false},

		// Invalid string
		{"stringInvalid", "12abc", false, 0, true},
		{"stringInvalidRoundTrue", "12abc", true, 0, true},

		// Float strict (round=false)
		{"float64ExactIntStrict", 42.0, false, 42, false},
		{"float64CloseIntBelowEpsilonStrict", 42.0000000001, false, 42, false},
		{"float64CloseIntAboveEpsilonStrict", 42.00000001, false, 0, true},
		{"float64NotIntStrict", 42.5, false, 0, true},
		{"float32ExactIntStrict", float32(100.0), false, 100, false},
		{"float32CloseIntBelowEpsilonStrict", float32(100.000000001), false, 100, false},
		{"float32CloseIntAboveEpsilonStrict", float32(100.00001), false, 0, true},
		{"float32NotIntStrict", float32(100.1), false, 0, true},

		// Float lenient (round=true)
		{"float64ExactIntLenient", 42.0, true, 42, false},
		{"float64CloseIntAboveEpsilonLenient", 42.00000001, true, 42, false},
		{"float64NotIntLenient", 42.5, true, 43, false},
		{"float32NotIntLenient", float32(100.1), true, 100, false},

		// String floats strict (round=false)
		{"stringFloatStrictValid", "42.00000000001", false, 42, false},
		{"stringFloatStrictInvalid", "42.00001", false, 0, true},
		{"stringFloatStrictNotInt", "42.5", false, 0, true},

		// String floats lenient (round=true)
		{"stringFloatLenientRound", "42.5", true, 43, false},

		// Unsupported type
		{"unsupportedType", struct{}{}, false, 0, true},
		{"unsupportedTypeRoundTrue", struct{}{}, true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := number.ParseFrom(tt.input, tt.round)
			if tt.expectErr {
				if err == nil {
					t.Errorf("ParseFrom(%v, round=%v) expected error but got nil", tt.input, tt.round)
				}
			} else {
				if err != nil {
					t.Errorf("ParseFrom(%v, round=%v) unexpected error: %v", tt.input, tt.round, err)
					return
				}
				if got != tt.want {
					t.Errorf("ParseFrom(%v, round=%v) = %d; want %d", tt.input, tt.round, got, tt.want)
				}
			}
		})
	}
}
