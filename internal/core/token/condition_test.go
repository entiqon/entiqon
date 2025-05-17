package token_test

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type ConditionTokenTestSuite struct {
	suite.Suite
}

func (s *ConditionTokenTestSuite) TestSet_WithParams() {
	c := token.Condition{}.Set(token.ConditionAnd, "id = ?", 123)
	s.Equal(token.ConditionAnd, c.Type)
	s.Equal("id = ?", c.Key)
	s.Equal([]any{123}, c.Params)
	s.Contains(c.Raw, "123")
	fmt.Printf("📦 Field resolved with %+v\n", c)
}

func (s *ConditionTokenTestSuite) TestSet_WithoutParams() {
	c := token.Condition{}.Set(token.ConditionSimple, "active = true")
	s.Equal("active = true", c.Key)
	s.Empty(c.Params)
	s.Contains(c.Raw, "active = true")
	fmt.Printf("📦 Field resolved with %+v\n", c)
}

func (s *ConditionTokenTestSuite) TestIsValid() {
	s.True(token.NewCondition(token.ConditionSimple, "status = ?", "active").IsValid())
	s.False(token.NewCondition(token.ConditionSimple, "", "active").IsValid())
}

func (s *ConditionTokenTestSuite) TestAppendCondition() {
	var initial []token.Condition
	cond := token.NewCondition(token.ConditionAnd, "status = ?", "active")
	result := token.AppendCondition(initial, cond)
	s.Len(result, 1)

	var invalids []token.Condition
	invalid := token.NewCondition(token.ConditionSimple, "")
	result2 := token.AppendCondition(result, invalid)
	if !invalid.IsValid() {
		invalids = append(invalids, invalid)
	}
	s.Len(result2, 1) // should not grow
	fmt.Printf("📦 Fields resolved with initial: %d, result: %d, last: %d, invalid: %d\n", len(initial), len(result), len(result2), len(invalids))
}

func TestConditionTokenTestSuite(t *testing.T) {
	suite.Run(t, new(ConditionTokenTestSuite))
}
