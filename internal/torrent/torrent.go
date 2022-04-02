package torrent

import (
	"os"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type TorrentDownloadStarted struct{}

var magnetPrefix = "magnet:"
var infohashPrefix = "infohash:"
var hashLength = 40

func AddTorrents(t *[]*string, dir *string, client *torrent.Client, l *list.Model) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, v := range *t {
		cmd, err := AddTorrent(v, dir, client, l)
		if err != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func AddTorrent(t, dir *string, client *torrent.Client, l *list.Model) (tea.Cmd, error) {
	store := getStorage(dir)
	if strings.HasPrefix(*t, magnetPrefix) {
		return addMagnetLink(t, &store, client, l)
	} else if strings.HasPrefix(*t, infohashPrefix) || len(*t) == hashLength {
		return addInfoHash(t, &store, client, l)
	} else {
		return addFromFile(t, &store, client, l)
	}
}

func downloadTorrent(t *torrent.Torrent) tea.Cmd {
	return func() tea.Msg {
		<-t.GotInfo()
		t.DownloadAll()
		return TorrentDownloadStarted{}
	}
}

func getStorage(dir *string) storage.ClientImpl {
	metadataDirectory := os.TempDir()
	if metadataStorage, err := storage.NewDefaultPieceCompletionForDir(metadataDirectory); err != nil {
		return storage.NewMMap(*dir)
	} else {
		return storage.NewMMapWithCompletion(*dir, metadataStorage)
	}
}
