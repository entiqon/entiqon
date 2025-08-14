// File: db/token/column_test.go

package token_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/token"
)

func TestColumn(t *testing.T) {
	t.Run("Constructors", func(t *testing.T) {
		t.Run("NewField", func(t *testing.T) {
			t.Run("Valid", func(t *testing.T) {
				in := "email"
				expr := "email"
				alias := ""

				f := token.NewField(in, expr, alias, false)
				if f == nil {
					t.Fatal("expected non-nil Field")
				}
				if f.Error != nil {
					t.Fatalf("unexpected error: %v", f.Error)
				}
				if !f.IsValid() {
					t.Fatal("expected field to be valid")
				}
				if f.IsAliased() {
					t.Error("expected IsAliased to be false")
				}
				if f.Name() != "email" {
					t.Errorf("Name() = %q, want %q", f.Name(), "email")
				}
				if f.IsRaw {
					t.Error("expected IsRaw to be false")
				}
				if got := f.Raw(); got != "email" {
					t.Errorf("Raw() = %q, want %q", got, "email")
				}
				if got := f.Render(); got != "email" {
					t.Errorf("Render() = %q, want %q", got, "email")
				}
			})

			t.Run("NoAlias", func(t *testing.T) {
				in := "email"
				expr := "email"
				alias := ""

				f := token.NewField(in, expr, alias, false)
				if f == nil {
					t.Fatal("expected non-nil Field")
				}
				if f.Error != nil {
					t.Fatalf("unexpected error: %v", f.Error)
				}
				if !f.IsValid() {
					t.Fatal("expected field to be valid")
				}
				if f.IsAliased() {
					t.Error("expected IsAliased to be false")
				}
			})

			t.Run("InlineAlias", func(t *testing.T) {
				t.Run("AS", func(t *testing.T) {
					input := "LOWER(m3_cuno || '-' || partnership_id) AS id"
					expr := "LOWER(m3_cuno || '-' || partnership_id)"
					alias := "id"

					f := token.NewField(input, expr, alias, true)
					if f == nil {
						t.Fatal("expected non-nil Field")
					}
					if f.Error != nil {
						t.Fatalf("unexpected error: %v", f.Error)
					}
					if !f.IsValid() {
						t.Fatal("expected field to be valid")
					}
					if !f.IsAliased() {
						t.Error("expected IsAliased to be true")
					}
					if f.Name() != "id" {
						t.Errorf("Name() = %q, want %q", f.Name(), "id")
					}
					if f.Input != input {
						t.Errorf("Input = %q, want %q", f.Input, input)
					}
					if f.Expr != expr {
						t.Errorf("Expr = %q, want %q", f.Expr, expr)
					}
					if f.Alias != alias {
						t.Errorf("Alias = %q, want %q", f.Alias, alias)
					}
					if !f.IsRaw {
						t.Error("expected IsRaw to be true")
					}
					if got, want := f.Raw(), expr+" AS "+alias; got != want {
						t.Errorf("Raw() = %q, want %q", got, want)
					}
					if got, want := f.Render(), expr+" AS "+alias; got != want {
						t.Errorf("Render() = %q, want %q", got, want)
					}
				})
			})

			t.Run("AliasTrimmed", func(t *testing.T) {
				in := "created_at AS ts"
				expr := "created_at"
				alias := "   Ts   "

				f := token.NewField(in, expr, alias, false)

				if f.Error != nil {
					t.Fatalf("unexpected error: %v", f.Error)
				}
				if f.Alias != "Ts" {
					t.Errorf("Alias trimmed = %q, want %q", f.Alias, "Ts")
				}
				if f.Name() != "ts" {
					t.Errorf("Name() = %q, want %q", f.Name(), "ts")
				}
				if got, want := f.Raw(), "created_at AS Ts"; got != want {
					t.Errorf("Raw() = %q, want %q", got, want)
				}
			})

			t.Run("Invalid", func(t *testing.T) {
				t.Run("EmptyExpr", func(t *testing.T) {
					f := token.NewField("", "   ", "", false)
					if f.Error == nil {
						t.Fatal("expected error for empty expression, got nil")
					}
					if f.IsValid() {
						t.Error("expected field to be invalid")
					}
					if got, want := f.Error.Error(), "expression cannot be empty"; got != want {
						t.Errorf("error = %q, want %q", got, want)
					}
				})

				t.Run("NameEmpty", func(t *testing.T) {
					// Expr with no alphanumeric chars -> derived name becomes empty
					f := token.NewField("!@#$", "!@#$", "", false)

					if f.Error == nil {
						t.Fatal("expected error for empty derived field name, got nil")
					}
					if f.IsValid() {
						t.Error("expected field to be invalid")
					}
					if got, want := f.Error.Error(), "derived field name cannot be empty"; got != want {
						t.Errorf("error = %q, want %q", got, want)
					}
				})
			})
		})
	})

	t.Run("Methods", func(t *testing.T) {
		t.Run("IsAliased", func(t *testing.T) {
			cases := []struct {
				alias string
				want  bool
			}{
				{"", false},
				{"  ", false},
				{"alias", true},
				{"ALIAS", true},
			}
			for _, tc := range cases {
				f := token.Field{Alias: tc.alias}
				got := f.IsAliased()
				if got != tc.want {
					t.Errorf("IsAliased() alias=%q got %v, want %v", tc.alias, got, tc.want)
				}
			}
		})

		t.Run("IsErrored", func(t *testing.T) {
			f := token.Field{Expr: "field", Alias: "f"}
			if f.IsErrored() {
				t.Error("expected IsErrored false when Error is nil")
			}
			f.Error = errors.New("some error")
			if !f.IsErrored() {
				t.Error("expected IsErrored true when Error set")
			}
		})

		t.Run("IsValid", func(t *testing.T) {
			f := token.Field{Expr: "field", Alias: "f"}
			if !f.IsValid() {
				t.Error("expected IsValid true when no Error and Expr non-empty")
			}
			f.Error = errors.New("some error")
			if f.IsValid() {
				t.Error("expected IsValid false when Error set")
			}
			f = token.Field{Expr: "  "}
			if f.IsValid() {
				t.Error("expected IsValid false when Expr is empty")
			}
			f = token.Field{Expr: "!@#$%^&*()"}
			if f.IsValid() {
				t.Error("expected IsValid false when derived Name is empty")
			}
		})

		t.Run("Name", func(t *testing.T) {
			cases := []struct {
				alias string
				expr  string
				want  string
			}{
				{"", "fieldName", "fieldname"},
				{"ALIAS", "ignored", "alias"},
				{"Id", "some_expr", "id"},
				{"", "Expr_With_123", "exprwith123"},
			}
			for _, tc := range cases {
				f := token.Field{Alias: tc.alias, Expr: tc.expr}
				got := f.Name()
				if got != tc.want {
					t.Errorf("Name() alias=%q expr=%q got %q, want %q", tc.alias, tc.expr, got, tc.want)
				}
			}
		})

		t.Run("Raw", func(t *testing.T) {
			f := token.Field{Expr: "field"}
			if got := f.Raw(); got != "field" {
				t.Errorf("Raw() without alias got %q, want %q", got, "field")
			}
			f.Alias = "alias"
			if got := f.Raw(); got != "field AS alias" {
				t.Errorf("Raw() with alias got %q, want %q", got, "field AS alias")
			}
		})

		t.Run("Render", func(t *testing.T) {
			f := token.Field{Expr: "LOWER(x)"}
			if got, want := f.Render(), "LOWER(x)"; got != want {
				t.Errorf("Render() got %q, want %q", got, want)
			}
			f = token.Field{Expr: "LOWER(x)", Alias: "id"}
			if got, want := f.Render(), "LOWER(x) AS id"; got != want {
				t.Errorf("Render() got %q, want %q", got, want)
			}
			f = token.Field{Expr: "created_at", Alias: "  ts  "}
			if got, want := f.Render(), "created_at AS ts"; got != want {
				t.Errorf("Render() trims alias spaces got %q, want %q", got, want)
			}
		})

		t.Run("String", func(t *testing.T) {
			f := token.Field{Expr: "field", Alias: ""}
			got := f.String()
			if !strings.HasPrefix(got, "✅ Field(") {
				t.Errorf("String() no alias got %q, want prefix %q", got, "✅ Field(")
			}
			f.Alias = "Alias"
			got = f.String()
			if !strings.Contains(got, "alias: true") {
				t.Errorf("String() with alias got %q, want alias: true", got)
			}
			f.Error = errors.New("some error")
			got = f.String()
			if !strings.Contains(got, "⛔️") || !strings.Contains(got, "some error") {
				t.Errorf("String() errored got %q, want icon and error message", got)
			}
		})
	})
}
