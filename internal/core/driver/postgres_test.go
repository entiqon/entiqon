package driver_test

import (
	"testing"

	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/stretchr/testify/suite"
)

type PostgresDialectTestSuite struct {
	suite.Suite
	d *driver.PostgresDialect
}

func (s *PostgresDialectTestSuite) SetupTest() {
	s.d = driver.NewPostgresDialect()
}

func (s *PostgresDialectTestSuite) TestPostgresDialect_NewInstance() {
	s.Equal("postgres", s.d.GetName())

	s.Run("RenderFrom", func() {
		raw := s.d.RenderFrom("user", "")
		s.Equal("\"user\"", raw)
	})
	s.Run("RenderFromWithAlias", func() {
		raw := s.d.RenderFrom("user", "u")
		s.Equal("\"user\" u", raw)
	})
}

func (s *PostgresDialectTestSuite) TestPostgresDialect_SupportsUpsert() {
	s.True(s.d.SupportsUpsert())
}

func (s *PostgresDialectTestSuite) TestPostgresDialect_SupportsReturning() {
	s.True(s.d.SupportsReturning())
}

func (s *PostgresDialectTestSuite) TestPostgresDialect_QuoteIdentifier() {
	s.Equal(`"username"`, s.d.QuoteIdentifier("username"))
}

func (s *PostgresDialectTestSuite) TestPostgresDialect_Placeholder() {
	s.Equal("$1", s.d.Placeholder(1))
	s.Equal("$99", s.d.Placeholder(99))
}

func TestPostgresDialectTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresDialectTestSuite))
}
