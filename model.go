package main

import (
	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"os"
)

type model struct {
	width, height int
	client        *torrent.Client
	help          help.Model
	addPrompt     modelAddPrompt
}

type modelAddPrompt struct {
	enabled bool
	input   textinput.Model
}

func initialAddPrompt() modelAddPrompt {
	input := textinput.NewModel()
	input.Width = 32

	s := modelAddPrompt{
		enabled: false,
		input:   input,
	}
	return s
}

func initialModel() model {
	metadataDirectory := os.TempDir()
	metadataStorage, _ := storage.NewDefaultPieceCompletionForDir(metadataDirectory)

	config := torrent.NewDefaultClientConfig()
	config.DefaultStorage = storage.NewMMapWithCompletion("", metadataStorage)
	config.Logger = log.Discard

	client, _ := torrent.NewClient(config)

	m := model{
		client:    client,
		help:      help.NewModel(),
		addPrompt: initialAddPrompt(),
	}
	return m
}
