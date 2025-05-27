package token

import (
	"fmt"
	"strings"
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

// NewErroredToken creates a BaseToken containing the provided error.
//
// This is used when a token (e.g., Column, Table) cannot be parsed or resolved
// and must be retained in the token stream for error reporting and validation.
//
// # Example:
//
//	col := Column{BaseToken: NewErroredToken(fmt.Errorf("empty input"))}
func NewErroredToken(err error) BaseToken {
	return BaseToken{Error: err}
}

// HasError reports whether the token has encountered a semantic or structural error.
//
// Typical causes include alias mismatches, unresolved references, or conflicting
// overrides detected during token construction or resolution.
func (b BaseToken) HasError() bool {
	return b.Error != nil
}

// IsValid returns true if the token has a usable identifier and no associated error.
//
// This method ensures that the token is structurally well-formed and ready
// to be included in a generated SQL query.
func (b BaseToken) IsValid() bool {
	return b.Error == nil && strings.TrimSpace(b.Name) != ""
}

// IsAliased reports whether the token has a defined alias.
//
// Returns true if the Alias field is non-empty, which indicates the token
// should appear in SQL output with an AS clause or similar aliasing logic.
func (b BaseToken) IsAliased() bool {
	return b.Alias != ""
}

// AliasOr returns the alias if defined, or falls back to the Name.
//
// This is helpful in rendering column headers or result labels where aliases
// take precedence, but a fallback is still required.
func (b BaseToken) AliasOr() string {
	if b.Alias != "" {
		return b.Alias
	}
	return b.Name
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
func (b BaseToken) Raw() string {
	if b.Alias != "" {
		return fmt.Sprintf("%s AS %s", b.Name, b.Alias)
	}
	return b.Name
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
func (b BaseToken) String(kind Kind) string {
	suffix := fmt.Sprintf("aliased: %v, errored: %v", b.IsAliased(), b.HasError())
	if b.HasError() {
		suffix += fmt.Sprintf(", error: %s", b.Error.Error())
	}

	return fmt.Sprintf("%s(\"%s\") [%s]", kind, b.Name, suffix)
}
