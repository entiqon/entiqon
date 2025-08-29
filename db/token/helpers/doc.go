// Package helpers provides utility functions for validating and
// classifying SQL identifiers, aliases, and expression fragments.
//
// # Purpose
//
// These helpers centralize low-level validation logic that is reused
// across multiple tokens (Field, Table, Join, etc.). Examples include
// checking whether a string is a valid identifier, whether an alias
// is acceptable, or whether an expression has a trailing alias.
//
// # Current Rules
//
// The current implementation applies simplified, dialect-agnostic rules:
//
//   - Identifiers must start with a letter or underscore and may
//     contain letters, digits, and underscores.
//   - Aliases must be valid identifiers and not match reserved keywords.
//   - Trailing alias detection is done heuristically by checking the
//     last token when "AS" is not present.
//
// These rules are intentionally strict and conservative to prevent
// invalid tokens from being accepted silently.
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
