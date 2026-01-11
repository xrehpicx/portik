package render

import "strings"

// trunc shortens a string to at most n characters, replacing newlines
// with spaces and adding a trailing ellipsis if truncated.
func trunc(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	if n <= 1 {
		return s[:n]
	}
	return s[:n-1] + "…"
}
