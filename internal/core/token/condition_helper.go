package token

import (
	"strings"
	"time"
)

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
	supportedOps := []string{
		"NOT IN", "IN", "BETWEEN", "<>", "!=", ">=", "<=", "LIKE", "=", ">", "<",
	}
	input = strings.TrimSpace(input)
	upper := strings.ToUpper(input)

	for _, op := range supportedOps {
		if idx := strings.Index(upper, op); idx > 0 {
			field = strings.TrimSpace(input[:idx])
			value = strings.TrimSpace(input[idx+len(op):])
			return field, op, value, true
		}
	}
	return "", "", "", false
}
