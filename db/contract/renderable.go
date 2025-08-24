// File: db/contract/renderable.go
//
// Renderable defines machine-facing SQL output. See package-level
// documentation in doc.go for an overview of all contracts and
// their distinct purposes.

package contract

// Renderable defines the contract for objects that can produce a canonical,
// machine-facing string representation of themselves.
//
// Render() is expected to return a stable, context-independent string
// suitable for use in query builders, SQL fragments, or any serialization
// where accuracy and consistency matter.
//
// Example:
//
//	// Machine-facing: safe for SQL generation
//	fmt.Println(field.Render()) // Output: "users.id"
type Renderable interface {
	// Render returns the canonical, machine-facing representation
	// of the object (e.g., SQL fragment). Must be stable and
	// consistent for use in query building.
	Render() string
}
