// File: db/builder/insert_test.go

package builder_test

import (
	"fmt"
	"testing"

	"github.com/entiqon/entiqon/db/builder"
	driver2 "github.com/entiqon/entiqon/db/driver"
	token2 "github.com/entiqon/entiqon/db/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type InsertBuilderTestSuite struct {
	suite.Suite
	qb *builder.InsertBuilder
}

func (s *InsertBuilderTestSuite) SetupTest() {
	s.qb = builder.NewInsert(nil)
}

// ─────────────────────────────────────────────
// 🧪 Columns()
// ─────────────────────────────────────────────
func (s *InsertBuilderTestSuite) TestBuildInsertOnly_NoColumns() {
	b := builder.NewInsert(nil).Into("users").Values(1, "Watson")
	_, _, err := b.BuildInsertOnly()
	s.ErrorContains(err, "at least one column is required")
}

// ─────────────────────────────────────────────
// 🧪 Values
// ─────────────────────────────────────────────
func (s *InsertBuilderTestSuite) TestBuildInsertOnly_NoValues() {
	b := builder.NewInsert(nil).Into("users").Columns("id", "name")
	_, _, err := b.BuildInsertOnly()
	s.ErrorContains(err, "at least one set of values is required")
}

// ─────────────────────────────────────────────
// 🧪 Returning
// ─────────────────────────────────────────────
func (s *InsertBuilderTestSuite) TestInsertBuilder_WithReturning() {
	q := builder.NewInsert(driver2.NewPostgresDialect()).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		Returning("id", "created_at")

	sql, args, err := q.Build()
	s.NoError(err)
	s.Equal(`INSERT INTO "users" ("id", "name") VALUES ($1, $2) RETURNING "id", "created_at"`, sql)
	s.Equal([]any{1, "Watson"}, args)
}

// ─────────────────────────────────────────────
// 🧪 Dialect Handling
// ─────────────────────────────────────────────
func (s *InsertBuilderTestSuite) TestInsertBuilder_WithDialect_Postgres() {
	sql, _, err := builder.NewInsert(driver2.NewPostgresDialect()).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		Build()

	s.NoError(err)
	s.Equal(`INSERT INTO "users" ("id", "name") VALUES ($1, $2)`, sql)
}

// ─────────────────────────────────────────────
// 🧪 Build()
// ─────────────────────────────────────────────
func (s *InsertBuilderTestSuite) TestBuildInsertOnly_MismatchedRowLength() {
	b := builder.NewInsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1)

	_, _, err := b.BuildInsertOnly()
	s.ErrorContains(err, "row 1 has 1 values, expected 2")
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_BuildErrors() {
	// Missing table
	_, _, err := builder.NewInsert(nil).
		Columns("id").
		Values(1).
		Build()
	s.Error(err)

	// Missing columns
	_, _, err = builder.NewInsert(nil).
		Into("users").
		Values(1).
		Build()
	s.Error(err)

	// Missing values
	_, _, err = builder.NewInsert(nil).
		Into("users").
		Columns("id").
		Build()
	s.Error(err)
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_MismatchedValueCount() {
	_, _, err := builder.NewInsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1).
		Build()
	s.Error(err)
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_MissingFieldError() {
	_, _, err := builder.NewInsert(nil).
		Into("users").
		Columns("id", "name", "email").
		Values(1, "Watson"). // Only 2 values for 3 columns
		Build()

	s.Error(err)
	s.Contains(err.Error(), "row 1 has 2 values")
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_WithAliasedColumn() {
	sql, args, err := builder.NewInsert(nil).
		Into("users").
		Columns("email AS contact").
		Values("watson@example.com").
		Build()

	s.Error(err)
	s.Empty(sql)
	s.Nil(args)
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_Build_ReturningWithoutDialectFails() {
	_, _, err := builder.NewInsert(nil).
		Into("users").
		Columns("id").
		Values(1).
		Returning("id").
		Build()

	s.ErrorContains(err, "builder validation failed")
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_Build_WithDialectNoReturning() {
	sql, args, err := builder.NewInsert(driver2.NewPostgresDialect()).
		Into("users").
		Columns("id").
		Values(1).
		Build()

	s.NoError(err)
	s.Equal(`INSERT INTO "users" ("id") VALUES ($1)`, sql)
	s.Equal([]any{1}, args)
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_Build_ReturningWithGenericDialectFails() {
	sql, args, err := builder.NewInsert(nil).
		Into("users").
		Columns("id").
		Values(1).
		Returning("id").
		Build()

	s.ErrorContains(err, "builder validation failed")
	s.Empty(sql)
	s.Nil(args)
}

// ─────────────────────────────────────────────
// 🧪 BuildInsertOnly()
// ─────────────────────────────────────────────
// BuildInsertOnly with valid input — ensures success case
func (s *InsertBuilderTestSuite) TestBuildInsertOnly_ValidInsert() {
	b := builder.NewInsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson")

	sql, args, err := b.BuildInsertOnly()
	s.NoError(err)
	s.Equal(`INSERT INTO users (id, name) VALUES (?, ?)`, sql)
	s.Equal([]any{1, "Watson"}, args)
}

func (s *InsertBuilderTestSuite) TestBuildInsertOnly_MultiRowSuccess() {
	sql, args, err := builder.NewInsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		Values(2, "Holmes").
		BuildInsertOnly()

	s.NoError(err)
	s.Equal(`INSERT INTO users (id, name) VALUES (?, ?), (?, ?)`, sql)
	s.Equal([]any{1, "Watson", 2, "Holmes"}, args)
}

func (s *InsertBuilderTestSuite) TestBuildInsertOnly_MissingTableFails() {
	_, _, err := builder.NewInsert(nil).
		Columns("id").
		Values(1).
		BuildInsertOnly()

	s.ErrorContains(err, "requires a target table")
}

func (s *InsertBuilderTestSuite) TestBuildInsertOnly_TableWithDialect() {
	sql, args, err := builder.NewInsert(driver2.NewPostgresDialect()).
		Into("users").
		Columns("id").
		Values(1).
		BuildInsertOnly()

	s.NoError(err)
	s.Equal(`INSERT INTO "users" ("id") VALUES ($1)`, sql)
	s.Equal([]any{1}, args)
}

func (s *InsertBuilderTestSuite) TestBuildInsertOnly_ColumnEscapingWithDialect() {
	s.Run("Generic", func() {
		sql, _, err := builder.NewInsert(driver2.NewGenericDialect()).
			Into("users").
			Columns("email").
			Values("x@example.com").
			BuildInsertOnly()

		s.NoError(err)
		s.Contains(sql, "email")
	})
	s.Run("MSSQL", func() {
		sql, _, err := builder.NewInsert(driver2.NewMSSQLDialect()).
			Into("users").
			Columns("email").
			Values("x@example.com").
			BuildInsertOnly()

		s.NoError(err)
		s.Contains(sql, "[email]")
	})
	s.Run("MySQL", func() {
		sql, _, err := builder.NewInsert(driver2.NewMySQLDialect()).
			Into("users").
			Columns("email").
			Values("x@example.com").
			BuildInsertOnly()

		s.NoError(err)
		s.Contains(sql, "`email`")
	})
	s.Run("PostgreSQL", func() {
		sql, _, err := builder.NewInsert(driver2.NewPostgresDialect()).
			Into("users").
			Columns("email").
			Values("x@example.com").
			BuildInsertOnly()

		s.NoError(err)
		s.Contains(sql, `"email"`)
	})
}

func (s *InsertBuilderTestSuite) TestBuild_BuildValidations() {
	c := token2.NewCondition(token2.ConditionSimple, "id = ?")

	s.Run("EmptyTable", func() {
		b := builder.NewInsert(nil)
		_, _, err := b.Build()
		s.Error(err)
		s.Contains(err.Error(), "requires a target table")
	})
	s.Run("HasDialect", func() {
		b := builder.NewInsert(nil).Into("users")
		_, _, err := b.Build()

		if !c.IsValid() {
			b.AddStageError("WHERE", fmt.Errorf("invalid clause"))
		}

		s.Error(err)
		s.Equal("generic", b.GetDialect().GetName())
	})
	s.Run("HasErrors", func() {
		b := builder.NewInsert(nil).Into("users")
		_, _, err := b.Into("").Build()
		s.Error(err)
		s.Contains(err.Error(), "builder validation failed")
	})

}

func TestInsertBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(InsertBuilderTestSuite))
}
