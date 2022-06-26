package full

import "github.com/charmbracelet/bubbles/key"

type CustomKeyMap struct {
	customHomeKeyMap struct {
		up       string
		down     string
		next     string
		prev     string
		add      string
		delete   string
		details  string
		deselect string
		help     string
		quit     string
	}
	customAddPromptKeyMap struct {
		quit    string
		back    string
		forward string
	}
}

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
	quit    key.Binding
	back    key.Binding
	forward key.Binding
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

func (c *CustomKeyMap) ToKeys() {
	if c.customHomeKeyMap.up != "" {
		homeKeys.up = key.NewBinding(
			key.WithKeys(c.customHomeKeyMap.up),
			key.WithHelp(c.customHomeKeyMap.up, "up"),
		)
	}
	if c.customHomeKeyMap.down != "" {
		homeKeys.down = key.NewBinding(
			key.WithKeys(c.customHomeKeyMap.down),
			key.WithHelp(c.customHomeKeyMap.down, "down"),
		)
	}
	if c.customHomeKeyMap.next != "" {
		homeKeys.next = key.NewBinding(
			key.WithKeys(c.customHomeKeyMap.next),
			key.WithHelp(c.customHomeKeyMap.next, "next page"),
		)
	}
	if c.customHomeKeyMap.prev != "" {
		homeKeys.prev = key.NewBinding(
			key.WithKeys(c.customHomeKeyMap.prev),
			key.WithHelp(c.customHomeKeyMap.prev, "prev page"),
		)
	}
	if c.customHomeKeyMap.add != "" {
		homeKeys.add = key.NewBinding(
			key.WithKeys(c.customHomeKeyMap.add),
			key.WithHelp(c.customHomeKeyMap.add, "add torrent"),
		)
	}
	if c.customHomeKeyMap.delete != "" {
		homeKeys.delete = key.NewBinding(
			key.WithKeys(c.customHomeKeyMap.delete),
			key.WithHelp(c.customHomeKeyMap.delete, "delete"),
		)
	}
	if c.customHomeKeyMap.deselect != "" {
		homeKeys.deselect = key.NewBinding(
			key.WithKeys(c.customHomeKeyMap.deselect),
			key.WithHelp(c.customHomeKeyMap.deselect, "deselect"),
		)
	}
	if c.customHomeKeyMap.details != "" {
		homeKeys.details = key.NewBinding(
			key.WithKeys(c.customHomeKeyMap.details),
			key.WithHelp(c.customHomeKeyMap.details, "torrent details"),
		)
	}
	if c.customHomeKeyMap.help != "" {
		homeKeys.help = key.NewBinding(
			key.WithKeys(c.customHomeKeyMap.help),
			key.WithHelp(c.customHomeKeyMap.help, "toggle help"),
		)
	}
	if c.customHomeKeyMap.quit != "" {
		homeKeys.quit = key.NewBinding(
			key.WithKeys(c.customHomeKeyMap.quit),
			key.WithHelp(c.customHomeKeyMap.quit, "quit"),
		)
	}
	if c.customAddPromptKeyMap.quit != "" {
		addPromptKeys.quit = key.NewBinding(
			key.WithKeys(c.customAddPromptKeyMap.quit),
			key.WithHelp(c.customAddPromptKeyMap.quit, "home"),
		)
	}
	if c.customAddPromptKeyMap.back != "" {
		addPromptKeys.back = key.NewBinding(
			key.WithKeys(c.customAddPromptKeyMap.back),
			key.WithHelp(c.customAddPromptKeyMap.back, "home"),
		)
	}
	if c.customAddPromptKeyMap.forward != "" {
		addPromptKeys.forward = key.NewBinding(
			key.WithKeys(c.customAddPromptKeyMap.forward),
			key.WithHelp(c.customAddPromptKeyMap.forward, "home"),
		)
	}
}
