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
		magnet  bool
	}
}
