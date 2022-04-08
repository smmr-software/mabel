package styles

import gloss "github.com/charmbracelet/lipgloss"

var (
	Bold = gloss.NewStyle().Bold(true)

	Border       = gloss.RoundedBorder()
	BorderWindow = gloss.NewStyle().
			Align(gloss.Center).
			BorderStyle(Border).
			BorderForeground(DefaultTheme.Primary)
)
