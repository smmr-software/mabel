package main

import (
	torrent "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	home "github.com/mitchellh/go-homedir"
	"strings"
	"time"
)

var interval = 500 * time.Millisecond

type tickMsg time.Time
type torrentDownloadStarted struct{}

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
		width := msg.Width - gloss.Width(border.Right+border.Left)
		height := msg.Height - gloss.Width(border.Bottom+border.Top)

		m.width = width
		m.help.Width = width
		m.height = height
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.client.Close()
			return m, tea.Quit
		} else if m.addPrompt.enabled {
			return addPromptKeyPress(m, msg)
		} else {
			return defaultKeyPress(m, msg)
		}
	case tickMsg:
		return m, tick()
	}
	return m, nil
}

func addPromptKeyPress(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.addPrompt = initialAddPrompt()
	case tea.KeyEnter:
		var t *torrent.Torrent
		input := m.addPrompt.input.Value()
		if strings.HasPrefix(input, "magnet:") {
			t, _ = m.client.AddMagnet(input)
			m.torrentMeta[t.InfoHash()] = time.Now()
		} else if strings.HasPrefix(input, "infohash:") {
			t, _ = m.client.AddTorrentInfoHash(metainfo.NewHashFromHex(strings.TrimPrefix(input, "infohash:")))
			m.torrentMeta[t.InfoHash()] = time.Now()
		} else {
			path, _ := home.Expand(input)
			t, _ = m.client.AddTorrentFromFile(path)
			m.torrentMeta[t.InfoHash()] = time.Now()
		}
		m.addPrompt = initialAddPrompt()
		return m, downloadTorrent(t)
	default:
		var cmd tea.Cmd
		m.addPrompt.input, cmd = m.addPrompt.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

func defaultKeyPress(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.quit):
		m.client.Close()
		return m, tea.Quit
	case key.Matches(msg, keys.help):
		m.help.ShowAll = !m.help.ShowAll
	case key.Matches(msg, keys.addTorrent):
		m.addPrompt.input.Focus()
		m.addPrompt.enabled = true
	}
	return m, nil
}

func tick() tea.Cmd {
	return tea.Tick(time.Duration(interval), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
