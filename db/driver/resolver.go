// File: db/driver/resolver.go

package driver

import "strings"

// ResolveDialect returns a predefined dialect instance based on the input name.
//
// Recognized values:
//   - "postgres", "postgresql" → NewPostgresDialect()
//   - "mysql", "mariadb"       → NewMySQLDialect()
//   - "mssql", "sqlserver"     → NewMSSQLDialect()
//   - default fallback         → NewGenericDialect()
//
// This function never returns nil.
//
// Since: v1.4.0
func ResolveDialect(name string) Dialect {
	name = strings.ToLower(strings.TrimSpace(name))

	switch name {
	case "postgres", "postgresql":
		return NewPostgresDialect()
	case "mysql", "mariadb":
		return NewMySQLDialect()
	case "mssql", "sqlserver":
		return NewMSSQLDialect()
	default:
		return NewGenericDialect()
	}
}
