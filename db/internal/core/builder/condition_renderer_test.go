// File: db/internal/core/builder/condition_renderer_test.go

package builder_test

import (
	"testing"

	driver2 "github.com/entiqon/entiqon/db/driver"
	"github.com/entiqon/entiqon/db/internal/core/builder"
	"github.com/entiqon/entiqon/db/internal/core/builder/bind"
	token2 "github.com/entiqon/entiqon/db/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type ConditionRendererTestSuite struct {
	suite.Suite
}

func (s *ConditionRendererTestSuite) TestAppendCondition() {
	base := token2.NewCondition(token2.ConditionSimple, "status", "active")
	and := token2.NewCondition(token2.ConditionAnd, "deleted", false)
	or := token2.NewCondition(token2.ConditionOr, "archived", false)

	result := builder.AppendCondition([]token2.Condition{base}, and)
	s.Len(result, 2)
	s.Equal(token2.ConditionAnd, result[1].Type)

	result = builder.AppendCondition(result, or)
	s.Len(result, 3)
	s.Equal(token2.ConditionOr, result[2].Type)

	s.Run("Invalid", func() {
		valid := token2.NewCondition(token2.ConditionSimple, "status", "active")
		invalid := token2.Condition{} // no field/operator/values → invalid

		result := builder.AppendCondition([]token2.Condition{valid}, invalid)

		s.Len(result, 1)
		s.Equal("status", result[0].Key)
	})
}

func (s *ConditionRendererTestSuite) TestRenderConditions_Generic() {
	s.Run("Valid", func() {
		conditions := []token2.Condition{
			token2.NewCondition(token2.ConditionSimple, "active", true),
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
		c := token2.NewCondition(token2.ConditionSimple, "status", "active")
		c.Type = token2.ConditionType(rune(999)) // simulate unsupported condition type

		binder := bind.NewParamBinder(driver2.NewGenericDialect())
		_, _, err := builder.RenderConditionsWithBinder(driver2.NewGenericDialect(), []token2.Condition{c}, binder)

		s.Error(err)
		s.Contains(err.Error(), "unsupported condition type")
	})

	s.Run("Invalid", func() {
		conditions := []token2.Condition{
			{}, // invalid: missing Key and Error is nil
		}

		binder := bind.NewParamBinder(driver2.NewGenericDialect())
		_, _, err := builder.RenderConditionsWithBinder(driver2.NewGenericDialect(), conditions, binder)

		s.Error(err)
		s.Contains(err.Error(), "invalid condition")
	})

	s.Run("WithAndCondition", func() {
		conditions := []token2.Condition{
			token2.NewCondition(token2.ConditionSimple, "status", "active"),
			token2.NewCondition(token2.ConditionAnd, "deleted", false),
		}

		binder := bind.NewParamBinder(driver2.NewGenericDialect())
		sql, args, err := builder.RenderConditionsWithBinder(driver2.NewGenericDialect(), conditions, binder)

		s.NoError(err)
		s.Equal("status = ? AND deleted = ?", sql)
		s.Equal([]any{"active", false}, args)
	})
}

func (s *ConditionRendererTestSuite) TestRenderConditionsWithBinder_Postgres() {
	conditions := []token2.Condition{
		token2.NewCondition(token2.ConditionSimple, "email_verified", true),
		token2.NewCondition(token2.ConditionOr, "email_verified", false),
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
