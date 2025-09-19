package table_test

import (
	"strings"
	"testing"

	"github.com/entiqon/db/token/table"
	"github.com/entiqon/db/token/types/identifier"
)

func TestTable(t *testing.T) {
	t.Run("Constructor", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			tbl := table.New()
			if tbl.Error() == nil {
				t.Error("expected error, got nil")
			}
			if !strings.Contains(tbl.Error().Error(), "empty input is not allowed") {
				t.Errorf("expected error 'empty input is not allowed', got %v", tbl.Error())
			}

			tbl = table.New(table.New("users"))
			if !strings.Contains(tbl.Error().Error(), "use Clone() instead") {
				t.Errorf("expected error contains 'unsupported type: field', got %v", tbl.Error())
			}

			tbl = table.New(123456)
			if tbl.Error().Error() != "expr has invalid format (type int)" {
				t.Errorf("expected error 'expr has invalid format (type int)', got %v", tbl.Error())
			}

			tbl = table.New("")
			if tbl.Error().Error() != "empty identifier is not allowed: \"\"" {
				t.Errorf("expected error 'empty identifier is not allowed', got %v", tbl.Error())
			}
		})

		t.Run("1-arg", func(t *testing.T) {
			t.Run("Default", func(t *testing.T) {
				tbl := table.New("table")
				if tbl.Error() != nil {
					t.Fatal("expected nil, got error")
				}
				if tbl.Input() != "table" {
					t.Errorf("expected table 'name', got %v", tbl.Name())
				}
				if tbl.ExpressionKind() != identifier.Identifier {
					t.Errorf("expected kind=Identifier, got %v", tbl.ExpressionKind())
				}
				if tbl.Name() != "table" {
					t.Errorf("expected table 'name', got %v", tbl.Name())
				}
			})
		})

		t.Run("2-args", func(t *testing.T) {
			t.Run("Default", func(t *testing.T) {
				src := table.New("users", "u")
				if src.IsErrored() {
					t.Fatalf("unexpected error: %v", src.Error())
				}
				if src.Alias() != "u" {
					t.Errorf("expected alias=u, got %q", src.Alias())
				}
			})

			t.Run("InvalidAlias", func(t *testing.T) {
				src := table.New("users", "123")
				if !src.IsErrored() {
					t.Fatal("expected error for invalid alias, got nil")
				}

				src = table.New("users", 123456)
				if !src.IsErrored() {
					t.Fatal("expected error for invalid alias, got nil")
				}

				src = table.New("users", false)
				if !src.IsErrored() {
					t.Fatal("expected error for invalid alias, got nil")
				}

				src = table.New("the craziness on the plate", false)
				if !src.IsErrored() {
					t.Fatal("expected error for invalid alias, got nil")
				}
			})

			t.Run("EmptyAlias", func(t *testing.T) {
				src := table.New("users", "")
				if !src.IsErrored() {
					t.Fatal("expected error for empty alias, got nil")
				}
			})
		})

		t.Run("TooManyArgs", func(t *testing.T) {
			src := table.New("users", "u", "extra")
			if !src.IsErrored() {
				t.Fatal("expected error for too many arguments, got nil")
			}
		})

		t.Run("Aggregate", func(t *testing.T) {
			src := table.New("COUNT(id) AS count")
			if !src.IsErrored() {
				t.Fatal("expected error for aggregate, got nil")
			}
		})

		t.Run("Literal", func(t *testing.T) {
			src := table.New("'users' AS u")
			if !src.IsErrored() {
				t.Fatal("expected error for literal, got nil")
			}
		})
	})

	t.Run("Contracts", func(t *testing.T) {
		t.Run("BaseToken", func(t *testing.T) {
			src := table.New("table", "t")
			if src.Error() != nil {
				t.Errorf("expected no error, got %v", src.Error())
			}
			if src.Input() != "table t" {
				t.Errorf("expected field, got %v", src.Input())
			}
			if src.ExpressionKind().String() != "Identifier" {
				t.Errorf("expected kind Identifier, got %s", src.ExpressionKind().String())
			}
			if src.Expr() != "table" {
				t.Errorf("expected field, got %v", src.Expr())
			}
			if !src.IsAliased() || src.Alias() != "t" {
				t.Errorf("expected alias, got %v", src.Alias())
			}
		})

		t.Run("TableToken", func(t *testing.T) {
			src := table.New("users u")
			if src.Name() != "users" {
				t.Errorf("expected name 'users', got %q", src.Name())
			}
		})

		t.Run("Clonable", func(t *testing.T) {
			src := table.New("users u")
			clone := src.Clone()
			if src.Render() != clone.Render() {
				t.Errorf("expected clone to render same, got %q vs %q", src.Render(), clone.Render())
			}
			if src == clone {
				t.Error("expected clone to be a different instance")
			}
		})

		t.Run("Debuggable", func(t *testing.T) {
			src := table.New("table", "t")
			if out := src.Debug(); !strings.Contains(out, "table") || !strings.Contains(out, "t") {
				t.Errorf("unexpected Debug output: %q", out)
			}

			// invalid case: should surface the error
			bad := table.New("")
			if out := bad.Debug(); !strings.Contains(out, "error") {
				t.Errorf("expected Debug() to mention error, got %q", out)
			}
		})

		t.Run("Errorable", func(t *testing.T) {
			src := table.New("the craziness in the plate is out of date") // garbage
			if !src.IsErrored() {
				t.Fatal("expected errored table, got valid")
			}
			if src.Error() == nil {
				t.Fatal("expected non-nil error")
			}
		})

		t.Run("Rawable", func(t *testing.T) {
			src := table.New("table")
			if src.IsRaw() {
				t.Errorf("expected Raw() to be false', got %t", src.IsRaw())
			}
			if !strings.Contains(src.Raw(), "table") {
				t.Errorf("expected Raw() to contain 'table', got %v", src.Raw())
			}

			src = table.New("(SELECT COUNT(field) FROM users WHERE id = 1) count")
			if !src.IsRaw() {
				t.Errorf("expected Raw() to be true, got %t", src.IsRaw())
			}
			if !strings.Contains(src.Raw(), "count") {
				t.Errorf("expected Raw() to contain 'count', got %v", src.Raw())
			}

			src = table.New("table", "t")
			if !strings.Contains(src.Raw(), "table") {
				t.Errorf("expected Raw() to contain 'table', got %v", src.Raw())
			}

			src = table.New("")
			if !strings.Contains(src.Raw(), "") {
				t.Errorf("expected Raw() to contain '', got %v", src.Raw())
			}
		})

		t.Run("Renderable", func(t *testing.T) {
			f := table.New("table", "t")
			if got := f.Render(); got != "table AS t" {
				t.Errorf("expected Render() 'table AS t', got %q", got)
			}

			bad := table.New("")
			if bad.Render() != "" {
				t.Errorf("expected %q, got %q", "", bad.Render())
			}
		})

		t.Run("Stringable", func(t *testing.T) {
			f := table.New("table")
			if got := f.String(); got != "Table(\"table\")" {
				t.Errorf("expected 'Table(\"table\")', got %q", got)
			}

			// Computed expression with alias
			f = table.New("table", "t")
			if got := f.String(); got != "Table(table AS t)" {
				t.Errorf("expected String() 'Table(\"table AS t\")', got %q", got)
			}

			// Invalid
			bad := table.New("")
			if !strings.Contains(bad.String(), "empty identifier is not allowed") {
				t.Errorf("expected contains 'empty identifier is not allowed' for invalid field, got %q", bad.String())
			}
		})

		t.Run("Validable", func(t *testing.T) {
			f := table.New()
			if f.IsValid() {
				t.Error("expected IsValid() to be true")
			}
		})
	})
}
