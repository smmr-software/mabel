package main

import (
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - gloss.Width(border.Right+border.Left)
		m.height = msg.Height - gloss.Width(border.Bottom+border.Top)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc", "q":
			if !m.addPrompt.enabled {
				return m, tea.Quit
			} else {
				m.addPrompt = initialAddPrompt()
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
		case "enter", "tab":
			if m.addPrompt.enabled {
				switch m.addPrompt.phase {
				case "tab-select":
					if m.addPrompt.state.download == "magnet" {
						m.addPrompt.phase = "magnet-link"
						m.addPrompt.state.magnetLink.Focus()
					} else {
						m.addPrompt.phase = "torrent-path"
						m.addPrompt.state.torrentPath.Focus()
					}
				case "magnet-link", "torrent-path":
					m.addPrompt.phase = "save-path"
					m.addPrompt.state.magnetLink.Blur()
					m.addPrompt.state.torrentPath.Blur()
					m.addPrompt.state.savePath.Focus()
				case "save-path":
					m.addPrompt.phase = "approval"
					m.addPrompt.state.savePath.Blur()
				}
			}
		case "shift+tab":
			if m.addPrompt.enabled {
				switch m.addPrompt.phase {
				case "magnet-link", "torrent-path":
					m.addPrompt.phase = "tab-select"
					m.addPrompt.state.magnetLink.Blur()
					m.addPrompt.state.torrentPath.Blur()
				case "save-path":
					if m.addPrompt.state.download == "magnet" {
						m.addPrompt.phase = "magnet-link"
						m.addPrompt.state.magnetLink.Focus()
					} else {
						m.addPrompt.phase = "torrent-path"
						m.addPrompt.state.torrentPath.Focus()
					}
					m.addPrompt.state.savePath.Blur()
				case "approval":
					m.addPrompt.phase = "save-path"
					m.addPrompt.state.savePath.Focus()
				}
			}
		}
	}
	return m, nil
}
