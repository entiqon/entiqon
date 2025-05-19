package builder

import (
	"testing"

	"github.com/ialopezg/entiqon/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type SelectBuilderTestSuite struct {
	suite.Suite
	qb *SelectBuilder
}

func (s *SelectBuilderTestSuite) SetupTest() {
	s.qb = NewSelect()
}

// ─────────────────────────────────────────────
// 🧪 Select
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestSelectBasicColumns() {
	sql, _, err := s.qb.
		Select("id", "name").
		From("users").
		Build()

	expected := "SELECT id, name FROM users"

	s.NoError(err)
	s.Equal(expected, sql)
}

func (s *SelectBuilderTestSuite) TestSelectCommaSeparated() {
	sql, _, err := s.qb.
		Select("id, name").
		From("users").
		Build()

	s.NoError(err)
	s.Equal("SELECT id, name FROM users", sql)
}

func (s *SelectBuilderTestSuite) TestSelectInlineAlias() {
	sql, _, err := s.qb.
		Select("email AS contact").
		From("users").
		Build()

	s.NoError(err)
	s.Equal("SELECT email AS contact FROM users", sql)
}

// ─────────────────────────────────────────────
// 🧪 AddSelect
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestSelect_AddSelectAppends() {
	sql, _, err := s.qb.
		Select("id").
		AddSelect("name AS full_name").
		From("users").
		Build()

	expected := "SELECT id, name AS full_name FROM users"
	s.NoError(err)
	s.Equal(expected, sql)
}

// ─────────────────────────────────────────────
// 🧪 From
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestFrom_SingleTable() {
	sql, _, err := s.qb.
		Select("id").
		From("customers").
		Build()

	expected := "SELECT id FROM customers"

	s.NoError(err)
	s.Equal(expected, sql)
}

func (s *SelectBuilderTestSuite) TestMissingFromClause() {
	sql, _, err := s.qb.
		Select("id").
		Build()

	s.Error(err)
	s.Empty(sql)
}

// ─────────────────────────────────────────────
// 🧪 Where, AndWhere, OrWhere
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestWhereAndOrConditions() {
	sql, _, err := s.qb.
		Select("id").
		From("customers").
		Where("active = ?", true).
		AndWhere("email_verified = ?", true).
		OrWhere("country = ?", "US").
		OrWhere("country = ?", "CA").
		Build()

	expected := "SELECT id FROM customers WHERE active = ? AND email_verified = ? OR country = ? OR country = ?"
	s.NoError(err)
	s.Equal(expected, sql)
}

func (s *SelectBuilderTestSuite) TestGroupedAndWhere() {
	sql, _, err := s.qb.
		From("invoices").
		Where("paid = ?", false).
		AndWhere("amount > ?", 100).
		AndWhere("overdue = ?", true).
		Build()

	expected := "SELECT * FROM invoices WHERE paid = ? AND amount > ? AND overdue = ?"
	s.NoError(err)
	s.Equal(expected, sql)
}

func (s *SelectBuilderTestSuite) TestSelectBuilderMultiParams() {
	sql, params, err := s.qb.
		Select("id", "email").
		From("users").
		Where("status = ?", "active").
		AndWhere("role = ?", "admin").
		AndWhere("created_at > ? AND region = ?", "2024-01-01", "NA").
		OrderBy("last_login DESC").
		Take(50).
		Skip(0).
		Build()

	expected := "SELECT id, email FROM users WHERE status = ? AND role = ? AND created_at > ? AND region = ? ORDER BY last_login DESC LIMIT 50 OFFSET 0"

	s.NoError(err)
	s.Equal(expected, sql)
	s.Equal([]any{"active", "admin", "2024-01-01", "NA"}, params)
}

// ─────────────────────────────────────────────
// 🧪 OrderBy, Take, Skip
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestOrderingTakeSkip() {
	sql, _, err := s.qb.
		Select("name").
		From("employees").
		OrderBy("created_at DESC").
		Take(10).
		Skip(5).
		Build()

	expected := "SELECT name FROM employees ORDER BY created_at DESC LIMIT 10 OFFSET 5"
	s.NoError(err)
	s.Equal(expected, sql)
}

// ─────────────────────────────────────────────
// 🧪 Build
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestBuild_InvalidConditionType() {
	b := NewSelect().
		Select("id").
		From("users")

	// Create a rogue condition
	rogue := token.Condition{
		Type: "💣",
		Key:  "is_admin = true",
	}

	// Directly inject invalid condition
	b.conditions = append(b.conditions, rogue)

	_, _, err := b.Build()
	s.Error(err)
	s.Contains(err.Error(), "invalid condition type")
}

func (s *SelectBuilderTestSuite) TestBuild_WithoutDialect_UsesRawLimitOffset() {
	sb := s.qb.
		Select("id").
		From("users").
		Take(20).
		Skip(10) // deliberately no .WithDialect()

	sql, args, err := sb.Build()

	s.Require().NoError(err)
	s.Contains(sql, "LIMIT 20")
	s.Contains(sql, "OFFSET 10")
	s.Empty(args)
}

func (s *SelectBuilderTestSuite) TestBuild_WithDialect_UsesDialectLimitOffset() {
	sb := s.qb.
		Select("id").
		From("users").
		Take(10).
		Skip(5).
		UseDialect("postgres")

	sql, args, err := sb.Build()

	s.Require().NoError(err)
	s.Contains(sql, "LIMIT 10 OFFSET 5")
	s.Empty(args)
}

// ─────────────────────────────────────────────
// 🧪 UseDialect
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestSelectBuilderUseDialectPostgres() {
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
func (s *SelectBuilderTestSuite) TestSelectBuilderWithDialect() {
	b := NewSelect().
		From("users").
		WithDialect("postgres")
	sql, args, err := b.Build()

	s.Require().NoError(err)
	s.Contains(sql, `SELECT * FROM`)
	s.Empty(args)
}

// ─────────────────────────────────────────────
// 🧪 UseDialect
// ─────────────────────────────────────────────

func TestSelectBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(SelectBuilderTestSuite))
}
