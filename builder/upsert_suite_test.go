package builder

import (
	"fmt"
	"testing"

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
	q := NewUpsert().
		WithDialect("postgres").
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
		`INSERT INTO "user profile" ("user id", "email") VALUES (?, ?) ON CONFLICT ("user id") DO UPDATE SET "email" = EXCLUDED.email`,
		sql,
	)
	s.Equal([]any{99, "hello@test.dev"}, args)
	fmt.Printf("📦 WithDialect → SQL: %s | Args: %+v\n", sql, args)
}

func (s *UpsertBuilderTestSuite) TestWithoutDialect() {
	q := NewUpsert().
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
	fmt.Printf("📦 WithDialect → SQL: %s | Args: %+v\n", sql, args)
}

func (s *UpsertBuilderTestSuite) TestReturning_WithoutDialectRawNames() {
	q := NewUpsert().
		Into("emails").
		Columns("id", "value").
		Values(101, "none@entiqon.dev").
		OnConflict("id").
		Returning("id", "value")

	sql, args, err := q.Build()
	s.Require().NoError(err)
	s.Equal(
		"INSERT INTO emails (id, value) VALUES (?, ?) ON CONFLICT (id) DO NOTHING RETURNING id, value",
		sql,
	)
	s.Equal([]any{101, "none@entiqon.dev"}, args)
	fmt.Printf("📦 WithDialect → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 Returning
// ─────────────────────────────────────────────
func (s *UpsertBuilderTestSuite) TestReturning_AppendsReturningClause() {
	q := NewUpsert().
		UseDialect("postgres").
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
		"INSERT INTO \"users\" (\"id\", \"email\") VALUES (?, ?) ON CONFLICT (\"id\") DO UPDATE SET \"email\" = EXCLUDED.email RETURNING \"id\", \"email\"",
		sql,
	)
	s.Equal([]any{1, "dev@entiqon.dev"}, args)
	fmt.Printf("📦 Returning → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 DoUpdateSet
// ─────────────────────────────────────────────
func (s *UpsertBuilderTestSuite) TestDoUpdateSet_AppendsAssignments() {
	q := NewUpsert().
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
	fmt.Printf("📦 DoUpdateSet → SQL: %s | Args: %+v\n", sql, args)
}

// ─────────────────────────────────────────────
// 🧪 OnConflict
// ─────────────────────────────────────────────
func (s *UpsertBuilderTestSuite) TestOnConflict_AppendsConflictColumns() {
	q := NewUpsert().
		UseDialect("postgres").
		Into("people").
		Columns("id", "email").
		Values(1, "someone@dev.com").
		OnConflict("id", "email")

	sql, args, err := q.Build()
	s.Require().NoError(err)
	s.Equal(
		`INSERT INTO "people" ("id", "email") VALUES (?, ?) ON CONFLICT ("id", "email") DO NOTHING`,
		sql,
	)
	s.Equal([]any{1, "someone@dev.com"}, args)
}

// ─────────────────────────────────────────────
// 🧪 Build
// ─────────────────────────────────────────────
func (s *UpsertBuilderTestSuite) TestBuild_DoUpdate() {
	q := NewUpsert().
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
	fmt.Printf("📦 Build → SQL: %s | Args: %+v\n", sql, args)
}

func (s *UpsertBuilderTestSuite) TestBuild_DoNothing() {
	q := NewUpsert().
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
	fmt.Printf("📦 Build → SQL: %s | Args: %+v\n", sql, args)
}

// helper to normalize sql from incomplete statements for inspection
func (s *UpsertBuilderTestSuite) normalizeSQL(q *UpsertBuilder) string {
	sql, _, _ := q.Build()
	return sql
}
