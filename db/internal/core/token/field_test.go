package token_test

import (
	"testing"

	"github.com/entiqon/db/internal/core/token"
)

func TestField_Basic(t *testing.T) {
	col := token.Field("id")
	if col.Name != "id" {
		t.Errorf("expected Name=%q, got %q", "id", col.Name)
	}
	if col.IsRaw {
		t.Errorf("expected IsRaw=false, got true")
	}
	if col.Alias != "" {
		t.Errorf("expected empty Alias, got %q", col.Alias)
	}
}

func TestFieldExpr_WithAlias(t *testing.T) {
	expr := token.FieldExpr("COUNT(*)", "total")
	if expr.Name != "COUNT(*)" {
		t.Errorf("expected Name=%q, got %q", "COUNT(*)", expr.Name)
	}
	if expr.Alias != "total" {
		t.Errorf("expected Alias=%q, got %q", "total", expr.Alias)
	}
	if !expr.IsRaw {
		t.Errorf("expected IsRaw=true, got false")
	}
}

func TestField_AsMethod(t *testing.T) {
	aliased := token.Field("created_at").As("created")
	if aliased.Name != "created_at" {
		t.Errorf("expected Name=%q, got %q", "created_at", aliased.Name)
	}
	if aliased.Alias != "created" {
		t.Errorf("expected Alias=%q, got %q", "created", aliased.Alias)
	}
}

func TestIsValid(t *testing.T) {
	if !token.Field("status").IsValid() {
		t.Errorf("expected Field(status) to be valid")
	}
	if token.Field("").IsValid() {
		t.Errorf("expected Field(\"\") to be invalid")
	}
}

func TestField_WithAliasInline(t *testing.T) {
	f := token.Field("first_name AS name")
	if f.Name != "first_name" {
		t.Errorf("expected Name=%q, got %q", "first_name", f.Name)
	}
	if f.Alias != "name" {
		t.Errorf("expected Alias=%q, got %q", "name", f.Alias)
	}
	if !f.IsValid() {
		t.Errorf("expected field to be valid")
	}
}

func TestField_WithAliasParams(t *testing.T) {
	f := token.Field("first_name", "name")
	if f.Name != "first_name" {
		t.Errorf("expected Name=%q, got %q", "first_name", f.Name)
	}
	if f.Alias != "name" {
		t.Errorf("expected Alias=%q, got %q", "name", f.Alias)
	}
	if !f.IsValid() {
		t.Errorf("expected field to be valid")
	}
}

func TestFieldsFromExpr_CommaSeparated(t *testing.T) {
	fields := token.FieldsFromExpr("id, first_name AS name, email AS contact")
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	if fields[0].Name != "id" {
		t.Errorf("expected fields[0].Name=%q, got %q", "id", fields[0].Name)
	}
	if fields[1].Name != "first_name" || fields[1].Alias != "name" {
		t.Errorf("expected fields[1] Name=%q Alias=%q, got Name=%q Alias=%q",
			"first_name", "name", fields[1].Name, fields[1].Alias)
	}
	if fields[2].Name != "email" || fields[2].Alias != "contact" {
		t.Errorf("expected fields[2] Name=%q Alias=%q, got Name=%q Alias=%q",
			"email", "contact", fields[2].Name, fields[2].Alias)
	}
}

func TestField_CommaSeparatedPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		} else if r != "Field: comma-separated values not allowed in a single call. Call Field(...) separately for each." {
			t.Errorf("unexpected panic value: %v", r)
		}
	}()
	_ = token.Field("id, name")
}

func TestFieldToken_WithValue(t *testing.T) {
	f := token.Field("email")
	fv := f.WithValue("x@entiqon.dev")

	if fv.Name != "email" {
		t.Errorf("expected Name=%q, got %q", "email", fv.Name)
	}
	if fv.Value != "x@entiqon.dev" {
		t.Errorf("expected Value=%q, got %q", "x@entiqon.dev", fv.Value)
	}
}
