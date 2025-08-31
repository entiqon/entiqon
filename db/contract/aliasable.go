package contract

// Aliasable is implemented by tokens that support aliasing.
//
// Aliasing allows a token's expression to be referenced by an alternative
// shorter name. For example:
//
//	SELECT u.id AS user_id
//
// In this case, "u.id" is the underlying expression and "user_id" is the alias.
//
// Aliasable provides two methods:
//
//   - Alias() — returns the alias string (empty if none).
//   - IsAliased() — reports whether an alias is set.
//
// Typical implementers include Field and Table tokens. Tokens that do not
// support aliasing should not implement this contract.
//
// See ExampleAliasable in example_test.go for a runnable demonstration.
type Aliasable interface {
	// Alias returns the alias name, if present.
	// Returns an empty string if the token has no alias.
	Alias() string

	// IsAliased reports whether the token has an alias.
	// Returns true if Alias() is non-empty.
	IsAliased() bool
}
