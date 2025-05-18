package builder

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/internal/core/dialect"
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
func (s *SelectBuilderTestSuite) TestSelect_BasicColumns() {
	sql, args, err := s.qb.
		Select("id", "name").
		From("users").
		Build()

	expected := "SELECT id, name FROM users"

	s.NoError(err)
	s.Equal(expected, sql)
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

func (s *SelectBuilderTestSuite) TestSelect_CommaSeparated() {
	sql, args, err := s.qb.
		Select("id, name").
		From("users").
		Build()

	s.NoError(err)
	s.Equal("SELECT id, name FROM users", sql)
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

func (s *SelectBuilderTestSuite) TestSelect_InlineAlias() {
	sql, args, err := s.qb.
		Select("email AS contact").
		From("users").
		Build()

	s.NoError(err)
	s.Equal("SELECT email AS contact FROM users", sql)
	fmt.Printf("📦 Select → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 AddSelect
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestSelect_AddSelectAppends() {
	sql, args, err := s.qb.
		Select("id").
		AddSelect("name AS full_name").
		From("users").
		Build()

	expected := "SELECT id, name AS full_name FROM users"
	s.NoError(err)
	s.Equal(expected, sql)
	fmt.Printf("📦 AddSelect → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 From
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestFrom_SingleTable() {
	sql, args, err := s.qb.
		Select("id").
		From("customers").
		Build()

	expected := "SELECT id FROM customers"

	s.NoError(err)
	s.Equal(expected, sql)
	fmt.Printf("📦 From → SQL: %s | Args: %+v\n", sql, args)
}

func (s *SelectBuilderTestSuite) TestMissingFromClause() {
	sql, params, err := s.qb.
		Select("id").
		Build()

	s.Error(err)
	s.Empty(sql)
	fmt.Printf("📦 Generated SQL Query: %s with params=%+v\n", sql, params)
}

// ─────────────────────────────────────────────
// 🧪 Where, AndWhere, OrWhere
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestWhereAndOrConditions() {
	sql, params, err := s.qb.
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
	fmt.Printf("📦 Generated SQL Query: %s with params=%+v\n", sql, params)
}

func (s *SelectBuilderTestSuite) TestGroupedAndWhere() {
	sql, params, err := s.qb.
		From("invoices").
		Where("paid = ?", false).
		AndWhere("amount > ?", 100).
		AndWhere("overdue = ?", true).
		Build()

	expected := "SELECT * FROM invoices WHERE paid = ? AND amount > ? AND overdue = ?"
	s.NoError(err)
	s.Equal(expected, sql)
	fmt.Printf("📦 Generated SQL Query: %s with params=%+v\n", sql, params)
}

func (s *SelectBuilderTestSuite) TestSelectBuilder_MultiParams() {
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

	fmt.Printf("📦 Generated SQL Query: %s with params=%+v\n", sql, params)
}

// ─────────────────────────────────────────────
// 🧪 OrderBy, Take, Skip
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestOrderingTakeSkip() {
	sql, params, err := s.qb.
		Select("name").
		From("employees").
		OrderBy("created_at DESC").
		Take(10).
		Skip(5).
		Build()

	expected := "SELECT name FROM employees ORDER BY created_at DESC LIMIT 10 OFFSET 5"
	s.NoError(err)
	s.Equal(expected, sql)
	fmt.Printf("📦 Generated SQL Query: %s with params=%+v\n", sql, params)
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

// ─────────────────────────────────────────────
// 🧪 WithDialect
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestSelectBuilder_WithDialect_Postgres() {
	sql, args, err := s.qb.
		Select("id", "created_at").
		From("users").
		Where("status = ?", "active").
		WithDialect(&dialect.PostgresEngine{}).
		Build()

	expectedSQL := `SELECT "id", "created_at" FROM users WHERE status = ?`
	s.NoError(err)
	s.Equal(expectedSQL, sql)
	s.Equal([]any{"active"}, args)

	fmt.Printf("📦 WithDialect → SQL: %s | Args: %+v\n", sql, args)
}

func TestSelectBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(SelectBuilderTestSuite))
}
