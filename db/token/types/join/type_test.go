package join_test

import (
	"testing"

	"github.com/entiqon/db/token/types/join"
)

func TestType(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		cases := []struct {
			typ  join.Type
			want string
		}{
			{join.Inner, "INNER JOIN"},
			{join.Left, "LEFT JOIN"},
			{join.Right, "RIGHT JOIN"},
			{join.Full, "FULL JOIN"},
			{join.Cross, "CROSS JOIN"},
			{join.Natural, "NATURAL JOIN"},
			{join.Invalid, "INVALID"},
		}

		for _, c := range cases {
			if got := c.typ.String(); got != c.want {
				t.Errorf("Type(%v).String() = %q, want %q", c.typ, got, c.want)
			}
		}
	})

	t.Run("IsValid", func(t *testing.T) {
		cases := []struct {
			typ  join.Type
			want bool
		}{
			{join.Inner, true},
			{join.Cross, true},
			{join.Natural, true},
			{join.Invalid, false},
			{join.Type(99), false},
		}

		for _, c := range cases {
			if got := c.typ.IsValid(); got != c.want {
				t.Errorf("Type(%v).IsValid() = %v, want %v", c.typ, got, c.want)
			}
		}
	})

	t.Run("ParseFrom", func(t *testing.T) {
		cases := []struct {
			in   string
			want join.Type
		}{
			{"inner", join.Inner},
			{"LEFT JOIN", join.Left},
			{"right", join.Right},
			{"FULL JOIN", join.Full},
			{"cross", join.Cross},
			{"NATURAL JOIN", join.Natural},
			{"", join.Invalid},
			{"weird", join.Invalid},
		}

		for _, c := range cases {
			if got := join.ParseFrom(c.in); got != c.want {
				t.Errorf("ParseFrom(%q) = %v, want %v", c.in, got, c.want)
			}
		}
	})
}
