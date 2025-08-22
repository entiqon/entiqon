package integer_test

import (
	"testing"

	"github.com/entiqon/entiqon/common/extension/integer"
)

func TestInteger(t *testing.T) {
	t.Run("Bool", func(t *testing.T) {
		tests := []struct {
			in   bool
			want int
		}{
			{true, 1},
			{false, 0},
		}
		for _, tt := range tests {
			got, err := integer.ParseFrom(tt.in)
			if err != nil {
				t.Errorf("ParseFrom(%v) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseFrom(%v) = %d, want %d", tt.in, got, tt.want)
			}
		}
	})

	t.Run("Float", func(t *testing.T) {
		tests := []struct {
			in   interface{}
			want int
		}{
			{1.9, 1},
			{-1.9, -1},
			{0.5, 0},
			{3.14159, 3},
		}
		for _, tt := range tests {
			got, err := integer.ParseFrom(tt.in)
			if err != nil {
				t.Errorf("ParseFrom(%v) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseFrom(%v) = %d, want %d", tt.in, got, tt.want)
			}
		}
	})

	t.Run("Integer", func(t *testing.T) {
		tests := []struct {
			in   interface{}
			want int
		}{
			{42, 42},
			{int8(-8), -8},
			{int16(123), 123},
			{int32(-456), -456},
			{int64(789), 789},
		}
		for _, tt := range tests {
			got, err := integer.ParseFrom(tt.in)
			if err != nil {
				t.Errorf("ParseFrom(%v) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseFrom(%v) = %d, want %d", tt.in, got, tt.want)
			}
		}
	})

	t.Run("Nil", func(t *testing.T) {
		_, err := integer.ParseFrom(nil)
		if err == nil {
			t.Errorf("ParseFrom(nil) expected error, got nil")
		}
	})

	t.Run("String", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			tests := []struct {
				in   string
				want int
			}{
				{"42", 42},
				{"  123  ", 123},
				{"3.99", 3},  // truncation
				{"-7.8", -7}, // truncation
			}
			for _, tt := range tests {
				got, err := integer.ParseFrom(tt.in)
				if err != nil {
					t.Errorf("ParseFrom(%q) unexpected error: %v", tt.in, err)
				}
				if got != tt.want {
					t.Errorf("ParseFrom(%q) = %d, want %d", tt.in, got, tt.want)
				}
			}
		})

		t.Run("Invalid", func(t *testing.T) {
			tests := []string{"abc", "", "  "}
			for _, in := range tests {
				_, err := integer.ParseFrom(in)
				if err == nil {
					t.Errorf("ParseFrom(%q) expected error, got nil", in)
				}
			}
		})
	})

	t.Run("Pointers", func(t *testing.T) {
		valInt := 42
		valFloat := 7.9
		valStr := "123"

		tests := []struct {
			in   interface{}
			want int
		}{
			{&valInt, 42},
			{&valFloat, 7}, // truncation
			{&valStr, 123},
		}
		for _, tt := range tests {
			got, err := integer.ParseFrom(tt.in)
			if err != nil {
				t.Errorf("ParseFrom(%v) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseFrom(%v) = %d, want %d", tt.in, got, tt.want)
			}
		}

		t.Run("NilPointer", func(t *testing.T) {
			var p *int
			_, err := integer.ParseFrom(p)
			if err == nil {
				t.Errorf("ParseFrom(nil pointer) expected error, got nil")
			}
		})
	})

	t.Run("UnsignedInteger", func(t *testing.T) {
		tests := []struct {
			in   interface{}
			want int
		}{
			{uint(42), 42},
			{uint8(8), 8},
			{uint16(123), 123},
			{uint32(456), 456},
			{uint64(789), 789},
		}
		for _, tt := range tests {
			got, err := integer.ParseFrom(tt.in)
			if err != nil {
				t.Errorf("ParseFrom(%v) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseFrom(%v) = %d, want %d", tt.in, got, tt.want)
			}
		}
	})

	t.Run("Unsupported", func(t *testing.T) {
		_, err := integer.ParseFrom(struct{}{})
		if err == nil {
			t.Errorf("ParseFrom(struct{}) expected error, got nil")
		}
	})
}
