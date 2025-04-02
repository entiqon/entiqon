package dialect

import "fmt"

type MySQLDialect struct{}

func (builder *MySQLDialect) EscapeIdentifier(identifier string) string {
	return fmt.Sprintf("`%s`", identifier)
}

func (builder *MySQLDialect) Placeholder(_ int) string {
	return "?"
}
