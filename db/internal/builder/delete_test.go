package builder

import (
	"strings"
	"testing"

	"github.com/entiqon/db/internal/core/token"
)

func TestDeleteBuilder_From(t *testing.T) {
	t.Run("BasicUsage", func(t *testing.T) {
		sql, _, err := NewDelete(nil).From("users").Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(sql, "DELETE") {
			t.Errorf("expected SQL to contain DELETE, got %q", sql)
		}
	})

	t.Run("EmptyFrom", func(t *testing.T) {
		_, _, err := NewDelete(nil).From("").Build()
		if err == nil || !strings.Contains(err.Error(), "builder validation failed") {
			t.Errorf("expected builder validation error, got %v", err)
		}
	})

	t.Run("MissingFrom", func(t *testing.T) {
		_, _, err := NewDelete(nil).Build()
		if err == nil || !strings.Contains(err.Error(), "builder validation failed") {
			t.Errorf("expected builder validation error, got %v", err)
		}
	})
}

func TestDeleteBuilder_Where(t *testing.T) {
	qb := NewDelete(nil).From("users").Where("id", 100)

	t.Run("BasicUsage", func(t *testing.T) {
		sql, args, err := qb.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(sql, "WHERE id = ?") {
			t.Errorf("expected WHERE id = ?, got %q", sql)
		}
		if len(args) != 1 || args[0] != 100 {
			t.Errorf("expected args [100], got %#v", args)
		}
	})

	t.Run("AndWhere", func(t *testing.T) {
		sql, args, _ := qb.AndWhere("status = active").Build()
		if !strings.Contains(sql, "WHERE id = ? AND status = ?") {
			t.Errorf("expected AND condition, got %q", sql)
		}
		if len(args) != 2 || args[0] != 100 || args[1] != "active" {
			t.Errorf("expected args [100, active], got %#v", args)
		}
	})

	t.Run("OrWhere", func(t *testing.T) {
		sql, args, _ := qb.OrWhere("email_confirmed", false).Build()
		if !strings.Contains(sql, "WHERE id = ? AND status = ? OR email_confirmed = ?") {
			t.Errorf("expected OR condition, got %q", sql)
		}
		if len(args) != 3 || args[0] != 100 || args[1] != "active" || args[2] != false {
			t.Errorf("expected args [100, active, false], got %#v", args)
		}
	})
}

func TestDeleteBuilder_LimitClause(t *testing.T) {
	sql, _, _ := NewDelete(nil).
		From("logs").
		Where("archived", true).
		Limit(100).
		Build()

	if !strings.Contains(sql, "LIMIT 100") {
		t.Errorf("expected LIMIT clause, got %q", sql)
	}
}

func TestDeleteBuilder_BuildValidations(t *testing.T) {
	t.Run("NoDialect", func(t *testing.T) {
		b := NewDelete(nil)
		b.Dialect = nil
		_, _, err := b.Build()
		if err == nil || !strings.Contains(err.Error(), "no dialect set") {
			t.Errorf("expected no dialect set error, got %v", err)
		}
	})

	t.Run("HasErrors", func(t *testing.T) {
		qb := NewDelete(nil).Where("id = ?")
		_, _, err := qb.Build()
		if err == nil || !qb.Validator.HasErrors() || !strings.Contains(err.Error(), "builder validation failed") {
			t.Errorf("expected builder validation error, got %v", err)
		}
	})

	t.Run("Basic", func(t *testing.T) {
		qb := NewDelete(nil).From("users").Where("status =")
		_, _, err := qb.Build()
		if err == nil || !strings.Contains(err.Error(), "builder validation failed") {
			t.Errorf("expected builder validation error, got %v", err)
		}
	})

	t.Run("InvalidAndWhere", func(t *testing.T) {
		_, _, err := NewDelete(nil).From("users").AndWhere("", 123, 456).Build()
		if err == nil || !strings.Contains(err.Error(), "builder validation failed") {
			t.Errorf("expected builder validation error, got %v", err)
		}
	})

	t.Run("InvalidOrWhere", func(t *testing.T) {
		_, _, err := NewDelete(nil).From("users").OrWhere("", 123).Build()
		if err == nil || !strings.Contains(err.Error(), "builder validation failed") {
			t.Errorf("expected builder validation error, got %v", err)
		}
	})

	t.Run("InvalidConditionType", func(t *testing.T) {
		qb := NewDelete(nil).From("users")
		_, _, err := qb.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		qb.conditions = append(qb.conditions, token.Condition{
			Type: "💥", Key: "status = 'active'",
		})
		_, _, err = qb.Build()
		if err == nil || !strings.Contains(err.Error(), "unsupported condition type") {
			t.Errorf("expected unsupported condition type error, got %v", err)
		}
	})
}
