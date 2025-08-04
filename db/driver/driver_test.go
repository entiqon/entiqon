// File: db/driver/driver.go

package driver_test

import (
	"testing"

	"github.com/entiqon/entiqon/db/driver"
	"github.com/entiqon/entiqon/db/driver/styling"
	"github.com/stretchr/testify/suite"
)

type DialectTestSuite struct {
	suite.Suite
	base     driver.Dialect
	generic  driver.Dialect
	postgres driver.Dialect
	mssql    driver.Dialect
	mysql    driver.Dialect
}

func (s *DialectTestSuite) SetupSuite() {
	s.base = &driver.BaseDialect{}
	s.generic = driver.NewGenericDialect()
	s.postgres = driver.NewPostgresDialect()
	s.mssql = driver.NewMSSQLDialect()
	s.mysql = driver.NewMySQLDialect()
}

func (s *DialectTestSuite) TestResolveDialect() {
	s.Equal("postgres", driver.ResolveDialect("postgres").GetName())
	s.Equal("generic", driver.ResolveDialect("unknown").GetName())
}

func (s *DialectTestSuite) TestBaseDialectDirectMethods() {
	base := &driver.BaseDialect{}

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

func (s *DialectTestSuite) TestGetName() {
	s.Run("base", func() {
		s.Equal("base", s.base.GetName())
	})
	s.Run("generic", func() {
		s.Equal("generic", s.generic.GetName())
	})
	s.Run("postgres", func() {
		s.Equal("postgres", s.postgres.GetName())
	})
	s.Run("mssql", func() {
		s.Equal("mssql", s.mssql.GetName())
	})
	s.Run("mysql", func() {
		s.Equal("mysql", s.mysql.GetName())
	})
}

func (s *DialectTestSuite) TestQuoteIdentifier() {
	s.Run("base", func() {
		s.Equal("user", s.base.QuoteIdentifier("user"))
	})
	s.Run("generic", func() {
		s.Equal("user", s.generic.QuoteIdentifier("user"))
	})
	s.Run("postgres", func() {
		s.Equal("\"user\"", s.postgres.QuoteIdentifier("user"))
	})
	s.Run("mssql", func() {
		s.Equal("[user]", s.mssql.QuoteIdentifier("user"))
	})
	s.Run("mysql", func() {
		s.Equal("`user`", s.mysql.QuoteIdentifier("user"))
	})
}

func (s *DialectTestSuite) TestQuoteTpeQuoteType() {
	s.Run("generic", func() {
		s.Equal(styling.QuoteNone, s.generic.QuoteType())
	})
	s.Run("postgres", func() {
		s.Equal(styling.QuoteDouble, s.postgres.QuoteType())
	})
	s.Run("mssql", func() {
		s.Equal(styling.QuoteBracket, s.mssql.QuoteType())
	})
	s.Run("mysql", func() {
		s.Equal(styling.QuoteBacktick, s.mysql.QuoteType())
	})
}

func (s *DialectTestSuite) TestPlaceholder() {
	s.Run("base", func() {
		s.Equal("?", s.base.Placeholder(1))
		s.Equal("?", s.base.Placeholder(99))
	})
	s.Run("generic", func() {
		s.Equal("?", s.generic.Placeholder(1))
		s.Equal("?", s.generic.Placeholder(99))
	})
	s.Run("postgres", func() {
		s.Equal("$1", s.postgres.Placeholder(1))
		s.Equal("$5", s.postgres.Placeholder(5))
	})
	s.Run("mssql", func() {
		s.Equal("?", s.mssql.Placeholder(1))
		s.Equal("?", s.mssql.Placeholder(99))
	})
	s.Run("mysql", func() {
		s.Equal("?", s.mysql.Placeholder(1))
		s.Equal("?", s.mysql.Placeholder(99))
	})
}

func (s *DialectTestSuite) TestSupportsRenderFrom() {
	s.Run("base", func() {
		s.Equal("users", s.base.RenderFrom("users", ""))
		s.Equal("users", s.base.RenderFrom("users", "u"))
	})
	s.Run("generic", func() {
		s.Equal("users", s.generic.RenderFrom("users", ""))
		s.Equal("users", s.generic.RenderFrom("users", "u"))
	})
	s.Run("postgres", func() {
		s.Equal("\"users\"", s.postgres.RenderFrom("users", ""))
		s.Equal("\"users\" u", s.postgres.RenderFrom("users", "u"))
	})
	s.Run("mssql", func() {
		s.Equal("[users]", s.mssql.RenderFrom("users", ""))
		s.Equal("[users] u", s.mssql.RenderFrom("users", "u"))
	})
	s.Run("mysql", func() {
		s.Equal("`users`", s.mysql.RenderFrom("users", ""))
		s.Equal("`users` u", s.mysql.RenderFrom("users", "u"))
	})
}

func (s *DialectTestSuite) TestSupportsReturning() {
	s.Run("base", func() {
		s.Equal(false, s.base.SupportsReturning())
	})
	s.Run("generic", func() {
		s.Equal(false, s.generic.SupportsReturning())
	})
	s.Run("postgres", func() {
		s.Equal(true, s.postgres.SupportsReturning())
	})
	s.Run("mssql", func() {
		s.Equal(false, s.mssql.SupportsReturning())
	})
	s.Run("mysql", func() {
		s.Equal(false, s.mysql.SupportsReturning())
	})
}

func (s *DialectTestSuite) TestSupportsUpsert() {
	s.Run("base", func() {
		s.Equal(false, s.base.SupportsUpsert())
	})
	s.Run("generic", func() {
		s.Equal(false, s.generic.SupportsUpsert())
	})
	s.Run("postgres", func() {
		s.Equal(true, s.postgres.SupportsUpsert())
	})
	s.Run("mssql", func() {
		s.Equal(false, s.mssql.SupportsUpsert())
	})
	s.Run("mysql", func() {
		s.Equal(false, s.mysql.SupportsUpsert())
	})
}

func (s *DialectTestSuite) TestValidate() {
	s.Run("Valid", func() {
		d := driver.BaseDialect{
			Name:             "test",
			QuoteStyle:       styling.QuoteNone,
			PlaceholderStyle: styling.PlaceholderQuestion,
		}
		err := d.Validate()
		s.NoError(err)
		s.Equal("?", d.Placeholder(0))
	})
	s.Run("MissingName", func() {
		d := driver.BaseDialect{Name: ""}
		err := d.Validate()
		s.Error(err)
		s.Contains(err.Error(), "dialect is not configured")
	})
	//s.Run("MissingPlaceholder", func() {
	//	d := driver.BaseDialect{Name: "test", PlaceholderStyle: styling.PlaceholderDollar}
	//	err := d.Validate()
	//	s.Error(err)
	//})
}

func TestDialectTestSuite(t *testing.T) {
	suite.Run(t, new(DialectTestSuite))
}
