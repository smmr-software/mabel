package main

import "github.com/anacrolix/torrent"

type model struct {
	width   int
	height  int
	torrent struct {
		client torrent.Client
		config torrent.ClientConfig
	}
	addPrompt struct {
		enabled bool
		phase   string
		state   struct {
			download string
		}
	}
}

func initialAddPrompt() struct {
	enabled bool
	phase   string
	state   struct {
		download string
	}
} {
	s := struct {
		enabled bool
		phase   string
		state   struct {
			download string
		}
	}{
		enabled: false,
		phase:   "tab-select",
		state: struct {
			download string
		}{
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
