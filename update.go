package main

import (
	"sort"
	"time"

	torrent "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

const interval = 500 * time.Millisecond

type tickMsg time.Time
type torrentDownloadStarted struct{}
type mabelError error

func downloadTorrent(t *torrent.Torrent) tea.Cmd {
	return func() tea.Msg {
		<-t.GotInfo()
		t.DownloadAll()
		return torrentDownloadStarted{}
	}
}

func reportError(err error) tea.Cmd {
	return func() tea.Msg {
		return mabelError(err)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width := msg.Width - gloss.Width(border.Right+border.Left)
		height := msg.Height - gloss.Width(border.Bottom+border.Top)

		m.width = width
		m.help.Width = width
		m.height = height
		return m, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.client.Close()
			return m, tea.Quit
		} else if m.addPrompt.enabled {
			return addPromptKeyPress(m, msg)
		} else if m.err != nil {
			m.err = nil
			return m, nil
		} else {
			return defaultKeyPress(m, msg)
		}
	case tickMsg:
		return m, tick()
	case mabelError:
		m.err = msg
		return m, nil
	default:
		return m, nil
	}
}

func addPromptKeyPress(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, addPromptKeys.quit):
		m.addPrompt = initialAddPrompt()
		return m, nil
	case key.Matches(msg, addPromptKeys.forward):
		if m.addPrompt.dir {
			return addTorrent(m, msg)
		} else {
			m.addPrompt.torrent.Blur()
			m.addPrompt.saveDir.Focus()
			m.addPrompt.dir = true
			return m, nil
		}
	case key.Matches(msg, addPromptKeys.back):
		if m.addPrompt.dir {
			m.addPrompt.saveDir.Blur()
			m.addPrompt.torrent.Focus()
			m.addPrompt.dir = false
		}
		return m, nil
	default:
		var cmd tea.Cmd
		if m.addPrompt.dir {
			m.addPrompt.saveDir, cmd = m.addPrompt.saveDir.Update(msg)
		} else {
			m.addPrompt.torrent, cmd = m.addPrompt.torrent.Update(msg)
		}
		return m, cmd
	}
}

func defaultKeyPress(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, homeKeys.quit):
		m.client.Close()
		return m, tea.Quit
	case key.Matches(msg, homeKeys.help):
		m.help.ShowAll = !m.help.ShowAll
	case key.Matches(msg, homeKeys.addTorrent):
		m.addPrompt.torrent.Focus()
		m.addPrompt.enabled = true
	case key.Matches(msg, homeKeys.down):
		torrents := m.client.Torrents()
		if len(torrents) == 0 {
			return m, nil
		}
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

		zero := metainfo.Hash{}
		if m.selected == zero {
			m.selected = torrents[0].InfoHash()
		} else {
			var index int
			for i := range torrents {
				if torrents[i].InfoHash() == m.selected {
					index = i + 1
					break
				}
			}
			if index == len(torrents) {
				m.selected = torrents[0].InfoHash()
			} else {
				m.selected = torrents[index].InfoHash()
			}
		}
	case key.Matches(msg, homeKeys.up):
		torrents := m.client.Torrents()
		if len(torrents) == 0 {
			return m, nil
		}
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

		zero := metainfo.Hash{}
		if m.selected == zero {
			m.selected = torrents[0].InfoHash()
		} else {
			var index int
			for i := range torrents {
				if torrents[i].InfoHash() == m.selected {
					index = i - 1
					break
				}
			}
			if index < 0 {
				m.selected = torrents[len(torrents)-1].InfoHash()
			} else {
				m.selected = torrents[index].InfoHash()
			}
		}
	case key.Matches(msg, homeKeys.delete):
		zero := metainfo.Hash{}
		if m.selected != zero {
			t, b := m.client.Torrent(m.selected)
			if b {
				t.Drop()
			}
		}
		return m, nil
	}
	return m, nil
}

func tick() tea.Cmd {
	return tea.Tick(time.Duration(interval), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
