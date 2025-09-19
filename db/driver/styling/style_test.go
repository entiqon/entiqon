// File: db/driver/styling/style_test.go

package styling_test

import (
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/driver/styling"
)

func TestStyleValidationCoverage(t *testing.T) {
	t.Run("AliasStyle", func(t *testing.T) {
		t.Run("IsValid", func(t *testing.T) {
			if !styling.AliasNone.IsValid() {
				t.Errorf("expected AliasNone to be valid")
			}
			if !styling.AliasWithoutKeyword.IsValid() {
				t.Errorf("expected AliasWithoutKeyword to be valid")
			}
			if !styling.AliasWithKeyword.IsValid() {
				t.Errorf("expected AliasWithKeyword to be valid")
			}
			if styling.AliasStyle(99).IsValid() {
				t.Errorf("expected AliasStyle(99) to be invalid")
			}
		})

		t.Run("Format", func(t *testing.T) {
			if got := styling.AliasWithKeyword.Format("users", "u"); got != "users AS u" {
				t.Errorf("expected %q, got %q", "users AS u", got)
			}
			if got := styling.AliasWithoutKeyword.Format("users", "u"); got != "users u" {
				t.Errorf("expected %q, got %q", "users u", got)
			}
			if got := styling.AliasNone.Format("users", "u"); got != "users" {
				t.Errorf("expected %q, got %q", "users", got)
			}
			if got := styling.AliasWithKeyword.Format("users", ""); got != "users" {
				t.Errorf("expected %q, got %q", "users", got)
			}

			// fallback
			if got := styling.AliasStyle(99).Format("users", "u"); got != "users" {
				t.Errorf("expected %q, got %q", "users", got)
			}
		})

		t.Run("FormatWith", func(t *testing.T) {
			d := driver.NewGenericDialect()

			if got := styling.AliasWithKeyword.FormatWith(d, "users", "u"); got != "users AS u" {
				t.Errorf("expected %q, got %q", "users AS u", got)
			}
			if got := styling.AliasWithoutKeyword.FormatWith(d, "users", "u"); got != "users u" {
				t.Errorf("expected %q, got %q", "users u", got)
			}
			if got := styling.AliasNone.FormatWith(d, "users", "u"); got != "users" {
				t.Errorf("expected %q, got %q", "users", got)
			}
			if got := styling.AliasWithKeyword.FormatWith(d, "users", ""); got != "users" {
				t.Errorf("expected %q, got %q", "users", got)
			}

			// fallback
			if got := styling.AliasStyle(99).FormatWith(d, "users", "u"); got != "users" {
				t.Errorf("expected %q, got %q", "users", got)
			}
			// nil dialect
			if got := styling.AliasStyle(99).FormatWith(nil, "users", "u"); got != "users" {
				t.Errorf("expected %q, got %q", "users", got)
			}
		})
	})

	t.Run("QuoteStyle", func(t *testing.T) {
		t.Run("IsValid", func(t *testing.T) {
			if styling.QuoteUnset.IsValid() {
				t.Errorf("expected QuoteUnset to be invalid")
			}
			if !styling.QuoteNone.IsValid() {
				t.Errorf("expected QuoteNone to be valid")
			}
			if !styling.QuoteDouble.IsValid() {
				t.Errorf("expected QuoteDouble to be valid")
			}
			if !styling.QuoteBacktick.IsValid() {
				t.Errorf("expected QuoteBacktick to be valid")
			}
			if !styling.QuoteBracket.IsValid() {
				t.Errorf("expected QuoteBracket to be valid")
			}
			if styling.QuoteStyle(99).IsValid() {
				t.Errorf("expected QuoteStyle(99) to be invalid")
			}
		})

		t.Run("Quote", func(t *testing.T) {
			if got := styling.QuoteDouble.Quote("name"); got != `"name"` {
				t.Errorf("expected %q, got %q", `"name"`, got)
			}
			if got := styling.QuoteBacktick.Quote("name"); got != "`name`" {
				t.Errorf("expected %q, got %q", "`name`", got)
			}
			if got := styling.QuoteBracket.Quote("name"); got != "[name]" {
				t.Errorf("expected %q, got %q", "[name]", got)
			}
			if got := styling.QuoteNone.Quote("name"); got != "name" {
				t.Errorf("expected %q, got %q", "name", got)
			}

			// fallback
			if got := styling.QuoteStyle(99).Quote("name"); got != "name" {
				t.Errorf("expected %q, got %q", "name", got)
			}
		})
	})

	t.Run("PlaceholderStyle", func(t *testing.T) {
		t.Run("IsValid", func(t *testing.T) {
			if !styling.PlaceholderQuestion.IsValid() {
				t.Errorf("expected PlaceholderQuestion to be valid")
			}
			if !styling.PlaceholderDollar.IsValid() {
				t.Errorf("expected PlaceholderDollar to be valid")
			}
			if !styling.PlaceholderNamed.IsValid() {
				t.Errorf("expected PlaceholderNamed to be valid")
			}
			if !styling.PlaceholderAt.IsValid() {
				t.Errorf("expected PlaceholderAt to be valid")
			}
			if styling.PlaceholderStyle(99).IsValid() {
				t.Errorf("expected PlaceholderStyle(99) to be invalid")
			}
		})

		t.Run("Format", func(t *testing.T) {
			if got := styling.PlaceholderQuestion.Format(1); got != "?" {
				t.Errorf("expected %q, got %q", "?", got)
			}
			if got := styling.PlaceholderDollar.Format(1); got != "$1" {
				t.Errorf("expected %q, got %q", "$1", got)
			}
			if got := styling.PlaceholderAt.Format(1); got != "?" {
				t.Errorf("expected %q, got %q", "?", got)
			}
			if got := styling.PlaceholderNamed.FormatNamed("param"); got != ":param" {
				t.Errorf("expected %q, got %q", ":param", got)
			}
			if got := styling.PlaceholderAt.FormatNamed("param"); got != "@param" {
				t.Errorf("expected %q, got %q", "@param", got)
			}
			if got := styling.PlaceholderQuestion.FormatNamed("param"); got != "?" {
				t.Errorf("expected %q, got %q", "?", got)
			}
			if got := styling.PlaceholderStyle(99).Format(1); got != "?" {
				t.Errorf("expected %q, got %q", "?", got)
			}
		})

		t.Run("PlaceholderStyle", func(t *testing.T) {
			if got := styling.PlaceholderQuestion.Format(1); got != "?" {
				t.Errorf("expected %q, got %q", "?", got)
			}
			if got := styling.PlaceholderDollar.Format(1); got != "$1" {
				t.Errorf("expected %q, got %q", "$1", got)
			}
			if got := styling.PlaceholderAt.Format(1); got != "?" {
				t.Errorf("expected %q, got %q", "?", got)
			}
			if got := styling.PlaceholderNamed.FormatNamed("param"); got != ":param" {
				t.Errorf("expected %q, got %q", ":param", got)
			}
			if got := styling.PlaceholderAt.FormatNamed("param"); got != "@param" {
				t.Errorf("expected %q, got %q", "@param", got)
			}
			if got := styling.PlaceholderQuestion.FormatNamed("param"); got != "?" {
				t.Errorf("expected %q, got %q", "?", got)
			}
			if got := styling.PlaceholderStyle(99).FormatNamed("id"); got != "?" {
				t.Errorf("expected %q, got %q", "?", got)
			}
		})
	})
}
