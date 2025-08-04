// File: db/driver/base_test.go

package driver_test

import (
	"testing"

	"github.com/entiqon/entiqon/db/driver"
	"github.com/entiqon/entiqon/db/driver/styling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseDialect(t *testing.T) {
	t.Run("BasicUsage", func(t *testing.T) {
		d := &driver.BaseDialect{
			Name:             "test",
			QuoteStyle:       styling.QuoteDouble,
			PlaceholderStyle: styling.PlaceholderDollar,
			EnableReturning:  true,
			EnableUpsert:     true,
		}

		assert.Equal(t, "test", d.GetName())
		assert.Equal(t, `"field"`, d.QuoteIdentifier("field"))
		assert.Equal(t, "$1", d.Placeholder(1))
		assert.Equal(t, true, d.SupportsReturning())
		assert.Equal(t, true, d.SupportsUpsert())
		assert.NoError(t, d.Validate())
	})

	t.Run("Validate", func(t *testing.T) {
		t.Run("EmptyName", func(t *testing.T) {
			d := &driver.BaseDialect{
				QuoteStyle:       styling.QuoteDouble,
				PlaceholderStyle: styling.PlaceholderQuestion,
			}
			err := d.Validate()
			require.Error(t, err)
			require.Contains(t, err.Error(), "dialect is not configured")
		})

		t.Run("PlaceholderUnset", func(t *testing.T) {
			d := &driver.BaseDialect{
				Name:       "test",
				QuoteStyle: styling.QuoteDouble,
			}
			err := d.Validate()
			require.Error(t, err)
			require.Contains(t, err.Error(), "placeholder style is not configured")
		})

		t.Run("QuoteUnset", func(t *testing.T) {
			t.Run("IsValid", func(t *testing.T) {
				d := &driver.BaseDialect{
					Name:             "test",
					QuoteStyle:       styling.QuoteNone,
					PlaceholderStyle: styling.PlaceholderQuestion,
				}
				err := d.Validate()
				require.NoError(t, err)
			})

			t.Run("InvalidPlaceholder", func(t *testing.T) {
				d := &driver.BaseDialect{
					Name:       "test",
					QuoteStyle: styling.QuoteDouble,
				}

				err := d.Validate()
				assert.ErrorContains(t, err, "placeholder style is not configured")
			})

			t.Run("InvalidQuote", func(t *testing.T) {
				d := &driver.BaseDialect{
					Name:             "test",
					PlaceholderStyle: styling.PlaceholderQuestion,
				}

				err := d.Validate()
				assert.ErrorContains(t, err, "quote style is not configured")
			})
		})
	})
}
