package torrent

import (
	"time"

	"github.com/smmr-software/mabel/internal/list"
	"github.com/smmr-software/mabel/internal/styles"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"

	"github.com/acarl005/stripansi"
	home "github.com/mitchellh/go-homedir"

	clist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// addFromFile takes a metainfo (.torrent) file adds the torrent to the
// client. If it is a new torrent, the torrent info is added to the
// Bubbles list.
func addFromFile(input *string, dir *storage.ClientImpl, client *torrent.Client, l *clist.Model, theme *styles.ColorTheme) (tea.Cmd, error) {
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

	t, nw, err := client.AddTorrentSpec(spec)
	if err != nil {
		return nil, err
	} else if l != nil && nw {
		l.SetItems(
			append(
				l.Items(),
				list.Item{
					Self:    t,
					Added:   time.Now(),
					Created: time.Unix(meta.CreationDate, 0),
					Comment: stripansi.Strip(meta.Comment),
					Program: stripansi.Strip(meta.CreatedBy),
					Theme:   theme,
				},
			),
		)
	}

	return downloadTorrent(t), nil
}
