package extension

import (
	"time"

	"github.com/entiqon/entiqon/common/extension/boolean"
	"github.com/entiqon/entiqon/common/extension/date"
	"github.com/entiqon/entiqon/common/extension/decimal"
	"github.com/entiqon/entiqon/common/extension/float"
	"github.com/entiqon/entiqon/common/extension/integer"
	"github.com/entiqon/entiqon/common/extension/number"
)

// Boolean parses a value into a bool.
// Returns false if parsing fails.
func Boolean(value any) bool {
	v, err := boolean.ParseFrom(value)
	if err != nil {
		return false
	}
	return v
}

// BooleanOr parses a value into a bool.
// Returns the provided default if parsing fails.
func BooleanOr(value any, def bool) bool {
	v, err := boolean.ParseFrom(value)
	if err != nil {
		return def
	}
	return v
}

// Date parses a value into a time.Time.
// Returns zero time if parsing fails.
func Date(value any) time.Time {
	v, err := date.ParseFrom(value)
	if err != nil {
		return time.Time{}
	}
	return v
}

// DateOr parses a value into a time.Time.
// Returns the provided default if parsing fails.
func DateOr(value any, def time.Time) time.Time {
	v, err := date.ParseFrom(value)
	if err != nil {
		return def
	}
	return v
}

// Decimal parses a value into a float64 with precision.
// Precision must be provided explicitly. Returns 0 if parsing fails.
func Decimal(value any, precision int) float64 {
	v, err := decimal.ParseFrom(value, precision)
	if err != nil {
		return 0
	}
	return v
}

// DecimalOr parses a value into a float64 with precision.
// Returns the provided default if parsing fails.
func DecimalOr(value any, precision int, def float64) float64 {
	v, err := decimal.ParseFrom(value, precision)
	if err != nil {
		return def
	}
	return v
}

// Float parses a value into a float64.
// Returns 0 if parsing fails.
func Float(value any) float64 {
	v, err := float.ParseFrom(value)
	if err != nil {
		return 0
	}
	return v
}

// FloatOr parses a value into a float64.
// Returns the provided default if parsing fails.
func FloatOr(value any, def float64) float64 {
	v, err := float.ParseFrom(value)
	if err != nil {
		return def
	}
	return v
}

// Integer parses a value into an int.
// Returns 0 if parsing fails.
func Integer(value any) int {
	v, err := integer.ParseFrom(value)
	if err != nil {
		return 0
	}
	return v
}

// IntegerOr parses a value into an int.
// Returns the provided default if parsing fails.
func IntegerOr(value any, def int) int {
	v, err := integer.ParseFrom(value)
	if err != nil {
		return def
	}
	return v
}

// Number parses a value into an int.
// Returns 0 if parsing fails.
func Number(value any) int {
	v, err := number.ParseFrom(value, false)
	if err != nil {
		return 0
	}
	return v
}

// NumberOr parses a value into an int.
// Returns the provided default if parsing fails.
func NumberOr(value any, def int) int {
	v, err := number.ParseFrom(value, false)
	if err != nil {
		return def
	}
	return v
}
