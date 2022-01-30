package main

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
	list                  list.Model
	help                  help.Model
	err                   error
	addPrompt             modelAddPrompt
	viewingTorrentDetails bool
	portStartupFailure    portStartupFailure
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

func initialAddPrompt() modelAddPrompt {
	torrent := textinput.New()
	torrent.Width = 32
	saveDir := torrent
	saveDir.SetValue(xdg.UserDirs.Download)
	saveDir.Blur()

	s := modelAddPrompt{
		enabled: false,
		dir:     false,
		torrent: torrent,
		saveDir: saveDir,
	}
	return s
}

func initialPortStartupFailure() portStartupFailure {
	input := textinput.New()
	input.Width = 32
	input.Focus()

	return portStartupFailure{port: input}
}

func genMabelConfig() *torrent.ClientConfig {
	config := torrent.NewDefaultClientConfig()
	config.Logger = log.Discard
	config.Seed = true

	metadataDirectory := os.TempDir()
	if metadataStorage, err := storage.NewDefaultPieceCompletionForDir(metadataDirectory); err != nil {
		config.DefaultStorage = storage.NewMMap("")
	} else {
		config.DefaultStorage = storage.NewMMapWithCompletion("", metadataStorage)
	}

	return config
}

func genList() list.Model {
	list := list.NewModel(make([]list.Item, 0), itemDelegate{}, 0, 0)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)
	return list
}

func initialModel() (model, error) {
	config := genMabelConfig()
	client, err := torrent.NewClient(config)

	m := model{
		client:             client,
		clientConfig:       config,
		list:               genList(),
		help:               help.New(),
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
