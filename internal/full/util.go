package full

import gloss "github.com/charmbracelet/lipgloss"

func truncateForMinimumSpacing(str string, spacing *int, min int) string {
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
