// File: db/internal/core/token/condition_ops.go

package token

import "strings"

// IsValid returns true if the condition has no validation or parsing errors.
//
// Since: v0.0.1
// Updated: v1.4.0
func (c Condition) IsValid() bool {
	return strings.TrimSpace(c.Key) != "" && c.Error == nil
}
