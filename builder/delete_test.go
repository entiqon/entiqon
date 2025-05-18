package builder

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/internal/core/dialect"
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
	sql, args, err := s.qb.From("users").Build()
	s.NoError(err)
	s.Contains(sql, "DELETE FROM users")
	fmt.Printf("📦 From → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 Where
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestWhere_AddsSimpleCondition() {
	sql, args, err := s.qb.
		From("users").
		Where("id = ?", 100).
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE id = ?")
	s.Equal([]any{100}, args)
	fmt.Printf("📦 Where → SQL: %s | Args: %+v\n", sql, args)
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
	s.Contains(sql, "WHERE active = true AND role = 'admin'")
	s.Equal(0, len(args)) // No parameterized args
	fmt.Printf("📦 AndWhere → SQL: %s | Args: %+v\n", sql, args)
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
	s.Contains(sql, "WHERE email_verified = false OR email_verified = true")
	s.Equal(0, len(args))
	fmt.Printf("📦 OrWhere → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 Build
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestBuild_WithDialect() {
	sql, args, err := NewDelete().
		From("users").
		WithDialect(&dialect.PostgresEngine{}).
		Build()

	s.NoError(err)
	s.Contains(sql, `DELETE FROM "users"`)
	fmt.Printf("📦 Build → SQL: %s | Args: %+v\n", sql, args)
}

func (s *DeleteBuilderTestSuite) TestBuild_MissingFrom_ReturnsError() {
	sql, args, err := NewDelete().
		Where("id = ?", 5).
		Build()

	s.Error(err)
	s.Contains(err.Error(), "DELETE requires a target table")
	fmt.Printf("📦 Build → SQL: %s | Args: %+v\n", sql, args)
}

func (s *DeleteBuilderTestSuite) TestBuild_InvalidConditionType_ReturnsError() {
	qb := NewDelete().
		From("users")

	qb.WithDialect(nil)
	qb.conditions = append(qb.conditions, token.Condition{
		Type: "💥", Key: "status = 'active'",
	})

	sql, args, err := qb.Build()
	s.Error(err)
	s.Contains(err.Error(), "invalid condition type")
	fmt.Printf("📦 Build → SQL: %s | Args: %+v\n", sql, args)
}

func TestDeleteBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteBuilderTestSuite))
	fmt.Println("🧪 DELETE tests complete, Holmes.")
}
