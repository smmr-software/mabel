package main

import (
	"fmt"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
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
		body.WriteString("Add Torrent\n")
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
				info := t.Info()

				var meta string
				if info == nil {
					meta = "getting torrent info..."
				} else {
					if t.BytesMissing() != 0 {
						meta = fmt.Sprintf(
							"%s/%s | %d/%d peers",
							humanize.Bytes(uint64(t.BytesCompleted())),
							humanize.Bytes(uint64(t.Length())),
							stats.ActivePeers,
							stats.TotalPeers,
						)
					} else {
						meta = "done!"
					}
				}

				spacerWidth := int(float64(m.width)*0.9) - gloss.Width(name) - gloss.Width(meta)

				body.WriteString(
					entry.Render(
						gloss.JoinHorizontal(
							gloss.Center,
							t.Name(),
							gloss.NewStyle().Width(spacerWidth).Render(""),
							meta,
						),
					) + "\n",
				)
			}
		} else {
			body.WriteString(entry.Render("You have no torrents!"))
		}
		body.WriteString("\n")

		content := body.String()
		help := m.help.View(keys)
		padding := m.height - gloss.Height(content) - gloss.Height(help)
		if padding < 0 {
			padding = 0
		}
		return fullscreen.Render(content + strings.Repeat("\n", padding) + help)
	}
}
