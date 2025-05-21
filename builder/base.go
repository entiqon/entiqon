// filename: builder/base.go

package builder

import (
	"github.com/ialopezg/entiqon/internal/core/builder"
	"github.com/ialopezg/entiqon/internal/core/driver"
)

// BaseBuilder provides shared dialect behavior for all query builders.
// It embeds the dialect reference and ensures default resolution when building queries.
//
// Since: v1.4.0
type BaseBuilder struct {
	dialect driver.Dialect
	errors  []builder.Error
}

// AddStageError appends a validation error to the given stage of the builder,
// such as "WHERE", "JOIN", or "LIMIT". Multiple errors may be collected per stage.
//
// Updated: v1.4.0
func (bb *BaseBuilder) AddStageError(stage string, err error) {
	for i := range bb.errors {
		if bb.errors[i].Token == stage {
			bb.errors[i].Errors = append(bb.errors[i].Errors, err)
			return
		}
	}
	bb.errors = append(bb.errors, builder.Error{
		Token:  stage,
		Errors: []error{err},
	})
}

// GetDialect returns the resolved dialect for the builder.
//
// If no dialect was explicitly set, it defaults to the generic dialect.
// This method guarantees that all builders can safely access a usable dialect at Build() time.
//
// Since: v1.4.0
func (bb *BaseBuilder) GetDialect() driver.Dialect {
	if bb.dialect == nil {
		bb.dialect = driver.NewGenericDialect()
	}
	return bb.dialect
}

// GetErrors returns the full list of stage-tagged errors collected during the build process.
//
// Updated: v1.4.0
func (bb *BaseBuilder) GetErrors() []builder.Error {
	return bb.errors
}

// HasDialect returns true if a dialect has been explicitly set.
//
// Since: v1.4.0
func (bb *BaseBuilder) HasDialect() bool {
	return bb.dialect != nil
}

// HasErrors returns true if any error has been recorded in the builder.
//
// Updated: v1.4.0
func (bb *BaseBuilder) HasErrors() bool {
	return len(bb.errors) > 0
}

// UseDialect sets a specific dialect to be used by the builder.
//
// If the provided dialect is nil or named "generic", the call has no effect.
// This method allows fluent overrides for specific SQL engines (e.g., PostgreSQL, ClickHouse).
//
// Updated: v1.4.0
func (bb *BaseBuilder) UseDialect(name string) *BaseBuilder {
	if name == "" || (bb.dialect != nil && bb.dialect.Name() == name) {
		return bb
	}
	if d := driver.ResolveDialect(name); d != nil {
		bb.dialect = d
	}
	return bb
}
