package main

import (
	gloss "github.com/charmbracelet/lipgloss"
	"strings"
)

func (m model) View() string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(borderWindow)
	popup := gloss.NewStyle().
		Width(m.width / 3).
		Height(m.height / 3).
		Inherit(borderWindow)
	entry := gloss.NewStyle().
		Align(gloss.Center).
		Width(int(float64(m.width)*0.9)).
		Border(gloss.NormalBorder(), false, false, true)

	var body strings.Builder

	if m.addPrompt.enabled {
		body.WriteString(title.Render("Add Torrent") + "\n\n")
		if m.addPrompt.state.download == "magnet" {
			body.WriteString(button.active.Render(" Magnet ") + "    ")
			body.WriteString(button.inactive.Render(" Torrent ") + "\n\n\n")
			body.WriteString(m.addPrompt.state.magnetLink.View() + "\n\n")
		} else {
			body.WriteString(button.inactive.Render(" Magnet ") + "    ")
			body.WriteString(button.active.Render(" Torrent ") + "\n\n\n")
			body.WriteString(m.addPrompt.state.torrentPath.View() + "\n\n")
		}
		body.WriteString(m.addPrompt.state.savePath.View() + "\n\n\n\n")
		if m.addPrompt.phase == "approval" {
			body.WriteString(button.active.Render(" Start Download "))
		} else {
			body.WriteString(button.inactive.Render(" Start Download "))
		}
		return fullscreen.Render(
			gloss.Place(
				m.width, m.height,
				gloss.Center, gloss.Center,
				popup.Render(body.String()),
				gloss.WithWhitespaceChars("⑀"),
				gloss.WithWhitespaceForeground(gloss.Color("#383838")),
			),
		)
	} else {
		if torrents := m.torrent.client.Torrents(); len(torrents) > 0 {
			body.WriteString(entry.Render("You have some torrents!"))
		} else {
			body.WriteString(entry.Render("You have no torrents!"))
		}
		body.WriteString("\n")
	}

	return fullscreen.Render(body.String())
}
