package render_test

import (
	"testing"

	"github.com/ialopezg/entiqon/driver"
	"github.com/ialopezg/entiqon/internal/build/render"
	"github.com/ialopezg/entiqon/internal/build/token"
	"github.com/stretchr/testify/assert"
)

func TestRenderColumn_ValidCases(t *testing.T) {
	d := driver.NewGenericDialect()

	col := token.NewColumn("id")
	assert.Equal(t, "id", render.Column(d, col))

	col = token.NewColumn("id", "uid")
	assert.Equal(t, "id AS uid", render.Column(d, col))

	col = token.NewColumn("users.id")
	assert.Equal(t, "users.id", render.Column(d, col))

	col = token.NewColumn("users.id", "uid")
	assert.Equal(t, "users.id AS uid", render.Column(d, col))

	// postgres dialect
	col = token.NewColumn("users.id", "uid")
	assert.Equal(t, `"users"."id" AS "uid"`, render.Column(driver.NewPostgresDialect(), col))
}

func TestRenderColumn_InvalidCases(t *testing.T) {
	d := driver.NewGenericDialect()

	// Invalid: empty name
	col := token.Column{}
	assert.Equal(t, "", render.Column(d, col))

	t.Run("NilDialect", func(t *testing.T) {
		col := token.NewColumn("users.id", "uid") // valid column
		result := render.Column(nil, col)
		assert.Equal(t, "users.id AS uid", result) // should match generic output
	})

	// Invalid: semantic alias mismatch
	col = token.NewColumn("id AS uid", "wrong")
	assert.Equal(t, "", render.Column(d, col))
}
