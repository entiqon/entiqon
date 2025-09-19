package field_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/entiqon/db/token/field"
	"github.com/entiqon/db/token/types/identifier"
)

func TestField(t *testing.T) {
	t.Run("Constructor", func(t *testing.T) {
		t.Run("New", func(t *testing.T) {
			t.Run("Error", func(t *testing.T) {
				f := field.New()
				if f.Error() == nil {
					t.Error("expected error, got nil")
				}
				if !strings.Contains(f.Error().Error(), "empty input is not allowed") {
					t.Errorf("expected error 'empty input is not allowed', got %v", f.Error())
				}

				f = field.New(field.New("id"))
				if !strings.Contains(f.Error().Error(), "use Clone() instead") {
					t.Errorf("expected error contains 'unsupported type: field', got %v", f.Error())
				}

				f = field.New(123456)
				if f.Error().Error() != "expr has invalid format (type int)" {
					t.Errorf("expected error 'expr has invalid format (type int)', got %v", f.Error())
				}

				f = field.New("")
				if f.Error().Error() != "empty identifier is not allowed: \"\"" {
					t.Errorf("expected error 'empty identifier is not allowed', got %v", f.Error())
				}
			})

			t.Run("1-arg", func(t *testing.T) {
				t.Run("Default", func(t *testing.T) {
					f := field.New("id")
					if f.Expr() != "id" {
						t.Errorf("expected id, got %v", f.Expr())
					}
				})

				t.Run("Aliased", func(t *testing.T) {
					f := field.New("field alias")
					if f.ExpressionKind() != identifier.Identifier || f.Expr() != "field" || f.Alias() != "alias" {
						t.Errorf("expected id, got %v", f.Expr())
					}

					f = field.New("field AS alias")
					if f.ExpressionKind() != identifier.Identifier || f.Expr() != "field" || f.Alias() != "alias" {
						t.Errorf("expected id, got %v", f.Expr())
					}
				})

				t.Run("StarField", func(t *testing.T) {
					f := field.New("*")
					if f.Error() != nil {
						t.Errorf("expected no error, got %v", f.Error())
					}

					f = field.New("* alias")
					if f.Error().Error() != "'*' cannot be aliased or raw" {
						t.Errorf("expected \"'*' cannot be aliased or raw\", got %v", f.Error())
					}

					f = field.New("* AS alias")
					if f.Error().Error() != "'*' cannot be aliased or raw" {
						t.Errorf("expected \"'*' cannot be aliased or raw\", got %v", f.Error())
					}
				})
			})

			t.Run("2-args", func(t *testing.T) {
				t.Run("Default", func(t *testing.T) {
					f := field.New("field", "alias")
					if f.ExpressionKind() != identifier.Identifier || f.Expr() != "field" || f.Alias() != "alias" {
						t.Errorf("expected id, got %v", f.Expr())
					}
				})

				t.Run("Error", func(t *testing.T) {
					f := field.New("field", "")
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "invalid alias identifier cannot be empty") {
						t.Errorf("expected error 'invalid alias identifier cannot be empty', got %v", f.Error())
					}

					f = field.New("field", 123)
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "alias must be a string, got int") {
						t.Errorf("expected error 'alias must be a string, got int', got %v", f.Error())
					}

					f = field.New("field", "123alias")
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "cannot start with digit") {
						t.Errorf("expected error contains 'cannot start with digit', got %v", f.Error())
					}

					f = field.New("field", "AS")
					if f.Error() == nil || !strings.Contains(f.Error().Error(), "alias is a reserved keyword") {
						t.Errorf("expected error contains 'alias is a reserved keyword', got %v", f.Error())
					}
				})

				t.Run("StarField", func(t *testing.T) {
					f := field.New("*", "alias")
					if f.Error().Error() != "'*' cannot be aliased or raw" {
						t.Errorf("expected \"'*' cannot be aliased or raw\", got %v", f.Error())
					}
				})
			})

			t.Run("TooManyArgs", func(t *testing.T) {
				f := field.New("field", "alias", "extra")
				if f.Error() == nil || !strings.Contains(f.Error().Error(), "invalid field constructor signature") {
					t.Errorf("expected error contains 'invalid field constructor signature', got %v", f.Error())
				}
			})
		})

		t.Run("NewWithTable", func(t *testing.T) {
			t.Run("Default", func(t *testing.T) {
				f := field.NewWithTable("users", "SUM(qty * price)", "line_total")
				if f.Error() != nil || f.Owner() == nil {
					t.Errorf("expected no error, got %v", f.Error())
				}
			})

			t.Run("Error", func(t *testing.T) {
				f := field.NewWithTable("", "field", "alias")
				if f.Error() == nil || f.Owner() != nil {
					t.Errorf("expected error to be non-nil, got %v", f.Error())
				}
			})
		})
	})

	t.Run("Contracts", func(t *testing.T) {
		t.Run("BaseToken", func(t *testing.T) {
			f := field.New("field AS alias")
			if f.Error() != nil {
				t.Errorf("expected no error, got %v", f.Error())
			}
			if f.Input() != "field AS alias" {
				t.Errorf("expected field, got %v", f.Input())
			}
			if f.Expr() != "field" {
				t.Errorf("expected field, got %v", f.Expr())
			}
			if !f.IsAliased() || f.Alias() != "alias" {
				t.Errorf("expected alias, got %v", f.Alias())
			}
		})

		t.Run("Token", func(t *testing.T) {
			// Valid field, no owner initially
			f := field.New("field", "alias")
			if f.HasOwner() || f.Owner() != nil {
				t.Errorf("expected nil, got %v", f.Owner())
			}

			owner := "orders"
			f.SetOwner(&owner)
			if !f.HasOwner() || *f.Owner() != owner {
				t.Error("expected HasOwner() to be true")
			}
		})

		t.Run("Clonable", func(t *testing.T) {
			orig := field.New("SUM(qty)", "total", true)
			cl := orig.Clone()
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
			f := field.New("id")
			if f.IsRaw() {
				t.Errorf("expected Raw() to be false', got %t", f.IsRaw())
			}

			f = field.New("COUNT(field) count")
			if !f.IsRaw() {
				t.Errorf("expected Raw() to be true, got %t", f.IsRaw())
			}
			if !strings.Contains(f.Raw(), "count") {
				t.Errorf("expected Raw() to contain 'count', got %v", f.Raw())
			}

			f = field.NewWithTable("table", "field", "alias")
			if !strings.Contains(f.Raw(), "table") {
				t.Errorf("expected Raw() to contain 'table', got %v", f.Raw())
			}
		})

		t.Run("Renderable", func(t *testing.T) {
			f := field.NewWithTable("users", "id")
			if got := f.Render(); got != "users.id" {
				t.Errorf("expected Render() 'users.id', got %q", got)
			}

			bad := field.New("")
			if bad.Render() != "" {
				t.Errorf("expected %q, got %q", "", bad.Render())
			}
		})

		t.Run("Stringable", func(t *testing.T) {
			f := field.New("field")
			if got := f.String(); got != "Field(\"field\")" {
				t.Errorf("expected String() 'field', got %q", got)
			}

			// Computed expression with alias
			f = field.New("field", "alias")
			if got := f.String(); got != "Field(\"field AS alias\")" {
				t.Errorf("expected String() 'field AS alias', got %q", got)
			}

			// With table owner
			f = field.NewWithTable("users", "id")
			if got := f.String(); got != "Field(\"users.id\")" {
				t.Errorf("expected String() 'users.id', got %q", got)
			}

			// Invalid
			bad := field.New("")
			if !strings.Contains(bad.String(), "empty identifier is not allowed") {
				t.Errorf("expected contains 'empty identifier is not allowed' for invalid field, got %q", bad.String())
			}
		})

		t.Run("Validable", func(t *testing.T) {
			f := field.New()
			if f.IsValid() {
				t.Error("expected IsValid() to be true")
			}
		})
	})
}
