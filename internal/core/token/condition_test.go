// filename: /internal/core/token/condition_test.go

package token_test

import (
	"testing"
	"time"

	"github.com/ialopezg/entiqon/internal/core/builder"
	"github.com/ialopezg/entiqon/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type ConditionTestSuite struct {
	suite.Suite
}

func (s *ConditionTestSuite) TestSet_WithParams() {
	c := token.NewCondition(token.ConditionSimple, "active", true)
	s.Equal("active", c.Key)
	s.Equal([]any{true}, c.Values)
	s.Contains(c.Raw, "active = :active")
	s.True(c.IsValid())
}

func (s *ConditionTestSuite) TestSet_WithoutParams() {
	c := token.NewCondition(token.ConditionSimple, "active = true")
	s.Equal("active", c.Key)
	s.Equal([]any{true}, c.Values)
	s.Contains(c.Raw, "active = :active")
	s.True(c.IsValid())
}

func (s *ConditionTestSuite) TestIsValid() {
	invalid := token.NewCondition(token.ConditionSimple, "active")
	s.False(invalid.IsValid())

	empty := token.NewCondition(token.ConditionSimple, "")
	s.False(empty.IsValid())
}

func (s *ConditionTestSuite) TestAppendCondition() {
	var conditions []token.Condition
	conditions = append(conditions, token.NewCondition(token.ConditionSimple, "id", 1))
	conditions = append(conditions, token.NewCondition(token.ConditionAnd, "status", "active"))
	s.Len(conditions, 2)
	s.True(conditions[0].IsValid())
	s.True(conditions[1].IsValid())
}

func (s *ConditionTestSuite) TestAllConstructors() {
	now := time.Now()
	s.Run("NewCondition_WithThreeParams", func() {
		c := token.NewCondition(token.ConditionSimple, "value", 1, 2, 3)
		s.True(c.IsValid())
		s.Nil(c.Error)
	})
	s.Run("EmptyInlineValue", func() {
		c := token.NewCondition(token.ConditionSimple, "status = ")
		s.False(c.IsValid())
		s.Contains(c.Error.Error(), "unable to parse condition")
	})
	s.Run("Between", func() {
		c := token.NewConditionBetween(token.ConditionSimple, "created_at", now.Add(-time.Hour), now)
		s.True(c.IsValid())
		s.Equal("BETWEEN", c.Operator)
	})
	s.Run("BetweenEmptyString", func() {
		c := token.NewConditionBetween(token.ConditionSimple, "created_at", nil, nil)
		s.False(c.IsValid())
	})
	s.Run("Between_EmptyStart", func() {
		c := token.NewConditionBetween(token.ConditionSimple, "created_at", "", "2024-01-01")
		s.False(c.IsValid())
		s.Contains(c.Error.Error(), "start value cannot be empty")
	})
	s.Run("Between_EmptyEnd", func() {
		c := token.NewConditionBetween(token.ConditionSimple, "created_at", "2024-01-01", "")
		s.False(c.IsValid())
		s.Contains(c.Error.Error(), "end value cannot be empty")
	})
	s.Run("Between_IncompatibleTypes", func() {
		c := token.NewConditionBetween(token.ConditionSimple, "created_at", "2024-01-01", 42)
		s.False(c.IsValid())
		s.Contains(c.Error.Error(), "compatible types")
	})
	s.Run("Equal", func() {
		c := token.NewCondition(token.ConditionSimple, "id", 123)
		s.True(c.IsValid())
		s.Equal("=", c.Operator)
	})
	s.Run("GreaterThan", func() {
		c := token.NewConditionGreaterThan(token.ConditionSimple, "score", 80)
		s.True(c.IsValid())
		s.Equal(">", c.Operator)
	})
	s.Run("GreaterThanOrEqual", func() {
		c := token.NewConditionGreaterThanOrEqual(token.ConditionSimple, "points", 100)
		s.True(c.IsValid())
		s.Equal(">=", c.Operator)
	})
	s.Run("In", func() {
		c := token.NewConditionIn(token.ConditionSimple, "region", "US", "CA")
		s.True(c.IsValid())
		s.Equal("IN", c.Operator)
	})
	s.Run("In_IncompatibleTypes", func() {
		c := token.NewConditionIn(token.ConditionSimple, "amount", 10, "twenty")
		s.False(c.IsValid())
	})
	s.Run("InvalidWithOperator_MissingField", func() {
		c := token.NewConditionWithOperator(token.ConditionSimple, "", "=", 1)
		s.False(c.IsValid())
	})
	s.Run("InvalidWithOperator_MissingValues", func() {
		c := token.NewConditionWithOperator(token.ConditionSimple, "status", "=")
		s.False(c.IsValid())
	})
	s.Run("LessThan", func() {
		c := token.NewConditionLessThan(token.ConditionSimple, "age", 65)
		s.True(c.IsValid())
		s.Equal("<", c.Operator)
	})
	s.Run("LessThanOrEqual", func() {
		c := token.NewConditionLessThanOrEqual(token.ConditionSimple, "price", 100.0)
		s.True(c.IsValid())
		s.Equal("<=", c.Operator)
	})
	s.Run("Like", func() {
		c := token.NewConditionLike(token.ConditionSimple, "name", "%John%")
		s.True(c.IsValid())
		s.Equal("LIKE", c.Operator)
	})
	s.Run("NotEqual", func() {
		c := token.NewConditionNotEqual(token.ConditionSimple, "status", "archived")
		s.True(c.IsValid())
		s.Equal("!=", c.Operator)
	})
	s.Run("NotIn", func() {
		c := token.NewConditionNotIn(token.ConditionSimple, "status", "inactive", "banned")
		s.True(c.IsValid())
		s.Equal("NOT IN", c.Operator)
	})
	s.Run("NotIn_IncompatibleTypes", func() {
		c := token.NewConditionNotIn(token.ConditionSimple, "id", "x", 1, true)
		s.False(c.IsValid())
		s.Contains(c.Error.Error(), "compatible types")
	})
}

func (s *ConditionTestSuite) TestAppendCondition_Mixed() {
	initial := []token.Condition{
		token.NewCondition(token.ConditionSimple, "id", 1),
	}

	valid := token.NewCondition(token.ConditionAnd, "status", "active")
	invalid := token.NewCondition(token.ConditionOr, "") // invalid: missing field

	// Append valid → should grow
	result := builder.AppendCondition(initial, valid)
	s.Len(result, 2)
	s.Equal("status", result[1].Key)

	// Append invalid → should not grow
	result = builder.AppendCondition(result, invalid)
	s.Len(result, 2) // unchanged

	// Check validity flags
	s.True(result[0].IsValid())
	s.True(result[1].IsValid())
	s.False(invalid.IsValid())
}

func (s *ConditionTestSuite) TestAreCompatibleTypes() {
	s.Run("AllStrings", func() {
		s.True(token.AreCompatibleTypes("a", "b", "c"))
	})

	s.Run("AllInts", func() {
		s.True(token.AreCompatibleTypes(1, 2, 3))
	})

	s.Run("IntAndFloat", func() {
		s.True(token.AreCompatibleTypes(1, 2.5))
	})

	s.Run("AllTime", func() {
		t1 := time.Now()
		t2 := t1.Add(time.Hour)
		s.True(token.AreCompatibleTypes(t1, t2))
	})

	s.Run("MixedStringAndInt", func() {
		s.False(token.AreCompatibleTypes("a", 1))
	})

	s.Run("MixedTimeAndString", func() {
		s.False(token.AreCompatibleTypes(time.Now(), "2023-01-01"))
	})

	s.Run("EmptyInput", func() {
		s.False(token.AreCompatibleTypes())
	})

	s.Run("SingleItem", func() {
		s.False(token.AreCompatibleTypes("only-one"))
	})

	s.Run("NilValues", func() {
		var a any
		var b any = nil
		s.False(token.AreCompatibleTypes(a, b))
	})
}

func TestConditionTestSuite(t *testing.T) {
	suite.Run(t, new(ConditionTestSuite))
}
