package styles

import gloss "github.com/charmbracelet/lipgloss"

var (
	Bold = gloss.NewStyle().Bold(true)

	PrimaryBlue = gloss.Color("#5FFFD7")
	LightBlue   = gloss.Color("#DFFFF7")
	DarkBlue    = gloss.Color("#00AC81")
	ErrorRed    = gloss.Color("#FF5E87")
	Tooltip     = gloss.NewStyle().Foreground(gloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	}).Padding(0, 2)

	Border       = gloss.RoundedBorder()
	BorderWindow = gloss.NewStyle().
			Align(gloss.Center).
			BorderStyle(Border).
			BorderForeground(PrimaryBlue)
)
