package driver_test

import (
	"strings"
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/driver/styling"
)

func TestResolveDialect(t *testing.T) {
	if got := driver.ResolveDialect("postgres").GetName(); got != "postgres" {
		t.Errorf("expected %q, got %q", "postgres", got)
	}
	if got := driver.ResolveDialect("unknown").GetName(); got != "generic" {
		t.Errorf("expected %q, got %q", "generic", got)
	}
}

func TestBaseDialectDirectMethods(t *testing.T) {
	base := &driver.BaseDialect{}

	if got := base.BuildLimitOffset(10, 20); got != "LIMIT 10 OFFSET 20" {
		t.Errorf("expected %q, got %q", "LIMIT 10 OFFSET 20", got)
	}
	if got := base.BuildLimitOffset(5, -1); got != "LIMIT 5" {
		t.Errorf("expected %q, got %q", "LIMIT 5", got)
	}
	if got := base.BuildLimitOffset(-1, 20); got != "OFFSET 20" {
		t.Errorf("expected %q, got %q", "OFFSET 20", got)
	}
	if got := base.BuildLimitOffset(-1, -1); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}

	if got := base.QuoteLiteral("value"); got != "'value'" {
		t.Errorf("expected %q, got %q", "'value'", got)
	}
	if got := base.QuoteLiteral(42); got != "42" {
		t.Errorf("expected %q, got %q", "42", got)
	}
	if got := base.QuoteLiteral(true); got != "true" {
		t.Errorf("expected %q, got %q", "true", got)
	}
	if got := base.QuoteLiteral([]int{1, 2, 3}); got != "'[1 2 3]'" {
		t.Errorf("expected %q, got %q", "'[1 2 3]'", got)
	}
}

func TestGetName(t *testing.T) {
	tests := []struct {
		name     string
		dialect  driver.Dialect
		expected string
	}{
		{"base", &driver.BaseDialect{}, "base"},
		{"generic", driver.NewGenericDialect(), "generic"},
		{"postgres", driver.NewPostgresDialect(), "postgres"},
		{"mssql", driver.NewMSSQLDialect(), "mssql"},
		{"mysql", driver.NewMySQLDialect(), "mysql"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dialect.GetName(); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestQuoteIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		dialect  driver.Dialect
		input    string
		expected string
	}{
		{"base", &driver.BaseDialect{}, "user", "user"},
		{"generic", driver.NewGenericDialect(), "user", "user"},
		{"postgres", driver.NewPostgresDialect(), "user", `"user"`},
		{"mssql", driver.NewMSSQLDialect(), "user", "[user]"},
		{"mysql", driver.NewMySQLDialect(), "user", "`user`"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dialect.QuoteIdentifier(tt.input); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestQuoteType(t *testing.T) {
	tests := []struct {
		name     string
		dialect  driver.Dialect
		expected styling.QuoteStyle
	}{
		{"generic", driver.NewGenericDialect(), styling.QuoteNone},
		{"postgres", driver.NewPostgresDialect(), styling.QuoteDouble},
		{"mssql", driver.NewMSSQLDialect(), styling.QuoteBracket},
		{"mysql", driver.NewMySQLDialect(), styling.QuoteBacktick},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dialect.QuoteType(); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestPlaceholder(t *testing.T) {
	tests := []struct {
		name     string
		dialect  driver.Dialect
		n        int
		expected string
	}{
		{"base1", &driver.BaseDialect{}, 1, "?"},
		{"base99", &driver.BaseDialect{}, 99, "?"},
		{"generic1", driver.NewGenericDialect(), 1, "?"},
		{"postgres1", driver.NewPostgresDialect(), 1, "$1"},
		{"postgres5", driver.NewPostgresDialect(), 5, "$5"},
		{"mssql", driver.NewMSSQLDialect(), 1, "?"},
		{"mysql", driver.NewMySQLDialect(), 1, "?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dialect.Placeholder(tt.n); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestRenderFrom(t *testing.T) {
	tests := []struct {
		name     string
		dialect  driver.Dialect
		table    string
		alias    string
		expected string
	}{
		{"base no alias", &driver.BaseDialect{}, "users", "", "users"},
		{"base alias", &driver.BaseDialect{}, "users", "u", "users"},
		{"postgres no alias", driver.NewPostgresDialect(), "users", "", `"users"`},
		{"postgres alias", driver.NewPostgresDialect(), "users", "u", `"users" u`},
		{"mssql no alias", driver.NewMSSQLDialect(), "users", "", "[users]"},
		{"mssql alias", driver.NewMSSQLDialect(), "users", "u", "[users] u"},
		{"mysql no alias", driver.NewMySQLDialect(), "users", "", "`users`"},
		{"mysql alias", driver.NewMySQLDialect(), "users", "u", "`users` u"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dialect.RenderFrom(tt.table, tt.alias); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestSupportsReturning(t *testing.T) {
	tests := []struct {
		name     string
		dialect  driver.Dialect
		expected bool
	}{
		{"base", &driver.BaseDialect{}, false},
		{"generic", driver.NewGenericDialect(), false},
		{"postgres", driver.NewPostgresDialect(), true},
		{"mssql", driver.NewMSSQLDialect(), false},
		{"mysql", driver.NewMySQLDialect(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dialect.SupportsReturning(); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestSupportsUpsert(t *testing.T) {
	tests := []struct {
		name     string
		dialect  driver.Dialect
		expected bool
	}{
		{"base", &driver.BaseDialect{}, false},
		{"generic", driver.NewGenericDialect(), false},
		{"postgres", driver.NewPostgresDialect(), true},
		{"mssql", driver.NewMSSQLDialect(), false},
		{"mysql", driver.NewMySQLDialect(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dialect.SupportsUpsert(); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		d := driver.BaseDialect{
			Name:             "test",
			QuoteStyle:       styling.QuoteNone,
			PlaceholderStyle: styling.PlaceholderQuestion,
		}
		if err := d.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if got := d.Placeholder(0); got != "?" {
			t.Errorf("expected %q, got %q", "?", got)
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		d := driver.BaseDialect{Name: ""}
		err := d.Validate()
		if err == nil {
			t.Errorf("expected error, got nil")
		} else if !strings.Contains(err.Error(), "dialect is not configured") {
			t.Errorf("expected error containing %q, got %q", "dialect is not configured", err.Error())
		}
	})

	// Uncomment and fix if Placeholder validation changes:
	// t.Run("MissingPlaceholder", func(t *testing.T) {
	// 	d := driver.BaseDialect{Name: "test", PlaceholderStyle: styling.PlaceholderDollar}
	// 	if err := d.Validate(); err == nil {
	// 		t.Errorf("expected error, got nil")
	// 	}
	// })
}
