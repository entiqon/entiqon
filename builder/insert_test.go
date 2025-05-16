package builder_test

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/builder"
	"github.com/stretchr/testify/require"
)

func TestInsertBuilder_WithReturning(t *testing.T) {
	q := builder.NewInsert().
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		Returning("id", "created_at")

	sql, args, err := q.Build()
	require.NoError(t, err)
	require.Equal(t,
		"INSERT INTO users (id, name) VALUES (?, ?) RETURNING id, created_at",
		sql,
	)
	require.Equal(t, []any{1, "Watson"}, args)
	fmt.Printf("📦 Generated SQL Query: %s with values %+v\n", sql, args)
	fmt.Println("🕵️ Verified by Watson: All is sound in the INSERT logic, Holmes.")
}
