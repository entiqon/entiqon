package builder_test

import (
	"entiqon/builder"
	"testing"
)

func TestSelectBuilder(t *testing.T) {
	t.Run("General", func(t *testing.T) {
		qb := new(builder.SelectQueryBuilder).
			Select("id", "email").
			From("users").
			Where("email = 'john.doe@example.com'", "active = true").
			OrderBy("created_at DESC").
			Take(10)

		sql, err := qb.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedSQL := `SELECT id, email FROM users WHERE email = 'john.doe@example.com' AND active = true ORDER BY created_at DESC LIMIT 10`

		if sql != expectedSQL {
			t.Errorf("expected SQL: %s, got: %s", expectedSQL, sql)
		}
	})
}
