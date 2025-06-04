package token_test

import (
	"strings"
	"testing"

	"github.com/entiqon/entiqon/driver"
	"github.com/entiqon/entiqon/internal/build/token"
)

func TestTable(t *testing.T) {
	t.Run("Usage", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			tbl := token.NewTable("users")
			if !tbl.IsValid() {
				t.Errorf("expected valid table")
			}
			if tbl.Name != "users" {
				t.Errorf("expected name 'users', got %q", tbl.Name)
			}
			if tbl.Alias != "" {
				t.Errorf("expected no alias")
			}
		})

		t.Run("Alias", func(t *testing.T) {
			t.Run("Inline", func(t *testing.T) {
				tbl := token.NewTable("users AS u")
				if !tbl.IsValid() {
					t.Errorf("expected valid table")
				}
				if tbl.Name != "users" || tbl.Alias != "u" {
					t.Errorf("expected users AS u, got %q AS %q", tbl.Name, tbl.Alias)
				}
			})
		})

		t.Run("Render", func(t *testing.T) {
			d := driver.NewPostgresDialect()

			t.Run("Unaliased", func(t *testing.T) {
				tbl := token.NewTable("users")
				got := tbl.Render(d)
				want := `"users"`
				if got != want {
					t.Errorf("Render mismatch: got %q, want %q", got, want)
				}
			})

			t.Run("Aliased", func(t *testing.T) {
				tbl := token.NewTable("users u")
				got := tbl.Render(d)
				want := `"users" AS "u"`
				if got != want {
					t.Errorf("Render mismatch: got %q, want %q", got, want)
				}
			})

			t.Run("EmptyTable", func(t *testing.T) {
				var tbl *token.Table
				if got := tbl.Render(d); got != "" {
					t.Errorf("expected empty render for nil table, got %q", got)
				}
			})
		})
	})

	t.Run("Validation", func(t *testing.T) {
		t.Run("ExplicitAliasMatch", func(t *testing.T) {
			tbl := token.NewTable("users AS u", "u")
			if !tbl.IsValid() {
				t.Errorf("expected valid table with matching alias")
			}
		})

		t.Run("AliasConflict", func(t *testing.T) {
			tbl := token.NewTable("users AS u", "x")
			if tbl.IsValid() {
				t.Errorf("expected invalid table due to alias conflict")
			}

			expected := `alias conflict: explicit alias "x" does not match inline alias "u"`
			if tbl.Error == nil || tbl.Error.Error() != expected {
				t.Errorf("unexpected alias conflict error: got %v, want %q", tbl.Error, expected)
			}
		})

		t.Run("EmptyInput", func(t *testing.T) {
			tbl := token.NewTable("")
			if tbl.IsValid() {
				t.Errorf("expected invalid table for empty input")
			}
			if tbl.Error == nil || !strings.Contains(tbl.Error.Error(), "table expression is empty") {
				t.Errorf("unexpected error: %v", tbl.Error)
			}
		})

		t.Run("InvalidInput", func(t *testing.T) {
			tbl := token.NewTable("users, orders")
			if tbl.IsValid() {
				t.Errorf("expected invalid table for comma-separated input")
			}
			if tbl.Error == nil || !strings.Contains(tbl.Error.Error(), "aliases must not be comma-separated") {
				t.Errorf("unexpected error: %v", tbl.Error)
			}
		})

		t.Run("StartsWithAS", func(t *testing.T) {
			tbl := token.NewTable("AS u")
			if tbl.IsValid() {
				t.Errorf("expected invalid table for expression starting with AS")
			}
			if tbl.Error == nil || !strings.Contains(tbl.Error.Error(), "cannot start with 'AS'") {
				t.Errorf("unexpected error: %v", tbl.Error)
			}
		})
	})
}
