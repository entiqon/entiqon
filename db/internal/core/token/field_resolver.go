// File: db/internal/core/token/field_resolver.go

package token

import (
	"strings"
)

// Field resolves a single column name, optionally with alias.
//
// Valid:
//
//	Field("id")                          → FieldToken{GetName: "id"}
//	Field("first_name AS name")         → FieldToken{GetName: "first_name", Alias: "name"}
//	Field("first_name", "name")         → FieldToken{GetName: "first_name", Alias: "name"}
//
// NOT allowed:
//
//	Field("id, name")                   → caller must split and call separately
func Field(parts ...string) FieldToken {
	if len(parts) == 2 {
		return FieldToken{
			Name:  strings.TrimSpace(parts[0]),
			Alias: strings.TrimSpace(parts[1]),
		}
	}

	expr := strings.TrimSpace(parts[0])
	if strings.Contains(expr, ",") {
		panic("Field: comma-separated values not allowed in a single call. Call Field(...) separately for each.")
	}

	sub := strings.SplitN(expr, " AS ", 2)
	if len(sub) == 2 {
		return FieldToken{
			Name:  strings.TrimSpace(sub[0]),
			Alias: strings.TrimSpace(sub[1]),
		}
	}

	return FieldToken{Name: expr}
}

// FieldsFromExpr splits a comma-separated string of field expressions
// and returns a slice of resolved FieldToken entries.
//
// Example:
//
//	FieldsFromExpr("id, name AS alias") =>
//	  []FieldToken{Field("id"), Field("name AS alias")}
func FieldsFromExpr(expr string) []FieldToken {
	var fields []FieldToken

	for _, raw := range strings.Split(expr, ",") {
		token := Field(strings.TrimSpace(raw))
		fields = append(fields, token)
	}

	return fields
}
