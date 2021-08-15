package main

import (
	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/charmbracelet/bubbles/textinput"
)

type model struct {
	width, height int
	typing        bool
	torrent       *torrent.Client
	addPrompt     modelAddPrompt
}

type modelAddPrompt struct {
	enabled bool
	phase   string
	state   modelAddPromptState
}

type modelAddPromptState struct {
	download    string
	magnetLink  textinput.Model
	torrentPath textinput.Model
	savePath    textinput.Model
}

func initialAddPrompt() modelAddPrompt {
	magnetLink := textinput.NewModel()
	magnetLink.Placeholder = "Magnet Link"
	magnetLink.Width = 32
	torrentPath := textinput.NewModel()
	torrentPath.Placeholder = "Path to Torrent File"
	torrentPath.Width = 32
	savePath := textinput.NewModel()
	savePath.Placeholder = "Path to Save Directory"
	savePath.Width = 32

	s := modelAddPrompt{
		enabled: false,
		phase:   "tab-select",
		state: modelAddPromptState{
			download:    "magnet",
			magnetLink:  magnetLink,
			torrentPath: torrentPath,
			savePath:    savePath,
		},
	}
	return s
}

func initialModel() model {
	config := torrent.NewDefaultClientConfig()
	config.DefaultStorage = storage.NewMMap("")
	config.Logger = log.Discard

	client, _ := torrent.NewClient(config)
	defer client.Close()

	m := model{
		addPrompt: initialAddPrompt(),
		torrent:   client,
	}
	return m
}
