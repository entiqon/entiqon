// file: common/errors/process_error.go

package errors

import (
	"fmt"

	"github.com/entiqon/common"
)

// ProcessStageError represents an error occurring at a specific process stage,
// with a cause, reason, and optional wrapped error.
type ProcessStageError struct {
	stage  common.ProcessStage
	cause  string
	reason string
	err    error
}

// Ensure ProcessStageError implements CausableError.
var _ CausableError = (*ProcessStageError)(nil)

// Error returns the error message including cause, stage, reason, and wrapped error if present.
func (e *ProcessStageError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("[%s] at stage %s: %s: %v", e.cause, e.stage, e.reason, e.err)
	}
	return fmt.Sprintf("[%s] at stage %s: %s", e.cause, e.stage, e.reason)
}

// Unwrap returns the underlying wrapped error.
func (e *ProcessStageError) Unwrap() error {
	return e.err
}

// Cause returns the cause or category of the error.
func (e *ProcessStageError) Cause() string {
	return e.cause
}

// Reason returns the detailed explanation of the error.
func (e *ProcessStageError) Reason() string {
	return e.reason
}

// Stage returns the process stage where the error occurred.
func (e *ProcessStageError) Stage() common.ProcessStage {
	return e.stage
}

// NewProcessStageError creates a new ProcessStageError instance with the given parameters.
func NewProcessStageError(stage common.ProcessStage, cause, reason string, err error) error {
	return &ProcessStageError{
		stage:  stage,
		cause:  cause,
		reason: reason,
		err:    err,
	}
}
