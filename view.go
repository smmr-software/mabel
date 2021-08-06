package main

import gloss "github.com/charmbracelet/lipgloss"

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

	var body string
	if torrents := m.torrent.client.Torrents(); len(torrents) > 0 {
		body += entry.Render("You have some torrents!")
	} else {
		body += entry.Render("You have no torrents!")
	}
	body += "\n"

	if m.addPrompt.enabled {
		body := title.Render("Add Torrent") + "\n\n"
		if m.addPrompt.state.download == "magnet" {
			body += button.active.Render(" Magnet ") + "    "
			body += button.inactive.Render(" Torrent ") + "\n\n\n"
			body += m.addPrompt.state.magnetLink.View() + "\n\n"
		} else {
			body += button.inactive.Render(" Magnet ") + "    "
			body += button.active.Render(" Torrent ") + "\n\n\n"
			body += m.addPrompt.state.torrentPath.View() + "\n\n"
		}
		body += m.addPrompt.state.savePath.View() + "\n\n\n\n"
		if m.addPrompt.phase == "approval" {
			body += button.active.Render(" Start Download ")
		} else {
			body += button.inactive.Render(" Start Download ")
		}
		return fullscreen.Render(
			gloss.Place(
				m.width, m.height,
				gloss.Center, gloss.Center,
				popup.Render(body),
				gloss.WithWhitespaceChars("⑀"),
				gloss.WithWhitespaceForeground(gloss.Color("#383838")),
			),
		)
	}

	return fullscreen.Render(body)
}
