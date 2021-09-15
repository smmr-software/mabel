package main

import (
	torrent "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	home "github.com/mitchellh/go-homedir"
	"os"
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
	case tea.KeyEnter, tea.KeyTab:
		if m.addPrompt.dir {
			var (
				input      = m.addPrompt.torrent.Value()
				saveDir, _ = home.Expand(m.addPrompt.saveDir.Value())

				metadataDirectory  = os.TempDir()
				metadataStorage, _ = storage.NewDefaultPieceCompletionForDir(metadataDirectory)
				dir                = storage.NewMMapWithCompletion(saveDir, metadataStorage)
			)

			var t *torrent.Torrent
			if strings.HasPrefix(input, "magnet:") {
				spec, _ := torrent.TorrentSpecFromMagnetUri(input)
				spec.Storage = dir
				t, _, _ = m.client.AddTorrentSpec(spec)
			} else if strings.HasPrefix(input, "infohash:") {
				hash := metainfo.NewHashFromHex(strings.TrimPrefix(input, "infohash:"))
				t, _ = m.client.AddTorrentInfoHashWithStorage(hash, dir)
			} else {
				var (
					path, _ = home.Expand(input)
					meta, _ = metainfo.LoadFromFile(path)
					spec    = torrent.TorrentSpecFromMetaInfo(meta)
				)
				spec.Storage = dir
				t, _, _ = m.client.AddTorrentSpec(spec)
			}
			m.torrentMeta[t.InfoHash()] = time.Now()
			m.addPrompt = initialAddPrompt()
			return m, downloadTorrent(t)
		} else {
			m.addPrompt.torrent.Blur()
			m.addPrompt.saveDir.Focus()
			m.addPrompt.dir = true
		}
	case tea.KeyShiftTab:
		if m.addPrompt.dir {
			m.addPrompt.saveDir.Blur()
			m.addPrompt.torrent.Focus()
			m.addPrompt.dir = false
		}
	default:
		var cmd tea.Cmd
		if m.addPrompt.dir {
			m.addPrompt.saveDir, cmd = m.addPrompt.saveDir.Update(msg)
		} else {
			m.addPrompt.torrent, cmd = m.addPrompt.torrent.Update(msg)
		}
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
		m.addPrompt.torrent.Focus()
		m.addPrompt.enabled = true
	}
	return m, nil
}

func tick() tea.Cmd {
	return tea.Tick(time.Duration(interval), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
