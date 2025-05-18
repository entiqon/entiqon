package builder

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/internal/core/dialect"
	"github.com/ialopezg/entiqon/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type UpdateBuilderTestSuite struct {
	suite.Suite
	qb *UpdateBuilder
}

func (s *UpdateBuilderTestSuite) SetupTest() {
	s.qb = NewUpdate()
}

// ─────────────────────────────────────────────
// 🧪 Table
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestTable_SetsTableName() {
	sql, args, err := s.qb.
		Table("users").
		Set("status", "active").
		Build()

	s.NoError(err)
	s.Contains(sql, "UPDATE users")
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 Set
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestSet_AppendsAssignment() {
	sql, args, err := s.qb.
		Table("users").
		Set("status", "active").
		Build()

	s.NoError(err)
	s.Contains(sql, "SET status = ?")
	s.Equal([]any{"active"}, args)
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 Where
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestWhere_SetsInitialCondition() {
	sql, args, err := s.qb.
		Table("users").
		Set("name", "Watson").
		Where("id = ?", 42).
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE id = ?")
	s.Equal([]any{"Watson", 42}, args)
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 AndWhere
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestAndWhere_AppendsAndCondition() {
	sql, args, err := s.qb.
		Table("users").
		Set("status", "inactive").
		Where("deleted = false").
		AndWhere("role = ?", "admin").
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE deleted = false AND role = ?")
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 OrWhere
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestOrWhere_AppendsOrCondition() {
	sql, args, err := s.qb.
		Table("users").
		Set("active", true).
		Where("email_verified = true").
		OrWhere("email_verified = false").
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE email_verified = true OR email_verified = false")
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 Build
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestBuild_WithAliasedColumn() {
	sql, args, err := s.qb.
		Table("users").
		Set("email AS contact", "watson@example.com").
		Where("id = ?", 1).
		Build()

	s.Error(err)
	s.Equal(sql, "")
	s.Nil(args)
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

func (s *UpdateBuilderTestSuite) TestBuild_MissingTableReturnsError() {
	sql, args, err := s.qb.
		Set("name", "Watson").
		Build()

	s.Error(err)
	s.Contains(err.Error(), "UPDATE requires a target table")
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

func (s *UpdateBuilderTestSuite) TestBuild_MissingAssignmentsReturnsError() {
	sql, args, err := s.qb.
		Table("users").
		Build()

	s.Error(err)
	s.Contains(err.Error(), "UPDATE must define at least one column assignment")
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

func (s *UpdateBuilderTestSuite) TestBuild_InvalidConditionType_ReturnsError() {
	q := NewUpdate().
		Table("users").
		Set("name", "Watson")

	// Inject invalid condition
	q.Set("x", "y") // keep Set valid
	q.Table("users")
	q.conditions = append(q.conditions, token.Condition{
		Type: "💣", Key: "broken = true",
	})

	sql, args, err := q.Build()
	s.Error(err)
	s.Contains(err.Error(), "invalid condition type")
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 WithDialect
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestSelectBuilder_WithDialect_Postgres() {
	sql, args, err := s.qb.
		Set("active", true).
		Table("users").
		Where("email_verified = true").
		OrWhere("email_verified = false").
		WithDialect(&dialect.PostgresEngine{}).
		Build()

	s.NoError(err)
	s.Equal([]any{true}, args)
	s.Contains(sql, "WHERE email_verified = true OR email_verified = false")
	fmt.Printf("📦 WithDialect → SQL: %s | Args: %+v\n", sql, args)
}

func TestUpdateBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateBuilderTestSuite))
}
