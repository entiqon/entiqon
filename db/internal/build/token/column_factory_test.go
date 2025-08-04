// File: db/internal/build/token/column_factory_test.go

package token_test

import (
	"testing"

	"github.com/entiqon/db/internal/build/token"
	"github.com/stretchr/testify/assert"
)

func TestParseColumns(t *testing.T) {

	//--------------------------------------------------
	// Usage
	//--------------------------------------------------

	t.Run("BasicUsage", func(t *testing.T) {
		cols := token.NewColumnsFrom("id")
		assert.Len(t, cols, 1)
		assert.Equal(t, "id", cols[0].GetName())
		assert.Nil(t, cols[0].GetError())
	})

	t.Run("CommaSeparatedInput", func(t *testing.T) {
		cols := token.NewColumnsFrom("id, name")
		assert.Len(t, cols, 2)
		assert.Equal(t, "id", cols[0].GetName())
		assert.Equal(t, "name", cols[1].GetName())
	})

	t.Run("InlineAlias", func(t *testing.T) {
		cols := token.NewColumnsFrom("user_id AS uid")
		assert.Len(t, cols, 1)
		assert.Equal(t, "user_id", cols[0].GetName())
		assert.Equal(t, "uid", cols[0].GetAlias())
		assert.Nil(t, cols[0].GetError())
	})

	t.Run("ExplicitAlias", func(t *testing.T) {
		cols := token.NewColumnsFrom("email", "contact AS primary_email")
		assert.Len(t, cols, 2)
		assert.Equal(t, "email", cols[0].GetName())
		assert.Nil(t, cols[0].GetError())
		assert.Equal(t, "contact", cols[1].GetName())
		assert.Equal(t, "primary_email", cols[1].GetAlias())
		assert.Nil(t, cols[1].GetError())
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
		assert.Error(t, cols[0].GetError())
	})
}
