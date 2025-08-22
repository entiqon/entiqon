package extension_test

import (
	"testing"
	"time"

	"github.com/entiqon/entiqon/common/extension"
)

func TestParser(t *testing.T) {
	t.Run("Boolean", func(t *testing.T) {
		cases := []struct {
			in   any
			want bool
		}{
			{true, true},
			{false, false},
			{"true", true},
			{"FALSE", false},
			{"1", true},
			{"0", false},
			{"yes", true},
			{"no", false},
			{"on", true},
			{"off", false},
			{"y", true},
			{"n", false},
			{"t", true},
			{"f", false},
			{"maybe", false}, // fallback to false
			{nil, false},     // fallback to false
		}
		for i, tc := range cases {
			got := extension.Boolean(tc.in)
			if got != tc.want {
				t.Fatalf("case %d: want %v, got %v", i, tc.want, got)
			}
		}

		// BooleanOr fallback
		t.Run("BooleanOr", func(t *testing.T) {
			if got := extension.BooleanOr("maybe", true); got != true {
				t.Errorf("expected fallback true, got %v", got)
			}
			if got := extension.BooleanOr("maybe", false); got != false {
				t.Errorf("expected fallback false, got %v", got)
			}
			if got := extension.BooleanOr("true", false); got != true {
				t.Errorf("expected parsed true, got %v", got)
			}
		})
	})

	t.Run("Number", func(t *testing.T) {
		cases := []struct {
			in   any
			want int
		}{
			{42, 42},
			{"99", 99},
			{"invalid", 0}, // fallback to 0
		}
		for i, tc := range cases {
			got := extension.Number(tc.in)
			if got != tc.want {
				t.Fatalf("case %d: want %v, got %v", i, tc.want, got)
			}
		}

		// NumberOr fallback
		t.Run("NumberOr", func(t *testing.T) {
			if got := extension.NumberOr("invalid", 77); got != 77 {
				t.Errorf("expected fallback 77, got %v", got)
			}
			if got := extension.NumberOr("invalid", -1); got != -1 {
				t.Errorf("expected fallback -1, got %v", got)
			}
			if got := extension.NumberOr("42", 0); got != 42 {
				t.Errorf("expected parsed 42, got %v", got)
			}
		})
	})

	t.Run("Float", func(t *testing.T) {
		cases := []struct {
			in   any
			want float64
		}{
			{3.14159, 3.14159},
			{float32(2.5), 2.5},
			{42, 42.0},
			{"42.75", 42.75},
			{"  -13.5  ", -13.5},
			{"not-a-number", 0}, // fallback
		}
		for i, tc := range cases {
			got := extension.Float(tc.in)
			if got != tc.want {
				t.Fatalf("case %d: want %v, got %v", i, tc.want, got)
			}
		}

		// FloatOr fallback
		t.Run("FloatOr", func(t *testing.T) {
			if got := extension.FloatOr("bad", 1.23); got != 1.23 {
				t.Errorf("expected fallback 1.23, got %v", got)
			}
			if got := extension.FloatOr("bad", -9.9); got != -9.9 {
				t.Errorf("expected fallback -9.9, got %v", got)
			}
			if got := extension.FloatOr("42.5", 0); got != 42.5 {
				t.Errorf("expected parsed 42.5, got %v", got)
			}
		})
	})

	t.Run("Decimal", func(t *testing.T) {
		cases := []struct {
			in        any
			precision int
			want      float64
		}{
			{"3.14159", 2, 3.14},
			{2.345, 2, 2.35},
			{int64(1234), 0, 1234.0},
			{"-0.0049", 2, -0.00},
			{"abc", 2, 0}, // fallback
		}
		for i, tc := range cases {
			got := extension.Decimal(tc.in, tc.precision)
			if got != tc.want {
				t.Fatalf("case %d: want %v, got %v", i, tc.want, got)
			}
		}

		// DecimalOr fallback
		t.Run("DecimalOr", func(t *testing.T) {
			if got := extension.DecimalOr("abc", 2, 9.99); got != 9.99 {
				t.Errorf("expected fallback 9.99, got %v", got)
			}
			if got := extension.DecimalOr("abc", 2, -1.11); got != -1.11 {
				t.Errorf("expected fallback -1.11, got %v", got)
			}
			if got := extension.DecimalOr("123.45", 2, 0); got != 123.45 {
				t.Errorf("expected parsed 123.45, got %v", got)
			}
		})
	})

	t.Run("Date", func(t *testing.T) {
		ymd := func(y int, m time.Month, d int) time.Time {
			return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		}

		cases := []struct {
			in   any
			want time.Time
		}{
			{"20240229", ymd(2024, 2, 29)}, // leap day
			{"bad-date", time.Time{}},      // fallback
		}
		for i, tc := range cases {
			got := extension.Date(tc.in)
			if !got.Equal(tc.want) {
				t.Fatalf("case %d: want %v, got %v", i, tc.want, got)
			}
		}

		// DateOr fallback
		t.Run("DateOr", func(t *testing.T) {
			def1 := ymd(2000, 1, 1)
			def2 := ymd(1999, 12, 31)

			if got := extension.DateOr("bad-date", def1); !got.Equal(def1) {
				t.Errorf("expected fallback def1, got %v", got)
			}
			if got := extension.DateOr("bad-date", def2); !got.Equal(def2) {
				t.Errorf("expected fallback def2, got %v", got)
			}
			if got := extension.DateOr("20230101", def1); !got.Equal(ymd(2023, 1, 1)) {
				t.Errorf("expected parsed 2023-01-01, got %v", got)
			}
		})
	})
}
