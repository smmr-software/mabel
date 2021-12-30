package main

import (
	"os"
	"strings"
	"time"

	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
)

type model struct {
	width, height         int
	client                *torrent.Client
	torrentMeta           map[metainfo.Hash]time.Time
	selected              metainfo.Hash
	help                  help.Model
	err                   error
	addPrompt             modelAddPrompt
	viewingTorrentDetails bool
}

type modelAddPrompt struct {
	enabled bool
	dir     bool
	torrent textinput.Model
	saveDir textinput.Model
}

func initialAddPrompt() modelAddPrompt {
	input := textinput.NewModel()
	input.Width = 32

	s := modelAddPrompt{
		enabled: false,
		dir:     false,
		torrent: input,
		saveDir: input,
	}
	return s
}

func iteratePorts(conf *torrent.ClientConfig) *torrent.Client {
	for search := 0; search < 5; search++ {
		conf.ListenPort += 1
		if client, err := torrent.NewClient(conf); err == nil {
			return client
		}
	}
	return nil
}

func initialModel() (model, error) {
	config := torrent.NewDefaultClientConfig()
	config.Logger = log.Discard
	config.Seed = true

	metadataDirectory := os.TempDir()
	if metadataStorage, err := storage.NewDefaultPieceCompletionForDir(metadataDirectory); err != nil {
		config.DefaultStorage = storage.NewMMap("")
	} else {
		config.DefaultStorage = storage.NewMMapWithCompletion("", metadataStorage)
	}

	var client *torrent.Client
	if cl, err := torrent.NewClient(config); err != nil {
		if strings.HasPrefix(err.Error(), "subsequent listen") {
			cl = iteratePorts(config)
			if cl != nil {
				client = cl
			} else {
				return model{}, err
			}
		} else {
			return model{}, err
		}
	}
	m := model{
		client:      client,
		torrentMeta: make(map[metainfo.Hash]time.Time),
		help:        help.NewModel(),
		addPrompt:   initialAddPrompt(),
	}
	return m, nil
}
