package full

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	torrent "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	tea "github.com/charmbracelet/bubbletea"
	home "github.com/mitchellh/go-homedir"
)

func addTorrent(m *model) (tea.Model, tea.Cmd) {
	input := m.addPrompt.torrent.Value()
	saveDir := getSaveLocation(m.addPrompt.saveDir.Value())
	storageDir := getStorage(&saveDir)

	if strings.HasPrefix(input, "magnet:") {
		return addMagnetLink(m, &input, &storageDir)
	} else if strings.HasPrefix(input, "infohash:") {
		return addInfoHash(m, &input, &storageDir)
	} else {
		return addFromFile(m, &input, &storageDir)
	}
}

func getSaveLocation(dir string) string {
	if d, err := home.Expand(dir); err != nil {
		return ""
	} else {
		cacheSaveDir(d)
		return d
	}
}

func cacheSaveDir(dir string) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return
	}
	cache += "/mabel"

	err = os.Mkdir(cache, os.ModePerm)
	if err != nil && err.Error() != fmt.Sprintf("mkdir %s: file exists", cache) {
		return
	}

	file, err := os.Create(cache + "/lastDownloadDir")
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(dir)
}

func getStorage(dir *string) storage.ClientImpl {
	metadataDirectory := os.TempDir()
	if metadataStorage, err := storage.NewDefaultPieceCompletionForDir(metadataDirectory); err != nil {
		return storage.NewMMap(*dir)
	} else {
		return storage.NewMMapWithCompletion(*dir, metadataStorage)
	}
}

func addMagnetLink(m *model, input *string, dir *storage.ClientImpl) (tea.Model, tea.Cmd) {
	spec, err := torrent.TorrentSpecFromMagnetUri(*input)
	if err != nil {
		m.addPrompt = initialAddPrompt()
		return m, reportError(err)
	}
	spec.Storage = *dir

	t, nw, err := m.client.AddTorrentSpec(spec)
	if err != nil {
		m.addPrompt = initialAddPrompt()
		return m, reportError(err)
	}

	var cmd tea.Cmd
	if nw {
		m.list.SetItems(
			append(
				m.list.Items(),
				item{
					self:  t,
					added: time.Now(),
				},
			),
		)
		cmd = downloadTorrent(t)
	}
	m.addPrompt = initialAddPrompt()
	return m, cmd
}

func addInfoHash(m *model, input *string, dir *storage.ClientImpl) (tea.Model, tea.Cmd) {
	hash := metainfo.NewHashFromHex(strings.TrimPrefix(*input, "infohash:"))
	t, nw := m.client.AddTorrentInfoHashWithStorage(hash, *dir)
	var cmd tea.Cmd
	if nw {
		m.list.SetItems(
			append(
				m.list.Items(),
				item{
					self:  t,
					added: time.Now(),
				},
			),
		)
		cmd = downloadTorrent(t)
	}
	m.addPrompt = initialAddPrompt()
	return m, cmd
}

func addFromFile(m *model, input *string, dir *storage.ClientImpl) (tea.Model, tea.Cmd) {
	path, err := home.Expand(*input)
	if err != nil {
		m.addPrompt = initialAddPrompt()
		return m, reportError(err)
	}

	meta, err := metainfo.LoadFromFile(path)
	if err != nil {
		m.addPrompt = initialAddPrompt()
		return m, reportError(err)
	}

	spec := torrent.TorrentSpecFromMetaInfo(meta)
	spec.Storage = *dir

	t, nw, err := m.client.AddTorrentSpec(spec)
	if err != nil {
		m.addPrompt = initialAddPrompt()
		return m, reportError(err)
	}

	var cmd tea.Cmd
	if nw {
		m.list.SetItems(
			append(
				m.list.Items(),
				item{
					self:    t,
					added:   time.Now(),
					created: time.Unix(meta.CreationDate, 0),
					comment: stripansi.Strip(meta.Comment),
					program: stripansi.Strip(meta.CreatedBy),
				},
			),
		)
		cmd = downloadTorrent(t)
	}
	m.addPrompt = initialAddPrompt()
	return m, cmd
}
