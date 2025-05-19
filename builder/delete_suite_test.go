package builder

import (
	"fmt"
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
	s.Contains(sql, "DELETE FROM users")
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
}

// ─────────────────────────────────────────────
// 🧪 UseDialect
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestDeleteBuilderUseDialectPostgres() {
	sql, args, err := s.qb.
		Select("id", "created_at").
		From("users").
		Where("status = ?", "active").
		UseDialect("postgres").
		Build()

	expectedSQL := `SELECT "id", "created_at" FROM "users" WHERE "status" = ?`
	s.NoError(err)
	s.Equal(expectedSQL, sql)
	s.Equal([]any{"active"}, args)
}

// ─────────────────────────────────────────────
// 🧪 WithDialect
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestWithDialect_Deprecated_Works() {
	b := NewDelete().
		From("users").
		WithDialect("postgres")
	sql, args, err := b.Build()

	s.Require().NoError(err)
	s.Contains(sql, `DELETE FROM "users"`)
	s.Empty(args)
}

// ─────────────────────────────────────────────
// 🧪 Build
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestBuild_WithDialect() {
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
	s.Contains(err.Error(), "DELETE requires a target table")
}

func (s *DeleteBuilderTestSuite) TestBuild_InvalidConditionType_ReturnsError() {
	qb := NewDelete().
		From("users")

	qb.UseDialect("")
	qb.conditions = append(qb.conditions, token.Condition{
		Type: "💥", Key: "status = 'active'",
	})

	_, _, err := qb.Build()
	s.Error(err)
	s.Contains(err.Error(), "invalid condition type")
}

func TestDeleteBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteBuilderTestSuite))
	fmt.Println("🧪 DELETE tests complete, Holmes.")
}
