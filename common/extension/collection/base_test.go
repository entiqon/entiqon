package collection_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/entiqon/entiqon/common/extension/collection"
)

func TestCollection(t *testing.T) {
	t.Run("Methods", func(t *testing.T) {
		// TestAdd verifies the Add method of Collection[T] across different cases.
		// It ensures values are appended correctly, multiple values are handled,
		// and that adding to an empty collection behaves as expected.
		t.Run("Add", func(t *testing.T) {
			t.Run("EmptyCollection", func(t *testing.T) {
				// Adding to an empty collection should create the first element.
				c := collection.New[int]().Add(42)
				got := c.Items()
				want := []int{42}
				if len(got) != len(want) {
					t.Fatalf("expected length %d, got %d", len(want), len(got))
				}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("SingleValue", func(t *testing.T) {
				// Adding a single value appends to the end.
				c := collection.FromSlice([]int{1, 2}).Add(3)
				got := c.Items()
				want := []int{1, 2, 3}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("MultipleValues", func(t *testing.T) {
				// Adding multiple values appends them in order.
				c := collection.FromSlice([]int{1}).Add(2, 3, 4)
				got := c.Items()
				want := []int{1, 2, 3, 4}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("NoValues", func(t *testing.T) {
				// Adding no values should leave collection unchanged.
				c := collection.FromSlice([]int{1, 2, 3}).Add()
				got := c.Items()
				want := []int{1, 2, 3}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("WithStructs", func(t *testing.T) {
				// Adding struct values demonstrates generic type support.
				type Item struct{ Name string }
				c := collection.New[Item]().
					Add(Item{"apple"}).
					Add(Item{"banana"})

				got := c.Items()
				want := []Item{{"apple"}, {"banana"}}
				if len(got) != len(want) {
					t.Fatalf("expected length %d, got %d", len(want), len(got))
				}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})
		})

		// TestAt verifies the At method of Collection[T] across valid and invalid indices.
		// It ensures correct values are returned when indices are valid and safe handling
		// (zero value + false) when indices are invalid.
		t.Run("At", func(t *testing.T) {
			t.Run("ValidIndex", func(t *testing.T) {
				c := collection.FromSlice([]int{10, 20, 30})
				val, ok := c.At(1)
				if !ok {
					t.Fatalf("expected ok=true, got false")
				}
				if val != 20 {
					t.Errorf("expected 20, got %v", val)
				}
			})

			t.Run("FirstIndex", func(t *testing.T) {
				c := collection.FromSlice([]int{99, 100})
				val, ok := c.At(0)
				if !ok {
					t.Fatalf("expected ok=true, got false")
				}
				if val != 99 {
					t.Errorf("expected 99, got %v", val)
				}
			})

			t.Run("LastIndex", func(t *testing.T) {
				c := collection.FromSlice([]string{"a", "b", "c"})
				val, ok := c.At(2)
				if !ok {
					t.Fatalf("expected ok=true, got false")
				}
				if val != "c" {
					t.Errorf("expected 'c', got %v", val)
				}
			})

			t.Run("NegativeIndex", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				val, ok := c.At(-1)
				if ok {
					t.Errorf("expected ok=false, got true with val=%v", val)
				}
				if val != 0 {
					t.Errorf("expected zero value (0), got %v", val)
				}
			})

			t.Run("OutOfRangeIndex", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				val, ok := c.At(10)
				if ok {
					t.Errorf("expected ok=false, got true with val=%v", val)
				}
				if val != 0 {
					t.Errorf("expected zero value (0), got %v", val)
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				c := collection.New[string]()
				val, ok := c.At(0)
				if ok {
					t.Errorf("expected ok=false, got true with val=%v", val)
				}
				if val != "" {
					t.Errorf("expected zero value (empty string), got %v", val)
				}
			})
		})

		// TestClear verifies the Clear method of Collection[T].
		// It ensures that Clear removes all elements regardless of initial state
		// and that subsequent operations on the collection remain valid.
		t.Run("Clear", func(t *testing.T) {
			t.Run("NonEmptyCollection", func(t *testing.T) {
				// Clear should remove all elements.
				c := collection.FromSlice([]int{1, 2, 3})
				c.Clear()
				if c.Length() != 0 {
					t.Errorf("expected length 0 after Clear, got %d", c.Length())
				}
				if len(c.Items()) != 0 {
					t.Errorf("expected Items() to be empty slice, got %v", c.Items())
				}
			})

			t.Run("AlreadyEmptyCollection", func(t *testing.T) {
				// Clearing an empty collection should be safe and idempotent.
				c := collection.New[string]()
				c.Clear()
				if c.Length() != 0 {
					t.Errorf("expected length 0, got %d", c.Length())
				}
			})

			t.Run("AfterClearAddNewValues", func(t *testing.T) {
				// Collection should remain usable after Clear.
				c := collection.FromSlice([]int{1, 2, 3})
				c.Clear()
				c.Add(10, 20)
				want := []int{10, 20}
				got := c.Items()
				if len(got) != len(want) {
					t.Fatalf("expected length %d, got %d", len(want), len(got))
				}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})
		})

		// TestClone verifies the Clone method of Collection[T].
		// It ensures that a new independent collection is created,
		// with the same elements but a distinct underlying slice.
		t.Run("Clone", func(t *testing.T) {
			t.Run("NonEmptyCollection", func(t *testing.T) {
				// Cloning should preserve values.
				c1 := collection.FromSlice([]int{1, 2, 3})
				c2 := c1.Clone()

				if c1.Length() != c2.Length() {
					t.Errorf("expected same length, got %d and %d", c1.Length(), c2.Length())
				}
				for i, v := range c1.Items() {
					if c2.Items()[i] != v {
						t.Errorf("expected %v at index %d, got %v", v, i, c2.Items()[i])
					}
				}
			})

			t.Run("Independence", func(t *testing.T) {
				// Modifications to the clone should not affect the original.
				c1 := collection.FromSlice([]int{1, 2, 3})
				c2 := c1.Clone().Add(4)

				if c1.Length() == c2.Length() {
					t.Errorf("expected different lengths, got both %d", c1.Length())
				}
				if c1.Contains(4) {
					t.Errorf("expected original not to contain 4")
				}
				if !c2.Contains(4) {
					t.Errorf("expected clone to contain 4")
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				// Cloning an empty collection should still return a distinct instance.
				c1 := collection.New[string]()
				c2 := c1.Clone()

				if c1.Length() != 0 || c2.Length() != 0 {
					t.Errorf("expected both collections empty, got %d and %d", c1.Length(), c2.Length())
				}
				if c1 == c2 {
					t.Errorf("expected distinct instances, got same reference")
				}
			})
		})

		// TestContains verifies the Contains method of Collection[T].
		// It ensures correct detection of existing and non-existing values
		// across different element types and edge cases.
		t.Run("Contains", func(t *testing.T) {
			t.Run("ValuePresent", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				if !c.Contains(2) {
					t.Errorf("expected collection to contain 2")
				}
			})

			t.Run("ValueAbsent", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				if c.Contains(99) {
					t.Errorf("did not expect collection to contain 99")
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				c := collection.New[string]()
				if c.Contains("x") {
					t.Errorf("did not expect empty collection to contain any value")
				}
			})

			t.Run("FirstElement", func(t *testing.T) {
				c := collection.FromSlice([]string{"a", "b", "c"})
				if !c.Contains("a") {
					t.Errorf("expected collection to contain 'a'")
				}
			})

			t.Run("LastElement", func(t *testing.T) {
				c := collection.FromSlice([]string{"a", "b", "c"})
				if !c.Contains("c") {
					t.Errorf("expected collection to contain 'c'")
				}
			})

			t.Run("Structs", func(t *testing.T) {
				type Item struct{ ID int }
				i1 := Item{ID: 1}
				i2 := Item{ID: 2}
				c := collection.FromSlice([]Item{i1, i2})

				if !c.Contains(i2) {
					t.Errorf("expected collection to contain %+v", i2)
				}
				if c.Contains(Item{ID: 3}) {
					t.Errorf("did not expect collection to contain {ID:3}")
				}
			})
		})

		// TestFilter verifies the Filter function of Collection[T].
		// Filter creates a new collection containing only the elements
		// for which the provided function returns true.
		//
		// Developer/User Guide Notes:
		//   - Filter never mutates the original collection; it always returns a new one.
		//   - The predicate function is applied to every element in order.
		//   - The result may be empty if no element satisfies the predicate.
		//   - Example:
		//
		//         nums := collection.FromSlice([]int{1, 2, 3, 4})
		//         evens := collection.Filter(nums, func(x int) bool { return x%2 == 0 })
		//         fmt.Println(evens.Items()) // [2 4]
		t.Run("Filter", func(t *testing.T) {
			t.Run("SelectEvens", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3, 4, 5, 6})
				evens := collection.Filter(c, func(x int) bool { return x%2 == 0 })

				want := []int{2, 4, 6}
				got := evens.Items()
				if len(got) != len(want) {
					t.Fatalf("expected %d items, got %d", len(want), len(got))
				}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("FilterNone", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 3, 5})
				result := collection.Filter(c, func(x int) bool { return x%2 == 0 })

				if result.Length() != 0 {
					t.Errorf("expected empty result, got %v", result.Items())
				}
			})

			t.Run("FilterAll", func(t *testing.T) {
				c := collection.FromSlice([]string{"a", "b", "c"})
				result := collection.Filter(c, func(s string) bool { return true })

				if result.Length() != 3 {
					t.Errorf("expected full collection, got length %d", result.Length())
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				c := collection.New[int]()
				result := collection.Filter(c, func(x int) bool { return x > 0 })

				if result.Length() != 0 {
					t.Errorf("expected empty collection, got %d", result.Length())
				}
			})

			t.Run("Structs", func(t *testing.T) {
				type User struct {
					Name string
					Age  int
				}
				c := collection.FromSlice([]User{
					{"Alice", 30},
					{"Bob", 20},
					{"Carol", 40},
				})

				adults := collection.Filter(c, func(u User) bool { return u.Age >= 30 })
				got := adults.Items()

				if len(got) != 2 || got[0].Name != "Alice" || got[1].Name != "Carol" {
					t.Errorf("expected [Alice, Carol], got %v", got)
				}
			})
		})

		// TestForEach verifies the ForEach method of Collection[T].
		// It ensures the provided function is invoked for every element in order
		// and that side effects are applied consistently across data types.
		t.Run("ForEach", func(t *testing.T) {
			t.Run("IntsAccumulation", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3, 4})
				sum := 0
				c.ForEach(func(v int) {
					sum += v
				})
				if sum != 10 {
					t.Errorf("expected sum=10, got %d", sum)
				}
			})

			t.Run("StringsConcatenation", func(t *testing.T) {
				c := collection.FromSlice([]string{"a", "b", "c"})
				var sb strings.Builder
				c.ForEach(func(s string) {
					sb.WriteString(s)
				})
				got := sb.String()
				want := "abc"
				if got != want {
					t.Errorf("expected %q, got %q", want, got)
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				c := collection.New[int]()
				called := false
				c.ForEach(func(v int) {
					called = true
				})
				if called {
					t.Errorf("did not expect function to be called on empty collection")
				}
			})

			t.Run("Structs", func(t *testing.T) {
				type User struct{ Name string }
				c := collection.FromSlice([]User{{"Alice"}, {"Bob"}})

				names := make([]string, 0)
				c.ForEach(func(u User) {
					names = append(names, u.Name)
				})

				if len(names) != 2 || names[0] != "Alice" || names[1] != "Bob" {
					t.Errorf("expected [Alice Bob], got %v", names)
				}
			})
		})

		// TestIndexOf verifies the IndexOf method of Collection[T].
		// It ensures the correct index of the first matching value is returned,
		// and -1 is returned when the value is not found.
		t.Run("IndexOf", func(t *testing.T) {
			t.Run("ValuePresent", func(t *testing.T) {
				c := collection.FromSlice([]int{10, 20, 30})
				idx := c.IndexOf(20)
				if idx != 1 {
					t.Errorf("expected index 1, got %d", idx)
				}
			})

			t.Run("ValueAbsent", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				idx := c.IndexOf(99)
				if idx != -1 {
					t.Errorf("expected -1 for missing value, got %d", idx)
				}
			})

			t.Run("FirstOccurrence", func(t *testing.T) {
				c := collection.FromSlice([]int{5, 10, 5, 10})
				idx := c.IndexOf(5)
				if idx != 0 {
					t.Errorf("expected first occurrence at index 0, got %d", idx)
				}
			})

			t.Run("LastElement", func(t *testing.T) {
				c := collection.FromSlice([]string{"a", "b", "c"})
				idx := c.IndexOf("c")
				if idx != 2 {
					t.Errorf("expected index 2, got %d", idx)
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				c := collection.New[string]()
				idx := c.IndexOf("x")
				if idx != -1 {
					t.Errorf("expected -1 for empty collection, got %d", idx)
				}
			})

			t.Run("Structs", func(t *testing.T) {
				type Item struct{ ID int }
				i1 := Item{ID: 1}
				i2 := Item{ID: 2}
				c := collection.FromSlice([]Item{i1, i2})

				if idx := c.IndexOf(i2); idx != 1 {
					t.Errorf("expected index 1, got %d", idx)
				}
				if idx := c.IndexOf(Item{ID: 3}); idx != -1 {
					t.Errorf("expected -1 for non-existing struct, got %d", idx)
				}
			})
		})

		// TestInsertAt verifies the InsertAt method behavior across different positions
		// and validates edge-case handling when indices are out of range or negative.
		t.Run("InsertAt", func(t *testing.T) {
			t.Run("Start", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3}).InsertAt(0, 0)
				want := []int{0, 1, 2, 3}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("Middle", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3}).InsertAt(1, 99)
				want := []int{1, 99, 2, 3}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("End", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3}).InsertAt(3, 42)
				want := []int{1, 2, 3, 42}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("OutOfRange", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3}).InsertAt(10, 77)
				want := []int{1, 2, 3, 77}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("Negative", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3}).InsertAt(-5, 11)
				want := []int{11, 1, 2, 3}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})
		})

		// TestItems verifies the Items method of Collection[T].
		// It ensures a shallow copy of the underlying slice is returned,
		// modifications to the result do not affect the original collection,
		// and values are preserved in correct order.
		t.Run("Items", func(t *testing.T) {
			t.Run("NonEmptyCollection", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				got := c.Items()
				want := []int{1, 2, 3}

				if len(got) != len(want) {
					t.Fatalf("expected length %d, got %d", len(want), len(got))
				}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v at %d, got %v", want[i], i, got[i])
					}
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				c := collection.New[string]()
				got := c.Items()
				if len(got) != 0 {
					t.Errorf("expected empty slice, got %v", got)
				}
			})

			t.Run("IndependenceOfCopy", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				items := c.Items()
				items[0] = 99 // mutate the returned slice

				// original collection should remain unchanged
				if c.Items()[0] != 1 {
					t.Errorf("expected original first element = 1, got %d", c.Items()[0])
				}
			})

			t.Run("Structs", func(t *testing.T) {
				type User struct{ Name string }
				c := collection.FromSlice([]User{{"Alice"}, {"Bob"}})

				got := c.Items()
				if len(got) != 2 || got[0].Name != "Alice" || got[1].Name != "Bob" {
					t.Errorf("expected [Alice Bob], got %v", got)
				}
			})
		})

		// TestLength verifies the Length method of Collection[T].
		// Length returns the number of elements currently stored in the collection.
		// It is equivalent to len(c.Items()) but avoids allocating a copy.
		//
		// User Guide Note:
		//   - Use Length() to check how many items are in a collection.
		//   - Length() is O(1) and safe on empty collections.
		//   - Example:
		//
		//         c := collection.FromSlice([]int{1, 2, 3})
		//         fmt.Println(c.Length()) // prints 3
		t.Run("Length", func(t *testing.T) {
			t.Run("EmptyCollection", func(t *testing.T) {
				// Length of a new collection must be 0.
				c := collection.New[int]()
				if c.Length() != 0 {
					t.Errorf("expected length 0, got %d", c.Length())
				}
			})

			t.Run("NonEmptyCollection", func(t *testing.T) {
				// Length should match number of elements in slice.
				c := collection.FromSlice([]string{"a", "b", "c"})
				if c.Length() != 3 {
					t.Errorf("expected length 3, got %d", c.Length())
				}
			})

			t.Run("AfterAdd", func(t *testing.T) {
				// Adding elements increases Length accordingly.
				c := collection.New[int]()
				c.Add(1, 2, 3)
				if c.Length() != 3 {
					t.Errorf("expected length 3 after Add, got %d", c.Length())
				}
			})

			t.Run("AfterRemove", func(t *testing.T) {
				// Removing elements decreases Length.
				c := collection.FromSlice([]int{1, 2, 2, 3})
				c.Remove(2)
				if c.Length() != 2 {
					t.Errorf("expected length 2 after Remove, got %d", c.Length())
				}
			})

			t.Run("AfterClear", func(t *testing.T) {
				// Clear resets Length to 0.
				c := collection.FromSlice([]int{1, 2, 3})
				c.Clear()
				if c.Length() != 0 {
					t.Errorf("expected length 0 after Clear, got %d", c.Length())
				}
			})
		})

		// TestMap verifies the Map function of Collection[T].
		// Map transforms each element into a new collection of type R
		// by applying the given function to every element.
		//
		// Developer/User Guide Notes:
		//   - Map always creates a new collection; the original is unchanged.
		//   - The transformation function is applied in order.
		//   - Output type can differ from input type.
		//   - Example:
		//
		//         nums := collection.FromSlice([]int{1, 2, 3})
		//         strs := collection.Map(nums, func(x int) string { return strconv.Itoa(x) })
		//         fmt.Println(strs.Items()) // ["1" "2" "3"]
		t.Run("Map", func(t *testing.T) {
			t.Run("IntToString", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				strs := collection.Map(c, func(x int) string { return strconv.Itoa(x) })

				want := []string{"1", "2", "3"}
				got := strs.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %q, got %q", want[i], got[i])
					}
				}
			})

			t.Run("SquareNumbers", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3, 4})
				squares := collection.Map(c, func(x int) int { return x * x })

				want := []int{1, 4, 9, 16}
				got := squares.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("StringsToLengths", func(t *testing.T) {
				c := collection.FromSlice([]string{"go", "lang"})
				lengths := collection.Map(c, func(s string) int { return len(s) })

				want := []int{2, 4}
				got := lengths.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				c := collection.New[int]()
				result := collection.Map(c, func(x int) int { return x * 2 })

				if result.Length() != 0 {
					t.Errorf("expected empty result, got %d", result.Length())
				}
			})

			t.Run("StructsToField", func(t *testing.T) {
				type User struct {
					Name string
					Age  int
				}
				c := collection.FromSlice([]User{
					{"Alice", 30},
					{"Bob", 20},
				})

				names := collection.Map(c, func(u User) string { return u.Name })
				got := names.Items()

				want := []string{"Alice", "Bob"}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %q, got %q", want[i], got[i])
					}
				}
			})
		})

		// TestRemove verifies the Remove method of Collection[T].
		// Remove deletes all occurrences of a given value from the collection.
		//
		// Developer/User Guide Notes:
		//   - Remove(value) traverses the collection and removes *all* matching items.
		//   - If the value does not exist, the collection is left unchanged.
		//   - Remove is safe on empty collections (no panic, no-op).
		//   - Returns the collection itself, so you can chain calls:
		//
		//         c := collection.FromSlice([]int{1, 2, 2, 3})
		//         c.Remove(2).Add(4)
		//         fmt.Println(c.Items()) // [1 3 4]
		t.Run("Remove", func(t *testing.T) {
			t.Run("SingleOccurrence", func(t *testing.T) {
				// Removing a value that appears once should delete it.
				c := collection.FromSlice([]int{1, 2, 3})
				c.Remove(2)
				want := []int{1, 3}
				got := c.Items()
				if len(got) != len(want) {
					t.Fatalf("expected length %d, got %d", len(want), len(got))
				}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("MultipleOccurrences", func(t *testing.T) {
				// All matching values should be removed.
				c := collection.FromSlice([]int{1, 2, 2, 3})
				c.Remove(2)
				want := []int{1, 3}
				got := c.Items()
				if len(got) != len(want) {
					t.Fatalf("expected length %d, got %d", len(want), len(got))
				}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("ValueAbsent", func(t *testing.T) {
				// Removing a non-existing value should leave collection unchanged.
				c := collection.FromSlice([]string{"a", "b", "c"})
				c.Remove("z")
				want := []string{"a", "b", "c"}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("FirstElement", func(t *testing.T) {
				// Removing first element should shift others down.
				c := collection.FromSlice([]int{10, 20, 30})
				c.Remove(10)
				want := []int{20, 30}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("LastElement", func(t *testing.T) {
				// Removing last element should shorten collection.
				c := collection.FromSlice([]int{10, 20, 30})
				c.Remove(30)
				want := []int{10, 20}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v, got %v", want, got)
					}
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				// Safe no-op on empty collections.
				c := collection.New[int]()
				c.Remove(42)
				if c.Length() != 0 {
					t.Errorf("expected empty collection, got %d", c.Length())
				}
			})

			t.Run("Structs", func(t *testing.T) {
				// Works with comparable struct types.
				type Item struct{ ID int }
				i1 := Item{ID: 1}
				i2 := Item{ID: 2}
				c := collection.FromSlice([]Item{i1, i2, i1})
				c.Remove(i1)

				got := c.Items()
				if len(got) != 1 || got[0] != i2 {
					t.Errorf("expected [%+v], got %v", i2, got)
				}
			})
		})

		// TestRemoveAt verifies the RemoveAt method of Collection[T].
		// RemoveAt deletes the element at a given index if valid.
		// If the index is out of range, the collection remains unchanged.
		//
		// Developer/User Guide Notes:
		//   - RemoveAt(idx) removes the element at the specified index.
		//   - If idx < 0 or idx >= Length(), nothing happens (safe no-op).
		//   - Returns the collection itself, so calls can be chained.
		//
		//         c := collection.FromSlice([]int{10, 20, 30})
		//         c.RemoveAt(1).Add(40)
		//         fmt.Println(c.Items()) // [10 30 40]
		t.Run("RemoveAt", func(t *testing.T) {
			t.Run("ValidIndex", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				c.RemoveAt(1)
				want := []int{1, 3}
				got := c.Items()
				if len(got) != len(want) {
					t.Fatalf("expected %d, got %d", len(want), len(got))
				}
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v at %d, got %v", want[i], i, got[i])
					}
				}
			})

			t.Run("FirstIndex", func(t *testing.T) {
				c := collection.FromSlice([]int{10, 20, 30})
				c.RemoveAt(0)
				want := []int{20, 30}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v at %d, got %v", want[i], i, got[i])
					}
				}
			})

			t.Run("LastIndex", func(t *testing.T) {
				c := collection.FromSlice([]string{"a", "b", "c"})
				c.RemoveAt(2)
				want := []string{"a", "b"}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v at %d, got %v", want[i], i, got[i])
					}
				}
			})

			t.Run("NegativeIndex", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				c.RemoveAt(-1) // no-op
				want := []int{1, 2, 3}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v at %d, got %v", want[i], i, got[i])
					}
				}
			})

			t.Run("OutOfRangeIndex", func(t *testing.T) {
				c := collection.FromSlice([]int{1, 2, 3})
				c.RemoveAt(10) // no-op
				want := []int{1, 2, 3}
				got := c.Items()
				for i := range want {
					if got[i] != want[i] {
						t.Errorf("expected %v at %d, got %v", want[i], i, got[i])
					}
				}
			})

			t.Run("EmptyCollection", func(t *testing.T) {
				c := collection.New[int]()
				c.RemoveAt(0) // safe no-op
				if c.Length() != 0 {
					t.Errorf("expected empty collection, got length %d", c.Length())
				}
			})
		})
	})
}
