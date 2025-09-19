// File: common/extension/date/cleaner.go

package date

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/entiqon/common/extension/date/internal"
)

// CleanParseOptions controls CleanAndParse behavior.
type CleanParseOptions struct {
	// If true, force parsing only the first 8 chars as YYYYMMDD (digits only),
	// ignoring any trailing content.
	RequireYYYYMMDDPrefix bool

	// Layouts to attempt in order after RFC3339/epoch. If nil/empty, callers
	// should use DefaultCleanParseOptions() to populate sensible defaults.
	AcceptLayouts []string

	// If true, allow epoch seconds (10 digits) and milliseconds (13 digits).
	AcceptEpoch bool

	// Location used for zoneless datetime layouts (e.g., "2006-01-02 15:04:05").
	// If nil, defaults to UTC.
	Location *time.Location
}

// CleanAndParse normalizes and parses a date string according to options.
// Errors are wrapped with "date.CleanAndParse:" for deterministic namespacing.
func CleanAndParse(raw string, opts *CleanParseOptions) (time.Time, error) {
	if opts == nil {
		opts = DefaultCleanParseOptions()
	}
	in := strings.TrimSpace(raw)
	if in == "" {
		return time.Time{}, errors.New("date.CleanAndParse: empty input")
	}

	// Strict prefix mode short-circuits.
	if opts.RequireYYYYMMDDPrefix {
		t, err := internal.ParseYYYYMMDDPrefix(in)
		if err != nil {
			return time.Time{}, fmt.Errorf("date.CleanAndParse: %w", err)
		}
		return t, nil
	}

	// RFC3339 (precise, with zone)
	if t, ok := internal.ParseRFC3339(in); ok {
		return t, nil
	}

	// Epoch (pure digits, 10 or 13)
	if opts.AcceptEpoch {
		if t, ok := internal.ParseEpoch(in); ok {
			return t, nil
		}
	}

	// Bare YYYYMMDD (8 digits) → reuse internal.ParseString for validation
	if internal.AllDigits(in) && len(in) == 8 {
		if t, err := internal.ParseString(in); err == nil {
			return t, nil
		}
	}

	// Layouts (deterministic order)
	loc := opts.Location
	if loc == nil {
		loc = time.UTC
	}
	for _, layout := range opts.AcceptLayouts {
		switch layout {
		// date-only → normalize to UTC midnight
		case "2006-01-02", "2006/01/02", "02 Jan 2006":
			t, err := time.Parse(layout, in)
			if err != nil {
				continue
			}
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil

		// zoneless → parse in provided location, then convert to UTC
		case "2006-01-02 15:04:05", "2006/01/02 15:04:05":
			t, err := time.ParseInLocation(layout, in, loc)
			if err != nil {
				continue
			}
			return t.UTC(), nil

		// zoned → preserve parsed zone
		case time.RFC1123:
			if t, ok := internal.ParseZoned(layout, in); ok {
				return t, nil
			}
		}
	}

	return time.Time{}, errors.New("date.CleanAndParse: unrecognized format")
}

// CleanAndParseAsString wraps CleanAndParse and formats the result.
// Returns "" if parsing fails. If layout is empty, it defaults to "2006-01-02".
func CleanAndParseAsString(raw, layout string) string {
	if strings.TrimSpace(layout) == "" {
		layout = "2006-01-02"
	}
	t, err := CleanAndParse(raw, nil)
	if err != nil {
		return ""
	}
	return t.Format(layout)
}
