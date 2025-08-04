// File: db/internal/core/token/field_token.go

package token_test

import (
	"testing"

	token2 "github.com/entiqon/db/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type FieldTokenTestSuite struct {
	suite.Suite
}

func (s *FieldTokenTestSuite) TestField_Basic() {
	col := token2.Field("id")
	s.Equal("id", col.Name)
	s.False(col.IsRaw)
	s.Empty(col.Alias)
}

func (s *FieldTokenTestSuite) TestFieldExpr_WithAlias() {
	expr := token2.FieldExpr("COUNT(*)", "total")
	s.Equal("COUNT(*)", expr.Name)
	s.Equal("total", expr.Alias)
	s.True(expr.IsRaw)
}

func (s *FieldTokenTestSuite) TestField_AsMethod() {
	aliased := token2.Field("created_at").As("created")
	s.Equal("created_at", aliased.Name)
	s.Equal("created", aliased.Alias)
}

func (s *FieldTokenTestSuite) TestIsValid() {
	s.True(token2.Field("status").IsValid())
	s.False(token2.Field("").IsValid())
}

func (s *FieldTokenTestSuite) TestField_WithAliasInline() {
	f := token2.Field("first_name AS name")
	s.Equal("first_name", f.Name)
	s.Equal("name", f.Alias)
	s.True(f.IsValid())
}

func (s *FieldTokenTestSuite) TestField_WithAliasParams() {
	f := token2.Field("first_name", "name")
	s.Equal("first_name", f.Name)
	s.Equal("name", f.Alias)
	s.True(f.IsValid())
}

func (s *FieldTokenTestSuite) TestFieldsFromExpr_CommaSeparated() {
	fields := token2.FieldsFromExpr("id, first_name AS name, email AS contact")
	s.Len(fields, 3)
	s.Equal("id", fields[0].Name)
	s.Equal("first_name", fields[1].Name)
	s.Equal("name", fields[1].Alias)
	s.Equal("email", fields[2].Name)
	s.Equal("contact", fields[2].Alias)
}

func (s *FieldTokenTestSuite) TestField_CommaSeparatedPanics() {
	s.PanicsWithValue(
		"Field: comma-separated values not allowed in a single call. Call Field(...) separately for each.",
		func() {
			_ = token2.Field("id, name")
		},
	)
}

func (s *FieldTokenTestSuite) TestFieldToken_WithValue() {
	f := token2.Field("email")
	fv := f.WithValue("x@entiqon.dev")

	s.Equal("email", fv.Name)
	s.Equal("x@entiqon.dev", fv.Value)
}

func TestFieldTokenTestSuite(t *testing.T) {
	suite.Run(t, new(FieldTokenTestSuite))
}
