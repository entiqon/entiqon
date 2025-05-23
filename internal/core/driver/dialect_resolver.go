package driver

import "strings"

// ResolveDialect maps a dialect name to a Dialect implementation.
//
// Valid names:
//   - "postgres": returns PostgresDialect
//   - "mysql": returns MySQLDialect (future support)
//   - any unknown: returns a generic BaseDialect
func ResolveDialect(name string) Dialect {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "postgres", "postgresql":
		return NewPostgresDialect()

	// Add more dialects here as needed:
	// case "mysql":
	//     return NewMySQLDialect()

	default:
		return &BaseDialect{name: "generic"}
	}
}
