package main

import (
	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/charmbracelet/bubbles/textinput"
)

type model struct {
	width, height int
	client        *torrent.Client
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
	config := torrent.NewDefaultClientConfig()
	config.DefaultStorage = storage.NewMMap("")
	config.Logger = log.Discard

	client, _ := torrent.NewClient(config)

	m := model{
		addPrompt: initialAddPrompt(),
		client:    client,
	}
	return m
}
