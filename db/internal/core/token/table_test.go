// File: db/internal/core/token/table_test.go
// Since: v1.5.0

package token_test

import (
	"testing"

	"github.com/entiqon/db/internal/core/token"
)

func TestNewTable_Basic(t *testing.T) {
	tbl := token.NewTable("users")
	if !tbl.IsValid() {
		t.Errorf("expected table to be valid")
	}
	if got := tbl.String(); got != "users" {
		t.Errorf("expected %q, got %q", "users", got)
	}
}

func TestNewTableWithAlias(t *testing.T) {
	tbl := token.NewTableWithAlias("orders", "o")
	if !tbl.IsValid() {
		t.Errorf("expected table to be valid")
	}
	if got := tbl.String(); got != "orders o" {
		t.Errorf("expected %q, got %q", "orders o", got)
	}
}

func TestNewTable_Empty(t *testing.T) {
	tbl := token.NewTable("")
	if tbl.IsValid() {
		t.Errorf("expected empty table to be invalid")
	}
	if got := tbl.String(); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestNewTableWithAlias_TrimmedInput(t *testing.T) {
	tbl := token.NewTableWithAlias("  logs  ", "  l ")
	if !tbl.IsValid() {
		t.Errorf("expected table to be valid")
	}
	if got := tbl.String(); got != "logs l" {
		t.Errorf("expected %q, got %q", "logs l", got)
	}
}
