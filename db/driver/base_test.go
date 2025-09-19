package driver_test

import (
	"strings"
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/driver/styling"
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

		if got := d.GetName(); got != "test" {
			t.Errorf("expected %q, got %q", "test", got)
		}
		if got := d.QuoteIdentifier("field"); got != `"field"` {
			t.Errorf("expected %q, got %q", `"field"`, got)
		}
		if got := d.Placeholder(1); got != "$1" {
			t.Errorf("expected %q, got %q", "$1", got)
		}
		if !d.SupportsReturning() {
			t.Errorf("expected SupportsReturning() to be true")
		}
		if !d.SupportsUpsert() {
			t.Errorf("expected SupportsUpsert() to be true")
		}
		if err := d.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Validate", func(t *testing.T) {
		t.Run("EmptyName", func(t *testing.T) {
			d := &driver.BaseDialect{
				QuoteStyle:       styling.QuoteDouble,
				PlaceholderStyle: styling.PlaceholderQuestion,
			}
			err := d.Validate()
			if err == nil {
				t.Errorf("expected error for empty name, got nil")
			} else if !strings.Contains(err.Error(), "dialect is not configured") {
				t.Errorf("expected error containing %q, got %q", "dialect is not configured", err.Error())
			}
		})

		t.Run("PlaceholderUnset", func(t *testing.T) {
			d := &driver.BaseDialect{
				Name:       "test",
				QuoteStyle: styling.QuoteDouble,
			}
			err := d.Validate()
			if err == nil {
				t.Errorf("expected error for unset placeholder, got nil")
			} else if !strings.Contains(err.Error(), "placeholder style is not configured") {
				t.Errorf("expected error containing %q, got %q", "placeholder style is not configured", err.Error())
			}
		})

		t.Run("QuoteUnset", func(t *testing.T) {
			t.Run("IsValid", func(t *testing.T) {
				d := &driver.BaseDialect{
					Name:             "test",
					QuoteStyle:       styling.QuoteNone,
					PlaceholderStyle: styling.PlaceholderQuestion,
				}
				if err := d.Validate(); err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			})

			t.Run("InvalidPlaceholder", func(t *testing.T) {
				d := &driver.BaseDialect{
					Name:       "test",
					QuoteStyle: styling.QuoteDouble,
				}
				err := d.Validate()
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if !strings.Contains(err.Error(), "placeholder style is not configured") {
					t.Errorf("expected error containing %q, got %q", "placeholder style is not configured", err.Error())
				}
			})

			t.Run("InvalidQuote", func(t *testing.T) {
				d := &driver.BaseDialect{
					Name:             "test",
					PlaceholderStyle: styling.PlaceholderQuestion,
				}
				err := d.Validate()
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if !strings.Contains(err.Error(), "quote style is not configured") {
					t.Errorf("expected error containing %q, got %q", "quote style is not configured", err.Error())
				}
			})
		})
	})
}
