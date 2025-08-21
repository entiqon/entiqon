// Package collection provides a generic, type-safe container for values.
//
// It offers common operations such as adding, inserting, removing, searching,
// cloning, and transforming items. Collections support method chaining
// for imperative-style code, and functional-style helpers like Filter and Map.
//
// Quick Reference:
//
//	Constructors:
//	  - New       → Create a new empty collection
//	  - FromSlice → Initialize from an existing slice
//
//	Mutators:
//	  - Add       → Append values
//	  - InsertAt  → Insert values at index
//	  - Remove    → Remove all occurrences of value
//	  - RemoveAt  → Remove value by index
//	  - Clear     → Remove all elements
//
//	Queries:
//	  - Contains  → Check if value exists
//	  - IndexOf   → Find index of first occurrence
//	  - At        → Safe access by index
//	  - Length    → Number of elements
//	  - Items     → Copy of slice
//
//	Functional:
//	  - ForEach   → Apply function to each element
//	  - Filter    → New collection matching predicate
//	  - Map       → Transform collection into another type
//
//	Utilities:
//	  - Clone     → Deep copy of collection
//
// Example usage:
//
//	nums := collection.New[int]().Add(1, 2, 3).Remove(2)
//	fmt.Println(nums.Items()) // [1 3]
//
//	evens := collection.Filter(nums, func(x int) bool { return x%2 == 0 })
//	fmt.Println(evens.Items()) // [2]
//
//	strs := collection.Map(nums, func(x int) string { return fmt.Sprintf("#%d", x) })
//	fmt.Println(strs.Items()) // ["#1" "#3"]
package collection
