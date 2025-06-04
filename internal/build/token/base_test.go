// File: internal/build/token/token_test.go

package token_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/driver"
	"github.com/entiqon/entiqon/internal/build/token"
)

func TestBaseToken(t *testing.T) {
	t.Run("Usage", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			b := &token.BaseToken{Alias: "user_id"}
			got := b.RenderAlias(driver.NewPostgresDialect(), "u.id")
			if got != `u.id AS "user_id"` {
				t.Errorf("expected quoted alias, got %q", got)
			}
		})

		t.Run("Members", func(t *testing.T) {
			t.Run("AliasOr", func(t *testing.T) {
				t.Run("Name", func(t *testing.T) {
					tok := token.BaseToken{Name: "id"}
					if got := tok.AliasOr(); got != "id" {
						t.Errorf("expected 'id', got %q", got)
					}
				})

				t.Run("Aliased", func(t *testing.T) {
					tok := token.BaseToken{Name: "id", Alias: "uid"}
					if got := tok.AliasOr(); got != "uid" {
						t.Errorf("expected 'uid', got %q", got)
					}
				})
			})

			t.Run("IsValid", func(t *testing.T) {
				t.Run("ValidToken", func(t *testing.T) {
					b := &token.BaseToken{Name: "id"}
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
					b := &token.BaseToken{}
					if b.IsValid() {
						t.Error("expected token with empty name to be invalid")
					}
				})

				t.Run("WhitespaceName", func(t *testing.T) {
					b := &token.BaseToken{Name: "   "}
					if b.IsValid() {
						t.Error("expected token with whitespace name to be invalid")
					}
				})

				t.Run("WithError", func(t *testing.T) {
					b := &token.BaseToken{Name: "id", Error: fmt.Errorf("something failed")}
					if b.IsValid() {
						t.Error("expected token with error to be invalid")
					}
				})
			})

			t.Run("Raw", func(t *testing.T) {
				t.Run("Name", func(t *testing.T) {
					base := &token.BaseToken{Name: "name"}
					if got := base.Raw(); got != "name" {
						t.Errorf("expected 'name', got %q", got)
					}
				})

				t.Run("Aliased", func(t *testing.T) {
					raw := token.BaseToken{Name: "name", Alias: "alias"}
					if got := raw.Raw(); got != "name AS alias" {
						t.Errorf("expected 'name AS alias', got %q", got)
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
					b := &token.BaseToken{Name: "id"}
					got := b.RenderAlias(driver.NewPostgresDialect(), "id")
					if got != "id" {
						t.Errorf("expected no alias fallback, got %q", got)
					}
				})

				t.Run("EmptyQualified", func(t *testing.T) {
					b := &token.BaseToken{Alias: "alias"}
					got := b.RenderAlias(driver.NewPostgresDialect(), "")
					if got != "" {
						t.Errorf("expected empty string, got %q", got)
					}
				})

				t.Run("NilDialect", func(t *testing.T) {
					b := &token.BaseToken{Alias: "plain_alias"}
					got := b.RenderAlias(nil, "id")
					if got != "id AS plain_alias" {
						t.Errorf("expected unquoted alias fallback, got %q", got)
					}
				})
			})

			t.Run("RenderName", func(t *testing.T) {
				t.Run("ValidToken", func(t *testing.T) {
					b := &token.BaseToken{Name: "id"}
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
					b := &token.BaseToken{}
					if got := b.RenderName(driver.NewPostgresDialect()); got != "" {
						t.Errorf("expected empty string, got %q", got)
					}
				})

				t.Run("NilDialect", func(t *testing.T) {
					b := &token.BaseToken{Name: "id"}
					if got := b.RenderName(nil); got != "id" {
						t.Errorf("expected unquoted name, got %q", got)
					}
				})
			})

			t.Run("SetErrorWith", func(t *testing.T) {
				t.Run("Source", func(t *testing.T) {
					b := &token.BaseToken{Name: "id", Alias: "uid"}
					err := fmt.Errorf("alias conflict")
					updated := b.SetErrorWith("id AS uid", err)

					if updated != b {
						t.Error("SetErrorWith should return the same token instance")
					}
					if b.Error == nil || b.Error.Error() != "alias conflict" {
						t.Errorf("expected error 'alias conflict', got %v", b.Error)
					}
					if b.Source != "id AS uid" {
						t.Errorf("expected source to be 'id AS uid', got %q", b.Source)
					}
				})

				t.Run("TokenErrored", func(t *testing.T) {
					b := &token.BaseToken{Source: "users.id"}
					b.SetErrorWith("ignored", fmt.Errorf("structural error"))

					if b.Source != "users.id" {
						t.Errorf("expected source to remain 'users.id', got %q", b.Source)
					}
				})
			})

			t.Run("String", func(t *testing.T) {
				t.Run("WithColumn", func(t *testing.T) {
					b := &token.BaseToken{Name: "id"}
					want := `Column("id") [aliased: false, errored: false]`
					if got := b.String(token.KindColumn); got != want {
						t.Errorf("got %q, want %q", got, want)
					}
				})

				t.Run("WithAliasedColumn", func(t *testing.T) {
					b := &token.BaseToken{Name: "id", Alias: "user_id"}
					want := `Column("id") [aliased: true, errored: false]`
					if got := b.String(token.KindColumn); got != want {
						t.Errorf("got %q, want %q", got, want)
					}
				})

				t.Run("WithAliasedTable", func(t *testing.T) {
					b := &token.BaseToken{Name: "users", Alias: "u"}
					want := `Table("users") [aliased: true, errored: false]`
					if got := b.String(token.KindTable); got != want {
						t.Errorf("got %q, want %q", got, want)
					}
				})

				t.Run("WithError", func(t *testing.T) {
					b := &token.BaseToken{
						Name:  "id",
						Alias: "uid",
						Error: fmt.Errorf("conflict"),
					}
					out := b.String(token.KindColumn)
					if !strings.Contains(out, "errored: true") || !strings.Contains(out, "error: conflict") {
						t.Errorf("expected error details, got %q", out)
					}
				})
			})
		})
	})
}
