package config

import "strings"

func splitComma(txt string) []string {
	if txt == "" {
		return []string{}
	}
	return strings.Split(txt, ",")
}
func Ternary[T any](statement bool, a, b T) T {
	if statement {
		return a
	}
	return b
}
