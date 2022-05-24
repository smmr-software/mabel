// Package torrent runs the torrent downloads of the client in
// metainfo, infohash, and magnet link formats.
package torrent

import (
	"os"
	"strings"

	"github.com/smmr-software/mabel/internal/styles"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	home "github.com/mitchellh/go-homedir"
)

type TorrentDownloadStarted struct{}

var magnetPrefix = "magnet:"
var infohashPrefix = "infohash:"
var hashLength = 40

// AddTorrents takes a group of torrents as strings, adds them to the
// client, and returns them batched.
func AddTorrents(t *[]string, dir *string, client *torrent.Client, l *list.Model, theme *styles.ColorTheme) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, v := range *t {
		cmd, err := AddTorrent(&v, dir, client, l, theme)
		if err == nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// AddTorrent takes a torrent as a string and adds it to the client
// based on its format.
func AddTorrent(t, dir *string, client *torrent.Client, l *list.Model, theme *styles.ColorTheme) (tea.Cmd, error) {
	store := getStorage(dir)
	if strings.HasPrefix(*t, magnetPrefix) {
		return addMagnetLink(t, &store, client, l, theme)
	} else if strings.HasPrefix(*t, infohashPrefix) || len(*t) == hashLength {
		return addInfoHash(t, &store, client, l, theme)
	} else {
		return addFromFile(t, &store, client, l, theme)
	}
}

// downloadTorrent asynchronously waits for torrent info to arrive,
// triggers the download, and returns a Bubble Tea start message.
func downloadTorrent(t *torrent.Torrent) tea.Cmd {
	return func() tea.Msg {
		<-t.GotInfo()
		t.DownloadAll()
		return TorrentDownloadStarted{}
	}
}

// getStorage returns a storage implementation that writes downloaded
// files to a user-defined directory, and writes unnecessary files to a
// temporary directory.
func getStorage(dir *string) storage.ClientImpl {
	var err error
	*dir, err = home.Expand(*dir)
	if err != nil {
		*dir = ""
	}

	metadataDirectory := os.TempDir()
	if metadataStorage, err := storage.NewDefaultPieceCompletionForDir(metadataDirectory); err != nil {
		return storage.NewMMap(*dir)
	} else {
		return storage.NewMMapWithCompletion(*dir, metadataStorage)
	}
}
