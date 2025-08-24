// File: db/contract/clonable.go
//
// Clonable defines semantic clone behavior for safe duplication of
// tokens or builders. See package-level documentation in doc.go
// for an overview of all contracts and their distinct purposes.

package contract

// Clonable is a generic contract for types that can produce a semantic clone
// of themselves. The type parameter T is the concrete return type of Clone.
//
// Notes:
//   - The name "Clonable" is intentional per project conventions.
//   - Implementations should return a deep or semantic clone suitable for
//     independent mutation by callers.
//
// Example:
//
//	c1 := field.New("users.id", "uid")
//	c2 := c1.Clone()
//	// c2 is safe to mutate independently of c1
type Clonable[T any] interface {
	Clone() T
}
