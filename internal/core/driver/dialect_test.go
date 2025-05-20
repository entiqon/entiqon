package driver_test

import (
	"testing"

	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/stretchr/testify/suite"
)

type DialectTestSuite struct {
	suite.Suite
	generic  driver.Dialect
	postgres driver.Dialect
}

func (s *DialectTestSuite) SetupSuite() {
	s.generic = driver.NewGenericDialect()
	s.postgres = driver.NewPostgresDialect()
}

func (s *DialectTestSuite) TestGenericDialect() {
	s.Equal("generic", s.generic.Name())
	s.Equal("?", s.generic.Placeholder(1))
	s.Equal("?", s.generic.Placeholder(99))
	s.Equal("users", s.generic.QuoteIdentifier("users"))
}

func (s *DialectTestSuite) TestPostgresDialect() {
	s.Equal("postgres", s.postgres.Name())
	s.Equal("$1", s.postgres.Placeholder(1))
	s.Equal("$5", s.postgres.Placeholder(5))
	s.Equal(`"users"`, s.postgres.QuoteIdentifier("users"))
	s.True(s.postgres.SupportsUpsert())
	s.True(s.postgres.SupportsReturning())
}

func (s *DialectTestSuite) TestResolveDialect() {
	s.Equal("postgres", driver.ResolveDialect("postgres").Name())
	s.Equal("generic", driver.ResolveDialect("unknown").Name())
}

func (s *DialectTestSuite) TestBaseDialectDirectMethods() {
	base := &driver.BaseDialect{DialectName: "test"}

	s.Equal("test", base.Name())
	s.Equal("?", base.Placeholder(999))
	s.Equal("field", base.QuoteIdentifier("field"))
	s.False(base.SupportsUpsert())
	s.False(base.SupportsReturning())
	s.Equal("LIMIT 10 OFFSET 20", base.BuildLimitOffset(10, 20))
	s.Equal("LIMIT 5", base.BuildLimitOffset(5, -1))
	s.Equal("OFFSET 20", base.BuildLimitOffset(-1, 20))
	s.Equal("", base.BuildLimitOffset(-1, -1))

	// QuoteLiteral coverage
	s.Equal("'value'", base.QuoteLiteral("value"))
	s.Equal("42", base.QuoteLiteral(42))
	s.Equal("true", base.QuoteLiteral(true))
	s.Equal("'[1 2 3]'", base.QuoteLiteral([]int{1, 2, 3}))
}

func TestDialectTestSuite(t *testing.T) {
	suite.Run(t, new(DialectTestSuite))
}
