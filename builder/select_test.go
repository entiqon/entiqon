package builder

import (
	"fmt"
	"testing"

	"github.com/entiqon/entiqon/driver"
	"github.com/entiqon/entiqon/internal/core/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SelectBuilderTestSuite struct {
	suite.Suite
	qb *SelectBuilder
}

func (s *SelectBuilderTestSuite) SetupTest() {
	s.qb = NewSelect(nil)
}

// ─────────────────────────────────────────────
// 🧪 Select
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestSelect() {
	s.Run("Basic", func() {
		sql, _, err := s.qb.
			Select("id", "name").
			From("users").
			Build()

		s.NoError(err)
		s.Contains(sql, "SELECT id, name FROM users")
	})
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
func (s *SelectBuilderTestSuite) TestFrom() {
	s.Run("EmptyTable", func() {
		s.qb.Select("*").From("")
		_, _, err := s.qb.Build()

		s.Error(err)
		s.Contains(err.Error(), "table is empty")
	})

	s.Run("MissingFromClause", func() {
		_, _, err := NewSelect(nil).
			Select("id").
			Build()

		s.Error(err)
		s.ErrorContains(err, "missing source")
	})
}

// ─────────────────────────────────────────────
// 🧪 Where, AndWhere, OrWhere
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestWhereAndOrConditions() {
	sql, _, err := NewSelect(nil).
		Select("id").
		From("customers").
		Where("active", true).
		AndWhere("email_verified", true).
		OrWhere("country", "US").
		OrWhere("country", "CA").
		Build()

	expected := "SELECT id FROM customers WHERE active = ? AND email_verified = ? OR country = ? OR country = ?"
	s.NoError(err)
	s.Equal(expected, sql)

	//s.Run("MultipleInline", func() {
	//	sql, _, err := NewSelect(nil).
	//		Select("id").
	//		From("customers").
	//		Where("active = true AND email_verified = true AND country = US OR country = CA").
	//		Build()
	//
	//	expected := "SELECT id FROM customers WHERE active = ? AND email_verified = ? OR country = ? OR country = ?"
	//	s.NoError(err)
	//	s.Equal(expected, sql)
	//})
}

func (s *SelectBuilderTestSuite) TestGroupedAndWhere() {
	sql, _, err := s.qb.
		From("invoices").
		Where("paid = false").
		AndWhere("amount > 100").
		AndWhere("overdue = true").
		Build()

	expected := "SELECT * FROM invoices WHERE paid = ? AND amount > ? AND overdue = ?"
	s.NoError(err)
	s.Equal(expected, sql)
}

func (s *SelectBuilderTestSuite) TestSelectBuilderMultiParams() {
	sql, params, err := s.qb.
		Select("id", "email").
		From("users").
		Where("status", "active").
		AndWhere("role", "admin").
		AndWhere("created_at > 2024-01-01").
		AndWhere("region = NA").
		OrderBy("last_login DESC").
		Take(50).
		Skip(0).
		Build()

	expected := "SELECT id, email FROM users WHERE status = ? AND role = ? AND created_at > ? AND region = ? ORDER BY last_login DESC LIMIT 50"

	s.NoError(err)
	s.Equal(expected, sql)
	s.Equal([]any{"active", "admin", "2024-01-01", "NA"}, params)
}

func (s *SelectBuilderTestSuite) TestBuild_UnsupportedConditionType() {
	c := token.NewCondition(token.ConditionSimple, "id", 1)
	c.Type = token.ConditionType(rune(999)) // force invalid

	s.qb.From("users").Select("*")
	s.qb.conditions = []token.Condition{c}

	_, _, err := s.qb.Build()
	s.Error(err)
	s.Contains(err.Error(), "unsupported condition type")
}

func (s *SelectBuilderTestSuite) TestOrWhere_InvalidCondition() {
	s.Run("AndWhere_InvalidCondition", func() {
		s.qb.Select("*").From("users").Where("active = ?", true).AndWhere("amount >")
		_, _, err := s.qb.Build()

		s.Error(err)
		s.Contains(err.Error(), "unable to parse condition")
	})
	s.Run("OrWhere_InvalidCondition", func() {
		s.qb.Select("*").From("users").Where("active = ?", true).OrWhere("amount >")
		_, _, err := s.qb.Build()

		s.Error(err)
		s.Contains(err.Error(), "unable to parse condition")
	})
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

func (s *SelectBuilderTestSuite) TestBuild_LimitOffsetWithoutDialect() {
	s.qb.Select("*").From("users").Take(10).Skip(5)

	sql, _, err := s.qb.Build()

	s.NoError(err)
	s.Contains(sql, "LIMIT 10")
	s.Contains(sql, "OFFSET 5")

	s.Run("Limit", func() {
		s.qb.Select("*").From("users").Take(10).Skip(0)
		sql, _, err := s.qb.Build()
		s.NoError(err)
		s.Contains(sql, "LIMIT 10")
	})
	s.Run("Offset", func() {
		s.qb.Select("*").From("users").Take(0).Skip(5)
		sql, _, err := s.qb.Build()
		s.NoError(err)
		s.Contains(sql, "OFFSET 5")
	})
}

// ─────────────────────────────────────────────
// 🧪 UseDialect
// ─────────────────────────────────────────────
func (s *SelectBuilderTestSuite) TestSelectBuilderUseDialectPostgres() {
	sql, args, err := NewSelect(driver.NewPostgresDialect()).
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
func (s *SelectBuilderTestSuite) TestBuild_InvalidConditionType() {
	b := NewSelect(nil).
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
	s.Contains(err.Error(), "unsupported condition type")
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

func (s *SelectBuilderTestSuite) TestBuild_BuildValidations() {
	c := token.NewCondition(token.ConditionSimple, "id = ?")

	s.Run("EmptyTable", func() {
		_, _, err := NewSelect(nil).Build()
		s.Error(err)
		s.ErrorContains(err, "missing source")
	})
	s.Run("HasDialect", func() {
		b := NewSelect(nil)
		if !c.IsValid() {
			b.AddStageError("WHERE clause", fmt.Errorf("invalid clause"))
		}
		_, _, err := b.Build()
		s.Error(err)
		s.Equal("generic", b.Dialect.GetName())
	})
	s.Run("HasErrors", func() {
		_, _, err := NewSelect(nil).Build()
		s.Error(err)
		s.Contains(err.Error(), "builder validation failed")
	})
}

func TestSelectBuilder(t *testing.T) {
	t.Run("Select", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			_, _, err := NewSelect(nil).
				Select("").
				From("users").
				Build()

			assert.Error(t, err)
			assert.ErrorContains(t, err, "empty column expression")
		})

		t.Run("WithTable", func(t *testing.T) {
			sql, _, err := NewSelect(nil).
				From("users AS u").
				AddSelect("id").
				Build()

			assert.NoError(t, err)
			assert.Equal(t, "SELECT u.id FROM users AS u", sql)
		})
	})

	t.Run("AddSelect", func(t *testing.T) {
		sql, _, err := NewSelect(nil).
			AddSelect("id", "created_at").
			From("users").
			Build()

		assert.NoError(t, err)
		assert.NotEmpty(t, sql)

		t.Run("WithoutPriorColumns", func(t *testing.T) {
			qb := SelectBuilder{BaseBuilder: BaseBuilder{Dialect: driver.NewGenericDialect()}}
			sql, _, err := qb.AddSelect("id", "email").
				From("users").
				Build()

			assert.NoError(t, err)
			assert.NotEmpty(t, sql)
		})
	})

	t.Run("From", func(t *testing.T) {
		t.Run("SingleTable", func(t *testing.T) {
			sql, _, err := NewSelect(nil).
				Select("id").
				From("customers").
				Build()

			assert.NoError(t, err)
			assert.Contains(t, sql, "FROM customers")
		})

		t.Run("WithAlias", func(t *testing.T) {
			sql, _, err := NewSelect(nil).
				Select("id", "name").
				From("users", "u").
				Build()

			assert.NoError(t, err)
			assert.Contains(t, sql, "users")
		})
	})

	t.Run("Where", func(t *testing.T) {
		_, _, err := NewSelect(nil).
			Where("status =").
			From("users").
			Build()

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid condition")
	})
}

func TestSelectBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(SelectBuilderTestSuite))
}
