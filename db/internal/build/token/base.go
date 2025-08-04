// File: db/internal/build/token/ase.go

package token

import (
	"fmt"
	"strings"

	contract2 "github.com/entiqon/db/internal/core/contract"
)

// BaseToken provides a normalized representation of a raw SQL-like expression,
// including name and optional alias parsing. It acts as a foundational structure
// shared across tokens like Column, Table, Join, or Condition.
//
// It is designed to be embedded in higher-level tokens, offering unified handling
// of identifier semantics, alias resolution, error reporting, and type classification.
//
// BaseToken performs *general* validation (e.g., empty input, malformed aliases),
// but does not resolve context-specific semantics such as:
//   - Table qualification (e.g., "users.id")
//   - Ownership or binding to other sources
//   - Dialect rendering preferences
//
// Qualification logic—like parsing "table.column"—is the responsibility of higher-level tokens.
// This keeps BaseToken adaptable, non-opinionated, and reusable across different token types.
//
// This struct should not be used standalone for generating SQL. It serves as an internal
// abstraction to simplify composition.
//
// # Example
//
//	type Column struct {
//	    *BaseToken
//	    Qualified string
//	}
//
//	b := NewBaseToken("users.id AS uid")
//	b.SetKind(ColumnKind)
//	fmt.Println(b.String())
//	// Output: Column("id") [aliased: true, errored: false]
type BaseToken struct {
	// input holds the original raw input string used to construct this token.
	// Unlike Raw(), this is not formatted or rendered—it is used for diagnostics only.
	input string

	// name represents the core identifier of the token (e.g., column or table name).
	// It should be a raw, unquoted SQL-safe identifier.
	name string

	// Alias is an optional alternative label for the token, used in SELECT or AS clauses.
	// If empty, the token will appear under its name.
	alias string

	// Error holds a semantic or structural conflict encountered during parsing,
	// such as an alias mismatch or invalid override. A nil value indicates no error.
	err error

	// kind classifies the token according to its role in a SQL query,
	// such as ColumnKind, TableKind, or ConditionKind.
	//
	// This field is assigned internally during token construction and is not exported,
	// preventing unintended modifications. It is used by GetKind() and String()
	// to support type-safe introspection and diagnostic output.
	//
	// Valid values include:
	//   - ColumnKind
	//   - TableKind
	//   - ConditionKind
	//   - UnknownKind (default)
	kind contract2.Kind
}

// NewBaseToken constructs a new BaseToken by parsing the input string and optional explicit alias.
// It performs the following steps:
//  1. Trim whitespace and ensure the input is non-empty.
//  2. Reject comma-separated inputs (commas not allowed in single token expressions).
//  3. Reject inputs that start with "AS " or whose base parses to "AS" alone.
//  4. Parse an inline alias (e.g., "users.id AS uid") if present.
//  5. Apply an explicit alias override if provided, detecting conflicts with inline alias.
//  6. Populate name, Alias, input, and Error as needed.
//
// If any validation step fails, Error will be non-nil, and input will remain set to the original input.
//
// # Examples
//
//	// Valid simple input
//	b := NewBaseToken("users.id")
//	fmt.Println(b.name)  // → "users.id"
//	fmt.Println(b.Alias) // → ""
//	fmt.Println(b.Error) // → <nil>
//
//	// Inline alias
//	b = NewBaseToken("users.id AS uid")
//	fmt.Println(b.name)  // → "users.id"
//	fmt.Println(b.Alias) // → "uid"
//	fmt.Println(b.Error) // → <nil>
//
//	// Explicit alias override
//	b = NewBaseToken("users.id", "uid")
//	fmt.Println(b.name)  // → "users.id"
//	fmt.Println(b.Alias) // → "uid"
//	fmt.Println(b.Error) // → <nil>
//
//	// Alias conflict: inline "user_id" vs. explicit "uid"
//	b = NewBaseToken("users.id AS user_id", "uid")
//	fmt.Println(b.name)  // → "users.id"
//	fmt.Println(b.Alias) // → "uid"
//	fmt.Println(b.Error) // → alias conflict: explicit alias "uid" does not match inline alias "user_id"
//
//	// Invalid: starts with AS
//	b = NewBaseToken("AS uid")
//	fmt.Println(b.GetError()) // → invalid input expression: cannot start with 'AS'
//
//	// Invalid: empty input
//	b = NewBaseToken("")
//	fmt.Println(b.GetError()) // → invalid input expression: expression is empty
//
//	// Invalid: comma-separated fields
//	b = NewBaseToken("id, name")
//	fmt.Println(b.GetError()) // → invalid input expression: aliases must not be comma-separated
func NewBaseToken(input string, alias ...string) *BaseToken {
	t := &BaseToken{input: input}

	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		t.SetError(input, fmt.Errorf("invalid input expression: expression is empty"))
		return t
	}

	if strings.Contains(input, ",") {
		t.SetError(input, fmt.Errorf("invalid input expression: aliases must not be comma-separated"))
		return t
	}

	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "AS ") {
		t.SetError(input, fmt.Errorf("invalid input expression: cannot start with 'AS'"))
		return t
	}

	base, parsedAlias := ParseAlias(input)
	if strings.TrimSpace(base) == "AS" {
		t.SetError(input, fmt.Errorf("invalid input expression: name cannot be AS keyword"))
		return t
	}

	t.SetName(base)
	t.SetAlias(parsedAlias)

	if len(alias) > 0 && alias[0] != "" {
		explicit := alias[0]

		// If there was already an inline alias, and it doesn’t match the explicit one,
		// record an error (but still honor the explicit alias).
		if parsedAlias != "" && explicit != parsedAlias {
			t.SetError(input, fmt.Errorf(
				"alias conflict: explicit alias %q does not match inline alias %q",
				explicit,
				parsedAlias,
			))
		}
		// Finally, override whatever alias we had with the explicit one.
		t.SetAlias(explicit)
	}

	return t
}

// AliasOr returns the alias if it is non-empty, or else returns the name.
// If the receiver is nil, it returns an empty string.
//
// This method is useful when rendering column headers or result labels
// where aliases take precedence, but a fallback to name is still required.
//
// # Example
//
//	var b *BaseToken
//	fmt.Println(b.AliasOr()) // → ""
//
//	b = NewBaseToken("users.id AS uid")
//	fmt.Println(b.AliasOr()) // → "uid"
func (b *BaseToken) AliasOr() string {
	if b == nil {
		return ""
	}
	if b.alias != "" {
		return b.alias
	}
	return b.name
}

// GetAlias returns the alias if it is non-empty, or else returns the name.
// If the receiver is nil, it returns an empty string.
//
// This method is useful when rendering column headers or result labels
// where aliases take precedence, but a fallback to name is still required.
//
// # Example
//
//	var b *BaseToken
//	fmt.Println(b.AliasOr()) // → ""
//
//	b = NewBaseToken("users.id AS uid")
//	fmt.Println(b.AliasOr()) // → "uid"
func (b *BaseToken) GetAlias() string {
	if b == nil {
		return ""
	}
	return b.alias
}

// GetError returns the underlying error associated with the token.
// If the receiver is nil or no error has been set, it returns nil.
//
// This allows callers to inspect the exact error message or type
// without risking a nil-pointer panic.
//
// # Example
//
//	var b *BaseToken
//	fmt.Println(b.GetError()) // → <nil>
//
//	b = NewBaseToken("id AS uid")
//	fmt.Println(b.GetError()) // → <nil>
//
//	b.SetError("id AS uid", fmt.Errorf("name is missing"))
//	fmt.Println(b.GetError()) // → name is missing
func (b *BaseToken) GetError() error {
	if b == nil {
		return nil
	}
	return b.err
}

// GetInput returns the original raw input expression used to construct the token in a nil-safe way.
// If the receiver is nil, it returns an empty string.
//
// This accessor helps decouple the internal representation from external usage,
// and is useful for diagnostic or error-reporting routines.
//
// # Example
//
//	var b *BaseToken
//	fmt.Println(b.GetInput()) // → ""
//
//	b = NewBaseToken("users.id AS uid")
//	fmt.Println(b.GetInput()) // → "users.id AS uid"
func (b *BaseToken) GetInput() string {
	if b == nil {
		return ""
	}
	return b.input
}

// GetKind returns the Kind classification assigned to this token.
// If the receiver is nil or no kind has been set, it returns UnknownKind.
//
// # Example
//
//	b := NewBaseToken("")
//	fmt.Println(b.GetKind()) // → UnknownKind
//
//	b.SetKind(TableKind)
//	fmt.Println(b.GetKind()) // → TableKind
func (b *BaseToken) GetKind() contract2.Kind {
	if b == nil {
		return contract2.UnknownKind
	}
	return b.kind
}

// GetName returns the parsed name of the token in a nil-safe way.
// If the receiver is nil, it returns an empty string.
//
// This is commonly used in higher-level tokens (e.g., Column, Table)
// to retrieve the logical identifier associated with the token.
//
// # Example
//
//	var b *BaseToken
//	fmt.Println(b.GetName()) // → ""
//
//	b = NewBaseToken("users.id")
//	fmt.Println(b.GetName()) // → "users.id"
func (b *BaseToken) GetName() string {
	if b == nil {
		return ""
	}
	return b.name
}

// GetRaw returns the SQL raw expression representation: "name" or "name AS alias".
// If the receiver is nil, returns an empty string.
//
// This method does not perform any quoting or validation—it simply concatenates
// name and Alias as provided.
//
// # Example
//
//	b := NewBaseToken("users.id")
//	fmt.Println(b.Raw()) // → "users.id"
//
//	b = NewBaseToken("users.id AS uid")
//	fmt.Println(b.Raw()) // → "users.id AS uid"
func (b *BaseToken) GetRaw() string {
	if b == nil {
		return ""
	}
	if b.alias != "" {
		return fmt.Sprintf("%s AS %s", b.name, b.alias)
	}
	return b.name
}

// HasError reports whether the token has encountered a semantic or structural error.
//
// Typical causes include alias mismatches, unresolved references, or conflicting
// overrides detected during token construction or resolution.
func (b *BaseToken) HasError() bool {
	return b.IsErrored()
}

// IsErrored reports whether the token contains a non-nil Error.
// If the receiver is nil, returns false.
//
// Use this to quickly check if parsing or alias resolution failed,
// without inspecting the actual error value.
//
// # Example
//
//	b := NewBaseToken("users.id")
//	fmt.Println(b.IsErrored()) // → false
//
//	b.SetError("users.id", fmt.Errorf("no permission"))
//	fmt.Println(b.IsErrored()) // → true
func (b *BaseToken) IsErrored() bool {
	return b != nil && b.err != nil
}

// IsAliased reports whether the token has a non-empty Alias.
// If the receiver is nil, returns false.
//
// This is useful when deciding whether to include an AS clause in SQL rendering.
//
// # Example
//
//	b := NewBaseToken("users.id")
//	fmt.Println(b.IsAliased()) // → false
//
//	b = NewBaseToken("users.id AS uid")
//	fmt.Println(b.IsAliased()) // → true
func (b *BaseToken) IsAliased() bool {
	return b != nil && b.alias != ""
}

// IsValid returns true if the token has a non-empty name and no associated Error.
// If the receiver is nil, returns false.
//
// This is commonly used before including the token in a generated SQL query.
//
// # Example
//
//	b := NewBaseToken("users.id")
//	fmt.Println(b.IsValid()) // → true
//
//	b = NewBaseToken("")
//	fmt.Println(b.IsValid()) // → false
func (b *BaseToken) IsValid() bool {
	return b != nil && b.err == nil && strings.TrimSpace(b.name) != ""
}

// Raw returns the SQL raw expression representation: "name" or "name AS alias".
// If the receiver is nil, returns an empty string.
//
// This method does not perform any quoting or validation—it simply concatenates
// name and Alias as provided.
//
// # Example
//
//	b := NewBaseToken("users.id")
//	fmt.Println(b.Raw()) // → "users.id"
//
//	b = NewBaseToken("users.id AS uid")
//	fmt.Println(b.Raw()) // → "users.id AS uid"
func (b *BaseToken) Raw() string {
	return b.GetRaw()
}

// RenderAlias returns a dialect-quoted alias expression if an Alias is set,
// otherwise returns the qualified input unchanged. If the receiver is nil or
// qualified is empty, it returns qualified unchanged.
//
// If dialect is nil, it uses an unquoted format: "qualified AS alias".
// Otherwise, it applies dialect.QuoteIdentifier() to the alias.
//
// # Example
//
//	b := NewBaseToken("users.id AS uid")
//	fmt.Println(b.RenderAlias(driver.NewPostgresDialect(), "u.id"))
//	// → u.id AS "uid"
//
//	b = NewBaseToken("users.id")
//	fmt.Println(b.RenderAlias(driver.NewPostgresDialect(), "u.id"))
//	// → u.id
//
//	b = NewBaseToken("")
//	fmt.Println(b.RenderAlias(driver.NewPostgresDialect(), "u.id"))
//	// → u.id
func (b *BaseToken) RenderAlias(q contract2.Quoter, qualified string) string {
	if b == nil || qualified == "" {
		return qualified
	}
	if b.alias == "" {
		return qualified
	}
	if q == nil {
		return fmt.Sprintf("%s AS %s", qualified, b.alias)
	}
	return fmt.Sprintf("%s AS %s", qualified, q.QuoteIdentifier(b.alias))
}

// RenderName returns the dialect-quoted name of the token if non-empty.
// If the receiver is nil or name is empty, it returns an empty string.
// If dialect is nil, it returns name unquoted.
//
// This is commonly used when constructing SELECT, FROM, or JOIN clauses.
//
// # Example
//
//	b := NewBaseToken("id")
//	fmt.Println(b.RenderName(driver.NewPostgresDialect())) // → "id"
//
//	b = NewBaseToken("")
//	fmt.Println(b.RenderName(driver.NewPostgresDialect())) // → ""
//
//	b = NewBaseToken("id")
//	fmt.Println(b.RenderName(nil)) // → "id"
func (b *BaseToken) RenderName(q contract2.Quoter) string {
	if b == nil || b.name == "" {
		return ""
	}
	if q == nil {
		return b.name
	}
	return q.QuoteIdentifier(b.name)
}

// SetAlias replaces the parsed alias of the BaseToken.
// It updates both the private `alias` field and the deprecated exported `Alias` field
// for backward compatibility. After calling SetAlias, GetAlias() and AliasOr() will
// return the new alias. If the receiver is nil, SetAlias is a no-op.
//
// Example:
//
//	b := NewBaseToken("users.id AS u")
//	fmt.Println(b.GetAlias()) // Output: u
//
//	// Override the alias to "uid":
//	b.SetAlias("uid")
//	fmt.Println(b.GetAlias())  // Output: uid
//	fmt.Println(b.Raw())       // Output: users.id AS uid
//	fmt.Println(b.AliasOr())   // Output: uid
//
// Since: v1.7.0
func (b *BaseToken) SetAlias(alias string) {
	if b == nil {
		return
	}
	b.alias = alias
}

// SetError assigns a semantic or structural error to this token,
// along with the source expression string. It does not return any value,
// matching the Errorable contract. Any existing error is overwritten.
//
// # Example
//
//	    var b *BaseToken
//	    b.SetError("id AS uid", fmt.Errorf("name is missing"))
//	    // → No panic, b remains nil
//
//		b := NewBaseToken("id AS uid") // invalid: missing name before AS
//		b.SetError("id AS uid", fmt.Errorf("name is missing"))
//		fmt.Println(b.IsErrored())     // → true
//		fmt.Println(b.GetError())      // → name is missing
//		fmt.Println(b.GetInput())     // → "id AS uid"
func (b *BaseToken) SetError(source string, err error) {
	if b == nil {
		return
	}
	b.err = err
	if b.input != source {
		b.input = source
	}
}

// SetErrorWith assigns a semantic or structural error to the token, along with the source expression.
// It returns the same *BaseToken instance to allow fluent chaining. If the receiver is nil, this does nothing.
//
// Note: This method is provided for backward compatibility but SetError is preferred for new code.
//
// # Example
//
//	b := NewBaseToken("AS uid") // invalid: missing name before AS
//	b.SetErrorWith("AS uid", fmt.Errorf("name is missing before 'AS'"))
//	fmt.Println(b.IsErrored())    // → true
//	fmt.Println(b.GetInput())    // → "AS uid"
//	fmt.Println(b.GetError())     // → name is missing before 'AS'
func (b *BaseToken) SetErrorWith(source string, err error) *BaseToken {
	b.SetError(source, err)
	return b
}

// SetKind assigns the internal Kind classification (e.g., ColumnKind, TableKind) to this token.
// It is nil-safe: if the receiver is nil, this method does nothing.
//
// This should be called by higher-level token constructors (e.g., NewColumn, NewTable).
//
// # Example
//
//	    var b *BaseToken
//	    b.SetKind(ColumnKind)	// no panic; b remains nil
//
//		b := NewBaseToken("id")
//		b.SetKind(ColumnKind)
//		fmt.Println(b.GetKind()) 	// → ColumnKind
func (b *BaseToken) SetKind(k contract2.Kind) {
	if b == nil {
		return
	}
	b.kind = k
}

// SetName replaces the parsed identifier of the BaseToken.
// It updates both the private `name` field and the deprecated exported `Name` field
// for backward compatibility. After calling SetName, GetName() will return the new value.
// If the receiver is nil, SetName is a no-op.
//
// Example:
//
//	b := NewBaseToken("users.email AS u")
//	fmt.Println(b.GetName()) // Output: users.email
//
//	// Override the parsed name to just "email":
//	b.SetName("email")
//	fmt.Println(b.GetName()) // Output: email
//	// Raw SQL form now reflects the new name with the original alias:
//	fmt.Println(b.GetRaw())  // Output: email AS u
//
// Since: v1.7.0
func (b *BaseToken) SetName(name string) {
	if b == nil {
		return
	}
	b.name = name
}

// String returns a diagnostic string representation of the token, including its Kind, name,
// alias status, and any error present. This method is intended for logging, debugging, and
// test assertions only—it does not produce valid SQL. If the receiver is nil, it returns "nil".
//
// Format: Kind("name") [aliased: true/false, errored: true/false, error: <message>]
//
// # Examples
//
//	b := NewBaseToken("id")
//	b.SetKind(ColumnKind)
//	fmt.Println(b.String())
//	// → Column("id") [aliased: false, errored: false]
//
//	b = NewBaseToken("users.id AS uid")
//	b.SetKind(ColumnKind)
//	fmt.Println(b.String())
//	// → Column("id") [aliased: true, errored: false]
//
//	b = NewBaseToken("id AS uid")
//	b.SetKind(ColumnKind)
//	b.SetError("id AS uid", fmt.Errorf("alias conflict"))
//	fmt.Println(b.String())
//	// → Column("id") [aliased: true, errored: true, error: alias conflict]
func (b *BaseToken) String() string {
	if b == nil {
		return "nil"
	}

	suffix := fmt.Sprintf("[aliased: %t, errored: %t]", b.IsAliased(), b.IsErrored())
	if b.IsErrored() {
		suffix += fmt.Sprintf(", error: %s", b.GetError().Error())
	}
	return fmt.Sprintf("%s(\"%s\") %s", b.kind.String(), b.GetName(), suffix)
}
