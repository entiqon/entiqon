package condition_test

import (
	"testing"

	"github.com/entiqon/db/token/types/condition"
)

func TestType(t *testing.T) {
	t.Run("IsValid", func(t *testing.T) {
		cases := []struct {
			typ     condition.Type
			isValid bool
		}{
			{condition.Invalid, false},
			{condition.Single, true},
			{condition.And, true},
			{condition.Or, true},
			{condition.Type(99), false},
		}

		for _, c := range cases {
			if got := c.typ.IsValid(); got != c.isValid {
				t.Errorf("IsValid(%v) = %v, want %v", c.typ, got, c.isValid)
			}
		}
	})

	t.Run("ParseFrom", func(t *testing.T) {
		t.Run("FromType", func(t *testing.T) {
			if got := condition.ParseFrom(condition.And); got != condition.And {
				t.Errorf("expected And, got %v", got)
			}
		})

		t.Run("FromInt", func(t *testing.T) {
			if got := condition.ParseFrom(3); got != condition.Or {
				t.Errorf("expected And, got %v", got)
			}
			if got := condition.ParseFrom(99); got != condition.Invalid {
				t.Errorf("expected Invalid, got %v", got)
			}
		})

		t.Run("FromString", func(t *testing.T) {
			if got := condition.ParseFrom("and"); got != condition.And {
				t.Errorf("expected And, got %v", got)
			}
			if got := condition.ParseFrom(" Or "); got != condition.Or {
				t.Errorf("expected Or, got %v", got)
			}
			if got := condition.ParseFrom("foobar"); got != condition.Invalid {
				t.Errorf("expected Invalid, got %v", got)
			}
		})

		t.Run("FromUnsupported", func(t *testing.T) {
			if got := condition.ParseFrom(3.14); got != condition.Invalid {
				t.Errorf("expected Invalid, got %v", got)
			}
		})
	})

	t.Run("String", func(t *testing.T) {
		cases := []struct {
			typ  condition.Type
			want string
		}{
			{condition.Invalid, "Invalid"},
			{condition.Single, ""},
			{condition.And, "AND"},
			{condition.Or, "OR"},
			{condition.Type(99), "Invalid"},
		}

		for _, c := range cases {
			if got := c.typ.String(); got != c.want {
				t.Errorf("String(%v) = %q, want %q", c.typ, got, c.want)
			}
		}
	})
}
