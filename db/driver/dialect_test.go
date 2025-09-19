package driver_test

import (
	"testing"

	"github.com/entiqon/db/driver"
)

func TestDialect(t *testing.T) {
	t.Run("generic", func(t *testing.T) {
		d := driver.NewGenericDialect()
		if got := d.Placeholder(1); got != "?" {
			t.Errorf("expected %q, got %q", "?", got)
		}
		if got := d.QuoteIdentifier("id"); got != "id" {
			t.Errorf("expected %q, got %q", "id", got)
		}
		if d.SupportsUpsert() {
			t.Errorf("expected SupportsUpsert=false")
		}
		if d.SupportsReturning() {
			t.Errorf("expected SupportsReturning=false")
		}
	})

	t.Run("db2", func(t *testing.T) {
		d := driver.NewDB2Dialect()
		if got := d.PlaceholderNamed("id"); got != ":id" {
			t.Errorf("expected %q, got %q", ":id", got)
		}
		if got := d.QuoteIdentifier("id"); got != "\"id\"" {
			t.Errorf("expected %q, got %q", "\"id\"", got)
		}
		if !d.SupportsReturning() {
			t.Errorf("expected SupportsReturning=true")
		}
		if !d.SupportsUpsert() {
			t.Errorf("expected SupportsUpsert=true")
		}
	})

	t.Run("firebird", func(t *testing.T) {
		d := driver.NewFirebirdDialect()
		if got := d.Placeholder(1); got != "?" {
			t.Errorf("expected %q, got %q", "?", got)
		}
		if got := d.QuoteIdentifier("id"); got != "\"id\"" {
			t.Errorf("expected %q, got %q", "\"id\"", got)
		}
		if !d.SupportsReturning() {
			t.Errorf("expected SupportsReturning=true")
		}
		if !d.SupportsUpsert() {
			t.Errorf("expected SupportsUpsert=true")
		}
	})

	t.Run("informix", func(t *testing.T) {
		d := driver.NewInformixDialect()
		if got := d.Placeholder(1); got != "?" {
			t.Errorf("expected %q, got %q", "?", got)
		}
		if got := d.QuoteIdentifier("id"); got != "\"id\"" {
			t.Errorf("expected %q, got %q", "\"id\"", got)
		}
		if !d.SupportsReturning() {
			t.Errorf("expected SupportsReturning=true")
		}
		if d.SupportsUpsert() {
			t.Errorf("expected SupportsUpsert=false")
		}
	})

	t.Run("mssql", func(t *testing.T) {
		d := driver.NewMSSQLDialect()
		if got := d.Placeholder(1); got != "?" {
			t.Errorf("expected %q, got %q", "?", got)
		}
		if got := d.QuoteIdentifier("id"); got != "[id]" {
			t.Errorf("expected %q, got %q", "[id]", got)
		}
		if d.SupportsReturning() {
			t.Errorf("expected SupportsReturning=false")
		}
		if d.SupportsUpsert() {
			t.Errorf("expected SupportsUpsert=false")
		}
	})

	t.Run("mysql", func(t *testing.T) {
		d := driver.NewMySQLDialect()
		if got := d.Placeholder(1); got != "?" {
			t.Errorf("expected %q, got %q", "?", got)
		}
		if got := d.QuoteIdentifier("id"); got != "`id`" {
			t.Errorf("expected %q, got %q", "`id`", got)
		}
		if d.SupportsReturning() {
			t.Errorf("expected SupportsReturning=false")
		}
		if d.SupportsUpsert() {
			t.Errorf("expected SupportsUpsert=false")
		}
	})

	t.Run("oracle", func(t *testing.T) {
		d := driver.NewOracleDialect()
		if got := d.PlaceholderNamed("id"); got != ":id" {
			t.Errorf("expected %q, got %q", ":id", got)
		}
		if got := d.QuoteIdentifier("id"); got != "\"id\"" {
			t.Errorf("expected %q, got %q", "\"id\"", got)
		}
		if !d.SupportsUpsert() {
			t.Errorf("expected SupportsUpsert=true")
		}
		if !d.SupportsReturning() {
			t.Errorf("expected SupportsReturning=true")
		}
	})

	t.Run("postgres", func(t *testing.T) {
		d := driver.NewPostgresDialect()
		if got := d.Placeholder(1); got != "$1" {
			t.Errorf("expected %q, got %q", "$1", got)
		}
		if got := d.QuoteIdentifier("id"); got != "\"id\"" {
			t.Errorf("expected %q, got %q", "\"id\"", got)
		}
		if !d.SupportsUpsert() {
			t.Errorf("expected SupportsUpsert=true")
		}
		if !d.SupportsReturning() {
			t.Errorf("expected SupportsReturning=true")
		}
	})

	t.Run("sqlite", func(t *testing.T) {
		d := driver.NewSQLiteDialect()
		if got := d.Placeholder(1); got != "?" {
			t.Errorf("expected %q, got %q", "?", got)
		}
		if got := d.QuoteIdentifier("id"); got != "\"id\"" {
			t.Errorf("expected %q, got %q", "\"id\"", got)
		}
		if !d.SupportsUpsert() {
			t.Errorf("expected SupportsUpsert=true")
		}
		if !d.SupportsReturning() {
			t.Errorf("expected SupportsReturning=true")
		}
	})
}
