package main

import (
	"fmt"
	gloss "github.com/charmbracelet/lipgloss"
	"strings"
)

func (m model) View() string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(borderWindow)
	entry := gloss.NewStyle().
		Width(int(float64(m.width)*0.9)).
		Border(gloss.NormalBorder(), false, false, true)

	var body strings.Builder

	if m.addPrompt.enabled {
		body.WriteString("Add Torrent from ")
		if m.addPrompt.magnet {
			body.WriteString("Magnet Link")
		} else {
			body.WriteString("File")
		}
		body.WriteString("\n")
		body.WriteString(borderWindow.Render(m.addPrompt.input.View()))
		return fullscreen.Render(
			gloss.Place(
				m.width, m.height,
				gloss.Center, gloss.Center,
				body.String(),
				gloss.WithWhitespaceChars("⑀"),
				gloss.WithWhitespaceForeground(gloss.Color("#383838")),
			),
		)
	} else {
		if torrents := m.client.Torrents(); len(torrents) > 0 {
			for _, t := range torrents {
				name := t.Name()
				stats := t.Stats()

				peers := fmt.Sprintf("%d peers", stats.ActivePeers)

				spacerWidth := int(float64(m.width)*0.9) - gloss.Width(name) - gloss.Width(peers)

				body.WriteString(
					entry.Render(
						gloss.JoinHorizontal(
							gloss.Center,
							t.Name(),
							gloss.NewStyle().Width(spacerWidth).Render(""),
							peers,
						),
					) + "\n",
				)
			}
		} else {
			body.WriteString(entry.Render("You have no torrents!"))
		}
		body.WriteString("\n")
	}

	return fullscreen.Render(body.String())
}
