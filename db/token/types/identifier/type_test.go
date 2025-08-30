package identifier_test

import (
	"testing"

	"github.com/entiqon/entiqon/db/token/types/identifier"
)

func TestType(t *testing.T) {
	t.Run("Alias", func(t *testing.T) {
		cases := []struct {
			typ  identifier.Type
			want string
		}{
			{identifier.Identifier, "id"},
			{identifier.Literal, "lt"},
			{identifier.Function, "fn"},
			{identifier.Aggregate, "ag"},
			{identifier.Computed, "cp"},
			{identifier.Subquery, "sq"},
			{identifier.Invalid, "ex"},
		}

		for _, c := range cases {
			if got := c.typ.Alias(); got != c.want {
				t.Errorf("Alias(%v) = %q, want %q", c.typ, got, c.want)
			}
		}
	})

	t.Run("IsValid", func(t *testing.T) {
		cases := []struct {
			typ  identifier.Type
			want bool
		}{
			{identifier.Invalid, false},
			{identifier.Identifier, true},
			{identifier.Literal, true},
			{identifier.Function, true},
			{identifier.Aggregate, true},
			{identifier.Computed, true},
			{identifier.Subquery, true},
			{identifier.Type(99), false}, // out of range
		}

		for _, c := range cases {
			if got := c.typ.IsValid(); got != c.want {
				t.Errorf("IsValid(%v) = %v, want %v", c.typ, got, c.want)
			}
		}
	})

	t.Run("ParseFrom", func(t *testing.T) {
		t.Run("FromInt", func(t *testing.T) {
			if got := identifier.Invalid.ParseFrom(int(identifier.Identifier)); got != identifier.Identifier {
				t.Errorf("ParseFrom(int Identifier) = %v, want Identifier", got)
			}
			if got := identifier.Invalid.ParseFrom(99); got != identifier.Invalid {
				t.Errorf("ParseFrom(99) = %v, want Invalid", got)
			}
		})

		t.Run("FromString", func(t *testing.T) {
			cases := []struct {
				in   string
				want identifier.Type
			}{
				{"identifier", identifier.Identifier},
				{"IDENTIFIER", identifier.Identifier},
				{"literal", identifier.Literal},
				{"function", identifier.Function},
				{"aggregate", identifier.Aggregate},
				{"computed", identifier.Computed},
				{"subquery", identifier.Subquery},
				{"unknown", identifier.Invalid},
				{"", identifier.Invalid},
			}

			for _, c := range cases {
				if got := identifier.Invalid.ParseFrom(c.in); got != c.want {
					t.Errorf("ParseFrom(%q) = %v, want %v", c.in, got, c.want)
				}
			}
		})

		t.Run("FromType", func(t *testing.T) {
			if got := identifier.Invalid.ParseFrom(identifier.Computed); got != identifier.Computed {
				t.Errorf("ParseFrom(Type Computed) = %v, want Computed", got)
			}
		})

		t.Run("FromUnsupported", func(t *testing.T) {
			if got := identifier.Invalid.ParseFrom([]string{"bad"}); got != identifier.Invalid {
				t.Errorf("ParseFrom([]string) = %v, want Invalid", got)
			}
		})
	})

	t.Run("String", func(t *testing.T) {
		cases := []struct {
			typ  identifier.Type
			want string
		}{
			{identifier.Invalid, "Invalid"},
			{identifier.Identifier, "Identifier"},
			{identifier.Computed, "Computed"},
			{identifier.Literal, "Literal"},
			{identifier.Subquery, "Subquery"},
			{identifier.Function, "Function"},
			{identifier.Aggregate, "Aggregate"},
			{identifier.Type(99), "Unknown"},
		}

		for _, c := range cases {
			if got := c.typ.String(); got != c.want {
				t.Errorf("String(%v) = %q, want %q", c.typ, got, c.want)
			}
		}
	})
}
