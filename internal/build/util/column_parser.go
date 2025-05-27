// File: internal/build/util/column_parser.go

package util

import (
	"fmt"
	"strings"

	"github.com/ialopezg/entiqon/internal/build/token"
)

// ParseColumns receives one or more column strings, each of which may contain
// a single column or a comma-separated list of columns. It returns a list of
// Column tokens, each validated and enriched by the token.NewColumn logic.
//
// Invalid tokens will still be included in the result with their Error field set.
// This allows the caller to decide how to handle malformed column definitions.
//
// Examples:
//
//	ParseColumns("id")
//	  → [Column{Name: "id"}]
//
//	ParseColumns("id, name")
//	  → [Column{Name: "id"}, Column{Name: "name"}]
//
//	ParseColumns("id", "name AS customer")
//	  → [Column{Name: "id"}, Column{Name: "name", Alias: "customer"}]
func ParseColumns(input ...string) []*token.Column {
	var fields []*token.Column
	for _, arg := range input {
		parts := strings.Split(arg, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				fields = append(fields, token.NewErroredColumn(fmt.Errorf("empty column expression")))
				continue
			}
			col := token.NewColumn(trimmed)
			fields = append(fields, col)
		}
	}
	return fields
}
