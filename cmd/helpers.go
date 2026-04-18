package cmd

import "strings"

// join is strings.Join without importing strings in every file.
func join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}
