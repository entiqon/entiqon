package builder

import "github.com/ialopezg/entiqon/dialect"

type ParamBinder struct {
	dialect  dialect.Dialect
	args     []interface{}
	position int
}

func NewParamBinder(d dialect.Dialect) *ParamBinder {
	return &ParamBinder{
		dialect: d,
		args:    make([]interface{}, 0),
	}
}

func (pb *ParamBinder) Bind(value interface{}) string {
	placeholder := pb.dialect.Placeholder(pb.position)
	pb.args = append(pb.args, value)
	pb.position++
	return placeholder
}

func (pb *ParamBinder) BindMany(values ...interface{}) []interface{} {
	result := make([]interface{}, len(values))
	for i, value := range values {
		result[i] = pb.Bind(value)
	}
	return result
}

func (pb *ParamBinder) Args() []interface{} {
	return pb.args
}
