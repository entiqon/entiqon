// File: db/internal/core/token/field.go

package token

// FieldToken represents a column or expression used in SQL statements.
//
// It allows customization of escaping, aliasing, and raw expression handling.
type FieldToken struct {
	// Name is the column name or raw expression (e.g., "id", "COUNT(*)").
	Name string

	// Alias is an optional alias for the field (used in SELECT as "AS <alias>").
	Alias string

	// IsRaw indicates whether the Name is a raw SQL expression and should not be escaped.
	IsRaw bool

	Value any
}
