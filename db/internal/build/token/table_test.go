// File: db/internal/build/token/table_test.go

package token_test

import (
	"strings"
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/internal/build/token"
)

func TestTable(t *testing.T) {
	t.Run("Usage", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			tbl := token.NewTable("users")
			if !tbl.IsValid() {
				t.Errorf("expected valid table")
			}
			if tbl.GetName() != "users" {
				t.Errorf("expected name 'users', got %q", tbl.GetName())
			}
			if tbl.GetAlias() != "" {
				t.Errorf("expected no alias")
			}
		})

		t.Run("Alias", func(t *testing.T) {
			t.Run("Inline", func(t *testing.T) {
				tbl := token.NewTable("users AS u")
				if !tbl.IsValid() {
					t.Errorf("expected valid table")
				}
				if tbl.GetName() != "users" || tbl.GetAlias() != "u" {
					t.Errorf("expected users AS u, got %q AS %q", tbl.GetName(), tbl.GetAlias())
				}
			})
		})

		t.Run("Raw", func(t *testing.T) {
			t.Run("Aliased", func(t *testing.T) {
				tbl := token.NewTable("users")
				got := tbl.Raw()
				if got != "users" {
					t.Errorf("expected 'users', got %q", got)
				}
			})

			t.Run("Unaliased", func(t *testing.T) {
				tbl := token.NewTable("users AS u")
				got := tbl.Raw()
				if got != "users AS u" {
					t.Errorf("expected 'users AS u', got %q", got)
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

		t.Run("String", func(t *testing.T) {
			t.Run("Unaliased", func(t *testing.T) {
				tbl := token.NewTable("users")
				want := `Table("users") [aliased: false, errored: false]`
				if got := tbl.String(); got != want {
					t.Errorf("expected %q, got %q", want, got)
				}
			})

			t.Run("Aliased", func(t *testing.T) {
				tbl := token.NewTable("users AS u")
				want := `Table("users") [aliased: true, errored: false]`
				if got := tbl.String(); got != want {
					t.Errorf("expected %q, got %q", want, got)
				}
			})

			t.Run("Errored", func(t *testing.T) {
				tbl := token.NewTable("users AS u", "x") // mismatch triggers alias conflict
				if !tbl.HasError() {
					t.Fatal("expected alias conflict error")
				}

				out := tbl.String()
				if !strings.Contains(out, `errored: true`) {
					t.Errorf("expected 'errored: true' in output, got %q", out)
				}
				if !strings.Contains(out, `alias conflict`) {
					t.Errorf("expected alias conflict in output, got %q", out)
				}
			})

			t.Run("NilTable", func(t *testing.T) {
				var tbl *token.Table
				if got := tbl.String(); got != "Table(nil)" {
					t.Errorf("expected 'Table(nil)', got %q", got)
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
			if tbl.GetError() == nil || tbl.GetError().Error() != expected {
				t.Errorf("unexpected alias conflict error: got %v, want %q", tbl.GetError(), expected)
			}
		})

		t.Run("EmptyInput", func(t *testing.T) {
			tbl := token.NewTable("")
			if tbl.IsValid() {
				t.Errorf("expected invalid table for empty input")
			}
			if tbl.GetError() == nil || !strings.Contains(tbl.GetError().Error(), "table expression is empty") {
				t.Errorf("unexpected error: %v", tbl.GetError())
			}
		})

		t.Run("InvalidInput", func(t *testing.T) {
			tbl := token.NewTable("users, orders")
			if tbl.IsValid() {
				t.Errorf("expected invalid table for comma-separated input")
			}
			if tbl.GetError() == nil || !strings.Contains(tbl.GetError().Error(), "aliases must not be comma-separated") {
				t.Errorf("unexpected error: %v", tbl.GetError())
			}
		})

		t.Run("StartsWithAS", func(t *testing.T) {
			tbl := token.NewTable("AS u")
			if tbl.IsValid() {
				t.Errorf("expected invalid table for expression starting with AS")
			}
			if tbl.GetError() == nil || !strings.Contains(tbl.GetError().Error(), "cannot start with 'AS'") {
				t.Errorf("unexpected error: %v", tbl.GetError())
			}
		})
	})
}
