package main

import (
	"os"
	"unicode"
)

// capitalize makes changes the first rune of a string to be in uppercase
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	firstRune := unicode.ToUpper(runes[0])
	return string(append([]rune{firstRune}, runes[1:]...))
}

// exists checks if the given path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
