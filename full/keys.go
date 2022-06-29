package full

import "github.com/charmbracelet/bubbles/key"

type CustomKeyMap struct {
	home struct {
		up, down, next,
		prev, add, delete,
		details, deselect,
		help, quit string
	}
	addPrompt struct {
		quit, back,
		forward string
	} `toml:"add-prompt"`
}

type homeKeyMap struct {
	up, down, next,
	prev, add, delete,
	details, deselect,
	help, quit key.Binding
}

// ShortHelp returns the key bindings for the help and quit actions in
// the home screen.
func (k homeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

// FullHelp returns the key bindings for all actions in the home
// screen.
func (k homeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down},
		{k.next, k.prev},
		{k.add, k.delete},
		{k.details, k.deselect},
		{k.help, k.quit},
	}
}

// Define the key bindings and help symbols for each action in the home
// screen.
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
	quit, back, forward key.Binding
}

// ShortHelp returns the key bindings for the forward, back, and quit
// actions in the the add prompt screen.
func (k addPromptKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.forward, k.back, k.quit}
}

// FullHelp returns the key bindings for all actions in the the add
// prompt screen.
func (k addPromptKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.forward},
		{k.back},
		{k.quit},
	}
}

// Define the key bindings and help symbols for each action in the add
// prompt screen.
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

func (c CustomKeyMap) ToKeys() {
	if c.home.up != "" {
		homeKeys.up = key.NewBinding(
			key.WithKeys(c.home.up),
			key.WithHelp(c.home.up, "up"),
		)
	}
	if c.home.down != "" {
		homeKeys.down = key.NewBinding(
			key.WithKeys(c.home.down),
			key.WithHelp(c.home.down, "down"),
		)
	}
	if c.home.next != "" {
		homeKeys.next = key.NewBinding(
			key.WithKeys(c.home.next),
			key.WithHelp(c.home.next, "next page"),
		)
	}
	if c.home.prev != "" {
		homeKeys.prev = key.NewBinding(
			key.WithKeys(c.home.prev),
			key.WithHelp(c.home.prev, "prev page"),
		)
	}
	if c.home.add != "" {
		homeKeys.add = key.NewBinding(
			key.WithKeys(c.home.add),
			key.WithHelp(c.home.add, "add torrent"),
		)
	}
	if c.home.delete != "" {
		homeKeys.delete = key.NewBinding(
			key.WithKeys(c.home.delete),
			key.WithHelp(c.home.delete, "delete"),
		)
	}
	if c.home.deselect != "" {
		homeKeys.deselect = key.NewBinding(
			key.WithKeys(c.home.deselect),
			key.WithHelp(c.home.deselect, "deselect"),
		)
	}
	if c.home.details != "" {
		homeKeys.details = key.NewBinding(
			key.WithKeys(c.home.details),
			key.WithHelp(c.home.details, "torrent details"),
		)
	}
	if c.home.help != "" {
		homeKeys.help = key.NewBinding(
			key.WithKeys(c.home.help),
			key.WithHelp(c.home.help, "toggle help"),
		)
	}
	if c.home.quit != "" {
		homeKeys.quit = key.NewBinding(
			key.WithKeys(c.home.quit),
			key.WithHelp(c.home.quit, "quit"),
		)
	}
	if c.addPrompt.quit != "" {
		addPromptKeys.quit = key.NewBinding(
			key.WithKeys(c.addPrompt.quit),
			key.WithHelp(c.addPrompt.quit, "home"),
		)
	}
	if c.addPrompt.back != "" {
		addPromptKeys.back = key.NewBinding(
			key.WithKeys(c.addPrompt.back),
			key.WithHelp(c.addPrompt.back, "home"),
		)
	}
	if c.addPrompt.forward != "" {
		addPromptKeys.forward = key.NewBinding(
			key.WithKeys(c.addPrompt.forward),
			key.WithHelp(c.addPrompt.forward, "home"),
		)
	}
}
