package condition_test

import (
	"errors"
	"testing"

	"github.com/entiqon/entiqon/db/token/condition"
	ct "github.com/entiqon/entiqon/db/token/types/condition"
	"github.com/entiqon/entiqon/db/token/types/operator"
)

func TestCondition(t *testing.T) {
	t.Run("Constructor", func(t *testing.T) {
		t.Run("New", func(t *testing.T) {
			t.Run("Error", func(t *testing.T) {
				c := condition.New(ct.ParseFrom(99))
				if c.Error() == nil {
					t.Error("expected error, got nil")
				}
				if c.Error().Error() != "invalid condition type" {
					t.Error("expected invalid condition type, got ", c.Error())
				}

				c = condition.New(ct.Single)
				if c.Error().Error() != "invalid condition input: []" {
					t.Error("expected error, got none")
				}

				c = condition.New(ct.Single, 123456)
				if c.Error().Error() != "expr must be string, got int" {
					t.Error("expected expr must be string, got int")
				}

				c = condition.New(ct.Single, "")
				if c.Error().Error() != "empty expression" {
					t.Error("expected 'empty expression', got none")
				}

				c = condition.New(ct.Single, "id ++ 1")
				if c.Error().Error() != "invalid condition expression: \"id ++ 1\"" {
					t.Error("expected 'invalid condition expression: \"id ++ 1\"', got none")
				}
			})

			t.Run("Inline", func(t *testing.T) {
				c := condition.New(ct.Single, "id = 1")
				if c.Error() != nil {
					t.Error("expected no error, got ", c.Error())
				}

				c = condition.New(ct.Single, "id IN (1, 2, 3)")
				values := c.Value().([]interface{})
				if len(values) != 3 || values[0].(int) != 1 || values[1].(int) != 2 || values[2].(int) != 3 {
					t.Error("expected '1', got ", c.Value())
				}

				c = condition.New(ct.Single, "id BETWEEN 1 AND 100")
				if c.Error() != nil {
					t.Error("expected no error, got ", c.Error())
				}

				c = condition.New(ct.Single, "id IS NULL")
				if c.Value() != nil {
					t.Error("expected no error, got ", c.Error())
				}
			})

			t.Run("Default", func(t *testing.T) {
				t.Run("Equal", func(t *testing.T) {
					c := condition.New(ct.Single, "id", "1")
					if c.Error() != nil {
						t.Error("expected no error, got ", c.Error())
					}
				})
			})

			t.Run("WithParam", func(t *testing.T) {
				t.Run("Named", func(t *testing.T) {
					c := condition.New(ct.Single, "id = :id", "1")
					if c.Error() != nil {
						t.Error("expected no error, got ", c.Error())
					}
					if c.Input() != "id = :id 1" {
						t.Error("expected 'id = 1', got ", c.Input())
					}
					if c.Name() != "id" {
						t.Error("expected 'id', got ", c.Name())
					}
					if c.Expr() != "id = :id" {
						t.Error("expected 'id = :id', got ", c.Expr())
					}
					if c.Operator() != operator.Equal {
						t.Error("expected '=', got ", c.Operator())
					}
					if c.Value() != "1" {
						t.Error("expected '1', got ", c.Expr())
					}
				})
			})

			t.Run("WithOperator", func(t *testing.T) {
				c := condition.New(ct.Single, "id", operator.In, []int{1, 2, 3})
				if c.Error() != nil {
					t.Error("expected no error, got ", c.Error())
				}

				c = condition.New(ct.Single, "id", operator.ParseFrom(false), []int{1, 2, 3})
				if c.Error() == nil {
					t.Error("expected error, got nil")
				}

				c = condition.New(ct.Single, "id", operator.In, []int{})
				if c.Error() == nil {
					t.Error("expected error, got nil")
				}
			})
		})

		t.Run("NewAnd", func(t *testing.T) {
			c := condition.NewAnd("id", "1")
			if c.Error() != nil {
				t.Error("expected no error, got ", c.Error())
			}
			if c.Kind() != ct.And {
				t.Error("expected 'And', got ", c.Kind())
			}
		})

		t.Run("NewOr", func(t *testing.T) {
			c := condition.NewOr("id", "1")
			if c.Error() != nil {
				t.Error("expected no error, got ", c.Error())
			}
			if c.Kind() != ct.Or {
				t.Error("expected 'And', got ", c.Kind())
			}
		})
	})

	t.Run("Methods", func(t *testing.T) {
		t.Run("Contracts", func(t *testing.T) {
			t.Run("Kindable", func(t *testing.T) {
				c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
				if c.Kind() != ct.Single {
					t.Error("expected 'Single', got ", c.Kind())
				}

				c.SetKind(ct.And)
				if c.Kind() != ct.And {
					t.Error("expected 'And', got ", c.Kind())
				}
			})

			t.Run("Identifiable", func(t *testing.T) {
				c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
				if len(c.Input()) != 3 && c.Input() != "id > 10" {
					t.Error("expected 'id 1', got ", c.Input())
				}
				if c.Expr() != "id > :id" {
					t.Error("expected 'id > :id', got ", c.Expr())
				}
			})

			t.Run("Errorable", func(t *testing.T) {
				c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
				if c.Error() != nil {
					t.Error("expected error, got nil")
				}
				if c.IsErrored() {
					t.Errorf("expected false, got %t", c.IsErrored())
				}
				c.SetError(errors.New("test error"))
				if c.Error() == nil {
					t.Error("expected error, got nil")
				}
				if !c.IsErrored() {
					t.Error("expected true, got false")
				}
			})

			t.Run("Debuggable", func(t *testing.T) {
				c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
				if c.Debug() != "Condition{Input=\"id > 10\", Type:\"\", Expression=\"id > :id\", Value=10, Error=<nil>}" {
					t.Error("expected 'Condition{Input=\"id > 10\", Type:\"\", Expression=\"id > :id\", Value=10, Error=<nil>}', got ", c.Debug())
				}
			})

			t.Run("Rawable", func(t *testing.T) {
				c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
				if c.IsRaw() {
					t.Errorf("expected false, got %t", c.IsRaw())
				}
				if c.Raw() != c.Expr() {
					t.Error("expected 'id > 10', got ", c.Raw())
				}
			})

			t.Run("Renderable", func(t *testing.T) {
				c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
				if c.Render() != c.Expr() {
					t.Error("expected 'id > 10', got ", c.Render())
				}
			})

			t.Run("Renderable", func(t *testing.T) {
				c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
				if c.String() != "Condition(\"id > :id\"): name=\"id\", value=10, errored=false" {
					t.Error("expected 'Condition(\"id > :id\"): name=\"id\", value=10, errored=false', got ", c.String())
				}
			})

			t.Run("Validable", func(t *testing.T) {
				c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
				if !c.IsValid() {
					t.Error("expected true, got true")
				}
			})
		})

		t.Run("Name", func(t *testing.T) {
			c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
			if len(c.Input()) != 3 && c.Input() != "id > 10" {
				t.Error("expected 'id 1', got ", c.Input())
			}
		})

		t.Run("Operator", func(t *testing.T) {
			c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
			if c.Operator() != operator.GreaterThan {
				t.Error("expected 'GreaterThan', got ", c.Operator())
			}
		})

		t.Run("Value", func(t *testing.T) {
			c := condition.New(ct.Single, "id", operator.GreaterThan, 10)
			if c.Value() != 10 {
				t.Error("expected 10, got ", c.Value())
			}
		})
	})
}
