// File: db/driver/dialect_test.go

package driver_test

import (
	"testing"

	"github.com/entiqon/entiqon/db/driver"
	"github.com/stretchr/testify/assert"
)

func TestDialect(t *testing.T) {
	t.Run("generic", func(t *testing.T) {
		d := driver.NewGenericDialect()
		assert.Equal(t, "?", d.Placeholder(1))
		assert.Equal(t, "id", d.QuoteIdentifier("id"))
		assert.False(t, d.SupportsUpsert())
		assert.False(t, d.SupportsReturning())
	})

	t.Run("db2", func(t *testing.T) {
		d := driver.NewDB2Dialect()
		assert.Equal(t, ":id", d.PlaceholderNamed("id"))
		assert.Equal(t, "\"id\"", d.QuoteIdentifier("id"))
		assert.True(t, d.SupportsReturning())
		assert.True(t, d.SupportsUpsert())
	})

	t.Run("firebird", func(t *testing.T) {
		d := driver.NewFirebirdDialect()
		assert.Equal(t, "?", d.Placeholder(1))
		assert.Equal(t, "\"id\"", d.QuoteIdentifier("id"))
		assert.True(t, d.SupportsReturning())
		assert.True(t, d.SupportsUpsert())
	})

	t.Run("informix", func(t *testing.T) {
		d := driver.NewInformixDialect()
		assert.Equal(t, "?", d.Placeholder(1))
		assert.Equal(t, "\"id\"", d.QuoteIdentifier("id"))
		assert.True(t, d.SupportsReturning())
		assert.False(t, d.SupportsUpsert())
	})

	t.Run("mssql", func(t *testing.T) {
		d := driver.NewMSSQLDialect()
		assert.Equal(t, "?", d.Placeholder(1))
		assert.Equal(t, "[id]", d.QuoteIdentifier("id"))
		assert.False(t, d.SupportsReturning())
		assert.False(t, d.SupportsUpsert())
	})

	t.Run("mysql", func(t *testing.T) {
		d := driver.NewMySQLDialect()
		assert.Equal(t, "?", d.Placeholder(1))
		assert.Equal(t, "`id`", d.QuoteIdentifier("id"))
		assert.False(t, d.SupportsReturning())
		assert.False(t, d.SupportsUpsert())
	})

	t.Run("oracle", func(t *testing.T) {
		d := driver.NewOracleDialect()
		assert.Equal(t, ":id", d.PlaceholderNamed("id"))
		assert.Equal(t, `"id"`, d.QuoteIdentifier("id"))
		assert.True(t, d.SupportsUpsert())
		assert.True(t, d.SupportsReturning())
	})

	t.Run("postgres", func(t *testing.T) {
		d := driver.NewPostgresDialect()
		assert.Equal(t, "$1", d.Placeholder(1))
		assert.Equal(t, "\"id\"", d.QuoteIdentifier("id"))
		assert.True(t, d.SupportsUpsert())
		assert.True(t, d.SupportsReturning())
	})

	t.Run("sqlite", func(t *testing.T) {
		d := driver.NewSQLiteDialect()
		assert.Equal(t, "?", d.Placeholder(1))
		assert.Equal(t, `"id"`, d.QuoteIdentifier("id"))
		assert.True(t, d.SupportsUpsert())
		assert.True(t, d.SupportsReturning())
	})
}
