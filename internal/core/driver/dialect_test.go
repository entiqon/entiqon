package driver_test

import (
	"testing"

	"github.com/ialopezg/entiqon/internal/core/driver"
	"github.com/ialopezg/entiqon/internal/core/token"
	"github.com/stretchr/testify/assert"
)

// ─────────────────────────────────────────────
// 🧪 BaseDialect
// ─────────────────────────────────────────────

func TestBaseDialect_Name(t *testing.T) {
	d := &driver.BaseDialect{DialectName: "test"}
	assert.Equal(t, "test", d.Name())
}

func TestBaseDialect_QuoteLiteral(t *testing.T) {
	d := &driver.BaseDialect{}

	assert.Equal(t, "'hello'", d.QuoteLiteral("hello"))
	assert.Equal(t, "42", d.QuoteLiteral(42))
	assert.Equal(t, "true", d.QuoteLiteral(true))
	assert.Equal(t, "'[1 2 3]'", d.QuoteLiteral([]int{1, 2, 3}))
}

func TestBaseDialect_QuoteIdentifier(t *testing.T) {
	d := &driver.BaseDialect{}
	assert.Equal(t, `"my_column"`, d.QuoteIdentifier("my_column"))
}

func TestBaseDialect_SupportsUpsert(t *testing.T) {
	d := &driver.BaseDialect{}
	assert.False(t, d.SupportsUpsert())
}

func TestBaseDialect_SupportsReturning(t *testing.T) {
	d := &driver.BaseDialect{}
	assert.False(t, d.SupportsReturning())
}

func TestBaseDialect_BuildLimitOffset(t *testing.T) {
	d := &driver.BaseDialect{}

	assert.Equal(t, "LIMIT 10 OFFSET 5", d.BuildLimitOffset(10, 5))
	assert.Equal(t, "LIMIT 20", d.BuildLimitOffset(20, -1))
	assert.Equal(t, "OFFSET 15", d.BuildLimitOffset(-1, 15))
	assert.Equal(t, "", d.BuildLimitOffset(-1, -1))
}

// ─────────────────────────────────────────────
// 🧪 PostgresDialect
// ─────────────────────────────────────────────

func TestPostgresDialect_NewInstance(t *testing.T) {
	pg := driver.NewPostgresDialect()
	assert.Equal(t, "postgres", pg.Name())
}

func TestPostgresDialect_SupportsUpsert(t *testing.T) {
	pg := driver.NewPostgresDialect()
	assert.True(t, pg.SupportsUpsert())
}

func TestPostgresDialect_SupportsReturning(t *testing.T) {
	pg := driver.NewPostgresDialect()
	assert.True(t, pg.SupportsReturning())
}

func TestPostgresDialect_QuoteIdentifier(t *testing.T) {
	pg := driver.NewPostgresDialect()
	assert.Equal(t, `"username"`, pg.QuoteIdentifier("username"))
}

// ─────────────────────────────────────────────
// 🧪 Dialect Resolver
// ─────────────────────────────────────────────

func TestResolveDialect_Postgres(t *testing.T) {
	d := driver.ResolveDialect("postgres")
	assert.Equal(t, "postgres", d.Name())
}

func TestResolveDialect_UnknownDialect(t *testing.T) {
	d := driver.ResolveDialect("not-real")
	assert.NotNil(t, d)
	assert.Equal(t, "generic", d.Name())
}

// ─────────────────────────────────────────────
// 🧪 FormatConditions
// ─────────────────────────────────────────────

func TestFormatConditions_MixedConditions(t *testing.T) {
	conditions := []token.Condition{
		token.NewCondition(token.ConditionSimple, "id = ?", 1),
		token.NewCondition(token.ConditionAnd, "email = ?", "x@test.dev"),
		token.NewCondition(token.ConditionOr, "active = true"),
	}
	sql, _, _ := token.FormatConditions(driver.NewPostgresDialect(), conditions)
	assert.Equal(t, "id = ? AND email = ? OR active = true", sql)
}

// ─────────────────────────────────────────────
// 🧪 WithValue (token.FieldToken)
// ─────────────────────────────────────────────

func TestFieldToken_WithValue(t *testing.T) {
	f := token.Field("email")
	f = f.WithValue("x@entiqon.dev")
	assert.Equal(t, "email", f.Name)
	assert.Equal(t, "x@entiqon.dev", f.Value)
}

// ─────────────────────────────────────────────
// 🧪 Test Dialect Mocks
// ─────────────────────────────────────────────

type testDialect struct{}

func (t *testDialect) Name() string                    { return "test" }
func (t *testDialect) QuoteIdentifier(s string) string { return "<<" + s + ">>" }
func (t *testDialect) QuoteLiteral(_ any) string       { return "!!" }

func TestTestDialect_Name(t *testing.T) {
	td := &testDialect{}
	assert.Equal(t, "test", td.Name())
}

func TestTestDialect_QuoteIdentifier(t *testing.T) {
	td := &testDialect{}
	assert.Equal(t, "<<foo>>", td.QuoteIdentifier("foo"))
}

func TestTestDialect_QuoteLiteral(t *testing.T) {
	td := &testDialect{}
	assert.Equal(t, "!!", td.QuoteLiteral("anything"))
}
