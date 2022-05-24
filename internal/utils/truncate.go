package utils

import gloss "github.com/charmbracelet/lipgloss"

// TruncateForMinimumSpacing ensures that a string fits into screen
// space with a given amount of padding, by possibly truncating the
// string until it fits in the width with the padding.
func TruncateForMinimumSpacing(str string, spacing *int, min int) string {
	runes := []rune(str)
	initial := len(runes)
	for *spacing-gloss.Width(string(runes)) < min {
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

	*spacing -= gloss.Width(final)
	if *spacing < 0 {
		*spacing = 0
	}

	return final
}
