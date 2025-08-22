// File: common/extension/date/parser.go

// Package date provides lightweight parsing utilities to convert heterogeneous
// inputs (strings, numbers, and time types) into time.Time values.
//
// Design goals:
//  1. Deterministic behavior: ISO-first parsing with explicit, documented fallbacks.
//  2. Minimal surprises: date-only inputs are normalized to 00:00:00 in UTC.
//  3. Zero reflection on the hot path: explicit type switches for performance and clarity.
//
// Supported inputs
//
//   - time.Time / *time.Time: returned as-is (nil pointers are rejected).
//   - string / *string:
//   - RFC3339 (e.g., "2025-08-16T13:45:00Z")
//   - ISO date "2006-01-02" and "2006/01/02" (normalized to 00:00:00 UTC)
//   - "2006-01-02 15:04:05" and "2006/01/02 15:04:05"
//   - "02 Jan 2006" (normalized to 00:00:00 UTC)
//   - RFC1123 ("Mon, 02 Jan 2006 15:04:05 MST")
//   - Pure digits:
//     · 10 digits → epoch seconds
//     · 13 digits → epoch milliseconds
//     · 8 digits  → YYYYMMDD (normalized to 00:00:00 UTC)
//     Unrecognized strings produce a descriptive error.
//   - Integers and unsigned integers of any width → epoch seconds in UTC.
//   - Floats (float32/float64) → epoch seconds in UTC; fractional part becomes nanoseconds.
//
// # Time zones
//
// Formats that embed a zone (e.g., RFC3339, RFC1123) preserve that zone. For
// date-only inputs (no time component) the result is normalized to midnight
// in UTC to avoid local time ambiguities.
//
// # Errors
//
// ParseFrom returns an error when the input is nil, a nil pointer for supported
// pointer types, an unsupported type, or a string that does not match any
// recognized format.
//
// Examples are provided in example_test.go.
package date

import (
	"errors"
	"math"
	"time"

	"github.com/entiqon/entiqon/common/extension/date/internal"
)

// ParseFrom attempts to parse value into a time.Time.
//
// Accepted types and behaviors:
//
//   - time.Time: returned as-is.
//   - *time.Time: dereferenced and returned; nil pointer → error.
//   - string: tried against a deterministic set of layouts in the following order:
//     1) RFC3339: time.RFC3339 (e.g., "2025-08-16T13:45:00Z").
//     2) Pure digits:
//   - 10 digits → epoch seconds (UTC)
//   - 13 digits → epoch milliseconds (UTC)
//   - 8 digits  → YYYYMMDD, normalized to 00:00:00 UTC
//     3) Fallback layouts (in order):
//   - "2006-01-02"                 → normalized to 00:00:00 UTC
//   - "2006/01/02"                 → normalized to 00:00:00 UTC
//   - "2006-01-02 15:04:05"
//   - "2006/01/02 15:04:05"
//   - "02 Jan 2006"                → normalized to 00:00:00 UTC
//   - time.RFC1123 ("Mon, 02 Jan 2006 15:04:05 MST")
//     Unrecognized strings produce an error.
//   - *string: dereferenced and treated as string; nil pointer → error.
//   - Integers/unsigned (any width): interpreted as epoch seconds (UTC).
//   - float32/float64: interpreted as epoch seconds (UTC) with fractional part
//     converted to nanoseconds.
//
// # Time zone handling
//
// If the successfully matched layout includes an offset/zone (e.g., RFC3339,
// RFC1123), that zone is preserved. For date-only inputs, the time is normalized
// to midnight (00:00:00) in UTC to avoid local time ambiguities.
//
// Parameters:
//   - value: The input to parse.
//
// Returns:
//   - time.Time: The parsed time value.
//   - error: A descriptive error if parsing fails or the type is unsupported.
func ParseFrom(value any) (time.Time, error) {
	if value == nil {
		return time.Time{}, errors.New("date.ParseFrom: value is nil")
	}

	switch v := value.(type) {
	case time.Time:
		return v, nil

	case *time.Time:
		if v == nil {
			return time.Time{}, errors.New("date.ParseFrom: *time.Time is nil")
		}
		return *v, nil

	case string:
		return internal.ParseString(v)

	case *string:
		if v == nil {
			return time.Time{}, errors.New("date.ParseFrom: *string is nil")
		}
		return internal.ParseString(*v)

	// Integers/unsigned → epoch seconds (UTC)
	case int:
		sec, err := internal.ToSecondsSigned(int64(v))
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(sec, 0).UTC(), nil

		// Signed (small)
	case int8:
		sec := int64(v) // cannot look like ms by range
		return time.Unix(sec, 0).UTC(), nil
	case int16:
		sec := int64(v)
		return time.Unix(sec, 0).UTC(), nil
	case int32:
		sec := int64(v)
		return time.Unix(sec, 0).UTC(), nil

	case int64:
		sec, err := internal.ToSecondsSigned(v) // has ms-guard
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(sec, 0).UTC(), nil

		// Unsigned (small)
	case uint8:
		sec := int64(v)
		return time.Unix(sec, 0).UTC(), nil
	case uint16:
		sec := int64(v)
		return time.Unix(sec, 0).UTC(), nil
	case uint32:
		sec := int64(v)
		return time.Unix(sec, 0).UTC(), nil

		// Unsigned (platform + wide)
	case uint:
		sec, err := internal.ToSecondsUnsigned(uint64(v)) // on 64-bit, ms-guard can fire
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(sec, 0).UTC(), nil
	case uint64:
		sec, err := internal.ToSecondsUnsigned(v) // overflow + ms-guard
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(sec, 0).UTC(), nil
	case uintptr:
		sec, err := internal.ToSecondsUnsigned(uint64(v)) // overflow + ms-guard (64-bit)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(sec, 0).UTC(), nil

	// Floats → epoch seconds + fractional nanoseconds (UTC)
	case float32:
		// float32 cannot represent sub-second precision near 1e9 seconds.
		// We intentionally drop the fractional part to avoid misleading results.
		sec := int64(v)
		return time.Unix(sec, 0).UTC(), nil

	case float64:
		// Define decimal semantics at millisecond precision:
		// round to nearest millisecond, then split to sec + nanoSecs.
		totalMillis := math.Round(v * 1e3)           // decimal rounding to ms
		totalNanos := int64(totalMillis) * 1_000_000 // exact integer nanos

		sec := totalNanos / 1_000_000_000
		nanoSecs := totalNanos % 1_000_000_000
		if nanoSecs < 0 { // normalize negatives into [0, 1e9)
			sec--
			nanoSecs += 1_000_000_000
		}
		return time.Unix(sec, nanoSecs).UTC(), nil
	}

	return time.Time{}, errors.New("date.ParseFrom: unsupported type")
}

// ParseAndFormat attempts to parse the input date string `value`
// using the specified `layout`. If `layout` is empty, it defaults to "2006-01-02".
//
// It tries to parse the input using multiple layouts, including the requested layout and common fallbacks:
//
//   - "2006-01-02" (ISO with dashes)
//   - "2006/01/02" (ISO with slashes)
//   - "20060102"   (compact numeric format)
//
// Upon successful parsing with any layout, the function reformats the date into
// the requested `layout` and returns the formatted string.
//
// If parsing fails for all layouts, it returns an empty string.
//
// Parameters:
//   - value: the date string to parse.
//   - layout: the desired output date format layout (e.g., "2006-01-02").
//     Defaults to "2006-01-02" if empty.
//
// Returns:
//   - string: the date formatted according to `layout`, or empty string if parsing fails.
func ParseAndFormat(value, layout string) string {
	return CleanAndParseAsString(value, layout)
}
