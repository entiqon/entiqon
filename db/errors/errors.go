package errors

import "errors"

var (
	// UnsupportedTypeError is returned by helpers.ValidateType when the input
	// is not a raw string, but instead a token or another unsupported type.
	//
	// For example:
	//   table.New(table.New("users"))
	// will return this error, suggesting the caller use Clone() instead
	// of nesting tokens directly.
	UnsupportedTypeError = errors.New("unsupported type")

	// InvalidIdentifierError is a sentinel error returned when an identifier
	// (table name, field name, alias, etc.) fails validation.
	//
	// This error is wrapped with context-specific messages, allowing callers
	// to distinguish identifier validation failures using errors.Is:
	//
	//   if errors.Is(err, errors.InvalidIdentifierError) {
	//       // Handle invalid identifier error
	//   }
	//
	// Example failures:
	//   - table.New("???")   → invalid table identifier
	//   - field.New("1abc") → invalid field identifier
	InvalidIdentifierError = errors.New("invalid identifier")
)
