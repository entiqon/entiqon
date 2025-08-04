// File: db/internal/core/builder/renderable.go

// Package contract defines core interfaces and types for Entiqon tokens,
// including rendering and error-handling abstractions. These interfaces live
// in the “contracts” layer to decouple token implementations from dialects
// and other build-time details.
//
// To regenerate documentation after coding sessions, use `godoc` or configure
// your Go documentation generator (e.g., `go doc`, `pkgsite`, or `godoc -http=:6060`).
package contract

// Quoter is the minimal interface required to quote SQL identifiers according
// to a specific SQL dialect. Any dialect implementation (Postgres, MySQL, etc.)
// must implement this method to be used as a Quoter.
//
// Example implementation in a dialect (Postgres):
//
//	func (p *PostgresDialect) QuoteIdentifier(id string) string {
//	    // Replace any internal double-quotes with two double-quotes:
//	    escaped := strings.ReplaceAll(id, `"`, `""`)
//	    return `"` + escaped + `"`
//	}
type Quoter interface {
	// QuoteIdentifier wraps the given identifier (e.g., table name, column name)
	// using the appropriate quoting mechanism for the dialect. Implementations
	// should escape internal occurrences of the quote character if required.
	//
	// Example:
	//   QuoteIdentifier("user")      → "\"user\""
	//   QuoteIdentifier("user\"name") → "\"user\"\"name\""
	QuoteIdentifier(id string) string
}

// Renderable defines the methods that any token-like entity should expose in order
// to be rendered into SQL or used for diagnostics. By centralizing these four methods,
// we ensure a uniform rendering strategy for tokens such as Column, Table, and Condition.
//
// Raw() and String() do not require any quoting, while RenderName and RenderAlias
// accept a Quoter to apply dialect-specific quoting. If the Quoter argument is nil,
// implementations must return unquoted results.
type Renderable interface {
	// RenderName returns the token’s identifier (alias if present, or name otherwise),
	// applying quoting only if the provided Quoter is non-nil. If Quoter is nil,
	// it returns the raw identifier (alias or name) without quotes. If the receiver
	// is nil or its Name is empty, RenderName must return an empty string.
	//
	// For example, given a token with Name="id", Alias="user_id":
	//   RenderName(q) with a Postgres Quoter → "\"user_id\""
	//   RenderName(nil)                       → "user_id"
	//
	// Examples of use:
	//   col := NewColumnToken("users.id", "u")
	//   quoted := col.RenderName(postgresDialect) // => "\"u\""
	//   plain  := col.RenderName(nil)            // => "u"
	RenderName(q Quoter) string

	// RenderAlias takes a fully-qualified identifier (qualified), and if the token has
	// a non-empty Alias, returns a string of the form "qualified AS alias", quoting
	// the alias if a Quoter is provided. If the token’s Alias is empty, or if the
	// receiver is nil, RenderAlias must return the qualified string unchanged.
	//
	// For example, given a token with Alias="u" and qualified="\"users\".\"id\"":
	//   RenderAlias(q) with a Postgres Quoter → "\"users\".\"id\" AS \"u\""
	//   RenderAlias(nil)                        → "users.id AS u"
	//   RenderAlias(any, "")                    → ""
	RenderAlias(q Quoter, qualified string) string

	// String returns a diagnostic or “pretty-print” representation of the token,
	// including its Kind (if available), Name, and metadata such as whether it is
	// aliased or errored. If the receiver is nil, String must return an empty string.
	//
	// Example output for a Column token (Kind=ColumnKind, Name="id", Alias="u", no error):
	//   Column("id") [aliased: true, errored: false]
	//
	// If the token’s kind is UnknownKind, String should print “Unknown” as the kind.
	//
	// Since: v1.7.0
	String() string
}
