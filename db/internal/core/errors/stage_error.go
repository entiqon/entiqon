// File: db/internal/core/error/stage_error.go
// Description: Defines the StageError struct and its error interface.
// Since: v1.5.0

package errors

import "fmt"

// StageError represents a validation failure tagged by a builder stage.
type StageError struct {
	Stage StageToken
	Error error
}

// String returns a string representation of the StageError.
func (e StageError) String() string {
	return fmt.Sprintf("[%s] %s", e.Stage, e.Error.Error())
}
