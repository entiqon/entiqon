package render_test

import (
	"fmt"
	"testing"

	"github.com/entiqon/entiqon/driver"
	"github.com/entiqon/entiqon/internal/build/render"
	"github.com/entiqon/entiqon/internal/build/token"
	"github.com/stretchr/testify/assert"
)

func TestRenderTable(t *testing.T) {
	t.Run("ValidCases", func(t *testing.T) {
		d := driver.NewGenericDialect()

		tbl := token.NewTable("users")
		assert.Equal(t, "users", render.Table(d, *tbl))

		tbl = token.NewTable("users AS u")
		assert.Equal(t, "users AS u", render.Table(d, *tbl))

		tbl = token.NewTable("users", "u")
		assert.Equal(t, "users AS u", render.Table(d, *tbl))
	})

	t.Run("ValidCases", func(t *testing.T) {
		d := driver.NewGenericDialect()

		tbl := (&token.Table{BaseToken: &token.BaseToken{}}).
			SetErrorWith("", fmt.Errorf("invalid"))
		assert.Equal(t, "", render.Table(d, *tbl))

		tbl = token.NewTable("users AS x", "y") // alias mismatch
		assert.Equal(t, "", render.Table(d, *tbl))

		tbl = token.NewTable("users, orders") // comma not allowed
		assert.Equal(t, "", render.Table(d, *tbl))
	})

	t.Run("With", func(t *testing.T) {
		t.Run("NilDialect", func(t *testing.T) {
			tbl := token.NewTable("users AS u")
			assert.Equal(t, "users AS u", render.Table(nil, *tbl)) // fallback to generic
		})

		t.Run("PostgresDialect", func(t *testing.T) {
			tbl := token.NewTable("users AS u")
			assert.Equal(t, "\"users\" AS \"u\"", render.Table(driver.NewPostgresDialect(), *tbl))
		})
	})
}
