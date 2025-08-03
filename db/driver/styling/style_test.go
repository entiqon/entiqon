// File: db/driver/styling/style_test.go

package styling_test

import (
	"testing"

	"github.com/entiqon/db/v2/driver"
	"github.com/entiqon/db/v2/driver/styling"
)

func TestStyleValidationCoverage(t *testing.T) {
	t.Run("AliasStyle", func(t *testing.T) {
		t.Run("IsValid", func(t *testing.T) {
			if !styling.AliasNone.IsValid() {
				t.Errorf("AliasNone should be valid")
			}
			if !styling.AliasWithoutKeyword.IsValid() {
				t.Errorf("AliasWithoutKeyword should be valid")
			}
			if !styling.AliasWithKeyword.IsValid() {
				t.Errorf("AliasWithKeyword should be valid")
			}
			if styling.AliasStyle(99).IsValid() {
				t.Errorf("AliasStyle(99) should be invalid")
			}
		})

		t.Run("Format", func(t *testing.T) {
			if got, want := styling.AliasWithKeyword.Format("users", "u"), "users AS u"; got != want {
				t.Errorf("AliasWithKeyword.Format() = %q; want %q", got, want)
			}
			if got, want := styling.AliasWithoutKeyword.Format("users", "u"), "users u"; got != want {
				t.Errorf("AliasWithoutKeyword.Format() = %q; want %q", got, want)
			}
			if got, want := styling.AliasNone.Format("users", "u"), "users"; got != want {
				t.Errorf("AliasNone.Format() = %q; want %q", got, want)
			}
			if got, want := styling.AliasWithKeyword.Format("users", ""), "users"; got != want {
				t.Errorf("AliasWithKeyword.Format() = %q; want %q", got, want)
			}

			// fallback
			if got, want := styling.AliasStyle(99).Format("users", "u"), "users"; got != want {
				t.Errorf("AliasStyle(99).Format() fallback = %q; want %q", got, want)
			}
		})

		t.Run("FormatWith", func(t *testing.T) {
			d := driver.NewGenericDialect()
			if got, want := styling.AliasWithKeyword.FormatWith(d, "users", "u"), "users AS u"; got != want {
				t.Errorf("AliasWithKeyword.FormatWith() = %q; want %q", got, want)
			}
			if got, want := styling.AliasWithoutKeyword.FormatWith(d, "users", "u"), "users u"; got != want {
				t.Errorf("AliasWithoutKeyword.FormatWith() = %q; want %q", got, want)
			}
			if got, want := styling.AliasNone.FormatWith(d, "users", "u"), "users"; got != want {
				t.Errorf("AliasNone.FormatWith() = %q; want %q", got, want)
			}
			if got, want := styling.AliasWithKeyword.FormatWith(d, "users", ""), "users"; got != want {
				t.Errorf("AliasWithKeyword.FormatWith() = %q; want %q", got, want)
			}

			// fallback
			if got, want := styling.AliasStyle(99).FormatWith(d, "users", "u"), "users"; got != want {
				t.Errorf("AliasStyle(99).FormatWith() fallback = %q; want %q", got, want)
			}
			// nil dialect
			if got, want := styling.AliasStyle(99).FormatWith(nil, "users", "u"), "users"; got != want {
				t.Errorf("AliasStyle(99).FormatWith(nil) = %q; want %q", got, want)
			}
		})
	})

	t.Run("QuoteStyle", func(t *testing.T) {
		t.Run("IsValid", func(t *testing.T) {
			if styling.QuoteUnset.IsValid() {
				t.Errorf("QuoteUnset should be invalid")
			}
			if !styling.QuoteNone.IsValid() {
				t.Errorf("QuoteNone should be valid")
			}
			if !styling.QuoteDouble.IsValid() {
				t.Errorf("QuoteDouble should be valid")
			}
			if !styling.QuoteBacktick.IsValid() {
				t.Errorf("QuoteBacktick should be valid")
			}
			if !styling.QuoteBracket.IsValid() {
				t.Errorf("QuoteBracket should be valid")
			}
			if styling.QuoteStyle(99).IsValid() {
				t.Errorf("QuoteStyle(99) should be invalid")
			}
		})

		t.Run("Quote", func(t *testing.T) {
			if got, want := styling.QuoteDouble.Quote("name"), `"name"`; got != want {
				t.Errorf("QuoteDouble.Quote() = %q; want %q", got, want)
			}
			if got, want := styling.QuoteBacktick.Quote("name"), "`name`"; got != want {
				t.Errorf("QuoteBacktick.Quote() = %q; want %q", got, want)
			}
			if got, want := styling.QuoteBracket.Quote("name"), "[name]"; got != want {
				t.Errorf("QuoteBracket.Quote() = %q; want %q", got, want)
			}
			if got, want := styling.QuoteNone.Quote("name"), "name"; got != want {
				t.Errorf("QuoteNone.Quote() = %q; want %q", got, want)
			}

			// fallback
			if got, want := styling.QuoteStyle(99).Quote("name"), "name"; got != want {
				t.Errorf("QuoteStyle(99).Quote() fallback = %q; want %q", got, want)
			}
		})
	})

	t.Run("PlaceholderStyle", func(t *testing.T) {
		t.Run("IsValid", func(t *testing.T) {
			if !styling.PlaceholderQuestion.IsValid() {
				t.Errorf("PlaceholderQuestion should be valid")
			}
			if !styling.PlaceholderDollar.IsValid() {
				t.Errorf("PlaceholderDollar should be valid")
			}
			if !styling.PlaceholderNamed.IsValid() {
				t.Errorf("PlaceholderNamed should be valid")
			}
			if !styling.PlaceholderAt.IsValid() {
				t.Errorf("PlaceholderAt should be valid")
			}
			if styling.PlaceholderStyle(99).IsValid() {
				t.Errorf("PlaceholderStyle(99) should be invalid")
			}
		})

		t.Run("Format", func(t *testing.T) {
			if got, want := styling.PlaceholderQuestion.Format(1), "?"; got != want {
				t.Errorf("PlaceholderQuestion.Format() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderDollar.Format(1), "$1"; got != want {
				t.Errorf("PlaceholderDollar.Format() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderAt.Format(1), "?"; got != want {
				t.Errorf("PlaceholderAt.Format() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderNamed.FormatNamed("param"), ":param"; got != want {
				t.Errorf("PlaceholderNamed.FormatNamed() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderAt.FormatNamed("param"), "@param"; got != want {
				t.Errorf("PlaceholderAt.FormatNamed() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderQuestion.FormatNamed("param"), "?"; got != want {
				t.Errorf("PlaceholderQuestion.FormatNamed() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderStyle(99).Format(1), "?"; got != want {
				t.Errorf("PlaceholderStyle(99).Format() fallback = %q; want %q", got, want)
			}
		})

		t.Run("PlaceholderStyle", func(t *testing.T) {
			if got, want := styling.PlaceholderQuestion.Format(1), "?"; got != want {
				t.Errorf("PlaceholderQuestion.Format() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderDollar.Format(1), "$1"; got != want {
				t.Errorf("PlaceholderDollar.Format() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderAt.Format(1), "?"; got != want {
				t.Errorf("PlaceholderAt.Format() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderNamed.FormatNamed("param"), ":param"; got != want {
				t.Errorf("PlaceholderNamed.FormatNamed() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderAt.FormatNamed("param"), "@param"; got != want {
				t.Errorf("PlaceholderAt.FormatNamed() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderQuestion.FormatNamed("param"), "?"; got != want {
				t.Errorf("PlaceholderQuestion.FormatNamed() = %q; want %q", got, want)
			}
			if got, want := styling.PlaceholderStyle(99).FormatNamed("id"), "?"; got != want {
				t.Errorf("PlaceholderStyle(99).FormatNamed() fallback = %q; want %q", got, want)
			}
		})
	})
}
