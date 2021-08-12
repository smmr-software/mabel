package main

import (
	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/textinput"
	"os"
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
	torrentPath := textinput.NewModel()
	torrentPath.Placeholder = "Path to Torrent File"
	savePath := textinput.NewModel()
	savePath.Placeholder = "Path to Save Directory"

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
	config.DataDir = os.TempDir()
	config.Logger = log.Discard

	client, _ := torrent.NewClient(config)
	defer client.Close()

	m := model{
		addPrompt: initialAddPrompt(),
		torrent:   client,
	}
	return m
}
