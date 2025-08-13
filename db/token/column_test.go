// File: db/token/column_test.go

package token_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/token"
)

func TestColumn(t *testing.T) {
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
				col := token.Column{Alias: tc.alias}
				got := col.IsAliased()
				if got != tc.want {
					t.Errorf("IsAliased() with alias=%q got %v, want %v", tc.alias, got, tc.want)
				}
			}
		})

		t.Run("IsErrored", func(t *testing.T) {
			col := token.Column{Expr: "field", Alias: "f"}
			if col.IsErrored() {
				t.Error("expected IsErrored false when Error is nil")
			}
			col.Error = errors.New("some error")
			if !col.IsErrored() {
				t.Error("expected IsErrored true when Error set")
			}
		})

		t.Run("IsValid", func(t *testing.T) {
			col := token.Column{Expr: "field", Alias: "f"}
			if !col.IsValid() {
				t.Error("expected IsValid true when no Error and Expr non-empty")
			}
			col.Error = errors.New("some error")
			if col.IsValid() {
				t.Error("expected IsValid false when Error set")
			}
			col = token.Column{Expr: "  "}
			if col.IsValid() {
				t.Error("expected IsValid false when Expr is empty")
			}
			col = token.Column{Expr: "!@#$%^&*()"}
			if col.IsValid() {
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
				col := token.Column{Alias: tc.alias, Expr: tc.expr}
				got := col.Name()
				if got != tc.want {
					t.Errorf("Name() with alias=%q expr=%q got %q, want %q", tc.alias, tc.expr, got, tc.want)
				}
			}
		})

		t.Run("Raw", func(t *testing.T) {
			col := token.Column{Expr: "field"}
			if got := col.Raw(); got != "field" {
				t.Errorf("Raw() without alias got %q, want %q", got, "field")
			}
			col.Alias = "alias"
			if got := col.Raw(); got != "field AS alias" {
				t.Errorf("Raw() with alias got %q, want %q", got, "field AS alias")
			}
		})

		t.Run("Render", func(t *testing.T) {
			// No alias => expr only
			col := token.Column{Expr: "LOWER(m3_cuno || '-' || partnership_id)"}
			got := col.Render()
			want := "LOWER(m3_cuno || '-' || partnership_id)"
			if got != want {
				t.Errorf("Render() without alias got %q, want %q", got, want)
			}

			// With alias => expr AS alias (no dialect quoting)
			col = token.Column{
				Expr:  "LOWER(m3_cuno || '-' || partnership_id)",
				Alias: "id",
			}
			got = col.Render()
			want = "LOWER(m3_cuno || '-' || partnership_id) AS id"
			if got != want {
				t.Errorf("Render() with alias got %q, want %q", got, want)
			}

			// Trim spaces in alias
			col = token.Column{
				Expr:  "created_at",
				Alias: "  ts  ",
			}
			got = col.Render()
			want = "created_at AS ts"
			if got != want {
				t.Errorf("Render() trims alias spaces got %q, want %q", got, want)
			}
		})

		t.Run("String", func(t *testing.T) {
			col := token.Column{Expr: "field", Alias: ""}
			got := col.String()
			if !strings.HasPrefix(got, "✅ Column(") {
				t.Errorf("String() no alias got %q, want prefix %q", got, "✅ Column(")
			}
			col.Alias = "Alias"
			got = col.String()
			if !strings.Contains(got, "alias: true") {
				t.Errorf("String() with alias got %q, want alias: true", got)
			}
			col.Error = errors.New("some error")
			got = col.String()
			if !strings.Contains(got, "⛔️") || !strings.Contains(got, "some error") {
				t.Errorf("String() errored got %q, want icon and error message", got)
			}
		})
	})
}
