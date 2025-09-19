// File: db/internal/core/builder/base_test.go

package builder_test

import (
	"testing"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/internal/builder"
)

func TestBaseBuilder_HasDialect(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		dialect := driver.BaseDialect{}
		if got := dialect.GetName(); got != "base" {
			t.Errorf("expected %q, got %q", "base", got)
		}

		qb := &builder.BaseBuilder{}
		if qb.HasDialect() {
			t.Errorf("expected HasDialect=false, got true")
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		qb := &builder.BaseBuilder{}
		if qb.HasDialect() {
			t.Errorf("expected HasDialect=false, got true")
		}
	})
}
