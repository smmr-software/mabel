package main

import (
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/textinput"
)

type modelTorrent struct {
	client torrent.Client
	config torrent.ClientConfig
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

type model struct {
	width, height int
	torrent       modelTorrent
	addPrompt     modelAddPrompt
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
	m := model{
		addPrompt: initialAddPrompt(),
	}
	return m
}
