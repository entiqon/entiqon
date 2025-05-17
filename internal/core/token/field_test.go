package token_test

import (
	"testing"

	"github.com/ialopezg/entiqon/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type FieldTokenTestSuite struct {
	suite.Suite
}

func (s *FieldTokenTestSuite) TestField_Basic() {
	col := token.Field("id")
	s.Equal("id", col.Name)
	s.False(col.IsRaw)
	s.Empty(col.Alias)
}

func (s *FieldTokenTestSuite) TestFieldExpr_WithAlias() {
	expr := token.FieldExpr("COUNT(*)", "total")
	s.Equal("COUNT(*)", expr.Name)
	s.Equal("total", expr.Alias)
	s.True(expr.IsRaw)
}

func (s *FieldTokenTestSuite) TestField_AsMethod() {
	aliased := token.Field("created_at").As("created")
	s.Equal("created_at", aliased.Name)
	s.Equal("created", aliased.Alias)
}

func (s *FieldTokenTestSuite) TestIsValid() {
	s.True(token.Field("status").IsValid())
	s.False(token.Field("").IsValid())
}

func TestFieldTokenTestSuite(t *testing.T) {
	suite.Run(t, new(FieldTokenTestSuite))
}
