package token

var reserved = map[string]struct{}{
	"AS": {}, "SELECT": {}, "FROM": {}, "WHERE": {}, "JOIN": {},
	"ON": {}, "GROUP": {}, "ORDER": {}, "BY": {}, "LIMIT": {},
	"INSERT": {}, "UPDATE": {}, "DELETE": {}, "INTO": {}, "VALUES": {},
	"CREATE": {}, "ALTER": {}, "DROP": {}, "TABLE": {}, "INDEX": {},
	// extend with the keywords you care about
}

// ExpressionKind classifies the semantic type of Field or Table.
//
// This unifies classification for tokens that represent values
// (Field) or sources (Table). Join tokens keep their own JoinKind.
//
//   - Invalid    — unrecognized expression
//   - Identifier — plain column reference (e.g. users.id)
//   - Computed   — computed expression (e.g. price * quantity)
//   - Literal    — constant value (e.g. 'abc', 42)
//   - Subquery   — nested SELECT used as a field or a table
//   - Function   — function call used as a field or a table
//   - Aggregate  — aggregate functions like SUM(), COUNT(), etc.
type ExpressionKind int

const (
	Invalid    ExpressionKind = iota // unrecognized expression
	Identifier                       // column reference
	Computed                         // computed expression
	Literal                          // constant value
	Subquery                         // subquery as field or table
	Function                         // function call as field or table
	Aggregate                        // aggregate functions like SUM(), COUNT(), etc.
)

// Alias returns the short two-letter code used when generating
// automatic aliases for this expression kind.
//
//	Identifier → "id"
//	Literal    → "lt"
//	Function   → "fn"
//	Aggregate  → "ag"
//	Computed   → "cp"
//	Subquery   → "sq"
//	Default    → "ex"
func (k ExpressionKind) Alias() string {
	switch k {
	case Identifier:
		return "id"
	case Literal:
		return "lt"
	case Function:
		return "fn"
	case Aggregate:
		return "ag"
	case Computed:
		return "cp"
	case Subquery:
		return "sq"
	default:
		return "ex"
	}
}

// String returns a human-readable label for the ExpressionKind.
func (k ExpressionKind) String() string {
	switch k {
	case Invalid:
		return "INVALID"
	case Identifier:
		return "IDENTIFIER"
	case Computed:
		return "COMPUTED"
	case Literal:
		return "LITERAL"
	case Subquery:
		return "SUBQUERY"
	case Function:
		return "FUNCTION"
	case Aggregate:
		return "AGGREGATE"
	default:
		return "INVALID"
	}
}

// IsValid reports whether the ExpressionKind is recognized.
func (k ExpressionKind) IsValid() bool {
	return k >= Identifier && k <= Aggregate
}
