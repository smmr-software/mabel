package full

import "github.com/charmbracelet/bubbles/key"

type CustomKeyMap struct {
	Home struct {
		Up, Down, Next,
		Prev, Add, Delete,
		Details, Deselect,
		Help, Quit string
	}
	AddPrompt struct {
		Quit, Back,
		Forward string
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
	if c.Home.Up != "" {
		homeKeys.up = key.NewBinding(
			key.WithKeys(c.Home.Up),
			key.WithHelp(c.Home.Up, "up"),
		)
	}
	if c.Home.Down != "" {
		homeKeys.down = key.NewBinding(
			key.WithKeys(c.Home.Down),
			key.WithHelp(c.Home.Down, "down"),
		)
	}
	if c.Home.Next != "" {
		homeKeys.next = key.NewBinding(
			key.WithKeys(c.Home.Next),
			key.WithHelp(c.Home.Next, "next page"),
		)
	}
	if c.Home.Prev != "" {
		homeKeys.prev = key.NewBinding(
			key.WithKeys(c.Home.Prev),
			key.WithHelp(c.Home.Prev, "prev page"),
		)
	}
	if c.Home.Add != "" {
		homeKeys.add = key.NewBinding(
			key.WithKeys(c.Home.Add),
			key.WithHelp(c.Home.Add, "add torrent"),
		)
	}
	if c.Home.Delete != "" {
		homeKeys.delete = key.NewBinding(
			key.WithKeys(c.Home.Delete),
			key.WithHelp(c.Home.Delete, "delete"),
		)
	}
	if c.Home.Deselect != "" {
		homeKeys.deselect = key.NewBinding(
			key.WithKeys(c.Home.Deselect),
			key.WithHelp(c.Home.Deselect, "deselect"),
		)
	}
	if c.Home.Details != "" {
		homeKeys.details = key.NewBinding(
			key.WithKeys(c.Home.Details),
			key.WithHelp(c.Home.Details, "torrent details"),
		)
	}
	if c.Home.Help != "" {
		homeKeys.help = key.NewBinding(
			key.WithKeys(c.Home.Help),
			key.WithHelp(c.Home.Help, "toggle help"),
		)
	}
	if c.Home.Quit != "" {
		homeKeys.quit = key.NewBinding(
			key.WithKeys(c.Home.Quit),
			key.WithHelp(c.Home.Quit, "quit"),
		)
	}
	if c.AddPrompt.Quit != "" {
		addPromptKeys.quit = key.NewBinding(
			key.WithKeys(c.AddPrompt.Quit),
			key.WithHelp(c.AddPrompt.Quit, "home"),
		)
	}
	if c.AddPrompt.Back != "" {
		addPromptKeys.back = key.NewBinding(
			key.WithKeys(c.AddPrompt.Back),
			key.WithHelp(c.AddPrompt.Back, "home"),
		)
	}
	if c.AddPrompt.Forward != "" {
		addPromptKeys.forward = key.NewBinding(
			key.WithKeys(c.AddPrompt.Forward),
			key.WithHelp(c.AddPrompt.Forward, "home"),
		)
	}
}
