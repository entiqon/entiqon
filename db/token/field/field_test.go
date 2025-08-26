// File: db/token/column_test.go

package field_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/token/field"
)

func TestField(t *testing.T) {
	t.Run("Constructors", func(t *testing.T) {
		t.Run("New", func(t *testing.T) {
			t.Run("InvalidCall", func(t *testing.T) {
				f := field.New()
				if f == nil {
					t.Fatal("expected non-nil field, got nil")
				}
				if f.Error() == nil || !strings.Contains(f.Error().Error(), "empty input is not allowed") {
					t.Errorf("expected error 'empty input is not allowed', got %v", f.Error())
				}
			})

			t.Run("StarField", func(t *testing.T) {
				f := field.New("*")
				if f == nil {
					t.Fatal("expected non-nil field, got nil")
				}
				got := f.String()
				if got != "✅ field(\"*\")" {
					t.Errorf("expected ✅ field(\"*\"), got %v", got)
				}

				f = field.New("*", "alias")
				got = f.String()
				if f.Error() == nil || got != "⛔ field(\"*\"): '*' cannot be aliased or raw" {
					t.Errorf("expected ⛔ field(\"*\"): '*' cannot be aliased or raw, got %v", got)
				}
			})

			t.Run("1-arg", func(t *testing.T) {
				t.Run("InvalidCall", func(t *testing.T) {
					f := field.New("")
					if f == nil && f.Error() == nil {
						t.Fatal("expected non-nil field, got nil")
					}
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "empty expr") {
						t.Errorf("expected error contains 'empty expr is not allowed', got %v", f.Error())
					}
					f = field.New(123)
					if f == nil {
						t.Fatal("expected non-nil field, got nil")
					}
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "expr has unsupported type") {
						t.Errorf("expected error contains 'expr has unsupported type', got %v", f.Error())
					}
					f = field.New(field.New("id"))
					if f == nil {
						t.Fatal("expected non-nil field, got nil")
					}
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "unsupported type: field") {
						t.Errorf("expected error contains 'unsupported type: field', got %v", f.Error())
					}
				})

				t.Run("Valid", func(t *testing.T) {
					t.Run("Default", func(t *testing.T) {
						f := field.New("id")
						if f.Error() != nil {
							t.Errorf("expected error 'nil', got %v", f.Error())
						}
					})

					t.Run("InlineAlias", func(t *testing.T) {
						f := field.New("id as user_id")
						if f == nil {
							t.Fatal("expected non-nil field, got nil")
						}
						f = field.New("(qty * price * discount) total")
						if f.Expr() != "(qty * price * discount)" || f.Alias() != "total" {
							t.Errorf("expected '(qty * price * discount)' AS total, got %q AS %q", f.Expr(), f.Alias())
						}
						f = field.New("(qty * price * discount) as total")
						if f.Expr() != "(qty * price * discount)" || f.Alias() != "total" {
							t.Errorf("expected '(qty * price * discount)' AS total, got %q AS %q", f.Expr(), f.Alias())
						}
					})

					t.Run("Computed", func(t *testing.T) {
						f := field.New("SUM(qty * price * discount) as total")
						if f == nil {
							t.Fatal("expected non-nil field, got nil")
						}
						if f.Expr() != "SUM(qty * price * discount)" || f.Alias() != "total" {
							t.Errorf("expected 'SUM(qty * price * discount)' AS total, got %q AS %q", f.Expr(), f.Alias())
						}
						f = field.New("SUM(qty * price * discount) total")
						if f.Expr() != "SUM(qty * price * discount)" || f.Alias() != "total" {
							t.Errorf("expected 'SUM(qty * price * discount)' AS total, got %q AS %q", f.Expr(), f.Alias())
						}
						f = field.New("SUM(qty * price * discount)")
						if f.Expr() != "SUM(qty * price * discount)" || !strings.Contains(f.Alias(), "expr_alias_") {
							t.Errorf("expected 'SUM(qty * price * discount)' with auto alias, got %q AS %q", f.Expr(), f.Alias())
						}
						f = field.New("col1 || '-' || col2")
						if f.Expr() != "col1 || '-' || col2" || !strings.Contains(f.Alias(), "expr_alias_") {
							t.Errorf("expected \"col1 || '-' || col2\" with alias like \"expr_alias_\", got %q AS %q", f.Expr(), f.Alias())
						}
					})
				})
			})

			t.Run("2-args", func(t *testing.T) {
				t.Run("InvalidCall", func(t *testing.T) {
					f := field.New("field", 123)
					if f == nil {
						t.Fatal("expected non-nil field, got nil")
					}
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "alias has unsupported type") {
						t.Errorf("expected error 'alias has unsupported type', got %v", f.Error())
					}
					f = field.New("field", "")
					if f == nil {
						t.Fatal("expected non-nil field, got nil")
					}
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "empty alias is not allowed") {
						t.Errorf("expected error 'empty alias is not allowed', got %v", f.Error())
					}
				})

				t.Run("Valid", func(t *testing.T) {
					t.Run("Default", func(t *testing.T) {
						f := field.New("field", "alias")
						if f == nil {
							t.Fatal("expected non-nil field, got nil")
						}
						if f.Error() != nil {
							t.Fatal("expected error to be nil, got non-nil")
						}
					})

					t.Run("Computed", func(t *testing.T) {
						f := field.New("(qty * price * discount)", "total")
						if f == nil {
							t.Fatal("expected non-nil field, got nil")
						}
						if f.Error() != nil {
							t.Fatal("expected error to be nil, got non-nil")
						}
					})
				})
			})

			t.Run("3-args", func(t *testing.T) {
				t.Run("InvalidCall", func(t *testing.T) {
					f := field.New("field", 123, true)
					if f == nil {
						t.Fatal("expected non-nil field, got nil")
					}
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "alias has unsupported type") {
						t.Errorf("expected error 'alias has unsupported type', got %v", f.Error())
					}
					f = field.New("field", "", true)
					if f == nil {
						t.Fatal("expected non-nil field, got nil")
					}
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "empty alias is not allowed") {
						t.Errorf("expected error 'empty alias is not allowed', got %v", f.Error())
					}
					f = field.New("field", "alias", 123)
					if f == nil {
						t.Fatal("expected non-nil field, got nil")
					}
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "isRaw has invalid type") {
						t.Errorf("expected error 'isRaw has invalid type', got %v", f.Error())
					}
				})

				t.Run("Valid", func(t *testing.T) {
					t.Run("Default", func(t *testing.T) {
						f := field.New("field", "alias", false)
						if f == nil {
							t.Fatal("expected non-nil field, got nil")
						}
						if f.Error() != nil {
							t.Fatal("expected error to be nil, got non-nil")
						}
					})

					t.Run("Computed", func(t *testing.T) {
						f := field.New("(qty * price * discount)", "total", true)
						if f == nil {
							t.Fatal("expected non-nil field, got nil")
						}
						if f.Error() != nil {
							t.Fatal("expected error to be nil, got non-nil")
						}
					})
				})
			})

			t.Run("n-args", func(t *testing.T) {
				f := field.New("field", "alias", false, 1234)
				if f == nil {
					t.Fatal("expected non-nil field, got nil")
				}
				if f.Error() == nil || f.Error().Error() != "invalid field constructor signature" {
					t.Errorf("expected error 'invalid field constructor signature', got %v", f.Error())
				}
			})
		})

		t.Run("NewWithTable", func(t *testing.T) {
			t.Run("InvalidCall", func(t *testing.T) {
				f := field.NewWithTable("", "SUM(qty * price)", "line_total")
				if f == nil {
					t.Fatal("expected non-nil field, got nil")
				}
				if f.Error() == nil || f.Error().Error() != "owner is empty" {
					t.Errorf("expected error to be non-nil, got %v", f.Error())
				}
			})

			t.Run("Valid", func(t *testing.T) {
				f := field.NewWithTable("order", "SUM(qty * price)", "line_total", true)
				if f == nil {
					t.Fatal("expected non-nil field, got nil")
				}
				if f.Error() != nil {
					t.Errorf("expected error to be nil, got %v", f.Error())
				}
				if !f.HasOwner() || *f.Owner() != "order" {
					t.Errorf("expected owner 'order', got %v", f.Owner())
				}
			})
		})
	})

	t.Run("Contracts", func(t *testing.T) {
		t.Run("BaseToken", func(t *testing.T) {
			f := field.New("SUM(qty) total")
			if !f.IsValid() {
				t.Error("expected IsValid() to be true")
			}
			if f.Input() != "SUM(qty) total" {
				t.Errorf("expected 'SUM(qty)', got %q", f.Input())
			}
			if f.Expr() != "SUM(qty)" {
				t.Errorf("expected 'SUM(qty)', got %q", f.Expr())
			}
			if !f.IsAliased() || f.Alias() != "total" {
				t.Errorf("expected 'total', got %q", f.Alias())
			}
			if !f.IsRaw() {
				t.Errorf("expected 'isRaw', got %t", f.IsRaw())
			}

			// Identifier case
			f = field.New("id")
			if f.IsRaw() {
				t.Errorf("expected IsRaw() to be false, got %t", f.IsRaw())
			}
			if f.IsAliased() {
				t.Errorf("expected IsAliased() to be false, got true with alias %q", f.Alias())
			}

			// Invalid field
			f = field.New("")
			if f.IsValid() {
				t.Error("expected IsValid() to be false")
			}
		})

		t.Run("Token", func(t *testing.T) {
			// Valid field, no owner initially
			f := field.New("SUM(qty) total")
			if f.HasOwner() {
				t.Error("expected HasOwner() to be false")
			}
			if f.Owner() != nil {
				t.Errorf("expected nil, got %v", *f.Owner())
			}

			owner := "orders"
			f.SetOwner(&owner)
			if !f.HasOwner() || *f.Owner() != owner {
				t.Error("expected HasOwner() to be true")
			}

			// Invalid field, owner should still be nil
			bad := field.New("")
			if bad.HasOwner() {
				t.Error("expected HasOwner() to be false on invalid field")
			}
			if bad.Owner() != nil {
				t.Errorf("expected nil owner on invalid field, got %v", *bad.Owner())
			}
		})

		t.Run("Clonable", func(t *testing.T) {
			orig := field.New("SUM(qty)", "total", true)
			cl := orig.Clone()
			if cl == orig {
				t.Fatal("expected different pointer")
			}
			if cl.Input() != orig.Input() || cl.Expr() != orig.Expr() ||
				cl.Alias() != orig.Alias() || cl.IsRaw() != orig.IsRaw() {
				t.Errorf("clone mismatch: got %+v, want %+v", cl, orig)
			}

			// clone preserves error state
			if field.New("").Clone().Error() == nil {
				t.Error("expected error to be preserved in clone")
			}

			owner := "orders"
			orig.SetOwner(&owner)
			cl = orig.Clone()
			if cl.Owner() == nil {
				t.Error("expected clone to preserve owner, got nil")
			}
			if cl.Owner() == &owner {
				t.Error("expected clone to deep-copy owner, got same pointer")
			}
		})

		t.Run("Debuggable", func(t *testing.T) {
			f := field.New("id user_id")
			if out := f.Debug(); !strings.Contains(out, "id") || !strings.Contains(out, "user_id") {
				t.Errorf("unexpected Debug output: %q", out)
			}

			// invalid case: should surface the error
			bad := field.New("")
			if out := bad.Debug(); !strings.Contains(out, "error") {
				t.Errorf("expected Debug() to mention error, got %q", out)
			}
		})

		t.Run("Errorable", func(t *testing.T) {
			// Valid field → no error
			f := field.New("id")
			if f.Error() != nil {
				t.Errorf("expected no error, got %v", f.Error())
			}
			if f.IsErrored() {
				t.Error("expected IsErrored() to be false")
			}

			// Invalid field → error set automatically
			bad := field.New("")
			if bad.Error() == nil || !bad.IsErrored() {
				t.Error("expected error to be non-nil and IsErrored() true")
			}

			// Manual SetError
			f.SetError(errors.New("forced error"))
			if f.Error() == nil || f.Error().Error() != "forced error" {
				t.Errorf("expected 'forced error', got %v", f.Error())
			}
			if !f.IsErrored() {
				t.Error("expected IsErrored() to be true after SetError")
			}
		})

		t.Run("Rawable", func(t *testing.T) {
			// Identifier
			f := field.New("id")
			if got := f.Raw(); got != "id" {
				t.Errorf("expected Raw() 'id', got %q", got)
			}

			// With table owner
			f = field.NewWithTable("users", "id")
			if got := f.Raw(); got != "users.id" {
				t.Errorf("expected Raw() 'users.id', got %q", got)
			}

			// Computed expression with alias
			f = field.New("SUM(qty)", "total", true)
			if got := f.Raw(); got != "SUM(qty) AS total" {
				t.Errorf("expected Raw() 'SUM(qty) AS total', got %q", got)
			}

			// Computed expression with alias and owner → owner should be ignored
			f = field.NewWithTable("orders", "SUM(qty)", "total", true)
			if got := f.Raw(); got != "SUM(qty) AS total" {
				t.Errorf("expected Raw() 'SUM(qty) AS total' (owner ignored), got %q", got)
			}

			// Subquery with alias
			f = field.New("(SELECT id FROM users) sub")
			if got := f.Raw(); got != "(SELECT id FROM users) AS sub" {
				t.Errorf("expected Raw() '(SELECT id FROM users) AS sub', got %q", got)
			}

			// Subquery with alias and owner
			f = field.NewWithTable("t", "(SELECT id FROM users)", "u", true)
			if got := f.Raw(); got != "(SELECT id FROM users) AS u" {
				t.Errorf("expected Raw() '(SELECT id FROM users) AS u', got %q", got)
			}

			// Invalid field → resolution is empty
			bad := field.New("")
			if got := bad.Raw(); got != "" {
				t.Errorf("expected Raw() to be empty string for invalid field, got %q", got)
			}
		})

		t.Run("Renderable", func(t *testing.T) {
			// Identifier with owner
			f := field.NewWithTable("users", "id")
			if got := f.String(); got != "✅ field(\"users.id\")" {
				t.Errorf("expected String() 'users.id', got %q", got)
			}

			// Computed expression with alias
			f = field.New("SUM(qty)", "total", true)
			if got := f.String(); got != "✅ field(\"SUM(qty) AS total\")" {
				t.Errorf("expected String() 'SUM(qty) AS total', got %q", got)
			}

			// Subquery with alias
			f = field.New("(SELECT id FROM users) sub")
			if got := f.String(); got != "✅ field(\"(SELECT id FROM users) AS sub\")" {
				t.Errorf("expected String() '(SELECT id FROM users) AS sub', got %q", got)
			}

			// Invalid field → String() must show diagnostic, never empty
			bad := field.New("")
			if got := bad.String(); got != "⛔ field(\"\"): empty expression is not allowed" {
				t.Errorf("expected %q, got %q", "⛔️ field(\"\"): empty expression is not allowed", got)
			}
		})

		t.Run("Stringable", func(t *testing.T) {
			// Identifier
			f := field.New("id")
			if got := f.String(); got != "✅ field(\"id\")" {
				t.Errorf("expected String() 'id', got %q", got)
			}

			// With table owner
			f = field.NewWithTable("users", "id")
			if got := f.String(); got != "✅ field(\"users.id\")" {
				t.Errorf("expected String() 'users.id', got %q", got)
			}

			// Computed expression with alias
			f = field.New("SUM(qty)", "total", true)
			if got := f.String(); got != "✅ field(\"SUM(qty) AS total\")" {
				t.Errorf("expected String() 'SUM(qty) AS total', got %q", got)
			}

			// Subquery
			f = field.New("(SELECT id FROM users) u")
			if got := f.String(); got != "✅ field(\"(SELECT id FROM users) AS u\")" {
				t.Errorf("expected String() '(SELECT id FROM users) AS u', got %q", got)
			}

			// Invalid
			bad := field.New("")
			if got := bad.String(); got != "⛔ field(\"\"): empty expression is not allowed" {
				t.Errorf("expected contains 'empty expression is not allowed' for invalid field, got %q", got)
			}
		})
	})

	t.Run("EdgeCases", func(t *testing.T) {
		t.Run("HasTrailingAliasWithoutAS", func(t *testing.T) {
			// Single token: false
			if field.HasTrailingAliasWithoutAS("id") {
				t.Error("expected false for single token expr")
			}

			// Valid trailing alias: true
			if !field.HasTrailingAliasWithoutAS("SUM(price) total") {
				t.Error("expected true when trailing alias looks like identifier")
			}

			// Invalid trailing alias (symbolic): false
			if field.HasTrailingAliasWithoutAS("col1 + col2 123") {
				t.Error("expected false when trailing token is not valid identifier")
			}
		})
	})
}
