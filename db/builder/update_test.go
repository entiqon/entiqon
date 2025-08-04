// File: db/builder/update_test.go

package builder

import (
	"fmt"
	"testing"

	core "github.com/entiqon/db/internal/core/errors"
	token2 "github.com/entiqon/db/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type UpdateBuilderTestSuite struct {
	suite.Suite
	qb *UpdateBuilder
}

func (s *UpdateBuilderTestSuite) SetupTest() {
	s.qb = NewUpdate(nil)
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
		Where("id = 42").
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
		Where("deleted", false).
		AndWhere("role", "admin").
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE deleted = ? AND role = ?")
}

// ─────────────────────────────────────────────
// 🧪 OrWhere
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestOrWhere_AppendsOrCondition() {
	sql, _, err := s.qb.
		Table("users").
		Set("active", true).
		Where("email_verified = true").
		OrWhere("status = ?", false).
		Build()

	s.NoError(err)
	s.Contains(sql, "WHERE email_verified = ? OR status = ?")
}

// ─────────────────────────────────────────────
// 🧪 Build
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestBuild_WithAliasedColumn() {
	sql, args, err := s.qb.
		Table("users").
		Set("email AS contact", "watson@example.com").
		Where("id", 1).
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
	s.Contains(err.Error(), "requires a target table")
}

func (s *UpdateBuilderTestSuite) TestBuild_MissingAssignmentsReturnsError() {
	_, _, err := s.qb.
		Table("users").
		Build()

	s.Error(err)
	s.Contains(err.Error(), "must define at least one column assignment")
}

func (s *UpdateBuilderTestSuite) TestBuild_InvalidConditionType_ReturnsError() {
	q := NewUpdate(nil).
		Table("users").
		Set("name", "Watson")

	// Inject invalid condition
	q.Set("x", "y") // keep Set valid
	q.Table("users")
	q.conditions = append(q.conditions, token2.Condition{
		Type: "💣", Key: "broken = true",
	})

	_, _, err := q.Build()
	s.Error(err)
	s.Contains(err.Error(), "UPDATE: unsupported condition type")
}

// ─────────────────────────────────────────────
// 🧪 UseDialect
// ─────────────────────────────────────────────
func (s *UpdateBuilderTestSuite) TestUpdateBuilder_UseDialect_Postgres() {
	sql, args, err := s.qb.
		Set("active", true).
		Table("users").
		Where("email_verified", true).
		OrWhere("email_verified", false).
		UseDialect("postgres").
		Build()

	s.NoError(err)
	s.Equal([]any{true, true, false}, args)
	s.Contains(sql, "WHERE \"email_verified\" = $2 OR \"email_verified\" = $3")
}

func (s *UpdateBuilderTestSuite) TestAddStageError_AppendsToExistingToken() {
	qb := NewUpdate(nil)
	qb.AddStageError(core.StageWhere, fmt.Errorf("first"))
	qb.AddStageError(core.StageWhere, fmt.Errorf("second"))

	errs := qb.Validator.GetErrors()
	s.Len(errs, 2)
	s.ErrorContains(errs[0].Error, "first")
	s.ErrorContains(errs[1].Error, "second")
}

func (s *UpdateBuilderTestSuite) TestAddStageError_CreatesNewTokenGroup() {
	initialLen := len(s.qb.Validator.GetErrors())
	s.qb.AddStageError("OR", fmt.Errorf("or error"))

	s.Len(s.qb.Validator.GetErrors(), initialLen+1)
	s.Equal(core.StageToken("OR"), s.qb.Validator.GetErrors()[len(s.qb.Validator.GetErrors())-1].Stage)
	s.ErrorContains(s.qb.Validator.GetErrors()[len(s.qb.Validator.GetErrors())-1].Error, "or error")
}

func (s *UpdateBuilderTestSuite) TestGetDialect_DefaultsToGeneric() {
	s.qb.BaseBuilder = BaseBuilder{}

	d := s.qb.GetDialect()

	s.NotNil(d)
	s.Equal("generic", d.GetName())
}

func (s *UpdateBuilderTestSuite) TestGetErrors_ReturnsCollectedErrors() {
	s.qb.AddStageError("WHERE", fmt.Errorf("invalid field"))
	errs := s.qb.Validator.GetErrors()

	s.Len(errs, 1)
	s.Equal(core.StageWhere, errs[0].Stage)
	s.ErrorContains(errs[0].Error, "invalid field")
}

func (s *UpdateBuilderTestSuite) TestUseDialect_ShortCircuitsOnEmptyOrSameName() {
	s.qb.UseDialect("generic")

	ptr1 := s.qb.UseDialect("generic")
	s.Equal(ptr1.Dialect.GetName(), s.qb.Dialect.GetName())

	ptr2 := s.qb.UseDialect("")
	s.Equal(ptr2.Dialect.GetName(), s.qb.Dialect.GetName())
}

func (s *UpdateBuilderTestSuite) TestUseDialect_ResolvesNamedDialect() {
	s.qb.UseDialect("postgres")

	d := s.qb.GetDialect()
	s.Equal("postgres", d.GetName())
}

func (s *UpdateBuilderTestSuite) TestBuild_BuildValidations() {
	c := token2.NewCondition(token2.ConditionSimple, "id = ?")

	b := UpdateBuilder{}
	if !c.IsValid() {
		b.AddStageError("WHERE clause", fmt.Errorf("invalid clause"))
	}
	b.Table("users").Set("name", "Watson")
	s.Run("HasDialect", func() {
		b := NewUpdate(nil)
		b.AddStageError("WHERE clause", fmt.Errorf("invalid clause"))
		b.conditions = []token2.Condition{c}
		_, _, err := b.Build()
		s.Error(err)
		s.Equal(true, b.HasDialect())
		s.Equal("generic", b.Dialect.GetName())
	})
	s.Run("HasErrors", func() {
		_, _, err := NewUpdate(nil).Build()

		s.Error(err)
		s.Contains(err.Error(), "must define at least one column assignment")
	})

}

func TestUpdateBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateBuilderTestSuite))
}
