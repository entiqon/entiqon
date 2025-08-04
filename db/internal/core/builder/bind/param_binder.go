// File: db/internal/core/builder/bind/param_binder.go

package bind

import "github.com/entiqon/entiqon/db/driver"

// ParamBinder is a utility that manages positional SQL placeholder generation
// and the corresponding argument list for dialect-aware query building.
type ParamBinder struct {
	dialect  driver.Dialect // dialect responsible for formatting placeholders
	args     []interface{}  // collected arguments to bind
	position int            // current positional index
}

// NewParamBinder creates a new ParamBinder with the given dialect.
// It initializes an empty argument list and sets the initial position to 0.
func NewParamBinder(dialect driver.Dialect) *ParamBinder {
	return &ParamBinder{
		dialect:  dialect,
		args:     make([]interface{}, 0),
		position: 1,
	}
}

// NewParamBinderWithPosition creates a ParamBinder with a custom starting
// position for placeholder indexing.
//
// This is useful when previous bindings (e.g., in SET clauses) have already
// consumed placeholder slots, and subsequent bindings (e.g., in WHERE clauses)
// must continue numbering without collisions.
//
// Example:
//
//	binder := NewParamBinderWithPosition(dialect, 2)
//	binder.Bind("active") // returns $3 (if dialect is Postgres)
//
// Since: v1.4.0
func NewParamBinderWithPosition(dialect driver.Dialect, position int) *ParamBinder {
	binder := NewParamBinder(dialect)
	binder.position = position
	return binder
}

// Bind registers a single value and returns its dialect-specific placeholder.
// For example, in Postgres dialect, the first call returns "$1", second "$2", etc.
func (pb *ParamBinder) Bind(value interface{}) string {
	placeholder := pb.dialect.Placeholder(pb.position)
	pb.args = append(pb.args, value)
	pb.position++
	return placeholder
}

// BindMany registers multiple values at once and returns a slice of placeholders,
// each corresponding to a value in order.
func (pb *ParamBinder) BindMany(values ...interface{}) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = pb.Bind(value)
	}
	return result
}

// Args returns the list of all arguments that have been bound so far.
// These should be passed to query execution methods alongside the built query.
func (pb *ParamBinder) Args() []interface{} {
	return pb.args
}
