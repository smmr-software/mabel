package torrent

import (
	"time"

	"github.com/smmr-software/mabel/internal/list"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"

	tea "github.com/charmbracelet/bubbletea"
)

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

	return downloadTorrent(t), nw, list.Item{Self: t, Added: time.Now()}, nil
}
