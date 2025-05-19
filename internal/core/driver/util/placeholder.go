// Package util provides reusable helper functions for SQL query construction.
// This module focuses on dialect-safe utilities used by query builders.
//
// Since: v1.4.0
package util

import (
	"strings"

	"github.com/ialopezg/entiqon/internal/core/driver"
)

// GeneratePlaceholders returns the dialect-specific SQL placeholders and arguments
// for a 2D matrix of row values. This is used in multi-row INSERT and UPSERT statements.
//
// It ensures:
//   - Positional placeholders are emitted in global order (e.g., $1, $2, $3, ...)
//   - Arguments are flattened across all rows
//   - Fallback to "?" if dialect is nil
//
// Example output (PostgreSQL, 2 rows × 2 columns):
//
//	placeholders = ["($1, $2)", "($3, $4)"]
//	args         = [val1, val2, val3, val4]
//
// Since: v1.4.0
func GeneratePlaceholders(values [][]any, dialect driver.Dialect) ([]string, []any) {
	argIndex := 1
	placeholders := make([]string, 0, len(values))
	args := make([]any, 0, len(values)*4)

	for _, row := range values {
		rowPlaceholders := make([]string, len(row))
		for i := range row {
			if dialect != nil {
				rowPlaceholders[i] = dialect.Placeholder(argIndex)
			} else {
				rowPlaceholders[i] = "?"
			}
			argIndex++
		}
		placeholders = append(placeholders, "("+strings.Join(rowPlaceholders, ", ")+")")
		args = append(args, row...)
	}

	return placeholders, args
}
