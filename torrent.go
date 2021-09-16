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

	if strings.HasPrefix(input, "magnet:") {
		return addMagnetLink(m, input, storageDir)
	} else if strings.HasPrefix(input, "infohash:") {
		return addInfoHash(m, input, storageDir)
	} else {
		return addFromFile(m, input, storageDir)
	}
}

func addMagnetLink(m model, input string, dir storage.ClientImpl) (tea.Model, tea.Cmd) {
	var spec *torrent.TorrentSpec
	if spc, err := torrent.TorrentSpecFromMagnetUri(input); err != nil {
		m.addPrompt = initialAddPrompt()
		return m, reportError(err)
	} else {
		spc.Storage = dir
		spec = spc
	}

	if t, _, err := m.client.AddTorrentSpec(spec); err != nil {
		m.addPrompt = initialAddPrompt()
		return m, reportError(err)
	} else {
		m.torrentMeta[t.InfoHash()] = time.Now()
		m.addPrompt = initialAddPrompt()
		return m, downloadTorrent(t)
	}
}

func addInfoHash(m model, input string, dir storage.ClientImpl) (tea.Model, tea.Cmd) {
	hash := metainfo.NewHashFromHex(strings.TrimPrefix(input, "infohash:"))
	t, _ := m.client.AddTorrentInfoHashWithStorage(hash, dir)
	m.torrentMeta[t.InfoHash()] = time.Now()
	m.addPrompt = initialAddPrompt()
	return m, downloadTorrent(t)
}

func addFromFile(m model, input string, dir storage.ClientImpl) (tea.Model, tea.Cmd) {
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
	spec.Storage = dir

	if t, _, err := m.client.AddTorrentSpec(spec); err != nil {
		m.addPrompt = initialAddPrompt()
		return m, reportError(err)
	} else {
		m.torrentMeta[t.InfoHash()] = time.Now()
		m.addPrompt = initialAddPrompt()
		return m, downloadTorrent(t)
	}
}
