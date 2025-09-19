package render_test

import (
	"fmt"
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/internal/build/render"
	"github.com/entiqon/db/internal/build/token"
)

func TestRenderTable(t *testing.T) {
	t.Run("ValidCases", func(t *testing.T) {
		d := driver.NewGenericDialect()

		tbl := token.NewTable("users")
		if got := render.Table(d, *tbl); got != "users" {
			t.Errorf("expected %q, got %q", "users", got)
		}

		tbl = token.NewTable("users AS u")
		if got := render.Table(d, *tbl); got != "users AS u" {
			t.Errorf("expected %q, got %q", "users AS u", got)
		}

		tbl = token.NewTable("users", "u")
		if got := render.Table(d, *tbl); got != "users AS u" {
			t.Errorf("expected %q, got %q", "users AS u", got)
		}
	})

	t.Run("InvalidCases", func(t *testing.T) {
		d := driver.NewGenericDialect()

		tbl := (&token.Table{BaseToken: &token.BaseToken{}}).
			SetErrorWith("", fmt.Errorf("invalid"))
		if got := render.Table(d, *tbl); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}

		tbl = token.NewTable("users AS x", "y") // alias mismatch
		if got := render.Table(d, *tbl); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}

		tbl = token.NewTable("users, orders") // comma not allowed
		if got := render.Table(d, *tbl); got != "" {
			t.Errorf("expected empty string, got %q", got)
		}
	})

	t.Run("With", func(t *testing.T) {
		t.Run("NilDialect", func(t *testing.T) {
			tbl := token.NewTable("users AS u")
			if got := render.Table(nil, *tbl); got != "users AS u" {
				t.Errorf("expected %q, got %q", "users AS u", got)
			}
		})

		t.Run("PostgresDialect", func(t *testing.T) {
			tbl := token.NewTable("users AS u")
			if got := render.Table(driver.NewPostgresDialect(), *tbl); got != `"users" AS "u"` {
				t.Errorf("expected %q, got %q", `"users" AS "u"`, got)
			}
		})
	})
}
