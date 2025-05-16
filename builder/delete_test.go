package builder_test

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/builder"
	"github.com/stretchr/testify/require"
)

func TestDeleteBuilder_Build(t *testing.T) {
	q := builder.NewDelete().
		From("users").
		Where("id = ?", 99).
		Returning("id")

	sql, args, err := q.Build()
	require.NoError(t, err)
	require.Equal(t,
		"DELETE FROM users WHERE id = ? RETURNING id",
		sql,
	)
	require.Equal(t, []any{99}, args)
	fmt.Printf("📦 Generated SQL Query: %s with values %+v\n", sql, args)
	fmt.Println("🕵️ Verified by Watson: All is sound in the DELETE logic, Holmes.")
}
