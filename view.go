package main

import (
	"fmt"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"sort"
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
	popup := gloss.NewStyle().
		Width(m.width / 3).
		Height(m.height / 3).
		Inherit(borderWindow)

	var body strings.Builder

	if m.addPrompt.enabled {
		body.WriteString("Add Torrent\n")
		body.WriteString(borderWindow.Render(m.addPrompt.torrent.View()))
		body.WriteString("\n\nSave Directory (Optional)\n")
		body.WriteString(borderWindow.Render(m.addPrompt.saveDir.View()))

		content := body.String()
		help := m.help.View(addPromptKeys)
		padding := ((m.height - gloss.Height(content)) / 2) - gloss.Height(help)
		if padding < 0 {
			padding = 0
		}
		body.WriteString(strings.Repeat("\n", padding) + help + "\n")

		return fullscreen.Render(
			gloss.Place(
				m.width, m.height,
				gloss.Center, gloss.Bottom,
				body.String(),
				gloss.WithWhitespaceChars("⑀"),
				gloss.WithWhitespaceForeground(gloss.Color("#383838")),
			),
		)
	} else if m.err != nil {
		body.WriteString(title.Render("Error") + "\n\n")
		body.WriteString(m.err.Error())
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
		if torrents := m.client.Torrents(); len(torrents) > 0 {
			sort.Slice(
				torrents,
				func(i, j int) bool {
					var (
						firstHash  = torrents[i].InfoHash()
						secondHash = torrents[j].InfoHash()

						firstTime  = m.torrentMeta[firstHash]
						secondTime = m.torrentMeta[secondHash]
					)

					return firstTime.Before(secondTime)
				},
			)
			for _, t := range torrents {
				selected := ""
				name := t.Name()
				stats := t.Stats()
				info := t.Info()

				if m.selected == t.InfoHash() {
					selected = "* "
				}

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

				spacerWidth := int(float64(m.width)*0.9) - gloss.Width(selected) - gloss.Width(name) - gloss.Width(meta)

				body.WriteString(
					entry.Render(
						gloss.JoinHorizontal(
							gloss.Center,
							selected,
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
		help := m.help.View(homeKeys)
		padding := m.height - gloss.Height(content) - gloss.Height(help)
		if padding < 0 {
			padding = 0
		}
		return fullscreen.Render(content + strings.Repeat("\n", padding) + help)
	}
}
