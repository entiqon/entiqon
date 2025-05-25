// File: internal/driver/alias_style.go

package styling

// AliasStyle defines how SQL aliases are rendered in dialect-specific implementations.
type AliasStyle int

const (
	// AliasNone disables aliasing entirely.
	AliasNone AliasStyle = iota

	// AliasWithoutKeyword renders alias without the AS keyword (e.g., "users u").
	AliasWithoutKeyword

	// AliasWithKeyword renders alias using the AS keyword (e.g., "users AS u").
	AliasWithKeyword
)
