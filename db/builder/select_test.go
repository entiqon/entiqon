// File: db/builder/select_test.go

package builder_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/builder"
	"github.com/entiqon/entiqon/db/token"
)

func TestSelectBuilder(t *testing.T) {
	t.Run("Methods", func(t *testing.T) {
		t.Run("Fields", func(t *testing.T) {
			t.Run("NilCollection", func(t *testing.T) {
				sb := &builder.SelectBuilder{} // fields is nil
				sb.Fields("id")
				fields := sb.GetFields()
				if fields.Length() != 1 || fields[0].Expr != "id" {
					t.Errorf("expected one field 'id', got %+v", fields)
				}
			})

			t.Run("NoArgs", func(t *testing.T) {
				sb := builder.NewSelect(nil).Fields()
				fields := sb.GetFields()
				if fields.Length() != 0 {
					t.Errorf("expected no fields, got %d", len(fields))
				}
			})

			t.Run("Add", func(t *testing.T) {
				sb := builder.NewSelect(nil).Fields("id")
				fields := sb.GetFields()
				if fields.Length() != 1 || fields[0].Expr != "id" {
					t.Errorf("expected reset only, got %+v", fields)
				}
			})

			t.Run("Reset", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Fields("reset") // should reset
				fields := sb.GetFields()
				if fields.Length() != 1 || fields[0].Expr != "reset" {
					t.Errorf("expected reset only, got %+v", fields)
				}
			})

			t.Run("AliasedField", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("expr", "alias")
				fields := sb.GetFields()
				if fields.Length() != 1 {
					t.Errorf("expected 1 field, got %d", fields.Length())
				}
			})

			t.Run("Pointer", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields(token.NewField("id")) // *token.Field
				fields := sb.GetFields()
				if fields.Length() != 1 {
					t.Errorf("expected 1 field, got %d", len(fields))
				}
			})

			t.Run("NotPointer", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields(*token.NewField("id")) // token.Field (value)
				fields := sb.GetFields()
				if fields.Length() != 1 {
					t.Errorf("expected 1 field, got %d", len(fields))
				}
			})

			t.Run("Invalid", func(t *testing.T) {
				sb := builder.NewSelect(nil).Fields(true)
				fields := sb.GetFields()
				if !fields[0].IsErrored() {
					t.Errorf("expected IsErrored to be true, got %v", fields[0].IsErrored())
				}
				if fields[0].Error == nil {
					t.Errorf("expected Error to be set, got nil")
				}
			})
		})

		t.Run("AddFields", func(t *testing.T) {
			sb := builder.NewSelect(nil).
				Fields("id").
				AddFields("name").
				AddFields("email")

			fields := sb.GetFields()
			if len(fields) != 3 {
				t.Errorf("expected 3 fields, got %d", len(fields))
			}
		})

		t.Run("Limit", func(t *testing.T) {
			sb := builder.NewSelect(nil).
				Fields("id").
				Source("users").
				Limit(10)
			sql, err := sb.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := "LIMIT 10"
			if !strings.Contains(sql, want) {
				t.Errorf("expected %q, got %q", want, sql)
			}
		})

		t.Run("Offset", func(t *testing.T) {
			sb := builder.NewSelect(nil).
				Fields("id").
				Source("users").
				Offset(20)
			sql, err := sb.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := "OFFSET 20"
			if !strings.Contains(sql, want) {
				t.Errorf("expected %q, got %q", want, sql)
			}
		})

		t.Run("Pagination", func(t *testing.T) {
			sb := builder.NewSelect(nil).
				Fields("id").
				Source("users").
				Limit(10).
				Offset(20)
			sql, err := sb.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := "LIMIT 10 OFFSET 20"
			if !strings.Contains(sql, want) {
				t.Errorf("expected %q, got %q", want, sql)
			}
		})

		t.Run("String", func(t *testing.T) {
			sb := builder.NewSelect(nil).
				Fields("id").
				Source("users")

			want := `Status ✅: SQL=SELECT id FROM users, Params=`
			got := sb.String()
			if got != want {
				t.Errorf("String() = %q, want %q", got, want)
			}

			// Test error case when no source
			sb = builder.NewSelect(nil).
				Fields("id")

			wantPrefix := "Status ❌: Error building SQL"
			got = sb.String()
			if !strings.HasPrefix(got, wantPrefix) {
				t.Errorf("String() error output = %q, want prefix %q", got, wantPrefix)
			}
		})

		t.Run("Build", func(t *testing.T) {
			t.Run("NilReceiver", func(t *testing.T) {
				var sb *builder.SelectBuilder = nil // nil receiver

				_, err := sb.Build()
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				want := "❌ [Build] - Wrong initialization. Cannot build on receiver nil"
				if err.Error() != want {
					t.Errorf("unexpected error: got %q, want %q", err.Error(), want)
				}
			})

			t.Run("EmptyFields", func(t *testing.T) {
				// Empty string should not add a field → defaults to SELECT *
				sql, _ := builder.NewSelect(nil).
					Source("users").
					Fields("").
					Build()
				if !strings.Contains(sql, "SELECT *") {
					t.Errorf("expected to contain '*', got %v", sql)
				}
			})

			t.Run("WithFields", func(t *testing.T) {
				sql, _ := builder.NewSelect(nil).
					Fields("id, name").
					Source("users").
					Build()

				if !strings.Contains(sql, "SELECT id, name") {
					t.Errorf("expected SQL to contain 'SELECT id, name', got %q", sql)
				}
			})

			t.Run("InvalidFields", func(t *testing.T) {
				_, err := builder.NewSelect(nil).
					Fields(true).     // rejected
					AddFields(false). // rejected
					AddFields(123).   // rejected
					Source("users").
					Build()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				out := err.Error()
				if !strings.Contains(out, "❌ [Build] - Invalid fields:") {
					t.Errorf("expected error to contain '❌ [Build] - Invalid fields:', got %q", out)
				}
				if !strings.Contains(out, "⛔️ Field(") {
					t.Errorf("expected detailed field errors, got %q", out)
				}
			})

			t.Run("Limit", func(t *testing.T) {
				sql, _ := builder.NewSelect(nil).
					Source("users").
					Limit(10).
					Build()
				if !strings.Contains(sql, "LIMIT 10") {
					t.Errorf("expected contains LIMIT 10, got %q", sql)
				}
			})

			t.Run("Offset", func(t *testing.T) {
				sql, _ := builder.NewSelect(nil).
					Source("users").
					Offset(20).
					Build()
				if !strings.Contains(sql, "OFFSET 20") {
					t.Errorf("expected contains OFFSET 20, got %q", sql)
				}
			})

			t.Run("Pagination", func(t *testing.T) {
				sql, _ := builder.NewSelect(nil).
					Source("users").
					Limit(10).
					Offset(20).
					Build()
				want := fmt.Sprintf("LIMIT %d OFFSET %d", 10, 20)
				if !strings.Contains(sql, want) {
					t.Errorf("expected contains %q, got %q", want, sql)
				}
			})

			t.Run("SQL", func(t *testing.T) {
				sql, _ := builder.NewSelect(nil).
					Fields("id, name").
					Source("users").
					Limit(10).
					Offset(20).
					Build()
				want := "SELECT id, name FROM users LIMIT 10 OFFSET 20"
				if sql != want {
					t.Errorf("expected %s, got %s", want, sql)
				}
			})

			t.Run("IgnoreUnsupportedTypes_All", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields(123, 4.5, true, struct{}{}). // all unsupported → ignored
					Source("users")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				// With no supported fields added, builder should default to "*".
				if !strings.Contains(sql, "SELECT * ") {
					t.Errorf("expected wildcard SELECT *, got %q", sql)
				}
			})
		})
	})
}
