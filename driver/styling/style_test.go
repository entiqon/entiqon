// File: driver/styling/style_test.go

package styling_test

import (
	"testing"

	"github.com/ialopezg/entiqon/driver"
	"github.com/ialopezg/entiqon/driver/styling"
	"github.com/stretchr/testify/assert"
)

func TestStyling(t *testing.T) {
	t.Run("AliasStyle", func(t *testing.T) {
		t.Run("Format", func(t *testing.T) {
			assert.Equal(t, "users AS u", styling.AliasWithKeyword.Format("users", "u"))
			assert.Equal(t, "users u", styling.AliasWithoutKeyword.Format("users", "u"))
			assert.Equal(t, "users", styling.AliasNone.Format("users", "u"))
			assert.Equal(t, "users", styling.AliasWithKeyword.Format("users", ""))
		})
		t.Run("FormatWith", func(t *testing.T) {
			type dialectCase struct {
				name   string
				driver driver.Dialect
				expect string
			}
			base, alias := "users", "u"

			cases := []dialectCase{
				{"postgres", driver.NewPostgresDialect(), `"users" AS "u"`},
				{"mysql", driver.NewMySQLDialect(), "`users` AS `u`"},
				{"mssql", driver.NewMSSQLDialect(), "[users] AS [u]"},
				//{"sqlite", driver.NewSQLiteDialect(), `"users" AS "u"`},
				//{"oracle", driver.NewOracleDialect(), `"users" "u"`},
				//{"db2", driver.NewDB2Dialect(), `"users" AS "u"`},
				//{"firebird", driver.NewFirebirdDialect(), `"users" AS "u"`},
				//{"informix", driver.NewInformixDialect(), `"users" AS "u"`},
				{"generic", driver.NewGenericDialect(), "users AS u"},
			}

			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					out := styling.AliasWithKeyword.FormatWith(c.driver, base, alias)
					assert.Equal(t, c.expect, out)
				})
			}
		})
	})

	t.Run("QuoteStyle", func(t *testing.T) {
		assert.Equal(t, `"name"`, styling.QuoteDouble.Quote("name"))
		assert.Equal(t, "`name`", styling.QuoteBacktick.Quote("name"))
		assert.Equal(t, "[name]", styling.QuoteBracket.Quote("name"))
		assert.Equal(t, "name", styling.QuoteNone.Quote("name"))
	})

	t.Run("PlaceholderStyle", func(t *testing.T) {
		assert.Equal(t, "?", styling.PlaceholderQuestion.Format(1))
		assert.Equal(t, "$1", styling.PlaceholderDollar.Format(1))
		assert.Equal(t, "?", styling.PlaceholderAt.Format(1))
		assert.Equal(t, ":param", styling.PlaceholderNamed.FormatNamed("param"))
		assert.Equal(t, "@param", styling.PlaceholderAt.FormatNamed("param"))
		assert.Equal(t, "?", styling.PlaceholderQuestion.FormatNamed("param"))
	})
}
