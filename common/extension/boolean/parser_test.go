// File: common/boolean/parser_test.go

package boolean

import "testing"

func TestParseFrom(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		var x any
		_, err := ParseFrom(x)
		if err == nil {
			t.Fatalf("ParseFrom(nil) expected error, got nil")
		}
	})

	t.Run("Bool", func(t *testing.T) {
		v, err := ParseFrom(true)
		if err != nil {
			t.Fatalf("ParseFrom(true) unexpected error: %v", err)
		}
		if v != true {
			t.Fatalf("ParseFrom(true) = %v, want %v", v, true)
		}
	})

	t.Run("Float", func(t *testing.T) {
		t.Run("Float32", func(t *testing.T) {
			cases := []struct {
				in   float32
				want bool
			}{
				{0, false},
				{0.0, false},
				{-0.0, false}, // tricky edge: -0.0 == 0.0
				{1.5, true},
				{-2.3, true},
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(float32(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(float32(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})

		t.Run("Float64", func(t *testing.T) {
			cases := []struct {
				in   float64
				want bool
			}{
				{0, false},
				{0.0, false},
				{-0.0, false}, // consistent with float32
				{3.14159, true},
				{-123.456, true},
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(float64(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(float64(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})
	})

	t.Run("Integer", func(t *testing.T) {
		t.Run("Int", func(t *testing.T) {
			cases := []struct {
				in   int
				want bool
			}{
				{0, false},
				{1, true},
				{-1, true},
				{42, true},
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(int(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(int(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})

		t.Run("Int8", func(t *testing.T) {
			cases := []struct {
				in   int8
				want bool
			}{
				{int8(0), false},
				{int8(1), true},
				{int8(-1), true},
				{int8(-128), true}, // min int8
				{int8(127), true},  // max int8
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(int8(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(int8(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})

		t.Run("Int16", func(t *testing.T) {
			cases := []struct {
				in   int16
				want bool
			}{
				{int16(0), false},
				{int16(1), true},
				{int16(-1), true},
				{int16(-32768), true}, // min int16
				{int16(32767), true},  // max int16
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(int16(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(int16(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})

		t.Run("Int32", func(t *testing.T) {
			cases := []struct {
				in   int32
				want bool
			}{
				{int32(0), false},
				{int32(1), true},
				{int32(-1), true},
				{int32(-2147483648), true}, // min int32
				{int32(2147483647), true},  // max int32
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(int32(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(int32(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})

		t.Run("Int64", func(t *testing.T) {
			cases := []struct {
				in   int64
				want bool
			}{
				{int64(0), false},
				{int64(1), true},
				{int64(-1), true},
				{-9223372036854775808, true}, // min int64
				{9223372036854775807, true},  // max int64
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(int64(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(int64(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			cases := map[string]bool{
				"1":     true,
				"t":     true,
				"true":  true,
				"on":    true,
				"y":     true,
				"yes":   true,
				"0":     false,
				"f":     false,
				"false": false,
				"off":   false,
				"n":     false,
				"no":    false,
			}
			for in, want := range cases {
				got, err := ParseFrom(in)
				if err != nil {
					t.Fatalf("ParseFrom(%q) unexpected error: %v", in, err)
				}
				if got != want {
					t.Fatalf("ParseFrom(%q) = %v, want %v", in, got, want)
				}
			}
		})

		t.Run("RejectsUnknownWords", func(t *testing.T) {
			invalids := []string{"enable", "disable", "ok", "", "  ", "sure", "nah"}
			for _, in := range invalids {
				_, err := ParseFrom(in)
				if err == nil {
					t.Fatalf("ParseFrom(%q) expected error, got nil", in)
				}
			}
		})
	})

	t.Run("UnsignedInteger", func(t *testing.T) {
		t.Run("Uint", func(t *testing.T) {
			cases := []struct {
				in   uint
				want bool
			}{
				{0, false},
				{1, true},
				{12345, true},
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(uint(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(uint(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})

		t.Run("Uint8", func(t *testing.T) {
			cases := []struct {
				in   uint8
				want bool
			}{
				{uint8(0), false},
				{uint8(1), true},
				{uint8(255), true}, // max uint8
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(uint8(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(uint8(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})

		t.Run("Uint16", func(t *testing.T) {
			cases := []struct {
				in   uint16
				want bool
			}{
				{uint16(0), false},
				{uint16(1), true},
				{uint16(65535), true}, // max uint16
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(uint16(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(uint16(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})

		t.Run("Uint32", func(t *testing.T) {
			cases := []struct {
				in   uint32
				want bool
			}{
				{uint32(0), false},
				{uint32(1), true},
				{uint32(4294967295), true}, // max uint32
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(uint32(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(uint32(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})

		t.Run("Uint64", func(t *testing.T) {
			cases := []struct {
				in   uint64
				want bool
			}{
				{uint64(0), false},
				{uint64(1), true},
				{18446744073709551615, true}, // max uint64
			}
			for _, c := range cases {
				got, err := ParseFrom(c.in)
				if err != nil {
					t.Fatalf("ParseFrom(uint64(%v)) unexpected error: %v", c.in, err)
				}
				if got != c.want {
					t.Fatalf("ParseFrom(uint64(%v)) = %v, want %v", c.in, got, c.want)
				}
			}
		})
	})

	t.Run("UnsupportedTypes", func(t *testing.T) {
		t.Run("Struct", func(t *testing.T) {
			_, err := ParseFrom(struct{}{})
			if err == nil {
				t.Fatalf("ParseFrom(struct{}) expected error, got nil")
			}
		})

		t.Run("Slice", func(t *testing.T) {
			_, err := ParseFrom([]int{1, 2, 3})
			if err == nil {
				t.Fatalf("ParseFrom([]int) expected error, got nil")
			}
		})

		t.Run("Map", func(t *testing.T) {
			_, err := ParseFrom(map[string]bool{"x": true})
			if err == nil {
				t.Fatalf("ParseFrom(map) expected error, got nil")
			}
		})

		t.Run("ComplexNumbers", func(t *testing.T) {
			_, err := ParseFrom(complex(1, 0))
			if err == nil {
				t.Fatalf("ParseFrom(complex) expected error, got nil")
			}
		})

		t.Run("Channel", func(t *testing.T) {
			ch := make(chan int)
			_, err := ParseFrom(ch)
			if err == nil {
				t.Fatalf("ParseFrom(chan) expected error, got nil")
			}
		})

		t.Run("Function", func(t *testing.T) {
			fn := func() {}
			_, err := ParseFrom(fn)
			if err == nil {
				t.Fatalf("ParseFrom(func) expected error, got nil")
			}
		})
	})
}
