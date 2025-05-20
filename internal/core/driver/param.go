// filename: /internal/core/driver/param.go

package driver

// ParamBinder manages SQL placeholder generation and parameter tracking in a dialect-safe way.
// It ensures consistent formatting for positional or named arguments depending on the dialect.
//
// Updated: v1.4.0
type ParamBinder struct {
	dialect  Dialect
	args     []any
	position int
}

// NewParamBinder creates a new parameter binder for the given dialect.
//
// Updated: v1.4.0
func NewParamBinder(dialect Dialect) *ParamBinder {
	return &ParamBinder{
		dialect:  dialect,
		args:     make([]any, 0),
		position: 1,
	}
}

// Args returns the collected bound parameters in insertion order.
//
// Updated: v1.4.0
func (pb *ParamBinder) Args() []any {
	return pb.args
}
