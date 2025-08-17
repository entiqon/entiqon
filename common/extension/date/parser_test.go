// File: common/extension/date/parser_test.go

package date_test

import (
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/entiqon/entiqon/common/extension/date"
	"github.com/entiqon/entiqon/common/extension/date/internal"
)

func TestDate(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var x any
		_, err := date.ParseFrom(x)
		if err == nil {
			t.Fatalf("ParseFrom(nil) expected error, got nil")
		}
	})

	t.Run("float", func(t *testing.T) {
		t.Run("32", func(t *testing.T) {
			t.Run("EpochSecondsFractionDropped", func(t *testing.T) {
				got, err := date.ParseFrom(float32(1_700_000_000.123))
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if want := time.Unix(1_700_000_000, 0).UTC(); !got.Equal(want) {
					t.Fatalf("expected %v, got %v", want, got)
				}
			})
		})

		t.Run("64", func(t *testing.T) {
			t.Run("EpochSeconds", func(t *testing.T) {
				exp := time.Unix(1_700_000_000, 123_000_000).UTC()
				got, err := date.ParseFrom(1_700_000_000.123)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !got.Equal(exp) {
					t.Fatalf("expected %v, got %v", exp, got)
				}
			})
		})
	})

	t.Run("float32", func(t *testing.T) {
		t.Run("EpochSecondsFractionDropped_Positive", func(t *testing.T) {
			in := float32(1_700_000_000.123)
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			want := time.Unix(1_700_000_000, 0).UTC()
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("EpochSecondsFractionDropped_Negative", func(t *testing.T) {
			in := float32(-1.987)
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			want := time.Unix(-1, 0).UTC()
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})
	})

	t.Run("float64", func(t *testing.T) {
		t.Run("EpochSecondsRoundedToMillis", func(t *testing.T) {
			in := 1_700_000_000.123 // -> 123 ms
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			want := time.Unix(1_700_000_000, 123_000_000).UTC()
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("RoundingHalfUpToNextMs", func(t *testing.T) {
			// Half up at ms precision should round to 124 ms
			in := 1_700_000_000.1235
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			want := time.Unix(1_700_000_000, 124_000_000).UTC()
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("CarryToNextSecond", func(t *testing.T) {
			// 0.9995 s → 1.000 s exactly after ms rounding
			in := 1_700_000_000.9995
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			want := time.Unix(1_700_000_001, 0).UTC()
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("NegativeNoMsChange", func(t *testing.T) {
			// -1.0004 → -1.000 after ms rounding (no carry)
			in := -1.0004
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			want := time.Unix(-1, 0).UTC()
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("NegativeRoundsAwayFromZeroWithNormalize", func(t *testing.T) {
			// -1.0005 → -1.001 s; normalize to sec=-2, nsec=999_000_000
			in := -1.0005
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			want := time.Unix(-2, 999_000_000).UTC()
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})
	})

	t.Run("int", func(t *testing.T) {
		t.Run("NumDigits", func(t *testing.T) {
			t.Run("Zero", func(t *testing.T) {
				got := internal.DecimalDigits(0)
				if got != 1 {
					t.Fatalf("expected 1 for input 0, got %d", got)
				}
			})

			t.Run("OneDigit", func(t *testing.T) {
				got := internal.DecimalDigits(7)
				if got != 1 {
					t.Fatalf("expected 1, got %d", got)
				}
			})

			t.Run("TwoDigits", func(t *testing.T) {
				got := internal.DecimalDigits(42)
				if got != 2 {
					t.Fatalf("expected 2, got %d", got)
				}
			})
		})
	})

	t.Run("int8", func(t *testing.T) {
		t.Run("SecondsPositive", func(t *testing.T) {
			in := int8(42)
			want := time.Unix(int64(in), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("SecondsNegative", func(t *testing.T) {
			in := int8(-7)
			want := time.Unix(int64(in), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})
	})

	t.Run("int16", func(t *testing.T) {
		t.Run("SecondsPositive", func(t *testing.T) {
			in := int16(12345)
			want := time.Unix(int64(in), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("SecondsNegative", func(t *testing.T) {
			in := int16(-12345)
			want := time.Unix(int64(in), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})
	})

	t.Run("int32", func(t *testing.T) {
		t.Run("SecondsPositive", func(t *testing.T) {
			in := int32(2_147_000_000) // near max
			want := time.Unix(int64(in), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("SecondsNegative", func(t *testing.T) {
			in := int32(-2_147_000_000)
			want := time.Unix(int64(in), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})
	})

	t.Run("int64", func(t *testing.T) {
		t.Run("SecondsPositive", func(t *testing.T) {
			in := int64(1_700_000_000)
			want := time.Unix(in, 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("SecondsNegative", func(t *testing.T) {
			in := int64(-1_000)
			want := time.Unix(in, 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("MillisecondsGuard", func(t *testing.T) {
			_, err := date.ParseFrom(int64(1_734_412_800_123)) // 13 digits
			if err == nil {
				t.Fatalf("expected ms-guard error")
			}
		})
	})

	t.Run("uint", func(t *testing.T) {
		t.Run("Seconds", func(t *testing.T) {
			in := uint(1_600_000_000)
			want := time.Unix(int64(in), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("MillisecondsGuard_OnlyOn64bit", func(t *testing.T) {
			if strconv.IntSize == 64 {
				in := uint(1_734_412_800_123)
				_, err := date.ParseFrom(in)
				if err == nil {
					t.Fatalf("expected ms-guard error")
				}
			}
		})
	})

	t.Run("uint8", func(t *testing.T) {
		t.Run("Seconds", func(t *testing.T) {
			in := uint8(200)
			want := time.Unix(int64(in), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})
	})

	t.Run("uint16", func(t *testing.T) {
		t.Run("Seconds", func(t *testing.T) {
			in := uint16(65000)
			want := time.Unix(int64(in), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})
	})

	t.Run("uint32", func(t *testing.T) {
		t.Run("Seconds", func(t *testing.T) {
			in := uint32(4_000_000_000)
			want := time.Unix(int64(in), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})
	})

	t.Run("uint64", func(t *testing.T) {
		t.Run("SecondsBorderlineMax", func(t *testing.T) {
			got, err := date.ParseFrom(uint64(math.MaxInt64))
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got.Unix() != math.MaxInt64 {
				t.Fatalf("unexpected seconds: %d", got.Unix())
			}
		})

		t.Run("MillisecondsGuard", func(t *testing.T) {
			_, err := date.ParseFrom(uint64(1_734_412_800_123))
			if err == nil {
				t.Fatalf("expected ms-guard error")
			}
		})

		t.Run("Overflow", func(t *testing.T) {
			_, err := date.ParseFrom(uint64(math.MaxInt64) + 1)
			if err == nil {
				t.Fatalf("expected overflow error")
			}
		})
	})

	t.Run("uintptr", func(t *testing.T) {
		// Always hits the success return (both 32/64-bit)
		t.Run("SecondsSmallValue", func(t *testing.T) {
			in := uintptr(42)
			want := time.Unix(int64(uint64(in)), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		// Another clean success: below ms-guard threshold (13 digits) and far from overflow
		t.Run("Seconds_NormalEpoch", func(t *testing.T) {
			in := uintptr(1_600_000_000) // 2020-09-13T12:26:40Z
			want := time.Unix(int64(uint64(in)), 0).UTC()
			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		// 64-bit only: success at the MaxInt64 boundary
		if strconv.IntSize == 64 {
			t.Run("Seconds_MaxInt64Boundary", func(t *testing.T) {
				in := uintptr(uint64(math.MaxInt64))
				got, err := date.ParseFrom(in)
				if err != nil {
					t.Fatalf("unexpected: %v", err)
				}
				if got.Unix() != math.MaxInt64 {
					t.Fatalf("unexpected seconds: %d", got.Unix())
				}
			})
		}
	})

	t.Run("integer", func(t *testing.T) {
		t.Run("SecondsPositive", func(t *testing.T) {
			// Safe for 32-bit and 64-bit
			in := 1_600_000_000 // 2020-09-13T12:26:40Z
			want := time.Unix(int64(in), 0).UTC()

			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("SecondsNegative", func(t *testing.T) {
			in := -1_000 // 1969-12-31T23:43:20Z
			want := time.Unix(int64(in), 0).UTC()

			got, err := date.ParseFrom(in)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("MillisecondsGuard", func(t *testing.T) {
			// Only meaningful if 'int' is 64-bit; otherwise constructing the value would overflow.
			if strconv.IntSize == 64 {
				in := 1_734_412_800_123 // 13 digits → "looks like ms"
				_, err := date.ParseFrom(in)
				if err == nil {
					t.Fatalf("expected ms-guard error for int with 13 digits")
				}
			}
		})

		t.Run("IntWithPositiveSeconds", func(t *testing.T) {
			want := time.Unix(1_700_000_000, 0).UTC()
			got, err := date.ParseFrom(1_700_000_000)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("Int64WithNegativeSeconds", func(t *testing.T) {
			want := time.Unix(-1_000, 0).UTC()
			got, err := date.ParseFrom(int64(-1_000))
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !got.Equal(want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})

		t.Run("Uint64WithOverflowToInt64", func(t *testing.T) {
			_, err := date.ParseFrom(uint64(math.MaxInt64) + 1)
			if err == nil {
				t.Fatalf("expected overflow error")
			}
		})

		t.Run("Uint64WithSecondsBorderlineMax", func(t *testing.T) {
			got, err := date.ParseFrom(uint64(math.MaxInt64))
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got.Unix() != math.MaxInt64 {
				t.Fatalf("unexpected seconds: %d", got.Unix())
			}
		})

		t.Run("Uintptr/Overflow", func(t *testing.T) {
			// Simulate by constructing a value > MaxInt64 when uintptr is 64-bit.
			// If your arch is 64-bit, this test makes sense; otherwise skip.
			if strconv.IntSize == 64 {
				uv := uint64(math.MaxInt64) + 1
				p := uintptr(uv) // NOTE: on some platforms this wraps; if so, skip test.
				if uint64(p) == uv {
					_, err := date.ParseFrom(p)
					if err == nil {
						t.Fatalf("expected overflow error")
					}
				}
			}
		})

		t.Run("LooksLikeMillis_Signed", func(t *testing.T) {
			t.Run("Signed", func(t *testing.T) {
				// 13 digits (ms) should be rejected for ints
				_, err := date.ParseFrom(int64(1_734_412_800_123))
				if err == nil {
					t.Fatalf("expected ms-accident error")
				}
			})

			t.Run("Unsigned", func(t *testing.T) {
				_, err := date.ParseFrom(uint64(1_734_412_800_123))
				if err == nil {
					t.Fatalf("expected ms-accident error")
				}
			})
		})

		t.Run("EpochSeconds", func(t *testing.T) {
			exp := time.Unix(1_700_000_000, 0).UTC()
			got, err := date.ParseFrom(int64(1_700_000_000))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(exp) {
				t.Fatalf("expected %v, got %v", exp, got)
			}
		})
	})

	t.Run("time.Time", func(t *testing.T) {
		t.Run("PointerTimeNil", func(t *testing.T) {
			var tp *time.Time = nil
			var x any = tp
			_, err := date.ParseFrom(x)
			if err == nil {
				t.Fatalf("expected error for *time.Time(nil)")
			}
		})

		t.Run("Time", func(t *testing.T) {
			now := time.Now().UTC()
			got, err := date.ParseFrom(now)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(now) {
				t.Fatalf("expected %v, got %v", now, got)
			}
		})

		t.Run("PointerTime", func(t *testing.T) {
			now := time.Now().UTC()
			got, err := date.ParseFrom(&now)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(now) {
				t.Fatalf("expected %v, got %v", now, got)
			}
		})
	})

	t.Run("string", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			_, err := date.ParseFrom("")
			if err == nil {
				t.Fatalf("expected empty string error")
			}
		})

		t.Run("Digits", func(t *testing.T) {
			t.Run("EpochSeconds10Digits", func(t *testing.T) {
				in := "1700000000" // 10 digits → epoch seconds
				got, err := date.ParseFrom(in)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				want := time.Unix(1_700_000_000, 0).UTC()
				if !got.Equal(want) {
					t.Fatalf("want %v, got %v", want, got)
				}
			})

			t.Run("EpochSecondsZeroWithLeadingZeros", func(t *testing.T) {
				in := "0000000000" // still 10 digits, parses to 0 seconds
				got, err := date.ParseFrom(in)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				want := time.Unix(0, 0).UTC()
				if !got.Equal(want) {
					t.Fatalf("want %v, got %v", want, got)
				}
			})

			t.Run("13Digits", func(t *testing.T) {
				t.Run("EpochMilliseconds", func(t *testing.T) {
					in := "1734412800123" // → 2024-12-17T00:00:00.123Z
					got, err := date.ParseFrom(in)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					want := time.UnixMilli(1_734_412_800_123).UTC()
					if !got.Equal(want) {
						t.Fatalf("want %v, got %v", want, got)
					}
				})

				t.Run("EpochMillisecondsTrimmed", func(t *testing.T) {
					in := "   1734412800123   "
					got, err := date.ParseFrom(in)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					want := time.UnixMilli(1_734_412_800_123).UTC()
					if !got.Equal(want) {
						t.Fatalf("want %v, got %v", want, got)
					}
				})

				t.Run("EpochMillisecondsZero", func(t *testing.T) {
					in := "0000000000000" // 13 zeros → epoch 0 in ms
					got, err := date.ParseFrom(in)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					want := time.Unix(0, 0).UTC()
					if !got.Equal(want) {
						t.Fatalf("want %v, got %v", want, got)
					}
				})
			})

			t.Run("Pointer13Digits", func(t *testing.T) {
				t.Run("EpochMilliseconds", func(t *testing.T) {
					s := "1734412800123"
					var x any = &s
					got, err := date.ParseFrom(x)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					want := time.UnixMilli(1_734_412_800_123).UTC()
					if !got.Equal(want) {
						t.Fatalf("want %v, got %v", want, got)
					}
				})
			})

		})

		t.Run("OnlySpaces", func(t *testing.T) {
			_, err := date.ParseFrom("   \t\n  ")
			if err == nil {
				t.Fatalf("expected empty string error after trim")
			}
		})

		t.Run("Pointer", func(t *testing.T) {
			t.Run("Empty", func(t *testing.T) {
				s := ""
				var x any = &s
				_, err := date.ParseFrom(x)
				if err == nil {
					t.Fatalf("expected empty string error via *string")
				}
			})

			t.Run("OnlySpaces", func(t *testing.T) {
				s := "  \n\t  "
				var x any = &s
				_, err := date.ParseFrom(x)
				if err == nil {
					t.Fatalf("expected empty string error via *string after trim")
				}
			})

			t.Run("Nil", func(t *testing.T) {
				var sp *string = nil
				var x any = sp
				_, err := date.ParseFrom(x)
				if err == nil {
					t.Fatalf("expected error for *string(nil)")
				}
			})

			t.Run("Valid", func(t *testing.T) {
				s := "2025-08-16T13:45:00Z"
				var x any = &s
				got, err := date.ParseFrom(x)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				want := time.Date(2025, 8, 16, 13, 45, 0, 0, time.UTC)
				if !got.Equal(want) {
					t.Fatalf("want %v, got %v", want, got)
				}
			})

			t.Run("EpochSeconds10Digits", func(t *testing.T) {
				s := "1700000000"
				var x any = &s
				got, err := date.ParseFrom(x)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				want := time.Unix(1_700_000_000, 0).UTC()
				if !got.Equal(want) {
					t.Fatalf("want %v, got %v", want, got)
				}
			})

			t.Run("YYYYMMDD", func(t *testing.T) {
				t.Run("InvalidDayRange", func(t *testing.T) {
					t.Run("DayZero", func(t *testing.T) {
						s := "20251200" // Dec 00
						var x any = &s
						_, err := date.ParseFrom(x)
						if err == nil {
							t.Fatalf("expected invalid YYYYMMDD day error via *string")
						}
					})

					t.Run("DayTooHigh", func(t *testing.T) {
						s := "20251133" // Nov 33
						var x any = &s
						_, err := date.ParseFrom(x)
						if err == nil {
							t.Fatalf("expected invalid YYYYMMDD day error via *string")
						}
					})
				})

				t.Run("Valid", func(t *testing.T) {
					s := "19991231" // Dec 31, 1999
					var x any = &s
					got, err := date.ParseFrom(x)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					want := time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)
					if !got.Equal(want) {
						t.Fatalf("want %v, got %v", want, got)
					}
				})
			})
		})

		t.Run("RFC3339", func(t *testing.T) {
			exp := time.Date(2025, 8, 16, 13, 45, 0, 0, time.UTC)
			got, err := date.ParseFrom("2025-08-16T13:45:00Z")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(exp) {
				t.Fatalf("expected %v, got %v", exp, got)
			}
		})

		t.Run("ISODate", func(t *testing.T) {
			exp := time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC)
			got, err := date.ParseFrom("2025-08-16")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(exp) {
				t.Fatalf("expected %v, got %v", exp, got)
			}
		})

		t.Run("YYYYMMDDInvalidMonth", func(t *testing.T) {
			_, err := date.ParseFrom("20251301")
			if err == nil {
				t.Fatalf("expected invalid YYYYMMDD month")
			}
		})

		t.Run("YYYYMMDD", func(t *testing.T) {
			t.Run("InvalidDay", func(t *testing.T) {
				_, err := date.ParseFrom("20250230")
				if err == nil {
					t.Fatalf("expected invalid YYYYMMDD date")
				}
			})

			t.Run("InvalidYearRange", func(t *testing.T) {
				t.Run("DayZero", func(t *testing.T) {
					_, err := date.ParseFrom("20250100") // Jan 00 → triggers d < 1
					if err == nil {
						t.Fatalf("expected invalid YYYYMMDD day error")
					}
				})

				t.Run("DayTooHigh", func(t *testing.T) {
					_, err := date.ParseFrom("20250132") // Jan 32 → triggers d > 31
					if err == nil {
						t.Fatalf("expected invalid YYYYMMDD day error")
					}
				})
			})

			t.Run("Valid", func(t *testing.T) {
				t.Run("NormalDate", func(t *testing.T) {
					in := "20250217" // 2025-02-17
					got, err := date.ParseFrom(in)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					want := time.Date(2025, 2, 17, 0, 0, 0, 0, time.UTC)
					if !got.Equal(want) {
						t.Fatalf("want %v, got %v", want, got)
					}
				})

				t.Run("LeapYear", func(t *testing.T) {
					in := "20240229" // 29th Feb 2024 is valid
					got, err := date.ParseFrom(in)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					want := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
					if !got.Equal(want) {
						t.Fatalf("want %v, got %v", want, got)
					}
				})
			})
		})

		t.Run("DateTimeLocal", func(t *testing.T) {
			t.Run("NoZone", func(t *testing.T) {
				t.Run("ConvertedToUTC", func(t *testing.T) {
					// Set a deterministic Local (UTC-05:00 for example)
					origLocal := time.Local
					time.Local = time.FixedZone("TEST-0500", -5*60*60)
					defer func() { time.Local = origLocal }()

					in := "2025-08-16 13:45:05" // interpreted as Local (TEST-0500)
					got, err := date.ParseFrom(in)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					// 13:45:05 at -05:00 → 18:45:05 UTC
					want := time.Date(2025, 8, 16, 18, 45, 5, 0, time.UTC)
					if !got.Equal(want) {
						t.Fatalf("want %v, got %v", want, got)
					}
				})

				t.Run("SlashConvertedToUTC", func(t *testing.T) {
					origLocal := time.Local
					time.Local = time.FixedZone("TEST-0530", -5*60*60-30*60) // UTC-05:30
					defer func() { time.Local = origLocal }()

					in := "2025/08/16 00:00:00"
					got, err := date.ParseFrom(in)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					// 00:00:00 at -05:30 → 05:30:00 UTC
					want := time.Date(2025, 8, 16, 5, 30, 0, 0, time.UTC)
					if !got.Equal(want) {
						t.Fatalf("want %v, got %v", want, got)
					}
				})
			})

			t.Run("RFC3339", func(t *testing.T) {
				t.Run("WithZone", func(t *testing.T) {
					t.Run("PreserveZoneNoUTCConversion", func(t *testing.T) {
						// RFC1123 carries a zone. We should NOT force UTC here; code returns t as-is.
						in := "Sat, 16 Aug 2025 13:45:05 GMT"
						got, err := date.ParseFrom(in)
						if err != nil {
							t.Fatalf("unexpected error: %v", err)
						}
						// GMT == UTC; equality still OK, but Location is not time.Local, so branch returns t directly
						want := time.Date(2025, 8, 16, 13, 45, 5, 0, time.FixedZone("GMT", 0))
						if got.Year() != want.Year() || got.Month() != want.Month() || got.Day() != want.Day() ||
							got.Hour() != want.Hour() || got.Minute() != want.Minute() || got.Second() != want.Second() {
							t.Fatalf("want %v, got %v", want, got)
						}
						// Optional: assert zone name (not strictly required):
						if got.Location().String() != "GMT" {
							t.Fatalf("expected GMT zone, got %s", got.Location().String())
						}
					})
				})
			})
		})

		t.Run("Unrecognized", func(t *testing.T) {
			t.Run("GarbageText", func(t *testing.T) {
				_, err := date.ParseFrom("definitely-not-a-date")
				if err == nil {
					t.Fatalf("expected unrecognized string format error")
				}
			})

			t.Run("TSeparatorWithoutZone", func(t *testing.T) {
				// Looks like ISO but no timezone → not RFC3339, not zoneless layouts
				_, err := date.ParseFrom("2025-08-17T14:30:00")
				if err == nil {
					t.Fatalf("expected unrecognized string format error")
				}
			})

			t.Run("IncompleteDate", func(t *testing.T) {
				// 7 digits → fails YYYYMMDD, also not in any known format
				_, err := date.ParseFrom("2025131")
				if err == nil {
					t.Fatalf("expected unrecognized string format error")
				}
			})
		})
	})

	t.Run("Unsupported", func(t *testing.T) {
		t.Run("Struct", func(t *testing.T) {
			type X struct{ A int }
			_, err := date.ParseFrom(X{A: 1})
			if err == nil {
				t.Fatalf("expected unsupported type error")
			}
		})

		// Optional extra shapes (still the same line, but nice sanity checks)
		t.Run("Slice", func(t *testing.T) {
			_, err := date.ParseFrom([]byte{1, 2, 3})
			if err == nil {
				t.Fatalf("expected unsupported type error")
			}
		})

		t.Run("Map", func(t *testing.T) {
			_, err := date.ParseFrom(map[string]int{"a": 1})
			if err == nil {
				t.Fatalf("expected unsupported type error")
			}
		})

		t.Run("Func", func(t *testing.T) {
			_, err := date.ParseFrom(func() {})
			if err == nil {
				t.Fatalf("expected unsupported type error")
			}
		})
	})
}

func TestParseAndFormatDate(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		layout string
		want   string
	}{
		{
			name:   "Valid date matches layout, returns unchanged",
			value:  "2025-08-01",
			layout: "2006-01-02",
			want:   "2025-08-01",
		},
		{
			name:   "Empty layout defaults to 2006-01-02",
			value:  "2025-08-01",
			layout: "",
			want:   "2025-08-01",
		},
		{
			name:   "Parse fallback layout with slashes",
			value:  "2025/08/01",
			layout: "2006-01-02",
			want:   "2025-08-01",
		},
		{
			name:   "Parse fallback layout compact",
			value:  "20250801",
			layout: "2006-01-02",
			want:   "2025-08-01",
		},
		{
			name:   "Invalid date returns empty string",
			value:  "not-a-date",
			layout: "2006-01-02",
			want:   "",
		},
		{
			name:   "Empty value returns empty string",
			value:  "",
			layout: "2006-01-02",
			want:   "",
		},
		{
			name:   "Different output layout",
			value:  "2025-08-01",
			layout: "02 Jan 2006",
			want:   "01 Aug 2025",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := date.ParseAndFormat(tc.value, tc.layout)
			if got != tc.want {
				t.Errorf("ParseAndFormatDate() = %q, want %q", got, tc.want)
			}
		})
	}
}
