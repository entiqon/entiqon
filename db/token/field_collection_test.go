package token_test

import (
	"testing"

	"github.com/entiqon/entiqon/db/token"
)

func TestFieldCollection(t *testing.T) {
	t.Run("Methods", func(t *testing.T) {
		t.Run("Add", func(t *testing.T) {
			t.Run("SingleElement", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}

				got := fc.Add(a)
				// Chainability check
				if got != &fc {
					t.Fatal("Add should return the same receiver for chaining")
				}
				// Expected: a
				want := []token.Field{a}
				if fc.Length() != len(want) {
					t.Fatalf("len = %d, want %d", fc.Length(), len(want))
				}
				for i := range want {
					if fc[i] != want[i] {
						t.Fatalf("Element mismatch at index %d: got %v, want %v", i, fc[i], want[i])
					}
				}
			})

			t.Run("MultipleElements", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}
				c := token.Field{Expr: "c"}

				fc.Add(a, b, c)
				want := []token.Field{a, b, c}
				if fc.Length() != len(want) {
					t.Fatalf("len = %d, want %d", fc.Length(), len(want))
				}
				for i := range want {
					if fc[i] != want[i] {
						t.Fatalf("Element mismatch at index %d: got %v, want %v", i, fc[i], want[i])
					}
				}
			})

			t.Run("AppendToExisting", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}
				c := token.Field{Expr: "c"}

				fc.Add(a).Add(b, c)

				want := []token.Field{a, b, c}
				if fc.Length() != len(want) {
					t.Fatalf("len = %d, want %d", fc.Length(), len(want))
				}
				for i := range want {
					if fc[i] != want[i] {
						t.Fatalf("Element mismatch at index %d: got %v, want %v", i, fc[i], want[i])
					}
				}
			})

			t.Run("EmptyArgs", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}

				fc.Add(a).Add() // second Add has no args, should be no-op

				want := []token.Field{a}
				if fc.Length() != len(want) {
					t.Fatalf("len = %d, want %d", fc.Length(), len(want))
				}
				if fc[0] != a {
					t.Fatalf("Element mismatch: got %v, want %v", fc[0], a)
				}
			})

			t.Run("NilReceiver", func(t *testing.T) {
				var nilFC *token.FieldCollection
				got := nilFC.Add(token.Field{Expr: "a"})
				if got != nil {
					t.Fatal("Add on nil receiver should return nil")
				}
			})
		})

		t.Run("Clear", func(t *testing.T) {
			t.Run("EmptiesCollection", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}

				fc.Add(a, b)
				if fc.Length() != 2 {
					t.Fatalf("precondition failed: len = %d, want 2", fc.Length())
				}

				fc.Clear()

				if fc.Length() != 0 {
					t.Fatalf("Clear should empty the collection, got len = %d", fc.Length())
				}
				if !fc.IsEmpty() {
					t.Fatal("IsEmpty should return true after Clear")
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				var fc token.FieldCollection
				if fc.Length() != 0 {
					t.Fatalf("precondition failed: len = %d, want 0", fc.Length())
				}

				fc.Clear() // should not panic or change anything

				if fc.Length() != 0 {
					t.Fatalf("Clear should keep collection empty, got len = %d", fc.Length())
				}
			})

			t.Run("NilReceiver", func(t *testing.T) {
				var nilFC *token.FieldCollection

				// Should not panic
				nilFC.Clear()
				// Nothing to assert beyond "no panic", as Clear has no return value
			})
		})

		t.Run("Clone", func(t *testing.T) {
			t.Run("CopiesElements", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}

				fc.Add(a, b)
				clone := fc.Clone()
				// Length check
				if len(clone) != fc.Length() {
					t.Fatalf("Clone length = %d, want %d", len(clone), fc.Length())
				}
				// Content check
				for i := range fc {
					if fc[i] != clone[i] {
						t.Fatalf("Element mismatch at index %d: got %v, want %v", i, clone[i], fc[i])
					}
				}
			})

			t.Run("IsIndependentFromOriginal", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}

				fc.Add(a, b)
				clone := fc.Clone()

				// Modify original
				fc.RemoveAt(0)
				if len(clone) != 2 {
					t.Fatalf("Clone should remain unaffected after original changes, got len = %d", len(clone))
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				var fc token.FieldCollection
				clone := fc.Clone()

				if clone != nil {
					t.Fatalf("Clone of empty collection should be nil, got %#v", clone)
				}
			})

			t.Run("NilReceiver", func(t *testing.T) {
				var nilFC *token.FieldCollection
				clone := nilFC.Clone()

				if clone != nil {
					t.Fatalf("Clone of nil receiver should be nil, got %#v", clone)
				}
			})
		})

		t.Run("Contains", func(t *testing.T) {
			t.Run("TargetFound", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}

				fc.Add(a, b)

				if !fc.Contains(a) {
					t.Fatalf("Contains(%v) = false, want true", a)
				}
				if !fc.Contains(b) {
					t.Fatalf("Contains(%v) = false, want true", b)
				}
			})

			t.Run("TargetNotFound", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}
				c := token.Field{Expr: "c"}

				fc.Add(a, b)

				if fc.Contains(c) {
					t.Fatalf("Contains(%v) = true, want false", c)
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}

				if fc.Contains(a) {
					t.Fatalf("Contains(%v) = true on empty collection, want false", a)
				}
			})

			t.Run("NilReceiver", func(t *testing.T) {
				var nilFC *token.FieldCollection
				a := token.Field{Expr: "a"}

				if nilFC.Contains(a) {
					t.Fatalf("Contains(%v) = true on nil receiver, want false", a)
				}
			})
		})

		t.Run("IndexOf", func(t *testing.T) {
			t.Run("TargetFound", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}
				c := token.Field{Expr: "c"}

				fc.Add(a, b, c)

				if idx := fc.IndexOf(a); idx != 0 {
					t.Fatalf("IndexOf(%v) = %d, want 0", a, idx)
				}
				if idx := fc.IndexOf(b); idx != 1 {
					t.Fatalf("IndexOf(%v) = %d, want 1", b, idx)
				}
				if idx := fc.IndexOf(c); idx != 2 {
					t.Fatalf("IndexOf(%v) = %d, want 2", c, idx)
				}
			})

			t.Run("TargetNotFound", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}

				fc.Add(a)

				if idx := fc.IndexOf(b); idx != -1 {
					t.Fatalf("IndexOf(%v) = %d, want -1", b, idx)
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}

				if idx := fc.IndexOf(a); idx != -1 {
					t.Fatalf("IndexOf(%v) = %d on empty collection, want -1", a, idx)
				}
			})

			t.Run("NilReceiver", func(t *testing.T) {
				var nilFC *token.FieldCollection
				a := token.Field{Expr: "a"}

				if idx := nilFC.IndexOf(a); idx != -1 {
					t.Fatalf("IndexOf(%v) = %d on nil receiver, want -1", a, idx)
				}
			})
		})

		t.Run("InsertAfter", func(t *testing.T) {
			t.Run("TargetFound", func(t *testing.T) {
				var fc token.FieldCollection
				a, b, c, x := token.Field{Expr: "a"}, token.Field{Expr: "b"}, token.Field{Expr: "c"}, token.Field{Expr: "x"}
				fc.Add(a, b, c).InsertAfter(b, x)
				// expect: a, b, x, c
				if fc.Length() != 4 {
					t.Fatalf("len = %d, want 4", fc.Length())
				}
				if fc[0] != a || fc[1] != b || fc[2] != x || fc[3] != c {
					t.Fatalf("order mismatch: got %#v", fc)
				}
			})

			t.Run("TargetNotFound", func(t *testing.T) {
				var fc token.FieldCollection
				a, b, x := token.Field{Expr: "a"}, token.Field{Expr: "b"}, token.Field{Expr: "x"}
				fc.Add(a, b).InsertAfter(token.Field{Expr: "z"}, x)
				// expect: a, b, x
				if fc.IndexOf(x) != 2 {
					t.Fatalf("x should be appended; got index %d", fc.IndexOf(x))
				}
			})

			t.Run("Variadic", func(t *testing.T) {
				var fc token.FieldCollection
				a, b, c, x, y := token.Field{Expr: "a"}, token.Field{Expr: "b"}, token.Field{Expr: "c"}, token.Field{Expr: "x"}, token.Field{Expr: "y"}
				fc.Add(a, c).InsertAfter(a, x, y).InsertAfter(y, b)
				// after first: a, x, y, c
				// after y: a, x, y, b, c
				want := []token.Field{a, x, y, b, c}
				if fc.Length() != len(want) {
					t.Fatalf("len = %d, want %d", fc.Length(), len(want))
				}
				for i := range want {
					if fc[i] != want[i] {
						t.Fatalf("pos %d mismatch: got %v, want %v", i, fc[i], want[i])
					}
				}
			})

			t.Run("NilReceiver", func(t *testing.T) {
				var fc *token.FieldCollection
				if got := fc.InsertAfter(token.Field{Expr: "a"}, token.Field{Expr: "b"}); got != nil {
					t.Fatal("InsertAfter on nil receiver should return nil")
				}
			})
		})

		t.Run("InsertAt", func(t *testing.T) {
			t.Run("ValidIndex", func(t *testing.T) {
				var f token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}
				c := token.Field{Expr: "c"}

				got := f.Add(a, c).InsertAt(1, b)
				// Chainability check
				if got != &f {
					t.Fatal("InsertAt should return the same receiver for chaining")
				}
				// Length check
				if f.Length() != 3 {
					t.Fatalf("Length = %d, want 3", f.Length())
				}
				// Order check
				want := []token.Field{a, b, c}
				for i := range want {
					if f[i] != want[i] {
						t.Fatalf("Order mismatch at index %d: got %v, want %v", i, f[i], want[i])
					}
				}
			})

			t.Run("NegativeIndex", func(t *testing.T) {
				var f token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}

				got := f.Add(a).InsertAt(-5, b)
				// Chainability check
				if got != &f {
					t.Fatal("InsertAt should return the same receiver for chaining")
				}
				// Length check
				if f.Length() != 2 {
					t.Fatalf("Length = %d, want 2", f.Length())
				}
				// Order check
				want := []token.Field{b, a}
				for i := range want {
					if f[i] != want[i] {
						t.Fatalf("Order mismatch at index %d: got %v, want %v", i, f[i], want[i])
					}
				}
			})

			t.Run("IndexGreaterThanLength", func(t *testing.T) {
				var f token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}

				got := f.Add(a).InsertAt(99, b)
				// Chainability check
				if got != &f {
					t.Fatal("InsertAt should return the same receiver for chaining")
				}
				// Length check
				if f.Length() != 2 {
					t.Fatalf("Length = %d, want 2", f.Length())
				}
				// Order check
				want := []token.Field{a, b}
				for i := range want {
					if f[i] != want[i] {
						t.Fatalf("Order mismatch at index %d: got %v, want %v", i, f[i], want[i])
					}
				}
			})

			t.Run("NilReceiver", func(t *testing.T) {
				var f *token.FieldCollection
				a := token.Field{Expr: "a"}

				// Should not panic, and should still be nil
				got := f.InsertAt(0, a)
				if got != nil {
					t.Fatal("InsertAt on nil receiver should return nil")
				}
			})
		})

		t.Run("InsertBefore", func(t *testing.T) {
			t.Run("TargetFound", func(t *testing.T) {
				var fc token.FieldCollection
				a, b, c, x := token.Field{Expr: "a"}, token.Field{Expr: "b"}, token.Field{Expr: "c"}, token.Field{Expr: "x"}
				fc.Add(a, b, c).InsertBefore(b, x)
				// expect: a, x, b, c
				if fc.Length() != 4 {
					t.Fatalf("len = %d, want 4", fc.Length())
				}
				if fc[0] != a || fc[1] != x || fc[2] != b || fc[3] != c {
					t.Fatalf("order mismatch: got %#v", fc)
				}
			})

			t.Run("TargetNotFound", func(t *testing.T) {
				var fc token.FieldCollection
				a, b, x := token.Field{Expr: "a"}, token.Field{Expr: "b"}, token.Field{Expr: "x"}
				fc.Add(a, b).InsertBefore(token.Field{Expr: "z"}, x)
				// expect: x, a, b
				if fc.IndexOf(x) != 0 {
					t.Fatalf("x should be prepended; got index %d", fc.IndexOf(x))
				}
			})

			t.Run("Variadic", func(t *testing.T) {
				var fc token.FieldCollection
				a, b, c, x, y := token.Field{Expr: "a"}, token.Field{Expr: "b"}, token.Field{Expr: "c"}, token.Field{Expr: "x"}, token.Field{Expr: "y"}
				fc.Add(a, c).InsertBefore(c, x, y).InsertBefore(a, b)
				// before c: a, x, y, c
				// before a: b, a, x, y, c
				want := []token.Field{b, a, x, y, c}
				if fc.Length() != len(want) {
					t.Fatalf("len = %d, want %d", fc.Length(), len(want))
				}
				for i := range want {
					if fc[i] != want[i] {
						t.Fatalf("pos %d mismatch: got %v, want %v", i, fc[i], want[i])
					}
				}
			})
			t.Run("NilSafety", func(t *testing.T) {
				var nilFC *token.FieldCollection
				if nilFC.InsertBefore(token.Field{Expr: "a"}).Length() != 0 {
					t.Fatal("nil receiver InsertBefore should no-op")
				}
			})
		})

		t.Run("Remove", func(t *testing.T) {
			t.Run("TargetFound", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}
				c := token.Field{Expr: "c"}

				got := fc.Add(a, b, c).Remove(b)

				// Chainability check
				if got != &fc {
					t.Fatal("Remove should return the same receiver for chaining")
				}

				// Expected: a, c
				want := []token.Field{a, c}
				if fc.Length() != len(want) {
					t.Fatalf("len = %d, want %d", fc.Length(), len(want))
				}
				for i := range want {
					if fc[i] != want[i] {
						t.Fatalf("Order mismatch at index %d: got %v, want %v", i, fc[i], want[i])
					}
				}
			})

			t.Run("TargetNotFound", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}

				fc.Add(a).Remove(b) // removing missing element should no-op

				want := []token.Field{a}
				if fc.Length() != len(want) {
					t.Fatalf("len = %d, want %d", fc.Length(), len(want))
				}
				if fc[0] != a {
					t.Fatalf("Unexpected modification: got %v, want %v", fc[0], a)
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}

				fc.Remove(a) // should no-op
				if fc.Length() != 0 {
					t.Fatalf("Remove from empty collection should keep it empty, got len = %d", fc.Length())
				}
			})

			t.Run("NilReceiver", func(t *testing.T) {
				var nilFC *token.FieldCollection
				a := token.Field{Expr: "a"}

				got := nilFC.Remove(a)
				if got != nil {
					t.Fatal("Remove on nil receiver should return nil")
				}
			})
		})

		t.Run("RemoveAt", func(t *testing.T) {
			t.Run("ValidIndex", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}
				c := token.Field{Expr: "c"}

				got := fc.Add(a, b, c).RemoveAt(1)

				// Chainability check
				if got != &fc {
					t.Fatal("RemoveAt should return the same receiver for chaining")
				}

				// Expected: a, c
				want := []token.Field{a, c}
				if fc.Length() != len(want) {
					t.Fatalf("len = %d, want %d", fc.Length(), len(want))
				}
				for i := range want {
					if fc[i] != want[i] {
						t.Fatalf("Order mismatch at index %d: got %v, want %v", i, fc[i], want[i])
					}
				}
			})

			t.Run("IndexZero", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}

				fc.Add(a, b).RemoveAt(0)
				// Expected: b
				if fc.Length() != 1 || fc[0] != b {
					t.Fatalf("Expected only %v remaining, got %#v", b, fc)
				}
			})

			t.Run("IndexLastElement", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}
				b := token.Field{Expr: "b"}

				fc.Add(a, b).RemoveAt(1)
				// Expected: a
				if fc.Length() != 1 || fc[0] != a {
					t.Fatalf("Expected only %v remaining, got %#v", a, fc)
				}
			})

			t.Run("IndexOutOfRangeNegative", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}

				fc.Add(a).RemoveAt(-1) // should no-op
				if fc.Length() != 1 || fc[0] != a {
					t.Fatalf("Negative index should no-op, got %#v", fc)
				}
			})

			t.Run("IndexOutOfRangeTooLarge", func(t *testing.T) {
				var fc token.FieldCollection
				a := token.Field{Expr: "a"}

				fc.Add(a).RemoveAt(5) // should no-op
				if fc.Length() != 1 || fc[0] != a {
					t.Fatalf("Index beyond length should no-op, got %#v", fc)
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				var fc token.FieldCollection
				fc.RemoveAt(0) // should no-op
				if fc.Length() != 0 {
					t.Fatalf("RemoveAt on empty collection should leave it empty, got len = %d", fc.Length())
				}
			})

			t.Run("NilReceiver", func(t *testing.T) {
				var nilFC *token.FieldCollection
				got := nilFC.RemoveAt(0)
				if got != nil {
					t.Fatal("RemoveAt on nil receiver should return nil")
				}
			})
		})

	})
}
