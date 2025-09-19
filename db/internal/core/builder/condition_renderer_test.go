// File: db/internal/core/builder/condition_renderer_test.go

package builder_test

import (
	"strings"
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/internal/core/builder"
	"github.com/entiqon/db/internal/core/builder/bind"
	"github.com/entiqon/db/internal/core/token"
)

func TestAppendCondition(t *testing.T) {
	base := token.NewCondition(token.ConditionSimple, "status", "active")
	and := token.NewCondition(token.ConditionAnd, "deleted", false)
	or := token.NewCondition(token.ConditionOr, "archived", false)

	result := builder.AppendCondition([]token.Condition{base}, and)
	if len(result) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(result))
	}
	if result[1].Type != token.ConditionAnd {
		t.Errorf("expected type %v, got %v", token.ConditionAnd, result[1].Type)
	}

	result = builder.AppendCondition(result, or)
	if len(result) != 3 {
		t.Fatalf("expected 3 conditions, got %d", len(result))
	}
	if result[2].Type != token.ConditionOr {
		t.Errorf("expected type %v, got %v", token.ConditionOr, result[2].Type)
	}

	t.Run("Invalid", func(t *testing.T) {
		valid := token.NewCondition(token.ConditionSimple, "status", "active")
		invalid := token.Condition{} // no key/operator

		result := builder.AppendCondition([]token.Condition{valid}, invalid)
		if len(result) != 1 {
			t.Errorf("expected 1 condition, got %d", len(result))
		}
		if result[0].Key != "status" {
			t.Errorf("expected key %q, got %q", "status", result[0].Key)
		}
	})
}

func TestRenderConditions_Generic(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		conditions := []token.Condition{
			token.NewCondition(token.ConditionSimple, "active", true),
		}
		sql, args, err := builder.RenderConditions(driver.NewGenericDialect(), conditions)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != "active = ?" {
			t.Errorf("expected %q, got %q", "active = ?", sql)
		}
		if len(args) != 1 || args[0] != true {
			t.Errorf("expected args [true], got %#v", args)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		binder := bind.NewParamBinder(driver.NewGenericDialect())
		sql, args, err := builder.RenderConditionsWithBinder(driver.NewGenericDialect(), nil, binder)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != "" {
			t.Errorf("expected empty sql, got %q", sql)
		}
		if args != nil {
			t.Errorf("expected nil args, got %#v", args)
		}
	})

	t.Run("Unsupported", func(t *testing.T) {
		c := token.NewCondition(token.ConditionSimple, "status", "active")
		c.Type = token.ConditionType(rune(999)) // simulate unsupported
		binder := bind.NewParamBinder(driver.NewGenericDialect())
		_, _, err := builder.RenderConditionsWithBinder(driver.NewGenericDialect(), []token.Condition{c}, binder)
		if err == nil || !contains(err.Error(), "unsupported condition type") {
			t.Errorf("expected unsupported condition type error, got %v", err)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		conditions := []token.Condition{{}} // invalid
		binder := bind.NewParamBinder(driver.NewGenericDialect())
		_, _, err := builder.RenderConditionsWithBinder(driver.NewGenericDialect(), conditions, binder)
		if err == nil || !contains(err.Error(), "invalid condition") {
			t.Errorf("expected invalid condition error, got %v", err)
		}
	})

	t.Run("WithAndCondition", func(t *testing.T) {
		conditions := []token.Condition{
			token.NewCondition(token.ConditionSimple, "status", "active"),
			token.NewCondition(token.ConditionAnd, "deleted", false),
		}
		binder := bind.NewParamBinder(driver.NewGenericDialect())
		sql, args, err := builder.RenderConditionsWithBinder(driver.NewGenericDialect(), conditions, binder)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "status = ? AND deleted = ?"
		if sql != expected {
			t.Errorf("expected %q, got %q", expected, sql)
		}
		if len(args) != 2 || args[0] != "active" || args[1] != false {
			t.Errorf("expected args [active,false], got %#v", args)
		}
	})
}

func TestRenderConditionsWithBinder_Postgres(t *testing.T) {
	conditions := []token.Condition{
		token.NewCondition(token.ConditionSimple, "email_verified", true),
		token.NewCondition(token.ConditionOr, "email_verified", false),
	}
	binder := bind.NewParamBinderWithPosition(driver.NewPostgresDialect(), 4)
	sql, args, err := builder.RenderConditionsWithBinder(driver.NewPostgresDialect(), conditions, binder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `"email_verified" = $4 OR "email_verified" = $5`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 || args[0] != true || args[1] != false {
		t.Errorf("expected args [true,false], got %#v", args)
	}
}

// helper contains checks substring membership
func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
