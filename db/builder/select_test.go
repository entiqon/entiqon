// File: db/builder/select_test.go

package builder_test

import (
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/builder"
	"github.com/entiqon/entiqon/db/driver"
)

func TestSelectBuilder(t *testing.T) {
	t.Run("Usage", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			sql, args, err := builder.NewSelect(nil).From("users").Build()
			if err != nil {
				t.Fatalf("unexpected error from Build(): %v", err)
			}

			// FROM clause should reference the table (no aliasing on columns)
			if !strings.Contains(sql, "FROM users") {
				t.Errorf("expected SQL to contain %q, got %q", "FROM users", sql)
			}

			if len(args) != 0 {
				t.Errorf("expected 1 args but got %d args", len(args))
			}
		})

		t.Run("Medium", func(t *testing.T) {
			sql, args, err := builder.NewSelect(nil).
				From("users").
				Where("id = 1").
				Build()
			if !strings.Contains(sql, "WHERE id = ?") {
				t.Errorf("expected sql contains %q, got %q", "WHERE", sql)
			}
			if len(args) != 1 {
				t.Errorf("expected 1 args but got %d args", len(args))
			}
			if err != nil {
				t.Errorf("expected err=nil but got %q", err)
			}
		})

		t.Run("Advanced", func(t *testing.T) {
			sb := builder.NewSelect(nil).
				From("users").
				Select("id AS uid", "username").
				AddSelect("email, password").
				Where("id = 1")

			sql, args, err := sb.Build()
			if err != nil {
				t.Fatalf("unexpected error from Build(): %v", err)
			}

			// Columns should appear exactly as rendered, unprefixed
			wantCols := "SELECT id AS uid, username, email, password"
			if !strings.Contains(sql, wantCols) {
				t.Errorf("expected SQL to contain %q, got %q", wantCols, sql)
			}

			// FROM clause should reference the table (no aliasing on columns)
			if !strings.Contains(sql, "FROM users") {
				t.Errorf("expected SQL to contain %q, got %q", "FROM users", sql)
			}

			// WHERE clause and parameter binding
			if !strings.Contains(sql, "WHERE id = ?") {
				t.Errorf("expected SQL to contain %q, got %q", "WHERE id = ?", sql)
			}

			if len(args) != 1 {
				t.Errorf("expected 1 args but got %d args", len(args))
			}
		})

		t.Run("AdvancedWithAliasedTable", func(t *testing.T) {
			sql, args, err := builder.NewSelect(nil).
				From("users", "u").
				Select("id AS uid", "username").
				AddSelect("email, password").
				Where("id = 1").
				Build()
			if err != nil {
				t.Fatalf("unexpected error from Build(): %v", err)
			}

			// Columns should appear exactly as rendered, unprefixed
			wantCols := "SELECT u.id AS uid, u.username, u.email, u.password"
			if !strings.Contains(sql, wantCols) {
				t.Errorf("expected SQL to contain %q, got %q", wantCols, sql)
			}

			// FROM clause should reference the table (no aliasing on columns)
			if !strings.Contains(sql, "FROM users AS u") {
				t.Errorf("expected SQL to contain %q, got %q", "FROM users AS u", sql)
			}

			// WHERE clause and parameter binding
			if !strings.Contains(sql, "WHERE id = ?") {
				t.Errorf("expected SQL to contain %q, got %q", "WHERE id = ?", sql)
			}

			if len(args) != 1 {
				t.Errorf("expected 1 args but got %d args", len(args))
			}
		})
	})

	t.Run("Members", func(t *testing.T) {
		t.Run("From", func(t *testing.T) {
			t.Run("Basic", func(t *testing.T) {
				sql, args, err := builder.NewSelect(nil).From("users").Build()
				if err != nil {
					t.Fatalf("unexpected error from Build(): %v", err)
				}
				if !strings.Contains(sql, "FROM users") {
					t.Errorf("expected SQL to contain %q, got %q", "FROM users", sql)
				}
				if len(args) != 0 {
					t.Errorf("expected 1 args but got %d args", len(args))
				}
			})

			t.Run("EmptyTable", func(t *testing.T) {
				_, _, err := builder.NewSelect(nil).
					From("").
					Select("*").
					Build()
				if err == nil || !strings.Contains(err.Error(), "table expression is empty") {
					t.Errorf("expected error about expected table expression is empty, got: %v", err)
				}
			})

			t.Run("MissingFromClause", func(t *testing.T) {
				_, _, err := builder.NewSelect(nil).
					Select("id").
					Build()

				if err == nil || !strings.Contains(err.Error(), "missing source") {
					t.Errorf("expected error about missing source, got: %v", err)
				}
			})

			t.Run("SingleTableNoAlias", func(t *testing.T) {
				sql, _, err := builder.NewSelect(nil).
					Select("id").
					From("customers").
					Build()

				expected := "FROM customers"
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !strings.Contains(sql, expected) {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("AliasedTable", func(t *testing.T) {
				sql, _, err := builder.NewSelect(nil).
					Select("id").
					From("users", "u").
					Build()

				expected := "FROM users AS u"
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !strings.Contains(sql, expected) {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("InvalidTable", func(t *testing.T) {
				b := builder.NewSelect(nil).From(" AS x")
				_, _, err := b.Select("id").Build()
				if err == nil || !strings.Contains(err.Error(), "cannot start with") {
					t.Errorf("expected error for invalid table expression, got %v", err)
				}
			})

			t.Run("QualifiedColumns", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					From("users").
					Select("users.id")

				sql, args, _ := sb.Build()
				if !strings.Contains(sql, "SELECT id") {
					t.Errorf("SQL mismatch: expected %s, got %s", "SELECT id", sql)
				}
				if len(args) != 0 {
					t.Errorf("expected no args, got %v", args)
				}
			})
		})

		t.Run("Select", func(t *testing.T) {
			t.Run("Basic", func(t *testing.T) {
				sql, _, err := builder.NewSelect(nil).
					Select("id", "name").
					From("users").
					Build()

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !strings.Contains(sql, "SELECT id") {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", "SELECT id", sql)
				}
			})

			t.Run("CommaSeparatedColumns", func(t *testing.T) {
				sql, _, err := builder.NewSelect(nil).
					Select("id, name").
					From("users").
					Build()

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !strings.Contains(sql, "SELECT id, name") {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", "SELECT id, name", sql)
				}
			})

			t.Run("InlineAlias", func(t *testing.T) {
				sql, _, err := builder.NewSelect(nil).
					Select("id AS uid").
					From("users").
					Build()

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !strings.Contains(sql, "SELECT id AS uid") {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", "SELECT id as uid", sql)
				}
			})
		})

		t.Run("AddSelect", func(t *testing.T) {
			t.Run("WithSelect", func(t *testing.T) {
				sql, _, err := builder.NewSelect(nil).
					From("users").
					Select("id").
					AddSelect("name AS full_name").
					AddSelect("username", "email", "password").
					Build()

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "SELECT id, name AS full_name, username, email, password"
				if !strings.Contains(sql, expected) {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", expected, sql)
				}
			})

			t.Run("WithNoSelect", func(t *testing.T) {
				sql, _, err := builder.NewSelect(nil).
					From("users").
					AddSelect("id").
					AddSelect("username", "email", "password").
					Build()

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "SELECT id"
				if !strings.Contains(sql, expected) {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", expected, sql)
				}
			})
		})

		t.Run("Where", func(t *testing.T) {
			t.Run("ValidCondition", func(t *testing.T) {
				sql, _, err := builder.NewSelect(nil).
					Select("id").
					From("invoices").
					Where("paid = false").
					Build()

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "WHERE paid = ?"
				if !strings.Contains(sql, expected) {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", expected, sql)
				}
			})

			t.Run("InvalidCondition", func(t *testing.T) {
				_, _, err := builder.NewSelect(nil).
					Select("*").
					From("users").
					Where("active =").
					Build()
				if err == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !strings.Contains(err.Error(), "unable to parse condition") {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", "unable to parse condition", err.Error())
				}
			})
		})

		t.Run("AndWhere", func(t *testing.T) {
			t.Run("ValidCondition", func(t *testing.T) {
				sql, _, err := builder.NewSelect(nil).
					Select("id").
					From("invoices").
					Where("paid = false").
					AndWhere("amount > 100").
					AndWhere("overdue = true").
					Build()

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "WHERE paid = ? AND amount > ? AND overdue = ?"
				if !strings.Contains(sql, expected) {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", expected, sql)
				}
			})

			t.Run("InvalidCondition", func(t *testing.T) {
				_, _, err := builder.NewSelect(nil).
					Select("*").
					From("users").
					Where("active = ?", true).
					AndWhere("amount >").
					Build()
				if err == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !strings.Contains(err.Error(), "unable to parse condition") {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", "unable to parse condition", err.Error())
				}
			})
		})

		t.Run("OrWhere", func(t *testing.T) {
			t.Run("ValidCondition", func(t *testing.T) {
				sql, _, err := builder.NewSelect(nil).
					Select("id").
					From("invoices").
					Where("paid = false").
					OrWhere("amount > 100").
					Build()

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "WHERE paid = ? OR amount > ?"
				if !strings.Contains(sql, expected) {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", expected, sql)
				}
			})

			t.Run("InvalidCondition", func(t *testing.T) {
				_, _, err := builder.NewSelect(nil).
					Select("*").
					From("users").
					Where("active = ?", true).
					OrWhere("amount >").
					Build()
				if err == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !strings.Contains(err.Error(), "unable to parse condition") {
					t.Errorf("SQL mismatch: expected to contain %q, got %q", "unable to parse condition", err.Error())
				}
			})
		})

		t.Run("Pagination", func(t *testing.T) {
			sql, _, err := builder.NewSelect(nil).
				Select("*").
				From("users").
				Take(10).
				Skip(5).
				Build()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(sql, "LIMIT 10") {
				t.Errorf("SQL mismatch: expected to contain %q, got %q", "LIMIT 10", sql)
			}
			if !strings.Contains(sql, "OFFSET 5") {
				t.Errorf("SQL mismatch: expected to contain %q, got %q", "OFFSET 5", sql)
			}
		})

		t.Run("Limit", func(t *testing.T) {
			sql, _, err := builder.NewSelect(nil).Select("*").From("users").Take(10).Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(sql, "LIMIT 10") {
				t.Errorf("SQL mismatch: expected to contain %q, got %q", "LIMIT 10", sql)
			}
		})

		t.Run("Offset", func(t *testing.T) {
			sql, _, err := builder.NewSelect(nil).Select("*").From("users").Skip(5).Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(sql, "OFFSET 5") {
				t.Errorf("SQL mismatch: expected to contain %q, got %q", "OFFSET 5", sql)
			}
		})

		t.Run("OrderBy", func(t *testing.T) {
			sql, _, err := builder.NewSelect(nil).
				Select("*").
				From("users").
				Skip(5).
				OrderBy("id").
				Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(sql, "ORDER BY ") {
				t.Errorf("SQL mismatch: expected to contain %q, got %q", "ORDER BY id", sql)
			}
		})

		t.Run("UseDialect", func(t *testing.T) {
			sql, args, err := builder.NewSelect(driver.NewPostgresDialect()).
				Select("id", "created_at").
				From("users").
				Where("status", "active").
				UseDialect("postgres").
				Build()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(sql, `WHERE "status" = $1`) {
				t.Errorf("SQL mismatch: expected to contain %q, got %q", `WHERE "status" = $1`, sql)
			}
			if len(args) != 1 {
				t.Errorf("expected 1 args, got %d", len(args))
			}
		})

		t.Run("Build", func(t *testing.T) {
			qb := &builder.SelectBuilder{}

			t.Run("withNilDialect", func(t *testing.T) {
				if qb.GetDialect() == nil {
					t.Fatalf("expected dialect builder, got nil")
				}
			})

			t.Run("withNilDialect", func(t *testing.T) {
				qb := &builder.SelectBuilder{}
				_, _, err := qb.Build()
				if err == nil {
					t.Fatal("expected an error because no sources were configured")
				}
				if !strings.Contains(err.Error(), "missing source") {
					t.Fatalf("expected missing‐source error, got %q", err.Error())
				}
			})
		})
	})

	t.Run("Validations", func(t *testing.T) {
		t.Run("InvalidCondition", func(t *testing.T) {
			_, _, err := builder.NewSelect(nil).Build()
			if err == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(err.Error(), "missing source") {
				t.Errorf("expected %q, got %q", "missing source", err.Error())
			}
		})

		t.Run("HasDialect", func(t *testing.T) {
			b := builder.NewSelect(nil).From("users")
			_, _, err := b.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if "generic" != b.Dialect.GetName() {
				t.Errorf("generic dialect should be resolved when no dialect is set")
			}
		})

		t.Run("HasErrors", func(t *testing.T) {
			_, _, err := builder.NewSelect(nil).From("AS u").Build()
			if err == nil {
				t.Fatalf("Build() should have failed with missing source")
			}
			if !strings.Contains(err.Error(), "invalid table expression") {
				t.Errorf("expected %q, got %q", "invalid table expression", err.Error())
			}
		})
	})
}
