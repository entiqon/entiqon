// File: common/extension/date/cleaner.go

package date_test

import (
	"testing"
	"time"

	"github.com/entiqon/entiqon/common/extension/date"
)

func TestCleaner(t *testing.T) {
	t.Run("Methods", func(t *testing.T) {
		t.Run("DefaultCleanParseOptions", func(t *testing.T) {
			opts := date.DefaultCleanParseOptions()
			if opts == nil || len(opts.AcceptLayouts) == 0 || !opts.AcceptEpoch {
				t.Errorf("expected populated defaults")
			}
		})

		t.Run("StrictYYYYMMDDOptions", func(t *testing.T) {
			opts := date.StrictYYYYMMDDOptions()
			if !opts.RequireYYYYMMDDPrefix {
				t.Errorf("expected RequireYYYYMMDDPrefix true")
			}
		})

		t.Run("CleanAndParse", func(t *testing.T) {
			t.Run("EpochSeconds", func(t *testing.T) {
				got, err := date.CleanAndParse("1700000000", nil)
				if err != nil || got.Year() != 2023 {
					t.Errorf("unexpected result: %v", got)
				}
			})

			t.Run("EpochMilliseconds", func(t *testing.T) {
				got, err := date.CleanAndParse("1700000000000", nil)
				if err != nil || got.Year() != 2023 {
					t.Errorf("unexpected result: %v", got)
				}
			})

			t.Run("YYYYMMDDValid", func(t *testing.T) {
				got, err := date.CleanAndParse("20240229", nil)
				if err != nil || got.Year() != 2024 {
					t.Errorf("unexpected result: %v", got)
				}
			})

			t.Run("LayoutDashed", func(t *testing.T) {
				got, err := date.CleanAndParse("2021-12-31", nil)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got.Hour() != 0 {
					t.Errorf("expected midnight UTC")
				}
			})

			t.Run("WithStrictYYYYMMDDPrefix", func(t *testing.T) {
				opts := date.StrictYYYYMMDDOptions()
				got, err := date.CleanAndParse("20250507.0000000000", opts)
				want := time.Date(2025, 5, 7, 0, 0, 0, 0, time.UTC)
				if err != nil || !got.Equal(want) {
					t.Errorf("want %s, got %s, err %v", want, got, err)
				}
			})

			t.Run("StrictInvalid", func(t *testing.T) {
				opts := date.StrictYYYYMMDDOptions()
				if _, err := date.CleanAndParse("2025A507", opts); err == nil {
					t.Errorf("expected error")
				}
			})

			t.Run("EmptyInput", func(t *testing.T) {
				if _, err := date.CleanAndParse("   ", nil); err == nil {
					t.Errorf("expected error")
				}
			})

			t.Run("Zoneless", func(t *testing.T) {

				t.Run("DateTime", func(t *testing.T) {
					t.Run("DefaultsUseUTC", func(t *testing.T) {
						// opts == nil → Location defaults to UTC inside CleanAndParse
						in := "2021-12-31 23:59:59" // matches layout "2006-01-02 15:04:05"
						got, err := date.CleanAndParse(in, nil)
						if err != nil {
							t.Fatalf("unexpected error: %v", err)
						}
						want := time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC)
						if !got.Equal(want) {
							t.Fatalf("want %s, got %s", want, got)
						}
					})

					t.Run("CustomLocationHonored", func(t *testing.T) {
						// Provide a non-nil Location and ensure it's used before converting to UTC
						opts := date.DefaultCleanParseOptions()
						opts.Location = time.FixedZone("X", -5*3600) // UTC-05:00

						in := "2021-12-31 23:00:00" // interpreted as 23:00 in UTC-05
						got, err := date.CleanAndParse(in, opts)
						if err != nil {
							t.Fatalf("unexpected error: %v", err)
						}
						// 23:00 at UTC-05 is 04:00 UTC next day
						want := time.Date(2022, 1, 1, 4, 0, 0, 0, time.UTC)
						if !got.Equal(want) {
							t.Fatalf("want %s, got %s", want, got)
						}
					})

					t.Run("LocationNilDefaultsToUTC", func(t *testing.T) {
						// Build options with Location == nil (not using DefaultCleanParseOptions).
						opts := &date.CleanParseOptions{
							RequireYYYYMMDDPrefix: false,
							AcceptLayouts:         date.DefaultLayouts(), // includes "2006-01-02 15:04:05"
							AcceptEpoch:           true,
							Location:              nil, // <-- this is key to cover the branch
						}

						in := "2021-12-31 23:59:59" // zoneless layout
						got, err := date.CleanAndParse(in, opts)
						if err != nil {
							t.Fatalf("unexpected error: %v", err)
						}

						// When Location is nil, cleaner should default to UTC, then convert to UTC (no change).
						want := time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC)
						if !got.Equal(want) {
							t.Fatalf("want %s, got %s", want, got)
						}
					})
				})

			})

			t.Run("RFC1123", func(t *testing.T) {
				t.Run("Valid", func(t *testing.T) {
					got, err := date.CleanAndParse("2023-11-14T22:13:20Z", nil)
					if err != nil || got.IsZero() {
						t.Errorf("unexpected failure: %v %v", got, err)
					}
				})

				t.Run("Zoned", func(t *testing.T) {
					t.Run("LayoutPreserved", func(t *testing.T) {
						// This string is RFC1123 (NOT RFC3339), so CleanAndParse will skip
						// RFC3339/epoch/8-digit paths and reach the layouts loop → time.RFC1123 case.
						in := "Tue, 14 Nov 2023 22:13:20 UTC"

						got, err := date.CleanAndParse(in, nil) // uses DefaultCleanParseOptions → includes RFC1123 in layouts
						if err != nil {
							t.Fatalf("unexpected error: %v", err)
						}

						// Should preserve the parsed zone (UTC) because ParseZoned is used.
						if got.Location() == time.UTC {
							// ok
						} else if got.Location().String() == "UTC" {
							// some Go builds name the location "UTC" rather than pointer-equal time.UTC; still fine
						} else {
							t.Fatalf("expected UTC location, got %q", got.Location().String())
						}

						// Sanity-check fields
						if got.Year() != 2023 || got.Month() != time.November || got.Day() != 14 ||
							got.Hour() != 22 || got.Minute() != 13 || got.Second() != 20 {
							t.Fatalf("unexpected parsed value: %v", got)
						}
					})
				})
			})
		})

		t.Run("CleanAndParseAsString", func(t *testing.T) {

			t.Run("DefaultLayout", func(t *testing.T) {
				got := date.CleanAndParseAsString("2023-11-14T22:13:20Z", "")
				if got != "2023-11-14" {
					t.Errorf("unexpected result: %s", got)
				}
			})

			t.Run("CustomLayout", func(t *testing.T) {
				got := date.CleanAndParseAsString("2023-11-14T22:13:20Z", "20060102")
				if got != "20231114" {
					t.Errorf("unexpected result: %s", got)
				}
			})

			t.Run("InvalidInput", func(t *testing.T) {
				got := date.CleanAndParseAsString("bad", "2006-01-02")
				if got != "" {
					t.Errorf("expected empty string, got %q", got)
				}
			})

		})
	})
}
