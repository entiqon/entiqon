// Package field defines the low-level primitives used by the SQL builder.
//
// # Overview
//
// The token package provides fundamental structures to represent SQL query
// fragments in a dialect-agnostic way. These tokens are consumed by higher-
// level builders (e.g. SelectBuilder) to construct safe SQL statements.
//
// The primary type in this package is Field, which models a column or
// expression in a SELECT clause.
//
// # Field
//
// A Field represents a single column or expression in a SELECT query.
// It can optionally have an alias and may be marked as raw.
//
// Field supports multiple instantiation forms through New:
//
//   - New("expr")
//     A single expression, e.g. "id"
//
//   - New("expr alias")
//     Expression with alias, parsed by space, e.g. "id user_id"
//
//   - New("expr AS alias")
//     Expression with alias using AS, e.g. "id AS user_id"
//
//   - New("expr", "alias")
//     Expression and alias provided separately
//
//   - New("expr", "alias", true)
//     Expression and alias with IsRaw set explicitly
//
//   - New(*Field)
//     Disallowed: users must call Clone() instead
//
// # Field Behavior
//
// A Field consists of:
//   - Input: the raw user input
//   - Expr: the resolved expression (e.g. "id")
//   - Alias: the optional alias (e.g. "user_id")
//   - IsRaw: whether this is a raw expression
//   - Error: set if construction fails
//
// Methods include:
//   - Render: returns a dialect-agnostic SQL fragment
//   - Clone: produces a deep copy
//   - IsAliased, IsErrored, IsValid: convenience checks
//   - HasOwner: reports whether the field is qualified by a table name or alias
//   - Owner: returns the owning table name (or alias) if one is set
//   - SetOwner: assigns or clears the owning table name (passing nil clears it)
//
// # Usage Example
//
// Select specific columns:
//
//	sb := builder.NewSelect(nil).
//	    Fields("id", "name").
//	    Source("users")
//
//	sql, _ := sb.Build()
//	// SELECT id, name FROM users
//
// Add aliases:
//
//	sb := builder.NewSelect(nil).
//	    Fields("id user_id", "COUNT(*) total").
//	    Source("users")
//
//	sql, _ := sb.Build()
//	// SELECT id AS user_id, COUNT(*) AS total FROM users
package field
