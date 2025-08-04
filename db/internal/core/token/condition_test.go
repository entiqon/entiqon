// File: db/internal/core/token/condition_test.go

package token_test

import (
	"testing"
	"time"

	"github.com/entiqon/db/internal/core/builder"
	token2 "github.com/entiqon/db/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type ConditionTestSuite struct {
	suite.Suite
}

func (s *ConditionTestSuite) TestSet_WithParams() {
	c := token2.NewCondition(token2.ConditionSimple, "active", true)
	s.Equal("active", c.Key)
	s.Equal([]any{true}, c.Values)
	s.Contains(c.Raw, "active = :active")
	s.True(c.IsValid())
}

func (s *ConditionTestSuite) TestSet_WithoutParams() {
	c := token2.NewCondition(token2.ConditionSimple, "active = true")
	s.Equal("active", c.Key)
	s.Equal([]any{true}, c.Values)
	s.Contains(c.Raw, "active = :active")
	s.True(c.IsValid())
}

func (s *ConditionTestSuite) TestIsValid() {
	invalid := token2.NewCondition(token2.ConditionSimple, "active")
	s.False(invalid.IsValid())

	empty := token2.NewCondition(token2.ConditionSimple, "")
	s.False(empty.IsValid())
}

func (s *ConditionTestSuite) TestAppendCondition() {
	var conditions []token2.Condition
	conditions = append(conditions, token2.NewCondition(token2.ConditionSimple, "id", 1))
	conditions = append(conditions, token2.NewCondition(token2.ConditionAnd, "status", "active"))
	s.Len(conditions, 2)
	s.True(conditions[0].IsValid())
	s.True(conditions[1].IsValid())
}

func (s *ConditionTestSuite) TestAllConstructors() {
	now := time.Now()
	s.Run("NewCondition_WithThreeParams", func() {
		c := token2.NewCondition(token2.ConditionSimple, "value", 1, 2, 3)
		s.True(c.IsValid())
		s.Nil(c.Error)
	})
	s.Run("EmptyInlineValue", func() {
		c := token2.NewCondition(token2.ConditionSimple, "status = ")
		s.False(c.IsValid())
		s.Contains(c.Error.Error(), "unable to parse condition")
	})
	s.Run("Between", func() {
		c := token2.NewConditionBetween(token2.ConditionSimple, "created_at", now.Add(-time.Hour), now)
		s.True(c.IsValid())
		s.Equal("BETWEEN", c.Operator)
	})
	s.Run("BetweenEmptyString", func() {
		c := token2.NewConditionBetween(token2.ConditionSimple, "created_at", nil, nil)
		s.False(c.IsValid())
	})
	s.Run("Between_EmptyStart", func() {
		c := token2.NewConditionBetween(token2.ConditionSimple, "created_at", "", "2024-01-01")
		s.False(c.IsValid())
		s.Contains(c.Error.Error(), "start value cannot be empty")
	})
	s.Run("Between_EmptyEnd", func() {
		c := token2.NewConditionBetween(token2.ConditionSimple, "created_at", "2024-01-01", "")
		s.False(c.IsValid())
		s.Contains(c.Error.Error(), "end value cannot be empty")
	})
	s.Run("Between_IncompatibleTypes", func() {
		c := token2.NewConditionBetween(token2.ConditionSimple, "created_at", "2024-01-01", 42)
		s.False(c.IsValid())
		s.Contains(c.Error.Error(), "compatible types")
	})
	s.Run("Equal", func() {
		c := token2.NewCondition(token2.ConditionSimple, "id", 123)
		s.True(c.IsValid())
		s.Equal("=", c.Operator)
	})
	s.Run("GreaterThan", func() {
		c := token2.NewConditionGreaterThan(token2.ConditionSimple, "score", 80)
		s.True(c.IsValid())
		s.Equal(">", c.Operator)
	})
	s.Run("GreaterThanOrEqual", func() {
		c := token2.NewConditionGreaterThanOrEqual(token2.ConditionSimple, "points", 100)
		s.True(c.IsValid())
		s.Equal(">=", c.Operator)
	})
	s.Run("In", func() {
		c := token2.NewConditionIn(token2.ConditionSimple, "region", "US", "CA")
		s.True(c.IsValid())
		s.Equal("IN", c.Operator)
	})
	s.Run("In_IncompatibleTypes", func() {
		c := token2.NewConditionIn(token2.ConditionSimple, "amount", 10, "twenty")
		s.False(c.IsValid())
	})
	s.Run("InvalidWithOperator_MissingField", func() {
		c := token2.NewConditionWithOperator(token2.ConditionSimple, "", "=", 1)
		s.False(c.IsValid())
	})
	s.Run("InvalidWithOperator_MissingValues", func() {
		c := token2.NewConditionWithOperator(token2.ConditionSimple, "status", "=")
		s.False(c.IsValid())
	})
	s.Run("LessThan", func() {
		c := token2.NewConditionLessThan(token2.ConditionSimple, "age", 65)
		s.True(c.IsValid())
		s.Equal("<", c.Operator)
	})
	s.Run("LessThanOrEqual", func() {
		c := token2.NewConditionLessThanOrEqual(token2.ConditionSimple, "price", 100.0)
		s.True(c.IsValid())
		s.Equal("<=", c.Operator)
	})
	s.Run("Like", func() {
		c := token2.NewConditionLike(token2.ConditionSimple, "name", "%John%")
		s.True(c.IsValid())
		s.Equal("LIKE", c.Operator)
	})
	s.Run("NotEqual", func() {
		c := token2.NewConditionNotEqual(token2.ConditionSimple, "status", "archived")
		s.True(c.IsValid())
		s.Equal("!=", c.Operator)
	})
	s.Run("NotIn", func() {
		c := token2.NewConditionNotIn(token2.ConditionSimple, "status", "inactive", "banned")
		s.True(c.IsValid())
		s.Equal("NOT IN", c.Operator)
	})
	s.Run("NotIn_IncompatibleTypes", func() {
		c := token2.NewConditionNotIn(token2.ConditionSimple, "id", "x", 1, true)
		s.False(c.IsValid())
		s.Contains(c.Error.Error(), "compatible types")
	})
}

func (s *ConditionTestSuite) TestAppendCondition_Mixed() {
	initial := []token2.Condition{
		token2.NewCondition(token2.ConditionSimple, "id", 1),
	}

	valid := token2.NewCondition(token2.ConditionAnd, "status", "active")
	invalid := token2.NewCondition(token2.ConditionOr, "") // invalid: missing field

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
		s.True(token2.AreCompatibleTypes("a", "b", "c"))
	})

	s.Run("AllInts", func() {
		s.True(token2.AreCompatibleTypes(1, 2, 3))
	})

	s.Run("IntAndFloat", func() {
		s.True(token2.AreCompatibleTypes(1, 2.5))
	})

	s.Run("AllTime", func() {
		t1 := time.Now()
		t2 := t1.Add(time.Hour)
		s.True(token2.AreCompatibleTypes(t1, t2))
	})

	s.Run("MixedStringAndInt", func() {
		s.False(token2.AreCompatibleTypes("a", 1))
	})

	s.Run("MixedTimeAndString", func() {
		s.False(token2.AreCompatibleTypes(time.Now(), "2023-01-01"))
	})

	s.Run("EmptyInput", func() {
		s.False(token2.AreCompatibleTypes())
	})

	s.Run("SingleItem", func() {
		s.False(token2.AreCompatibleTypes("only-one"))
	})

	s.Run("NilValues", func() {
		var a any
		var b any = nil
		s.False(token2.AreCompatibleTypes(a, b))
	})
}

func TestConditionTestSuite(t *testing.T) {
	suite.Run(t, new(ConditionTestSuite))
}
