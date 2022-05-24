package utils

import gloss "github.com/charmbracelet/lipgloss"

// TruncateForMinimumPadding ensures that a string fits into screen
// space with a given amount of padding, by truncating the string until
// it fits in the width with the padding when necessary.
func TruncateForMinimumPadding(str string, width *int, padding int) string {
	runes := []rune(str)
	initial := len(runes)
	for *width-gloss.Width(string(runes)) < padding {
		if index := len(runes) - 1; index > 0 {
			runes = runes[:index]
		} else {
			break
		}
	}
	if initial > len(runes) {
		runes[len(runes)-1] = '…'
	}

	final := string(runes)

	*width -= gloss.Width(final)
	if *width < 0 {
		*width = 0
	}

	return final
}
