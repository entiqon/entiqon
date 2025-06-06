package token

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/driver"
)

// BaseToken provides a normalized representation of a raw SQL-like expression,
// including name and optional alias parsing. It acts as a foundational structure
// shared across tokens like Column and Table.
//
// It is designed to be embedded in higher-level tokens such as Column,
// Table, Join, or Condition, offering unified handling of identifier
// semantics, alias resolution, and error reporting.
//
// BaseToken performs *general* validation (e.g., empty input, malformed aliases),
// but does not resolve context-specific semantics such as:
//
//   - Table qualification (e.g., "users.id")
//   - Ownership or binding to other sources
//   - Dialect rendering preferences
//
// Qualification logic — like parsing "table.column" — is the responsibility of
// higher-level tokens such as Column, since only some SQL expressions are validly qualified.
//
// This keeps BaseToken adaptable, non-opinionated, and reusable across different token types.
//
// This struct should not be used standalone to represent SQL elements,
// but is intended as an internal abstraction to simplify composition.
//
// # Usage
//
//	type Column struct {
//	    BaseToken
//	    Qualified string
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

	// kind classifies the token according to its role in a SQL query,
	// such as Column, Table, or Condition.
	//
	// This field is assigned internally during token construction and is not exported,
	// preventing unintended modifications. It is used by the Kind() method to support
	// type-safe introspection and rendering logic.
	//
	// Valid values include:
	//   - ColumnKind
	//   - TableKind
	//   - ConditionKind
	//   - UnknownKind (default)
	kind Kind
}

// NewBaseToken constructs a new BaseToken by parsing the input string.
// It supports optional aliasing, both inline (e.g., "field AS alias") and explicit via variadic arguments.
//
// If the input expression is malformed, empty, or structurally invalid (e.g., starts with "AS" or uses commas),
// the returned BaseToken will carry an appropriate error.
//
// If both inline and explicit aliases are provided and conflict, the explicit one takes precedence,
// and an alias conflict error is assigned.
//
// This function is intended to centralize alias parsing and validation logic used by higher-level
// tokens such as Column or Table.
//
// # Examples
//
//	b := NewBaseToken("users.id")
//	fmt.Println(b.Name)   // "users.id"
//	fmt.Println(b.Alias)  // ""
//	fmt.Println(b.Error)  // <nil>
//
//	b = NewBaseToken("users.id AS uid")
//	fmt.Println(b.Name)   // "users.id"
//	fmt.Println(b.Alias)  // "uid"
//	fmt.Println(b.Error)  // <nil>
//
//	b = NewBaseToken("users.id", "uid")
//	fmt.Println(b.Name)   // "users.id"
//	fmt.Println(b.Alias)  // "uid"
//	fmt.Println(b.Error)  // <nil>
//
//	b = NewBaseToken("users.id AS user_id", "uid")
//	fmt.Println(b.Name)   // "users.id"
//	fmt.Println(b.Alias)  // "uid"
//	fmt.Println(b.Error)  // alias conflict: explicit alias "uid" does not match inline alias "user_id"
//
//	b = NewBaseToken("AS uid")
//	fmt.Println(b.Error)  // invalid input expression: cannot start with 'AS'
//
//	b = NewBaseToken("")
//	fmt.Println(b.Error)  // invalid input expression: expression is empty
//
//	b = NewBaseToken("id, name")
//	fmt.Println(b.Error)  // invalid input expression: aliases must not be comma-separated
func NewBaseToken(input string, alias ...string) *BaseToken {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return &BaseToken{
			Source: input,
			Error:  fmt.Errorf("invalid input expression: expression is empty"),
		}
	}

	if strings.Contains(input, ",") {
		return &BaseToken{
			Source: input,
			Error:  fmt.Errorf("invalid input expression: aliases must not be comma-separated"),
		}
	}

	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "AS ") {
		return &BaseToken{
			Source: input,
			Error:  fmt.Errorf("invalid input expression: cannot start with 'AS'"),
		}
	}

	base, parsedAlias := ParseAlias(input)
	if strings.TrimSpace(base) == "AS" {
		return &BaseToken{
			Source: input,
			Error:  fmt.Errorf("invalid input expression: name cannot be AS keyword"),
		}
	}

	b := &BaseToken{
		Source: input,
		Name:   base,
		Alias:  parsedAlias,
	}

	if len(alias) > 0 && alias[0] != "" {
		if parsedAlias != "" && alias[0] != parsedAlias {
			b.Error = fmt.Errorf(
				"alias conflict: explicit alias %q does not match inline alias %q",
				alias[0], parsedAlias,
			)
		}
		b.Alias = alias[0]
	}

	return b
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

// GetName safely returns the Name of the BaseToken.
//
// This method is a defensive accessor used to avoid nil pointer dereference
// when accessing the Name field of a potential nil *BaseToken.
//
// It is commonly used in higher-level tokens (e.g., Column, Table) to extract
// the logical identifier associated with the token, while maintaining stability
// when BaseToken may not have been initialized.
//
// Returns:
//   - The Name string if BaseToken is non-nil
//   - An empty string ("") if BaseToken is nil
//
// Example:
//
//	var b *BaseToken = nil
//	name := b.GetName() // safely returns ""
//
//	b = NewBaseToken("id")
//	name = b.GetName() // returns "id"
//
// Usage in Column:
//
//	if col.BaseToken.GetName() == "id" {
//	    // Perform logic using column name
//	}
func (b *BaseToken) GetName() string {
	if b == nil {
		return ""
	}
	return b.Name
}

// GetSource returns the original input string associated with the token.
//
// This method safely retrieves the `input` field of the BaseToken, which was
// previously known as `Source`. The `input` typically represents the raw
// expression string used to construct the token (e.g., "users.id AS user_id").
//
// If the BaseToken is nil, this method returns an empty string.
//
// This accessor ensures safe and consistent use of the underlying expression,
// and helps decouple the internal representation (`input`) from external usage.
//
// Example:
//
//	b := token.NewBaseToken("id", "users.id AS user_id", "user_id")
//	raw := b.GetSource() // returns "users.id AS user_id"
//
//	var b2 *token.BaseToken
//	raw = b2.GetSource() // safely returns ""
func (b *BaseToken) GetSource() string {
	if b == nil {
		return ""
	}
	return b.Source
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

// GetKind returns the classification kind associated with this token.
//
// If the token is nil, it safely returns UnknownKind.
// If no kind has been explicitly set, UnknownKind is also returned by default.
//
// Concrete token types (e.g., Column, Table) should assign a specific kind
// via SetKind(...) during instantiation.
//
// # Example
//
//	b := NewBaseToken("id")
//	b.SetKind(ColumnKind)
//	fmt.Println(b.GetKind()) // → ColumnKind
func (b *BaseToken) GetKind() Kind {
	if b == nil {
		return UnknownKind
	}
	return b.kind
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
// If the alias is empty or the qualified name is empty, aliasing is skipped
// and the qualified name is returned as-is. This avoids emitting invalid SQL
// like `AS ""` or aliasing blank expressions.
//
// If the dialect is nil, a basic unquoted alias is used.
// Otherwise, the alias will be quoted using the dialect’s identifier rules.
//
// # Example
//
//	input := "u.id AS user_id"
//	b := NewBaseToken(input)
//	b.Alias = "user_id"
//
//	fmt.Println(b.RenderAlias(postgres, "u.id")) // → u.id AS "user_id"
//
//	b = NewBaseToken(input)
//	fmt.Println(b.RenderAlias(postgres, "u.id")) // → u.id
//
//	b = NewBaseToken("")
//	fmt.Println(b.RenderAlias(postgres, ""))     // → ""
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
//
// If the token is nil or the Name field is empty, it returns an empty string.
// If the dialect is nil, the name is returned as-is without quoting.
// Otherwise, it applies the dialect's identifier quoting rules.
//
// This is commonly used for rendering column or table names in
// SELECT, FROM, or JOIN clauses.
//
// # Example
//
//	b := NewBaseToken("id")
//	b.Name = "id"
//
//	fmt.Println(b.RenderName(driver.NewPostgresDialect())) // → "id"
//
//	b = NewBaseToken("")
//	fmt.Println(b.RenderName(nil))                         // → ""
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
// along with the original input string that triggered the failure.
//
// This method sets the Error and internal input fields (if unset),
// preserving the original expression for diagnostic purposes.
//
// It returns the same *BaseToken instance, allowing fluent chaining
// from higher-level token wrappers such as Column or Table.
//
// This method does not attempt to parse or correct the input — it is
// purely used for marking invalid or unsupported expressions during
// token construction or validation.
//
// # Example
//
//	b := NewBaseToken("AS uid") // Invalid: missing name before 'AS'
//	b.SetErrorWith("AS uid", fmt.Errorf("name is missing before 'AS'"))
//
//	fmt.Println(b.Error)        // name is missing before 'AS'
//	fmt.Println(b.GetSource())  // AS uid
//
//	// Output:
//	name is missing before 'AS'
//	AS uid
func (b *BaseToken) SetErrorWith(source string, err error) *BaseToken {
	b.Error = err
	if b.Source == "" {
		b.Source = source
	}
	return b
}

// SetKind assigns the internal Kind classification to the token.
// This should be used by higher-level token constructors (e.g., Column, Table)
// to explicitly declare the token type.
//
// This method is safe to call during instantiation, but not intended
// for dynamic reclassification at runtime.
//
// # Example
//
//	b := NewBaseToken("id")
//	b.SetKind(ColumnKind)
//	fmt.Println(b.GetKind()) // → ColumnKind
func (b *BaseToken) SetKind(k Kind) {
	if b != nil {
		b.kind = k
	}
}

// String returns a diagnostic string representation of the token,
// including its internal kind, alias status, and any error state.
//
// This method is intended for logging and testing — it does not produce SQL output.
//
// # Example
//
//	b := NewBaseToken("id")
//	b.SetKind(ColumnKind)
//	fmt.Println(b.String()) // → Column("id") [aliased: false, errored: false]
//
//	b.SetAlias("uid")
//	fmt.Println(b.String()) // → Column("id") [aliased: true, errored: false]
//
//	b.SetErrorWith("id AS uid", fmt.Errorf("alias conflict"))
//	fmt.Println(b.String()) // → Column("id") [aliased: true, errored: true, error: alias conflict]
func (b *BaseToken) String() string {
	if b == nil {
		return "nil"
	}

	suffix := fmt.Sprintf("aliased: %v, errored: %v", b.IsAliased(), b.HasError())
	if b.HasError() {
		suffix += fmt.Sprintf(", error: %s", b.Error.Error())
	}

	return fmt.Sprintf("%s(\"%s\") [%s]", b.kind.String(), b.Name, suffix)
}
