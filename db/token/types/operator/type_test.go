package operator_test

import (
	"reflect"
	"testing"

	"github.com/entiqon/db/token/types/operator"
)

// TestTypeMethods verifies Alias, IsValid, and String basic behavior.
func TestTypeMethods(t *testing.T) {
	t.Run("Alias", func(t *testing.T) {
		if operator.Equal.Alias() != "eq" {
			t.Fatalf("expected alias 'eq', got %q", operator.Equal.Alias())
		}
		if operator.NotIn.Alias() != "nin" {
			t.Fatalf("expected alias 'nin', got %q", operator.NotIn.Alias())
		}
	})

	t.Run("IsValid", func(t *testing.T) {
		if operator.Invalid.IsValid() {
			t.Fatalf("Invalid should not be valid")
		}
		if !operator.GreaterThan.IsValid() || !operator.IsNull.IsValid() {
			t.Fatalf("expected valid operators to be valid")
		}
	})

	t.Run("String", func(t *testing.T) {
		if operator.Equal.String() != "=" {
			t.Fatalf("expected '=', got %q", operator.Equal.String())
		}
		if operator.NotEqual.String() != "!=" {
			t.Fatalf("expected '!=', got %q", operator.NotEqual.String())
		}
		if operator.IsDistinctFrom.String() != "IS DISTINCT FROM" {
			t.Fatalf("expected 'IS DISTINCT FROM', got %q", operator.IsDistinctFrom.String())
		}
		if operator.NotIsDistinctFrom.String() != "IS NOT DISTINCT FROM" {
			t.Fatalf("expected 'IS NOT DISTINCT FROM', got %q", operator.NotIsDistinctFrom.String())
		}
	})
}

// TestParseFrom covers symbols, words, aliases, []byte inputs, and invalids.
func TestParseFrom(t *testing.T) {
	t.Run("Symbols", func(t *testing.T) {
		cases := map[string]operator.Type{
			"=":  operator.Equal,
			"!=": operator.NotEqual,
			"<>": operator.NotEqual,
			">=": operator.GreaterThanOrEqual,
			"<=": operator.LessThanOrEqual,
			">":  operator.GreaterThan,
			"<":  operator.LessThan,
		}
		for in, want := range cases {
			if got := operator.ParseFrom(in); got != want {
				t.Fatalf("ParseFrom(%q) => %v, want %v", in, got, want)
			}
		}
	})

	t.Run("Words", func(t *testing.T) {
		cases := map[string]operator.Type{
			"in":                         operator.In,
			"NOT IN":                     operator.NotIn,
			"between":                    operator.Between,
			"LIKE":                       operator.Like,
			"nOt    LiKe":                operator.NotLike,
			" is   null   ":              operator.IsNull,
			" Is  NoT   NuLl ":           operator.IsNotNull,
			"is distinct from":           operator.IsDistinctFrom,
			"IS   NOT   DISTINCT   FROM": operator.NotIsDistinctFrom,
		}
		for in, want := range cases {
			if got := operator.ParseFrom(in); got != want {
				t.Fatalf("ParseFrom(%q) => %v, want %v", in, got, want)
			}
		}
	})

	t.Run("Aliases", func(t *testing.T) {
		cases := map[string]operator.Type{
			"eq":          operator.Equal,
			"neq":         operator.NotEqual,
			"gt":          operator.GreaterThan,
			"gte":         operator.GreaterThanOrEqual,
			"lt":          operator.LessThan,
			"lte":         operator.LessThanOrEqual,
			"in":          operator.In,
			"nin":         operator.NotIn,
			"between":     operator.Between,
			"like":        operator.Like,
			"nlike":       operator.NotLike,
			"isnull":      operator.IsNull,
			"notnull":     operator.IsNotNull,
			"isdistinct":  operator.IsDistinctFrom,
			"notdistinct": operator.NotIsDistinctFrom,
		}
		for in, want := range cases {
			if got := operator.ParseFrom(in); got != want {
				t.Fatalf("ParseFrom(%q) => %v, want %v", in, got, want)
			}
		}
	})

	t.Run("Bytes", func(t *testing.T) {
		if got := operator.ParseFrom([]byte(" IN ")); got != operator.In {
			t.Fatalf("ParseFrom([]byte(\" IN \")) => %v, want %v", got, operator.In)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if got := operator.ParseFrom(false); got != operator.Invalid {
			t.Fatalf("ParseFrom(false) => %v, want Invalid", got)
		}
	})
}

// TestGetKnownOperators asserts deterministic order and copy semantics.
func TestGetKnownOperators(t *testing.T) {
	t.Run("DeterministicOrder", func(t *testing.T) {
		got := operator.GetKnownOperators()
		want := []string{
			"IS NOT DISTINCT FROM",
			"IS DISTINCT FROM",
			"IS NOT NULL",
			"NOT LIKE",
			"BETWEEN",
			"IS NULL",
			"NOT IN",
			"LIKE",
			"IN",
			"!=",
			">=",
			"<=",
			">",
			"<",
			"=",
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected order:\n got: %v\nwant: %v", got, want)
		}
	})

	t.Run("CopyOnReturn", func(t *testing.T) {
		a := operator.GetKnownOperators()
		b := operator.GetKnownOperators()
		if !reflect.DeepEqual(a, b) {
			t.Fatalf("fresh calls should be equal:\n a: %v\n b: %v", a, b)
		}
		// Mutate a, ensure b remains unchanged (function returns a copy).
		a[0] = "HACK"
		if reflect.DeepEqual(a, b) {
			t.Fatalf("expected slices to diverge after mutation of the first")
		}
	})
}

// TestRoundTrip ensures String(ParseFrom(s)) == s for canonical spellings.
func TestRoundTrip(t *testing.T) {
	known := operator.GetKnownOperators()
	for _, s := range known {
		got := operator.ParseFrom(s)
		if !got.IsValid() {
			t.Fatalf("ParseFrom(%q) returned Invalid", s)
		}
		if got.String() != s {
			t.Fatalf("round-trip mismatch: String(%v)=%q, want %q", got, got.String(), s)
		}
	}

	// Special check: "<>" parses to NotEqual but String() normalizes to "!="
	if operator.ParseFrom("<>").String() != "!=" {
		t.Fatalf("expected String(ParseFrom(\"<>\")) to be \"!=\"")
	}
}
