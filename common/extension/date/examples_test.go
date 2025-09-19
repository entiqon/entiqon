package date_test

import (
	"fmt"
	"time"

	"github.com/entiqon/common/extension/date"
)

func ExampleParseFrom() {
	t, _ := date.ParseFrom("2025-08-21")
	fmt.Println(t.Format("2006-01-02"))
	// Output: 2025-08-21
}

func ExampleParseAndFormat() {
	// Using a supported format (YYYY/MM/DD)
	out := date.ParseAndFormat("2025/08/21", "2006-01-02")
	fmt.Println(out)
	// Output: 2025-08-21
}

func ExampleCleanAndParse_default() {
	t, _ := date.CleanAndParse("2025/08/21", nil)
	fmt.Println(t.Format("2006-01-02"))
	// Output: 2025-08-21
}

func ExampleCleanAndParse_epochSeconds() {
	// CleanAndParse accepts epoch timestamps (seconds or milliseconds).
	t, _ := date.CleanAndParse("1700000000", nil)
	fmt.Println(t.Year())
	// Output: 2023
}

func ExampleCleanAndParseAsString() {
	out := date.CleanAndParseAsString("2023-11-14T22:13:20Z", "20060102")
	fmt.Println(out)
	// Output: 20231114
}

func ExampleStrictYYYYMMDDOptions() {
	opts := date.StrictYYYYMMDDOptions()
	t, _ := date.CleanAndParse("20240229", opts)
	fmt.Println(t.Format(time.RFC3339))
	// Output: 2024-02-29T00:00:00Z
}
