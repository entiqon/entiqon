package builder_test

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/builder"
	"github.com/stretchr/testify/require"
)

func TestUpdateBuilder_Build(t *testing.T) {
	q := builder.NewUpdate().
		Table("users").
		Set("name", "Watson").
		Set("status", "active").
		Where("id = ?", 101)

	sql, args, err := q.Build()
	require.NoError(t, err)
	require.Equal(t,
		"UPDATE users SET name = ?, status = ? WHERE id = ?",
		sql,
	)
	require.Equal(t, []any{"Watson", "active", 101}, args)
	fmt.Printf("📦 Generated SQL Query: %s\n", sql)
}
