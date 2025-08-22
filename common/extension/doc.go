// Package extension provides strongly typed parsers and helpers for primitive
// and structured values. It serves as the foundation for consistent data
// normalization across Entiqon.
//
// Each subpackage (boolean, date, decimal, float, number, integer, object, collection)
// provides type-specific parsing utilities with full test coverage.
//
// In addition to subpackages, extension exposes **shortcut functions** that
// simplify parsing with default fallbacks:
//
//   - BooleanOr(value any, def bool) bool
//   - NumberOr(value any, def float64) float64
//   - FloatOr(value any, def float64) float64
//   - DecimalOr(value any, def string) decimal.Decimal
//   - DateOr(value any, def time.Time) time.Time
//
// These shortcuts wrap the respective ParseFrom functions but return a caller-
// supplied default value if parsing fails, avoiding error handling boilerplate.
//
// Example:
//
//	send := extension.BooleanOr(c.QueryParam("send"), false)
//
// This will attempt to parse the "send" query parameter into a boolean. If the
// value is invalid, it falls back to `false`.
package extension
