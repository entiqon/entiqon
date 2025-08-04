// File: db/internal/core/error/error_collector.go
// Description: Manages accumulation and grouping of StageErrors.
// Since: v1.4.0

package errors

import (
	stdErrors "errors"
	"fmt"
)

// StageErrorCollector accumulates StageErrors and formats them for builder validation.
type StageErrorCollector struct {
	errors []StageError
}

// AddStageError adds a non-nil error associated with a specific stage.
func (c *StageErrorCollector) AddStageError(stage StageToken, err error) {
	if err != nil {
		c.errors = append(c.errors, StageError{Stage: stage, Error: err})
	}
}

// HasErrors reports whether any stage errors were recorded.
func (c *StageErrorCollector) HasErrors() bool {
	return len(c.errors) > 0
}

// GetErrors returns all recorded stage errors.
func (c *StageErrorCollector) GetErrors() []StageError {
	return c.errors
}

// ErrorsByStage groups all stage errors by their associated stage.
func (c *StageErrorCollector) ErrorsByStage() map[string][]error {
	grouped := make(map[string][]error)
	for _, e := range c.errors {
		grouped[e.Stage.String()] = append(grouped[e.Stage.String()], e.Error)
	}
	return grouped
}

// CombineErrors returns a grouped and indented summary of all stage errors.
func (c *StageErrorCollector) CombineErrors() error {
	if len(c.errors) == 0 {
		return nil
	}
	grouped := c.ErrorsByStage()
	msg := "builder validation failed:"
	for stage, errs := range grouped {
		if len(errs) == 1 {
			msg += fmt.Sprintf("\n  - [%s] %s", stage, errs[0])
		} else {
			msg += fmt.Sprintf("\n  - [%s]", stage)
			for _, err := range errs {
				msg += fmt.Sprintf("\n     - %s", err)
			}
		}
	}
	return stdErrors.New(msg)
}

// String returns the same grouped view as CombineErrors but implements fmt.Stringer.
func (c *StageErrorCollector) String() string {
	err := c.CombineErrors()
	if err == nil {
		return ""
	}
	return err.Error()
}
