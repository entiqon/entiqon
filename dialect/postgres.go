package dialect

import "fmt"

type PostgresDialect struct{}

func (builder *PostgresDialect) EscapeIdentifier(name string) string {
	return fmt.Sprintf(`"%s"`, name)
}

func (builder *PostgresDialect) Placeholder(index int) string {
	return fmt.Sprintf("$%d", index+1)
}
