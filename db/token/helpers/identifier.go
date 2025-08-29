// Package helpers provides utility functions for validating and
// classifying SQL identifiers, aliases, and expression fragments.
//
// These helpers are intentionally small, pure functions so they can
// be reused by multiple token types (Field, Table, etc.) and tested
// independently without involving higher-level builders.
package helpers

// IsValidIdentifier reports whether the string is a valid SQL identifier.
//
// Rules:
//   - Must not be empty
//   - First character must be a letter (A–Z, a–z) or underscore (_)
//   - Remaining characters may be letters, digits (0–9), or underscores (_)
//
// This helper does not check for reserved keywords; see IsValidAlias
// for alias-specific validation.
func IsValidIdentifier(s string) bool {
	if s == "" {
		return false
	}

	// first char: letter or underscore
	first := s[0]
	if !(first == '_' ||
		(first >= 'A' && first <= 'Z') ||
		(first >= 'a' && first <= 'z')) {
		return false
	}

	// remaining chars: letter, digit, underscore
	for i := 1; i < len(s); i++ {
		ch := s[i]
		if !(ch == '_' ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9')) {
			return false
		}
	}
	return true
}
