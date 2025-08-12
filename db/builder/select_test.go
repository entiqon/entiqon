// File: db/builder/select_test.go

package builder_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/builder"
)

func TestSelectBuilder_Build(t *testing.T) {
	table := "cuno_partnership"
	sb := builder.NewSelect(nil)

	// No columns default to * and no error
	sb.Source(table)
	sql, err := sb.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `SELECT * FROM "cuno_partnership"`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	fmt.Println(sb.String())

	// No table specified error
	sb = builder.NewSelect(nil).
		Columns("LOWER(m3_cuno || '-' || partnership_id) id", "m3_cuno", "partnership_id")
	_, err = sb.Build()
	if err == nil || !strings.Contains(err.Error(), "no table specified") {
		t.Errorf("expected 'no table specified' error, got %v", err)
	}
	fmt.Println(sb.String())

	// Multiple columns
	sb = builder.NewSelect(nil).
		Columns("id", "m3_cuno", "partnership_id").
		Source("cuno_partnership")
	sql, err = sb.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = `SELECT "id", "m3_cuno", "partnership_id" FROM "cuno_partnership"`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	fmt.Println(sb.String())

	// Limit and Offset
	sb = builder.NewSelect(nil).
		Columns("id").
		Source("users").
		Limit(10).
		Offset(20)
	sql, err = sb.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = `SELECT "id" FROM "users" LIMIT 10 OFFSET 20`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	fmt.Println(sb.String())
}

func TestSelectBuilder_String(t *testing.T) {
	sb := builder.NewSelect(nil).
		Columns("id").
		Source("users")

	want := `Status ✅: SQL=SELECT "id" FROM "users", Params=`
	got := sb.String()
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	// Test error case when no table
	sb = builder.NewSelect(nil).
		Columns("id")

	wantErrPrefix := "Status ❌: Error building SQL"
	got = sb.String()
	if !strings.HasPrefix(got, wantErrPrefix) {
		t.Errorf("String() error output = %q, want prefix %q", got, wantErrPrefix)
	}
}
