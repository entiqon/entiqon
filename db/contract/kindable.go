package contract

// Kindable is a generic contract for tokens or builders that expose a "kind"
// classification.
//
// A kind is an enumerated value that qualifies the nature of the token,
// such as whether a condition is Single, And, Or, or whether a field is
// an Identifier, Subquery, Literal, etc.
//
// Kindable is intentionally generic so each token defines its own strongly-typed
// enumeration and implements Kindable with that type parameter.
//
// Example:
//
//	type ConditionKind int
//
//	const (
//	    Invalid ConditionKind = iota
//	    Single
//	    And
//	    Or
//	)
//
//	type ConditionToken struct {
//	    kind ConditionKind
//	}
//
//	func (c *ConditionToken) Kind() ConditionKind   { return c.kind }
//	func (c *ConditionToken) SetKind(k ConditionKind) { c.kind = k }
//
//	var _ contract.Kindable[ConditionKind] = (*ConditionToken)(nil)
//
// This pattern allows builders to handle any token that can expose and mutate
// its classification in a uniform, type-safe way.
type Kindable[T any] interface {
	// Kind returns the current classification value.
	Kind() T

	// SetKind assigns the classification value.
	SetKind(T)
}
