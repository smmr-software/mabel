package torrent

import (
	"os"
	"strings"
	"time"

	"github.com/smmr-software/mabel/internal/list"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"

	"github.com/acarl005/stripansi"
	home "github.com/mitchellh/go-homedir"

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

func DownloadTorrent(t *torrent.Torrent) tea.Cmd {
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

func addMagnetLink(input *string, dir *storage.ClientImpl, client *torrent.Client) (tea.Cmd, bool, list.Item, error) {
	spec, err := torrent.TorrentSpecFromMagnetUri(*input)
	if err != nil {
		return nil, false, list.Item{}, err
	}
	spec.Storage = *dir

	t, nw, err := client.AddTorrentSpec(spec)
	if err != nil {
		return nil, false, list.Item{}, err
	}

	return DownloadTorrent(t), nw, list.Item{Self: t, Added: time.Now()}, nil
}

func addInfoHash(input *string, dir *storage.ClientImpl, client *torrent.Client) (tea.Cmd, bool, list.Item, error) {
	if strings.HasPrefix(*input, infohashPrefix) {
		*input = strings.TrimPrefix(*input, infohashPrefix)
	}

	hash := metainfo.NewHashFromHex(*input)
	t, nw := client.AddTorrentInfoHashWithStorage(hash, *dir)

	return DownloadTorrent(t), nw, list.Item{Self: t, Added: time.Now()}, nil
}

func addFromFile(input *string, dir *storage.ClientImpl, client *torrent.Client) (tea.Cmd, bool, list.Item, error) {
	path, err := home.Expand(*input)
	if err != nil {
		return nil, false, list.Item{}, err
	}

	meta, err := metainfo.LoadFromFile(path)
	if err != nil {
		return nil, false, list.Item{}, err
	}

	spec := torrent.TorrentSpecFromMetaInfo(meta)
	spec.Storage = *dir

	t, nw, err := client.AddTorrentSpec(spec)
	if err != nil {
		return nil, false, list.Item{}, err
	}

	return DownloadTorrent(t), nw,
		list.Item{
			Self:    t,
			Added:   time.Now(),
			Created: time.Unix(meta.CreationDate, 0),
			Comment: stripansi.Strip(meta.Comment),
			Program: stripansi.Strip(meta.CreatedBy),
		}, nil
}
