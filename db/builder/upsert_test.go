// File: db/builder/upsert_test.go

package builder

import (
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/stretchr/testify/suite"
)

type UpsertBuilderTestSuite struct {
	suite.Suite
}

func TestUpsertBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(UpsertBuilderTestSuite))
}

// ─────────────────────────────────────────────
// 🧪 WithDialect
// ─────────────────────────────────────────────
func (s *UpsertBuilderTestSuite) TestWithDialect_EscapesIdentifiers() {
	q := NewUpsert(driver.NewPostgresDialect()).
		Into("user profile").
		Columns("user id", "email").
		Values(99, "hello@test.dev").
		OnConflict("user id").
		DoUpdateSet(
			Assignment{Column: "email", Expr: "EXCLUDED.email"},
		)

	sql, args, err := q.Build()
	s.Require().NoError(err)
	s.Equal(
		`INSERT INTO "user profile" ("user id", "email") VALUES ($1, $2) ON CONFLICT ("user id") DO UPDATE SET "email" = EXCLUDED.email`,
		sql,
	)
	s.Equal([]any{99, "hello@test.dev"}, args)
}

func (s *UpsertBuilderTestSuite) TestWithoutDialect() {
	q := NewUpsert(nil).
		Into("user profile").
		Columns("user id", "email").
		Values(99, "hello@test.dev").
		OnConflict("user id").
		DoUpdateSet(
			Assignment{Column: "email", Expr: "EXCLUDED.email"},
		)

	sql, args, err := q.Build()
	s.Require().NoError(err)
	s.Equal(
		`INSERT INTO user profile (user id, email) VALUES (?, ?) ON CONFLICT (user id) DO UPDATE SET email = EXCLUDED.email`,
		sql,
	)
	s.Equal([]any{99, "hello@test.dev"}, args)
}

func (s *UpsertBuilderTestSuite) TestReturning_WithoutDialectRawNames() {
	q := NewUpsert(nil).
		Into("emails").
		Columns("id", "value").
		Values(101, "none@entiqon.dev").
		OnConflict("id").
		Returning("id", "value")

	sql, args, err := q.Build()
	s.Require().Error(err)
	s.Empty(sql)
	s.Nil(args)
}

// ─────────────────────────────────────────────
// 🧪 Returning
// ─────────────────────────────────────────────
func (s *UpsertBuilderTestSuite) TestReturning_AppendsReturningClause() {
	q := NewUpsert(driver.NewPostgresDialect()).
		Into("users").
		Columns("id", "email").
		Values(1, "dev@entiqon.dev").
		OnConflict("id").
		DoUpdateSet(
			Assignment{Column: "email", Expr: "EXCLUDED.email"},
		).
		Returning("id", "email")

	sql, args, err := q.Build()
	s.Require().NoError(err)
	s.Equal(
		"INSERT INTO \"users\" (\"id\", \"email\") VALUES ($1, $2) ON CONFLICT (\"id\") DO UPDATE SET \"email\" = EXCLUDED.email RETURNING \"id\", \"email\"",
		sql,
	)
	s.Equal([]any{1, "dev@entiqon.dev"}, args)
}

// ─────────────────────────────────────────────
// 🧪 DoUpdateSet
// ─────────────────────────────────────────────
func (s *UpsertBuilderTestSuite) TestDoUpdateSet_AppendsAssignments() {
	q := NewUpsert(nil).
		DoUpdateSet(
			Assignment{Column: "name", Expr: "EXCLUDED.name"},
			Assignment{Column: "email", Expr: "EXCLUDED.email"},
		)

	sql, args, _ := q.Build()
	s.Empty(sql)
	s.Len(args, 0) // ensure no values are injected yet
	s.Equal([]Assignment{
		{Column: "name", Expr: "EXCLUDED.name"},
		{Column: "email", Expr: "EXCLUDED.email"},
	}, q.DoUpdateSet().updateSet)
}

// ─────────────────────────────────────────────
// 🧪 OnConflict
// ─────────────────────────────────────────────
func (s *UpsertBuilderTestSuite) TestOnConflict_AppendsConflictColumns() {
	q := NewUpsert(driver.NewPostgresDialect()).
		Into("people").
		Columns("id", "email").
		Values(1, "someone@dev.com").
		OnConflict("id", "email")

	sql, args, err := q.Build()
	s.Require().NoError(err)
	s.Equal(
		`INSERT INTO "people" ("id", "email") VALUES ($1, $2) ON CONFLICT ("id", "email") DO NOTHING`,
		sql,
	)
	s.Equal([]any{1, "someone@dev.com"}, args)
}

// ─────────────────────────────────────────────
// 🧪 Build
// ─────────────────────────────────────────────
func (s *UpsertBuilderTestSuite) TestBuild_DoUpdate() {
	q := NewUpsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		OnConflict("id").
		DoUpdateSet(
			Assignment{Column: "name", Expr: "EXCLUDED.name"},
		)

	sql, args, err := q.Build()
	s.Require().NoError(err)
	s.Equal(
		"INSERT INTO users (id, name) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name",
		sql,
	)
	s.Equal([]any{1, "Watson"}, args)
}

func (s *UpsertBuilderTestSuite) TestBuild_DoNothing() {
	q := NewUpsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		OnConflict("id")

	sql, args, err := q.Build()
	s.Require().NoError(err)
	s.Equal(
		"INSERT INTO users (id, name) VALUES (?, ?) ON CONFLICT (id) DO NOTHING",
		sql,
	)
	s.Equal([]any{1, "Watson"}, args)
}

func (s *UpsertBuilderTestSuite) TestBuild_BuildValidations() {
	b := UpsertBuilder{}
	s.Run("EmptyTable", func() {
		_, _, err := NewUpsert(nil).Build()
		s.Error(err)
		s.Contains(err.Error(), "requires a target table")
	})
	s.Run("HasDialect", func() {
		_, _, err := NewUpsert(nil).Into("users").Columns("id").Build()
		s.Error(err)
		s.Equal("generic", b.GetDialect().GetName())
	})
	s.Run("HasErrors", func() {
		_, _, err := NewUpsert(nil).Into("users").Columns("").Build()
		s.Error(err)
		s.Contains(err.Error(), "at least one set of values is required")
	})
	s.Run("Returning", func() {
		_, _, err := NewUpsert(nil).Into("users").Columns("id").Values(1).Returning("id").Build()
		s.Error(err)
		s.Contains(err.Error(), "RETURNING not supported in dialect")
	})
	s.Run("ColumnWithAlias", func() {
		_, _, err := NewUpsert(nil).Into("users").
			Columns("id AS IDENTIFIER").
			Values(1).
			Build()
		s.Error(err)
		s.Contains(err.Error(), "row 1 has 1 values")
	})
}

// helper to normalize sql sources incomplete statements for inspection
func (s *UpsertBuilderTestSuite) normalizeSQL(q *UpsertBuilder) string {
	sql, _, _ := q.Build()
	return sql
}
