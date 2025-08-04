// file: common/errors/causable.go

package errors_test

import (
	"testing"

	"github.com/entiqon/entiqon/common/errors"
)

type CausableErrorSuite struct {
	t *testing.T
}

func (s *CausableErrorSuite) TestImplementsError(_ *testing.T) {
	var _ error = errors.NewCausableError("cause", "reason")
}

func (s *CausableErrorSuite) TestNewCausableError(t *testing.T) {
	cause := "Database"
	reason := "Connection failed"

	err := errors.NewCausableError(cause, reason)

	if err.Cause() != cause {
		t.Errorf("Cause() = %q; want %q", err.Cause(), cause)
	}

	if err.Reason() != reason {
		t.Errorf("Reason() = %q; want %q", err.Reason(), reason)
	}

	if err.Error() != reason {
		t.Errorf("Error() = %q; want %q", err.Error(), reason)
	}
}

func (s *CausableErrorSuite) TestEmptyFields(t *testing.T) {
	err := errors.NewCausableError("", "")

	if err.Cause() != "" {
		t.Errorf("Cause() = %q; want empty string", err.Cause())
	}
	if err.Reason() != "" {
		t.Errorf("Reason() = %q; want empty string", err.Reason())
	}
	if err.Error() != "" {
		t.Errorf("Error() = %q; want empty string", err.Error())
	}
}

func TestCausableErrorSuite(t *testing.T) {
	suite := &CausableErrorSuite{}

	t.Run("ImplementsError", suite.TestImplementsError)
	t.Run("NewCausableError", suite.TestNewCausableError)
	t.Run("EmptyFields", suite.TestEmptyFields)
}
