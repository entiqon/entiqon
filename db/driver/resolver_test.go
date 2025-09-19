// File: db/driver/resolver_test.go

package driver_test

import (
	"testing"

	"github.com/entiqon/db/driver"
)

func TestResolveDialect_MySQL(t *testing.T) {
	mysql := driver.ResolveDialect("mysql")
	if got := mysql.GetName(); got != "mysql" {
		t.Errorf("expected %q, got %q", "mysql", got)
	}

	mariadb := driver.ResolveDialect("mariadb")
	if got := mariadb.GetName(); got != "mysql" {
		t.Errorf("expected %q, got %q", "mysql", got)
	}
}

func TestResolveDialect_MSSQL(t *testing.T) {
	mssql := driver.ResolveDialect("mssql")
	if got := mssql.GetName(); got != "mssql" {
		t.Errorf("expected %q, got %q", "mssql", got)
	}

	sqlserver := driver.ResolveDialect("sqlserver")
	if got := sqlserver.GetName(); got != "mssql" {
		t.Errorf("expected %q, got %q", "mssql", got)
	}
}

func TestResolveDialect_UnknownFallback(t *testing.T) {
	d := driver.ResolveDialect("unknown_db")
	if got := d.GetName(); got != "generic" {
		t.Errorf("expected %q, got %q", "generic", got)
	}

	if err := d.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
