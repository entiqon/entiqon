// File: db/builder/base.go
// Description: Provides shared dialect and error handling across query builders.
// Since: v1.4.0

package builder

import (
	"fmt"
	"strings"

	"github.com/entiqon/db/driver"
	error2 "github.com/entiqon/db/internal/core/errors"
)

// BaseBuilder provides shared dialect behavior and error handling for all query builders.
type BaseBuilder struct {
	Dialect   driver.Dialect
	Name      string
	Validator error2.StageErrorCollector

	dialect driver.Dialect
	// Stage-tagged error collector
	errors error2.StageErrorCollector
}

// NewBaseBuilder creates a properly initialized BaseBuilder with a specified name.
//
// This constructor ensures:
//   - The builder name is always set and lowercased (e.g., "select", "delete")
//   - If the provided dialect is nil, the builder defaults to using the generic dialect.
//   - The Validator is initialized (ready for AddStageError, etc.)
//   - Safe behavior in Validate() and error diagnostics
//
// Example:
//
//	b := NewBaseBuilder("select", driver.NewPostgresDialect())
//
// Since: v1.4.0
func NewBaseBuilder(name string, dialect driver.Dialect) BaseBuilder {
	if dialect == nil {
		dialect = driver.NewGenericDialect()
	}
	return BaseBuilder{
		Name:      strings.ToLower(name),
		Dialect:   dialect,
		Validator: error2.StageErrorCollector{},
	}
}

// AddStageError records an error tied to a specific logical stage (e.g. "FROM", "WHERE").
func (b *BaseBuilder) AddStageError(stage error2.StageToken, err error) {
	b.Validator.AddStageError(stage, err)
}

// CombineErrors returns a grouped and readable error summary across all stages.
func (b *BaseBuilder) CombineErrors() error {
	return b.Validator.CombineErrors()
}

// ErrorsByStage groups errors by their logical builder stage.
func (b *BaseBuilder) ErrorsByStage() map[string][]error {
	return b.Validator.ErrorsByStage()
}

// GetDialect returns the resolved dialect, falling back to the generic default if unset.
func (b *BaseBuilder) GetDialect() driver.Dialect {
	if b.Dialect == nil {
		b.Dialect = driver.NewGenericDialect()
	}
	return b.Dialect
}

// HasDialect checks if a specific dialect has been set.
func (b *BaseBuilder) HasDialect() bool {
	return b.Dialect != nil
}

// HasErrors returns true if any stage errors exist.
func (b *BaseBuilder) HasErrors() bool {
	return b.Validator.HasErrors()
}

func (b *BaseBuilder) RenderFrom(table string, alias string) string {
	if alias != "" {
		return fmt.Sprintf("%s %s", b.Dialect.QuoteIdentifier(table), alias)
	}
	return b.Dialect.QuoteIdentifier(table)
}

// UseDialect sets a specific dialect to be used by the builder.
func (b *BaseBuilder) UseDialect(name string) *BaseBuilder {
	if name == "" || (b.Dialect != nil && b.Dialect.GetName() == name) {
		return b
	}
	if d := driver.ResolveDialect(name); d != nil {
		b.Dialect = d
	}
	return b
}

// Validate checks that a dialect is assigned to the builder.
// If not, returns a formatted error using the builder's Name field.
//
// Usage:
//
//	if err := b.Validate(); err != nil {
//	    return "", nil, err
//	}
func (b *BaseBuilder) Validate() error {
	if b.Dialect == nil {
		return fmt.Errorf(
			"%s: no dialect set — please assign one (e.g., driver.NewGenericDriver())",
			strings.ToUpper(b.Name),
		)
	}
	if b.Validator.HasErrors() {
		return b.Validator.CombineErrors()
	}
	return nil
}
