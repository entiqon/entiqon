// File: common/bools.go

// Package common provides shared utility functions used across the Entiqon project.
package common

// BoolToStr returns trueStr if b is true, otherwise returns falseStr.
// This is a helper to convert boolean values into custom string representations.
//
// Example:
//
//	BoolToStr(true, "enabled", "disabled")  // returns "enabled"
//	BoolToStr(false, "yes", "no")           // returns "no"
func BoolToStr(b bool, trueStr, falseStr string) string {
	if b {
		return trueStr
	}
	return falseStr
}
