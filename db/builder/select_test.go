// File: db/builder/select_test.go

package builder_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/builder"
	"github.com/entiqon/entiqon/db/token/field"
	"github.com/entiqon/entiqon/db/token/table"
)

func TestSelectBuilder(t *testing.T) {
	t.Run("Methods", func(t *testing.T) {
		t.Run("Fields", func(t *testing.T) {
			t.Run("NilCollection", func(t *testing.T) {
				sb := &builder.SelectBuilder{} // fields is nil
				sb.Fields("id")
				fields := sb.GetFields()
				f, _ := fields.At(0)
				if fields.Length() != 1 || f.Expr() != "id" {
					t.Errorf("expected one field 'id', got %+v", fields)
				}
			})

			t.Run("NoArgs", func(t *testing.T) {
				sb := builder.NewSelect(nil).Fields()
				fields := sb.GetFields()
				if fields.Length() != 0 {
					t.Errorf("expected no fields, got %d", fields.Length())
				}
			})

			t.Run("Add", func(t *testing.T) {
				sb := builder.NewSelect(nil).Fields("id")
				fields := sb.GetFields()
				f, _ := fields.At(0)
				if fields.Length() != 1 || f.Expr() != "id" {
					t.Errorf("expected reset only, got %+v", fields)
				}
			})

			t.Run("Reset", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Fields("reset") // should reset
				fields := sb.GetFields()
				f, _ := fields.At(0)
				if fields.Length() != 1 || f.Expr() != "reset" {
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
					Fields(field.New("id")) // *token.field
				fields := sb.GetFields()
				if fields.Length() != 1 {
					t.Errorf("expected 1 field, got %d", fields.Length())
				}
			})

			t.Run("NotPointer", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields(field.New("id")) // token.field (value)
				fields := sb.GetFields()
				if fields.Length() != 1 {
					t.Errorf("expected 1 field, got %d", fields.Length())
				}
			})

			t.Run("Invalid", func(t *testing.T) {
				sb := builder.NewSelect(nil).Fields(true)
				fields := sb.GetFields()
				f, _ := fields.At(0)
				if !f.IsErrored() {
					t.Errorf("expected IsErrored to be true, got %v", f.IsErrored())
				}
				if f.Error() == nil {
					t.Errorf("expected Error to be set, got nil")
				}
			})
		})

		t.Run("From", func(t *testing.T) {
			t.Run("InlineAlias", func(t *testing.T) {
				sql, err := builder.NewSelect(nil).
					Source("users u").
					Build()
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				want := "SELECT * FROM users AS u"
				if sql != want {
					t.Errorf("got %q, want %q", sql, want)
				}
			})

			t.Run("ExplicitAlias", func(t *testing.T) {
				sql, err := builder.NewSelect(nil).
					Source("users", "u").
					Build()
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				want := "SELECT * FROM users AS u"
				if sql != want {
					t.Errorf("got %q, want %q", sql, want)
				}
			})

			t.Run("Table", func(t *testing.T) {
				sql, err := builder.NewSelect(nil).
					Source(table.New("users")).
					Build()
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				want := "SELECT * FROM users"
				if sql != want {
					t.Errorf("got %q, want %q", sql, want)
				}
			})

			t.Run("Error", func(t *testing.T) {
				t.Run("MultipleSource", func(t *testing.T) {
					got := builder.NewSelect(nil).
						Source(table.New("users"), "account").
						String()
					want := "no source specified"
					if !strings.Contains(got, want) {
						t.Errorf("got %q, want %q", got, want)
					}
				})

				t.Run("InvalidName", func(t *testing.T) {
					_, err := builder.NewSelect(nil).
						Source(12345).
						Build()
					if err == nil {
						t.Errorf("expected error, got nil")
					}
					if !strings.Contains(err.Error(), "cannot be used as a table source") {
						t.Errorf("got %q, want %q", err, "cannot be used as a table source")
					}
				})

				t.Run("InvalidAlias", func(t *testing.T) {
					_, err := builder.NewSelect(nil).
						Source("users", false).
						Build()
					if err == nil {
						t.Errorf("expected no error, got %v", err)
					}
				})
			})
		})

		t.Run("AddFields", func(t *testing.T) {
			sb := builder.NewSelect(nil).
				Fields("id").
				AddFields("name").
				AddFields("email")

			fields := sb.GetFields()
			if fields.Length() != 3 {
				t.Errorf("expected 3 fields, got %d", fields.Length())
			}
		})

		t.Run("Conditions", func(t *testing.T) {
			t.Run("NilCollection", func(t *testing.T) {
				sb := &builder.SelectBuilder{}
				sb.Fields("id")
				sb.Source("users")
				sb.Where("age > 18")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				expected := "SELECT id FROM users WHERE age > 18"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("Where", func(t *testing.T) {
				t.Run("SingleCondition", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						Where("age > 18")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users WHERE age > 18"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})

				t.Run("ResetCollection", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						Where("age > 45").
						Where("age < 50")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users WHERE age < 50"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})

				t.Run("MultipleConditions", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						Where("age > 18", "status = 'active'")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users WHERE age > 18 AND status = 'active'"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})

				t.Run("IgnoreEmptyConditions", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						Where("   ", "role = 'admin'")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users WHERE role = 'admin'"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})
			})

			t.Run("And", func(t *testing.T) {
				t.Run("Appends", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						Where("age > 18").
						And("status = 'active'")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users WHERE age > 18 AND status = 'active'"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})

				t.Run("AsFirstCondition", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						And("status = 'active'")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users WHERE status = 'active'"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})
			})

			t.Run("Or", func(t *testing.T) {
				t.Run("Appends", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						Where("age > 18").
						Or("status = 'active'").
						Or("role = 'admin'")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users WHERE age > 18 OR status = 'active' OR role = 'admin'"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})

				t.Run("AsFirstCondition", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						Or("status = 'active'")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users WHERE status = 'active'"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})
			})
		})

		t.Run("Grouping", func(t *testing.T) {
			t.Run("NilCollection", func(t *testing.T) {
				sb := &builder.SelectBuilder{} // groupings is nil
				sb.Fields("id")
				sb.Source("users")
				sb.GroupBy("role")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				expected := "SELECT id FROM users GROUP BY role"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("GroupBy", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					GroupBy("department")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "SELECT id FROM users GROUP BY department"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("ThenGroupBy", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					GroupBy("department").
					ThenGroupBy("role")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "SELECT id FROM users GROUP BY department, role"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("ResetWithGroupBy", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					GroupBy("department").
					GroupBy("role") // should reset

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "SELECT id FROM users GROUP BY role"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("IgnoreEmptyGrouping", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					GroupBy("   ", "department")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "SELECT id FROM users GROUP BY department"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})
		})

		t.Run("Having", func(t *testing.T) {
			t.Run("NilCollection", func(t *testing.T) {
				sb := &builder.SelectBuilder{}
				sb.Fields("id")
				sb.Source("users")
				sb.GroupBy("role")
				sb.Having("COUNT(*) > 5")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				expected := "SELECT id FROM users GROUP BY role HAVING COUNT(*) > 5"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("SingleCondition", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					GroupBy("role").
					Having("COUNT(*) > 5")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				expected := "SELECT id FROM users GROUP BY role HAVING COUNT(*) > 5"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("ResetCollection", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					GroupBy("role").
					Having("COUNT(*) > 5").
					Having("COUNT(*) < 10")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				expected := "SELECT id FROM users GROUP BY role HAVING COUNT(*) < 10"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("MultipleConditions", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					GroupBy("role").
					Having("COUNT(*) > 5", "AVG(age) > 30")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				expected := "SELECT id FROM users GROUP BY role HAVING COUNT(*) > 5 AND AVG(age) > 30"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("IgnoreEmptyConditions", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					GroupBy("role").
					Having("   ", "SUM(salary) > 1000")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				expected := "SELECT id FROM users GROUP BY role HAVING SUM(salary) > 1000"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("And", func(t *testing.T) {
				t.Run("Appends", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						GroupBy("role").
						Having("COUNT(*) > 5").
						AndHaving("AVG(age) > 30")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users GROUP BY role HAVING COUNT(*) > 5 AND AVG(age) > 30"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})

				t.Run("AsFirstCondition", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						GroupBy("role").
						AndHaving("AVG(age) > 30")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users GROUP BY role HAVING AVG(age) > 30"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})
			})

			t.Run("Or", func(t *testing.T) {
				t.Run("Appends", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						GroupBy("role").
						Having("COUNT(*) > 5").
						OrHaving("AVG(age) > 30").
						OrHaving("SUM(salary) > 1000")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users GROUP BY role HAVING COUNT(*) > 5 OR AVG(age) > 30 OR SUM(salary) > 1000"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})

				t.Run("AsFirstCondition", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id").
						Source("users").
						GroupBy("role").
						OrHaving("SUM(salary) > 1000")

					sql, err := sb.Build()
					if err != nil {
						t.Fatalf("expected no error, got %v", err)
					}
					expected := "SELECT id FROM users GROUP BY role HAVING SUM(salary) > 1000"
					if sql != expected {
						t.Errorf("expected %q, got %q", expected, sql)
					}
				})
			})
		})

		t.Run("Sorting", func(t *testing.T) {
			t.Run("NilCollection", func(t *testing.T) {
				sb := &builder.SelectBuilder{} // sorting is nil
				sb.Fields("id")
				sb.Source("users")
				sb.OrderBy("name ASC")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				expected := "SELECT id FROM users ORDER BY name ASC"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("OrderBy", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					OrderBy("created_at DESC")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "SELECT id FROM users ORDER BY created_at DESC"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("ThenSort", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					OrderBy("created_at DESC").
					ThenOrderBy("id ASC")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "SELECT id FROM users ORDER BY created_at DESC, id ASC"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("ResetWithSortBy", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					OrderBy("name ASC").
					OrderBy("email DESC") // should reset

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "SELECT id FROM users ORDER BY email DESC"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})

			t.Run("IgnoreEmptySorting", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users").
					OrderBy("   ", "updated_at ASC")

				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expected := "SELECT id FROM users ORDER BY updated_at ASC"
				if sql != expected {
					t.Errorf("expected %q, got %q", expected, sql)
				}
			})
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

		t.Run("Debug", func(t *testing.T) {
			t.Run("WithSource", func(t *testing.T) {
				sb := builder.NewSelect(nil).Source("users")
				got := sb.Debug()
				want := "✅ SelectBuilder{fields:0, source: ✅ Table(users), where:0, groupBy:0, having:0, orderBy:0}"
				if got != want {
					t.Errorf("expected %q, got %q", want, got)
				}
			})

			t.Run("WithHaving", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Source("orders").
					Having("SUM(quantity) > 2")
				got := sb.Debug()
				want := "✅ SelectBuilder{fields:0, source: ✅ Table(orders), where:0, groupBy:0, having:1, orderBy:0}"
				if got != want {
					t.Errorf("expected %q, got %q", want, got)
				}
			})

			t.Run("Errors", func(t *testing.T) {
				t.Run("NilReceiver", func(t *testing.T) {
					var sb *builder.SelectBuilder = nil // nil receiver
					got := sb.Debug()
					want := "❌ SelectBuilder(nil)"
					if got != want {
						t.Errorf("expected %q, got %q", want, got)
					}
				})

				t.Run("EmptySource", func(t *testing.T) {
					sb := builder.NewSelect(nil)
					got := sb.Debug()
					want := "❌ SelectBuilder{fields:0, source:<nil>, where:0, groupBy:0, having:0, orderBy:0}"
					if got != want {
						t.Errorf("expected %q, got %q", want, got)
					}
				})
			})
		})

		t.Run("String", func(t *testing.T) {
			t.Run("Default", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id").
					Source("users")

				want := "✅ SelectBuilder: status: ready, fields=1, no conditions, grouped=false, sorted=false"
				got := sb.String()
				if got != want {
					t.Errorf("String() = %q, want %q", got, want)
				}

				// Test error case when no source
				sb = builder.NewSelect(nil).
					Fields("id")

				wantPrefix := "❌"
				got = sb.String()
				if !strings.HasPrefix(got, wantPrefix) {
					t.Errorf("String() error output = %q, want prefix %q", got, wantPrefix)
				}
			})

			t.Run("Conditions", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("SUM(quantity)").
					Source("orders").
					Having("SUM(quantity) > 0")

				want := "✅ SelectBuilder: status: ready, fields=1, conditions=1, grouped=false, sorted=false"
				got := sb.String()
				if got != want {
					t.Errorf("String() = %q, want %q", got, want)
				}
			})

			t.Run("Errors", func(t *testing.T) {
				t.Run("NilReceiver", func(t *testing.T) {
					var sb *builder.SelectBuilder = nil // nil receiver

					want := sb.String()
					if want == "" {
						t.Errorf("String() error output = %q, want non-empty", want)
					}
				})
			})
		})

		t.Run("Build", func(t *testing.T) {
			t.Run("EmptyFields", func(t *testing.T) {
				// Empty string should not add a field → defaults to SELECT *
				sql, _ := builder.NewSelect(nil).
					Source("users").
					Fields("").
					Build()
				if !strings.Contains(sql, "SELECT * FROM users") {
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
				if !strings.Contains(out, "⛔️ field(") {
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

			t.Run("Errors", func(t *testing.T) {
				t.Run("NilReceiver", func(t *testing.T) {
					var sb *builder.SelectBuilder = nil // nil receiver

					_, err := sb.Build()
					if err == nil {
						t.Fatal("expected error, got nil")
					}

					want := "❌ [Build] – Wrong initialization. Cannot build on receiver nil"
					if err.Error() != want {
						t.Errorf("unexpected error: got %q, want %q", err.Error(), want)
					}
				})

				t.Run("NoSource", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("id")
					_, err := sb.Build()
					if err == nil {
						t.Fatal("expected error, got nil")
					}
					if !strings.Contains(err.Error(), "Source is nil") {
						t.Errorf("expected error, got nil")
					}
				})
			})
		})
	})
}
