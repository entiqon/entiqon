// File: internal/build/token/column_test.go

package token_test

import (
	"fmt"
	"testing"

	"github.com/entiqon/entiqon/internal/build/token"
	"github.com/stretchr/testify/assert"
)

func TestColumn(t *testing.T) {
	t.Run("Usage", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			col := token.NewColumn("id")
			assert.Equal(t, "id", col.Name)
			assert.False(t, col.IsAliased())
			assert.False(t, col.IsQualified())
			assert.True(t, col.IsValid())
		})

		t.Run("Inline", func(t *testing.T) {
			assert.Equal(t, "uid", token.NewColumn("id AS uid").Alias)
			assert.Equal(t, "uid", token.NewColumn("id uid").Alias)
			assert.Equal(t, "uid", token.NewColumn("id, uid").Alias) // TODO: check
		})
	})

	t.Run("AliasOverriding", func(t *testing.T) {
		t.Run("Match", func(t *testing.T) {
			col := token.NewColumn("id AS uid", "uid")
			assert.True(t, col.IsAliased())
			assert.True(t, col.IsValid())
		})

		t.Run("Mismatch", func(t *testing.T) {
			col := token.NewColumn("id AS internal", "external")
			assert.True(t, col.IsAliased())
			assert.True(t, col.HasError())
			assert.False(t, col.IsValid())
		})
	})

	t.Run("Table", func(t *testing.T) {
		t.Run("QualifiedParsing", func(t *testing.T) {
			col := token.NewColumn("users.id")
			assert.Equal(t, "users", col.Table.Name)
			assert.True(t, col.IsQualified())
			assert.True(t, col.IsValid())
		})

		t.Run("Conflict", func(t *testing.T) {
			col := token.NewColumn("users.id").WithTable(token.NewTable("orders"))
			assert.Equal(t, "users", col.Table.Name)
			assert.True(t, col.HasError())
			assert.False(t, col.IsValid())
		})
	})

	t.Run("WithTable", func(t *testing.T) {
		col := token.NewColumn("users.id")
		assert.Equal(t, "users", col.Table.Name)
		assert.False(t, col.HasError())

		col = col.WithTable(token.NewTable("orders")) // should produce error
		assert.Equal(t, "users", col.Table.Name)      // should not change
		assert.True(t, col.HasError())
		assert.Contains(t, col.Error.Error(), "table mismatch")

		t.Run("NoTable", func(t *testing.T) {
			col := token.NewColumn("id", "user_id").WithTable(token.NewTable("users"))
			assert.Equal(t, "users", col.Table.Name)
			assert.False(t, col.HasError())
			assert.Equal(t, "users.user_id", col.Raw())
		})
	})

	t.Run("Rendering", func(t *testing.T) {
		col := token.NewColumn("users.id", "uid")
		assert.Equal(t, "users.uid", col.Raw())
		assert.Contains(t, col.String(), `Column("id") [aliased: true, qualified: true, errored: false`)
	})

	t.Run("Raw", func(t *testing.T) {
		col := token.NewColumn("id")
		assert.False(t, col.HasError())
		assert.Equal(t, "id", col.Raw())
	})

	t.Run("String", func(t *testing.T) {
		col := token.NewColumn("users.id").WithTable(token.NewTable("orders"))
		assert.True(t, col.HasError())
		assert.Contains(t, col.String(), "table mismatch")
	})

	t.Run("ColumnResolution", func(t *testing.T) {
		cases := []struct {
			columnExpr string
			tableExpr  string
			expectErr  bool
			expectName string
		}{
			// Unqualified + alias → should resolve with alias
			{"id", "users AS u", false, "u.id"},

			// Qualified and table match → trusted
			{"users.id", "users", false, "users.id"},

			// Qualified and alias match → trusted
			{"u.id", "users AS u", false, "u.id"},

			// Qualified but table mismatch → error
			{"orders.id", "users AS u", true, "orders.id"},

			// Qualified but no table → error
			{"users.id", "", false, "users.id"},
		}

		for _, tc := range cases {
			col := token.NewColumn(tc.columnExpr)
			if tc.tableExpr != "" {
				col.WithTable(token.NewTable(tc.tableExpr))
			}

			if tc.expectErr && !col.HasError() {
				t.Errorf("expected error for %q with table %q, but got none", tc.columnExpr, tc.tableExpr)
			} else if !tc.expectErr && col.HasError() {
				t.Errorf("unexpected error for %q with table %q: %v", tc.columnExpr, tc.tableExpr, col.Error)
			}

			if tc.expectName != "" {
				got := col.Raw()
				if got != tc.expectName {
					t.Errorf("expected resolved name %q, got %q", tc.expectName, got)
				}
			}
		}
	})
}

func TestNewErroredColumn(t *testing.T) {
	col := token.NewErroredColumn(fmt.Errorf("errored column"))
	assert.ErrorContains(t, col.Error, "errored column")
}
