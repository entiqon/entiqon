// File: db/internal/build/render/column.go

package render

import (
	"fmt"

	"github.com/entiqon/entiqon/db/driver"
	"github.com/entiqon/entiqon/db/internal/build/token"
)

// Column renders a dialect-safe SQL column expression for use in SELECT, INSERT, or UPDATE clauses.
//
// This function handles:
//   - Quoting of the column name via the provided Dialect
//   - Table prefixing if the column is qualified
//   - Alias formatting using "AS" if an alias is present
//
// If the column is not valid or has semantic errors, an empty string is returned.
// Validation and error tracking are assumed to have occurred prior to rendering.
//
// This function is used only for output and does not perform logging or validation.
//
// # Examples
//
//	Column(driver, Column{Name: "id"})                     → "id"
//	Column(driver, Column{Name: "id", Alias: "uid"})       → "id AS uid"
//	Column(driver, Column{Name: "id", Qualified: "u"})     → "u.id"
//	Column(driver, Column{..., Qualified: "u", Alias: "uid"}) → "u.id AS uid"
func Column(d driver.Dialect, column token.Column) string {
	if !column.IsValid() || column.HasError() {
		return ""
	}

	if d == nil {
		d = driver.NewGenericDialect()
	}

	qualified := column.RenderName(d)
	if column.IsQualified() && column.Table != nil {
		base := d.QuoteIdentifier(column.Table.AliasOr())
		if base != "" {
			qualified = fmt.Sprintf("%s.%s", base, qualified)
		}
	}

	if column.IsAliased() {
		qualified = column.RenderAlias(d, qualified)
	}

	return qualified
}
