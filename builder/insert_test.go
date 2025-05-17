package builder_test

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/builder"
	"github.com/ialopezg/entiqon/internal/core/dialect"
	"github.com/stretchr/testify/suite"
)

type InsertBuilderTestSuite struct {
	suite.Suite
	qb *builder.InsertBuilder
}

func (s *InsertBuilderTestSuite) SetupTest() {
	s.qb = builder.NewInsert()
}

func (s *InsertBuilderTestSuite) TestInsertBuilder_WithReturning() {
	q := builder.NewInsert().
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		Returning("id", "created_at")

	sql, args, err := q.Build()
	s.NoError(err)
	s.Equal("INSERT INTO users (id, name) VALUES (?, ?) RETURNING id, created_at", sql)
	s.Equal([]any{1, "Watson"}, args)
	fmt.Printf("📦 Generated SQL Query: %s with values %+v\n", sql, args)
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

func (s *InsertBuilderTestSuite) TestInsertBuilder_WithDialect_Postgres() {
	sql, _, err := builder.NewInsert().
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		WithDialect(&dialect.PostgresEngine{}).
		Build()

	s.NoError(err)
	s.Equal(`INSERT INTO users ("id", "name") VALUES (?, ?)`, sql)
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

	s.NoError(err)
	s.Contains(sql, "email AS contact")
	s.Equal([]any{"watson@example.com"}, args)
}

func TestInsertBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(InsertBuilderTestSuite))
}
