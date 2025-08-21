// File: common/extension/collection/base.go

// Package collection provides a generic, type-safe container for values.
// It offers common operations such as adding, inserting, removing, searching,
// cloning, and transforming items. Collections support method chaining
// for imperative-style code, and functional-style helpers like Filter and Map.
//
// Example:
//
//	nums := collection.New[int]().Add(1, 2, 3).Remove(2)
//	fmt.Println(nums.Items()) // [1 3]
//
//	evens := collection.Filter(nums, func(x int) bool { return x%2 == 0 })
//	fmt.Println(evens.Items()) // [2]
//
//	squares := collection.Map(nums, func(x int) int { return x * x })
//	fmt.Println(squares.Items()) // [1 9]
package collection

// Collection is a generic container for elements of type T.
// It provides common methods for managing and querying its items.
// All operations are type-safe and preserve order.
type Collection[T comparable] struct {
	items []T
}

// New creates and returns an empty collection.
func New[T comparable]() *Collection[T] {
	return &Collection[T]{items: []T{}}
}

// FromSlice initializes a collection from the given slice.
// The slice is copied to avoid modifying the caller’s data.
func FromSlice[T comparable](src []T) *Collection[T] {
	cp := make([]T, len(src))
	copy(cp, src)
	return &Collection[T]{items: cp}
}

// Add appends one or more values to the end of the collection.
// Returns the collection itself for method chaining.
func (c *Collection[T]) Add(values ...T) *Collection[T] {
	c.items = append(c.items, values...)
	return c
}

// At returns the element at the given index and true if the index is valid.
// If the index is out of range, it returns the zero value of T and false.
func (c *Collection[T]) At(idx int) (T, bool) {
	if idx < 0 || idx >= len(c.items) {
		var zero T
		return zero, false
	}
	return c.items[idx], true
}

// Clear removes all elements from the collection.
func (c *Collection[T]) Clear() {
	c.items = []T{}
}

// Clone returns a deep copy of the collection.
// The new collection has its own backing slice.
func (c *Collection[T]) Clone() *Collection[T] {
	cp := make([]T, len(c.items))
	copy(cp, c.items)
	return &Collection[T]{items: cp}
}

// Contains reports whether the given value exists in the collection.
func (c *Collection[T]) Contains(value T) bool {
	for _, v := range c.items {
		if v == value {
			return true
		}
	}
	return false
}

// Filter creates a new collection containing only the elements
// for which the given function returns true.
func Filter[T comparable](c *Collection[T], fn func(T) bool) *Collection[T] {
	res := make([]T, 0)
	for _, v := range c.Items() {
		if fn(v) {
			res = append(res, v)
		}
	}
	return &Collection[T]{items: res}
}

// ForEach applies the given function to each element in the collection.
// It is primarily used for side effects, as it does not return a new collection.
func (c *Collection[T]) ForEach(fn func(T)) {
	for _, v := range c.items {
		fn(v)
	}
}

// IndexOf returns the index of the first occurrence of the given value,
// or -1 if the value is not present.
func (c *Collection[T]) IndexOf(value T) int {
	for i, v := range c.items {
		if v == value {
			return i
		}
	}
	return -1
}

// InsertAt inserts the given values starting at the specified index.
// If index < 0, values are inserted at the start.
// If index > len, values are appended at the end.
// Returns the collection itself for method chaining.
func (c *Collection[T]) InsertAt(idx int, values ...T) *Collection[T] {
	if idx < 0 {
		idx = 0
	}
	if idx > len(c.items) {
		idx = len(c.items)
	}
	c.items = append(c.items[:idx], append(values, c.items[idx:]...)...)
	return c
}

// Items returns a shallow copy of the underlying slice.
// Modifying the returned slice does not affect the collection.
func (c *Collection[T]) Items() []T {
	cp := make([]T, len(c.items))
	copy(cp, c.items)
	return cp
}

// Length returns the number of elements in the collection.
func (c *Collection[T]) Length() int {
	return len(c.items)
}

// Map transforms each element in the collection using the given function,
// returning a new collection of type R.
func Map[T comparable, R comparable](c *Collection[T], fn func(T) R) *Collection[R] {
	res := make([]R, 0, c.Length())
	for _, v := range c.Items() {
		res = append(res, fn(v))
	}
	return &Collection[R]{items: res}
}

// Remove deletes all occurrences of the given value from the collection.
// Returns the collection itself for method chaining.
func (c *Collection[T]) Remove(value T) *Collection[T] {
	out := c.items[:0]
	for _, v := range c.items {
		if v != value {
			out = append(out, v)
		}
	}
	c.items = out
	return c
}

// RemoveAt deletes the element at the specified index.
// If the index is invalid, the collection remains unchanged.
// Returns the collection itself for method chaining.
func (c *Collection[T]) RemoveAt(idx int) *Collection[T] {
	if _, ok := c.At(idx); !ok {
		return c
	}
	c.items = append(c.items[:idx], c.items[idx+1:]...)
	return c
}
