/**
 * @Author: Isidro Lopez isidro.lopezg@live.com
 * @Date: 2025-08-21 22:27:49
 * @LastEditors: Isidro Lopez isidro.lopezg@live.com
 * @LastEditTime: 2025-08-21 22:40:48
 * @FilePath: common/extension/boolean/to_string.go
 * @Description: 这是默认设置,可以在设置》工具》File Description中进行配置
 */
// File: common/to_string.go

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

// BoolToStr is an alias of BoolToString.
//
// Deprecated: use BoolToString instead.
func BoolToStr(b bool, trueStr, falseStr string) string {
	return BoolToString(b, trueStr, falseStr)
}
