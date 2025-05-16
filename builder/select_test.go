package builder_test

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/builder"
	"github.com/stretchr/testify/suite"
)

type SelectQueryBuilderTestSuite struct {
	suite.Suite
	qb *builder.SelectBuilder
}

func (s *SelectQueryBuilderTestSuite) SetupTest() {
	s.qb = builder.NewSelect()
}

func (s *SelectQueryBuilderTestSuite) TestBasicSelect() {
	sql, params, err := s.qb.Select("id", "name").From("users").Build()
	s.NoError(err)
	s.Equal("SELECT id, name FROM users", sql)
	fmt.Printf("📦 Generated SQL Query: %s with params=%+v\n", sql, params)
}

func (s *SelectQueryBuilderTestSuite) TestWhereAndOrConditions() {
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

func (s *SelectQueryBuilderTestSuite) TestOrderingTakeSkip() {
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

func (s *SelectQueryBuilderTestSuite) TestMissingFromClause() {
	sql, params, err := s.qb.
		Select("id").
		Build()

	s.Error(err)
	s.Empty(sql)
	fmt.Printf("📦 Generated SQL Query: %s with params=%+v\n", sql, params)
}

func (s *SelectQueryBuilderTestSuite) TestGroupedAndWhere() {
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

func (s *SelectQueryBuilderTestSuite) TestSelectBuilder_MultiParams() {
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

func TestSelectQueryBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(SelectQueryBuilderTestSuite))
	fmt.Println("🕵️ Verified by Watson: All is sound in the SELECT logic, Holmes.")
}
