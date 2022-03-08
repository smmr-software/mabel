package main

import "github.com/charmbracelet/bubbles/key"

type homeKeyMap struct {
	up       key.Binding
	down     key.Binding
	next     key.Binding
	prev     key.Binding
	add      key.Binding
	delete   key.Binding
	details  key.Binding
	deselect key.Binding
	help     key.Binding
	quit     key.Binding
}

func (k homeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

func (k homeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down},
		{k.next, k.prev},
		{k.add, k.delete},
		{k.details, k.deselect},
		{k.help, k.quit},
	}
}

var homeKeys = homeKeyMap{
	up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "up"),
	),
	down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "down"),
	),
	next: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l/→", "next page"),
	),
	prev: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h/←", "prev page"),
	),
	add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add torrent"),
	),
	delete: key.NewBinding(
		key.WithKeys("d", "backspace"),
		key.WithHelp("d/⌦ ", "delete"),
	),
	deselect: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("⎋", "deselect"),
	),
	details: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "torrent details"),
	),
	help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

type addPromptKeyMap struct {
	quit    key.Binding
	back    key.Binding
	forward key.Binding
}

func (k addPromptKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.forward, k.back, k.quit}
}

func (k addPromptKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.back},
		{k.forward},
		{k.quit},
	}
}

var addPromptKeys = addPromptKeyMap{
	quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("⎋", "home"),
	),
	back: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("⇧ ↹", "previous"),
	),
	forward: key.NewBinding(
		key.WithKeys("enter", "tab"),
		key.WithHelp("↵", "next"),
	),
}
