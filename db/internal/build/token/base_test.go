// File: db/internal/build/token/base_test.go

package token_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/internal/build/token"
	"github.com/entiqon/db/internal/core/contract"
)

func TestBaseToken(t *testing.T) {
	t.Run("Usage", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			b := token.NewBaseToken("id", "user_id")
			got := b.RenderAlias(driver.NewPostgresDialect(), "u.id")
			if got != `u.id AS "user_id"` {
				t.Errorf("expected quoted alias, got %q", got)
			}
		})

		t.Run("ExplicitAlias", func(t *testing.T) {
			t.Run("NoConflict", func(t *testing.T) {
				b := token.NewBaseToken("users.id", "uid")
				if b.GetName() != "users.id" {
					t.Errorf("expected name 'users.id', got %q", b.GetName())
				}
				if b.GetAlias() != "uid" {
					t.Errorf("expected alias 'uid', got %q", b.GetAlias())
				}
				if b.GetError() != nil {
					t.Errorf("unexpected error: %v", b.GetError())
				}
			})

			t.Run("AliasConflict", func(t *testing.T) {
				b := token.NewBaseToken("users.id AS user_id", "uid")
				if b.GetAlias() != "uid" {
					t.Errorf("expected alias to be overridden to 'uid', got %q", b.GetAlias())
				}
				if b.GetError() == nil || b.GetError().Error() != `alias conflict: explicit alias "uid" does not match inline alias "user_id"` {
					t.Errorf("expected alias conflict error, got %v", b.GetError())
				}
			})
		})
	})

	t.Run("Members", func(t *testing.T) {
		t.Run("AliasOr", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var b *token.BaseToken = nil
				b.SetError("ignored", fmt.Errorf("structural error"))
				if got := b.AliasOr(); got != "" {
					t.Errorf("expected nil error, got %v", got)
				}
			})

			t.Run("EmptyAlias", func(t *testing.T) {
				tok := token.NewBaseToken("id")
				if got := tok.AliasOr(); got != "id" {
					t.Errorf("expected 'id', got %q", got)
				}
			})

			t.Run("Aliased", func(t *testing.T) {
				tok := token.NewBaseToken("id", "uid")
				if got := tok.AliasOr(); got != "uid" {
					t.Errorf("expected 'uid', got %q", got)
				}
			})
		})

		t.Run("GetAlias", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var b *token.BaseToken = nil
				b.SetError("ignored", fmt.Errorf("structural error"))
				if got := b.GetAlias(); got != "" {
					t.Errorf("expected nil error, got %v", got)
				}
			})

			t.Run("EmptyAlias", func(t *testing.T) {
				tok := token.NewBaseToken("id")
				if got := tok.GetAlias(); got != "" {
					t.Errorf("expected 'id', got %q", got)
				}
			})

			t.Run("Aliased", func(t *testing.T) {
				tok := token.NewBaseToken("id", "uid")
				if got := tok.GetAlias(); got != "uid" {
					t.Errorf("expected 'uid', got %q", got)
				}
			})
		})

		t.Run("GetError", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var b *token.BaseToken = nil
				if got := b.GetError(); got != nil {
					t.Errorf("expected nil error, got %v", got)
				}
			})

			t.Run("Full", func(t *testing.T) {
				b := token.NewBaseToken("id", "uid")
				err := fmt.Errorf("alias conflict")
				b.SetError("id AS uid", err)
				if b.GetError() == nil || b.GetError().Error() != "alias conflict" {
					t.Errorf("expected error 'alias conflict', got %v", b.GetError())
				}
				if b.GetInput() != "id AS uid" {
					t.Errorf("expected source to be 'id AS uid', got %q", b.GetInput())
				}
				if !b.IsErrored() {
					t.Errorf("expected IsErrored() to return true after SetError")
				}
			})
		})

		t.Run("GetInput", func(t *testing.T) {
			t.Run("NilToken", func(t *testing.T) {
				var b *token.BaseToken
				if b.GetInput() != "" {
					t.Errorf("expected empty string from nil receiver, got %q", b.GetInput())
				}
			})

			t.Run("ValidToken", func(t *testing.T) {
				b := token.NewBaseToken("users.id")
				if b.GetInput() != "users.id" {
					t.Errorf("expected 'users.id', got %q", b.GetInput())
				}
			})

			t.Run("InvalidSource", func(t *testing.T) {
				b := token.NewBaseToken("")
				if b.GetInput() != "" {
					t.Errorf("expected empty string from token without source, got %q", b.GetInput())
				}
			})
		})

		t.Run("GetName", func(t *testing.T) {
			t.Run("NilToken", func(t *testing.T) {
				var b *token.BaseToken
				if b.GetName() != "" {
					t.Errorf("expected empty string from nil receiver, got %q", b.GetName())
				}
			})

			t.Run("ValidToken", func(t *testing.T) {
				b := token.NewBaseToken("id")
				if b.GetName() != "id" {
					t.Errorf("expected 'id', got %q", b.GetName())
				}
			})

			t.Run("InvalidTokenName", func(t *testing.T) {
				b := token.NewBaseToken("")
				if b.GetName() != "" {
					t.Errorf("expected empty string from token without name, got %q", b.GetName())
				}
			})
		})

		t.Run("GetRaw", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var b *token.BaseToken = nil
				b.SetError("ignored", fmt.Errorf("structural error"))
				if got := b.GetRaw(); got != "" {
					t.Errorf("expected nil error, got %v", got)
				}
			})

			t.Run("WithName", func(t *testing.T) {
				base := token.NewBaseToken("name")
				if got := base.GetRaw(); got != "name" {
					t.Errorf("expected 'name', got %q", got)
				}
			})

			t.Run("Aliased", func(t *testing.T) {
				raw := token.NewBaseToken("name", "alias")
				if got := raw.GetRaw(); got != "name AS alias" {
					t.Errorf("expected 'name AS alias', got %q", got)
				}
			})
		})

		t.Run("HasError", func(t *testing.T) {
			input := "AS"
			b := token.NewBaseToken(input)
			b.SetError(input, fmt.Errorf("error"))
			errored := b.HasError()
			if b == nil && !errored {
				t.Errorf("expected no panic and no error on nil receiver")
			}
		})

		t.Run("IsErrored", func(t *testing.T) {
			input := "AS"
			b := token.NewBaseToken(input)
			b.SetError(input, fmt.Errorf("error"))
			if b == nil && !b.IsErrored() {
				t.Errorf("expected no panic and no error on nil receiver")
			}
		})

		t.Run("IsValid", func(t *testing.T) {
			t.Run("ValidToken", func(t *testing.T) {
				b := token.NewBaseToken("id")
				if !b.IsValid() {
					t.Error("expected token to be valid")
				}
			})

			t.Run("NilToken", func(t *testing.T) {
				var b *token.BaseToken
				if b.IsValid() {
					t.Error("expected nil token to be invalid")
				}
			})

			t.Run("EmptyName", func(t *testing.T) {
				b := token.NewBaseToken("")
				if b.IsValid() {
					t.Error("expected token with empty name to be invalid")
				}
			})

			t.Run("WhitespaceName", func(t *testing.T) {
				b := token.NewBaseToken("   ")
				if b.IsValid() {
					t.Error("expected token with whitespace name to be invalid")
				}
			})

			t.Run("WithError", func(t *testing.T) {
				b := token.NewBaseToken("id")
				b.SetError("id", fmt.Errorf("something failed"))
				if b.IsValid() {
					t.Error("expected token with error to be invalid")
				}
			})
		})

		t.Run("Raw", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var b *token.BaseToken = nil
				b.SetError("ignored", fmt.Errorf("structural error"))
				if got := b.Raw(); got != "" {
					t.Errorf("expected nil error, got %v", got)
				}
			})
		})

		t.Run("RenderAlias", func(t *testing.T) {
			t.Run("NilToken", func(t *testing.T) {
				var b *token.BaseToken
				if got := b.RenderAlias(driver.NewPostgresDialect(), "id"); got != "id" {
					t.Errorf("expected fallback on nil token, got %q", got)
				}
			})

			t.Run("EmptyAlias", func(t *testing.T) {
				b := token.NewBaseToken("id")
				got := b.RenderAlias(driver.NewPostgresDialect(), "id")
				if got != "id" {
					t.Errorf("expected no alias fallback, got %q", got)
				}
			})

			t.Run("EmptyQualified", func(t *testing.T) {
				b := token.NewBaseToken("", "alias")
				got := b.RenderAlias(driver.NewPostgresDialect(), "")
				if got != "" {
					t.Errorf("expected empty string, got %q", got)
				}
			})

			t.Run("NilDialect", func(t *testing.T) {
				b := token.NewBaseToken("id", "plain_alias")
				got := b.RenderAlias(nil, "id")
				if got != "id AS plain_alias" {
					t.Errorf("expected unquoted alias fallback, got %q", got)
				}
			})
		})

		t.Run("RenderName", func(t *testing.T) {
			t.Run("ValidToken", func(t *testing.T) {
				b := token.NewBaseToken("id")
				if got := b.RenderName(driver.NewPostgresDialect()); got != `"id"` {
					t.Errorf("expected quoted name, got %q", got)
				}
			})

			t.Run("NilToken", func(t *testing.T) {
				var b *token.BaseToken
				if got := b.RenderName(driver.NewPostgresDialect()); got != "" {
					t.Errorf("expected empty string, got %q", got)
				}
			})

			t.Run("EmptyName", func(t *testing.T) {
				b := token.NewBaseToken("")
				if got := b.RenderName(driver.NewPostgresDialect()); got != "" {
					t.Errorf("expected empty string, got %q", got)
				}
			})

			t.Run("NilDialect", func(t *testing.T) {
				b := token.NewBaseToken("id")
				if got := b.RenderName(nil); got != "id" {
					t.Errorf("expected unquoted name, got %q", got)
				}
			})
		})

		t.Run("SetAlias", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var b *token.BaseToken = nil
				b.SetAlias("uid")
				if name := b.GetName(); name != "" {
					t.Errorf("expected name %q, got %q", "", name)
				}
			})

			t.Run("WithValue", func(t *testing.T) {
				b := token.NewBaseToken("id")
				b.SetAlias("uid")
				if name := b.GetAlias(); name != "uid" {
					t.Errorf("expected name %q, got %q", "uid", name)
				}
			})
		})

		t.Run("SetError", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var b *token.BaseToken = nil
				b.SetError("ignored", fmt.Errorf("structural error"))
				if got := b.GetError(); got != nil {
					t.Errorf("expected nil error, got %v", got)
				}
			})

			t.Run("WithInput", func(t *testing.T) {
				b := token.NewBaseToken("is", "uid")
				err := fmt.Errorf("alias conflict")
				b.SetError("id AS uid", err)

				if b.GetError() == nil || b.GetError().Error() != "alias conflict" {
					t.Errorf("expected error 'alias conflict', got %v", b.GetError())
				}
				if b.GetInput() != "id AS uid" {
					t.Errorf("expected source to be 'id AS uid', got %q", b.GetInput())
				}
			})

			t.Run("TokenErrored", func(t *testing.T) {
				b := token.NewBaseToken("users.id")
				b.SetError("ignored", fmt.Errorf("structural error"))

				if b.GetInput() != "ignored" {
					t.Errorf("expected source to remain 'ignored', got %q", b.GetInput())
				}
			})
		})

		t.Run("SetErrorWith", func(t *testing.T) {
			t.Run("input", func(t *testing.T) {
				b := token.NewBaseToken("id", "uid")
				err := fmt.Errorf("alias conflict")
				b.SetErrorWith("id AS uid", err)

				if b.GetError() == nil || b.GetError().Error() != "alias conflict" {
					t.Errorf("expected error 'alias conflict', got %v", b.GetError())
				}
				if b.GetInput() != "id AS uid" {
					t.Errorf("expected source to be 'id AS uid', got %q", b.GetInput())
				}
			})
		})

		t.Run("SetKind", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var b *token.BaseToken = nil
				b.SetKind(contract.ColumnKind)
				if got := b.GetKind(); got != contract.UnknownKind {
					t.Errorf("expected UnknownKind for nil receiver, got %v", got)
				}
			})

			t.Run("UnknownKind", func(t *testing.T) {
				b := token.NewBaseToken("")
				if got := b.GetKind(); got != contract.UnknownKind {
					t.Errorf("expected UnknownKind, got %v", got)
				}
			})

			t.Run("Full", func(t *testing.T) {
				b := token.NewBaseToken("users", "u")
				b.SetKind(contract.ColumnKind)
				if got := b.GetKind(); got != contract.ColumnKind {
					t.Errorf("expected ColumnKind, got %v", got)
				}
			})
		})

		t.Run("SetName", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var b *token.BaseToken = nil
				b.SetName("id")
				if name := b.GetName(); name != "" {
					t.Errorf("expected name %q, got %q", "id", name)
				}
			})

			t.Run("WithValue", func(t *testing.T) {
				b := token.NewBaseToken("id")
				b.SetName("uid")
				if name := b.GetName(); name != "uid" {
					t.Errorf("expected name %q, got %q", "id", name)
				}
			})
		})

		t.Run("String", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var b *token.BaseToken = nil
				_ = b.String() // just ensure no panic
			})

			t.Run("WithColumn", func(t *testing.T) {
				b := token.NewBaseToken("id")
				want := `Unknown("id") [aliased: false, errored: false]`
				if got := b.String(); got != want {
					t.Errorf("got %q, want %q", got, want)
				}
			})

			t.Run("WithAliasedColumn", func(t *testing.T) {
				b := token.NewBaseToken("id", "user_id")
				b.SetKind(contract.ColumnKind)
				want := `Column("id") [aliased: true, errored: false]`
				if got := b.String(); got != want {
					t.Errorf("got %q, want %q", got, want)
				}
			})

			t.Run("WithAliasedTable", func(t *testing.T) {
				b := token.NewBaseToken("users", "u")
				b.SetKind(contract.TableKind)
				want := `Table("users") [aliased: true, errored: false]`
				if got := b.String(); got != want {
					t.Errorf("got %q, want %q", got, want)
				}
			})

			t.Run("WithError", func(t *testing.T) {
				b := token.NewBaseToken("id AS uid", "user_id")
				out := b.String()
				if !strings.Contains(out, "errored: true") || !strings.Contains(out, "error: alias conflict") {
					t.Errorf("expected error details, got %q", out)
				}
			})
		})
	})

	t.Run("Validations", func(t *testing.T) {
		t.Run("EmptyInput", func(t *testing.T) {
			b := token.NewBaseToken("")
			if b.GetError() == nil {
				t.Errorf("expected error for empty input, got nil")
			} else if b.GetError().Error() != "invalid input expression: expression is empty" {
				t.Errorf("unexpected error message: %s", b.GetError())
			}
			if b.GetInput() != "" {
				t.Errorf("expected input to be empty, got %q", b.GetInput())
			}
		})

		t.Run("InvalidInput", func(t *testing.T) {
			b := token.NewBaseToken("id, name")
			if b.GetError() == nil {
				t.Errorf("expected error, got nil")
			} else if b.GetError().Error() != "invalid input expression: aliases must not be comma-separated" {
				t.Errorf("unexpected error message: %s", b.GetError())
			}
		})

		t.Run("InputStartsWithAS", func(t *testing.T) {
			b := token.NewBaseToken("AS uid")
			if b.GetError() == nil {
				t.Errorf("expected error, got nil")
			} else if b.GetError().Error() != "invalid input expression: cannot start with 'AS'" {
				t.Errorf("unexpected error message: %s", b.GetError())
			}
		})

		t.Run("ReservedWordASOnly", func(t *testing.T) {
			b := token.NewBaseToken("AS")
			if b.GetError() == nil {
				t.Errorf("expected error for reserved word 'AS', got nil")
			} else if b.GetError().Error() != "invalid input expression: name cannot be AS keyword" {
				t.Errorf("unexpected error message: %s", b.GetError())
			}
		})
	})
}
