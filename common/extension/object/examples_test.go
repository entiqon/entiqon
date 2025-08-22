package object_test

import (
	"fmt"

	"github.com/entiqon/entiqon/common/extension/object"
)

type Product struct {
	ID   int
	Name string
}

func ExampleExists() {
	m := map[string]any{"Foo": 123}
	fmt.Println(object.Exists(m, "foo"))
	fmt.Println(object.Exists(m, "bar"))
	// Output:
	// true
	// false
}

func ExampleGetValue() {
	m := map[string]any{"Foo": 123}
	val := object.GetValue[int](m, "foo", 0)
	fmt.Println(val)
	// Output: 123
}

func ExampleSetValue_map() {
	m := map[string]any{}
	m, _ = object.SetValue(m, "Key", "Value")
	fmt.Println(m["Key"])
	// Output: Value
}

func ExampleSetValue_struct() {
	p := &Product{ID: 1, Name: "Book"}
	p, _ = object.SetValue(p, "Name", "Notebook")
	fmt.Println(p.Name)
	// Output: Notebook
}
