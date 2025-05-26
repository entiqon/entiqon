// File: internal/build/token/column_test.go

package token_test

import (
	"testing"

	"github.com/ialopezg/entiqon/internal/build/token"
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
			assert.Equal(t, "users", col.TableName)
			assert.True(t, col.IsQualified())
			assert.True(t, col.IsValid())
		})

		t.Run("Conflict", func(t *testing.T) {
			col := token.NewColumn("users.id").WithTable("orders")
			assert.Equal(t, "users", col.TableName)
			assert.True(t, col.HasError())
			assert.False(t, col.IsValid())
		})
	})

	t.Run("WithTable", func(t *testing.T) {
		col := token.NewColumn("users.id")
		assert.Equal(t, "users", col.TableName)
		assert.False(t, col.HasError())

		col = col.WithTable("orders")           // should produce error
		assert.Equal(t, "users", col.TableName) // should not change
		assert.True(t, col.HasError())
		assert.Contains(t, col.Error.Error(), "table mismatch")

		t.Run("NoTable", func(t *testing.T) {
			col := token.NewColumn("id", "user_id").WithTable("users")
			assert.Equal(t, "users", col.TableName)
			assert.False(t, col.HasError())
			assert.Equal(t, "users.id AS user_id", col.Raw())
		})
	})

	t.Run("Rendering", func(t *testing.T) {
		col := token.NewColumn("users.id", "uid")
		assert.Equal(t, "users.id AS uid", col.Raw())
		assert.Contains(t, col.String(), `Column("id") [aliased: true, qualified: true, errored: false`)
	})

	t.Run("Raw", func(t *testing.T) {
		col := token.NewColumn("id")
		assert.False(t, col.HasError())
		assert.Equal(t, "id", col.Raw())
	})

	t.Run("String", func(t *testing.T) {
		col := token.NewColumn("users.id").WithTable("orders")
		assert.True(t, col.HasError())
		assert.Contains(t, col.String(), "table mismatch")
	})
}
