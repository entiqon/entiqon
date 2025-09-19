// File: db/internal/errors/error/stage_error_test.go
// Since: v1.5.0

package errors_test

import (
	stdErrors "errors"
	"strings"
	"testing"

	"github.com/entiqon/db/internal/core/errors"
)

func TestStageErrorCollector_AddAndHasErrors(t *testing.T) {
	c := &errors.StageErrorCollector{}
	c.AddStageError("FROM", stdErrors.New("table missing"))
	c.AddStageError("SELECT", stdErrors.New("no fields"))

	if !c.HasErrors() {
		t.Errorf("expected HasErrors=true")
	}
	if got := len(c.GetErrors()); got != 2 {
		t.Errorf("expected 2 errors, got %d", got)
	}
}

func TestStageErrorCollector_ErrorsByStage(t *testing.T) {
	c := &errors.StageErrorCollector{}
	c.AddStageError("SELECT", stdErrors.New("missing column"))
	c.AddStageError("SELECT", stdErrors.New("alias invalid"))
	c.AddStageError("WHERE", stdErrors.New("bad condition"))

	grouped := c.ErrorsByStage()
	if got := len(grouped["SELECT"]); got != 2 {
		t.Errorf("expected 2 SELECT errors, got %d", got)
	}
	if got := len(grouped["WHERE"]); got != 1 {
		t.Errorf("expected 1 WHERE error, got %d", got)
	}
}

func TestStageErrorCollector_CombineErrorsFormat(t *testing.T) {
	c := &errors.StageErrorCollector{}
	c.AddStageError("FROM", stdErrors.New("table empty"))
	c.AddStageError("SELECT", stdErrors.New("missing fields"))
	c.AddStageError("SELECT", stdErrors.New("bad alias"))

	output := c.CombineErrors().Error()
	if !strings.Contains(output, "[FROM] table empty") {
		t.Errorf("expected output to contain %q, got %q", "[FROM] table empty", output)
	}
	if !strings.Contains(output, "[SELECT]") {
		t.Errorf("expected output to contain [SELECT], got %q", output)
	}
	if !strings.Contains(output, "missing fields") {
		t.Errorf("expected output to contain 'missing fields', got %q", output)
	}
	if !strings.Contains(output, "bad alias") {
		t.Errorf("expected output to contain 'bad alias', got %q", output)
	}
}

func TestStageErrorCollector_StringMethod(t *testing.T) {
	c := &errors.StageErrorCollector{}
	c.AddStageError("ORDER", stdErrors.New("invalid direction"))

	if got := c.String(); got != c.CombineErrors().Error() {
		t.Errorf("expected String() and CombineErrors().Error() to match, got %q vs %q",
			got, c.CombineErrors().Error())
	}
}
