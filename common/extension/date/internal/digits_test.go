// File: common/extension/date/internal/digits_test.go

package internal_test

import (
	"math"
	"testing"

	"github.com/entiqon/common/extension/date/internal"
)

func TestDigits(t *testing.T) {
	t.Run("Methods", func(t *testing.T) {

		t.Run("AllDigits", func(t *testing.T) {
			t.Run("EmptyString", func(t *testing.T) {
				// NOTE: Current implementation returns true for empty string.
				if ok := internal.AllDigits(""); !ok {
					t.Fatalf("expected true for empty string per current implementation")
				}
			})
			t.Run("OnlyDigits", func(t *testing.T) {
				cases := []string{"0", "9", "123", "00123", "9876543210"}
				for _, in := range cases {
					if !internal.AllDigits(in) {
						t.Fatalf("expected true for %q", in)
					}
				}
			})
			t.Run("NonDigits", func(t *testing.T) {
				cases := []string{"12a3", " 123", "123 ", "12.3", "+123", "-123", "١٢٣"} // non-ASCII numerals too
				for _, in := range cases {
					if internal.AllDigits(in) {
						t.Fatalf("expected false for %q", in)
					}
				}
			})
		})

		t.Run("DecimalDigits", func(t *testing.T) {
			t.Run("ZeroIsOneDigit", func(t *testing.T) {
				if got := internal.DecimalDigits(0); got != 1 {
					t.Fatalf("want 1, got %d", got)
				}
			})
			t.Run("Positive", func(t *testing.T) {
				if got := internal.DecimalDigits(7); got != 1 {
					t.Fatalf("want 1, got %d", got)
				}
				if got := internal.DecimalDigits(42); got != 2 {
					t.Fatalf("want 2, got %d", got)
				}
				if got := internal.DecimalDigits(1234567890); got != 10 {
					t.Fatalf("want 10, got %d", got)
				}
			})
			t.Run("Negative", func(t *testing.T) {
				if got := internal.DecimalDigits(-5); got != 1 {
					t.Fatalf("want 1, got %d", got)
				}
				if got := internal.DecimalDigits(-1234); got != 4 {
					t.Fatalf("want 4, got %d", got)
				}
			})
		})

		t.Run("LooksLikeEpochMillis", func(t *testing.T) {
			t.Run("TrueFor12Or13Digits", func(t *testing.T) {
				if !internal.LooksLikeEpochMillis(999999999999) { // 12 digits
					t.Fatalf("expected true for 12-digit number")
				}
				if !internal.LooksLikeEpochMillis(1_700_000_000_000) { // 13 digits
					t.Fatalf("expected true for 13-digit number")
				}
				if !internal.LooksLikeEpochMillis(-1_700_000_000_000) { // negative 13 digits
					t.Fatalf("expected true for negative 13-digit number")
				}
			})
			t.Run("FalseOtherwise", func(t *testing.T) {
				if internal.LooksLikeEpochMillis(0) {
					t.Fatalf("did not expect true for 0")
				}
				if internal.LooksLikeEpochMillis(1700000000) { // 10 digits (seconds)
					t.Fatalf("did not expect true for 10-digit seconds")
				}
				if internal.LooksLikeEpochMillis(123456789) { // 9 digits
					t.Fatalf("did not expect true for 9 digits")
				}
			})
		})

		t.Run("ToSecondsSigned", func(t *testing.T) {
			t.Run("AcceptsSeconds", func(t *testing.T) {
				got, err := internal.ToSecondsSigned(1_700_000_000) // 10 digits
				if err != nil || got != 1_700_000_000 {
					t.Fatalf("unexpected result: got %d, err=%v", got, err)
				}
			})
			t.Run("RejectsMillis", func(t *testing.T) {
				if _, err := internal.ToSecondsSigned(1_700_000_000_000); err == nil {
					t.Fatalf("expected error for milliseconds input")
				}
				if _, err := internal.ToSecondsSigned(-1_700_000_000_000); err == nil {
					t.Fatalf("expected error for negative milliseconds input")
				}
			})
		})

		t.Run("ToSecondsUnsigned", func(t *testing.T) {
			t.Run("AcceptsSeconds", func(t *testing.T) {
				got, err := internal.ToSecondsUnsigned(1_700_000_000) // 10 digits
				if err != nil || got != 1_700_000_000 {
					t.Fatalf("unexpected result: got %d, err=%v", got, err)
				}
			})
			t.Run("RejectsMillis", func(t *testing.T) {
				if _, err := internal.ToSecondsUnsigned(1_700_000_000_000); err == nil {
					t.Fatalf("expected error for milliseconds-looking input")
				}
			})
			t.Run("RejectsOverflow", func(t *testing.T) {
				overflow := uint64(math.MaxInt64) + 1
				if _, err := internal.ToSecondsUnsigned(overflow); err == nil {
					t.Fatalf("expected overflow error")
				}
			})
		})
	})
}
