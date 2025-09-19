package errors_test

import (
	stdErrors "errors"
	"fmt"
	"testing"

	"github.com/entiqon/db/errors"
)

func ExampleInvalidIdentifierError() {
	// Simulate a validation failure for a table identifier.
	err := fmt.Errorf("%w: name contains invalid characters", errors.InvalidIdentifierError)

	if stdErrors.Is(err, errors.InvalidIdentifierError) {
		fmt.Println("caught invalid identifier error")
	}

	// Output:
	// caught invalid identifier error
}

func ExampleUnsupportedTypeError() {
	// Simulate passing an unsupported type into a constructor.
	err := fmt.Errorf("%w: type %T not allowed", errors.UnsupportedTypeError, 123)

	if stdErrors.Is(err, errors.UnsupportedTypeError) {
		fmt.Println("caught unsupported type error")
	}

	// Output:
	// caught unsupported type error
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name     string
		baseErr  error
		wrapped  error
		expected error
	}{
		{
			name:     "InvalidIdentifier",
			baseErr:  errors.InvalidIdentifierError,
			wrapped:  fmt.Errorf("%w: bad format", errors.InvalidIdentifierError),
			expected: errors.InvalidIdentifierError,
		},
		{
			name:     "UnsupportedType",
			baseErr:  errors.UnsupportedTypeError,
			wrapped:  fmt.Errorf("%w: got int", errors.UnsupportedTypeError),
			expected: errors.UnsupportedTypeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !stdErrors.Is(tt.wrapped, tt.expected) {
				t.Errorf("expected errors.Is(%v, %v) to be true", tt.wrapped, tt.expected)
			}
		})
	}
}
