// File: common/extension/parser_test.go

package extension_test

import (
	"testing"
	"time"

	"github.com/entiqon/entiqon/common/extension"
)

func TestParser(t *testing.T) {
	t.Run("Methods", func(t *testing.T) {

		t.Run("Boolean", func(t *testing.T) {
			type TC struct {
				in   any
				want bool
				ok   bool
			}
			cases := []TC{
				{true, true, true},
				{false, false, true},
				{"true", true, true},
				{"FALSE", false, true},
				{"1", true, true},
				{"0", false, true},
				{"yes", true, true},
				{"no", false, true},
				{"on", true, true},
				{"off", false, true},
				{"y", true, true},
				{"n", false, true},
				{"t", true, true},
				{"f", false, true},
				{"maybe", false, false},
				{nil, false, false},
			}
			for i, tc := range cases {
				got, err := extension.Boolean(tc.in)
				if tc.ok && err != nil {
					t.Fatalf("case %d: unexpected error: %v", i, err)
				}
				if !tc.ok && err == nil {
					t.Fatalf("case %d: expected error, got nil", i)
				}
				if tc.ok && got != tc.want {
					t.Fatalf("case %d: want %v, got %v", i, tc.want, got)
				}
			}
		})

		t.Run("Float", func(t *testing.T) {
			type TC struct {
				in   any
				want float64
				ok   bool
			}
			cases := []TC{
				{3.14159, 3.14159, true},
				{float32(2.5), 2.5, true},
				{int(42), 42.0, true},
				{"42.75", 42.75, true},
				{"  -13.5  ", -13.5, true},
				{"not-a-number", 0, false},
			}
			for i, tc := range cases {
				got, err := extension.Float(tc.in)
				if tc.ok && err != nil {
					t.Fatalf("case %d: unexpected error: %v", i, err)
				}
				if !tc.ok && err == nil {
					t.Fatalf("case %d: expected error, got nil", i)
				}
				if tc.ok && got != tc.want {
					t.Fatalf("case %d: want %v, got %v", i, tc.want, got)
				}
			}
		})

		t.Run("Decimal", func(t *testing.T) {
			type TC struct {
				in        any
				precision int
				want      float64
				ok        bool
			}
			cases := []TC{
				{"3.14159", 2, 3.14, true},
				{2.345, 2, 2.35, true},
				{int64(1234), 0, 1234.0, true},
				{"-0.0049", 2, -0.00, true}, // rounds toward nearest even depending on impl
				{"abc", 2, 0, false},
			}
			for i, tc := range cases {
				got, err := extension.Decimal(tc.in, tc.precision)
				if tc.ok && err != nil {
					t.Fatalf("case %d: unexpected error: %v", i, err)
				}
				if !tc.ok && err == nil {
					t.Fatalf("case %d: expected error, got nil", i)
				}
				if tc.ok && got != tc.want {
					t.Fatalf("case %d: want %v, got %v", i, tc.want, got)
				}
			}
		})

		t.Run("Date", func(t *testing.T) {
			type TC struct {
				in    any
				want  time.Time
				ok    bool
				equal func(a, b time.Time) bool
			}

			// helpers
			ymd := func(y int, m time.Month, d int) time.Time {
				return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
			}
			eq := func(a, b time.Time) bool { return a.Equal(b) }

			cases := []TC{
				// RFC3339
				{"2023-11-14T22:13:20Z", time.Date(2023, 11, 14, 22, 13, 20, 0, time.UTC), true, eq},
				// Epoch seconds / milliseconds
				{"1700000000", time.Unix(1700000000, 0).UTC(), true, eq},
				{"1700000000000", time.Unix(1700000000, 0).UTC(), true, eq},
				// YYYYMMDD
				{"20240229", ymd(2024, 2, 29), true, eq},
				// Fallback layouts (date-only → UTC midnight)
				{"2021-12-31", ymd(2021, 12, 31), true, eq},
				{"2021/12/31", ymd(2021, 12, 31), true, eq},
				{"02 Jan 2006", ymd(2006, 1, 2), true, eq},
				// RFC1123 (zoned)
				{"Tue, 14 Nov 2023 22:13:20 UTC", time.Date(2023, 11, 14, 22, 13, 20, 0, time.FixedZone("UTC", 0)), true, func(a, b time.Time) bool {
					// Allow different UTC location pointers; compare instant & offset
					_, offA := a.Zone()
					_, offB := b.Zone()
					return a.UTC().Equal(b.UTC()) && offA == offB
				}},
				// Invalid
				{"bad-date", time.Time{}, false, eq},
			}

			for i, tc := range cases {
				got, err := extension.Date(tc.in)
				if tc.ok && err != nil {
					t.Fatalf("case %d: unexpected error: %v", i, err)
				}
				if !tc.ok && err == nil {
					t.Fatalf("case %d: expected error, got nil (got=%v)", i, got)
				}
				if tc.ok {
					if tc.equal == nil {
						tc.equal = eq
					}
					if !tc.equal(got, tc.want) {
						t.Fatalf("case %d: want %v, got %v", i, tc.want, got)
					}
				}
			}
		})
	})
}
