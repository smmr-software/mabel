package main

import (
	"sort"
	"strconv"
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
type selectedTorrentChanged struct{}
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
		switch {
		case msg.Type == tea.KeyCtrlC:
			if m.client != nil {
				m.client.Close()
			}
			return m, tea.Quit
		case m.portStartupFailure.enabled:
			return portStartupFailureKeyPress(&m, &msg)
		case m.addPrompt.enabled:
			return addPromptKeyPress(&m, &msg)
		case m.err != nil:
			m.err = nil
			return m, nil
		case m.viewingTorrentDetails:
			m.viewingTorrentDetails = false
			return m, nil
		default:
			return defaultKeyPress(&m, &msg)
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

func portStartupFailureKeyPress(m *model, msg *tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "backspace":
		var cmd tea.Cmd
		m.portStartupFailure.port, cmd = m.portStartupFailure.port.Update(*msg)
		return m, cmd
	default:
		port, err := strconv.Atoi(m.portStartupFailure.port.Value())
		if err != nil {
			return m, reportError(err)
		}

		config := genMabelConfig()
		config.ListenPort = port
		client, err := torrent.NewClient(config)
		if err != nil {
			return m, reportError(err)
		}

		m.client = client
		m.clientConfig = config
	}

	return m, nil
}

func addPromptKeyPress(m *model, msg *tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(*msg, addPromptKeys.quit):
		m.addPrompt = initialAddPrompt()
		return m, nil
	case key.Matches(*msg, addPromptKeys.forward):
		if m.addPrompt.dir {
			return addTorrent(m)
		} else {
			m.addPrompt.torrent.Blur()
			m.addPrompt.saveDir.Focus()
			m.addPrompt.dir = true
			return m, nil
		}
	case key.Matches(*msg, addPromptKeys.back):
		if m.addPrompt.dir {
			m.addPrompt.saveDir.Blur()
			m.addPrompt.torrent.Focus()
			m.addPrompt.dir = false
		}
		return m, nil
	default:
		var cmd tea.Cmd
		if m.addPrompt.dir {
			m.addPrompt.saveDir, cmd = m.addPrompt.saveDir.Update(*msg)
		} else {
			m.addPrompt.torrent, cmd = m.addPrompt.torrent.Update(*msg)
		}
		return m, cmd
	}
}

func defaultKeyPress(m *model, msg *tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(*msg, homeKeys.quit):
		m.client.Close()
		return m, tea.Quit
	case key.Matches(*msg, homeKeys.help):
		m.help.ShowAll = !m.help.ShowAll
	case key.Matches(*msg, homeKeys.add):
		m.addPrompt.torrent.Focus()
		m.addPrompt.enabled = true
	case key.Matches(*msg, homeKeys.details):
		if t, ok := m.client.Torrent(m.selected); ok && t.Info() != nil {
			m.viewingTorrentDetails = true
		}
	case key.Matches(*msg, homeKeys.down, homeKeys.up):
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

					firstTime  = m.torrentMeta[firstHash].added
					secondTime = m.torrentMeta[secondHash].added
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
					if key.Matches(*msg, homeKeys.down) {
						index = i + 1
					} else {
						index = i - 1
					}
					break
				}
			}

			if key.Matches(*msg, homeKeys.down) && index == len(torrents) {
				m.selected = torrents[0].InfoHash()
			} else if key.Matches(*msg, homeKeys.up) && index < 0 {
				m.selected = torrents[len(torrents)-1].InfoHash()
			} else {
				m.selected = torrents[index].InfoHash()
			}
		}
		return m, func() tea.Msg { return selectedTorrentChanged{} }
	case key.Matches(*msg, homeKeys.delete):
		zero := metainfo.Hash{}
		if m.selected != zero {
			t, b := m.client.Torrent(m.selected)
			if b {
				t.Drop()
			}
		}
	case key.Matches(*msg, homeKeys.deselect):
		m.selected = metainfo.Hash{}
	}
	return m, nil
}

func tick() tea.Cmd {
	return tea.Tick(time.Duration(interval), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
