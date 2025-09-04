package helpers

import (
	"reflect"
	"testing"
)

func TestCoerceScalar(t *testing.T) {
	cases := map[string]any{
		"42":    42,    // int
		"3.14":  3.14,  // float
		"NULL":  nil,   // null
		"'abc'": "abc", // quoted string
		`"xyz"`: "xyz", // double-quoted string
		"raw":   "raw", // fallback
	}
	for in, want := range cases {
		if got := coerceScalar(in); !reflect.DeepEqual(got, want) {
			t.Errorf("coerceScalar(%q)=%v, want %v", in, got, want)
		}
	}
}

func TestIsBoundary(t *testing.T) {
	tests := []struct {
		in   byte
		want bool
	}{
		{' ', true},
		{'\t', true},
		{'\n', true},
		{'\r', true},
		{'(', true},
		{')', true},
		{',', true},
		{';', true},
		{'x', false}, // hit default branch
	}
	for _, tt := range tests {
		if got := isBoundary(tt.in); got != tt.want {
			t.Errorf("isBoundary(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestIndexOf(t *testing.T) {
	// symbol operator branch
	if got := indexOf("a != b", "!="); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}

	// word operator branch, found at start
	if got := indexOf("in (1,2,3)", "in"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}

	// word operator branch, found later → triggers i = j+1 loop advance
	if got := indexOf("x in y in z", "in"); got < 2 {
		t.Errorf("expected >=2, got %d", got)
	}

	// not found
	if got := indexOf("foo", "in"); got != -1 {
		t.Errorf("expected -1, got %d", got)
	}

	// "in" occurs at pos=3 in "print", but not bounded by space/parens → must skip
	got := indexOf("print", "in")
	if got != -1 {
		t.Errorf("expected -1 (no valid boundary match), got %d", got)
	}
}
