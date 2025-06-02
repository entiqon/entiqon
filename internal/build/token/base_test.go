// File: internal/build/token/token_test.go

package token_test

import (
	"fmt"
	"testing"

	"github.com/entiqon/entiqon/internal/build/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseToken(t *testing.T) {
	t.Run("AliasOr", func(t *testing.T) {
		tok := token.BaseToken{Name: "id"}
		assert.Equal(t, "id", tok.AliasOr())

		tok = token.BaseToken{Name: "id", Alias: "uid"}
		assert.Equal(t, "uid", tok.AliasOr())
	})

	t.Run("IsValid", func(t *testing.T) {
		valid := &token.BaseToken{Name: "id"}
		assert.True(t, valid.IsValid())
		invalid := &token.BaseToken{Name: ""}
		assert.False(t, invalid.IsValid())
		errored := &token.BaseToken{Name: "id", Error: assert.AnError}
		assert.False(t, errored.IsValid())
	})

	t.Run("Raw", func(t *testing.T) {
		base := &token.BaseToken{Name: "name"}
		assert.Equal(t, "name", base.Raw())
		raw := token.BaseToken{Name: "name", Alias: "alias"}
		assert.Equal(t, "name AS alias", raw.Raw())
	})

	t.Run("String", func(t *testing.T) {
		singleCol := &token.BaseToken{Name: "id"}
		assert.Equal(t, `Column("id") [aliased: false, errored: false]`, singleCol.String(token.KindColumn))
		aliasedCol := &token.BaseToken{Name: "id", Alias: "user_id"}
		assert.Equal(t, `Column("id") [aliased: true, errored: false]`, aliasedCol.String(token.KindColumn))
		tblAliased := token.BaseToken{Name: "users", Alias: "u"}
		assert.Equal(t, `Table("users") [aliased: true, errored: false]`, tblAliased.String(token.KindTable))

		t.Run("Error", func(t *testing.T) {
			tok := token.BaseToken{
				Name:  "id",
				Alias: "uid",
				Error: fmt.Errorf("conflict"),
			}
			str := tok.String(token.KindColumn)
			assert.Contains(t, str, `errored: true`)
			assert.Contains(t, str, `error: conflict`)
		})
	})

	t.Run("WithError", func(t *testing.T) {
		b := &token.BaseToken{Name: "x"}
		res := b.WithError(fmt.Errorf("invalid"))

		require.Same(t, b, res)
		require.True(t, res.HasError())
		require.EqualError(t, res.Error, "invalid")
	})
}
