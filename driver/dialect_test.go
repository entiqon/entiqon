package driver_test

import (
	"testing"

	core "github.com/ialopezg/entiqon/driver"
	"github.com/ialopezg/entiqon/driver/styling"
	"github.com/stretchr/testify/assert"
)

func TestBaseDialectInterfaceSatisfaction(t *testing.T) {
	d := &core.BaseDialect{
		Name:             "test",
		Quotation:        styling.QuoteDouble,
		PlaceholderStyle: styling.PlaceholderDollar,
		EnableReturning:  true,
		EnableUpsert:     true,
	}

	// Satisfies interface
	var _ core.Dialect = d

	assert.Equal(t, "test", d.GetName())
	assert.Equal(t, `"field"`, d.QuoteIdentifier("field"))
	assert.Equal(t, "$1", d.Placeholder(1))
	assert.Equal(t, true, d.SupportsReturning())
	assert.Equal(t, true, d.SupportsUpsert())
	assert.NoError(t, d.Validate())
}
