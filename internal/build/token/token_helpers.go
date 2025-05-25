package token

import "strings"

// trim returns the input string with leading and trailing spaces removed.
func trim(s string) string {
	return strings.TrimSpace(s)
}

// splitAlias attempts to split a raw expression using "AS" as a delimiter.
// It returns two parts if the alias keyword is found.
func splitAlias(expr string) []string {
	return strings.SplitN(expr, "AS", 2)
}

// firstOrEmpty returns the first string in the slice if it exists,
// otherwise it returns an empty string.
func firstOrEmpty(a []string) string {
	if len(a) > 0 {
		return a[0]
	}
	return ""
}
