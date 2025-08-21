package collection_test

import (
	"fmt"

	"github.com/entiqon/entiqon/common/extension/collection"
)

// ExampleNew demonstrates how to create a new collection and add items.
func ExampleNew() {
	c := collection.New[int]().Add(1, 2, 3)
	fmt.Println(c.Items())
	// Output: [1 2 3]
}

// ExampleFromSlice demonstrates creating a collection from an existing slice.
func ExampleFromSlice() {
	c := collection.FromSlice([]string{"a", "b"})
	fmt.Println(c.Items())
	// Output: [a b]
}

// ExampleFilter demonstrates filtering a collection by predicate.
func ExampleFilter() {
	nums := collection.FromSlice([]int{1, 2, 3, 4, 5})
	evens := collection.Filter(nums, func(x int) bool { return x%2 == 0 })
	fmt.Println(evens.Items())
	// Output: [2 4]
}

// ExampleMap demonstrates mapping a collection to another type.
func ExampleMap() {
	nums := collection.FromSlice([]int{1, 2, 3})
	strs := collection.Map(nums, func(x int) string { return fmt.Sprintf("#%d", x) })
	fmt.Println(strs.Items())
	// Output: [#1 #2 #3]
}

// ExampleClone demonstrates cloning a collection.
func ExampleClone() {
	c1 := collection.FromSlice([]int{1, 2, 3})
	c2 := c1.Clone().Add(4)

	fmt.Println(c1.Items())
	fmt.Println(c2.Items())
	// Output:
	// [1 2 3]
	// [1 2 3 4]
}
