// File: db/internal/core/builder/bind/param_binder_test.go

package bind_test

import (
	"testing"

	driver2 "github.com/entiqon/entiqon/db/driver"
	"github.com/entiqon/entiqon/db/internal/core/builder/bind"
	"github.com/stretchr/testify/suite"
)

type ParamBinderTestSuite struct {
	suite.Suite
	binder *bind.ParamBinder
}

func (s *ParamBinderTestSuite) SetupTest() {
	s.binder = bind.NewParamBinder(driver2.NewGenericDialect())
}

func (s *ParamBinderTestSuite) TestBind() {
	placeholder := s.binder.Bind("admin")
	s.Equal("?", placeholder)
	s.Equal([]any{"admin"}, s.binder.Args())

	s.Run("WithPostgres", func() {
		binder := bind.NewParamBinder(driver2.NewPostgresDialect())
		placeholder = binder.Bind("alpha")
		if placeholder != "$1" {
			s.T().Errorf("expected $1, got %s", placeholder)
		}

		args := binder.Args()
		if len(args) != 1 || args[0] != "alpha" {
			s.T().Errorf("expected args [alpha], got %#v", args)
		}
	})
}

func (s *ParamBinderTestSuite) TestBindMany() {
	placeholders := s.binder.BindMany(42, true, "active")
	s.Equal([]string{"?", "?", "?"}, placeholders)
	s.Equal([]any{42, true, "active"}, s.binder.Args())

	s.Run("WithPostgres", func() {
		binder := bind.NewParamBinder(driver2.NewPostgresDialect())
		placeholders := binder.BindMany(42, true, "active")
		s.Equal([]string{"$1", "$2", "$3"}, placeholders)
		s.Equal([]any{42, true, "active"}, binder.Args())
	})
}

func (s *ParamBinderTestSuite) TestArgsReturnsBoundValues() {
	s.binder.Bind("first")
	s.binder.Bind("second")
	s.Equal([]any{"first", "second"}, s.binder.Args())
}

func (s *ParamBinderTestSuite) TestParamBinderWithPosition() {
	binder := bind.NewParamBinderWithPosition(driver2.NewGenericDialect(), 4)
	placeholder := binder.Bind("next")
	if placeholder != "?" {
		s.T().Errorf("expected ?, got %s", placeholder)
	}

	args := binder.Args()
	if len(args) != 1 || args[0] != "next" {
		s.T().Errorf("unexpected args: %#v", args)
	}

	s.Run("WithPostgres", func() {
		binder = bind.NewParamBinderWithPosition(driver2.NewPostgresDialect(), 4)

		placeholder := binder.Bind("next")
		if placeholder != "$4" {
			s.T().Errorf("expected $4, got %s", placeholder)
		}

		args := binder.Args()
		if len(args) != 1 || args[0] != "next" {
			s.T().Errorf("unexpected args: %#v", args)
		}
	})
}

func TestParamBinderTestSuite(t *testing.T) {
	suite.Run(t, new(ParamBinderTestSuite))
}
