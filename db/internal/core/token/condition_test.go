// File: db/internal/core/token/condition_test.go

package token_test

import (
	"strings"
	"testing"
	"time"

	"github.com/entiqon/db/internal/core/builder"
	"github.com/entiqon/db/internal/core/token"
)

func TestCondition_SetWithParams(t *testing.T) {
	c := token.NewCondition(token.ConditionSimple, "active", true)
	if c.Key != "active" {
		t.Errorf("expected key %q, got %q", "active", c.Key)
	}
	if len(c.Values) != 1 || c.Values[0] != true {
		t.Errorf("expected values [true], got %#v", c.Values)
	}
	if !strings.Contains(c.Raw, "active = :active") {
		t.Errorf("expected Raw to contain %q, got %q", "active = :active", c.Raw)
	}
	if !c.IsValid() {
		t.Errorf("expected condition to be valid")
	}
}

func TestCondition_SetWithoutParams(t *testing.T) {
	c := token.NewCondition(token.ConditionSimple, "active = true")
	if c.Key != "active" {
		t.Errorf("expected key %q, got %q", "active", c.Key)
	}
	if len(c.Values) != 1 || c.Values[0] != true {
		t.Errorf("expected values [true], got %#v", c.Values)
	}
	if !strings.Contains(c.Raw, "active = :active") {
		t.Errorf("expected Raw to contain %q, got %q", "active = :active", c.Raw)
	}
	if !c.IsValid() {
		t.Errorf("expected condition to be valid")
	}
}

func TestCondition_IsValid(t *testing.T) {
	invalid := token.NewCondition(token.ConditionSimple, "active")
	if invalid.IsValid() {
		t.Errorf("expected invalid condition")
	}

	empty := token.NewCondition(token.ConditionSimple, "")
	if empty.IsValid() {
		t.Errorf("expected invalid empty condition")
	}
}

func TestCondition_AppendCondition(t *testing.T) {
	conditions := []token.Condition{
		token.NewCondition(token.ConditionSimple, "id", 1),
		token.NewCondition(token.ConditionAnd, "status", "active"),
	}
	if len(conditions) != 2 {
		t.Errorf("expected 2 conditions, got %d", len(conditions))
	}
	if !conditions[0].IsValid() || !conditions[1].IsValid() {
		t.Errorf("expected both conditions to be valid")
	}
}

func TestCondition_Constructors(t *testing.T) {
	now := time.Now()

	c := token.NewCondition(token.ConditionSimple, "value", 1, 2, 3)
	if !c.IsValid() {
		t.Errorf("expected valid condition")
	}

	c = token.NewCondition(token.ConditionSimple, "status = ")
	if c.IsValid() || c.Error == nil {
		t.Errorf("expected invalid condition with error")
	}

	c = token.NewConditionBetween(token.ConditionSimple, "created_at", now.Add(-time.Hour), now)
	if !c.IsValid() || c.Operator != "BETWEEN" {
		t.Errorf("expected BETWEEN operator")
	}

	c = token.NewConditionBetween(token.ConditionSimple, "created_at", nil, nil)
	if c.IsValid() {
		t.Errorf("expected invalid BETWEEN with nils")
	}

	c = token.NewConditionBetween(token.ConditionSimple, "created_at", "", "2024-01-01")
	if c.IsValid() || !strings.Contains(c.Error.Error(), "start value cannot be empty") {
		t.Errorf("expected invalid BETWEEN start value")
	}

	c = token.NewConditionBetween(token.ConditionSimple, "created_at", "2024-01-01", "")
	if c.IsValid() || !strings.Contains(c.Error.Error(), "end value cannot be empty") {
		t.Errorf("expected invalid BETWEEN end value")
	}

	c = token.NewConditionBetween(token.ConditionSimple, "created_at", "2024-01-01", 42)
	if c.IsValid() || !strings.Contains(c.Error.Error(), "compatible types") {
		t.Errorf("expected incompatible types error")
	}

	if token.NewCondition(token.ConditionSimple, "id", 123).Operator != "=" {
		t.Errorf("expected = operator")
	}
	if token.NewConditionGreaterThan(token.ConditionSimple, "score", 80).Operator != ">" {
		t.Errorf("expected > operator")
	}
	if token.NewConditionGreaterThanOrEqual(token.ConditionSimple, "points", 100).Operator != ">=" {
		t.Errorf("expected >= operator")
	}
	if token.NewConditionIn(token.ConditionSimple, "region", "US", "CA").Operator != "IN" {
		t.Errorf("expected IN operator")
	}
	if token.NewConditionIn(token.ConditionSimple, "amount", 10, "twenty").IsValid() {
		t.Errorf("expected invalid incompatible IN")
	}
	if token.NewConditionWithOperator(token.ConditionSimple, "", "=", 1).IsValid() {
		t.Errorf("expected invalid missing field")
	}
	if token.NewConditionWithOperator(token.ConditionSimple, "status", "=").IsValid() {
		t.Errorf("expected invalid missing values")
	}
	if token.NewConditionLessThan(token.ConditionSimple, "age", 65).Operator != "<" {
		t.Errorf("expected < operator")
	}
	if token.NewConditionLessThanOrEqual(token.ConditionSimple, "price", 100.0).Operator != "<=" {
		t.Errorf("expected <= operator")
	}
	if token.NewConditionLike(token.ConditionSimple, "name", "%John%").Operator != "LIKE" {
		t.Errorf("expected LIKE operator")
	}
	if token.NewConditionNotEqual(token.ConditionSimple, "status", "archived").Operator != "!=" {
		t.Errorf("expected != operator")
	}
	if token.NewConditionNotIn(token.ConditionSimple, "status", "inactive", "banned").Operator != "NOT IN" {
		t.Errorf("expected NOT IN operator")
	}
	if token.NewConditionNotIn(token.ConditionSimple, "id", "x", 1, true).IsValid() {
		t.Errorf("expected invalid NOT IN with incompatible types")
	}
}

func TestCondition_AppendConditionMixed(t *testing.T) {
	initial := []token.Condition{
		token.NewCondition(token.ConditionSimple, "id", 1),
	}

	valid := token.NewCondition(token.ConditionAnd, "status", "active")
	invalid := token.NewCondition(token.ConditionOr, "") // invalid

	result := builder.AppendCondition(initial, valid)
	if len(result) != 2 || result[1].Key != "status" {
		t.Errorf("expected second condition with key 'status'")
	}

	result = builder.AppendCondition(result, invalid)
	if len(result) != 2 {
		t.Errorf("expected unchanged length, got %d", len(result))
	}

	if !result[0].IsValid() || !result[1].IsValid() {
		t.Errorf("expected conditions to be valid")
	}
	if invalid.IsValid() {
		t.Errorf("expected invalid condition")
	}
}

func TestCondition_AreCompatibleTypes(t *testing.T) {
	if !token.AreCompatibleTypes("a", "b", "c") {
		t.Errorf("expected true for all strings")
	}
	if !token.AreCompatibleTypes(1, 2, 3) {
		t.Errorf("expected true for all ints")
	}
	if !token.AreCompatibleTypes(1, 2.5) {
		t.Errorf("expected true for int and float")
	}
	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	if !token.AreCompatibleTypes(t1, t2) {
		t.Errorf("expected true for time values")
	}
	if token.AreCompatibleTypes("a", 1) {
		t.Errorf("expected false for string and int")
	}
	if token.AreCompatibleTypes(time.Now(), "2023-01-01") {
		t.Errorf("expected false for time and string")
	}
	if token.AreCompatibleTypes() {
		t.Errorf("expected false for empty input")
	}
	if token.AreCompatibleTypes("only-one") {
		t.Errorf("expected false for single item")
	}
	var a any
	var b any = nil
	if token.AreCompatibleTypes(a, b) {
		t.Errorf("expected false for nil values")
	}
}
