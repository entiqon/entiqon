// File: db/internal/core/token/table_test.go
// Since: v1.5.0

package token_test

import (
	"testing"

	"github.com/entiqon/db/internal/core/token"
	"github.com/stretchr/testify/suite"
)

type TableSuite struct {
	suite.Suite
}

func (s *TableSuite) TestBasicTable() {
	t := token.NewTable("users")
	s.True(t.IsValid())
	s.Equal("users", t.String())
}

func (s *TableSuite) TestTableWithAlias() {
	t := token.NewTableWithAlias("orders", "o")
	s.True(t.IsValid())
	s.Equal("orders o", t.String())
}

func (s *TableSuite) TestEmptyTable() {
	t := token.NewTable("")
	s.False(t.IsValid())
	s.Equal("", t.String())
}

func (s *TableSuite) TestTrimmedInput() {
	t := token.NewTableWithAlias("  logs  ", "  l ")
	s.True(t.IsValid())
	s.Equal("logs l", t.String())
}

func TestTableSuite(t *testing.T) {
	suite.Run(t, new(TableSuite))
}
