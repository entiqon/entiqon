package helpers

import "fmt"

func Stringify(input []any) []string {
	parts := make([]string, len(input))
	for i, v := range input {
		parts[i] = fmt.Sprint(v)
	}
	return parts
}
