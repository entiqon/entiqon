// File: db/dialect/postgres_test.go

package dialect_test

import (
	"testing"

	"github.com/entiqon/db/dialect"
)

func TestPostgresDialect(t *testing.T) {
	pg := &dialect.PostgresDialect{}

	// Test Name
	if got := pg.Name(); got != "postgres" {
		t.Errorf("Name() = %q; want 'postgres'", got)
	}

	// Test QuoteIdentifier with embedded quotes
	input := `user"name`
	want := `"user""name"`
	if got := pg.QuoteIdentifier(input); got != want {
		t.Errorf("QuoteIdentifier(%q) = %q; want %q", input, got, want)
	}

	// Test Placeholder
	for i, want := range []string{"$1", "$2", "$3"} {
		if got := pg.Placeholder(i + 1); got != want {
			t.Errorf("Placeholder(%d) = %q; want %q", i+1, got, want)
		}
	}

	// Test SupportsReturning
	if !pg.SupportsReturning() {
		t.Error("SupportsReturning() = false; want true")
	}

	// Test PaginationSyntax
	tests := []struct {
		limit, offset int
		want          string
	}{
		{10, 20, "LIMIT 10 OFFSET 20"},
		{10, 0, "LIMIT 10"},
		{0, 20, "OFFSET 20"},
		{0, 0, ""},
	}

	for _, tt := range tests {
		if got := pg.PaginationSyntax(tt.limit, tt.offset); got != tt.want {
			t.Errorf("PaginationSyntax(%d, %d) = %q; want %q", tt.limit, tt.offset, got, tt.want)
		}
	}
}
