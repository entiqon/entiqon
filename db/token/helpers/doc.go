// Package helpers provides utility functions for validating and
// classifying SQL identifiers and aliases.
//
// # Purpose
//
// These helpers centralize low-level validation logic that is reused
// across multiple tokens (Field, Table, etc.). Examples include
// checking whether a string is a valid identifier, whether an alias
// is acceptable, whether a trailing alias is present, or generating
// deterministic aliases when none are provided.
//
// # Current Rules
//
// The current implementation applies simplified, dialect-agnostic rules:
//
//   - Identifiers must start with a letter or underscore and may
//     contain letters, digits, and underscores.
//   - Non-ASCII identifiers (e.g. café, mañana, niño) are rejected.
//   - Aliases must be valid identifiers and must not be reserved
//     keywords (case-insensitive).
//   - Trailing aliases (e.g. "(price * qty) total") are valid if the
//     last token is a valid alias and not part of the expression.
//   - Explicit AS aliases are handled by the resolver, not helpers.
//   - Deterministic aliases can be generated with GenerateAlias(),
//     which combines a two-letter code with a SHA-1 hash of the
//     expression string.
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
//   - GenerateAlias(prefix, expr) string → produces safe, deterministic aliases.
//
// This ensures consistent usage across identifiers, aliases, trailing
// alias detection, and generated aliases.
//
// # Future Dialect-Specific Rules
//
// In the future, dialect packages (e.g. Postgres, MySQL) will provide
// their own validation rules to reflect the full grammar of each SQL
// dialect. At that point, helpers may delegate to dialect-specific
// implementations while preserving the same external contract.
//
// # Auto-Alias Rules
//
// If an expression is not a plain identifier and has no alias, it may
// receive a generated alias using GenerateAlias() together with the
// alias code provided by ExpressionKind.Alias() (e.g. "fn_a1b2c3").
// This ensures all non-identifier expressions can be referenced
// reliably downstream. Aliases that are explicitly invalid (bad syntax,
// reserved keyword) will still be rejected.
//
// # Reserved Keywords
//
// The ReservedKeywords function returns the dialect-agnostic set of
// keywords currently disallowed as aliases. Dialect packages may extend
// or override this list.
//
// # Testing
//
// Each helper is tested independently in its own *_test.go file with
// exhaustive cases to ensure correctness and 100% coverage.
package helpers
