package boolean_test

import (
	"fmt"

	"github.com/entiqon/entiqon/common/extension/boolean"
)

func ExampleParseFrom() {
	b, _ := boolean.ParseFrom("yes")
	fmt.Println(b)
	// Output: true
}

func ExampleBoolToString() {
	fmt.Println(boolean.BoolToString(true, "yes", "no"))
	fmt.Println(boolean.BoolToString(false, "enabled", "disabled"))
	// Output:
	// yes
	// disabled
}

func ExampleBoolToStr_deprecated() {
	fmt.Println(boolean.BoolToStr(true, "Y", "N"))
	// Output: Y
}
