package dialect

type Dialect interface {
	// EscapeIdentifier like table/column names
	EscapeIdentifier(identifier string) string
	// Placeholder syntax: ?, $1, :name, etc.
	Placeholder(index int) string
}
