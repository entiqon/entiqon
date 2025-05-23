// File: internal/core/builder/base.go
// Description: Provides shared dialect and error handling across query builders.
// Since: v1.4.0

package builder

import (
	"fmt"

	"github.com/ialopezg/entiqon/internal/core/driver"
	core "github.com/ialopezg/entiqon/internal/core/error"
)

// BaseBuilder provides shared dialect behavior and error handling for all query builders.
type BaseBuilder struct {
	dialect driver.Dialect
	// Stage-tagged error collector
	errors core.StageErrorCollector
}

// AddStageError records an error tied to a specific logical stage (e.g. "FROM", "WHERE").
func (b *BaseBuilder) AddStageError(stage core.StageToken, err error) {
	b.errors.AddStageError(stage, err)
}

// CombineErrors returns a grouped and readable error summary across all stages.
func (b *BaseBuilder) CombineErrors() error {
	return b.errors.CombineErrors()
}

// ErrorsByStage groups errors by their logical builder stage.
func (b *BaseBuilder) ErrorsByStage() map[string][]error {
	return b.errors.ErrorsByStage()
}

// GetDialect returns the resolved dialect, falling back to the generic default if unset.
func (b *BaseBuilder) GetDialect() driver.Dialect {
	if b.dialect == nil {
		b.dialect = driver.NewGenericDialect()
	}
	return b.dialect
}

// GetErrors returns all collected stage-tagged errors.
func (b *BaseBuilder) GetErrors() []core.StageError {
	return b.errors.Errors()
}

// GetErrorsString returns a stringified grouped summary of all collected errors.
func (b *BaseBuilder) GetErrorsString() string {
	return b.errors.String()
}

// HasDialect checks if a specific dialect has been set.
func (b *BaseBuilder) HasDialect() bool {
	return b.dialect != nil
}

// HasErrors returns true if any stage errors exist.
func (b *BaseBuilder) HasErrors() bool {
	return b.errors.HasErrors()
}

func (b *BaseBuilder) RenderFrom(table string, alias string) string {
	if alias != "" {
		return fmt.Sprintf("%s %s", b.dialect.QuoteIdentifier(table), alias)
	}
	return b.dialect.QuoteIdentifier(table)
}

// UseDialect sets a specific dialect to be used by the builder.
func (b *BaseBuilder) UseDialect(name string) *BaseBuilder {
	if name == "" || (b.dialect != nil && b.dialect.GetName() == name) {
		return b
	}
	if d := driver.ResolveDialect(name); d != nil {
		b.dialect = d
	}
	return b
}
