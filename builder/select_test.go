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
	sql, err := s.qb.Select("id", "name").From("users").Build()
	s.NoError(err)
	s.Equal("SELECT id, name FROM users", sql)
	fmt.Printf("📦 Generated SQL Query: %s\n", sql)
}

func (s *SelectQueryBuilderTestSuite) TestWhereAndOrConditions() {
	sql, err := s.qb.
		Select("id").
		From("customers").
		Where("active = true").
		AndWhere("email_verified = true").
		OrWhere("country = 'US'", "country = 'CA'").
		Build()

	expected := "SELECT id FROM customers WHERE active = true AND email_verified = true OR (country = 'US' OR country = 'CA')"
	s.NoError(err)
	s.Equal(expected, sql)
	fmt.Printf("📦 Generated SQL Query: %s\n", sql)
}

func (s *SelectQueryBuilderTestSuite) TestOrderingTakeSkip() {
	sql, err := s.qb.
		Select("name").
		From("employees").
		OrderBy("created_at DESC").
		Take(10).
		Skip(5).
		Build()

	expected := "SELECT name FROM employees ORDER BY created_at DESC LIMIT 10 OFFSET 5"
	s.NoError(err)
	s.Equal(expected, sql)
	fmt.Printf("📦 Generated SQL Query: %s\n", sql)
}

func (s *SelectQueryBuilderTestSuite) TestMissingFromClause() {
	sql, err := s.qb.
		Select("id").
		Build()

	s.Error(err)
	s.Empty(sql)
}

func (s *SelectQueryBuilderTestSuite) TestGroupedAndWhere() {
	sql, err := s.qb.
		From("invoices").
		Where("paid = false").
		AndWhere("amount > 100", "overdue = true").
		Build()

	expected := "SELECT * FROM invoices WHERE paid = false AND (amount > 100 AND overdue = true)"
	s.NoError(err)
	s.Equal(expected, sql)
	fmt.Printf("📦 Generated SQL Query: %s\n", sql)
}

func TestSelectQueryBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(SelectQueryBuilderTestSuite))
	fmt.Println("🕵️ Verified by Watson: All is sound in the SELECT logic, Holmes.")
}
