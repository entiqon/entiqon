// File: db/builder/insert_test.go

package builder_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/internal/builder"
	"github.com/entiqon/db/internal/core/token"
)

func TestInsertBuilder_NoColumns(t *testing.T) {
	b := builder.NewInsert(nil).Into("users").Values(1, "Watson")
	_, _, err := b.BuildInsertOnly()
	if err == nil || !strings.Contains(err.Error(), "at least one column is required") {
		t.Errorf("expected error about required columns, got %v", err)
	}
}

func TestInsertBuilder_NoValues(t *testing.T) {
	b := builder.NewInsert(nil).Into("users").Columns("id", "name")
	_, _, err := b.BuildInsertOnly()
	if err == nil || !strings.Contains(err.Error(), "at least one set of values is required") {
		t.Errorf("expected error about required values, got %v", err)
	}
}

func TestInsertBuilder_WithReturning(t *testing.T) {
	q := builder.NewInsert(driver.NewPostgresDialect()).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		Returning("id", "created_at")

	sql, args, err := q.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO "users" ("id", "name") VALUES ($1, $2) RETURNING "id", "created_at"`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != "Watson" {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestInsertBuilder_WithDialect_Postgres(t *testing.T) {
	sql, _, err := builder.NewInsert(driver.NewPostgresDialect()).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO "users" ("id", "name") VALUES ($1, $2)`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
}

func TestInsertBuilder_MismatchedRowLength(t *testing.T) {
	b := builder.NewInsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1)

	_, _, err := b.BuildInsertOnly()
	if err == nil || !strings.Contains(err.Error(), "row 1 has 1 values") {
		t.Errorf("expected mismatch error, got %v", err)
	}
}

func TestInsertBuilder_BuildErrors(t *testing.T) {
	_, _, err := builder.NewInsert(nil).Columns("id").Values(1).Build()
	if err == nil {
		t.Errorf("expected error for missing table")
	}

	_, _, err = builder.NewInsert(nil).Into("users").Values(1).Build()
	if err == nil {
		t.Errorf("expected error for missing columns")
	}

	_, _, err = builder.NewInsert(nil).Into("users").Columns("id").Build()
	if err == nil {
		t.Errorf("expected error for missing values")
	}
}

func TestInsertBuilder_MismatchedValueCount(t *testing.T) {
	_, _, err := builder.NewInsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1).
		Build()
	if err == nil {
		t.Errorf("expected mismatched value count error")
	}
}

func TestInsertBuilder_MissingFieldError(t *testing.T) {
	_, _, err := builder.NewInsert(nil).
		Into("users").
		Columns("id", "name", "email").
		Values(1, "Watson").
		Build()

	if err == nil || !strings.Contains(err.Error(), "row 1 has 2 values") {
		t.Errorf("expected row mismatch error, got %v", err)
	}
}

func TestInsertBuilder_WithAliasedColumn(t *testing.T) {
	sql, args, err := builder.NewInsert(nil).
		Into("users").
		Columns("email AS contact").
		Values("watson@example.com").
		Build()

	if err == nil {
		t.Errorf("expected error for aliased column")
	}
	if sql != "" {
		t.Errorf("expected empty SQL, got %q", sql)
	}
	if args != nil {
		t.Errorf("expected nil args, got %#v", args)
	}
}

func TestInsertBuilder_ReturningWithoutDialectFails(t *testing.T) {
	_, _, err := builder.NewInsert(nil).
		Into("users").
		Columns("id").
		Values(1).
		Returning("id").
		Build()

	if err == nil || !strings.Contains(err.Error(), "builder validation failed") {
		t.Errorf("expected validation error, got %v", err)
	}
}

func TestInsertBuilder_Build_WithDialectNoReturning(t *testing.T) {
	sql, args, err := builder.NewInsert(driver.NewPostgresDialect()).
		Into("users").
		Columns("id").
		Values(1).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO "users" ("id") VALUES ($1)`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 1 || args[0] != 1 {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestInsertBuilder_ReturningWithGenericDialectFails(t *testing.T) {
	sql, args, err := builder.NewInsert(nil).
		Into("users").
		Columns("id").
		Values(1).
		Returning("id").
		Build()

	if err == nil || !strings.Contains(err.Error(), "builder validation failed") {
		t.Errorf("expected validation error, got %v", err)
	}
	if sql != "" {
		t.Errorf("expected empty SQL, got %q", sql)
	}
	if args != nil {
		t.Errorf("expected nil args, got %#v", args)
	}
}

func TestInsertBuilder_BuildInsertOnly_ValidInsert(t *testing.T) {
	b := builder.NewInsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson")

	sql, args, err := b.BuildInsertOnly()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO users (id, name) VALUES (?, ?)`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != "Watson" {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestInsertBuilder_BuildInsertOnly_MultiRowSuccess(t *testing.T) {
	sql, args, err := builder.NewInsert(nil).
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		Values(2, "Holmes").
		BuildInsertOnly()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO users (id, name) VALUES (?, ?), (?, ?)`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 4 || args[0] != 1 || args[1] != "Watson" || args[2] != 2 || args[3] != "Holmes" {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestInsertBuilder_BuildInsertOnly_MissingTableFails(t *testing.T) {
	_, _, err := builder.NewInsert(nil).
		Columns("id").
		Values(1).
		BuildInsertOnly()

	if err == nil || !strings.Contains(err.Error(), "requires a target table") {
		t.Errorf("expected missing table error, got %v", err)
	}
}

func TestInsertBuilder_BuildInsertOnly_TableWithDialect(t *testing.T) {
	sql, args, err := builder.NewInsert(driver.NewPostgresDialect()).
		Into("users").
		Columns("id").
		Values(1).
		BuildInsertOnly()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `INSERT INTO "users" ("id") VALUES ($1)`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 1 || args[0] != 1 {
		t.Errorf("unexpected args: %#v", args)
	}
}

func TestInsertBuilder_BuildInsertOnly_ColumnEscapingWithDialect(t *testing.T) {
	tests := []struct {
		name     string
		dialect  driver.Dialect
		expected string
	}{
		{"Generic", driver.NewGenericDialect(), "email"},
		{"MSSQL", driver.NewMSSQLDialect(), "[email]"},
		{"MySQL", driver.NewMySQLDialect(), "`email`"},
		{"PostgreSQL", driver.NewPostgresDialect(), `"email"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, _, err := builder.NewInsert(tt.dialect).
				Into("users").
				Columns("email").
				Values("x@example.com").
				BuildInsertOnly()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(sql, tt.expected) {
				t.Errorf("expected SQL to contain %q, got %q", tt.expected, sql)
			}
		})
	}
}

func TestInsertBuilder_BuildValidations(t *testing.T) {
	c := token.NewCondition(token.ConditionSimple, "id = ?")

	t.Run("EmptyTable", func(t *testing.T) {
		_, _, err := builder.NewInsert(nil).Build()
		if err == nil || !strings.Contains(err.Error(), "requires a target table") {
			t.Errorf("expected missing table error, got %v", err)
		}
	})

	t.Run("HasDialect", func(t *testing.T) {
		b := builder.NewInsert(nil).Into("users")
		_, _, err := b.Build()

		if !c.IsValid() {
			b.AddStageError("WHERE", fmt.Errorf("invalid clause"))
		}

		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if got := b.GetDialect().GetName(); got != "generic" {
			t.Errorf("expected generic dialect, got %q", got)
		}
	})

	t.Run("HasErrors", func(t *testing.T) {
		b := builder.NewInsert(nil).Into("users")
		_, _, err := b.Into("").Build()
		if err == nil || !strings.Contains(err.Error(), "builder validation failed") {
			t.Errorf("expected builder validation error, got %v", err)
		}
	})
}
