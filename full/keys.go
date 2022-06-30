package full

import "github.com/charmbracelet/bubbles/key"

type CustomKeyMap struct {
	Home struct {
		Up, Down, Next,
		Prev, Add, Delete,
		Details, Deselect,
		Help, Quit CustomKey
	}
	AddPrompt struct {
		Quit, Prev,
		Next CustomKey
	} `toml:"add-prompt"`
}

type CustomKey struct {
	Key, Icon, Desc string
}

// isZero checks if a CustomKey object is of nil value.
func (k CustomKey) isZero() bool {
	return k.Key == "" && k.Icon == "" && k.Desc == ""
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
	quit, prev, next key.Binding
}

// ShortHelp returns the key bindings for the forward, back, and quit
// actions in the the add prompt screen.
func (k addPromptKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.next, k.prev, k.quit}
}

// FullHelp returns the key bindings for all actions in the the add
// prompt screen.
func (k addPromptKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.next},
		{k.prev},
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
	prev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("⇧ ↹", "previous"),
	),
	next: key.NewBinding(
		key.WithKeys("enter", "tab"),
		key.WithHelp("↵", "next"),
	),
}

// ToKeys changes the default homeKeys and addPromptKeys key bindings
// if there is user configuration.
func (c CustomKeyMap) ToKeys() {
	if !c.Home.Up.isZero() {
		homeKeys.up = key.NewBinding(
			key.WithKeys(c.Home.Up.Key),
			key.WithHelp(c.Home.Up.Icon, c.Home.Up.Desc),
		)
	}
	if !c.Home.Down.isZero() {
		homeKeys.down = key.NewBinding(
			key.WithKeys(c.Home.Down.Key),
			key.WithHelp(c.Home.Down.Icon, c.Home.Down.Desc),
		)
	}
	if !c.Home.Next.isZero() {
		homeKeys.next = key.NewBinding(
			key.WithKeys(c.Home.Next.Key),
			key.WithHelp(c.Home.Next.Icon, c.Home.Next.Desc),
		)
	}
	if !c.Home.Prev.isZero() {
		homeKeys.prev = key.NewBinding(
			key.WithKeys(c.Home.Prev.Key),
			key.WithHelp(c.Home.Prev.Icon, c.Home.Prev.Desc),
		)
	}
	if !c.Home.Add.isZero() {
		homeKeys.add = key.NewBinding(
			key.WithKeys(c.Home.Add.Key),
			key.WithHelp(c.Home.Add.Icon, c.Home.Add.Desc),
		)
	}
	if !c.Home.Delete.isZero() {
		homeKeys.delete = key.NewBinding(
			key.WithKeys(c.Home.Delete.Key),
			key.WithHelp(c.Home.Delete.Icon, c.Home.Delete.Desc),
		)
	}
	if !c.Home.Deselect.isZero() {
		homeKeys.deselect = key.NewBinding(
			key.WithKeys(c.Home.Deselect.Key),
			key.WithHelp(c.Home.Deselect.Icon, c.Home.Deselect.Desc),
		)
	}
	if !c.Home.Details.isZero() {
		homeKeys.details = key.NewBinding(
			key.WithKeys(c.Home.Details.Key),
			key.WithHelp(c.Home.Details.Icon, c.Home.Details.Desc),
		)
	}
	if !c.Home.Help.isZero() {
		homeKeys.help = key.NewBinding(
			key.WithKeys(c.Home.Help.Key),
			key.WithHelp(c.Home.Help.Icon, c.Home.Help.Desc),
		)
	}
	if !c.Home.Quit.isZero() {
		homeKeys.quit = key.NewBinding(
			key.WithKeys(c.Home.Quit.Key),
			key.WithHelp(c.Home.Quit.Icon, c.Home.Quit.Desc),
		)
	}
	if !c.AddPrompt.Quit.isZero() {
		addPromptKeys.quit = key.NewBinding(
			key.WithKeys(c.AddPrompt.Quit.Key),
			key.WithHelp(c.AddPrompt.Quit.Icon, c.AddPrompt.Quit.Desc),
		)
	}
	if !c.AddPrompt.Prev.isZero() {
		addPromptKeys.prev = key.NewBinding(
			key.WithKeys(c.AddPrompt.Prev.Key),
			key.WithHelp(c.AddPrompt.Prev.Icon, c.AddPrompt.Prev.Desc),
		)
	}
	if !c.AddPrompt.Next.isZero() {
		addPromptKeys.next = key.NewBinding(
			key.WithKeys(c.AddPrompt.Next.Key),
			key.WithHelp(c.AddPrompt.Next.Icon, c.AddPrompt.Next.Desc),
		)
	}
}
