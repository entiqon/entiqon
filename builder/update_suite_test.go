package builder

import (
	"testing"

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
	sql, _, err := s.qb.
		Table("users").
		Set("status", "active").
		Build()

	s.NoError(err)
	s.Contains(sql, "UPDATE users")
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
}

// 🧪 Set (Multiple)
func (s *UpdateBuilderTestSuite) TestSet_MultipleAssignments() {
	sql, args, err := s.qb.
		Table("users").
		Set("name", "Alice").
		Set("status", "verified").
		Build()

	s.NoError(err)
	s.Contains(sql, "SET name = ?")
	s.Contains(sql, "status = ?")
	s.Equal([]any{"Alice", "verified"}, args)
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
}

// ─────────────────────────────────────────────
// 🧪 AndWhere
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestAndWhere_AppendsAndCondition() {
	sql, _, err := s.qb.
		Table("users").
		Set("status", "inactive").
		Where("deleted = false").
		AndWhere("role = ?", "admin").
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE deleted = false AND role = ?")
}

// ─────────────────────────────────────────────
// 🧪 OrWhere
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestOrWhere_AppendsOrCondition() {
	sql, _, err := s.qb.
		Table("users").
		Set("active", true).
		Where("email_verified = true").
		OrWhere("email_verified = false").
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE email_verified = true OR email_verified = false")
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
}

func (s *UpdateBuilderTestSuite) TestBuild_MissingTableReturnsError() {
	_, _, err := s.qb.
		Set("name", "Watson").
		Build()

	s.Error(err)
	s.Contains(err.Error(), "UPDATE requires a target table")
}

func (s *UpdateBuilderTestSuite) TestBuild_MissingAssignmentsReturnsError() {
	_, _, err := s.qb.
		Table("users").
		Build()

	s.Error(err)
	s.Contains(err.Error(), "UPDATE must define at least one column assignment")
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

	_, _, err := q.Build()
	s.Error(err)
	s.Contains(err.Error(), "invalid condition type")
}

// ─────────────────────────────────────────────
// 🧪 UseDialect
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestUpdateBuilder_UseDialect_Postgres() {
	sql, args, err := s.qb.
		Set("active", true).
		Table("users").
		Where("email_verified = true").
		OrWhere("email_verified = false").
		UseDialect("postgres").
		Build()

	s.NoError(err)
	s.Equal([]any{true}, args)
	s.Contains(sql, "WHERE email_verified = true OR email_verified = false")
}

// ─────────────────────────────────────────────
// 🧪 WithDialect
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestUpdateBuilder_WithDialect_Postgres() {
	sql, args, err := s.qb.
		Set("active", true).
		Table("users").
		Where("email_verified = true").
		OrWhere("email_verified = false").
		WithDialect("postgres").
		Build()

	s.NoError(err)
	s.Equal([]any{true}, args)
	s.Contains(sql, "WHERE email_verified = true OR email_verified = false")
}

func TestUpdateBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateBuilderTestSuite))
}
