// File: db/builder/delete_test.go

package builder

import (
	"strings"
	"testing"

	"github.com/entiqon/db/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type DeleteBuilderTestSuite struct {
	suite.Suite
	qb *DeleteBuilder
}

func (s *DeleteBuilderTestSuite) SetupTest() {
	s.qb = NewDelete(nil)
}

// ─────────────────────────────────────────────
// 🧪 From
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestFrom() {
	s.Run("BasicUsage", func() {
		sql, _, err := s.qb.From("users").Build()
		s.NoError(err)
		s.Contains(sql, "DELETE")
	})
	s.Run("EmptyFrom", func() {
		_, _, err := s.qb.From("").Build()

		s.Error(err)
		s.ErrorContains(err, "builder validation failed")
	})
	s.Run("MissingFrom", func() {
		_, _, err := NewDelete(nil).Build()
		s.Error(err)
		s.ErrorContains(err, "builder validation failed")
	})
}

// ─────────────────────────────────────────────
// 🧪 Where
// ─────────────────────────────────────────────
func (s *DeleteBuilderTestSuite) TestWhere() {
	qb := NewDelete(nil).From("users").Where("id", 100)

	s.Run("BasicUsage", func() {
		sql, args, err := qb.Build()
		s.NoError(err)
		s.Contains(sql, "WHERE id = ?")
		s.Equal([]any{100}, args)
	})

	s.Run("AndWhere", func() {
		sql, args, _ := qb.AndWhere("status = active").Build()

		s.Contains(sql, "WHERE id = ? AND status = ?")
		s.Equal([]any{100, "active"}, args)
	})
	s.Run("OrWhere", func() {
		sql, args, _ := qb.OrWhere("email_confirmed", false).Build()

		s.Contains(sql, "WHERE id = ? AND status = ? OR email_confirmed = ?")
		s.Equal([]any{100, "active", false}, args)
	})
}

// ─────────────────────────────────────────────
// 🧪 Build
// ─────────────────────────────────────────────
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

func (s *DeleteBuilderTestSuite) TestBuild_Validations() {
	s.Run("NoDialect", func() {
		b := NewDelete(nil)
		b.Dialect = nil

		_, _, err := b.Build()
		s.Error(err)
		s.ErrorContains(err, "no dialect set")
	})
	s.Run("HasErrors", func() {
		qb := NewDelete(nil).Where("id = ?")
		_, _, err := qb.Build()

		s.Error(err)
		s.True(qb.Validator.HasErrors(), true)
		s.ErrorContains(err, "builder validation failed")
	})
	s.Run("Basic", func() {
		s.qb.From("users").Where("status =")
		_, _, err := s.qb.Build()

		s.Error(err)
		s.Contains(err.Error(), "builder validation failed")
	})
	s.Run("InvalidAndWhere", func() {
		_, _, err := NewDelete(nil).From("users").AndWhere("", 123, 456).Build()

		s.Error(err)
		s.ErrorContains(err, "builder validation failed")
	})
	s.Run("InvalidOrWhere", func() {
		_, _, err := NewDelete(nil).From("users").OrWhere("", 123).Build()

		s.Error(err)
		s.ErrorContains(err, "builder validation failed")
	})
	s.Run("InvalidConditionType", func() {
		qb := NewDelete(nil).
			From("users")

		_, _, err := qb.Build()
		s.Nil(err)
		qb.conditions = append(qb.conditions, token.Condition{
			Type: "💥", Key: "status = 'active'",
		})

		_, _, err = qb.Build()
		s.Error(err)
		s.Contains(err.Error(), "unsupported condition type")
	})
}

func TestDeleteBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteBuilderTestSuite))
}
