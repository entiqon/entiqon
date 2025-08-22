/**
 * @Author: Isidro Lopez isidro.lopezg@live.com
 * @Date: 2025-08-21 22:27:49
 * @LastEditors: Isidro Lopez isidro.lopezg@live.com
 * @LastEditTime: 2025-08-21 22:40:48
 * @FilePath: common/extension/boolean/to_string_test.go
 * @Description: 这是默认设置,可以在设置》工具》File Description中进行配置
 */
// File: common/to_string_test.go

package boolean_test

import (
	"testing"

	"github.com/entiqon/entiqon/common/extension/boolean"
)

func TestBoolToString(t *testing.T) {
	tests := []struct {
		input    bool
		trueStr  string
		falseStr string
		want     string
	}{
		{true, "yes", "no", "yes"},
		{false, "yes", "no", "no"},
		{true, "enabled", "disabled", "enabled"},
		{false, "enabled", "disabled", "disabled"},
	}

	for _, tt := range tests {
		got := boolean.BoolToString(tt.input, tt.trueStr, tt.falseStr)
		if got != tt.want {
			t.Errorf("BoolToString(%v, %q, %q) = %q; want %q",
				tt.input, tt.trueStr, tt.falseStr, got, tt.want)
		}
	}
}

func TestBoolToStrDeprecated(t *testing.T) {
	if boolean.BoolToStr(true, "yes", "no") != "yes" {
		t.Errorf("BoolToStr did not return expected value")
	}
}
