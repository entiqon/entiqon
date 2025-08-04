// File: db/internal/core/token/field_ops.go

package token

import "strings"

// As sets the alias for the FieldToken.
func (f FieldToken) As(alias string) FieldToken {
	f.Alias = alias
	return f
}

// FieldExpr creates a raw SQL expression with an optional alias (unescaped).
func FieldExpr(expression string, alias string) FieldToken {
	return FieldToken{
		Name:  expression,
		Alias: alias,
		IsRaw: true,
	}
}

// IsValid returns true if the field has a non-empty Name.
func (f FieldToken) IsValid() bool {
	return strings.TrimSpace(f.Name) != ""
}

// WithValue sets the bound value for the field and returns it.
func (f FieldToken) WithValue(value any) FieldToken {
	f.Value = value
	return f
}
