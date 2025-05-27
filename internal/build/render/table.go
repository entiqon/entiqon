package render

import (
	"fmt"

	"github.com/ialopezg/entiqon/driver"
	"github.com/ialopezg/entiqon/internal/build/token"
)

// Table renders a dialect-safe SQL representation of a table reference,
// including optional aliasing using "AS". This does not prepend "FROM", "INTO", etc.
//
// If the table is not valid or has an error, an empty string is returned.
//
// # Examples
//
//	Table(driver, Table{Name: "users"})                     → "users"
//	Table(driver, Table{Name: "users", Alias: "u"})         → "users AS u"
func Table(d driver.Dialect, tbl token.Table) string {
	if !tbl.IsValid() || tbl.HasError() {
		return ""
	}

	if d == nil {
		d = driver.NewGenericDialect()
	}

	name := d.QuoteIdentifier(tbl.Name)

	if tbl.IsAliased() {
		alias := d.QuoteIdentifier(tbl.Alias)
		return fmt.Sprintf("%s AS %s", name, alias)
	}

	return name
}
