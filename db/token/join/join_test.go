package join_test

import (
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/token/join"
	"github.com/entiqon/entiqon/db/token/table"
)

func TestJoin(t *testing.T) {
	t.Run("Constructor", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			t.Run("InvalidJoinType", func(t *testing.T) {
				j := join.New("CRAsiNESS", nil, nil, "id = 1")
				if j.IsValid() {
					t.Errorf("expected invalid join, got valid: %v", j)
				}
				j = join.New(999, nil, nil, "id = 1")
				if j.IsValid() {
					t.Errorf("expected invalid join, got valid: %v", j)
				}
			})

			t.Run("TwoNil", func(t *testing.T) {
				j := join.New("INNER", nil, nil, "id = 1")
				if j.IsValid() {
					t.Errorf("expected invalid join, got valid: %v", j)
				}
			})

			t.Run("LeftNil", func(t *testing.T) {
				j := join.New("INNER", nil, "orders", "id = 1")
				if j.IsValid() {
					t.Errorf("expected invalid join, got valid: %v", j)
				}
			})

			t.Run("RightNil", func(t *testing.T) {
				j := join.New("INNER", "users", nil, "id = 1")
				if j.IsValid() {
					t.Errorf("expected invalid join, got valid: %v", j)
				}
			})

			t.Run("TwoErrored", func(t *testing.T) {
				bad1 := table.New("") // empty string → errored token
				bad2 := table.New("")
				j := join.New("INNER", bad1, bad2, "id = 1")
				if j.IsValid() {
					t.Errorf("expected invalid join, got valid: %v", j)
				}
			})

			t.Run("LeftErrored", func(t *testing.T) {
				bad := table.New("")
				j := join.New("INNER", bad, "orders", "id = 1")
				if j.IsValid() {
					t.Errorf("expected invalid join, got valid: %v", j)
				}
			})

			t.Run("RightErrored", func(t *testing.T) {
				bad := table.New("")
				j := join.New("INNER", "users", bad, "id = 1")
				if j.IsValid() {
					t.Errorf("expected invalid join, got valid: %v", j)
				}
			})

			t.Run("EmptyCondition", func(t *testing.T) {
				j := join.New("INNER", "users", "orders", "")
				if j.IsValid() {
					t.Errorf("expected invalid join, got valid: %v", j)
				}
			})

			t.Run("InvalidCall", func(t *testing.T) {
				j := join.New("INNER", 12345, "orders", "id = 1") // unsupported type
				if j.IsValid() {
					t.Errorf("expected invalid join due to unsupported type, got valid: %v", j)
				}
			})
		})

		t.Run("Default", func(t *testing.T) {
			j := join.New("INNER", "users", "orders", "users.id = orders.user_id")
			if j == nil {
				t.Fatal("expected non-nil Join")
			}
			if j.Kind() != join.Inner {
				t.Errorf("expected Inner, got %v", j.Kind())
			}
		})

		t.Run("Inner", func(t *testing.T) {
			j := join.NewInner("users", "orders", "users.id = orders.user_id")
			if j.Kind() != join.Inner {
				t.Errorf("expected Inner, got %v", j.Kind())
			}
		})

		t.Run("Left", func(t *testing.T) {
			j := join.NewLeft("users", "orders", "users.id = orders.user_id")
			if j.Kind() != join.Left {
				t.Errorf("expected Left, got %v", j.Kind())
			}
			j = join.New("LEFT JOIN", "users", "orders", "users.id = orders.user_id")
			if j.Kind() != join.Left {
				t.Errorf("expected Left, got %v", j.Kind())
			}
		})

		t.Run("Right", func(t *testing.T) {
			j := join.NewRight("users", "orders", "users.id = orders.user_id")
			if j.Kind() != join.Right {
				t.Errorf("expected Right, got %v", j.Kind())
			}
			j = join.New("RIGHT JOIN", "users", "orders", "users.id = orders.user_id")
			if j.Kind() != join.Right {
				t.Errorf("expected Right, got %v", j.Kind())
			}
		})

		t.Run("Full", func(t *testing.T) {
			j := join.NewFull("users", "orders", "users.id = orders.user_id")
			if j.Kind() != join.Full {
				t.Errorf("expected Full, got %v", j.Kind())
			}
			j = join.New("FULL JOIN", "users", "orders", "users.id = orders.user_id")
			if j.Kind() != join.Full {
				t.Errorf("expected Full, got %v", j.Kind())
			}
		})
	})

	t.Run("Contract", func(t *testing.T) {
		t.Run("Clonable", func(t *testing.T) {
			// construct valid join
			j := join.NewInner("users", "orders", "users.id = orders.user_id")
			clone := j.Clone()

			if clone == j {
				t.Error("expected Clone to return a different instance")
			}
			if clone.Kind() != j.Kind() {
				t.Errorf("expected same kind, got %v vs %v", clone.Kind(), j.Kind())
			}
		})

		t.Run("Debuggable", func(t *testing.T) {
			j := join.NewFull("users", "orders", "users.id = orders.user_id")
			d := j.Debug()
			if !strings.Contains(d, "Join{Kind:") {
				t.Errorf("expected debug string, got %q", d)
			}
		})

		t.Run("Errorable", func(t *testing.T) {
			j := join.New("SIDEWAYS", "users", "orders", "id = 1")
			if j.Error() == nil {
				t.Error("expected error on invalid kind")
			}
			if !j.IsErrored() {
				t.Error("expected IsErrored() = true")
			}
		})

		t.Run("Rawable", func(t *testing.T) {
			j := join.NewLeft("users", "orders", "users.id = orders.user_id")
			if !strings.HasPrefix(j.Raw(), "LEFT JOIN") {
				t.Errorf("expected LEFT JOIN in raw, got %q", j.Raw())
			}
		})

		t.Run("Renderable", func(t *testing.T) {
			j := join.NewInner("users u", "orders o", "u.id = o.user_id")
			sql := j.Render()
			if !strings.Contains(sql, "INNER JOIN") {
				t.Errorf("expected INNER JOIN in render, got %q", sql)
			}
		})

		t.Run("Stringable", func(t *testing.T) {
			j := join.NewRight("users", "orders", "users.id = orders.user_id")
			s := j.String()
			if !strings.Contains(s, "join(") {
				t.Errorf("expected string representation with 'join(', got %q", s)
			}
		})

		t.Run("Validable", func(t *testing.T) {
			j := join.NewInner("users", "orders", "users.id = orders.user_id")
			if !j.IsValid() {
				t.Error("expected valid join")
			}
		})
	})
}
