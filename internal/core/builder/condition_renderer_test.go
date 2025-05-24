package builder_test

import (
	"testing"

	driver2 "github.com/ialopezg/entiqon/driver"
	"github.com/ialopezg/entiqon/internal/core/builder"
	"github.com/ialopezg/entiqon/internal/core/builder/bind"
	"github.com/ialopezg/entiqon/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type ConditionRendererTestSuite struct {
	suite.Suite
}

func (s *ConditionRendererTestSuite) TestAppendCondition() {
	base := token.NewCondition(token.ConditionSimple, "status", "active")
	and := token.NewCondition(token.ConditionAnd, "deleted", false)
	or := token.NewCondition(token.ConditionOr, "archived", false)

	result := builder.AppendCondition([]token.Condition{base}, and)
	s.Len(result, 2)
	s.Equal(token.ConditionAnd, result[1].Type)

	result = builder.AppendCondition(result, or)
	s.Len(result, 3)
	s.Equal(token.ConditionOr, result[2].Type)

	s.Run("Invalid", func() {
		valid := token.NewCondition(token.ConditionSimple, "status", "active")
		invalid := token.Condition{} // no field/operator/values → invalid

		result := builder.AppendCondition([]token.Condition{valid}, invalid)

		s.Len(result, 1)
		s.Equal("status", result[0].Key)
	})
}

func (s *ConditionRendererTestSuite) TestRenderConditions_Generic() {
	s.Run("Valid", func() {
		conditions := []token.Condition{
			token.NewCondition(token.ConditionSimple, "active", true),
		}

		sql, args, err := builder.RenderConditions(driver2.NewGenericDialect(), conditions)
		s.NoError(err)
		s.Equal("active = ?", sql)
		s.Equal([]any{true}, args)
	})

	s.Run("Empty", func() {
		binder := bind.NewParamBinder(driver2.NewGenericDialect())
		sql, args, err := builder.RenderConditionsWithBinder(driver2.NewGenericDialect(), nil, binder)

		s.NoError(err)
		s.Equal("", sql)
		s.Nil(args)
	})

	s.Run("Unsupported", func() {
		c := token.NewCondition(token.ConditionSimple, "status", "active")
		c.Type = token.ConditionType(rune(999)) // simulate unsupported condition type

		binder := bind.NewParamBinder(driver2.NewGenericDialect())
		_, _, err := builder.RenderConditionsWithBinder(driver2.NewGenericDialect(), []token.Condition{c}, binder)

		s.Error(err)
		s.Contains(err.Error(), "unsupported condition type")
	})

	s.Run("Invalid", func() {
		conditions := []token.Condition{
			{}, // invalid: missing Key and Error is nil
		}

		binder := bind.NewParamBinder(driver2.NewGenericDialect())
		_, _, err := builder.RenderConditionsWithBinder(driver2.NewGenericDialect(), conditions, binder)

		s.Error(err)
		s.Contains(err.Error(), "invalid condition")
	})

	s.Run("WithAndCondition", func() {
		conditions := []token.Condition{
			token.NewCondition(token.ConditionSimple, "status", "active"),
			token.NewCondition(token.ConditionAnd, "deleted", false),
		}

		binder := bind.NewParamBinder(driver2.NewGenericDialect())
		sql, args, err := builder.RenderConditionsWithBinder(driver2.NewGenericDialect(), conditions, binder)

		s.NoError(err)
		s.Equal("status = ? AND deleted = ?", sql)
		s.Equal([]any{"active", false}, args)
	})
}

func (s *ConditionRendererTestSuite) TestRenderConditionsWithBinder_Postgres() {
	conditions := []token.Condition{
		token.NewCondition(token.ConditionSimple, "email_verified", true),
		token.NewCondition(token.ConditionOr, "email_verified", false),
	}

	binder := bind.NewParamBinderWithPosition(driver2.NewPostgresDialect(), 4)
	sql, args, err := builder.RenderConditionsWithBinder(driver2.NewPostgresDialect(), conditions, binder)

	s.NoError(err)
	s.Equal("\"email_verified\" = $4 OR \"email_verified\" = $5", sql)
	s.Equal([]any{true, false}, args)
}

func (s *ConditionRendererTestSuite) TestRenderConditionsWithBinder_InvalidCondition() {

}

func TestConditionRendererTestSuite(t *testing.T) {
	suite.Run(t, new(ConditionRendererTestSuite))
}
