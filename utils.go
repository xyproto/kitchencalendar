package kitchencalendar

import (
	"os"
	"unicode"
)

// capitalize makes it so that the first run of a given string is in uppercase
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
