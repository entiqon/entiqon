// File: db/builder/select_test.go

package builder_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/builder"
	"github.com/entiqon/entiqon/db/token"
)

// Group 1: fields (between SELECT and FROM)
// Group 2: from clause (after FROM, until LIMIT/OFFSET or end)
// Group 3: "LIMIT <n>" (optional, full string with keyword)
// Group 4: "OFFSET <n>" (optional, full string with keyword)
var re = regexp.MustCompile(`(?i)^\s*SELECT\s+(.*?)\s+FROM\s+(.+?)(\s+LIMIT\s+\d+)?(\s+OFFSET\s+\d+)?\s*$`)

func TestSelectBuilder(t *testing.T) {
	t.Run("Methods", func(t *testing.T) {
		table := "users"

		t.Run("Build", func(t *testing.T) {
			t.Run("EmptyFields", func(t *testing.T) {
				// Empty string should not add a field → defaults to SELECT *
				sql, _ := builder.NewSelect(nil).
					Source(table).
					Fields("").
					Build()
				if !strings.Contains(sql, "SELECT *") {
					t.Errorf("expected to contain '*', got %v", sql)
				}
			})

			t.Run("EmptyInput", func(t *testing.T) {
				// Calling Fields("") should produce SELECT * as no fields were added
				sql, _ := builder.NewSelect(nil).
					Source(table).
					Fields("").
					Build()
				if !strings.Contains(sql, "SELECT *") {
					t.Errorf("expected to contain '*', got %v", sql)
				}
			})

			t.Run("WithFields", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("expr", "alias", true)
				fields := sb.GetFields()
				if len(fields) != 1 {
					t.Errorf("expected 1 field, got %d", len(fields))
				}
				if fields[0].Expr != "expr" || fields[0].Alias != "alias" || !fields[0].IsRaw {
					t.Errorf("expected field content, got %v", fields[0])
				}
			})
		})

		t.Run("Fields", func(t *testing.T) {
			t.Run("SingleField", func(t *testing.T) {
				t.Run("NoAlias", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("expr")
					fields := sb.GetFields()
					if len(fields) != 1 {
						t.Errorf("expected 1 field, got %d", len(fields))
					}
				})

				t.Run("Aliased", func(t *testing.T) {
					t.Run("BySpace", func(t *testing.T) {
						sb := builder.NewSelect(nil).
							Fields("id user_id")
						fields := sb.GetFields()
						if fields[0].Expr != "id" || fields[0].Alias != "user_id" {
							t.Errorf("expected field content, got %v", fields[0])
						}
					})

					t.Run("ByASKeyword", func(t *testing.T) {
						sb := builder.NewSelect(nil).
							Fields("id AS user_id")
						fields := sb.GetFields()
						if !fields[0].IsAliased() {
							t.Errorf("expected field to be aliased, got %v", fields[0].IsAliased())
						}
					})
				})

				t.Run("Raw", func(t *testing.T) {
					t.Run("GeneratedAlias", func(t *testing.T) {
						sb := builder.NewSelect(nil).
							Fields("UPPER(firstname || '-' || lastname)")
						fields := sb.GetFields()
						if len(fields) != 1 || !fields[0].IsAliased() {
							t.Errorf("expected exactly 1 aliased field, got %+v", fields)
						}
					})

					t.Run("WithAlias", func(t *testing.T) {
						sb := builder.NewSelect(nil).
							Fields("UPPER(firstname || '-' || lastname) fullname")
						fields := sb.GetFields()
						if len(fields) != 1 || !fields[0].IsErrored() {
							t.Errorf("expected exactly 1 errored field, got %+v", fields)
						}
						if fields[0].Error == nil {
							t.Errorf("expected error, got nil")
						}
					})
				})

				t.Run("Invalid", func(t *testing.T) {
					sb := builder.NewSelect(nil).Fields("")
					fields := sb.GetFields()
					if len(fields) != 0 {
						t.Errorf("expected 0 fields, got %d", len(fields))
					}
				})
			})

			t.Run("AliasedField", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("expr", "alias")
				fields := sb.GetFields()
				if len(fields) != 1 {
					t.Errorf("expected 1 field, got %d", len(fields))
				}
			})

			t.Run("RawField", func(t *testing.T) {

				t.Run("IsRawResolved", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("firstname || '-' || lastname", "fullname")
					fields := sb.GetFields()
					if len(fields) != 1 {
						t.Errorf("expected 1 field, got %d", len(fields))
					}
					if !fields[0].IsRaw {
						t.Errorf("expected IsRaw resolved to true, got %v", fields[0].IsRaw)
					}
				})

				t.Run("IsRawSpecific", func(t *testing.T) {
					sb := builder.NewSelect(nil).
						Fields("firstname || '-' || lastname", "fullname", true)
					fields := sb.GetFields()
					if len(fields) != 1 {
						t.Errorf("expected 1 field, got %d", len(fields))
					}
					if !fields[0].IsRaw {
						t.Errorf("expected IsRaw to be true, got %v", fields[0].IsRaw)
					}
				})

			})

			t.Run("Field", func(t *testing.T) {
				t.Run("Nil", func(t *testing.T) {
					var f *token.Field
					sb := builder.NewSelect(nil).
						Fields(f)
					fields := sb.GetFields()
					if len(fields) != 0 {
						t.Errorf("expected 0 fields, got %d", len(fields))
					}
				})

				t.Run("NonNil", func(t *testing.T) {
					f := token.NewField("id")
					sb := builder.NewSelect(nil).
						Fields(f) // *token.Field (non-nil)
					fields := sb.GetFields()
					if len(fields) != 1 {
						t.Errorf("expected 1 field, got %d", len(fields))
					}
					if fields[0].Expr != "id" {
						t.Errorf("expected Expr 'id', got %q", fields[0].Expr)
					}
				})

				t.Run("NotPointer", func(t *testing.T) {
					field := token.NewField("id")
					sb := builder.NewSelect(nil).
						Fields(*field) // pass by value, not pointer
					fields := sb.GetFields()
					if len(fields) != 1 {
						t.Errorf("expected 1 field, got %d", len(fields))
					}
					if fields[0].Expr != "id" {
						t.Errorf("expected Expr 'id', got %q", fields[0].Expr)
					}
					if fields[0].IsAliased() {
						t.Errorf("expected field not to be aliased, got %v", fields[0].IsAliased())
					}
				})
			})

			t.Run("InvalidType", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields(true)
				fields := sb.GetFields()
				if !fields[0].IsErrored() {
					t.Errorf("expected IsErrored to be true, got %v", fields[0].IsErrored())
				}
				if fields[0].Error == nil {
					t.Errorf("expected Error to be set, got nil")
				}
			})
		})
	})
}

func TestSelectBuilderLegacy(t *testing.T) {
	t.Run("Methods", func(t *testing.T) {
		table := "users"
		limit := 10
		offset := 20

		t.Run("Build", func(t *testing.T) {

			t.Run("InlineAliasWithSpace", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("UPPER(firstname || '-' || lastname) fullname").
					Source(table)
				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				matches := re.FindStringSubmatch(sql)
				if len(matches) != 5 {
					t.Errorf("expected 5 parts, got %d", len(matches))
				}
				if !strings.Contains(sql, matches[1]) {
					t.Errorf("expected %q, got %q", matches[1], sql)
				}
			})

			t.Run("InlineAliasWithASKeyword", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("UPPER(firstname || '-' || lastname) AS fullname").
					Source(table)
				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				matches := re.FindStringSubmatch(sql)
				if len(matches) != 5 {
					t.Errorf("expected 5 parts, got %d", len(matches))
				}
				if !strings.Contains(sql, matches[1]) {
					t.Errorf("expected %q, got %q", matches[1], sql)
				}
			})

			t.Run("Aliased", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields("id", "user_id").
					Source(table)
				sql, err := sb.Build()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				matches := re.FindStringSubmatch(sql)
				if len(matches) != 5 {
					t.Errorf("expected 5 parts, got %d", len(matches))
				}
				if !strings.Contains(sql, matches[1]) {
					t.Errorf("expected contains %q, got %q", matches[1], sql)
				}
			})

			t.Run("ColumnToken", func(t *testing.T) {
				sql, _ := builder.NewSelect(nil).
					Fields(
						token.Field{
							Expr: "id",
						},
						token.Field{
							Expr:  "name",
							Alias: "firstname",
						},
						token.Field{
							Expr: "lastname",
						},
						token.Field{
							Expr:  "(country || ', ' || state)",
							Alias: "country",
							IsRaw: true,
						},
					).
					Source(table).
					Limit(limit).
					Offset(offset).
					Build()
				matches := re.FindStringSubmatch(sql)
				if len(matches) != 5 {
					t.Errorf("expected 5 parts, got %d", len(matches))
				}
				if !strings.Contains(sql, matches[1]) {
					t.Errorf("expected %q, got %q", matches[1], sql)
				}
			})

			t.Run("Limit", func(t *testing.T) {
				sql, _ := builder.NewSelect(nil).
					Source(table).
					Limit(10).
					Build()
				matches := re.FindStringSubmatch(sql)
				if len(matches) != 5 {
					t.Errorf("expected 5 parts, got %d", len(matches))
				}
				if !strings.Contains(sql, matches[3]) {
					t.Errorf("expected contains %q, got %q", matches[3], sql)
				}
			})

			t.Run("Offset", func(t *testing.T) {
				sql, _ := builder.NewSelect(nil).
					Source(table).
					Offset(offset).
					Build()
				matches := re.FindStringSubmatch(sql)
				if len(matches) != 5 {
					t.Errorf("expected 5 parts, got %d", len(matches))
				}
				if !strings.Contains(sql, matches[4]) {
					t.Errorf("expected contains %q, got %q", matches[4], sql)
				}
			})

			t.Run("Pagination", func(t *testing.T) {
				sql, _ := builder.NewSelect(nil).
					Source(table).
					Limit(limit).
					Offset(offset).
					Build()
				matches := re.FindStringSubmatch(sql)
				if len(matches) != 5 {
					t.Errorf("expected 5 parts, got %d", len(matches))
				}
				want := fmt.Sprintf("%s %s", strings.TrimSpace(matches[3]), strings.TrimSpace(matches[4]))
				if !strings.Contains(sql, want) {
					t.Errorf("expected contains %q, got %q", want, sql)
				}
			})

			t.Run("SQL", func(t *testing.T) {
				sql, _ := builder.NewSelect(nil).
					Fields("id, name").
					Source(table).
					Limit(limit).
					Offset(offset).
					Build()
				matches := re.FindStringSubmatch(sql)
				if len(matches) != 5 {
					t.Errorf("expected 5 parts, got %d", len(matches))
				}
			})

			t.Run("IgnoreUnsupportedTypes_All", func(t *testing.T) {
				sb := builder.NewSelect(nil).
					Fields(123, 4.5, true, struct{}{}). // all unsupported → ignored
					Source(table)

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

		t.Run("Limit", func(t *testing.T) {
			sb := builder.NewSelect(nil).
				Fields("id").
				Source(table).
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
				Source(table).
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
				Source(table).
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
				Source(table)

			want := `Status ✅: SQL=SELECT id FROM "users", Params=`
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
	})
}
