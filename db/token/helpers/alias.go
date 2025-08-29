package helpers

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

// reserved defines a minimal, dialect-agnostic set of SQL keywords
// that cannot be used as identifiers or aliases.
// Future dialect-specific packages may extend or override this list.
var reserved = map[string]struct{}{
	"AS": {}, "SELECT": {}, "FROM": {}, "WHERE": {}, "JOIN": {},
	"ON": {}, "GROUP": {}, "ORDER": {}, "BY": {}, "LIMIT": {},
	"INSERT": {}, "UPDATE": {}, "DELETE": {}, "INTO": {}, "VALUES": {},
	"CREATE": {}, "ALTER": {}, "DROP": {}, "TABLE": {}, "INDEX": {},
}

// ValidateAlias checks whether s is a valid SQL alias.
// Returns nil if valid, or an error describing why it is invalid.
//
// Rules:
//   - Must be a valid identifier (validated via ValidateIdentifier).
//     This may produce errors such as:
//   - "identifier cannot be empty"
//   - "identifier cannot start with digit"
//   - "invalid identifier syntax: ..."
//   - Must not be a reserved keyword (case-insensitive).
//
// Future: reserved keyword rules will be extended or delegated
// to dialect-specific packages (e.g. Postgres, MySQL).
func ValidateAlias(s string) error {
	if err := ValidateIdentifier(s); err != nil {
		return fmt.Errorf("invalid alias %w", err)
	}
	if _, found := reserved[strings.ToUpper(s)]; found {
		return fmt.Errorf("alias is a reserved keyword: %q", s)
	}
	return nil
}

// IsValidAlias reports whether s is a valid alias.
// It is a convenience wrapper around ValidateAlias, returning true
// if the alias is valid and false otherwise.
// Prefer ValidateAlias when the reason for invalidation is needed.
func IsValidAlias(s string) bool {
	return ValidateAlias(s) == nil
}

// ValidateTrailingAlias checks if the last token in expr is a valid
// trailing alias (when no explicit AS is used).
//
// Returns the alias string if valid, or an error if invalid.
//
// Rules:
//   - Expressions with "AS" are rejected (explicit aliases handled elsewhere).
//   - A single-token expression has no alias candidate.
//   - If the token before the last is an operator, the last token is
//     part of the expression, not an alias.
//   - Otherwise, the last token must be a valid alias.
func ValidateTrailingAlias(expr string) (string, error) {
	up := strings.ToUpper(expr)
	if strings.Contains(up, " AS ") {
		return "", fmt.Errorf("explicit AS found, not a trailing alias")
	}

	tokens := strings.Fields(expr)
	if len(tokens) <= 1 {
		return "", fmt.Errorf("no trailing alias candidate")
	}

	last := tokens[len(tokens)-1]
	penultimate := tokens[len(tokens)-2]

	operators := map[string]bool{"||": true, "+": true, "-": true, "*": true, "/": true}
	if operators[penultimate] {
		return "", fmt.Errorf("last token %q is part of expression, not alias", last)
	}

	if err := ValidateAlias(last); err != nil {
		return "", fmt.Errorf("invalid trailing alias %q: %w", last, err)
	}
	return last, nil
}

// HasTrailingAlias reports whether expr has a valid trailing alias.
// It is a convenience wrapper over ValidateTrailingAlias, returning
// true if a trailing alias exists and false otherwise.
func HasTrailingAlias(expr string) bool {
	_, err := ValidateTrailingAlias(expr)
	return err == nil
}

// ReservedKeywords returns the set of reserved SQL keywords used
// by ValidateAlias. The returned slice is a copy of the internal set
// and can be used safely in tests or diagnostics without risk of
// modification.
//
// The list is dialect-agnostic and intentionally minimal; dialect
// packages may extend it with additional keywords.
func ReservedKeywords() []string {
	out := make([]string, 0, len(reserved))
	for k := range reserved {
		out = append(out, k)
	}
	return out
}

// GenerateAlias returns a deterministic alias for an expression.
//
// The alias is constructed using the provided short code (kind alias)
// and a SHA-1 hash of the expression string (first 6 hex chars).
//
// Example:
//
//	helpers.GenerateAlias("fn", "SUM(price)") → "fn_a1b2c3"
//
// Rules:
//   - The prefix should come from ExpressionKind.Alias().
//   - The result is always a safe SQL identifier (letters, digits, underscore).
//   - Same input → same alias (deterministic).
//   - Different inputs → different aliases with high collision resistance.
//
// Future: dialect-specific packages may override this generator
// with shorter, descriptive, or quoted variants.
func GenerateAlias(prefix, expr string) string {
	h := sha1.New()
	h.Write([]byte(expr))
	sum := hex.EncodeToString(h.Sum(nil))
	return prefix + "_" + sum[:6]
}
