package boolean

// BoolToString returns trueStr if b is true, otherwise returns falseStr.
// This helper converts a boolean value into custom string representations.
//
// Example:
//
//	BoolToString(true, "enabled", "disabled")  // returns "enabled"
//	BoolToString(false, "yes", "no")           // returns "no"
func BoolToString(b bool, trueStr, falseStr string) string {
	if b {
		return trueStr
	}
	return falseStr
}
