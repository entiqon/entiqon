// Package helpers provides utility functions for validating and
// classifying SQL identifiers.
//
// # Purpose
//
// These helpers centralize low-level validation logic that is reused
// across multiple tokens (Field, Table, etc.). Examples include
// checking whether a string is a valid identifier.
//
// # Current Rules
//
// The current implementation applies simplified, dialect-agnostic rules:
//
//   - Identifiers must start with a letter or underscore and may
//     contain letters, digits, and underscores.
//   - Non-ASCII identifiers (e.g. café, mañana, niño) are rejected.
//
// These rules are intentionally strict and conservative to prevent
// invalid tokens from being accepted silently.
//
// # Consistency
//
// All helpers follow the same validation pattern:
//
//   - ValidateXxx(s string) error → returns a detailed error if invalid.
//   - IsValidXxx(s string) bool   → returns true/false as a convenience wrapper.
//
// This ensures consistent usage across all helpers.
//
// # Future Dialect-Specific Rules
//
// In the future, dialect packages (e.g. Postgres, MySQL) will provide
// their own validation rules to reflect the full grammar of each SQL
// dialect. At that point, helpers may delegate to dialect-specific
// implementations while preserving the same external contract.
//
// # Testing
//
// Each helper is tested independently in its own *_test.go file with
// exhaustive cases to ensure correctness and 100% coverage.
package helpers
