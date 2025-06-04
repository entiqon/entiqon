package token

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/driver"
)

// BaseToken provides a reusable building block for SQL token types that
// carry a core name, optional alias, and possible validation error.
//
// It is designed to be embedded in higher-level tokens such as Column,
// Table, Join, or Condition, offering unified handling of identifier
// semantics, alias resolution, and error reporting.
//
// This struct should not be used standalone to represent SQL elements,
// but is intended as an internal abstraction to simplify composition.
//
// # Usage
//
//	type Column struct {
//	    BaseToken
//	    TableName string
//	}
//
//	c := Column{
//	    BaseToken: BaseToken{Name: "id", Alias: "user_id"},
//	}
//
//	fmt.Println(c.Raw())      // "id"
//	fmt.Println(c.IsAliased()) // true
//	fmt.Println(c.String())   // BaseToken("id") [aliased: true]
//
// File: internal/build/token/base.go
// Since: v1.6.0
type BaseToken struct {
	// Source holds the original input string used to construct the column.
	// Unlike Raw(), this is not formatted or rendered — it's used for diagnostics only.
	Source string

	// Name represents the core identifier of the token (e.g., column or table name).
	// It should be a raw, unquoted SQL-safe identifier.
	Name string

	// Alias is an optional alternative label for the token, used in SELECT or AS clauses.
	// If empty, the token will appear under its Name.
	Alias string

	// Error holds a semantic or structural conflict encountered during parsing,
	// such as an alias mismatch or invalid override. A nil value indicates no error.
	Error error
}

// AliasOr returns the alias if defined, or falls back to the Name.
//
// This is helpful in rendering column headers or result labels where aliases
// take precedence, but a fallback is still required.
func (b *BaseToken) AliasOr() string {
	if b.Alias != "" {
		return b.Alias
	}
	return b.Name
}

// HasError reports whether the token has encountered a semantic or structural error.
//
// Typical causes include alias mismatches, unresolved references, or conflicting
// overrides detected during token construction or resolution.
func (b *BaseToken) HasError() bool {
	return b.Error != nil
}

// IsAliased reports whether the token has a defined alias.
//
// Returns true if the Alias field is non-empty, which indicates the token
// should appear in SQL output with an AS clause or similar aliasing logic.
func (b *BaseToken) IsAliased() bool {
	return b.Alias != ""
}

// IsValid returns true if the token has a usable identifier and no associated error.
//
// This method ensures that the token is structurally well-formed and ready
// to be included in a generated SQL query.
func (b *BaseToken) IsValid() bool {
	return b != nil && b.Error == nil && strings.TrimSpace(b.Name) != ""
}

// Raw returns the base SQL expression, optionally including aliasing.
//
// If the token has an alias, the returned string will follow the format:
//
//	"name AS alias"
//
// If no alias is set, only the base name is returned.
//
// # Examples
//
//	BaseToken{Name: "id"}.Raw()                 → "id"
//	BaseToken{Name: "id", Alias: "user_id"}.Raw() → "id AS user_id"
func (b *BaseToken) Raw() string {
	if b.Alias != "" {
		return fmt.Sprintf("%s AS %s", b.Name, b.Alias)
	}
	return b.Name
}

// RenderAlias returns a dialect-quoted alias expression if an alias is set,
// otherwise returns the qualified name unchanged.
//
// If dialect is nil, it returns an unquoted fallback format.
//
// If the qualified name is empty, aliasing is suppressed entirely and an
// empty string is returned, as aliasing an empty expression produces invalid SQL.
//
// # Example
//
//	qualified := "u.id"
//	b := &BaseToken{Alias: "user_id"}
//	fmt.Println(b.RenderAlias(postgres, qualified)) // → u.id AS "user_id"
func (b *BaseToken) RenderAlias(d driver.Dialect, qualified string) string {
	if b == nil || qualified == "" {
		return qualified
	}

	if b.Alias == "" {
		return qualified
	}

	if d == nil {
		return fmt.Sprintf("%s AS %s", qualified, b.Alias)
	}

	return fmt.Sprintf("%s AS %s", qualified, d.QuoteIdentifier(b.Alias))
}

// RenderName returns the dialect-quoted name of the token.
// If the token is nil or the name is empty, it returns an empty string.
// If the dialect is nil, the name is returned unquoted.
//
// # Example
//
//	b := &BaseToken{Name: "id"}
//	fmt.Println(b.RenderName(driver.NewPostgresDialect())) → `"id"`
func (b *BaseToken) RenderName(d driver.Dialect) string {
	if b == nil || b.Name == "" {
		return ""
	}
	if d == nil {
		return b.Name
	}
	return d.QuoteIdentifier(b.Name)
}

// SetErrorWith assigns a semantic or structural error to the token,
// along with the original expression that triggered the failure.
//
// This method sets both the Error and Source fields if they are not already set,
// and returns the updated *BaseToken. It is typically called during parsing
// or validation when an invalid or unsupported expression is encountered.
//
// This method does not return a new token, but mutates the existing one
// in place. It is intended to be called by higher-level tokens such as
// Column or Table, which may wrap the result and return themselves for
// fluent chaining.
//
// # Example
//
//	b := &BaseToken{Name: "id", Alias: "uid"}
//	b.SetErrorWith("id AS uid", fmt.Errorf("alias conflict"))
//	fmt.Println(b.Error)  // alias conflict
//	fmt.Println(b.Source) // id AS uid
//
//	// Output:
//	alias conflict
//	id AS uid
func (b *BaseToken) SetErrorWith(source string, err error) *BaseToken {
	b.Error = err
	if b.Source == "" {
		b.Source = source
	}
	return b
}

// String returns a diagnostic string representation of the token,
// including the token type label, alias status, and any error.
//
// This method is intended for logging and test assertions only — it does
// not produce SQL output for execution.
//
// # Examples
//
//	b := BaseToken{Name: "id"}
//	fmt.Println(b.String(KindColumn)) → Column("id") [aliased: false]
//
//	b = BaseToken{Name: "id", Alias: "user_id"}
//	fmt.Println(b.String(KindColumn)) → Column("id") [aliased: true]
//
//	b = BaseToken{Name: "id", Alias: "uid", Error: fmt.Errorf("conflict")}
//	fmt.Println(b.String(KindColumn)) → Column("id") [aliased: true, error: conflict]
func (b *BaseToken) String(kind Kind) string {
	suffix := fmt.Sprintf("aliased: %v, errored: %v", b.IsAliased(), b.HasError())
	if b.HasError() {
		suffix += fmt.Sprintf(", error: %s", b.Error.Error())
	}

	return fmt.Sprintf("%s(\"%s\") [%s]", kind, b.Name, suffix)
}
