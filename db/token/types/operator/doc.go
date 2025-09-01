// Package operator provides the set of SQL comparison and predicate operators
// supported by Entiqon condition tokens (WHERE, ON, HAVING). It offers:
//
//   - A typed enum (Type) for operators like "=", "!=", "IN", "BETWEEN",
//     "IS NULL", "IS DISTINCT FROM", etc.
//   - Robust parsing via ParseFrom (symbols, words, and short aliases).
//   - Canonical string forms via String() and short mnemonics via Alias().
//   - Discovery of supported operators via GetKnownOperators(), returned in
//     longest-first order to aid text scanning resolvers.
//
// Typical usage:
//
//	// Parse a user-supplied operator (case/space-insensitive):
//	op := operator.ParseFrom("  not   in ")
//	if !op.IsValid() {
//	    // handle unsupported operator
//	}
//	fmt.Println(op.String()) // "NOT IN"
//	fmt.Println(op.Alias())  // "nin"
//
// Integration notes:
//
//   - GetKnownOperators returns operator spellings sorted by descending length.
//     This is useful when scanning raw expressions so multi-word operators
//     (e.g., "IS NOT DISTINCT FROM") are matched before shorter prefixes
//     (e.g., "IS").
//
//   - The resolver layer (e.g., helpers.ResolveCondition) is expected to keep
//     its own normalization model. A common pattern is to canonicalize the LHS
//     as "field = :field" while preserving the actual operator in a separate
//     Type value and the RHS as a typed payload (scalar, []any, or nil).
package operator
