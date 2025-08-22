package extension_test

import (
	"fmt"
	"time"

	"github.com/entiqon/entiqon/common/extension"
)

func ExampleBoolean() {
	// Valid input
	v := extension.Boolean("yes")
	fmt.Println(v)

	// Invalid input → falls back to default `false`
	v2 := extension.Boolean("maybe")
	fmt.Println(v2)

	// Output:
	// true
	// false
}

func ExampleBooleanOr() {
	// Valid input
	v := extension.BooleanOr("yes", false)
	fmt.Println(v)

	// Invalid input → caller-supplied default (true) is returned
	v2 := extension.BooleanOr("maybe", true)
	fmt.Println(v2)

	// Output:
	// true
	// true
}

func ExampleNumber() {
	// Valid input (string integer)
	v := extension.Number("42")
	fmt.Println(v)

	// Valid input (negative string integer)
	v2 := extension.Number("-7")
	fmt.Println(v2)

	// Invalid input (non-integer string) → falls back to default 0
	v3 := extension.Number("oops")
	fmt.Println(v3)

	// Output:
	// 42
	// -7
	// 0
}

func ExampleNumberOr() {
	// Valid input
	v := extension.NumberOr("100", 99)
	fmt.Println(v)

	// Invalid input → caller-supplied default 99 is returned
	v2 := extension.NumberOr("oops", 99)
	fmt.Println(v2)

	// Output:
	// 100
	// 99
}

func ExampleFloat() {
	// Valid input
	v := extension.Float("2.718")
	fmt.Printf("%.2f\n", v)

	// Invalid input → falls back to default 0.0
	v2 := extension.Float("oops")
	fmt.Println(v2)

	// Output:
	// 2.72
	// 0
}

func ExampleFloatOr() {
	// Valid input
	v := extension.FloatOr("2.718", 1.23)
	fmt.Printf("%.2f\n", v)

	// Invalid input → caller-supplied default 1.23 is returned
	v2 := extension.FloatOr("oops", 1.23)
	fmt.Println(v2)

	// Output:
	// 2.72
	// 1.23
}

func ExampleDecimal() {
	// Valid input
	v := extension.Decimal("123.45", 2)
	fmt.Println(v)

	// Invalid input → falls back to default "0"
	v2 := extension.Decimal("oops", 2)
	fmt.Println(v2)

	// Output:
	// 123.45
	// 0
}

func ExampleDecimalOr() {
	// Valid input
	v := extension.DecimalOr("123.45", 2, 99.9)
	fmt.Println(v)

	// Invalid input → caller-supplied default "99.9" is returned
	v2 := extension.DecimalOr("oops", 2, 99.9)
	fmt.Println(v2)

	// Output:
	// 123.45
	// 99.9
}

func ExampleDate() {
	// Valid input
	v := extension.Date("2025-08-21")
	fmt.Println(v.Format("2006-01-02"))

	// Invalid input → falls back to default zero time
	v2 := extension.Date("oops")
	fmt.Println(v2.IsZero())

	// Output:
	// 2025-08-21
	// true
}

func ExampleDateOr() {
	def := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	// Valid input
	v := extension.DateOr("2025-08-21", def)
	fmt.Println(v.Format("2006-01-02"))

	// Invalid input → caller-supplied default date is returned
	v2 := extension.DateOr("oops", def)
	fmt.Println(v2.Format("2006-01-02"))

	// Output:
	// 2025-08-21
	// 2000-01-01
}
