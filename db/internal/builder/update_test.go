package builder

import (
	"fmt"
	"strings"
	"testing"

	"github.com/entiqon/db/internal/core/errors"
	"github.com/entiqon/db/internal/core/token"
)

func TestUpdateBuilder_Table(t *testing.T) {
	sql, _, err := NewUpdate(nil).
		Table("users").
		Set("status", "active").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sql, "UPDATE users") {
		t.Errorf("expected SQL to contain %q, got %q", "UPDATE users", sql)
	}
}

func TestUpdateBuilder_Set(t *testing.T) {
	sql, args, err := NewUpdate(nil).
		Table("users").
		Set("status", "active").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sql, "SET status = ?") {
		t.Errorf("expected SQL to contain %q, got %q", "SET status = ?", sql)
	}
	if len(args) != 1 || args[0] != "active" {
		t.Errorf("expected args [active], got %#v", args)
	}
}

func TestUpdateBuilder_SetMultiple(t *testing.T) {
	sql, args, err := NewUpdate(nil).
		Table("users").
		Set("name", "Alice").
		Set("status", "verified").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sql, "SET name = ?") || !strings.Contains(sql, "status = ?") {
		t.Errorf("expected SQL to contain assignments, got %q", sql)
	}
	if len(args) != 2 || args[0] != "Alice" || args[1] != "verified" {
		t.Errorf("expected args [Alice verified], got %#v", args)
	}
}

func TestUpdateBuilder_Where(t *testing.T) {
	sql, args, err := NewUpdate(nil).
		Table("users").
		Set("name", "Watson").
		Where("id = 42").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sql, "WHERE id = ?") {
		t.Errorf("expected SQL to contain WHERE clause, got %q", sql)
	}
	if len(args) != 2 || args[0] != "Watson" || args[1] != 42 {
		t.Errorf("expected args [Watson 42], got %#v", args)
	}
}

func TestUpdateBuilder_AndWhere(t *testing.T) {
	sql, _, err := NewUpdate(nil).
		Table("users").
		Set("status", "inactive").
		Where("deleted", false).
		AndWhere("role", "admin").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sql, "WHERE deleted = ? AND role = ?") {
		t.Errorf("expected SQL to contain combined WHERE, got %q", sql)
	}
}

func TestUpdateBuilder_OrWhere(t *testing.T) {
	sql, _, err := NewUpdate(nil).
		Table("users").
		Set("active", true).
		Where("email_verified = true").
		OrWhere("status = ?", false).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sql, "WHERE email_verified = ? OR status = ?") {
		t.Errorf("expected OR condition, got %q", sql)
	}
}

func TestUpdateBuilder_BuildErrors(t *testing.T) {
	// Aliased column
	sql, args, err := NewUpdate(nil).
		Table("users").
		Set("email AS contact", "watson@example.com").
		Where("id", 1).
		Build()
	if err == nil {
		t.Errorf("expected error for aliased column")
	}
	if sql != "" || args != nil {
		t.Errorf("expected empty results on error, got sql=%q args=%#v", sql, args)
	}

	// Missing table
	_, _, err = NewUpdate(nil).
		Set("name", "Watson").
		Build()
	if err == nil || !strings.Contains(err.Error(), "requires a target table") {
		t.Errorf("expected table missing error, got %v", err)
	}

	// Missing assignments
	_, _, err = NewUpdate(nil).
		Table("users").
		Build()
	if err == nil || !strings.Contains(err.Error(), "must define at least one column assignment") {
		t.Errorf("expected assignment error, got %v", err)
	}

	// Invalid condition type
	q := NewUpdate(nil).Table("users").Set("name", "Watson")
	q.conditions = append(q.conditions, token.Condition{Type: "💣", Key: "broken = true"})
	_, _, err = q.Build()
	if err == nil || !strings.Contains(err.Error(), "unsupported condition type") {
		t.Errorf("expected unsupported condition type error, got %v", err)
	}
}

func TestUpdateBuilder_UseDialect_Postgres(t *testing.T) {
	sql, args, err := NewUpdate(nil).
		Set("active", true).
		Table("users").
		Where("email_verified", true).
		OrWhere("email_verified", false).
		UseDialect("postgres").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 3 {
		t.Errorf("expected 3 args, got %#v", args)
	}
	if !strings.Contains(sql, "WHERE \"email_verified\" = $2 OR \"email_verified\" = $3") {
		t.Errorf("expected Postgres quoting, got %q", sql)
	}
}

func TestUpdateBuilder_AddStageError(t *testing.T) {
	qb := NewUpdate(nil)
	qb.AddStageError(errors.StageWhere, fmt.Errorf("first"))
	qb.AddStageError(errors.StageWhere, fmt.Errorf("second"))

	errs := qb.Validator.GetErrors()
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
	if !strings.Contains(errs[0].Error.Error(), "first") ||
		!strings.Contains(errs[1].Error.Error(), "second") {
		t.Errorf("expected errors to contain first and second, got %#v", errs)
	}
}

func TestUpdateBuilder_AddStageError_NewGroup(t *testing.T) {
	qb := NewUpdate(nil)
	initialLen := len(qb.Validator.GetErrors())
	qb.AddStageError("OR", fmt.Errorf("or error"))

	errs := qb.Validator.GetErrors()
	if len(errs) != initialLen+1 {
		t.Errorf("expected %d errors, got %d", initialLen+1, len(errs))
	}
	last := errs[len(errs)-1]
	if last.Stage != "OR" {
		t.Errorf("expected Stage OR, got %v", last.Stage)
	}
	if !strings.Contains(last.Error.Error(), "or error") {
		t.Errorf("expected error message, got %v", last.Error)
	}
}

func TestUpdateBuilder_GetDialect_DefaultsToGeneric(t *testing.T) {
	u := NewUpdate(nil)
	u.BaseBuilder = BaseBuilder{} // reset

	d := u.GetDialect()
	if d == nil {
		t.Fatalf("expected non-nil dialect")
	}
	if got := d.GetName(); got != "generic" {
		t.Errorf("expected generic dialect, got %q", got)
	}
}

func TestUpdateBuilder_GetErrors(t *testing.T) {
	u := NewUpdate(nil)
	u.AddStageError("WHERE", fmt.Errorf("invalid field"))

	errs := u.Validator.GetErrors()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Stage != errors.StageWhere {
		t.Errorf("expected stage WHERE, got %v", errs[0].Stage)
	}
	if !strings.Contains(errs[0].Error.Error(), "invalid field") {
		t.Errorf("expected error message, got %v", errs[0].Error)
	}
}

func TestUpdateBuilder_UseDialectShortCircuit(t *testing.T) {
	u := NewUpdate(nil).UseDialect("generic")
	ptr1 := u.UseDialect("generic")
	if ptr1.Dialect.GetName() != u.Dialect.GetName() {
		t.Errorf("expected same dialect, got %q vs %q", ptr1.Dialect.GetName(), u.Dialect.GetName())
	}

	ptr2 := u.UseDialect("")
	if ptr2.Dialect.GetName() != u.Dialect.GetName() {
		t.Errorf("expected unchanged dialect on empty name")
	}
}

func TestUpdateBuilder_UseDialect_ResolvesNamed(t *testing.T) {
	u := NewUpdate(nil).UseDialect("postgres")
	if got := u.GetDialect().GetName(); got != "postgres" {
		t.Errorf("expected postgres dialect, got %q", got)
	}
}

func TestUpdateBuilder_BuildValidations(t *testing.T) {
	c := token.NewCondition(token.ConditionSimple, "id = ?")

	b := UpdateBuilder{}
	if !c.IsValid() {
		b.AddStageError("WHERE clause", fmt.Errorf("invalid clause"))
	}
	b.Table("users").Set("name", "Watson")

	t.Run("HasDialect", func(t *testing.T) {
		b := NewUpdate(nil)
		b.AddStageError("WHERE clause", fmt.Errorf("invalid clause"))
		b.conditions = []token.Condition{c}
		_, _, err := b.Build()
		if err == nil {
			t.Errorf("expected error")
		}
		if !b.HasDialect() {
			t.Errorf("expected HasDialect=true")
		}
		if b.Dialect.GetName() != "generic" {
			t.Errorf("expected dialect=generic, got %q", b.Dialect.GetName())
		}
	})

	t.Run("HasErrors", func(t *testing.T) {
		_, _, err := NewUpdate(nil).Build()
		if err == nil || !strings.Contains(err.Error(), "must define at least one column assignment") {
			t.Errorf("expected assignment error, got %v", err)
		}
	})
}
