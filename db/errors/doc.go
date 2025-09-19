// Package errors defines sentinel errors used across Entiqon’s SQL
// builder tokens (tables, fields, joins, conditions).
//
// These errors provide consistent classification for common failure modes
// such as unsupported constructor types or invalid identifier names.
//
// # Overview
//
// The package exposes sentinel error values that can be used with
// [errors.Is] for robust error handling:
//
//   - [UnsupportedTypeError] is returned by constructors (e.g. table.New)
//     when an input type is not allowed. For example:
//
//     table.New(table.New("users"))
//     // → error: unsupported type; if you want to create a copy, use Clone() instead
//
//   - [InvalidIdentifierError] is returned when an identifier fails
//     validation (e.g. invalid characters or format). For example:
//
//     table.New("???")
//     // → error: invalid table identifier
//
// # Usage
//
// Callers should prefer [errors.Is] to detect sentinel values:
//
//	if errors.Is(err, errors.InvalidIdentifierError) {
//	    log.Printf("invalid identifier: %v", err)
//	}
//
// This allows context-specific constructors (table, field, etc.) to wrap
// the base sentinel with their own diagnostic messages while preserving
// consistent error typing.
package errors
