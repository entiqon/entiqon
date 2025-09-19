// File: db/builder/select_test.go

package selects_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/builder/selects"
	"github.com/entiqon/entiqon/db/token/condition"
	"github.com/entiqon/entiqon/db/token/field"
	"github.com/entiqon/entiqon/db/token/table"
	ct "github.com/entiqon/entiqon/db/token/types/condition"
	"github.com/entiqon/entiqon/db/token/types/operator"
)

func TestSelectBuilder(t *testing.T) {
	t.Run("Constructor", func(t *testing.T) {
		t.Run("Default", func(t *testing.T) {
			sb := selects.New(nil)
			if sb == nil {
				t.Fatal("expected a Select, got nil")
			}
		})
	})

	t.Run("Methods", func(t *testing.T) {
		t.Run("Fields", func(t *testing.T) {
			t.Run("Add", func(t *testing.T) {
				t.Run("Empty", func(t *testing.T) {
					sb := selects.New(nil).Fields()
					fields := sb.GetFields()
					if len(fields) > 0 {
						t.Fatal("expected no fields, got ", sb.GetFields())
					}
				})

				t.Run("1-arg", func(t *testing.T) {
					t.Run("String", func(t *testing.T) {
						sb := selects.New(nil).Fields("")
						fields := sb.GetFields()
						if len(fields) != 1 || !fields[0].IsErrored() {
							t.Fatal("expected 1 errored field, got ", fields[0].IsErrored())
						}

						sb = selects.New(nil).Fields("field")
						fields = sb.GetFields()
						if len(fields) != 1 || fields[0].Expr() != "field" {
							t.Errorf("expected one field 'id', got %+v", fields[0].Expr())
						}

						sb = selects.New(nil).Fields("field1, field2")
						fields = sb.GetFields()
						if len(fields) != 2 {
							t.Errorf("expected two fields, got %+v", fields)
						}
					})

					t.Run("Field", func(t *testing.T) {
						// add by value
						sb := selects.New(nil).Fields(field.New("field1"))

						// add by pointer
						f := field.New("field2")
						sb = sb.AppendFields(&f)

						fields := sb.GetFields()
						if len(fields) != 2 {
							t.Fatalf("expected 2 fields, got %d", len(fields))
						}

						for i, fld := range fields {
							if fld.IsErrored() {
								t.Fatalf("field[%d] should not be errored", i)
							}
						}
					})

					t.Run("InvalidType", func(t *testing.T) {
						sb := selects.New(nil).Fields(123456)
						fields := sb.GetFields()
						if len(fields) != 1 {
							t.Fatal("expected 1 field, got ", fields[0].IsErrored())
						}
						if !fields[0].IsErrored() {
							t.Fatal("expected field to be errored; got ", fields[0].String())
						}
					})
				})

				t.Run("2-arg", func(t *testing.T) {
					sb := selects.New(nil).Fields("field", "alias")
					fields := sb.GetFields()
					if len(fields) != 1 || fields[0].IsErrored() {
						t.Fatal("expected 1 field, got ", fields[0].IsErrored())
					}
				})
			})

			t.Run("Append", func(t *testing.T) {
				sb := selects.New(nil).Fields("id")
				fields := sb.GetFields()
				if len(fields) != 1 || fields[0].Render() != "id" {
					t.Errorf("expected one field 'id', got %+v", fields[0].Render())
				}

				sb.AppendFields("name")
				fields = sb.GetFields()
				if len(fields) != 2 || fields[1].Render() != "name" {
					t.Errorf("expected one field 'name', got %+v", fields[1].Expr())
				}

				sb.Fields("las_name")
				fields = sb.GetFields()
				if len(fields) != 1 || fields[0].Render() != "las_name" {
					t.Errorf("expected one field 'las_name', got %+v", fields[0].Render())
				}
			})

			t.Run("GetFields", func(t *testing.T) {
				sb := selects.New(nil)
				if len(sb.GetFields()) != 0 {
					t.Fatal("expected no fields, got ", len(sb.GetFields()))
				}

				sb.Fields("field")
				fields := sb.GetFields()
				if len(fields) != 1 || fields[0].Render() != "field" {
					t.Fatal("expected one field 'field', got ", fields[0].Render())
				}
			})
		})

		t.Run("From", func(t *testing.T) {
			t.Run("Error", func(t *testing.T) {
				sb := selects.New(nil).From()
				tbl := sb.Table()
				if tbl == nil {
					t.Fatal("expected table token, got nil")
				}
				if !tbl.IsErrored() {
					t.Fatal("expected table to be errored")
				}
			})

			t.Run("String", func(t *testing.T) {
				t.Run("NoAlias", func(t *testing.T) {
					sb := selects.New(nil).From("table")
					tbl := sb.Table()
					if tbl == nil {
						t.Fatal("expected table token, got nil")
					}
				})

				t.Run("InlineAlias", func(t *testing.T) {
					sb := selects.New(nil).From("users u")
					tbl := sb.Table()
					if tbl.Name() != "users" {
						t.Fatal("expected table name 'users', got ", tbl.Name())
					}
					if tbl.Alias() != "u" {
						t.Fatal("expected table name 'users', got ", tbl.Name())
					}
				})

				t.Run("ExplicitAlias", func(t *testing.T) {
					sb := selects.New(nil).From("users", "u")
					tbl := sb.Table()
					if tbl.Name() != "users" {
						t.Fatal("expected table name 'users', got ", tbl.Name())
					}
					if tbl.Alias() != "u" {
						t.Fatal("expected table name 'users', got ", tbl.Name())
					}
				})
			})

			t.Run("Table", func(t *testing.T) {
				sb := selects.New(nil).From(table.New("users"))
				tbl := sb.Table()
				if tbl == nil {
					t.Fatal("expected table token, got nil")
				}

				tbl2 := table.New("users")
				sb = selects.New(nil).From(&tbl2)
				tbl = sb.Table()
				if tbl != tbl2 {
					t.Fatal("expected table name 'users', got ", tbl.Name())
				}
			})
		})

		t.Run("Join", func(t *testing.T) {
			t.Run("InnerJoin", func(t *testing.T) {
				t.Run("Error", func(t *testing.T) {
					sb := selects.New(nil).From("users", "u").
						InnerJoin(123456, "accounts a", "a.user_id = u.id")
					joins := sb.Joins()
					if len(joins) != 1 || !joins[0].IsErrored() {
						t.Fatal("expected 1 errored join, got ", len(joins))
					}
				})

				t.Run("String", func(t *testing.T) {
					sb := selects.New(nil).From("users", "u").
						InnerJoin("users u", "accounts a", "a.user_id = u.id")
					joins := sb.Joins()
					expected := "INNER JOIN accounts AS a ON a.user_id = u.id"
					if len(joins) != 1 || joins[0].Render() != expected {
						t.Fatalf("expected 1 join %q, got %v", expected, joins)
					}
				})

				t.Run("Table", func(t *testing.T) {
					tbl := table.New("users", "u")
					sb := selects.New(nil).From("users", "u").
						InnerJoin(tbl, "accounts a", fmt.Sprintf("a.user_id = %s.id", tbl.Alias()))
					joins := sb.Joins()
					expected := "INNER JOIN accounts AS a ON a.user_id = u.id"
					if len(joins) != 1 || joins[0].Render() != expected {
						t.Fatalf("expected 1 join %q, got %v", expected, joins)
					}

					tbl2 := joins[0].Right()
					sb.InnerJoin(tbl2, "permissions p", fmt.Sprintf("p.account_id = %s.id", tbl2.Alias()))
					joins = sb.Joins()
					expected = "INNER JOIN permissions AS p ON p.account_id = a.id"
					if len(joins) != 2 || joins[1].Render() != expected {
						t.Fatalf("expected 2 join %q, got %v", expected, joins)
					}
				})

				t.Run("Pointer", func(t *testing.T) {
					base := table.New("users", "u")
					related := table.New("accounts", "a")
					sb := selects.New(nil).From(base).
						InnerJoin(
							&base,
							related,
							fmt.Sprintf("%s.user_id = %s.id", related.Alias(), base.Alias()),
						)
					joins := sb.Joins()
					expected := "INNER JOIN accounts AS a ON a.user_id = u.id"
					if len(joins) != 1 || joins[0].Render() != expected {
						t.Fatalf("expected 1 join %q, got %v", expected, joins)
					}
				})

				t.Run("NotSource", func(t *testing.T) {
					tbl1 := table.New("users", "u")
					tbl2 := table.New("accounts", "a")
					tbl3 := table.New("permissions", "p")
					sb := selects.New(nil).From(tbl1).
						InnerJoin(
							tbl1,
							tbl2,
							fmt.Sprintf("%s.user_id = %s.id", tbl2.Alias(), tbl1.Alias()),
						)
					sb.InnerJoin(&tbl2, tbl3, fmt.Sprintf("%s.account_id = %s.id", tbl3.Alias(), tbl2.Alias()))
					joins := sb.Joins()
					expected := "INNER JOIN permissions AS p ON p.account_id = a.id"
					if len(joins) != 2 || joins[1].Render() != expected {
						t.Fatalf("expected 2 join %q, got %v", expected, joins[1].Render())
					}
				})

				t.Run("Nil", func(t *testing.T) {
					base := table.New("users", "u")
					var nilTbl table.Token
					tbl3 := table.New("permissions", "p")
					sb := selects.New(nil).From(base).
						InnerJoin(&nilTbl, tbl3, fmt.Sprintf("%s.account_id = %s.id", tbl3.Alias(), ""))
					joins := sb.Joins()
					if len(joins) != 1 || !joins[0].IsErrored() {
						t.Fatalf("expected 1 errored join, got %v", joins)
					}
				})
			})

			t.Run("LeftJoin", func(t *testing.T) {
				tbl := table.New("users", "u")
				sb := selects.New(nil).From("users", "u").
					LeftJoin(tbl, "accounts a", fmt.Sprintf("a.user_id = %s.id", tbl.Alias()))
				joins := sb.Joins()
				expected := "LEFT JOIN accounts AS a ON a.user_id = u.id"
				if len(joins) != 1 || joins[0].Render() != expected {
					t.Fatalf("expected 1 join %q, got %v", expected, joins)
				}
			})

			t.Run("RightJoin", func(t *testing.T) {
				tbl := table.New("users", "u")
				sb := selects.New(nil).From("users", "u").
					RightJoin(tbl, "accounts a", fmt.Sprintf("a.user_id = %s.id", tbl.Alias()))
				joins := sb.Joins()
				expected := "RIGHT JOIN accounts AS a ON a.user_id = u.id"
				if len(joins) != 1 || joins[0].Render() != expected {
					t.Fatalf("expected 1 join %q, got %v", expected, joins)
				}
			})

			t.Run("FullJoin", func(t *testing.T) {
				sb := selects.New(nil).
					Fields("u.id").
					From("users u").
					FullJoin("users u", "accounts a", "a.user_id = u.id")
				joins := sb.Joins()
				expected := "FULL JOIN accounts AS a ON a.user_id = u.id"
				if len(joins) != 1 || joins[0].Render() != expected {
					t.Fatalf("expected 1 join %q, got %v", expected, joins)
				}
			})

			t.Run("CrossJoin", func(t *testing.T) {
				sb := selects.New(nil).
					Fields("u.id").
					From("users u").
					CrossJoin("countries c")
				joins := sb.Joins()
				expected := "CROSS JOIN countries AS c"
				if len(joins) != 1 || joins[0].Render() != expected {
					t.Fatalf("expected 1 join %q, got %v", expected, joins)
				}
			})

			t.Run("NaturalJoin", func(t *testing.T) {
				sb := selects.New(nil).
					Fields("u.id").
					From("users u").
					NaturalJoin("countries c")
				joins := sb.Joins()
				expected := "NATURAL JOIN countries AS c"
				if len(joins) != 1 || joins[0].Render() != expected {
					t.Fatalf("expected 1 join %q, got %v", expected, joins)
				}
			})

			t.Run("Joins", func(t *testing.T) {
				sb := selects.New(nil)
				if len(sb.Joins()) != 0 {
					t.Fatal("expected no joins, got ", sb.Joins())
				}

				sb = selects.New(nil).
					Fields("u.id").
					Fields("o.id").
					AppendFields("p.amount").
					From("users u").
					InnerJoin("users u", "orders o", "u.id = o.user_id").
					LeftJoin("orders o", "payments p", "o.id = p.order_id").
					CrossJoin("countries c").
					NaturalJoin("states s")
				if len(sb.Joins()) != 4 {
					t.Fatal("expected 4 joins, got ", len(sb.Joins()))
				}
			})
		})

		t.Run("Conditions", func(t *testing.T) {
			t.Run("NoConditions", func(t *testing.T) {
				sb := selects.New(nil).From("users")
				conditions := sb.Conditions()
				if len(conditions) != 0 {
					t.Fatal("expected no conditions, got ", conditions)
				}

				sb = selects.New(nil).From("users").Where()
				conditions = sb.Conditions()
				if len(conditions) != 0 {
					t.Fatal("expected no conditions, got ", conditions)
				}
			})

			t.Run("String", func(t *testing.T) {
				t.Run("EmptyArg", func(t *testing.T) {
					sb := selects.New(nil).From("users").Where("")
					conditions := sb.Conditions()
					if len(conditions) != 1 || !conditions[0].IsErrored() {
						t.Fatal("expected 1 errored conditions, got ", conditions)
					}
				})

				t.Run("1-arg", func(t *testing.T) {
					sb := selects.New(nil).
						From("users").
						Where("age >= 45")

					conditions := sb.Conditions()
					if len(conditions) != 1 || conditions[0].IsErrored() {
						t.Fatalf("expected 1 non-errored condition, got %v", conditions)
					}
				})

				t.Run("2-arg", func(t *testing.T) {
					sb := selects.New(nil).
						From("users").
						Where("age", 45)

					conditions := sb.Conditions()
					if len(conditions) != 1 {
						t.Fatalf("expected 1 condition, got %v", conditions)
					}
					c := conditions[0]
					if c.Name() != "age" || c.Operator() != operator.Equal || c.Value() != 45 {
						t.Fatalf("expected 1 condition, got %v", c)
					}
				})

				t.Run("3-arg", func(t *testing.T) {
					sb := selects.New(nil).
						From("users").
						Where("age", operator.GreaterThanOrEqual, 45)

					conditions := sb.Conditions()
					if len(conditions) != 1 {
						t.Fatalf("expected 1 condition, got %v", conditions)
					}
					c := conditions[0]
					if c.Name() != "age" || c.Operator() != operator.GreaterThanOrEqual || c.Value() != 45 {
						t.Fatalf("expected 1 condition, got %v", c)
					}
				})
			})

			t.Run("Condition", func(t *testing.T) {
				sb := selects.New(nil).From("users").
					Where(condition.New(ct.Single, "age >= 45"))
				conditions := sb.Conditions()
				if len(conditions) != 1 {
					t.Fatal("expected 6 conditions, got ", len(conditions))
				}

				c := condition.New(ct.Or, "state", operator.Like, "New Jersey")
				sb.Where(&c)
				conditions = sb.Conditions()
				if len(conditions) != 1 {
					t.Fatal("expected 6 conditions, got ", len(conditions))
				}
			})

			t.Run("MixedArgTypes", func(t *testing.T) {
				c := condition.New(ct.Or, "country", operator.Like, "USA")
				sb := selects.New(nil).From("users").
					Where(
						condition.New(ct.Single, "age >= 45"),
						condition.New(ct.And, "age <= 50"),
						&c,
						condition.New(ct.And, "state = 'New York'"),
						condition.New(ct.Or, "state", "'New Jersey'"),
						"state = 'Texas'",
					)
				conditions := sb.Conditions()
				if len(conditions) != 6 {
					t.Fatal("expected 6 conditions, got ", len(conditions))
				}
			})

			t.Run("MixedConditionTypes", func(t *testing.T) {
				c := condition.New(ct.Or, "country", operator.Like, "USA")
				sb := selects.New(nil).
					Where(condition.New(ct.Single, "age >= 45")).
					AndWhere("age", operator.LessThanOrEqual, 50).
					OrWhere(&c, condition.New(ct.And, "state = 'New York'"), "state = 'Texas'")
				conditions := sb.Conditions()
				if len(conditions) != 5 {
					t.Fatalf("expected 5 condition, got %v", len(conditions))
				}
			})
		})

		t.Run("Groupings", func(t *testing.T) {
			t.Run("NoGrouping", func(t *testing.T) {
				sb := selects.New(nil).From("users")
				groupings := sb.Groupings()
				if len(groupings) != 0 {
					t.Fatal("expected no conditions, got ", len(groupings))
				}

				sb.GroupBy()
				groupings = sb.Groupings()
				if len(groupings) != 0 {
					t.Fatal("expected no conditions, got ", len(groupings))
				}
			})

			t.Run("GroupBy", func(t *testing.T) {
				sb := selects.New(nil).
					Fields("id").
					From("users").
					GroupBy("department")
				groupings := sb.Groupings()
				if len(groupings) != 1 {
					t.Fatal("expected 1 grouping, got ", len(groupings))
				}

				sb.GroupBy("role")
				if len(groupings) != 1 {
					t.Fatal("expected 1 grouping, got ", len(groupings))
				}
			})

			t.Run("ThenGroupBy", func(t *testing.T) {
				sb := selects.New(nil).
					Fields("id").
					From("users").
					GroupBy("department").
					ThenGroupBy("role")
				groupings := sb.Groupings()
				if len(groupings) != 2 {
					t.Fatal("expected 2 grouping, got ", len(groupings))
				}
			})
		})

		t.Run("Sorting", func(t *testing.T) {
			t.Run("NoSorting", func(t *testing.T) {
				sb := selects.New(nil).From("users")
				sorting := sb.Sorting()
				if len(sorting) != 0 {
					t.Fatal("expected no sorting, got ", len(sorting))
				}

				sb.OrderBy()
				sorting = sb.Sorting()
				if len(sorting) != 0 {
					t.Fatal("expected no sorting, got ", len(sorting))
				}
			})

			t.Run("OrderBy", func(t *testing.T) {
				sb := selects.New(nil).
					Fields("id").
					From("users").
					OrderBy("created_at DESC")
				sorting := sb.Sorting()
				if len(sorting) != 1 {
					t.Fatal("expected 1 sorting token, got ", len(sorting))
				}

				sb.OrderBy("updated_at DESC")
				sorting = sb.Sorting()
				if len(sorting) != 1 {
					t.Fatal("expected 1 sorting token, got ", len(sorting))
				}
			})

			t.Run("ThenOrderBy", func(t *testing.T) {
				sb := selects.New(nil).
					Fields("id").
					From("users").
					OrderBy("created_at DESC").
					ThenOrderBy("id ASC")
				sorting := sb.Sorting()
				if len(sorting) != 2 {
					t.Fatal("expected 2 sorting tokens, got ", len(sorting))
				}
			})
		})

		t.Run("Having", func(t *testing.T) {
			t.Run("NoConditions", func(t *testing.T) {
				sb := selects.New(nil).From("users")
				conditions := sb.HavingConditions()
				if len(conditions) != 0 {
					t.Fatal("expected no aggregate condition tokens, got ", len(conditions))
				}

				sb.OrderBy()
				conditions = sb.HavingConditions()
				if len(conditions) != 0 {
					t.Fatal("expected no aggregate condition tokens, got ", len(conditions))
				}
			})

			t.Run("EmptyArgs", func(t *testing.T) {
				sb := selects.New(nil).Having("")
				conditions := sb.HavingConditions()
				if len(conditions) != 0 {
					t.Fatal("expected 0 aggregate condition, got ", len(conditions))
				}
			})

			t.Run("Having", func(t *testing.T) {
				t.Run("1-arg", func(t *testing.T) {
					sb := selects.New(nil).
						Fields("department_id").
						From("collaborators").
						GroupBy("department_id").
						Having("COUNT(id) = 5")
					conditions := sb.HavingConditions()
					if len(conditions) != 1 {
						t.Fatal("expected 1 aggregate condition, got ", len(conditions))
					}

					sb.Having("COUNT(id) >= 5")
					conditions = sb.HavingConditions()
					if len(conditions) != 1 {
						t.Fatal("expected 1 aggregate condition, got ", len(conditions))
					}
				})

				t.Run("2-args", func(t *testing.T) {
					sb := selects.New(nil).Having("COUNT(id) >= 18", "COUNT(id) >= 65")
					conditions := sb.HavingConditions()
					if len(conditions) != 2 {
						t.Fatal("expected no aggregate condition tokens, got ", len(conditions))
					}
				})
			})

			t.Run("MultipleHavingTypes", func(t *testing.T) {
				sb := selects.New(nil).
					Fields("department_id").
					From("collaborators").
					GroupBy("department_id").
					Having("COUNT(id) >= 5").
					AndHaving("COUNT(id) <= 10").
					OrHaving("SUM(salary) >= 100000")
				conditions := sb.HavingConditions()
				if len(conditions) != 3 {
					t.Fatal("expected 3 having condition tokens, got ", len(conditions))
				}
			})
		})

		t.Run("Pagination", func(t *testing.T) {
			sb := selects.New(nil).Take(10)
			take := sb.Limit()
			if take != 10 {
				t.Fatal("expected take 10, got ", take)
			}

			sb.Skip(20)
			skip := sb.Offset()
			if skip != 20 {
				t.Fatal("expected skip 20, got ", skip)
			}

			sb.Take(10).Skip(0) // page 1
			take, skip = sb.Pagination()
			if take != 10 || skip != 0 {
				t.Fatal("expected take 10 skip 0, got ", skip)
			}
		})

		t.Run("Debug", func(t *testing.T) {
			sb := selects.New(nil)
			got := sb.Debug()
			want := "SelectBuilder{table:<nil>, fields:0, join:0, where:0, groupBy:0, having:0, orderBy:0, limit:0, offset:0}"
			if got != want {
				t.Fatal("expected ", want, " got ", got)
			}

			sb = selects.New(nil).From("users").
				Fields("COUNT(id)", "collaborators").
				LeftJoin("users u", "accounts a", "a.user_id = u.id").
				Where("department_id", 144).
				GroupBy("department_id").
				OrderBy("created_at DESC").
				Having("SUM(qty) > 2").
				Take(10).
				Skip(0)
			got = sb.Debug()
			want = "SelectBuilder{table:Table(\"users\"), fields:1, join:1, where:1, groupBy:1, having:1, orderBy:1, limit:10, offset:0}"
			if got != want {
				t.Errorf("expected %q, got %q", want, got)
			}
		})

		t.Run("String", func(t *testing.T) {
			sb := selects.New(nil)
			got := sb.String()
			want := "SelectBuilder: status=invalid – no table specified"
			if got != want {
				t.Fatal("expected ", want, " got ", got)
			}

			sb = selects.New(nil).From("users").
				Fields("COUNT(id)", "collaborators").
				LeftJoin("users u", "accounts a", "a.user_id = u.id").
				Where("department_id", 144).
				GroupBy("department_id").
				OrderBy("created_at DESC").
				Having("SUM(qty) > 2").
				Take(10).
				Skip(0)
			got = sb.String()
			want = "SelectBuilder: status:ready, table:Table(\"users\"), fields=1, joined=true, conditions=2, grouped=true, sorted=true"
			if got != want {
				t.Errorf("expected %q, got %q", want, got)
			}
		})

		t.Run("Build", func(t *testing.T) {
			t.Run("Errored", func(t *testing.T) {
				base := table.New("users", "u")

				_, _, err := selects.New(nil).Build()
				if err == nil {
					t.Fatal("expected error, got none")
				}
				if !strings.Contains(err.Error(), "no table specified") {
					t.Fatal("expected expr has invalid format, got ", err)
				}

				_, _, err = selects.New(nil).From(123456).Build()
				if !strings.Contains(err.Error(), "expr has invalid format (type int)") {
					t.Fatal("expected expr has invalid format, got ", err)
				}

				_, _, err = selects.New(nil).From(base).
					Fields(111111).
					Build()
				if !strings.Contains(err.Error(), "expr has invalid format") {
					t.Fatal("expected expr has invalid format, got ", err)
				}

				_, _, err = selects.New(nil).
					From(base).
					Fields("id, name, email").
					InnerJoin(base, 123456, "").
					Build()
				if !strings.Contains(err.Error(), "token invalid") {
					t.Fatal("expected contains token invalid, got ", err)
				}

				_, _, err = selects.New(nil).
					From(base).
					Fields("id, name, email").
					Where("id", operator.Invalid, 99).
					Build()
				if !strings.Contains(err.Error(), "invalid operator") {
					t.Fatal("expected contains token invalid, got ", err)
				}
			})

			t.Run("WithFields", func(t *testing.T) {
				sb := selects.New(nil)

				// Start field
				sql, _, _ := sb.From("users").Build()
				if !strings.Contains(sql, "SELECT * FROM users") {
					t.Errorf("expected to contain '*', got %v", sql)
				}

				// Valid fields
				sql, _, _ = selects.New(nil).
					Fields("id, name, email").
					From("users").
					Build()

				want := "SELECT id, name, email FROM users"
				if sql != want {
					t.Errorf("expected SQL to contain 'SELECT id, name', got %q", sql)
				}
			})

			t.Run("WithJoin", func(t *testing.T) {
				sb := selects.New(nil)

				base := table.New("users", "u")
				sql, _, _ := sb.Fields("id, name, email").
					From(base).
					InnerJoin(base, "accounts a", "a.user_id = u.id").
					Build()
				want := "SELECT id, name, email FROM users AS u INNER JOIN accounts AS a ON a.user_id = u.id"
				if sql != want {
					t.Errorf("expected %s, got %q", want, sql)
				}
			})

			t.Run("WithConditions", func(t *testing.T) {
				sql, params, _ := selects.New(nil).
					From("users").
					Where("id", operator.GreaterThan, 10).
					Build()
				want := "SELECT * FROM users WHERE id > :id"
				if sql != want {
					t.Errorf("expected `%s`, got `%s`", want, sql)
				}
				if len(params) != 1 {
					t.Fatal("expected 1 param, got ", len(params))
				}
			})

			t.Run("WithGroupBy", func(t *testing.T) {
				sb := selects.New(nil).
					Fields("COUNT(id)", "collaborators").
					From("users").
					GroupBy()
				sql, _, _ := sb.Build()
				want := "SELECT COUNT(id) AS collaborators FROM users"
				if sql != want {
					t.Errorf("expected `%s`, got `%s`", want, sql)
				}

				sb.GroupBy("department_id")
				sql, _, _ = sb.Build()
				want += " GROUP BY department_id"
				if sql != want {
					t.Errorf("expected `%s`, got `%s`", want, sql)
				}
			})

			t.Run("WithOrderBy", func(t *testing.T) {
				sb := selects.New(nil).
					Fields("id, name, email").
					From("users").
					OrderBy()
				sql, _, _ := sb.Build()
				want := "SELECT id, name, email FROM users"
				if sql != want {
					t.Errorf("expected `%s`, got `%s`", want, sql)
				}

				sb.OrderBy("id").ThenOrderBy("name")
				sql, _, _ = sb.Build()
				want += " ORDER BY id, name"
				if sql != want {
					t.Errorf("expected `%s`, got `%s`", want, sql)
				}
			})

			t.Run("WithHaving", func(t *testing.T) {
				sql, _, _ := selects.New(nil).
					Fields("department_id", "department").
					AppendFields("COUNT(id)", "collaborators").
					From("users").
					Having("COUNT(id) > 5").
					Build()
				want := "SELECT department_id AS department, COUNT(id) AS collaborators FROM users HAVING COUNT(id) > 5"
				if sql != want {
					t.Errorf("expected `%s`, got `%s`", want, sql)
				}
			})

			t.Run("WithPagination", func(t *testing.T) {
				sql, _, _ := selects.New(nil).
					From("users").
					Take(10).
					Build()
				want := "SELECT * FROM users LIMIT 10"
				if sql != want {
					t.Errorf("expected `%s`, got `%s`", want, sql)
				}

				sql, _, _ = selects.New(nil).
					From("users").
					Skip(20).
					Build()
				want = "SELECT * FROM users OFFSET 20"
				if sql != want {
					t.Errorf("expected `%s`, got `%s`", want, sql)
				}

				sql, _, _ = selects.New(nil).
					From("users").
					Take(10).
					Skip(20).
					Build()
				want = "SELECT * FROM users LIMIT 10 OFFSET 20"
				if want != sql {
					t.Errorf("expected `%s`, got `%s`", want, sql)
				}
			})
		})
	})
}
