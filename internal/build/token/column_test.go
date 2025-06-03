// File: internal/build/token/column_test.go

package token_test

import (
	"strings"
	"testing"

	"github.com/entiqon/entiqon/internal/build/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColumn(t *testing.T) {
	t.Run("Usage", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			col := token.NewColumn("id")
			assert.Equal(t, "id", col.Name)
			assert.False(t, col.IsAliased())
			assert.False(t, col.IsQualified())
			assert.True(t, col.IsValid())
		})

		t.Run("Table", func(t *testing.T) {
			t.Run("Inline", func(t *testing.T) {
				col := token.NewColumn("users.id")
				assert.Equal(t, "users", col.Table.Name)
				assert.False(t, col.HasError())
			})
		})

		t.Run("Alias", func(t *testing.T) {
			t.Run("Inline", func(t *testing.T) {
				col := token.NewColumn("id AS email")
				require.True(t, col.IsValid())
				assert.Equal(t, "id", col.Name)
				assert.Equal(t, "email", col.Alias)
			})

			t.Run("ExplicitAlias", func(t *testing.T) {
				col := token.NewColumn("id", "email")
				require.True(t, col.IsValid())
				assert.Equal(t, "id", col.Name)
				assert.Equal(t, "email", col.Alias)
			})

			t.Run("PostgresEmptyName", func(t *testing.T) {
				col := token.NewColumn("'' AS name")
				require.True(t, col.IsValid())
				assert.Equal(t, "''", col.Name)
				assert.Equal(t, "name", col.Alias)
			})

			t.Run("AliasOverriding", func(t *testing.T) {
				t.Run("Match", func(t *testing.T) {
					col := token.NewColumn("id AS uid", "uid")
					assert.True(t, col.IsAliased())
					assert.True(t, col.IsValid())
				})

				t.Run("AliasConflict", func(t *testing.T) {
					col := token.NewColumn("user_id AS internal", "external")

					require.False(t, col.IsValid())
					require.Error(t, col.Error)
					require.Contains(t, col.Error.Error(), "alias conflict")

					assert.Equal(t, "user_id", col.Name)
					assert.Equal(t, "external", col.Alias) // resolved to explicit
					assert.Equal(t, "user_id AS internal", col.Source)
				})
			})

			t.Run("ParsingErrors", func(t *testing.T) {
				t.Run("StartsWithAS", func(t *testing.T) {
					col := token.NewColumn("AS email")
					assert.False(t, col.IsValid())
					assert.ErrorContains(t, col.Error, "cannot start with 'AS'")
					assert.Equal(t, "AS email", col.Source)
				})

				t.Run("OnlyAlias", func(t *testing.T) {
					col := token.NewColumn(" AS alias")
					assert.False(t, col.IsValid())
					assert.ErrorContains(t, col.Error, "cannot start with 'AS'")
					assert.Equal(t, " AS alias", col.Source)
				})

				t.Run("DotWithoutColumn", func(t *testing.T) {
					col := token.NewColumn("u.")
					assert.False(t, col.IsValid())
					assert.ErrorContains(t, col.Error, "column name is required")
					assert.Equal(t, "u.", col.Source)
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
			assert.Equal(t, "users", col.Table.Name)
		})

		t.Run("Conflict", func(t *testing.T) {
			col := token.NewColumn("users.id").WithTable(token.NewTable("orders")) // should produce error
			assert.Equal(t, "users", col.Table.Name)                               // should not change
			assert.True(t, col.HasError())
			assert.ErrorContains(t, col.Error, "table mismatch")
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

	t.Run("String", func(t *testing.T) {
		col := token.NewColumn("users.id", "uid")
		assert.Contains(t, col.String(), `Column("id") [aliased: true, qualified: true, errored: false`)

		t.Run("WithError", func(t *testing.T) {
			col := token.NewColumn("users.id").WithTable(token.NewTable("orders"))
			assert.True(t, col.HasError())
			assert.Contains(t, col.String(), "table mismatch")
		})
	})
}
