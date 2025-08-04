// File: common/bools_test.go

package common_test

import (
	"testing"

	"github.com/entiqon/common"
)

func TestBoolToStr(t *testing.T) {
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
		got := common.BoolToStr(tt.input, tt.trueStr, tt.falseStr)
		if got != tt.want {
			t.Errorf("BoolToStr(%v, %q, %q) = %q; want %q", tt.input, tt.trueStr, tt.falseStr, got, tt.want)
		}
	}
}
