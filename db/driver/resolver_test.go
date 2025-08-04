// File: db/driver/resolver_test.go

package driver_test

import (
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/stretchr/testify/assert"
)

func TestResolveDialect_MySQL(t *testing.T) {
	mysql := driver.ResolveDialect("mysql")
	assert.Equal(t, "mysql", mysql.GetName())

	mariadb := driver.ResolveDialect("mariadb")
	assert.Equal(t, "mysql", mariadb.GetName())
}

func TestResolveDialect_MSSQL(t *testing.T) {
	mssql := driver.ResolveDialect("mssql")
	assert.Equal(t, "mssql", mssql.GetName())

	sqlserver := driver.ResolveDialect("sqlserver")
	assert.Equal(t, "mssql", sqlserver.GetName())
}

func TestResolveDialect_UnknownFallback(t *testing.T) {
	d := driver.ResolveDialect("unknown_db")
	assert.Equal(t, "generic", d.GetName())
	assert.NoError(t, d.Validate())
}
