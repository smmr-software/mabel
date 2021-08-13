package main

import (
	torrent "github.com/anacrolix/torrent"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	home "github.com/mitchellh/go-homedir"
)

type magnetLinkMsg string
type torrentFromFileMsg string
type torrentDownloadStarted struct{}

func addMagnetLink(uri string) tea.Cmd {
	return func() tea.Msg { return magnetLinkMsg(uri) }
}
func addTorrentFromFile(path string) tea.Cmd {
	return func() tea.Msg { return torrentFromFileMsg(path) }
}
func downloadTorrent(t *torrent.Torrent) tea.Cmd {
	return func() tea.Msg {
		<-t.GotInfo()
		t.DownloadAll()
		return torrentDownloadStarted{}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - gloss.Width(border.Right+border.Left)
		m.height = msg.Height - gloss.Width(border.Bottom+border.Top)
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		} else if m.typing {
			return typingKeyPress(m, msg)
		} else if m.addPrompt.enabled {
			return addPromptKeyPress(m, msg)
		} else {
			return defaultKeyPress(m, msg)
		}
	case magnetLinkMsg:
		t, _ := m.torrent.AddMagnet(string(msg))
		m.addPrompt = initialAddPrompt()
		return m, downloadTorrent(t)
	case torrentFromFileMsg:
		t, _ := m.torrent.AddTorrentFromFile(string(msg))
		m.addPrompt = initialAddPrompt()
		return m, downloadTorrent(t)
	}
	return m, nil
}

func typingKeyPress(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter, tea.KeyTab:
		if m.addPrompt.enabled {
			switch m.addPrompt.phase {
			case "magnet-link", "torrent-path":
				m.addPrompt.phase = "save-path"
				m.addPrompt.state.magnetLink.Blur()
				m.addPrompt.state.torrentPath.Blur()
				m.addPrompt.state.savePath.Focus()
			case "save-path":
				m.addPrompt.phase = "approval"
				m.addPrompt.state.savePath.Blur()
				m.typing = false
			}
		}
	case tea.KeyShiftTab, tea.KeyEsc:
		if m.addPrompt.enabled {
			switch m.addPrompt.phase {
			case "magnet-link", "torrent-path":
				m.addPrompt.phase = "tab-select"
				m.addPrompt.state.magnetLink.Blur()
				m.addPrompt.state.torrentPath.Blur()
				m.typing = false
			case "save-path":
				m.addPrompt.state.savePath.Blur()
				if m.addPrompt.state.download == "magnet" {
					m.addPrompt.phase = "magnet-link"
					m.addPrompt.state.magnetLink.Focus()
				} else {
					m.addPrompt.phase = "torrent-path"
					m.addPrompt.state.torrentPath.Focus()
				}
			}
		}
	default:
		if m.addPrompt.enabled {
			var cmd tea.Cmd
			switch m.addPrompt.phase {
			case "magnet-link":
				m.addPrompt.state.magnetLink, cmd = m.addPrompt.state.magnetLink.Update(msg)
			case "torrent-path":
				m.addPrompt.state.torrentPath, cmd = m.addPrompt.state.torrentPath.Update(msg)
			case "save-path":
				m.addPrompt.state.savePath, cmd = m.addPrompt.state.savePath.Update(msg)
			}
			return m, cmd
		}
	}
	return m, nil
}

func addPromptKeyPress(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.addPrompt = initialAddPrompt()
	case "h":
		if m.addPrompt.phase == "tab-select" {
			m.addPrompt.state.download = "magnet"
		}
	case "l":
		if m.addPrompt.phase == "tab-select" {
			m.addPrompt.state.download = "torrent"
		}
	case "enter", "tab":
		switch m.addPrompt.phase {
		case "tab-select":
			m.typing = true
			if m.addPrompt.state.download == "magnet" {
				m.addPrompt.phase = "magnet-link"
				m.addPrompt.state.magnetLink.Focus()
			} else {
				m.addPrompt.phase = "torrent-path"
				m.addPrompt.state.torrentPath.Focus()
			}
		case "approval":
			if m.addPrompt.state.download == "magnet" {
				return m, addMagnetLink(m.addPrompt.state.magnetLink.Value())
			} else {
				path, _ := home.Expand(m.addPrompt.state.torrentPath.Value())
				return m, addTorrentFromFile(path)
			}
		}
	case "shift+tab":
		if m.addPrompt.phase == "approval" {
			m.addPrompt.phase = "save-path"
			m.addPrompt.state.savePath.Focus()
			m.typing = true
		}
	}
	return m, nil
}

func defaultKeyPress(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "a":
		m.addPrompt.enabled = true
	}
	return m, nil
}
