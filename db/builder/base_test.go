package builder_test

import (
	"testing"

	"github.com/entiqon/db/builder"
	"github.com/entiqon/db/driver"
	"github.com/stretchr/testify/suite"
)

type BaseBuilderTestSuite struct {
	suite.Suite
	qb *builder.BaseBuilder
}

func (s *BaseBuilderTestSuite) SetupTest() {
	s.qb = &builder.BaseBuilder{}
}

func (s *BaseBuilderTestSuite) TestHasDialect() {
	s.Run("Valid", func() {
		dialect := driver.BaseDialect{}

		s.Equal("base", dialect.GetName())
		s.False(s.qb.HasDialect())
	})
	s.Run("Invalid", func() {
		s.False(s.qb.HasDialect())
	})
}

func TestBaseBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BaseBuilderTestSuite))
}
