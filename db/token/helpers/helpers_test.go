package helpers_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/entiqon/db/token/field"
	"github.com/entiqon/db/token/helpers"
	"github.com/entiqon/db/token/types/operator"
)

func TestHelpers(t *testing.T) {
	t.Run("ValidationConsistency", func(t *testing.T) {
		// Shared invalid cases across identifiers and aliases.
		invalidCases := map[string]string{
			"":        "empty",
			"123abc":  "digit",
			"9":       "digit",
			"-name":   "syntax",
			"$var":    "syntax",
			"user id": "syntax",
			"name!":   "syntax",
			"中":       "syntax",
			"café":    "syntax",
			"mañana":  "syntax",
			"niño":    "syntax",
		}

		// Reserved keywords only apply to aliases.
		reserved := helpers.ReservedKeywords()

		tests := []struct {
			name     string
			validate func(string) error
			isValid  func(string) bool
			valid    []string
			invalid  map[string]string
			reserved []string
		}{
			{
				name:     "Identifier",
				validate: helpers.ValidateIdentifier,
				isValid:  helpers.IsValidIdentifier,
				valid:    []string{"id", "user_id", "_col123", "U1", "XYZ"},
				invalid:  invalidCases, // all shared invalid cases
			},
			{
				name:     "Alias",
				validate: helpers.ValidateAlias,
				isValid:  helpers.IsValidAlias,
				valid:    []string{"id", "alias1", "_col", "total_sum", "user123"},
				invalid:  invalidCases, // same base invalid cases
				reserved: reserved,     // plus reserved keywords
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name+"/Valid", func(t *testing.T) {
				for _, s := range tt.valid {
					if !tt.isValid(s) {
						t.Errorf("expected %s %q to be valid", tt.name, s)
					}
					if err := tt.validate(s); err != nil {
						t.Errorf("expected %s %q valid, got error: %v", tt.name, s, err)
					}
				}
			})

			t.Run(tt.name+"/Invalid", func(t *testing.T) {
				for s, expected := range tt.invalid {
					if tt.isValid(s) {
						t.Errorf("expected %s %q to be invalid", tt.name, s)
					}
					err := tt.validate(s)
					if err == nil {
						t.Errorf("expected error for %s %q, got nil", tt.name, s)
						continue
					}
					if !strings.Contains(strings.ToLower(err.Error()), expected) {
						t.Errorf("%s %q: expected error about %s, got %v", tt.name, s, expected, err)
					}
				}
			})

			if len(tt.reserved) > 0 {
				t.Run(tt.name+"/Reserved", func(t *testing.T) {
					for _, kw := range tt.reserved {
						if tt.isValid(kw) {
							t.Errorf("expected %s %q to be invalid (reserved)", tt.name, kw)
						}
						err := tt.validate(kw)
						if err == nil {
							t.Errorf("expected reserved error for %s %q, got nil", tt.name, kw)
						}
						if !strings.Contains(strings.ToLower(err.Error()), "reserved") {
							t.Errorf("%s %q: expected reserved keyword error, got %v", tt.name, kw, err)
						}
					}
				})
			}
		}
	})

	t.Run("TrailingAlias", func(t *testing.T) {
		tests := []struct {
			name     string
			expr     string
			expected string
			valid    bool
		}{
			// ✅ Valid trailing aliases
			{"SubqueryWithAlias", "(SELECT * FROM users) u", "u", true},
			{"ExpressionWithAlias", "(price * quantity) total", "total", true},

			// ❌ Invalid cases
			{"ExplicitAS", "(SELECT * FROM users) AS u", "", false},
			{"NoAlias", "(price * quantity)", "", false},
			{"OperatorCase", "(price * quantity) * discount", "", false},
			{"ReservedAlias", "(SELECT * FROM users) SELECT", "", false},
			{"InvalidSyntaxAlias", "(SELECT * FROM users) 123alias", "", false},
			{"NonASCIIAlias", "(SELECT * FROM users) café", "", false},

			// ⚠️ Edge cases: not enough tokens to form an alias
			{"SingleToken", "users", "", false},
			{"EmptyExpr", "", "", false},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				alias, err := helpers.ValidateTrailingAlias(tt.expr)
				if tt.valid {
					if err != nil {
						t.Errorf("expected valid alias for %q, got error: %v", tt.expr, err)
					}
					if alias != tt.expected {
						t.Errorf("expected alias %q, got %q", tt.expected, alias)
					}
					if !helpers.HasTrailingAlias(tt.expr) {
						t.Errorf("expected HasTrailingAlias true for %q", tt.expr)
					}
				} else {
					if err == nil {
						t.Errorf("expected error for %q, got nil", tt.expr)
					}
					if helpers.HasTrailingAlias(tt.expr) {
						t.Errorf("expected HasTrailingAlias false for %q", tt.expr)
					}
				}
			})
		}
	})

	t.Run("GenerateAlias", func(t *testing.T) {
		tests := []struct {
			name   string
			prefix string
			expr   string
		}{
			{"Literal", "lt", "42"},
			{"Function", "fn", "SUM(price)"},
			{"Aggregate", "ag", "COUNT(*)"},
			{"Computed", "cp", "(price * quantity)"},
			{"Subquery", "sq", "(SELECT * FROM users)"},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				alias1 := helpers.GenerateAlias(tt.prefix, tt.expr)
				alias2 := helpers.GenerateAlias(tt.prefix, tt.expr)

				// Deterministic
				if alias1 != alias2 {
					t.Errorf("expected deterministic alias, got %q vs %q", alias1, alias2)
				}

				// Must start with the prefix
				if gotPrefix := alias1[:len(tt.prefix)]; gotPrefix != tt.prefix {
					t.Errorf("expected alias to start with prefix %q, got %q", tt.prefix, gotPrefix)
				}

				// Must have underscore + hash part
				if !strings.Contains(alias1, "_") {
					t.Errorf("expected alias %q to contain an underscore", alias1)
				}
				if len(alias1) <= len(tt.prefix)+1 {
					t.Errorf("expected alias %q to have prefix and hash part", alias1)
				}
			})
		}

		t.Run("DifferentExpressions", func(t *testing.T) {
			a1 := helpers.GenerateAlias("fn", "SUM(price)")
			a2 := helpers.GenerateAlias("fn", "AVG(price)")
			if a1 == a2 {
				t.Errorf("expected different aliases for different expressions, got %q", a1)
			}
		})
	})

	t.Run("ValidateWildcard", func(t *testing.T) {
		tests := []struct {
			name  string
			expr  string
			alias string
			valid bool
		}{
			{"BareStarNoAlias", "*", "", true},         // allowed
			{"BareStarWithAlias", "*", "total", false}, // invalid
			{"NotWildcard", "id", "alias", true},       // not a wildcard, ignore
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				err := helpers.ValidateWildcard(tt.expr, tt.alias)
				if tt.valid {
					if err != nil {
						t.Errorf("expected valid for expr=%q alias=%q, got error: %v", tt.expr, tt.alias, err)
					}
				} else {
					if err == nil {
						t.Errorf("expected error for expr=%q alias=%q, got nil", tt.expr, tt.alias)
					}
				}
			})
		}
	})

	t.Run("ResolveExpressionType", func(t *testing.T) {
		tests := []struct {
			name string
			expr string
			want string
		}{
			{"Empty", "", "Invalid"},
			{"Subquery", "(SELECT * FROM users)", "Subquery"},
			{"Computed", "(a+b)", "Computed"},
			{"AggregateSUM", "SUM(qty)", "Aggregate"},
			{"AggregateCOUNT", "COUNT(*)", "Aggregate"},
			{"Function", "JSON_EXTRACT(data, '$.id')", "Function"},
			{"LiteralString", "'abc'", "Literal"},
			{"LiteralNumber", "42", "Literal"},
			{"Identifier", "users", "Identifier"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := helpers.ResolveExpressionType(tt.expr)
				if got.String() != tt.want {
					t.Errorf("ResolveExpressionType(%q) = %v, want %v", tt.expr, got, tt.want)
				}
			})
		}
	})

	t.Run("ResolveExpression", func(t *testing.T) {
		tests := []struct {
			name       string
			input      string
			allowAlias bool
			wantKind   string
			wantExpr   string
			wantAlias  string
			wantErr    bool
		}{
			//
			// === Invalid ===
			{"EmptyInput", "", true, "Invalid", "", "", true},
			{"GarbageInput", "foo bar baz qux", true, "Invalid", "", "", true},

			//=== Identifiers ===
			{"Identifier", "field", true, "Identifier", "field", "", false},
			{"IdentifierWithAlias", "field alias", true, "Identifier", "field", "alias", false},
			{"IdentifierWithNotAllowAlias", "field alias", false, "Identifier", "field", "alias", true},
			{"IdentifierWithInvalidAlias", "field 123invalid", true, "Invalid", "field", "123invalid", true},
			{"IdentifierWithASAlias", "field AS alias", true, "Identifier", "field", "alias", false},
			{"IdentifierWithASAliasNotAllowAlias", "field AS alias", false, "Identifier", "field", "alias", true},
			{"IdentifierWithASAliasInvalidAlias", "field AS 123invalid", true, "Invalid", "field", "123invalid", true},
			{"IdentifierInvalidForm", "field alias extra", true, "Invalid", "", "", true},
			{"IdentifierTooManyTokens", "field AS alias extra", true, "Invalid", "", "", true},

			// === Subqueries ===
			{"SubqueryNoAlias", "(SELECT * FROM users)", true, "Subquery", "(SELECT * FROM users)", "", false},
			{"SubqueryNoAlias", "(SELECT * FROM users", true, "Subquery", "(SELECT * FROM users", "", true},
			{"SubqueryWithAlias", "(SELECT * FROM users) u", true, "Subquery", "(SELECT * FROM users)", "u", false},
			{"SubqueryWithAliasNotAllowAlias", "(SELECT * FROM users) u", false, "Subquery", "(SELECT * FROM users)", "u", true},
			{"SubqueryWithAliasInvalidAlias", "(SELECT * FROM users) 123u", true, "Subquery", "(SELECT * FROM users)", "123u", true},
			{"SubqueryWithASAlias", "(SELECT * FROM users) AS u", true, "Subquery", "(SELECT * FROM users)", "u", false},
			{"SubqueryWithASAliasNotAllowAlias", "(SELECT * FROM users) AS u", false, "Subquery", "(SELECT * FROM users)", "u", true},
			{"SubqueryWithASAliasInvalidAlias", "(SELECT * FROM users) AS 123u", true, "Subquery", "(SELECT * FROM users)", "123u", true},
			{"BareSelectRejected", "SELECT * FROM users", true, "Invalid", "", "", true},
			{"SubqueryInvalidAlias", "(SELECT * FROM users) AS abc 123", true, "Invalid", "", "", true},

			// === Computed ===
			{"ComputedNoAlias", "(price * qty)", true, "Computed", "(price * qty)", "", false},
			{"ComputedWithAlias", "(price * qty) total", true, "Computed", "(price * qty)", "total", false},
			{"ComputedWithASAlias", "(price * qty) AS total", true, "Computed", "(price * qty)", "total", false},
			{"ComputedWithASAlias", "(price * qty) AS abc 123", true, "Computed", "", "", true},
			{"ComputedBareRejected", "price * qty", true, "Invalid", "", "", true},

			// === Aggregates ===
			{"AggregateCount", "COUNT(*)", true, "Aggregate", "COUNT(*)", "", false},
			{"AggregateCountWithAlias", "COUNT(*) total", true, "Aggregate", "COUNT(*)", "total", false},
			{"AggregateCountWithASAlias", "COUNT(*) AS total", true, "Aggregate", "COUNT(*)", "total", false},
			{"AggregateSum", "SUM(price * qty)", true, "Aggregate", "SUM(price * qty)", "", false},
			{"AggregateSumWithAlias", "SUM(price * qty) total", true, "Aggregate", "SUM(price * qty)", "total", false},
			{"AggregateSumWithASAlias", "SUM(price * qty) AS total", true, "Aggregate", "SUM(price * qty)", "total", false},
			{"ComputedInvalidAlias", "SUM(price * qty) AS abc 123", true, "Computed", "", "", true},

			// === Functions ===
			{"FunctionNoAlias", "JSON_EXTRACT(data,'$.id')", true, "Function", "JSON_EXTRACT(data,'$.id')", "", false},
			{"FunctionWithAlias", "LOWER(name) alias", true, "Function", "LOWER(name)", "alias", false},
			{"FunctionWithASAlias", "LOWER(name) AS alias", true, "Function", "LOWER(name)", "alias", false},
			{"FunctionWithReservedAlias", "LOWER(name) AS SELECT", true, "Invalid", "", "", true},
			{"ComputedInvalidAlias", "LOWER(name) AS abc 123", true, "Computed", "", "", true},

			// === Literals ===
			{"LiteralString", "'abc'", true, "Literal", "'abc'", "", false},
			{"LiteralStringWithAlias", "'abc' val", true, "Literal", "'abc'", "val", false},
			{"LiteralStringWithASAlias", "'abc' AS val", true, "Literal", "'abc'", "val", false},
			{"LiteralStringWithReservedAlias", "'abc' AS SELECT", true, "Invalid", "", "", true},
			{"LiteralInvalidAlias", "'abc' AS abc 123", true, "Computed", "", "", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				kind, expr, alias, err := helpers.ResolveExpression(tt.input, tt.allowAlias)

				if tt.wantErr {
					if err == nil {
						t.Errorf("%s: expected error, got nil (kind=%v expr=%q alias=%q)", tt.name, kind, expr, alias)
					}
					return
				}

				if err != nil {
					t.Errorf("%s: unexpected error: %v", tt.name, err)
					return
				}

				if kind.String() != tt.wantKind {
					t.Errorf("%s: expected kind=%q, got %q", tt.name, tt.wantKind, kind.String())
				}
				if expr != tt.wantExpr {
					t.Errorf("%s: expected expr=%q, got %q", tt.name, tt.wantExpr, expr)
				}
				if alias != tt.wantAlias {
					t.Errorf("%s: expected alias=%q, got %q", tt.name, tt.wantAlias, alias)
				}
			})
		}
	})
}

func TestValidateType(t *testing.T) {
	err := helpers.ValidateType("string")
	if err != nil {
		t.Error("expected error, got nil")
	}
	err = helpers.ValidateType(123456)
	if err == nil {
		t.Error("expected error, got nil")
	}
	err = helpers.ValidateType(field.New("id"))
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCondition(t *testing.T) {
	t.Run("ResolveCondition", func(t *testing.T) {
		tests := []struct {
			name      string
			input     string
			wantOp    operator.Type
			wantField string
			wantVal   any
			wantErr   bool
		}{
			{
				name:      "EqualCondition",
				input:     "id = 1",
				wantOp:    operator.Equal,
				wantField: "id",
				wantVal:   1,
			},
			{
				name:      "ValidCondition",
				input:     "id",
				wantOp:    operator.Equal,
				wantField: "id",
				wantVal:   nil,
				wantErr:   false,
			},
			{
				name:      "InvalidCondition",
				input:     "",
				wantOp:    operator.Equal,
				wantField: "",
				wantVal:   nil,
				wantErr:   true,
			},
			{
				name:      "InList",
				input:     "lastname IN ('smith','brown')",
				wantOp:    operator.In,
				wantField: "lastname",
				wantVal:   []any{"smith", "brown"},
			},
			{
				name:      "InInvalid",
				input:     "lastname IN ()",
				wantOp:    operator.In,
				wantField: "lastname",
				wantVal:   []any{},
				wantErr:   true,
			},
			{
				name:      "BetweenRange",
				input:     "price BETWEEN 1 AND 10",
				wantOp:    operator.Between,
				wantField: "price",
				wantVal:   []any{1, 10},
			},
			{
				name:      "BetweenInvalid",
				input:     "price BETWEEN 1 AND",
				wantOp:    operator.Invalid,
				wantField: "price",
				wantVal:   nil,
				wantErr:   true,
			},
			{
				name:      "IsNullCondition",
				input:     "deleted_at IS NULL",
				wantOp:    operator.IsNull,
				wantField: "deleted_at",
				wantVal:   nil,
			},
			{
				name:      "InvalidOperator",
				input:     "id ++ 1",
				wantOp:    operator.Invalid,
				wantField: "",
				wantErr:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				f, op, val, err := helpers.ResolveCondition(tt.input)
				if (err != nil) != tt.wantErr {
					t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
				}
				if !tt.wantErr {
					if op != tt.wantOp {
						t.Errorf("expected op=%v, got %v", tt.wantOp, op)
					}
					if f != tt.wantField {
						t.Errorf("expected field=%q, got %q", tt.wantField, f)
					}
					if !reflect.DeepEqual(val, tt.wantVal) {
						t.Errorf("expected val=%v, got %v", tt.wantVal, val)
					}
				}
			})
		}
	})

	t.Run("IsValidSlice", func(t *testing.T) {
		tests := []struct {
			name     string
			op       operator.Type
			value    any
			expected bool
		}{
			// IN operator
			{"In_NonEmptyInts", operator.In, []int{1, 2, 3}, true},
			{"In_EmptyInts", operator.In, []int{}, false},
			{"In_EmptyInt64s", operator.In, []int64{}, false},
			{"In_EmptyFloat64s", operator.In, []float64{}, false},
			{"In_NonEmptyStrings", operator.In, []string{"a", "b"}, true},
			{"In_EmptyStrings", operator.In, []string{}, false},
			{"In_Nil", operator.In, []struct{}{}, false},
			{"In_Nil", operator.In, nil, false},

			// NOT IN operator
			{"NotIn_NonEmpty", operator.NotIn, []any{"x"}, true},
			{"NotIn_Empty", operator.NotIn, []any{}, false},

			// BETWEEN operator
			{"Between_ExactlyTwo", operator.Between, []int{1, 10}, true},
			{"Between_OneValue", operator.Between, []int{1}, false},
			{"Between_ThreeValues", operator.Between, []int{1, 2, 3}, false},
			{"Between_Nil", operator.Between, nil, false},

			// Other operators should always fail
			{"Equal_WithSlice", operator.Equal, []int{1, 2}, false},
			{"InvalidOperator", operator.Invalid, []int{1}, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ok := helpers.IsValidSlice(tt.op, tt.value)
				if ok != tt.expected {
					t.Errorf("IsValidSlice(%v, %v) = %v, want %v",
						tt.op, tt.value, ok, tt.expected)
				}
			})
		}
	})

	t.Run("ToParamKey", func(t *testing.T) {
		cases := map[string]string{
			`users.id`:         "users_id",
			`"last name"`:      "last_name",
			`u."last name"`:    "u_last_name",
			``:                 "field",
			`some-weird$field`: "some_weird_field",
		}
		for in, want := range cases {
			if got := helpers.ToParamKey(in); got != want {
				t.Errorf("ToParamKey(%q)=%q, want %q", in, got, want)
			}
		}
	})
}
