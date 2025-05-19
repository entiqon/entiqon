package builder_test

import (
	"testing"

	"github.com/ialopezg/entiqon/builder"
	"github.com/stretchr/testify/suite"
)

type InsertBuilderTestSuite struct {
	suite.Suite
	qb *builder.InsertBuilder
}

func (s *InsertBuilderTestSuite) SetupTest() {
	s.qb = builder.NewInsert()
}

// ─────────────────────────────────────────────
// 🧪 Columns()
// ─────────────────────────────────────────────
func (s *InsertBuilderTestSuite) TestBuildInsertOnly_NoColumns() {
	b := builder.NewInsert().Into("users").Values(1, "Watson")
	_, _, err := b.BuildInsertOnly()
	s.ErrorContains(err, "at least one column is required")
}

// ─────────────────────────────────────────────
// 🧪 Values
// ─────────────────────────────────────────────
func (s *InsertBuilderTestSuite) TestBuildInsertOnly_NoValues() {
	b := builder.NewInsert().Into("users").Columns("id", "name")
	_, _, err := b.BuildInsertOnly()
	s.ErrorContains(err, "at least one set of values is required")
}

// ─────────────────────────────────────────────
// 🧪 Returning
// ─────────────────────────────────────────────
func (s *InsertBuilderTestSuite) TestInsertBuilder_WithReturning() {
	q := s.qb.
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		Returning("id", "created_at").
		UseDialect("postgres")

	sql, args, err := q.Build()
	s.NoError(err)
	s.Equal(`INSERT INTO "users" ("id", "name") VALUES (?, ?) RETURNING "id", "created_at"`, sql)
	s.Equal([]any{1, "Watson"}, args)
}

// ─────────────────────────────────────────────
// 🧪 Dialect Handling
// ─────────────────────────────────────────────
func (s *InsertBuilderTestSuite) TestInsertBuilder_WithDialect_Postgres() {
	sql, _, err := builder.NewInsert().
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		UseDialect("postgres").
		Build()

	s.NoError(err)
	s.Equal(`INSERT INTO "users" ("id", "name") VALUES (?, ?)`, sql)
}

// WithDialect is deprecated; this test ensures backward compatibility.
func (s *InsertBuilderTestSuite) TestInsertBuilder_WithDialect_LegacySupport() {
	sql, args, err := builder.NewInsert().
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		WithDialect("postgres"). // direct injection
		Build()

	s.NoError(err)
	s.Equal(`INSERT INTO "users" ("id", "name") VALUES (?, ?)`, sql)
	s.Equal([]any{1, "Watson"}, args)
}

// ─────────────────────────────────────────────
// 🧪 Build()
// ─────────────────────────────────────────────
func (s *InsertBuilderTestSuite) TestBuildInsertOnly_MismatchedRowLength() {
	b := builder.NewInsert().
		Into("users").
		Columns("id", "name").
		Values(1) // 👈 Only one value instead of two

	_, _, err := b.BuildInsertOnly()
	s.ErrorContains(err, "expected 2 values, got 1")
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_BuildErrors() {
	// Missing table
	_, _, err := builder.NewInsert().
		Columns("id").
		Values(1).
		Build()
	s.Error(err)

	// Missing columns
	_, _, err = builder.NewInsert().
		Into("users").
		Values(1).
		Build()
	s.Error(err)

	// Missing values
	_, _, err = builder.NewInsert().
		Into("users").
		Columns("id").
		Build()
	s.Error(err)
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_MismatchedValueCount() {
	_, _, err := builder.NewInsert().
		Into("users").
		Columns("id", "name").
		Values(1). // only one value for two columns
		Build()
	s.Error(err)
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_MissingFieldError() {
	_, _, err := builder.NewInsert().
		Into("users").
		Columns("id", "name", "email").
		Values(1, "Watson"). // Only 2 values for 3 columns
		Build()

	s.Error(err)
	s.Contains(err.Error(), "row 1 has 2 values")
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_WithAliasedColumn() {
	sql, args, err := builder.NewInsert().
		Into("users").
		Columns("email AS contact").
		Values("watson@example.com").
		Build()

	s.Error(err)
	s.Empty(sql)
	s.Nil(args)
}

// Build with RETURNING but no dialect — should return a safe error
func (s *InsertBuilderTestSuite) TestInsertBuilder_Build_ReturningWithoutDialectFails() {
	_, _, err := builder.NewInsert().
		Into("users").
		Columns("id").
		Values(1).
		Returning("id").
		Build()

	s.ErrorContains(err, "RETURNING is not supported by the active dialect")
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_Build_WithDialectNoReturning() {
	sql, args, err := builder.NewInsert().
		Into("users").
		Columns("id").
		Values(1).
		UseDialect("postgres").
		Build()

	s.NoError(err)
	s.Equal(`INSERT INTO "users" ("id") VALUES (?)`, sql)
	s.Equal([]any{1}, args)
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_Build_ReturningWithGenericDialectFails() {
	sql, args, err := builder.NewInsert().
		Into("users").
		Columns("id").
		Values(1).
		Returning("id").
		UseDialect("generic"). // simulate a dialect that doesn't support RETURNING
		Build()

	s.ErrorContains(err, `RETURNING is not supported by the active dialect: "generic"`)
	s.Empty(sql)
	s.Nil(args)
}

// ─────────────────────────────────────────────
// 🧪 BuildInsertOnly()
// ─────────────────────────────────────────────
// BuildInsertOnly with valid input — ensures success case
func (s *InsertBuilderTestSuite) TestBuildInsertOnly_ValidInsert() {
	b := builder.NewInsert().
		Into("users").
		Columns("id", "name").
		Values(1, "Watson")

	sql, args, err := b.BuildInsertOnly()
	s.NoError(err)
	s.Equal(`INSERT INTO users (id, name) VALUES (?, ?)`, sql)
	s.Equal([]any{1, "Watson"}, args)
}

func (s *InsertBuilderTestSuite) TestBuildInsertOnly_MultiRowSuccess() {
	sql, args, err := builder.NewInsert().
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
	_, _, err := builder.NewInsert().
		Columns("id").
		Values(1).
		BuildInsertOnly()

	s.ErrorContains(err, "table name is required")
}

func (s *InsertBuilderTestSuite) TestBuildInsertOnly_TableWithDialect() {
	sql, args, err := builder.NewInsert().
		Into("users").
		Columns("id").
		Values(1).
		UseDialect("postgres").
		BuildInsertOnly()

	s.NoError(err)
	s.Equal(`INSERT INTO "users" ("id") VALUES (?)`, sql)
	s.Equal([]any{1}, args)
}

func (s *InsertBuilderTestSuite) TestBuildInsertOnly_ColumnEscapingWithDialect() {
	sql, _, err := builder.NewInsert().
		Into("users").
		Columns("user_id", "email").
		Values(1, "x@example.com").
		UseDialect("postgres").
		BuildInsertOnly()

	s.NoError(err)
	s.Contains(sql, `"user_id"`)
	s.Contains(sql, `"email"`)
}

func TestInsertBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(InsertBuilderTestSuite))
}
