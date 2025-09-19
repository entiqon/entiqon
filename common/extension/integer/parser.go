// File: common/extension/float/integer/parser.go

package integer

import (
	"fmt"

	"github.com/entiqon/common/extension/float"
)

// ParseFrom converts a variety of input types into an int,
// truncating any fractional part toward zero.
//
// Internally it delegates parsing to float.ParseFrom.
func ParseFrom(value interface{}) (int, error) {
	f, err := float.ParseFrom(value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse into integer: %w", err)
	}
	return int(f), nil
}
