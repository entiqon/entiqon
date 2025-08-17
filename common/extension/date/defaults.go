// File: common/extension/date/defaults.go

package date

import (
	"time"

	"github.com/entiqon/entiqon/common/extension/date/internal"
)

// DefaultLayouts returns a fresh copy so callers can’t mutate the source.
func DefaultLayouts() []string {
	src := internal.DefaultLayouts()
	return append([]string(nil), src...)
}

// DefaultCleanParseOptions returns the recommended default options.
// - Allows epoch seconds/milliseconds
// - Accepts DefaultLayouts
// - Normalizes to UTC if no explicit zone is present
func DefaultCleanParseOptions() *CleanParseOptions {
	return &CleanParseOptions{
		RequireYYYYMMDDPrefix: false,
		AcceptLayouts:         DefaultLayouts(),
		AcceptEpoch:           true,
		Location:              time.UTC,
	}
}

// StrictYYYYMMDDOptions returns options that require the string
// to start with an 8-digit YYYYMMDD prefix, ignoring trailing content.
// Useful for fixed-format feeds like "20250507.0000000000".
func StrictYYYYMMDDOptions() *CleanParseOptions {
	return &CleanParseOptions{
		RequireYYYYMMDDPrefix: true,
		Location:              time.UTC,
	}
}
