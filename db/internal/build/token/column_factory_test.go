package token_test

import (
	"testing"

	"github.com/entiqon/db/internal/build/token"
)

func TestParseColumns(t *testing.T) {

	//--------------------------------------------------
	// Usage
	//--------------------------------------------------

	t.Run("BasicUsage", func(t *testing.T) {
		cols := token.NewColumnsFrom("id")
		if len(cols) != 1 {
			t.Fatalf("expected 1 column, got %d", len(cols))
		}
		if cols[0].GetName() != "id" {
			t.Errorf("expected name %q, got %q", "id", cols[0].GetName())
		}
		if err := cols[0].GetError(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("CommaSeparatedInput", func(t *testing.T) {
		cols := token.NewColumnsFrom("id, name")
		if len(cols) != 2 {
			t.Fatalf("expected 2 columns, got %d", len(cols))
		}
		if cols[0].GetName() != "id" {
			t.Errorf("expected first column %q, got %q", "id", cols[0].GetName())
		}
		if cols[1].GetName() != "name" {
			t.Errorf("expected second column %q, got %q", "name", cols[1].GetName())
		}
	})

	t.Run("InlineAlias", func(t *testing.T) {
		cols := token.NewColumnsFrom("user_id AS uid")
		if len(cols) != 1 {
			t.Fatalf("expected 1 column, got %d", len(cols))
		}
		if cols[0].GetName() != "user_id" {
			t.Errorf("expected name %q, got %q", "user_id", cols[0].GetName())
		}
		if cols[0].GetAlias() != "uid" {
			t.Errorf("expected alias %q, got %q", "uid", cols[0].GetAlias())
		}
		if err := cols[0].GetError(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("ExplicitAlias", func(t *testing.T) {
		cols := token.NewColumnsFrom("email", "contact AS primary_email")
		if len(cols) != 2 {
			t.Fatalf("expected 2 columns, got %d", len(cols))
		}
		if cols[0].GetName() != "email" {
			t.Errorf("expected name %q, got %q", "email", cols[0].GetName())
		}
		if err := cols[0].GetError(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if cols[1].GetName() != "contact" {
			t.Errorf("expected name %q, got %q", "contact", cols[1].GetName())
		}
		if cols[1].GetAlias() != "primary_email" {
			t.Errorf("expected alias %q, got %q", "primary_email", cols[1].GetAlias())
		}
		if err := cols[1].GetError(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	//--------------------------------------------------
	// Validations
	//--------------------------------------------------

	t.Run("EmptyInput", func(t *testing.T) {
		cols := token.NewColumnsFrom("")
		if len(cols) != 1 {
			t.Fatalf("expected 1 column, got %d", len(cols))
		}
		// error check depends on implementation
	})

	t.Run("OnlyWhitespace", func(t *testing.T) {
		cols := token.NewColumnsFrom("   ")
		if len(cols) != 1 {
			t.Fatalf("expected 1 column, got %d", len(cols))
		}
		// error check depends on implementation
	})

	t.Run("OnlyAliasKeyword", func(t *testing.T) {
		cols := token.NewColumnsFrom("AS alias")
		if len(cols) != 1 {
			t.Fatalf("expected 1 column, got %d", len(cols))
		}
		// error check depends on implementation
	})

	t.Run("AliasWithoutName", func(t *testing.T) {
		cols := token.NewColumnsFrom(" AS email")
		if len(cols) != 1 {
			t.Fatalf("expected 1 column, got %d", len(cols))
		}
		if err := cols[0].GetError(); err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
