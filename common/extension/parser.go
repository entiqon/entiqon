// File: common/extension/parser.go

// Package extension provides convenience wrappers for parsing values into
// common types. These are thin shortcuts around the subpackages:
//
//	Boolean(v) → boolean.ParseFrom(v)
//	Float(v)   → float.ParseFrom(v)
//	Decimal(v, precision) → decimal.ParseFrom(v, precision)
//	Date(v)    → date.ParseFrom(v)
//
// Import subpackages directly if you need advanced features.
package extension

import (
	"time"

	"github.com/entiqon/entiqon/common/extension/boolean"
	"github.com/entiqon/entiqon/common/extension/date"
	"github.com/entiqon/entiqon/common/extension/decimal"
	"github.com/entiqon/entiqon/common/extension/float"
)

func Boolean(value any) (bool, error) {
	return boolean.ParseFrom(value)
}

func Float(value any) (float64, error) {
	return float.ParseFrom(value)
}

//func Integer(value any) (int, error) {
//	return number.ParseFrom(value)
//}

func Decimal(value any, precision int) (float64, error) {
	return decimal.ParseFrom(value, precision)
}

func Date(value any) (time.Time, error) {
	return date.ParseFrom(value)
}

//func String(value any) (string, error) {
//	return string.ParseFrom(value)
//}
