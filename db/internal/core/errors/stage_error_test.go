// File: db/internal/core/error/stage_error_test.go
// Since: v1.5.0

package errors_test

import (
	"errors"
	"testing"

	core "github.com/entiqon/entiqon/db/internal/core/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StageErrorSuite struct {
	suite.Suite
	collector *core.StageErrorCollector
}

func (s *StageErrorSuite) SetupTest() {
	s.collector = &core.StageErrorCollector{}
}

func (s *StageErrorSuite) TestAddAndHasErrors() {
	s.collector.AddStageError("FROM", errors.New("table missing"))
	s.collector.AddStageError("SELECT", errors.New("no fields"))
	assert.True(s.T(), s.collector.HasErrors())
	assert.Len(s.T(), s.collector.GetErrors(), 2)
}

func (s *StageErrorSuite) TestErrorsByStage() {
	s.collector.AddStageError("SELECT", errors.New("missing column"))
	s.collector.AddStageError("SELECT", errors.New("alias invalid"))
	s.collector.AddStageError("WHERE", errors.New("bad condition"))

	grouped := s.collector.ErrorsByStage()
	assert.Len(s.T(), grouped["SELECT"], 2)
	assert.Len(s.T(), grouped["WHERE"], 1)
}

func (s *StageErrorSuite) TestCombineErrorsFormat() {
	s.collector.AddStageError("FROM", errors.New("table empty"))
	s.collector.AddStageError("SELECT", errors.New("missing fields"))
	s.collector.AddStageError("SELECT", errors.New("bad alias"))

	output := s.collector.CombineErrors()
	assert.ErrorContains(s.T(), output, "[FROM] table empty")
	assert.ErrorContains(s.T(), output, "[SELECT]")
	assert.ErrorContains(s.T(), output, "missing fields")
	assert.ErrorContains(s.T(), output, "bad alias")
}

func (s *StageErrorSuite) TestStringMethod() {
	s.collector.AddStageError("ORDER", errors.New("invalid direction"))
	assert.Equal(s.T(), s.collector.String(), s.collector.CombineErrors().Error())
}

func TestStageErrorSuite(t *testing.T) {
	suite.Run(t, new(StageErrorSuite))
}
