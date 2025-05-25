// File: internal/build/token/token_test.go

package token_test

import (
	"testing"

	"github.com/ialopezg/entiqon/internal/build/token"
	"github.com/stretchr/testify/assert"
)

func TestAliasableToken(t *testing.T) {
	t.Run("IsValid", func(t *testing.T) {
		assert.True(t, token.AliasableToken{Name: "id"}.IsValid())
		assert.False(t, token.AliasableToken{Name: ""}.IsValid())
		assert.False(t, token.AliasableToken{Name: "id", Error: assert.AnError}.IsValid())
	})

	t.Run("Raw", func(t *testing.T) {
		assert.Equal(t, "name", token.AliasableToken{Name: "name"}.Raw())
		assert.Equal(t, "name AS alias", token.AliasableToken{Name: "name", Alias: "alias"}.Raw())
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, `Column("id")`, token.AliasableToken{Name: "id"}.String("Column"))
		assert.Equal(t, `Column("id" AS "user_id")`, token.AliasableToken{Name: "id", Alias: "user_id"}.String("Column"))
		assert.Equal(t, `Table("users" AS "u")`, token.AliasableToken{Name: "users", Alias: "u"}.String("Table"))
	})
}
