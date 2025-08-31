// Package field defines the token.Field type, a low-level primitive used by the SQL builder.
//
// # Overview
//
// The token package provides fundamental structures to represent SQL query
// fragments in a dialect-agnostic way. These tokens are consumed by higher-
// level builders (e.g. SelectBuilder) to construct safe, expressive, and
// auditable SQL statements.
//
// The primary type in this package is Field, which models a column,
// subquery, function, computed expression, or literal in a SELECT clause.
//
// A Field is built on top of the BaseToken contract, ensuring consistent
// handling of input, expression, aliasing, raw detection, classification,
// validation, and error state.
//
// # Construction
//
// Fields are created using New(...) or NewWithTable(...):
//
//   - No argument:
//     field.New() → errored (empty input)
//
//   - Plain field:
//     field.New("id") → id
//
//   - Aliased (inline):
//     field.New("id user_id")    → id AS user_id
//     field.New("id AS user_id") → id AS user_id
//
//   - Aliased (explicit arguments):
//     field.New("id", "user_id") → id AS user_id
//     The second argument may also be any fmt.Stringer implementation.
//     Aliases are validated via identifier.IsValidAlias.
//
//   - Wildcard:
//     field.New("*") → *
//     field.New("* alias") → errored (wildcard cannot be aliased)
//
//   - Subquery (must have alias):
//     field.New("(SELECT COUNT(*) FROM users) AS t") → (SELECT COUNT(*) FROM users) AS t
//     field.New("(SELECT COUNT(*) FROM users)", "t") → (SELECT COUNT(*) FROM users) AS t
//
//   - Computed / Function / Literal:
//     field.New("price * quantity", "total") → (price * quantity) AS total
//     field.New("SUM(price)", "sum_price")  → SUM(price) AS sum_price
//     field.New("'hello'", "greeting")      → 'hello' AS greeting
//
//   - Errors:
//     field.New("")                  → errored
//     field.New("field alias extra") → errored (too many tokens)
//     field.New(field.New("id"))     → errored (use Clone() instead)
//     field.New(123)                 → errored (invalid type)
//
// # Contracts
//
// Field implements the following contracts from db/contract:
//
//   - BaseToken   → Input(), Expr(), Alias(), IsAliased(), ExpressionKind()
//   - Renderable  → Render()
//   - Rawable     → Raw(), IsRaw()
//   - Stringable  → String()
//   - Debuggable  → Debug()
//   - Clonable    → Clone()
//   - Errorable   → IsErrored(), Error()
//   - Validable   → IsValid()
//
// # Usage Example
//
// Select plain fields:
//
//	sb := builder.NewSelect(nil).
//	    Fields("id", "name").
//	    Source("users")
//
//	sql, _ := sb.Build()
//	// SELECT id, name FROM users
//
// With aliases:
//
//	sb := builder.NewSelect(nil).
//	    Fields("id user_id", "COUNT(*) total").
//	    Source("users")
//
//	sql, _ := sb.Build()
//	// SELECT id AS user_id, COUNT(*) AS total FROM users
//
// With subquery:
//
//	f := field.New("(SELECT COUNT(*) FROM users) AS t")
//	fmt.Println(f.String())
//	// field((SELECT COUNT(*) FROM users) AS t)
//
// Wildcard:
//
//	f := field.New("*")
//	fmt.Println(f.String())
//	// field(*)
//
// Invalid input:
//
//	f := field.New("id as user_id foo")
//	fmt.Println(f.Error())
//	// invalid format "id as user_id foo"
package field
