// Package strutil offers common string utility functions.
package strutil

// Stringify converts a string pointer into a string.
func Stringify(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
