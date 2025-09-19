package boolean_test

import (
	"testing"

	"github.com/entiqon/common/extension/boolean"
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
