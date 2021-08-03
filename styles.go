package main

import gloss "github.com/charmbracelet/lipgloss"

var border = gloss.RoundedBorder()
var borderWindow = gloss.NewStyle().
	Align(gloss.Center).
	BorderStyle(border).
	BorderForeground(gloss.Color("86"))
var title = gloss.NewStyle().
	Bold(true).
	Underline(true).
	Foreground(gloss.Color("#FFF"))

var button = struct {
	inactive gloss.Style
	active   gloss.Style
}{
	inactive: gloss.NewStyle().Foreground(gloss.Color("#000")).Background(gloss.Color("#FFF")),
	active:   gloss.NewStyle().Foreground(gloss.Color("#000")).Background(gloss.Color("86")),
}
