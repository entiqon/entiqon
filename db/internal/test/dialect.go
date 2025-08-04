// File: db/internal/core/test/dialect.go

package test

// TestDialect is a mock dialect that escapes identifiers using double quotes,
// matching ANSI SQL and PostgreSQL expectations.
type TestDialect struct{}

// Escape wraps identifiers in double quotes, e.g., "users"
func (TestDialect) Escape(identifier string) string {
	return `"` + identifier + `"`
}

// Name returns the name of the test dialect.
func (TestDialect) Name() string {
	return "test"
}
