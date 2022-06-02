package utils

import gloss "github.com/charmbracelet/lipgloss"

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
