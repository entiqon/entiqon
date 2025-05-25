// File: internal/build/token/column_test.go

package token_test

import (
	"testing"

	"github.com/ialopezg/entiqon/internal/build/token"
	"github.com/stretchr/testify/assert"
)

//--------------------------------------------------
// Basic Usage
//--------------------------------------------------

func TestNewColumn_ValidCases(t *testing.T) {
	t.Run("BasicName", func(t *testing.T) {
		col := token.NewColumn("id")
		assert.Equal(t, "id", col.Name)
		assert.Empty(t, col.Alias)
		assert.Nil(t, col.Error)
		assert.True(t, col.IsValid())
		assert.Equal(t, "id", col.Raw())
		assert.Equal(t, `Column("id")`, col.String())
	})
}

//--------------------------------------------------
// Usage Variants
//--------------------------------------------------

func TestNewColumn_AliasingVariants(t *testing.T) {
	t.Run("InlineAlias", func(t *testing.T) {
		col := token.NewColumn("user_id AS id")
		assert.Equal(t, "user_id", col.Name)
		assert.Equal(t, "id", col.Alias)
		assert.Nil(t, col.Error)
		assert.True(t, col.IsValid())
		assert.Equal(t, "user_id AS id", col.Raw())
		assert.Equal(t, `Column("user_id" AS "id")`, col.String())
	})

	t.Run("ExplicitAlias", func(t *testing.T) {
		col := token.NewColumn("email", "contact_email")
		assert.Equal(t, "email", col.Name)
		assert.Equal(t, "contact_email", col.Alias)
		assert.Nil(t, col.Error)
		assert.True(t, col.IsValid())
		assert.Equal(t, "email AS contact_email", col.Raw())
		assert.Equal(t, `Column("email" AS "contact_email")`, col.String())
	})

	t.Run("WhitespaceHandling", func(t *testing.T) {
		col := token.NewColumn("  user_id  AS  id  ")
		assert.Equal(t, "user_id", col.Name)
		assert.Equal(t, "id", col.Alias)
		assert.True(t, col.IsValid())
		assert.Equal(t, "user_id AS id", col.Raw())
		assert.Equal(t, `Column("user_id" AS "id")`, col.String())
	})
}

//--------------------------------------------------
// Validations
//--------------------------------------------------

func TestNewColumn_InvalidCases(t *testing.T) {
	t.Run("EmptyExpression", func(t *testing.T) {
		col := token.NewColumn("")
		assert.Error(t, col.Error)
		assert.False(t, col.IsValid())
	})

	t.Run("OnlyAliasKeyword", func(t *testing.T) {
		col := token.NewColumn("AS name")
		assert.Error(t, col.Error)
		assert.False(t, col.IsValid())
	})

	t.Run("MissingNameBeforeAlias", func(t *testing.T) {
		col := token.NewColumn(" AS alias")
		assert.Error(t, col.Error)
		assert.False(t, col.IsValid())
	})

	t.Run("MissingAliasAfterAs", func(t *testing.T) {
		col := token.NewColumn("id AS ")
		assert.Equal(t, "id", col.Name)
		assert.Equal(t, "", col.Alias)
		assert.Nil(t, col.Error)
		assert.True(t, col.IsValid())
	})

	t.Run("OnlyAs", func(t *testing.T) {
		col := token.NewColumn("AS")
		assert.Error(t, col.Error)
		assert.False(t, col.IsValid())
	})
}

//--------------------------------------------------
// Other
//--------------------------------------------------

func TestColumn_InterfaceSatisfaction(t *testing.T) {
	var col token.BaseToken = token.NewColumn("x")
	assert.NotNil(t, col)
	assert.True(t, col.IsValid())
	assert.Equal(t, "x", col.Raw())
	assert.Contains(t, col.String(), "Column")
}
