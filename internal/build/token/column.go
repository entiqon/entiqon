package token

import (
	"fmt"
	"strings"
)

// Column represents a SQL column token used in SELECT and other clauses.
//
// Fields:
//   - Name:  the name or expression of the column (required)
//   - Alias: an optional alias to rename the column in the result set
//   - Error: if non-nil, indicates a validation or parsing error encountered during construction
//
// Examples:
//   Column{Name: "user_id"}                   → SELECT user_id
//   Column{Name: "user_id", Alias: "id"}     → SELECT user_id AS id
//   Column{Error: err}                         → represents an invalid column

type Column struct {
	Name  string // Name is the column identifier or expression (required).
	Alias string // Alias is the optional column alias used in query output.
	Error error  // Error indicates a problem encountered during parsing or validation.
}

// NewColumn creates and returns a new Column token instance with optional alias.
//
// If no alias is explicitly passed, it attempts to parse inline aliasing from
// the expression using the SQL "AS" keyword. For example:
//
//	NewColumn("id")
//	  → Column{Name: "id"}
//
//	NewColumn("id AS user_id")
//	  → Column{Name: "id", Alias: "user_id"}
//
//	NewColumn("id", "user_id")
//	  → Column{Name: "id", Alias: "user_id"} (explicit alias overrides inline parsing)
func NewColumn(expr string, alias ...string) Column {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return Column{Error: fmt.Errorf("column name cannot be empty")}
	}

	if len(alias) == 0 {
		parts := strings.SplitN(expr, "AS", 2)
		if len(parts) == 2 {
			name := strings.TrimSpace(parts[0])
			aliasName := strings.TrimSpace(parts[1])
			if name == "" {
				return Column{Alias: aliasName, Error: fmt.Errorf("missing column name before alias")}
			}
			return Column{Name: name, Alias: aliasName}
		}
	}

	c := Column{Name: expr}
	if len(alias) > 0 {
		c.Alias = alias[0]
	}
	return c
}

// IsValid performs a basic check to ensure the column name is not empty or just whitespace.
// It returns false if any internal Error has been set.
//
// Example:
//
//	col := NewColumn("id")
//	valid := col.IsValid() // true
func (c Column) IsValid() bool {
	return c.Error == nil && strings.TrimSpace(c.Name) != ""
}

// Raw returns the raw SQL representation of the column, including aliasing via "AS" if set.
// It does not apply any quoting or dialect-specific formatting.
//
// Example:
//
//	col := NewColumn("id AS user_id")
//	sql := col.Raw() // "id AS user_id"
func (c Column) Raw() string {
	if c.Alias != "" {
		return fmt.Sprintf("%s AS %s", c.Name, c.Alias)
	}
	return c.Name
}

// String returns a developer-friendly representation of the column for logging or debugging.
// It includes quotes around names and aliases to make output unambiguous in logs.
//
// Example:
//
//	col := NewColumn("id AS user_id")
//	debug := col.String() // "Column(\"id\" AS \"user_id\")"
func (c Column) String() string {
	if c.Alias != "" {
		return fmt.Sprintf("Column(%q AS %q)", c.Name, c.Alias)
	}
	return fmt.Sprintf("Column(%q)", c.Name)
}

// Compile-time interface satisfaction check
var _ BaseToken = Column{}
