package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"os"
)

type model struct {
	width     int
	height    int
	addPrompt struct {
		enabled bool
		magnet  bool
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
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
			if m.addPrompt.enabled && m.addPrompt.magnet {
				m.addPrompt.magnet = false
			}
		case "l":
			if m.addPrompt.enabled && !m.addPrompt.magnet {
				m.addPrompt.magnet = true
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	fullscreen := gloss.NewStyle().
		Width(m.width - gloss.Width(border.Right+border.Left)).
		Height(m.height - gloss.Width(border.Bottom+border.Top)).
		Inherit(borderWindow)
	popup := gloss.NewStyle().
		Width(m.width / 3).
		Height(m.height / 3).
		MarginTop(m.height / 3).
		Inherit(borderWindow)

	var display string

	if m.addPrompt.enabled {
		body := title.Render("Add Torrent") + "\n\n"

		if m.addPrompt.magnet {
			body += button.inactive.Render(" Torrent ") + "    "
			body += button.active.Render(" Magnet ") + "\n\n"
		} else {
			body += button.active.Render(" Torrent ") + "    "
			body += button.inactive.Render(" Magnet ") + "\n\n"
		}

		display = popup.Render(body)
	}
	return fmt.Sprint(fullscreen.Render(display))
}

func main() {
	if err := tea.NewProgram(model{}, tea.WithAltScreen()).Start(); err != nil {
		os.Exit(1)
	}
}
