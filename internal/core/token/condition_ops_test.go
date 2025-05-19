package token_test

import (
	"testing"

	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/ialopezg/entiqon/internal/core/token"
	"github.com/stretchr/testify/assert"
)

// ─────────────────────────────────────────────
// 🧪 ConditionBuilder Set
// ─────────────────────────────────────────────

func TestConditionBuilder_Set_NilReceiver(t *testing.T) {
	b := token.Condition{}
	out := b.Set(token.ConditionSimple, "x = ?", 1)
	assert.NotNil(t, out)
}

// ─────────────────────────────────────────────
// 🧪 FormatConditions
// ─────────────────────────────────────────────

func TestFormatConditions_Empty(t *testing.T) {
	var conditions []token.Condition
	sql, _, _ := token.FormatConditions(driver.NewPostgresDialect(), conditions)
	assert.Equal(t, "", sql)
}

func TestFormatConditions_Simple(t *testing.T) {
	conditions := []token.Condition{
		token.NewCondition(token.ConditionSimple, "id = ?", 1),
	}
	sql, _, _ := token.FormatConditions(driver.NewPostgresDialect(), conditions)
	assert.Equal(t, "id = ?", sql)
}

func TestFormatConditions_Complex(t *testing.T) {
	conditions := []token.Condition{
		token.NewCondition(token.ConditionSimple, "id = ?", 1),
		token.NewCondition(token.ConditionAnd, "email = ?", "x@test.dev"),
		token.NewCondition(token.ConditionOr, "status = 'active'"),
	}
	sql, _, _ := token.FormatConditions(driver.NewPostgresDialect(), conditions)
	assert.Equal(t, "id = ? AND email = ? OR status = 'active'", sql)
}

func TestFormatConditions_InvalidConditionType_ReturnsError(t *testing.T) {
	conditions := []token.Condition{
		{Type: "🔥", Key: "unexpected = true"},
	}
	_, _, err := token.FormatConditions(nil, conditions)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid condition type")
}
