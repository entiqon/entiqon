// File: db/internal/build/token/token_helper.go

package token

import "strings"

// ParseAlias extracts a base identifier and optional alias from a string.
//
// It supports multiple input formats:
//   - Standard SQL aliasing:     "id AS user_id"
//   - Loose shorthand formats:   "id user_id" or "id, user_id"
//
// Examples:
//
//	ParseAlias("id")               → ("id", "")
//	ParseAlias("id AS user_id")    → ("id", "user_id")
//	ParseAlias("id user_id")       → ("id", "user_id")
//	ParseAlias("id, user_id")      → ("id", "user_id")
//
// Whitespace is trimmed automatically.
//
// This function is used internally by NewColumn, NewTable, etc.
// It is part of the token parsing utility layer.
func ParseAlias(input string) (string, string) {
	trimmed := strings.TrimSpace(input)

	// Handle standard SQL aliasing: "name AS alias"
	if parts := strings.SplitN(trimmed, " AS ", 2); len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	// Handle loose formats: "name alias" or "name, alias"
	if parts := strings.FieldsFunc(trimmed, func(r rune) bool {
		return r == ',' || r == ' '
	}); len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	// No alias found
	return trimmed, ""
}

// ParseTableColumn splits a SQL identifier into its optional table prefix and column name.
//
// It supports formats like:
//
//	"column"             → table: "",     column: "column"
//	"table.column"       → table: "table", column: "column"
//	"  table . column  " → table: "table", column: "column" (whitespace is trimmed)
//
// If no dot is present, the entire input is treated as the column name, with no table.
//
// This function is tolerant of formatting inconsistencies and is used to support expressions
// like "users.id", "orders.total", or just "id".
//
// It is commonly used by NewColumn, NewTable, and other SQL tokens to normalize identifier parsing.
func ParseTableColumn(input string) (table string, column string) {
	parts := strings.SplitN(strings.TrimSpace(input), ".", 2)

	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	return "", strings.TrimSpace(parts[0])
}
