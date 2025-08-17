// File:

package internal

import "time"

// DefaultLayouts returns a fresh slice with the deterministic parsing layouts.
// Keep this private to avoid cross-package mutation; callers should copy if needed.
func DefaultLayouts() []string {
	return []string{
		"2006-01-02",          // date-only → UTC midnight
		"2006/01/02",          // date-only → UTC midnight
		"2006-01-02 15:04:05", // no zone → Local→UTC
		"2006/01/02 15:04:05", // no zone → Local→UTC
		"02 Jan 2006",         // date-only → UTC midnight
		time.RFC1123,          // has zone → preserve
	}
}
