package token

import (
	"errors"
	"strings"
)

// ResolveExpr splits an input into expr + alias, then classifies the expr.
// If allowAlias=false, alias parts are rejected.
func ResolveExpr(input string, allowAlias bool) (kind ExpressionKind, expr, alias string, err error) {
	in := strings.TrimSpace(input)
	if in == "" {
		return -1, "", "", errors.New("empty input")
	}

	// ✅ Subquery detection: treat whole input as one expression
	// even if it contains spaces.
	if strings.HasPrefix(in, "(") && strings.HasSuffix(in, ")") {
		return Subquery, in, "", nil
	}

	// If explicit "AS" exists → split at it
	up := strings.ToUpper(in)
	if strings.Contains(up, " AS ") {
		parts := strings.SplitN(in, " AS ", 2)
		expr = strings.TrimSpace(parts[0])
		alias = strings.TrimSpace(parts[1])
		if !allowAlias {
			return -1, "", "", errors.New("alias not allowed: " + input)
		}
		if !IsValidAlias(alias) {
			return -1, "", "", errors.New("invalid alias: " + alias)
		}
		return ClassifyExpression(expr), expr, alias, nil
	}

	// Otherwise, split by whitespace
	parts := strings.Fields(in)
	switch len(parts) {
	case 1:
		expr = parts[0]
	case 2:
		if !allowAlias {
			return -1, "", "", errors.New("alias not allowed: " + input)
		}
		expr, alias = parts[0], parts[1]
		if !IsValidAlias(alias) {
			return -1, "", "", errors.New("invalid alias: " + alias)
		}
	default:
		// Handle cases like "(SELECT * FROM users) u"
		if HasTrailingAliasWithoutAS(in) {
			alias = parts[len(parts)-1]
			expr = strings.Join(parts[:len(parts)-1], " ")
			if !allowAlias {
				return -1, "", "", errors.New("alias not allowed: " + input)
			}
			if !IsValidAlias(alias) {
				return -1, "", "", errors.New("invalid alias: " + alias)
			}
		} else {
			return -1, "", "", errors.New("invalid input: " + input)
		}
	}

	kind = ClassifyExpression(expr)

	// Strict check: Identifiers must be a single token
	if kind == Identifier {
		coreParts := strings.Fields(expr)
		if len(coreParts) != 1 {
			return -1, expr, alias, errors.New("invalid format " + expr)
		}
	}

	return kind, expr, alias, nil
}
