// File: db/internal/core/token/condition_helper.go

package token

import (
	"regexp"
	"strings"
	"time"
)

var conditionRegex = regexp.MustCompile(`(?i)^(.+?)\s+(NOT IN|IN|BETWEEN|<>|!=|>=|<=|LIKE|=|>|<)\s+(.+)$`)

// AreCompatibleTypes checks whether all provided values belong to the same compatible type group.
//
// Since: v1.4.0
func AreCompatibleTypes(values ...any) bool {
	if len(values) < 2 {
		return false
	}

	switch values[0].(type) {
	case string:
		for _, v := range values[1:] {
			if _, ok := v.(string); !ok {
				return false
			}
		}
	case int, int64, float32, float64:
		for _, v := range values[1:] {
			switch v.(type) {
			case int, int64, float32, float64:
				// ok
			default:
				return false
			}
		}
	case time.Time:
		for _, v := range values[1:] {
			if _, ok := v.(time.Time); !ok {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func extractConditionParts(input string) (field string, operator string, value string, ok bool) {
	matches := conditionRegex.FindStringSubmatch(input)
	if len(matches) == 4 {
		return strings.TrimSpace(matches[1]), strings.ToUpper(matches[2]), strings.TrimSpace(matches[3]), true
	}
	return "", "", "", false
}
