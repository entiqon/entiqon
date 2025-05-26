// File: internal/build/token/token_test.go

package token_test

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/internal/build/token"
	"github.com/stretchr/testify/assert"
)

func TestBaseToken(t *testing.T) {
	t.Run("IsValid", func(t *testing.T) {
		assert.True(t, token.BaseToken{Name: "id"}.IsValid())
		assert.False(t, token.BaseToken{Name: ""}.IsValid())
		assert.False(t, token.BaseToken{Name: "id", Error: assert.AnError}.IsValid())
	})

	t.Run("Raw", func(t *testing.T) {
		assert.Equal(t, "name", token.BaseToken{Name: "name"}.Raw())
		assert.Equal(t, "name AS alias", token.BaseToken{Name: "name", Alias: "alias"}.Raw())
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, `Column("id") [aliased: false, errored: false]`, token.BaseToken{Name: "id"}.String(token.KindColumn))
		assert.Equal(t, `Column("id") [aliased: true, errored: false]`, token.BaseToken{Name: "id", Alias: "user_id"}.String(token.KindColumn))
		assert.Equal(t, `Table("users") [aliased: true, errored: false]`, token.BaseToken{Name: "users", Alias: "u"}.String(token.KindTable))

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

	t.Run("AliasOr", func(t *testing.T) {
		tok := token.BaseToken{Name: "id"}
		assert.Equal(t, "id", tok.AliasOr())

		tok = token.BaseToken{Name: "id", Alias: "uid"}
		assert.Equal(t, "uid", tok.AliasOr())
	})
}
