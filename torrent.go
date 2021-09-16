package main

import (
	torrent "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	tea "github.com/charmbracelet/bubbletea"
	home "github.com/mitchellh/go-homedir"
	"os"
	"strings"
	"time"
)

func addTorrent(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
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
}
