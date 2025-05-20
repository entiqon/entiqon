// filename: internal/core/driver/dialect_generic.go

package driver

// NewGenericDialect returns a base dialect instance with the name "generic".
// This is used as a default fallback when no dialect is explicitly set.
//
// It provides neutral behavior using "?" as placeholder and basic identifier quoting.
//
// This is NOT an engine-specific dialect, but a non-op placeholder.
//
// Since: v1.4.0
// Updated: v1.4.0
func NewGenericDialect() Dialect {
	return &BaseDialect{DialectName: "generic"}
}
