// file: common/errors/causable.go

package errors

// CausableError is an error interface that extends the built-in error
// by providing additional context through Cause and Reason methods.
//
// Cause returns a short identifier or category indicating the origin
// or source of the error, useful for programmatic error handling.
//
// Reason returns a human-readable explanation describing the error,
// suitable for logging or user messages.
type CausableError interface {
	error

	// Cause returns the cause or category of the error.
	Cause() string

	// Reason returns the detailed explanation of the error.
	Reason() string
}

// causableError is a concrete implementation of the CausableError interface.
type causableError struct {
	// cause is a short string representing the error's origin or category.
	cause string

	// reason is a human-readable explanation describing the error.
	reason string
}

// Error implements the error interface, returning the error message.
func (e *causableError) Error() string {
	return e.reason
}

// Cause returns the cause or category of the error.
func (e *causableError) Cause() string {
	return e.cause
}

// Reason returns the detailed explanation of the error.
func (e *causableError) Reason() string {
	return e.reason
}

// NewCausableError creates a new CausableError with the given cause and reason.
//
// The cause should be a short string indicating the source or category of the error.
// The reason should provide a more descriptive message explaining the error.
func NewCausableError(cause, reason string) CausableError {
	return &causableError{
		cause:  cause,
		reason: reason,
	}
}
