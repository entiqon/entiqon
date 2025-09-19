// file: common/errors/process_stage_test.go

package errors_test

import (
	stdErrors "errors"
	"testing"

	"github.com/entiqon/common"
	"github.com/entiqon/common/errors"
)

func TestProcessStageError(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		stage := common.ProcessStage("Init")
		cause := "Loader"
		reason := "Failed to load resource"
		wrapped := stdErrors.New("file not found")

		err := errors.NewProcessStageError(stage, cause, reason, wrapped)

		var pse *errors.ProcessStageError
		ok := stdErrors.As(err, &pse)
		if !ok {
			t.Fatalf("expected *ProcessStageError, got %T", err)
		}

		if pse.Stage() != stage {
			t.Errorf("Stage() = %v; want %v", pse.Stage(), stage)
		}
		if pse.Cause() != cause {
			t.Errorf("Cause() = %q; want %q", pse.Cause(), cause)
		}
		if pse.Reason() != reason {
			t.Errorf("Reason() = %q; want %q", pse.Reason(), reason)
		}
		if !stdErrors.Is(wrapped, pse.Unwrap()) {
			t.Errorf("Unwrap() = %v; want %v", pse.Unwrap(), wrapped)
		}

		wantMsg := "[Loader] at stage Init: Failed to load resource: file not found"
		if pse.Error() != wantMsg {
			t.Errorf("Error() = %q; want %q", pse.Error(), wantMsg)
		}
	})

	t.Run("NoWrapped", func(t *testing.T) {
		stage := common.ProcessStage("Parse")
		cause := "Parser"
		reason := "Syntax error"

		err := errors.NewProcessStageError(stage, cause, reason, nil)

		var pse *errors.ProcessStageError
		stdErrors.As(err, &pse)

		wantMsg := "[Parser] at stage Parse: Syntax error"
		if pse.Error() != wantMsg {
			t.Errorf("Error() = %q; want %q", pse.Error(), wantMsg)
		}
		if pse.Unwrap() != nil {
			t.Errorf("Unwrap() = %v; want nil", pse.Unwrap())
		}
	})
}
