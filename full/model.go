package full

import (
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type model struct {
	width, height         int
	client                *torrent.Client
	clientConfig          *torrent.ClientConfig
	list                  *list.Model
	help                  *help.Model
	err                   error
	addPrompt             *modelAddPrompt
	viewingTorrentDetails bool
	portStartupFailure    *portStartupFailure
}

type modelAddPrompt struct {
	enabled bool
	dir     bool
	torrent textinput.Model
	saveDir textinput.Model
}

type portStartupFailure struct {
	enabled bool
	port    textinput.Model
}

func initialAddPrompt() *modelAddPrompt {
	torrent := textinput.New()
	torrent.Width = 32
	saveDir := torrent

	cache, _ := os.UserCacheDir()
	cache += "/mabel/lastDownloadDir"
	bytes, err := os.ReadFile(cache)
	if err != nil {
		saveDir.SetValue(xdg.UserDirs.Download)
	} else {
		saveDir.SetValue(string(bytes))
	}

	saveDir.Blur()

	s := modelAddPrompt{
		enabled: false,
		dir:     false,
		torrent: torrent,
		saveDir: saveDir,
	}
	return &s
}

func initialPortStartupFailure() *portStartupFailure {
	input := textinput.New()
	input.Width = 32
	input.Focus()

	port := portStartupFailure{port: input}
	return &port
}

func genMabelConfig() *torrent.ClientConfig {
	config := torrent.NewDefaultClientConfig()
	config.Logger = log.Default
	config.Logger.Handlers = []log.Handler{log.DiscardHandler}
	config.Seed = true

	metadataDirectory := os.TempDir()
	if metadataStorage, err := storage.NewDefaultPieceCompletionForDir(metadataDirectory); err != nil {
		config.DefaultStorage = storage.NewMMap("")
	} else {
		config.DefaultStorage = storage.NewMMapWithCompletion("", metadataStorage)
	}

	return config
}

func genList() *list.Model {
	list := list.New(make([]list.Item, 0), itemDelegate{}, 0, 0)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)
	return &list
}

func initialModel() (model, error) {
	config := genMabelConfig()
	client, err := torrent.NewClient(config)
	hlp := help.New()

	m := model{
		client:             client,
		clientConfig:       config,
		list:               genList(),
		help:               &hlp,
		addPrompt:          initialAddPrompt(),
		portStartupFailure: initialPortStartupFailure(),
	}

	if err != nil {
		msg := err.Error()
		switch {
		case strings.HasPrefix(msg, "subsequent listen"), strings.HasPrefix(msg, "first listen"):
			m.portStartupFailure.enabled = true
		default:
			return model{}, err
		}
	}
	return m, nil
}
