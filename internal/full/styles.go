package full

import gloss "github.com/charmbracelet/lipgloss"

var (
	bold = gloss.NewStyle().Bold(true)

	primaryBlue = gloss.Color("#5FFFD7")
	lightBlue   = gloss.Color("#DFFFF7")
	darkBlue    = gloss.Color("#00AC81")
	errorRed    = gloss.Color("#FF5E87")
	tooltip     = gloss.NewStyle().Foreground(gloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	}).Padding(0, 2)

	border       = gloss.RoundedBorder()
	borderWindow = gloss.NewStyle().
			Align(gloss.Center).
			BorderStyle(border).
			BorderForeground(primaryBlue)
)
