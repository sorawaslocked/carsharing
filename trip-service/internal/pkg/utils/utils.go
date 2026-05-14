package utils

import "strings"

func Uncapitalize(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
