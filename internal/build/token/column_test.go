// File: internal/build/token/column_test.go

package token_test

import (
	"strings"
	"testing"

	"github.com/entiqon/entiqon/driver"
	"github.com/entiqon/entiqon/internal/build/token"
)

func TestColumn(t *testing.T) {
	t.Run("Usage", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			col := token.NewColumn("id")
			if col.Name != "id" {
				t.Errorf("expected column name 'id', got %q", col.Name)
			}
			if col.IsAliased() {
				t.Errorf("expected column to not be aliased")
			}
			if col.IsQualified() {
				t.Errorf("expected column to not be qualified")
			}
			if !col.IsValid() {
				t.Errorf("expected column to be valid")
			}
		})

		t.Run("Table", func(t *testing.T) {
			t.Run("Inline", func(t *testing.T) {
				col := token.NewColumn("users.id")
				if col.Table == nil || col.Table.Name != "users" {
					t.Errorf("expected table name to be 'users', got %v", col.Table)
				}
				if col.HasError() {
					t.Errorf("unexpected error: %v", col.Error)
				}
			})
		})

		t.Run("Alias", func(t *testing.T) {
			t.Run("Inline", func(t *testing.T) {
				col := token.NewColumn("id AS email")
				if !col.IsValid() {
					t.Fatalf("expected column to be valid, but got error: %v", col.Error)
				}
				if col.Name != "id" {
					t.Errorf("expected column name 'id', got %q", col.Name)
				}
				if col.Alias != "email" {
					t.Errorf("expected alias 'email', got %q", col.Alias)
				}
			})

			t.Run("ExplicitAlias", func(t *testing.T) {
				col := token.NewColumn("id", "email")
				if !col.IsValid() {
					t.Fatalf("expected column to be valid, but got error: %v", col.Error)
				}
				if col.Name != "id" {
					t.Errorf("expected column name 'id', got %q", col.Name)
				}
				if col.Alias != "email" {
					t.Errorf("expected alias 'email', got %q", col.Alias)
				}
			})

			t.Run("PostgresEmptyName", func(t *testing.T) {
				col := token.NewColumn("'' AS name")
				if !col.IsValid() {
					t.Fatalf("expected column to be valid, but got error: %v", col.Error)
				}
				if col.Name != "''" {
					t.Errorf("expected column name \"''\", got %q", col.Name)
				}
				if col.Alias != "name" {
					t.Errorf("expected alias 'name', got %q", col.Alias)
				}
			})

			t.Run("AliasOverriding", func(t *testing.T) {
				t.Run("Match", func(t *testing.T) {
					col := token.NewColumn("id AS uid", "uid")
					if !col.IsAliased() {
						t.Errorf("expected column to be aliased")
					}
					if !col.IsValid() {
						t.Errorf("expected column to be valid, got error: %v", col.Error)
					}
				})

				t.Run("AliasConflict", func(t *testing.T) {
					col := token.NewColumn("user_id AS internal", "external")
					if col.IsValid() {
						t.Fatalf("expected column to be invalid due to alias conflict")
					}
					if col.Error == nil {
						t.Fatalf("expected error due to alias conflict, got nil")
					}
					if !strings.Contains(col.Error.Error(), "alias conflict") {
						t.Errorf("expected error message to contain 'alias conflict', got: %v", col.Error)
					}
					if col.Name != "user_id" {
						t.Errorf("expected column name 'user_id', got %q", col.Name)
					}
					if col.Alias != "external" {
						t.Errorf("expected alias 'external', got %q", col.Alias)
					}
					if col.Source != "user_id AS internal" {
						t.Errorf("expected source 'user_id AS internal', got %q", col.Source)
					}
				})
			})

			t.Run("ParsingErrors", func(t *testing.T) {
				t.Run("StartsWithAS", func(t *testing.T) {
					col := token.NewColumn("AS email")
					if col.IsValid() {
						t.Fatalf("expected column to be invalid due to starting with 'AS'")
					}
					if col.Error == nil || !strings.Contains(col.Error.Error(), "cannot start with 'AS'") {
						t.Errorf("expected error containing 'cannot start with 'AS'', got: %v", col.Error)
					}
					if col.Source != "AS email" {
						t.Errorf("expected source to be 'AS email', got %q", col.Source)
					}
				})

				t.Run("OnlyAlias", func(t *testing.T) {
					col := token.NewColumn(" AS alias")
					if col.IsValid() {
						t.Fatalf("expected column to be invalid due to starting with 'AS'")
					}
					if col.Error == nil || !strings.Contains(col.Error.Error(), "cannot start with 'AS'") {
						t.Errorf("expected error containing 'cannot start with 'AS'', got: %v", col.Error)
					}
					if col.Source != " AS alias" {
						t.Errorf("expected source to be ' AS alias', got %q", col.Source)
					}
				})

				t.Run("DotWithoutColumn", func(t *testing.T) {
					col := token.NewColumn("u.")
					if col.IsValid() {
						t.Fatalf("expected column to be invalid due to missing column name after dot")
					}
					if col.Error == nil || !strings.Contains(col.Error.Error(), "column name is required") {
						t.Errorf("expected error containing 'column name is required', got: %v", col.Error)
					}
					if col.Source != "u." {
						t.Errorf("expected source to be 'u.', got %q", col.Source)
					}
				})
			})
		})

		t.Run("Invalid", func(t *testing.T) {
			col := token.NewColumn("id, uid")
			if col.IsValid() {
				t.Errorf("expected column to be invalid due to comma-separated input")
			}
			if col.Error == nil {
				t.Errorf("expected error due to comma-separated alias, got nil")
			} else if !strings.Contains(col.Error.Error(), "aliases must not be comma-separated") {
				t.Errorf(
					"unexpected error message: got %q, want message containing %q",
					col.Error.Error(),
					"aliases must not be comma-separated",
				)
			}
		})
	})

	t.Run("WithTable", func(t *testing.T) {
		t.Run("Explicit", func(t *testing.T) {
			col := token.NewColumn("id", "user_id").WithTable(token.NewTable("users"))
			if col.Table == nil || col.Table.Name != "users" {
				t.Errorf("expected table name to be 'users', got %v", col.Table)
			}
		})

		t.Run("Conflict", func(t *testing.T) {
			col := token.NewColumn("users.id").WithTable(token.NewTable("orders")) // should produce error
			if col.Table == nil || col.Table.Name != "users" {
				t.Errorf("expected table name to remain 'users', got %v", col.Table)
			}
			if !col.HasError() {
				t.Fatalf("expected error due to table mismatch, but got none")
			}
			if !strings.Contains(col.Error.Error(), "table mismatch") {
				t.Errorf("expected error to contain 'table mismatch', got: %v", col.Error)
			}
		})

		t.Run("InvalidColumn", func(t *testing.T) {
			col := token.NewColumn("") // empty = invalid
			if col.IsValid() {
				t.Errorf("expected column to be invalid")
			}
			result := col.WithTable(token.NewTable("users"))
			if result != col {
				t.Errorf("expected WithTable to return original column")
			}
		})

		t.Run("NilColumn", func(t *testing.T) {
			var col *token.Column
			result := col.WithTable(token.NewTable("users"))
			if result != nil {
				t.Errorf("expected nil result for nil column")
			}
		})
	})

	t.Run("Raw", func(t *testing.T) {
		col := token.NewColumn("id")
		if col.HasError() {
			t.Errorf("unexpected error: %v", col.Error)
		}
		if got := col.Raw(); got != "id" {
			t.Errorf("Raw mismatch: got %q, want %q", got, "id")
		}

		t.Run("WithTable", func(t *testing.T) {
			col := token.NewColumn("u.id AS user_id").WithTable(token.NewTable("users u"))

			if got := col.Source; got != "u.id AS user_id" {
				t.Errorf("Source mismatch: got %q, want %q", got, "id")
			}
			if got := col.Raw(); got != "u.user_id" {
				t.Errorf("Raw mismatch: got %q, want %q", got, "u.id AS user_id")
			}
			if got := col.String(); got != `Column("id") [aliased: true, qualified: true, errored: false]` {
				t.Errorf("String mismatch: got %q, want %q", got, "u.user_id")
			}
			if col.Alias != "user_id" {
				t.Errorf("Alias mismatch: got %q, want %q", col.Alias, "user_id")
			}
			if !col.IsAliased() {
				t.Errorf("expected column to be aliased")
			}
			if !col.IsQualified() {
				t.Errorf("expected column to be qualified")
			}
			if col.HasError() {
				t.Errorf("unexpected error: %v", col.Error)
			}
		})

		t.Run("IsQualified", func(t *testing.T) {
			col := token.NewColumn("users.id")
			if col.IsAliased() {
				t.Errorf("expected column to not be aliased")
			}
			if !col.IsQualified() {
				t.Errorf("expected column to be qualified")
			}
			if got := col.Raw(); got != "users.id" {
				t.Errorf("Raw mismatch: got %q, want %q", got, "users.id")
			}
		})

		t.Run("EdgeCases", func(t *testing.T) {
			t.Run("IsAliased", func(t *testing.T) {
				col := token.NewColumn("id AS user_id")
				if got := col.Raw(); got != "id AS user_id" {
					t.Errorf("Raw mismatch: got %q, want %q", got, "id AS user_id")
				}
			})

			t.Run("NotTableIsAliased", func(t *testing.T) {
				table := token.NewTable("users") // no alias
				col := token.NewColumn("users.id AS user_id").WithTable(table)
				if got := col.Raw(); got != "users.user_id" {
					t.Errorf("Raw mismatch: got %q, want %q", got, "users.user_id")
				}
			})
		})
	})

	t.Run("Render", func(t *testing.T) {
		postgres := driver.NewPostgresDialect()

		t.Run("Basic", func(t *testing.T) {
			col := token.NewColumn("email")
			got := col.Render(postgres)
			want := `"email"`

			if got != want {
				t.Errorf("Render mismatch: got %q, want %q", got, want)
			}
		})

		t.Run("Aliased", func(t *testing.T) {
			col := token.NewColumn("email AS user_email")
			got := col.Render(postgres)
			want := `"email" AS "user_email"`

			if got != want {
				t.Errorf("Render mismatch: got %q, want %q", got, want)
			}
		})

		t.Run("Qualified", func(t *testing.T) {
			col := token.NewColumn("users.email").WithTable(token.NewTable("users"))
			got := col.Render(postgres)
			want := `"users"."email"`

			if got != want {
				t.Errorf("Render mismatch: got %q, want %q", got, want)
			}
		})

		t.Run("Full", func(t *testing.T) {
			col := token.NewColumn("u.email AS mail").WithTable(token.NewTable("users u"))
			got := col.Render(postgres)
			want := `"u"."email" AS "mail"`

			if got != want {
				t.Errorf("Render mismatch: got %q, want %q", got, want)
			}
		})

		t.Run("NoTableAliased", func(t *testing.T) {
			col := token.NewColumn("users.email AS mail").WithTable(token.NewTable("users"))
			got := col.Render(postgres)
			want := `"users"."email" AS "mail"`

			if got != want {
				t.Errorf("Render mismatch: got %q, want %q", got, want)
			}
		})

		t.Run("InvalidColumn", func(t *testing.T) {
			col := token.NewColumn("") // invalid
			got := col.Render(postgres)

			if got != "" {
				t.Errorf("Expected empty render for invalid column, got %q", got)
			}
		})
	})

	t.Run("String", func(t *testing.T) {
		col := token.NewColumn("users.id", "uid")
		want := `Column("id") [aliased: true, qualified: true, errored: false]`
		if got := col.String(); got != want {
			t.Errorf("String mismatch: got %q, want %q", got, want)
		}

		t.Run("WithError", func(t *testing.T) {
			col := token.NewColumn("users.id").WithTable(token.NewTable("orders"))
			if !col.HasError() {
				t.Errorf("expected column to have error due to table mismatch")
			}
			if got := col.String(); !strings.Contains(got, "table mismatch") {
				t.Errorf("expected String to contain 'table mismatch', got %q", got)
			}
		})
	})
}
