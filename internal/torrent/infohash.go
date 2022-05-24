package torrent

import (
	"strings"
	"time"

	"github.com/smmr-software/mabel/internal/list"
	"github.com/smmr-software/mabel/internal/styles"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"

	clist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// addInfoHash takes an infohash and adds the torrent to the client. If
// it is a new torrent, the torrent info is added to the Bubbles list.
func addInfoHash(input *string, dir *storage.ClientImpl, client *torrent.Client, l *clist.Model, theme *styles.ColorTheme) (tea.Cmd, error) {
	if strings.HasPrefix(*input, infohashPrefix) {
		*input = strings.TrimPrefix(*input, infohashPrefix)
	}

	hash := metainfo.NewHashFromHex(*input)
	t, nw := client.AddTorrentInfoHashWithStorage(hash, *dir)
	if l != nil && nw {
		l.SetItems(
			append(
				l.Items(),
				list.Item{
					Self:  t,
					Added: time.Now(),
					Theme: theme,
				},
			),
		)
	}

	return downloadTorrent(t), nil
}
