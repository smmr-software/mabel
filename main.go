package main

import (
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"os"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - gloss.Width(border.Right+border.Left)
		m.height = msg.Height - gloss.Width(border.Bottom+border.Top)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if !m.addPrompt.enabled {
				return m, tea.Quit
			} else {
				m.addPrompt.enabled = false
			}
		case "a":
			if !m.addPrompt.enabled {
				m.addPrompt.enabled = true
			}
		case "h":
			if m.addPrompt.enabled && m.addPrompt.phase == "tab-select" {
				m.addPrompt.state.download = "magnet"
			}
		case "l":
			if m.addPrompt.enabled && m.addPrompt.phase == "tab-select" {
				m.addPrompt.state.download = "torrent"
			}
		}
	}
	return m, nil
}

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
			body += button.inactive.Render(" Torrent ") + "\n\n"
		} else {
			body += button.inactive.Render(" Magnet ") + "    "
			body += button.active.Render(" Torrent ") + "\n\n"
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

func main() {
	if err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start(); err != nil {
		os.Exit(1)
	}
}
