// File: common/extension/date/string.go

package internal

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// ParseRFC3339 attempts to parse an RFC3339 timestamp (precise, includes zone).
// Returns (time, true) on success; (zero, false) on failure.
func ParseRFC3339(in string) (time.Time, bool) {
	t, err := time.Parse(time.RFC3339, in)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// ParseYYYYMMDDPrefix parses the first 8 characters as YYYYMMDD,
// requiring all 8 to be digits. Trailing content is ignored.
// Returns UTC midnight on success.
func ParseYYYYMMDDPrefix(raw string) (time.Time, error) {
	in := strings.TrimSpace(raw)
	if len(in) < 8 {
		return time.Time{}, errors.New("date.ParseYYYYMMDDPrefix: input shorter than 8 characters")
	}
	ymd := in[:8]
	if !AllDigits(ymd) {
		return time.Time{}, errors.New("date.ParseYYYYMMDDPrefix: non-digit in YYYYMMDD prefix")
	}

	// Validate ranges & calendar (leap years).
	y, _ := strconv.Atoi(ymd[0:4])
	m, _ := strconv.Atoi(ymd[4:6])
	d, _ := strconv.Atoi(ymd[6:8])

	if m < 1 || m > 12 {
		return time.Time{}, errors.New("date.ParseYYYYMMDDPrefix: invalid month")
	}
	if d < 1 || d > 31 {
		return time.Time{}, errors.New("date.ParseYYYYMMDDPrefix: invalid day")
	}
	if !ValidYMD(y, time.Month(m), d) {
		return time.Time{}, errors.New("date.ParseYYYYMMDDPrefix: invalid calendar date")
	}
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), nil
}

// ParseEpoch parses Unix epochs when input is pure digits:
//   - length 10 → seconds
//   - length 13 → milliseconds
//
// Returns (time, true) on success; (zero, false) otherwise.
func ParseEpoch(s string) (time.Time, bool) {
	if !AllDigits(s) {
		return time.Time{}, false
	}
	switch len(s) {
	case 10: // seconds
		sec, _ := strconv.ParseInt(s, 10, 64)
		return time.Unix(sec, 0).UTC(), true
	case 13: // milliseconds
		ms, _ := strconv.ParseInt(s, 10, 64)
		return time.Unix(ms/1000, (ms%1000)*1_000_000).UTC(), true
	default:
		return time.Time{}, false
	}
}

// ParseZoned parses using a layout that includes zone information
// (e.g., time.RFC1123). On success, preserves the parsed location.
func ParseZoned(layout, in string) (time.Time, bool) {
	t, err := time.Parse(layout, in)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// ParseString applies deterministic string parsing rules used across the package.
// Order:
//  1. RFC3339 (precise, includes zone)
//  2. Pure digits fast-path: epoch seconds/ms, or YYYYMMDD
//  3. Fallback layouts from DefaultLayouts() (date-only, zoneless, zoned)
func ParseString(s string) (time.Time, error) {
	in := strings.TrimSpace(s)
	if in == "" {
		return time.Time{}, errors.New("date.ParseFrom: empty string")
	}

	// 1) RFC3339
	if t, ok := ParseRFC3339(in); ok {
		return t, nil
	}

	// 2) Pure digits
	if AllDigits(in) {
		switch l := len(in); l {
		case 10, 13: // epoch seconds/ms
			if t, ok := ParseEpoch(in); ok {
				return t, nil
			}
		case 8: // YYYYMMDD (with validation)
			// Delegate to the already fully covered YYYYMMDD validator
			// NOTE: ParseYYYYMMDDPrefix handles just the first 8 chars, so plain 8-digit works.
			t, err := ParseYYYYMMDDPrefix(in)
			if err != nil {
				// Keep message style consistent with the prior code
				return time.Time{}, errors.New("date.ParseFrom: invalid YYYYMMDD date")
			}
			return t, nil
		}
	}

	// Fallback layouts (deterministic)
	for _, layout := range DefaultLayouts() {
		switch layout {
		// date-only → normalize to UTC midnight
		case "2006-01-02", "2006/01/02", "02 Jan 2006":
			t, err := time.Parse(layout, in)
			if err != nil {
				continue
			}
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil

		// zoneless → parse in Local, then convert to UTC
		case "2006-01-02 15:04:05", "2006/01/02 15:04:05":
			t, err := time.ParseInLocation(layout, in, time.Local)
			if err != nil {
				continue
			}
			return t.UTC(), nil

		// zoned → preserve zone
		case time.RFC1123:
			if t, ok := ParseZoned(layout, in); ok {
				return t, nil
			}
		}
	}

	return time.Time{}, errors.New("date.ParseFrom: unrecognized string format")
}

// ValidYMD returns true if (y, m, d) is a real calendar date.
func ValidYMD(y int, m time.Month, d int) bool {
	t := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return t.Year() == y && t.Month() == m && t.Day() == d
}
