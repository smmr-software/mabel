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
	}
	return m, nil
}

func addPromptKeyPress(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, addPromptKeys.quit):
		m.addPrompt = initialAddPrompt()
	case key.Matches(msg, addPromptKeys.forward):
		if m.addPrompt.dir {
			var (
				saveDir    string
				storageDir storage.ClientImpl
			)

			input := m.addPrompt.torrent.Value()
			if dir, err := home.Expand(m.addPrompt.saveDir.Value()); err != nil {
				m.addPrompt = initialAddPrompt()
				return m, reportError(err)
			} else {
				saveDir = dir
			}

			metadataDirectory := os.TempDir()
			if metadataStorage, err := storage.NewDefaultPieceCompletionForDir(metadataDirectory); err != nil {
				storageDir = storage.NewMMap(saveDir)
			} else {
				storageDir = storage.NewMMapWithCompletion(saveDir, metadataStorage)
			}

			var t *torrent.Torrent
			if strings.HasPrefix(input, "magnet:") {
				var spec *torrent.TorrentSpec
				if spc, err := torrent.TorrentSpecFromMagnetUri(input); err != nil {
					m.addPrompt = initialAddPrompt()
					return m, reportError(err)
				} else {
					spc.Storage = storageDir
					spec = spc
				}

				if tr, _, err := m.client.AddTorrentSpec(spec); err != nil {
					m.addPrompt = initialAddPrompt()
					return m, reportError(err)
				} else {
					t = tr
				}
			} else if strings.HasPrefix(input, "infohash:") {
				hash := metainfo.NewHashFromHex(strings.TrimPrefix(input, "infohash:"))
				t, _ = m.client.AddTorrentInfoHashWithStorage(hash, storageDir)
			} else {
				var (
					path string
					meta *metainfo.MetaInfo
				)

				if p, err := home.Expand(input); err != nil {
					m.addPrompt = initialAddPrompt()
					return m, reportError(err)
				} else {
					path = p
				}

				if mt, err := metainfo.LoadFromFile(path); err != nil {
					m.addPrompt = initialAddPrompt()
					return m, reportError(err)
				} else {
					meta = mt
				}

				spec := torrent.TorrentSpecFromMetaInfo(meta)
				spec.Storage = storageDir

				if tr, _, err := m.client.AddTorrentSpec(spec); err != nil {
					m.addPrompt = initialAddPrompt()
					return m, reportError(err)
				} else {
					t = tr
				}
			}
			m.torrentMeta[t.InfoHash()] = time.Now()
			m.addPrompt = initialAddPrompt()
			return m, downloadTorrent(t)
		} else {
			m.addPrompt.torrent.Blur()
			m.addPrompt.saveDir.Focus()
			m.addPrompt.dir = true
		}
	case key.Matches(msg, addPromptKeys.back):
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
	case key.Matches(msg, homeKeys.quit):
		m.client.Close()
		return m, tea.Quit
	case key.Matches(msg, homeKeys.help):
		m.help.ShowAll = !m.help.ShowAll
	case key.Matches(msg, homeKeys.addTorrent):
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
