package main

import gloss "github.com/charmbracelet/lipgloss"

var (
	bold = gloss.NewStyle().Bold(true)

	mainColor  = gloss.Color("#5FFFD7")
	errorColor = gloss.Color("#FF5E87")
	tooltip    = gloss.NewStyle().Foreground(gloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	}).Padding(0, 2)

	border       = gloss.RoundedBorder()
	borderWindow = gloss.NewStyle().
			Align(gloss.Center).
			BorderStyle(border).
			BorderForeground(mainColor)
)
