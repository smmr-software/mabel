package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	quit       key.Binding
	addTorrent key.Binding
}

var keys = keyMap{
	quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	addTorrent: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add torrent"),
	),
}
