package torrent

import (
	"os"
	"strings"

	"github.com/smmr-software/mabel/internal/list"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type TorrentDownloadStarted struct{}

var magnetPrefix = "magnet:"
var infohashPrefix = "infohash:"
var hashLength = 40

func AddTorrent(t, dir *string, client *torrent.Client) (tea.Cmd, bool, list.Item, error) {
	store := getStorage(dir)
	if strings.HasPrefix(*t, magnetPrefix) {
		return addMagnetLink(t, &store, client)
	} else if strings.HasPrefix(*t, infohashPrefix) || len(*t) == hashLength {
		return addInfoHash(t, &store, client)
	} else {
		return addFromFile(t, &store, client)
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
