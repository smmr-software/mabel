// Package styles defines the color themes and Lip Gloss styles for
// Mabel based on the user config file.
package styles

import gloss "github.com/charmbracelet/lipgloss"

// Define the Lip Gloss styles for borders and bold text
var (
	Bold = gloss.NewStyle().Bold(true)

	Border = gloss.RoundedBorder()
	Window = gloss.NewStyle().
		Align(gloss.Center).
		BorderStyle(Border)
	Fullscreen = Window.Copy()
)
