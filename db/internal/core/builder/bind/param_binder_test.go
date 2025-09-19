// File: db/internal/core/builder/bind/param_binder_test.go

package bind_test

import (
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/internal/core/builder/bind"
)

func TestParamBinder_Bind_Generic(t *testing.T) {
	binder := bind.NewParamBinder(driver.NewGenericDialect())
	placeholder := binder.Bind("admin")

	if placeholder != "?" {
		t.Errorf("expected %q, got %q", "?", placeholder)
	}
	if args := binder.Args(); len(args) != 1 || args[0] != "admin" {
		t.Errorf("expected args [admin], got %#v", args)
	}
}

func TestParamBinder_Bind_Postgres(t *testing.T) {
	binder := bind.NewParamBinder(driver.NewPostgresDialect())
	placeholder := binder.Bind("alpha")

	if placeholder != "$1" {
		t.Errorf("expected %q, got %q", "$1", placeholder)
	}
	if args := binder.Args(); len(args) != 1 || args[0] != "alpha" {
		t.Errorf("expected args [alpha], got %#v", args)
	}
}

func TestParamBinder_BindMany_Generic(t *testing.T) {
	binder := bind.NewParamBinder(driver.NewGenericDialect())
	placeholders := binder.BindMany(42, true, "active")

	expected := []string{"?", "?", "?"}
	for i, p := range placeholders {
		if p != expected[i] {
			t.Errorf("expected %q at index %d, got %q", expected[i], i, p)
		}
	}

	args := binder.Args()
	expectedArgs := []any{42, true, "active"}
	for i, v := range args {
		if v != expectedArgs[i] {
			t.Errorf("expected arg[%d]=%v, got %v", i, expectedArgs[i], v)
		}
	}
}

func TestParamBinder_BindMany_Postgres(t *testing.T) {
	binder := bind.NewParamBinder(driver.NewPostgresDialect())
	placeholders := binder.BindMany(42, true, "active")

	expected := []string{"$1", "$2", "$3"}
	for i, p := range placeholders {
		if p != expected[i] {
			t.Errorf("expected %q at index %d, got %q", expected[i], i, p)
		}
	}

	args := binder.Args()
	expectedArgs := []any{42, true, "active"}
	for i, v := range args {
		if v != expectedArgs[i] {
			t.Errorf("expected arg[%d]=%v, got %v", i, expectedArgs[i], v)
		}
	}
}

func TestParamBinder_ArgsReturnsBoundValues(t *testing.T) {
	binder := bind.NewParamBinder(driver.NewGenericDialect())
	binder.Bind("first")
	binder.Bind("second")

	args := binder.Args()
	expected := []any{"first", "second"}
	for i, v := range args {
		if v != expected[i] {
			t.Errorf("expected arg[%d]=%v, got %v", i, expected[i], v)
		}
	}
}

func TestParamBinder_WithPosition_Generic(t *testing.T) {
	binder := bind.NewParamBinderWithPosition(driver.NewGenericDialect(), 4)
	placeholder := binder.Bind("next")

	if placeholder != "?" {
		t.Errorf("expected %q, got %q", "?", placeholder)
	}
	args := binder.Args()
	if len(args) != 1 || args[0] != "next" {
		t.Errorf("expected args [next], got %#v", args)
	}
}

func TestParamBinder_WithPosition_Postgres(t *testing.T) {
	binder := bind.NewParamBinderWithPosition(driver.NewPostgresDialect(), 4)
	placeholder := binder.Bind("next")

	if placeholder != "$4" {
		t.Errorf("expected %q, got %q", "$4", placeholder)
	}
	args := binder.Args()
	if len(args) != 1 || args[0] != "next" {
		t.Errorf("expected args [next], got %#v", args)
	}
}
