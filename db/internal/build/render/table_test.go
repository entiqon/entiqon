// File: db/internal/build/render/table_test.go

package render_test

import (
	"fmt"
	"testing"

	"github.com/entiqon/entiqon/db/driver"
	"github.com/entiqon/entiqon/db/internal/build/render"
	token2 "github.com/entiqon/entiqon/db/internal/build/token"
	"github.com/stretchr/testify/assert"
)

func TestRenderTable(t *testing.T) {
	t.Run("ValidCases", func(t *testing.T) {
		d := driver.NewGenericDialect()

		tbl := token2.NewTable("users")
		assert.Equal(t, "users", render.Table(d, *tbl))

		tbl = token2.NewTable("users AS u")
		assert.Equal(t, "users AS u", render.Table(d, *tbl))

		tbl = token2.NewTable("users", "u")
		assert.Equal(t, "users AS u", render.Table(d, *tbl))
	})

	t.Run("ValidCases", func(t *testing.T) {
		d := driver.NewGenericDialect()

		tbl := (&token2.Table{BaseToken: &token2.BaseToken{}}).
			SetErrorWith("", fmt.Errorf("invalid"))
		assert.Equal(t, "", render.Table(d, *tbl))

		tbl = token2.NewTable("users AS x", "y") // alias mismatch
		assert.Equal(t, "", render.Table(d, *tbl))

		tbl = token2.NewTable("users, orders") // comma not allowed
		assert.Equal(t, "", render.Table(d, *tbl))
	})

	t.Run("With", func(t *testing.T) {
		t.Run("NilDialect", func(t *testing.T) {
			tbl := token2.NewTable("users AS u")
			assert.Equal(t, "users AS u", render.Table(nil, *tbl)) // fallback to generic
		})

		t.Run("PostgresDialect", func(t *testing.T) {
			tbl := token2.NewTable("users AS u")
			assert.Equal(t, "\"users\" AS \"u\"", render.Table(driver.NewPostgresDialect(), *tbl))
		})
	})
}
