package torrent

import (
	"time"

	"github.com/smmr-software/mabel/internal/list"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"

	"github.com/acarl005/stripansi"
	home "github.com/mitchellh/go-homedir"

	tea "github.com/charmbracelet/bubbletea"
)

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

	return downloadTorrent(t), nw,
		list.Item{
			Self:    t,
			Added:   time.Now(),
			Created: time.Unix(meta.CreationDate, 0),
			Comment: stripansi.Strip(meta.Comment),
			Program: stripansi.Strip(meta.CreatedBy),
		}, nil
}
