package builder

type Error struct {
	Token  string // e.g., "WHERE", "JOIN"
	Errors []error
}
