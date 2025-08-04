// File: db/driver/styling/style_test.go

package styling_test

import (
	"testing"

	"github.com/entiqon/entiqon/db/driver"
	"github.com/entiqon/entiqon/db/driver/styling"
	"github.com/stretchr/testify/assert"
)

func TestStyleValidationCoverage(t *testing.T) {
	t.Run("AliasStyle", func(t *testing.T) {
		t.Run("IsValid", func(t *testing.T) {
			assert.True(t, styling.AliasNone.IsValid())
			assert.True(t, styling.AliasWithoutKeyword.IsValid())
			assert.True(t, styling.AliasWithKeyword.IsValid())
			assert.False(t, styling.AliasStyle(99).IsValid())
		})

		t.Run("Format", func(t *testing.T) {
			assert.Equal(t, "users AS u", styling.AliasWithKeyword.Format("users", "u"))
			assert.Equal(t, "users u", styling.AliasWithoutKeyword.Format("users", "u"))
			assert.Equal(t, "users", styling.AliasNone.Format("users", "u"))
			assert.Equal(t, "users", styling.AliasWithKeyword.Format("users", ""))

			// fallback
			assert.Equal(t, "users", styling.AliasStyle(99).Format("users", "u"))
		})

		t.Run("FormatWith", func(t *testing.T) {
			d := driver.NewGenericDialect()
			assert.Equal(t, "users AS u", styling.AliasWithKeyword.FormatWith(d, "users", "u"))
			assert.Equal(t, "users u", styling.AliasWithoutKeyword.FormatWith(d, "users", "u"))
			assert.Equal(t, "users", styling.AliasNone.FormatWith(d, "users", "u"))
			assert.Equal(t, "users", styling.AliasWithKeyword.FormatWith(d, "users", ""))

			// fallback
			assert.Equal(t, "users", styling.AliasStyle(99).FormatWith(d, "users", "u"))
			// nil dialect
			assert.Equal(t, "users", styling.AliasStyle(99).FormatWith(nil, "users", "u"))
		})
	})

	t.Run("QuoteStyle", func(t *testing.T) {
		t.Run("IsValid", func(t *testing.T) {
			assert.False(t, styling.QuoteUnset.IsValid())
			assert.True(t, styling.QuoteNone.IsValid())
			assert.True(t, styling.QuoteDouble.IsValid())
			assert.True(t, styling.QuoteBacktick.IsValid())
			assert.True(t, styling.QuoteBracket.IsValid())
			assert.False(t, styling.QuoteStyle(99).IsValid())
		})

		t.Run("Quote", func(t *testing.T) {
			assert.Equal(t, `"name"`, styling.QuoteDouble.Quote("name"))
			assert.Equal(t, "`name`", styling.QuoteBacktick.Quote("name"))
			assert.Equal(t, "[name]", styling.QuoteBracket.Quote("name"))
			assert.Equal(t, "name", styling.QuoteNone.Quote("name"))

			// fallback
			assert.Equal(t, "name", styling.QuoteStyle(99).Quote("name"))
		})
	})

	t.Run("PlaceholderStyle", func(t *testing.T) {
		t.Run("IsValid", func(t *testing.T) {
			assert.True(t, styling.PlaceholderQuestion.IsValid())
			assert.True(t, styling.PlaceholderDollar.IsValid())
			assert.True(t, styling.PlaceholderNamed.IsValid())
			assert.True(t, styling.PlaceholderAt.IsValid())
			assert.False(t, styling.PlaceholderStyle(99).IsValid())
		})

		t.Run("Format", func(t *testing.T) {
			assert.Equal(t, "?", styling.PlaceholderQuestion.Format(1))
			assert.Equal(t, "$1", styling.PlaceholderDollar.Format(1))
			assert.Equal(t, "?", styling.PlaceholderAt.Format(1))
			assert.Equal(t, ":param", styling.PlaceholderNamed.FormatNamed("param"))
			assert.Equal(t, "@param", styling.PlaceholderAt.FormatNamed("param"))
			assert.Equal(t, "?", styling.PlaceholderQuestion.FormatNamed("param"))
			assert.Equal(t, "?", styling.PlaceholderStyle(99).Format(1)) // fallback
		})

		t.Run("PlaceholderStyle", func(t *testing.T) {
			assert.Equal(t, "?", styling.PlaceholderQuestion.Format(1))
			assert.Equal(t, "$1", styling.PlaceholderDollar.Format(1))
			assert.Equal(t, "?", styling.PlaceholderAt.Format(1)) // fallback for positional
			assert.Equal(t, ":param", styling.PlaceholderNamed.FormatNamed("param"))
			assert.Equal(t, "@param", styling.PlaceholderAt.FormatNamed("param"))
			assert.Equal(t, "?", styling.PlaceholderQuestion.FormatNamed("param")) // fallback
			assert.Equal(t, "?", styling.PlaceholderStyle(99).FormatNamed("id"))   // fallback
		})
	})
}
