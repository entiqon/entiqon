// File: db/contract/renderable.go

// Package contract defines small, reusable behavioral contracts (interfaces)
// that core tokens and builders can implement to enable polymorphic behavior
// without tight coupling between packages.
package contract

// Cloanable is a generic contract for types that can produce a semantic clone
// of themselves. The type parameter T is the concrete return type of Clone.
//
// Notes:
//   - The name "Cloanable" is intentional per project conventions.
//   - Implementations should return a deep or semantic clone suitable for
//     independent mutation by callers.
type Cloanable[T any] interface {
	// Clone returns a new instance that is a semantic copy of the receiver.
	Clone() T
}
