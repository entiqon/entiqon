package builder

import (
	"strings"
	"testing"

	"github.com/entiqon/db/driver"
)

func TestUpsertBuilder_WithDialect_EscapesIdentifiers(t *testing.T) {
	q := NewUpsert(driver.NewPostgresDialect()).
		Into("user profile").
		Columns("user id", "email").
		Values(99, "hello@test.dev").
		OnConflict("user id").
		DoUpdateSet(Assignment{Column: "email", Expr: "EXCLUDED.email"})

	sql, args, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO "user profile" ("user id", "email") VALUES ($1, $2) ON CONFLICT ("user id") DO UPDATE SET "email" = EXCLUDED.email`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 || args[0] != 99 || args[1] != "hello@test.dev" {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestUpsertBuilder_WithoutDialect(t *testing.T) {
	q := NewUpsert(nil).
		Into("user profile").
		Columns("user id", "email").
		Values(99, "hello@test.dev").
		OnConflict("user id").
		DoUpdateSet(Assignment{Column: "email", Expr: "EXCLUDED.email"})

	sql, args, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO user profile (user id, email) VALUES (?, ?) ON CONFLICT (user id) DO UPDATE SET email = EXCLUDED.email`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 || args[0] != 99 || args[1] != "hello@test.dev" {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestUpsertBuilder_Returning_WithoutDialectRawNames(t *testing.T) {
	q := NewUpsert(nil).
		Into("emails").
		Columns("id", "value").
		Values(101, "none@entiqon.dev").
		OnConflict("id").
		Returning("id", "value")

	sql, args, err := q.Build()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if sql != "" {
		t.Errorf("expected empty SQL, got %q", sql)
	}
	if args != nil {
		t.Errorf("expected nil args, got %#v", args)
	}
}

func TestUpsertBuilder_Returning_AppendsReturningClause(t *testing.T) {
	q := NewUpsert(driver.NewPostgresDialect()).
		Into("users").
		Columns("id", "email").
		Values(1, "dev@entiqon.dev").
		OnConflict("id").
		DoUpdateSet(Assignment{Column: "email", Expr: "EXCLUDED.email"}).
		Returning("id", "email")

	sql, args, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO "users" ("id", "email") VALUES ($1, $2) ON CONFLICT ("id") DO UPDATE SET "email" = EXCLUDED.email RETURNING "id", "email"`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != "dev@entiqon.dev" {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestUpsertBuilder_DoUpdateSet_AppendsAssignments(t *testing.T) {
	q := NewUpsert(nil).
		DoUpdateSet(
			Assignment{Column: "name", Expr: "EXCLUDED.name"},
			Assignment{Column: "email", Expr: "EXCLUDED.email"},
		)

	sql, args, _ := q.Build()
	if sql != "" {
		t.Errorf("expected empty SQL, got %q", sql)
	}
	if len(args) != 0 {
		t.Errorf("expected empty args, got %#v", args)
	}
	if got := q.DoUpdateSet().updateSet; len(got) != 2 {
		t.Errorf("expected 2 assignments, got %#v", got)
	}
}

func TestUpsertBuilder_OnConflict_AppendsConflictColumns(t *testing.T) {
	q := NewUpsert(driver.NewPostgresDialect()).
		Into("people").
		Columns("id", "email").
		Values(1, "someone@dev.com").
		OnConflict("id", "email")

	sql, args, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO "people" ("id", "email") VALUES ($1, $2) ON CONFLICT ("id", "email") DO NOTHING`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != "someone@dev.com" {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestUpsertBuilder_Build_DoUpdate(t *testing.T) {
	q := NewUpsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		OnConflict("id").
		DoUpdateSet(Assignment{Column: "name", Expr: "EXCLUDED.name"})

	sql, args, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO users (id, name) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != "Watson" {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestUpsertBuilder_Build_DoNothing(t *testing.T) {
	q := NewUpsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		OnConflict("id")

	sql, args, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO users (id, name) VALUES (?, ?) ON CONFLICT (id) DO NOTHING`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != "Watson" {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestUpsertBuilder_BuildValidations(t *testing.T) {
	t.Run("EmptyTable", func(t *testing.T) {
		_, _, err := NewUpsert(nil).Build()
		if err == nil || !strings.Contains(err.Error(), "requires a target table") {
			t.Errorf("expected table error, got %v", err)
		}
	})
	t.Run("HasDialect", func(t *testing.T) {
		_, _, err := NewUpsert(nil).Into("users").Columns("id").Build()
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		b := UpsertBuilder{}
		if got := b.GetDialect().GetName(); got != "generic" {
			t.Errorf("expected generic dialect, got %q", got)
		}
	})
	t.Run("HasErrors", func(t *testing.T) {
		_, _, err := NewUpsert(nil).Into("users").Columns("").Build()
		if err == nil || !strings.Contains(err.Error(), "at least one set of values is required") {
			t.Errorf("expected values error, got %v", err)
		}
	})
	t.Run("Returning", func(t *testing.T) {
		_, _, err := NewUpsert(nil).Into("users").Columns("id").Values(1).Returning("id").Build()
		if err == nil || !strings.Contains(err.Error(), "RETURNING not supported in dialect") {
			t.Errorf("expected returning error, got %v", err)
		}
	})
	t.Run("ColumnWithAlias", func(t *testing.T) {
		_, _, err := NewUpsert(nil).Into("users").
			Columns("id AS IDENTIFIER").
			Values(1).
			Build()
		if err == nil || !strings.Contains(err.Error(), "row 1 has 1 values") {
			t.Errorf("expected alias/value mismatch, got %v", err)
		}
	})
}
