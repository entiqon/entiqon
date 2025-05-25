package token

import (
	"fmt"
	"strings"
)

// AliasableToken provides shared fields and behavior for SQL tokens
// that support aliasing, such as columns and tables.
//
// It is intended to be embedded in higher-level token structs to
// standardize alias resolution and error handling.
//
// Since: v1.5.0
type AliasableToken struct {
	// Name is the unquoted identifier used in SQL (e.g., column or table name).
	Name string

	// Alias is the optional alias for the identifier.
	Alias string

	// Error holds a parsing or validation error associated with the token.
	Error error
}

// IsValid returns true if the token has a valid, non-empty Name and no parsing Error.
//
// It does not check for SQL safety or reserved words.
//
// Example:
//
//	a := AliasableToken{Name: "id"}
//	valid := a.IsValid() // true
func (a AliasableToken) IsValid() bool {
	return a.Error == nil && strings.TrimSpace(a.Name) != ""
}

// Raw returns the SQL representation of the token.
//
// If an alias is present, returns: "name AS alias".
// Otherwise, returns only the Name.
//
// Example:
//
//	a := AliasableToken{Name: "users", Alias: "u"}
//	raw := a.Raw() // "users AS u"
func (a AliasableToken) Raw() string {
	if a.Alias != "" {
		return fmt.Sprintf("%s AS %s", a.Name, a.Alias)
	}
	return a.Name
}

// String returns a debug-friendly string representation of the token,
// with an explicit type prefix for clarity.
//
// The typ argument should indicate the struct embedding the token
// (e.g., "Column", "Table").
//
// Example:
//
//	a := AliasableToken{Name: "id", Alias: "user_id"}
//	log := a.String("Column") // Column("id" AS "user_id")
func (a AliasableToken) String(typ string) string {
	if a.Alias != "" {
		return fmt.Sprintf(`%s("%s" AS "%s")`, typ, a.Name, a.Alias)
	}
	return fmt.Sprintf(`%s("%s")`, typ, a.Name)
}
