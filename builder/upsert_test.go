package builder_test

import (
	"fmt"
	"testing"

	"github.com/ialopezg/entiqon/builder"
	"github.com/stretchr/testify/require"
)

func TestUpsertBuilder_Build(t *testing.T) {
	q := builder.NewUpsert().
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		OnConflict("id").
		DoUpdateSet(map[string]string{
			"name": "EXCLUDED.name",
		})

	sql, args, err := q.Build()
	require.NoError(t, err)
	require.Equal(t,
		"INSERT INTO users (id, name) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name",
		sql,
	)
	require.Equal(t, []any{1, "Watson"}, args)
	fmt.Printf("📦 Generated SQL Query: %s\n with params: %+v", sql, args)
}
