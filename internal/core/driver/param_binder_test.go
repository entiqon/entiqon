package driver_test

import (
	"testing"

	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/stretchr/testify/suite"
)

type ParamBinderTestSuite struct {
	suite.Suite
	binder *driver.ParamBinder
}

func (s *ParamBinderTestSuite) SetupTest() {
	s.binder = driver.NewParamBinder(driver.NewPostgresDialect())
}

func (s *ParamBinderTestSuite) TestBind() {
	placeholder := s.binder.Bind("admin")
	s.Equal("$1", placeholder)
	s.Equal([]any{"admin"}, s.binder.Args())
}

func (s *ParamBinderTestSuite) TestBindMany() {
	placeholders := s.binder.BindMany(42, true, "active")
	s.Equal([]string{"$1", "$2", "$3"}, placeholders)
	s.Equal([]any{42, true, "active"}, s.binder.Args())
}

func (s *ParamBinderTestSuite) TestArgsReturnsBoundValues() {
	s.binder.Bind("first")
	s.binder.Bind("second")
	s.Equal([]any{"first", "second"}, s.binder.Args())
}

func TestParamBinderTestSuite(t *testing.T) {
	suite.Run(t, new(ParamBinderTestSuite))
}
