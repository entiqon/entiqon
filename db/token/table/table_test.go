package table_test

import (
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/token/table"
)

func TestTable(t *testing.T) {
	t.Run("Constructors", func(t *testing.T) {
		t.Run("Plain", func(t *testing.T) {
			src := table.New("users")
			if src.IsErrored() {
				t.Fatalf("unexpected error: %v", src.Error())
			}
			if got := src.Render(); got != "users" {
				t.Errorf("expected Render=users, got %q", got)
			}
			if got := src.Raw(); got != "users" {
				t.Errorf("expected Raw=users, got %q", got)
			}
			if !src.IsValid() {
				t.Errorf("expected valid table")
			}
		})

		t.Run("Aliased", func(t *testing.T) {
			t.Run("BySpace", func(t *testing.T) {
				src := table.New("users u")
				if src.IsErrored() {
					t.Fatalf("unexpected error: %v", src.Error())
				}
				if got := src.Render(); got != "users AS u" {
					t.Errorf("expected Render=users AS u, got %q", got)
				}
				if got := src.Raw(); got != "users AS u" {
					t.Errorf("expected Raw=users AS u, got %q", got)
				}
				if !src.IsAliased() {
					t.Errorf("expected aliased=true")
				}
			})

			t.Run("ByAS", func(t *testing.T) {
				src := table.New("users AS u")
				if src.IsErrored() {
					t.Fatalf("unexpected error: %v", src.Error())
				}
				if got := src.Render(); got != "users AS u" {
					t.Errorf("expected Render=users AS u, got %q", got)
				}
				if !src.IsAliased() {
					t.Errorf("expected aliased=true")
				}
			})

			t.Run("ExplicitAlias", func(t *testing.T) {
				src := table.New("users", "u")
				if src.IsErrored() {
					t.Fatalf("unexpected error: %v", src.Error())
				}
				if !src.IsRaw() {
					t.Errorf("expected IsRaw=true")
				}
				if got := src.Render(); got != "users AS u" {
					t.Errorf("expected Render=users AS u, got %q", got)
				}
			})
		})

		t.Run("Raw", func(t *testing.T) {
			t.Run("InlineAlias", func(t *testing.T) {
				src := table.New("(SELECT COUNT(id) FROM users) AS t")
				if src.IsErrored() {
					t.Fatalf("unexpected error: %v", src.Error())
				}
				if !src.IsRaw() {
					t.Errorf("expected IsRaw=true for subquery")
				}
				if got := src.Render(); got != "(SELECT COUNT(id) FROM users) AS t" {
					t.Errorf("unexpected Render=%q", got)
				}
			})

			t.Run("ExplicitAlias", func(t *testing.T) {
				src := table.New("(SELECT COUNT(id) FROM users)", "t")
				if src.IsErrored() {
					t.Fatalf("unexpected error: %v", src.Error())
				}
				if !src.IsRaw() {
					t.Errorf("expected IsRaw=true for subquery")
				}
				if got := src.Render(); got != "(SELECT COUNT(id) FROM users) AS t" {
					t.Errorf("unexpected Render=%q", got)
				}
			})
		})

		t.Run("Errors", func(t *testing.T) {
			t.Run("NoArgs", func(t *testing.T) {
				src := table.New()
				if !src.IsErrored() {
					t.Fatal("expected error, got none")
				}
			})

			t.Run("Empty", func(t *testing.T) {
				src := table.New("")
				if !src.IsErrored() {
					t.Fatal("expected error, got none")
				}
			})

			t.Run("EmptyAlias", func(t *testing.T) {
				src := table.New("users", "")
				if !src.IsErrored() {
					t.Fatal("expected error when alias is empty")
				}
				if got := src.Error().Error(); got != "table and alias must be non-empty" {
					t.Errorf("unexpected error: %v", got)
				}
			})

			t.Run("EmptyName", func(t *testing.T) {
				src := table.New("", "u")
				if !src.IsErrored() {
					t.Fatal("expected error when name is empty")
				}
				if got := src.Error().Error(); got != "table and alias must be non-empty" {
					t.Errorf("unexpected error: %v", got)
				}
			})

			t.Run("InvalidFormat", func(t *testing.T) {
				src := table.New("users AS")
				if !src.IsErrored() {
					t.Fatal("expected error, got none")
				}
			})

			t.Run("TooManyArgs", func(t *testing.T) {
				src := table.New("users", "u", "extra")
				if !src.IsErrored() {
					t.Fatal("expected error, got none")
				}
			})

			t.Run("TooManyTokens", func(t *testing.T) {
				src := table.New("users one two three")
				if !src.IsErrored() {
					t.Fatal("expected error for too many tokens")
				}
				if got := src.Error().Error(); got != `too many tokens in "users one two three"` {
					t.Errorf("unexpected error: %v", got)
				}
			})

			t.Run("NoAlias", func(t *testing.T) {
				src := table.New("(SELECT COUNT(id) FROM users)")
				if !src.IsErrored() {
					t.Fatal("expected error for subquery without alias")
				}
				if !src.IsRaw() {
					t.Errorf("expected IsRaw=true even when errored")
				}
			})

			t.Run("InvalidFormatTwoTokens", func(t *testing.T) {
				src := table.New("users AS")
				if !src.IsErrored() {
					t.Fatal("expected error for invalid format")
				}
				if got := src.Error().Error(); got != `invalid format "users AS"` {
					t.Errorf("unexpected error: %v", got)
				}
			})

			t.Run("InvalidFormatThreeTokens", func(t *testing.T) {
				src := table.New("users WTF u")
				if !src.IsErrored() {
					t.Fatal("expected error for invalid format")
				}
				if got := src.Error().Error(); got != `invalid format "users WTF u"` {
					t.Errorf("unexpected error: %v", got)
				}
			})
		})
	})

	t.Run("Methods", func(t *testing.T) {
		src := table.New("users", "u")

		t.Run("Error", func(t *testing.T) {
			src := table.New("")
			err := src.Error()
			if !src.IsErrored() && err == nil {
				t.Fatal("expected error, got none")
			}
		})

		t.Run("Raw", func(t *testing.T) {
			t.Run("Valid", func(t *testing.T) {
				if got := src.Raw(); got != "users AS u" {
					t.Errorf("expected Raw=users AS u, got %q", got)
				}
			})

			t.Run("Invalid", func(t *testing.T) {
				src := table.New("", "u")
				raw := src.Raw()
				if raw != "" {
					t.Error("expected raw, got none")
				}
			})
		})

		t.Run("Render", func(t *testing.T) {
			t.Run("Valid", func(t *testing.T) {
				if got := src.Render(); got != "users AS u" {
					t.Errorf("expected Render=users AS u, got %q", got)
				}
			})

			t.Run("Invalid", func(t *testing.T) {
				src := table.New("", "u")
				raw := src.Render()
				if raw != "" {
					t.Error("expected raw, got none")
				}
			})
		})

		t.Run("String", func(t *testing.T) {
			t.Run("Valid", func(t *testing.T) {
				src := table.New("users")
				want := "✅ Table(users)"
				if got := src.String(); got != want {
					t.Errorf("expected String=%q, got %q", want, got)
				}
			})

			t.Run("ValidWithAlias", func(t *testing.T) {
				want := "✅ Table(users AS u)"
				if got := src.String(); got != want {
					t.Errorf("expected String=%q, got %q", want, got)
				}
			})

			t.Run("Invalid", func(t *testing.T) {
				src := table.New("", "u")
				want := "❌ Table(\" u\"): table and alias must be non-empty"
				if got := src.String(); got != want {
					t.Errorf("expected String=%q, got %q", want, got)
				}
			})
		})

		t.Run("Debug", func(t *testing.T) {
			t.Run("Valid", func(t *testing.T) {
				got := src.Debug()
				if got == "" || !strings.Contains(got, "✅") {
					t.Errorf("expected Debug to start with ✅, got %q", got)
				}
			})

			t.Run("Invalid", func(t *testing.T) {
				src := table.New("")
				got := src.Debug()
				if got == "" || !strings.Contains(got, "❌") {

				}
			})
		})

		t.Run("Clone", func(t *testing.T) {
			clone := src.Clone()
			if clone == nil {
				t.Fatal("expected non-nil clone")
			}
			if clone.Render() != src.Render() {
				t.Errorf("expected clone.Render=%q, got %q", src.Render(), clone.Render())
			}
		})
	})
}
