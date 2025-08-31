package errors

import "errors"

var (
	// UnsupportedTypeError is returned by ValidateType when the input
	// is a token or another unsupported value instead of a raw string.
	UnsupportedTypeError = errors.New("unsupported type")
)
