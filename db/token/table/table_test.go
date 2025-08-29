package table_test

import (
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/token"
	"github.com/entiqon/entiqon/db/token/table"
)

type testStringer string

func (s testStringer) String() string { return string(s) }

func TestTable(t *testing.T) {
	t.Run("Constructor", func(t *testing.T) {
		t.Run("1-arg", func(t *testing.T) {
			t.Run("Error", func(t *testing.T) {
				// No args → must error
				src := table.New()
				if src.Error() == nil {
					t.Fatal("expected error for no args, got nil")
				}

				// Unsupported type (int) → must error
				src = table.New(123)
				if src.Error() == nil {
					t.Fatal("expected error for unsupported type (int), got nil")
				}

				tbl := table.New("users")
				// Passing a token directly → must error with Clone() hint
				src = table.New(tbl)
				if src.Error() == nil {
					t.Fatal("expected error for passing token directly, got nil")
				}
				if !strings.Contains(src.Error().Error(), "use Clone() instead") {
					t.Errorf("expected error to suggest Clone(), got %v", src.Error())
				}

				// Empty string → must error
				src = table.New("")
				if src.Error() == nil {
					t.Fatal("expected error for empty string input, got nil")
				}
			})

			t.Run("Identifier", func(t *testing.T) {
				t.Run("NoAlias", func(t *testing.T) {
					src := table.New("users")
					if src.IsErrored() {
						t.Fatalf("unexpected error: %v", src.Error())
					}
					if src.ExpressionKind() != token.Identifier {
						t.Errorf("expected kind=Identifier, got %v", src.ExpressionKind())
					}
					if src.Name() != "users" {
						t.Errorf("expected name=users, got %q", src.Name())
					}
				})

				t.Run("WithAlias", func(t *testing.T) {
					src := table.New("users u")
					if src.IsErrored() {
						t.Fatalf("unexpected error: %v", src.Error())
					}

					src = table.New("users AS u")
					if src.IsErrored() {
						t.Fatalf("unexpected error: %v", src.Error())
					}
					if src.ExpressionKind() != token.Identifier {
						t.Errorf("expected kind=Identifier, got %v", src.ExpressionKind())
					}
					if src.Name() != "users" {
						t.Errorf("expected name=users, got %q", src.Name())
					}
				})

				t.Run("Errored", func(t *testing.T) {
					src := table.New("users 123456")
					if !src.IsErrored() {
						t.Fatalf("unexpected error: %v", src.Error())
					}

					src = table.New("users AS 123456")
					if !src.IsErrored() {
						t.Fatalf("unexpected error: %v", src.Error())
					}
				})

				t.Run("GarbageInput", func(t *testing.T) {
					src := table.New("the craziness in the plate is out of date")
					if !src.IsErrored() {
						t.Fatalf("expected error for garbage input, got nil")
					}
					if !strings.Contains(src.Error().Error(), "invalid format") {
						t.Errorf(
							"expected error to contain 'invalid format', got %v",
							src.Error(),
						)
					}
				})
			})

			t.Run("Subquery", func(t *testing.T) {
				t.Run("WithAlias", func(t *testing.T) {
					src := table.New("(SELECT customer, order_id FROM orders WHERE customer_id = 1234) orders")
					if src.Error() != nil {
						t.Fatalf("expected no error, got %v", src.Error())
					}

					src = table.New("(SELECT customer, order_id FROM orders WHERE customer_id = 1234) AS orders")
					if src.Error() != nil {
						t.Fatalf("expected no error, got %v", src.Error())
					}
					if src.ExpressionKind() != token.Subquery {
						t.Errorf("expected kind=Subquery, got %v", src.ExpressionKind())
					}
					if src.Alias() != "orders" {
						t.Errorf("expected alias=orders, got %v", src.Alias())
					}
				})

				t.Run("Error", func(t *testing.T) {
					src := table.New("(SELECT customer, order_id FROM orders WHERE customer_id = 1234) 123456")
					if src.Error() == nil {
						t.Fatal("expected error, got nil")
					}
				})
			})

			t.Run("FunctionOrComputed", func(t *testing.T) {
				t.Run("NoAlias", func(t *testing.T) {
					src := table.New("JSON_EACH(data)")
					if src.Error() != nil {
						t.Fatalf("expected no error, got %v", src.Error())
					}
					if src.ExpressionKind() != token.Function {
						t.Errorf("expected kind=Function, got %v", src.ExpressionKind())
					}
				})

				t.Run("WithAlias", func(t *testing.T) {
					src := table.New("JSON_EACH(data) j")
					if src.Error() != nil {
						t.Fatalf("expected no error, got %v", src.Error())
					}

					src = table.New("JSON_EACH(data) AS j")
					if src.ExpressionKind() != token.Function {
						t.Errorf("expected kind=Function, got %v", src.ExpressionKind())
					}
					if src.Name() != "JSON_EACH(data)" {
						t.Errorf("expected name=JSON_EACH(data), got %q", src.Name())
					}
					if src.Alias() != "j" {
						t.Errorf("expected alias=j, got %q", src.Alias())
					}
				})

				t.Run("Error", func(t *testing.T) {
					src := table.New("JSON_EACH(data) 123456")
					if !src.IsErrored() {
						t.Fatalf("expected error for invalid alias, got nil")
					}
					if !strings.Contains(src.Error().Error(), "invalid alias") {
						t.Errorf("expected invalid alias error, got %v", src.Error())
					}

					src = table.New("JSON_EACH(data) AS json is awesome")
					if !src.IsErrored() {
						t.Fatalf("expected error for malformed function expr, got nil")
					}
				})
			})

			t.Run("Literal", func(t *testing.T) {
				src := table.New("'foo'")
				if !src.IsErrored() {
					t.Fatal("expected error for string literal, got nil")
				}

				src = table.New("42")
				if !src.IsErrored() {
					t.Fatal("expected error for numeric literal, got nil")
				}

				src = table.New("\"users\"")
				if !src.IsErrored() {
					t.Fatal("expected error for quoted literal, got nil")
				}
				if !strings.Contains(src.Error().Error(), "literal") {
					t.Errorf("expected error mentioning literal, got %v", src.Error())
				}
			})

			t.Run("Aggregated", func(t *testing.T) {
				src := table.New("SUM(qty)")
				if !src.IsErrored() {
					t.Fatal("expected error for aggregate function, got nil")
				}
				if !strings.Contains(src.Error().Error(), "aggregate") {
					t.Errorf("expected aggregate function error, got %v", src.Error())
				}
			})
		})

		t.Run("2-args", func(t *testing.T) {
			t.Run("ValidAlias", func(t *testing.T) {
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

			t.Run("StringerAlias", func(t *testing.T) {
				src := table.New("users", testStringer("u"))
				if src.IsErrored() {
					t.Fatalf("unexpected error: %v", src.Error())
				}
				if src.Alias() != "u" {
					t.Errorf("expected alias 'u', got %q", src.Alias())
				}
				if src.Expr() != "users" {
					t.Errorf("expected expr 'users', got %q", src.Expr())
				}
			})

			t.Run("EmptyAlias", func(t *testing.T) {
				src := table.New("users", "")
				if !src.IsErrored() {
					t.Fatal("expected error for empty alias, got nil")
				}
			})

			t.Run("TooManyArgs", func(t *testing.T) {
				src := table.New("users", "u", "extra")
				if !src.IsErrored() {
					t.Fatal("expected error for too many arguments, got nil")
				}
			})
		})
	})

	t.Run("Methods", func(t *testing.T) {
		t.Run("Contract", func(t *testing.T) {
			t.Run("TableToken", func(t *testing.T) {
				src := table.New("users u")
				if src.Name() != "users" {
					t.Errorf("expected name 'users', got %q", src.Name())
				}
			})

			t.Run("BaseToken", func(t *testing.T) {
				src := table.New("users u")
				if src.Input() != "users u" {
					t.Errorf("expected input 'users u', got %q", src.Input())
				}
				if src.Expr() != "users" {
					t.Errorf("expected expr 'users', got %q", src.Expr())
				}
				if src.Alias() != "u" {
					t.Errorf("expected alias 'u', got %q", src.Alias())
				}
				if !src.IsAliased() {
					t.Error("expected table to be aliased")
				}
				if src.ExpressionKind().String() != "IDENTIFIER" {
					t.Errorf("expected kind IDENTIFIER, got %s", src.ExpressionKind().String())
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
				// valid case
				valid := table.New("users u")
				got := valid.Debug()
				if !strings.Contains(got, "raw:") {
					t.Errorf("expected debug output with flags, got %q", got)
				}
				if !strings.Contains(got, "✅ Table") {
					t.Errorf("expected valid marker in debug output, got %q", got)
				}

				// invalid case
				invalid := table.New("the craziness in the plate is out of date")
				got = invalid.Debug()
				if !strings.Contains(got, "❌ Table") {
					t.Errorf("expected error marker in debug output, got %q", got)
				}
				if invalid.Error() == nil {
					t.Fatal("expected error in invalid table")
				}
				if !strings.Contains(got, invalid.Error().Error()) {
					t.Errorf("expected debug output to include error message, got %q", got)
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
				// valid case with alias
				withAlias := table.New("users u")
				if withAlias.Raw() != "users AS u" {
					t.Errorf("expected `users AS u`, got %q", withAlias.Raw())
				}

				// valid case without alias
				noAlias := table.New("users")
				if noAlias.Raw() != "users" {
					t.Errorf("expected `users`, got %q", noAlias.Raw())
				}

				// invalid case
				invalid := table.New("the craziness in the plate is out of date")
				if invalid.Raw() != "" {
					t.Errorf("expected empty string for invalid table, got %q", invalid.Raw())
				}
				if !invalid.IsErrored() {
					t.Fatal("expected errored table for invalid input")
				}
			})

			t.Run("Renderable", func(t *testing.T) {
				// valid case with alias
				withAlias := table.New("users u")
				if withAlias.Render() != "users AS u" {
					t.Errorf("expected `users AS u`, got %q", withAlias.Render())
				}

				// valid case without alias
				noAlias := table.New("users")
				if noAlias.Render() != "users" {
					t.Errorf("expected `users`, got %q", noAlias.Render())
				}

				// invalid case
				invalid := table.New("the craziness in the plate is out of date")
				if invalid.Render() != "" {
					t.Errorf("expected empty string for invalid table, got %q", invalid.Render())
				}
				if !invalid.IsErrored() {
					t.Fatal("expected errored table for invalid input")
				}
			})

			t.Run("Stringable", func(t *testing.T) {
				// valid case with alias
				withAlias := table.New("users u")
				got := withAlias.String()
				if !strings.Contains(got, "✅ Table(users AS u)") {
					t.Errorf("expected success marker with alias, got %q", got)
				}

				// valid case without alias
				noAlias := table.New("users")
				got = noAlias.String()
				if !strings.Contains(got, "✅ Table(users)") {
					t.Errorf("expected success marker without alias, got %q", got)
				}

				// invalid case
				invalid := table.New("the craziness in the plate is out of date")
				got = invalid.String()
				if !strings.Contains(got, "❌ Table") {
					t.Errorf("expected error marker, got %q", got)
				}
				if invalid.Error() == nil {
					t.Fatal("expected error in invalid table")
				}
			})

			t.Run("Validable", func(t *testing.T) {
				valid := table.New("users u")
				if !valid.IsValid() {
					t.Error("expected valid table")
				}

				invalid := table.New("the craziness in the plate is out of date")
				if invalid.IsValid() {
					t.Error("expected invalid table")
				}
			})
		})
	})
}
