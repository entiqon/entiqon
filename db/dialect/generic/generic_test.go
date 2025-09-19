package generic_test

import (
	"testing"
	"time"

	"github.com/entiqon/db/dialect/generic"
)

func TestGenericDialect(t *testing.T) {
	d := generic.New()

	t.Run("Name", func(t *testing.T) {
		if got, want := d.Name(), "generic"; got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("Options", func(t *testing.T) {
		opts := d.Options()
		if opts.Name != "generic" {
			t.Errorf("unexpected Options.Name = %q", opts.Name)
		}
		if opts.PlaceholderStyle != "?" {
			t.Errorf("unexpected PlaceholderStyle = %q", opts.PlaceholderStyle)
		}
	})

	t.Run("QuoteIdentifier", func(t *testing.T) {
		cases := []struct {
			in   string
			want string
		}{
			{"users", "users"},
			{"UserData", `"UserData"`},
			{"order items", `"order items"`},
			{"", ""},
		}

		for _, c := range cases {
			if got := d.QuoteIdentifier(c.in); got != c.want {
				t.Errorf("QuoteIdentifier(%q) = %q, want %q", c.in, got, c.want)
			}
		}
	})

	t.Run("QuoteLiteral", func(t *testing.T) {
		now := time.Date(2025, 9, 19, 3, 30, 0, 0, time.UTC)

		cases := []struct {
			in   any
			want string
		}{
			{nil, "NULL"},
			{"O'Reilly", "'O''Reilly'"},
			{true, "TRUE"},
			{false, "FALSE"},
			{42, "42"},
			{uint(99), "99"},
			{3.14, "3.14"},        // float64
			{float32(2.5), "2.5"}, // float32 → toFloat64
			{now, "'2025-09-19 03:30:00'"},
		}

		for _, c := range cases {
			if got := d.QuoteLiteral(c.in); got != c.want {
				t.Errorf("QuoteLiteral(%v) = %q, want %q", c.in, got, c.want)
			}
		}

		// Default branch: unhandled type
		got := d.QuoteLiteral(struct{ Foo string }{"bar"})
		if got != "'{bar}'" { // fmt.Sprintf("%v") on struct{Foo string}{"bar"} → {bar}
			t.Errorf("QuoteLiteral(struct) = %q, want %q", got, "'{bar}'")
		}

		// Default branch of toFloat64: unhandled type
		got = d.QuoteLiteral([]byte{1, 2})
		if got != "'[1 2]'" {
			t.Errorf("QuoteLiteral([]byte) = %q, want %q", got, "'[1 2]'")
		}
	})

	t.Run("PaginationSyntax", func(t *testing.T) {
		cases := []struct {
			limit, offset int
			want          string
		}{
			{0, 0, ""},
			{10, 0, " LIMIT 10"},
			{0, 20, " OFFSET 20"},
			{10, 20, " LIMIT 10 OFFSET 20"},
		}

		for _, c := range cases {
			if got := d.PaginationSyntax(c.limit, c.offset); got != c.want {
				t.Errorf("PaginationSyntax(%d,%d) = %q, want %q",
					c.limit, c.offset, got, c.want)
			}
		}
	})

	t.Run("Placeholder", func(t *testing.T) {
		if got := d.Placeholder(1); got != "?" {
			t.Errorf("Placeholder(1) = %q, want '?'", got)
		}
		if got := d.Placeholder(42); got != "?" {
			t.Errorf("Placeholder(42) = %q, want '?'", got)
		}
	})
}
