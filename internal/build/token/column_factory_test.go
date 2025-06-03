// File: internal/build/common/parser_test.go

package token_test

import (
	"testing"

	"github.com/entiqon/entiqon/internal/build/token"
	"github.com/stretchr/testify/assert"
)

func TestParseColumns(t *testing.T) {

	//--------------------------------------------------
	// Usage
	//--------------------------------------------------

	t.Run("BasicUsage", func(t *testing.T) {
		cols := token.NewColumnsFrom("id")
		assert.Len(t, cols, 1)
		assert.Equal(t, "id", cols[0].Name)
		assert.Nil(t, cols[0].Error)
	})

	t.Run("CommaSeparatedInput", func(t *testing.T) {
		cols := token.NewColumnsFrom("id, name")
		assert.Len(t, cols, 2)
		assert.Equal(t, "id", cols[0].Name)
		assert.Equal(t, "name", cols[1].Name)
	})

	t.Run("InlineAlias", func(t *testing.T) {
		cols := token.NewColumnsFrom("user_id AS uid")
		assert.Len(t, cols, 1)
		assert.Equal(t, "user_id", cols[0].Name)
		assert.Equal(t, "uid", cols[0].Alias)
		assert.Nil(t, cols[0].Error)
	})

	t.Run("ExplicitAlias", func(t *testing.T) {
		cols := token.NewColumnsFrom("email", "contact AS primary_email")
		assert.Len(t, cols, 2)
		assert.Equal(t, "email", cols[0].Name)
		assert.Nil(t, cols[0].Error)
		assert.Equal(t, "contact", cols[1].Name)
		assert.Equal(t, "primary_email", cols[1].Alias)
		assert.Nil(t, cols[1].Error)
	})

	//--------------------------------------------------
	// Validations
	//--------------------------------------------------

	t.Run("EmptyInput", func(t *testing.T) {
		cols := token.NewColumnsFrom("")
		assert.Len(t, cols, 1)
		//assert.Error(t, cols[0].Error)
	})

	t.Run("OnlyWhitespace", func(t *testing.T) {
		cols := token.NewColumnsFrom("   ")
		assert.Len(t, cols, 1)
		//assert.Error(t, cols[0].Error)
	})

	t.Run("OnlyAliasKeyword", func(t *testing.T) {
		cols := token.NewColumnsFrom("AS alias")
		assert.Len(t, cols, 1)
		//assert.Error(t, cols[0].Error)
	})

	t.Run("AliasWithoutName", func(t *testing.T) {
		cols := token.NewColumnsFrom(" AS email")
		assert.Len(t, cols, 1)
		assert.Error(t, cols[0].Error)
	})
}
