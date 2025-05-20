// filename: /internal/core/driver/param_ops.go

package driver

// Bind adds a value to the parameter list and returns the appropriate placeholder.
//
// Updated: v1.4.0
func (pb *ParamBinder) Bind(value any) string {
	placeholder := pb.dialect.Placeholder(pb.position)
	pb.args = append(pb.args, value)
	pb.position++
	return placeholder
}

// BindMany adds multiple values to the parameter list and returns a slice of placeholders.
//
// Updated: v1.4.0
func (pb *ParamBinder) BindMany(values ...any) []string {
	placeholders := make([]string, len(values))
	for i, val := range values {
		placeholders[i] = pb.Bind(val)
	}
	return placeholders
}
