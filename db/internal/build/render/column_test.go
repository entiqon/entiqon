// File: db/internal/build/render/column_test.go

package render_test

import (
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/internal/build/render"
	"github.com/entiqon/db/internal/build/token"
)

func TestRenderColumn_ValidCases(t *testing.T) {
	d := driver.NewGenericDialect()

	col := token.NewColumn("id")
	if got := render.Column(d, *col); got != "id" {
		t.Errorf("got %s, want %s", got, "id")
	}

	col = token.NewColumn("id", "uid")
	if got := render.Column(d, *col); got != "id AS uid" {
		t.Errorf("got %s, want %s", got, "id AS uid")
	}

	col = token.NewColumn("users.id")
	if got := render.Column(d, *col); got != "users.id" {
		t.Errorf("got %s, want %s", got, "users.id")
	}

	col = token.NewColumn("users.id", "uid")
	if got := render.Column(d, *col); got != "users.id AS uid" {
		t.Errorf("got %s, want %s", got, "users.id AS uid")
	}

	// postgres dialect
	col = token.NewColumn("users.id", "uid")
	if got := render.Column(driver.NewPostgresDialect(), *col); got != `"users"."id" AS "uid"` {
		t.Errorf("got %s, want %s", got, `"users"."id" AS "uid"`)
	}
}

func TestRenderColumn_InvalidCases(t *testing.T) {
	d := driver.NewGenericDialect()

	// Invalid: empty name
	col := &token.Column{}
	if got := render.Column(d, *col); got != "" {
		t.Errorf("got %s, want '%s'", got, "")
	}

	t.Run("NilDialect", func(t *testing.T) {
		col := token.NewColumn("users.id", "uid") // valid column
		if got := render.Column(nil, *col); got != "users.id AS uid" {
			t.Errorf("got %s, want '%s'", got, "users.id AS uid")
		}
	})

	t.Run("InvalidAlias", func(t *testing.T) {
		col = token.NewColumn("id AS uid", "wrong")
		if got := render.Column(d, *col); got != "" {
			t.Errorf("got %s, want '%s'", got, "")
		}
	})
}
