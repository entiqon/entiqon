// File: db/token/column_test.go

package field_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/token/field"
)

func TestField(t *testing.T) {
	t.Run("Constructors", func(t *testing.T) {
		t.Run("NoInputs", func(t *testing.T) {
			f := field.NewField()
			if f == nil || f.Expr != "" || f.Alias != "" {
				t.Errorf("expected empty Field, got %+v", f)
			}
		})

		t.Run("Expr", func(t *testing.T) {
			t.Run("Field", func(t *testing.T) {
				src := field.NewField("col1")
				f := field.NewField(src)
				if f.Error() == nil {
					t.Errorf("expected error about Clone, got %+v", f.Error())
				}
			})

			t.Run("IsRaw", func(t *testing.T) {

				t.Run("WithOperator", func(t *testing.T) {
					operators := []string{"||", "+", "-", "*", "/"}
					for _, op := range operators {
						expr := fmt.Sprintf("col1 %s col2", op)
						f := field.NewField(expr)
						if !f.IsRaw {
							t.Errorf("expected IsRaw = true for expr %q, got false", expr)
						}
					}
				})

				t.Run("WithParenthesis", func(t *testing.T) {
					src := field.NewField("(col1 || '-' || col2)", " mixed")
					f := field.NewField(src)
					if f.Error() == nil {
						t.Errorf("expected error about Clone, got %+v", f.Error())
					}
				})

				t.Run("WithAliasRules", func(t *testing.T) {
					tests := []struct {
						name        string
						input       string
						expectErr   bool
						expectRaw   bool
						expectAlias string // expected alias if not error
					}{
						{
							name:        "NoAliasAutoGenerate",
							input:       "col1 || '-' || col2",
							expectErr:   false,
							expectRaw:   true,
							expectAlias: "",
						},
						{
							name:        "ExplicitAS",
							input:       "col1 || '-' || col2 AS full_name",
							expectErr:   false,
							expectRaw:   true,
							expectAlias: "full_name",
						},
						{
							name:      "SpaceAliasError",
							input:     "col1 || '-' || col2 full_name",
							expectErr: true,
							expectRaw: true,
						},
						{
							name:        "ComplexExprNoAlias",
							input:       "(1*rate(principal-balance))",
							expectErr:   false,
							expectRaw:   true,
							expectAlias: "",
						},
						{
							name:        "ComplexExprWithAS",
							input:       "(1*rate(principal-balance)) AS calc_rate",
							expectErr:   false,
							expectRaw:   true,
							expectAlias: "calc_rate",
						},
						{
							name:      "ComplexExprWithSpaceAliasError",
							input:     "(1*rate(principal-balance)) calc_rate",
							expectErr: true,
							expectRaw: true,
						},
					}

					for _, tt := range tests {
						t.Run(tt.name, func(t *testing.T) {
							f := field.NewField(tt.input)

							if tt.expectErr {
								if f.Error() == nil {
									t.Errorf("expected error, got nil")
								}
								return
							}

							if f.Error() != nil {
								t.Errorf("unexpected error: %v", f.Error())
							}

							if f.IsRaw != tt.expectRaw {
								t.Errorf("expected IsRaw=%v, got %v", tt.expectRaw, f.IsRaw)
							}

							if tt.expectAlias != "" {
								if f.Alias != tt.expectAlias {
									t.Errorf("expected alias %q, got %q", tt.expectAlias, f.Alias)
								}
							} else {
								if !strings.HasPrefix(f.Alias, "raw_expr_") {
									t.Errorf("expected auto-generated alias to start with raw_expr_, got %q", f.Alias)
								}
								if len(f.Alias) != len("raw_expr_")+6 {
									t.Errorf("expected alias length %d, got %d", len("raw_expr_")+6, len(f.Alias))
								}
							}
						})
					}
				})

				t.Run("HasTrailingAliasWithoutAS", func(t *testing.T) {
					cases := []struct {
						Name     string
						Expr     string
						Expected bool
					}{
						{
							Name:     "ExplicitAS",
							Expr:     "col1 || '-' || col2 AS full_name",
							Expected: false,
						},
						{
							Name:     "SingleToken",
							Expr:     "col1",
							Expected: false,
						},
						{
							Name:     "PenultimateOperator",
							Expr:     "col1 || col2",
							Expected: false,
						},
						{
							Name:     "AliasWithoutAS",
							Expr:     "col1 full_name",
							Expected: true,
						},
						{
							Name:     "NonIdentifierLastToken",
							Expr:     "col1 'abc'",
							Expected: false,
						},
					}

					for _, c := range cases {
						t.Run(c.Name, func(t *testing.T) {
							got := field.HasTrailingAliasWithoutAS(c.Expr)
							if got != c.Expected {
								t.Errorf("expected %v, got %v", c.Expected, got)
							}
						})
					}
				})
			})

			t.Run("UnsupportedType", func(t *testing.T) {
				f := field.NewField(123)
				if f.Error() == nil || !strings.Contains(f.Error().Error(), "input type unsupported: int") {
					t.Errorf("expected unsupported expr type error, got %+v", f.Error())
				}
			})

		})

		t.Run("OneArgSimpleExpr", func(t *testing.T) {
			f := field.NewField("col1")
			if f.Expr != "col1" || f.Alias != "" || f.IsRaw {
				t.Errorf("unexpected field %+v", f)
			}
		})

		t.Run("OneArgWithAS", func(t *testing.T) {
			f := field.NewField("col1 AS c1")
			if f.Expr != "col1" || f.Alias != "c1" {
				t.Errorf("expected expr=col1 alias=c1, got %+v", f)
			}
		})

		t.Run("OneArgWithSpaceAlias", func(t *testing.T) {
			f := field.NewField("col1 c1")
			if f.Expr != "col1" || f.Alias != "c1" {
				t.Errorf("expected expr=col1 alias=c1, got %+v", f)
			}
		})

		t.Run("TwoArgsExprAlias", func(t *testing.T) {
			f := field.NewField("col1", "alias1")
			if f.Expr != "col1" || f.Alias != "alias1" {
				t.Errorf("expected expr=col1 alias=alias1, got %+v", f)
			}
		})

		t.Run("TwoArgsAliasUnsupportedType", func(t *testing.T) {
			f := field.NewField("col1", 123)
			if f.Error() == nil || !strings.Contains(f.Error().Error(), "input type unsupported: int") {
				t.Errorf("expected unsupported alias type error, got %+v", f.Error())
			}
		})

		t.Run("ThreeArgsValid", func(t *testing.T) {
			f := field.NewField("input1", "alias1", true)
			if f.Raw() != "input1 AS alias1" || f.Expr != "input1" || f.Alias != "alias1" || !f.IsRaw {
				t.Errorf("unexpected field %+v", f)
			}
		})

		t.Run("ThreeArgsAliasWrongType", func(t *testing.T) {
			f := field.NewField("input1", 123, true)
			if f.Error() == nil || !strings.Contains(f.Error().Error(), "input type unsupported: int") {
				t.Errorf("expected unsupported type error, got %+v", f.Error())
			}
		})

		t.Run("ThreeArgsIsRawWrongType", func(t *testing.T) {
			f := field.NewField("input1", "alias1", "yes")
			if f.Error() == nil || f.Error().Error() != "isRaw must be a bool" {
				t.Errorf("expected isRaw must be a bool error, got %+v", f.Error())
			}
		})

		t.Run("InvalidSignatureFourArgs", func(t *testing.T) {
			// expr, alias, isRaw
			f1 := field.NewField("col1", "alias1", true, "yes")
			if f1.Error() == nil || f1.Error().Error() != "invalid NewField signature" {
				t.Errorf("expected invalid NewField signature error, got %+v", f1.Error())
			}
		})

	})

	t.Run("Methods", func(t *testing.T) {
		t.Run("Clone", func(t *testing.T) {
			t.Run("Success", func(t *testing.T) {
				orig := &field.Field{Input: "i", Expr: "e", Alias: "a", IsRaw: true}
				cl := orig.Clone()
				if cl == orig {
					t.Fatal("expected different pointer")
				}
				if *cl != *orig {
					t.Errorf("clone differs: got %+v, want %+v", *cl, *orig)
				}
			})

			t.Run("NilReceiver", func(t *testing.T) {
				var f *field.Field = nil
				got := f.Clone()
				if got != nil {
					t.Errorf("cloned field = %+v, want %+v", got, f)
				}
			})
		})

		t.Run("IsAliased", func(t *testing.T) {
			cases := []struct {
				alias string
				want  bool
			}{
				{"", false},
				{"  ", false},
				{"alias", true},
				{"ALIAS", true},
			}
			for _, tc := range cases {
				f := field.Field{Alias: tc.alias}
				got := f.IsAliased()
				if got != tc.want {
					t.Errorf("IsAliased() alias=%q got %v, want %v", tc.alias, got, tc.want)
				}
			}
		})

		t.Run("IsErrored", func(t *testing.T) {
			f := field.Field{Expr: "field", Alias: "f"}
			if f.IsErrored() {
				t.Error("expected IsErrored false when Error is nil")
			}
			f.SetError(errors.New("some error"))
			if !f.IsErrored() {
				t.Error("expected IsErrored true when Error set")
			}
		})

		t.Run("IsValid", func(t *testing.T) {
			f := field.Field{Expr: "field", Alias: "f"}
			if !f.IsValid() {
				t.Error("expected IsValid true when no Error and Expr non-empty")
			}
			f.SetError(errors.New("some error"))
			if f.IsValid() {
				t.Error("expected IsValid false when Error set")
			}
			f = field.Field{Expr: "  "}
			if f.IsValid() {
				t.Error("expected IsValid false when Expr is empty")
			}
			f = field.Field{Expr: "!@#$%^&*()"}
			if f.IsValid() {
				t.Error("expected IsValid false when derived Name is empty")
			}
		})

		t.Run("Name", func(t *testing.T) {
			cases := []struct {
				alias string
				expr  string
				want  string
			}{
				{"", "fieldName", "fieldname"},
				{"ALIAS", "ignored", "alias"},
				{"Id", "some_expr", "id"},
				{"", "Expr_With_123", "exprwith123"},
			}
			for _, tc := range cases {
				f := field.Field{Alias: tc.alias, Expr: tc.expr}
				got := f.Name()
				if got != tc.want {
					t.Errorf("Name() alias=%q expr=%q got %q, want %q", tc.alias, tc.expr, got, tc.want)
				}
			}
		})

		t.Run("Raw", func(t *testing.T) {
			f := field.Field{Expr: "field"}
			if got := f.Raw(); got != "field" {
				t.Errorf("Raw() without alias got %q, want %q", got, "field")
			}
			f.Alias = "alias"
			if got := f.Raw(); got != "field AS alias" {
				t.Errorf("Raw() with alias got %q, want %q", got, "field AS alias")
			}
		})

		t.Run("Render", func(t *testing.T) {
			f := field.Field{Expr: "LOWER(x)"}
			if got, want := f.Render(), "LOWER(x)"; got != want {
				t.Errorf("Render() got %q, want %q", got, want)
			}
			f = field.Field{Expr: "LOWER(x)", Alias: "id"}
			if got, want := f.Render(), "LOWER(x) AS id"; got != want {
				t.Errorf("Render() got %q, want %q", got, want)
			}
			f = field.Field{Expr: "created_at", Alias: "  ts  "}
			if got, want := f.Render(), "created_at AS ts"; got != want {
				t.Errorf("Render() trims alias spaces got %q, want %q", got, want)
			}
		})

		t.Run("String", func(t *testing.T) {
			f := field.Field{Expr: "field", Alias: ""}
			got := f.String()
			if !strings.HasPrefix(got, "✅ Field(") {
				t.Errorf("String() no alias got %q, want prefix %q", got, "✅ Field(")
			}
			f.Alias = "Alias"
			got = f.String()
			if !strings.Contains(got, "field AS Alias") {
				t.Errorf("String() with alias got %q, want alias: true", got)
			}
			f.SetError(errors.New("some error"))
			got = f.String()
			if !strings.Contains(got, "⛔️") || !strings.Contains(got, "some error") {
				t.Errorf("String() errored got %q, want icon and error message", got)
			}
		})
	})
}
