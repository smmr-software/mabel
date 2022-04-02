package torrent

import (
	"time"

	"github.com/smmr-software/mabel/internal/list"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"

	clist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func addMagnetLink(input *string, dir *storage.ClientImpl, client *torrent.Client, l *clist.Model) (tea.Cmd, error) {
	spec, err := torrent.TorrentSpecFromMagnetUri(*input)
	if err != nil {
		return nil, err
	}
	spec.Storage = *dir

	t, nw, err := client.AddTorrentSpec(spec)
	if err != nil {
		return nil, err
	} else if l != nil && nw {
		l.SetItems(
			append(
				l.Items(),
				list.Item{Self: t, Added: time.Now()},
			),
		)
	}

	return downloadTorrent(t), nil
}
