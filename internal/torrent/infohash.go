package torrent

import (
	"strings"
	"time"

	"github.com/smmr-software/mabel/internal/list"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func addInfoHash(input *string, dir *storage.ClientImpl, client *torrent.Client) (tea.Cmd, bool, list.Item, error) {
	if strings.HasPrefix(*input, infohashPrefix) {
		*input = strings.TrimPrefix(*input, infohashPrefix)
	}

	hash := metainfo.NewHashFromHex(*input)
	t, nw := client.AddTorrentInfoHashWithStorage(hash, *dir)

	return downloadTorrent(t), nw, list.Item{Self: t, Added: time.Now()}, nil
}
