package builder

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ialopezg/entiqon/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type DeleteBuilderTestSuite struct {
	suite.Suite
	qb *DeleteBuilder
}

func (s *DeleteBuilderTestSuite) SetupTest() {
	s.qb = NewDelete()
}

// ─────────────────────────────────────────────
// 🧪 From
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestFrom_SetsTargetTable() {
	sql, _, err := s.qb.From("users").Build()
	s.NoError(err)
	s.Contains(sql, "DELETE")
}

// ─────────────────────────────────────────────
// 🧪 Where
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestWhere_AddsSimpleCondition() {
	sql, args, err := s.qb.
		From("users").
		Where("id", 100).
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE id = ?")
	s.Equal([]any{100}, args)
}

func (s *SelectBuilderTestSuite) TestWhere_InvalidCondition() {
	s.qb.Select("*").From("users").Where("status =")
	_, _, err := s.qb.Build()

	s.Error(err)
	s.Contains(err.Error(), "1 invalid condition(s)")
}

// ─────────────────────────────────────────────
// 🧪 AndWhere
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestAndWhere_AppendsAND() {
	sql, args, err := s.qb.
		From("users").
		Where("active = true").
		AndWhere("role = 'admin'").
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE active = ? AND role = ?")
	s.Equal(2, len(args))
}

// ─────────────────────────────────────────────
// 🧪 OrWhere
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestOrWhere_AppendsOR() {
	sql, args, err := s.qb.
		From("users").
		Where("email_verified = false").
		OrWhere("email_verified = true").
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE email_verified = ? OR email_verified = ?")
	s.Equal(2, len(args))
}

// ─────────────────────────────────────────────
// 🧪 UseDialect
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestDeleteBuilderUseDialectPostgres() {
	sql, args, err := s.qb.
		Select("id", "created_at").
		From("users").
		Where("status", "active").
		UseDialect("postgres").
		Build()

	expectedSQL := `SELECT "id", "created_at" FROM "users" WHERE "status" = $1`
	s.NoError(err)
	s.Equal(expectedSQL, sql)
	s.Equal([]any{"active"}, args)
}

// ─────────────────────────────────────────────
// 🧪 Build
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestBuild_UseDialect() {
	sql, _, err := NewDelete().
		From("users").
		UseDialect("postgres").
		Build()

	s.NoError(err)
	s.Contains(sql, `DELETE FROM "users"`)
}

func (s *DeleteBuilderTestSuite) TestBuild_MissingFrom_ReturnsError() {
	_, _, err := NewDelete().
		Where("id = ?", 5).
		Build()

	s.Error(err)
	s.Contains(err.Error(), "requires a target table")
}

func (s *DeleteBuilderTestSuite) TestBuild_InvalidConditionType_ReturnsError() {
	qb := NewDelete().
		From("users").
		UseDialect("postgres")

	_, _, err := qb.Build()
	s.Nil(err)
	qb.conditions = append(qb.conditions, token.Condition{
		Type: "💥", Key: "status = 'active'",
	})

	_, _, err = qb.Build()
	s.Error(err)
	s.Contains(err.Error(), "unsupported condition type")
}

func (s *DeleteBuilderTestSuite) TestDeleteBuilder_LimitClause() {
	sql, _, _ := s.qb.
		From("logs").
		Where("archived", true).
		Limit(100).
		Build()

	if !strings.Contains(sql, "LIMIT 100") {
		s.Failf("TestDeleteBuilder_LimitClause", "expected LIMIT clause, got: %s", sql)
	}
}

func (s *DeleteBuilderTestSuite) TestBuild_WithInvalidConditions_ShouldFail() {
	s.qb.AddStageError("WHERE", errors.New("invalid condition"))
	_, _, err := s.qb.From("users").Build()

	s.Require().Error(err)
	s.Contains(err.Error(), "invalid condition")
}

func (s *DeleteBuilderTestSuite) TestWhere_InvalidCondition_ShouldAppendError() {
	s.qb.From("users").Where("", 123) // empty condition string

	errs := s.qb.GetErrors()

	s.Require().Len(errs, 1)
	s.Equal("WHERE", errs[0].Token)
	s.Contains(errs[0].Errors[0].Error(), "invalid")
}

func (s *DeleteBuilderTestSuite) TestAndWhere_InvalidCondition_ShouldAppendError() {
	s.qb.From("users").AndWhere("", 123) // Invalid condition

	errs := s.qb.GetErrors()

	s.Require().Len(errs, 1)
	s.Equal("WHERE", errs[0].Token)
	s.Contains(errs[0].Errors[0].Error(), "invalid")
}

func (s *DeleteBuilderTestSuite) TestOrWhere_InvalidCondition_ShouldAppendError() {
	s.qb.From("users").OrWhere("", 123) // Invalid condition

	errs := s.qb.GetErrors()

	s.Require().Len(errs, 1)
	s.Equal("WHERE", errs[0].Token)
	s.Contains(errs[0].Errors[0].Error(), "invalid")
}

func (s *DeleteBuilderTestSuite) TestBuild_BuildValidations() {
	c := token.NewCondition(token.ConditionSimple, "id = ?")

	b := DeleteBuilder{}
	s.Run("EmptyTable", func() {
		_, _, err := b.Build()
		s.Error(err)
		s.ErrorContains(err, "requires a target table")
	})
	if !c.IsValid() {
		b.AddStageError("WHERE clause", fmt.Errorf("invalid clause"))
	}
	b.From("users")
	s.Run("HasDialect", func() {
		b.conditions = []token.Condition{c}
		_, _, err := b.Build()
		s.Error(err)
		s.Equal("generic", b.dialect.Name())
	})
	s.Run("HasErrors", func() {
		_, _, err := b.Build()
		s.Error(err)
		s.Contains(err.Error(), "invalid condition(s)")
	})

}

func TestDeleteBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteBuilderTestSuite))
	fmt.Println("🧪 DELETE tests complete, Holmes.")
}
