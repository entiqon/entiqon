package errors_test

import (
	stdErrors "errors"
	"fmt"

	"github.com/entiqon/db/errors"
)

func Example() {
	// Simulate an identifier validation failure
	err := fmt.Errorf("%w: name contains invalid characters", errors.InvalidIdentifierError)

	if stdErrors.Is(err, errors.InvalidIdentifierError) {
		fmt.Println("invalid identifier caught")
	}

	// Simulate passing unsupported type into constructor
	err = fmt.Errorf("%w: got int", errors.UnsupportedTypeError)

	if stdErrors.Is(err, errors.UnsupportedTypeError) {
		fmt.Println("unsupported type caught")
	}

	// Output:
	// invalid identifier caught
	// unsupported type caught
}
