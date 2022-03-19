package shared

import (
	"os"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"

	home "github.com/mitchellh/go-homedir"

	tea "github.com/charmbracelet/bubbletea"
)

type TorrentDownloadStarted struct{}

var magnetPrefix = "magnet:"
var infohashPrefix = "infohash:"
var hashLength = 40

func AddTorrent(t, dir *string, client *torrent.Client) (tea.Cmd, error) {
	store := getStorage(dir)
	if strings.HasPrefix(*t, magnetPrefix) {
		return addMagnetLink(t, &store, client)
	} else if strings.HasPrefix(*t, infohashPrefix) || len(*t) == hashLength {
		return addInfoHash(t, &store, client), nil
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

func addMagnetLink(input *string, dir *storage.ClientImpl, client *torrent.Client) (tea.Cmd, error) {
	spec, err := torrent.TorrentSpecFromMagnetUri(*input)
	if err != nil {
		return nil, err
	}
	spec.Storage = *dir

	t, _, err := client.AddTorrentSpec(spec)
	if err != nil {
		return nil, err
	}

	return DownloadTorrent(t), nil
}

func addInfoHash(input *string, dir *storage.ClientImpl, client *torrent.Client) tea.Cmd {
	if strings.HasPrefix(*input, infohashPrefix) {
		*input = strings.TrimPrefix(*input, infohashPrefix)
	}

	hash := metainfo.NewHashFromHex(*input)
	t, _ := client.AddTorrentInfoHashWithStorage(hash, *dir)

	return DownloadTorrent(t)
}

func addFromFile(input *string, dir *storage.ClientImpl, client *torrent.Client) (tea.Cmd, error) {
	path, err := home.Expand(*input)
	if err != nil {
		return nil, err
	}

	meta, err := metainfo.LoadFromFile(path)
	if err != nil {
		return nil, err
	}

	spec := torrent.TorrentSpecFromMetaInfo(meta)
	spec.Storage = *dir

	t, _, err := client.AddTorrentSpec(spec)
	if err != nil {
		return nil, err
	}

	return DownloadTorrent(t), nil
}
