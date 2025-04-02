package builder_test

import (
	"entiqon/builder"
	"entiqon/dialect"
	"testing"
)

func TestSelectBuilder(t *testing.T) {
	t.Run("Postgres", func(t *testing.T) {
		// setup
		pg := &dialect.PostgresDialect{}
		qb := builder.NewSelect(pg).
			Columns("id", "email").
			From("users").
			Where("email = ? AND active = ?", "joe@example.com", true).
			OrderBy("created_at DESC").
			Limit(10)

		query, args := qb.Build()

		expectedSQL := `SELECT "id", "email" FROM "users" WHERE email = $1 AND active = $2 ORDER BY created_at DESC LIMIT 10`
		expectedArgs := []any{"joe@example.com", true}

		if query != expectedSQL {
			t.Errorf("expected SQL: %s, got: %s", expectedSQL, query)
		}

		if len(args) != len(expectedArgs) {
			t.Fatalf("expected %d args, got %d", len(expectedArgs), len(args))
		}

		for i := range args {
			if args[i] != expectedArgs[i] {
				t.Errorf("arg %d mismatch: expected %v, got %v", i+1, expectedArgs[i], args[i])
			}
		}
	})

	t.Run("MySQL", func(t *testing.T) {
		mysql := &dialect.MySQLDialect{}
		qb := builder.NewSelect(mysql).
			Columns("id", "email").
			From("users").
			Where("email = ? AND active = ?", "joe@example.com", true).
			OrderBy("created_at DESC").
			Limit(10)

		query, args := qb.Build()

		expectedSQL := "SELECT `id`, `email` FROM `users` WHERE email = ? AND active = ? ORDER BY created_at DESC LIMIT 10"
		expectedArgs := []any{"joe@example.com", true}

		if query != expectedSQL {
			t.Errorf("expected SQL: %s, got: %s", expectedSQL, query)
		}

		if len(args) != len(expectedArgs) {
			t.Fatalf("expected %d args, got %d", len(expectedArgs), len(args))
		}

		for i := range args {
			if args[i] != expectedArgs[i] {
				t.Errorf("arg %d mismatch: expected %v, got %v", i+1, expectedArgs[i], args[i])
			}
		}
	})
}
