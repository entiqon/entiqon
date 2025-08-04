// File: db/internal/test/test_dialect_test.go

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestDialect_Escape(t *testing.T) {
	d := &TestDialect{}
	assert.Equal(t, `"escape"`, d.Escape("escape"))
}

func TestTestDialect_Name(t *testing.T) {
	d := &TestDialect{}
	assert.Equal(t, "test", d.Name())
}
