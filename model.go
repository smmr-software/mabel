package main

import "github.com/anacrolix/torrent"

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
	download string
}

type model struct {
	width     int
	height    int
	torrent   modelTorrent
	addPrompt modelAddPrompt
}

func initialAddPrompt() modelAddPrompt {
	s := modelAddPrompt{
		enabled: false,
		phase:   "tab-select",
		state: modelAddPromptState{
			download: "magnet",
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
