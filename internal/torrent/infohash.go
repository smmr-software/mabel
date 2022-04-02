package torrent

import (
	"strings"
	"time"

	"github.com/smmr-software/mabel/internal/list"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"

	clist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func addInfoHash(input *string, dir *storage.ClientImpl, client *torrent.Client, l *clist.Model) (tea.Cmd, error) {
	if strings.HasPrefix(*input, infohashPrefix) {
		*input = strings.TrimPrefix(*input, infohashPrefix)
	}

	hash := metainfo.NewHashFromHex(*input)
	t, nw := client.AddTorrentInfoHashWithStorage(hash, *dir)
	if l != nil && nw {
		l.SetItems(
			append(
				l.Items(),
				list.Item{Self: t, Added: time.Now()},
			),
		)
	}

	return downloadTorrent(t), nil
}
