package full

import (
	"strings"

	"github.com/smmr-software/mabel/internal/styles"
	"github.com/smmr-software/mabel/internal/torrent"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type addTorrent struct {
	width, height    int
	dir              bool
	torrent, saveDir textinput.Model
	theme            *styles.ColorTheme
	help             *help.Model
	main             *model
}

// Init starts ticking to refresh the UI without user interaction.
func (m addTorrent) Init() tea.Cmd {
	return tick()
}

// Update responds to messages by refreshing and resizing the view. It
// responds to several kinds of key presses in order to progress through
// adding the torrent.
func (m addTorrent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w := msg.Width - styles.BorderWindow.GetHorizontalBorderSize()
		h := msg.Height - styles.BorderWindow.GetHorizontalBorderSize()

		m.width = w
		m.main.width = w
		m.height = h
		m.main.height = h

		return m, nil
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC:
			m.main.client.Close()
			return m, tea.Quit
		case key.Matches(msg, addPromptKeys.quit):
			return m.main, nil
		case key.Matches(msg, addPromptKeys.forward):
			if m.dir {
				input := m.torrent.Value()
				dir := m.saveDir.Value()

				cmd, err := torrent.AddTorrent(&input, &dir, m.main.client, m.main.list, m.theme)
				if err != nil {
					cmd = reportError(err)
				}

				return m.main, cmd
			} else {
				m.torrent.Blur()
				m.saveDir.Focus()
				m.dir = true
				return m, nil
			}
		case key.Matches(msg, addPromptKeys.back):
			if m.dir {
				m.saveDir.Blur()
				m.torrent.Focus()
				m.dir = false
				return m, nil
			}
			return m.main, nil
		default:
			var cmd tea.Cmd
			if m.dir {
				m.saveDir, cmd = m.saveDir.Update(msg)
			} else {
				m.torrent, cmd = m.torrent.Update(msg)
			}
			return m, cmd
		}
	case tickMsg:
		return m, tick()
	}
	return m, nil
}

// View renders a text box for the torrent and a second for a save
// directory to facilitate adding a torrent to the client.
func (m addTorrent) View() string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(styles.BorderWindow)

	var body strings.Builder
	body.WriteString("Add Torrent\n")
	body.WriteString(styles.BorderWindow.Render(m.torrent.View()))
	body.WriteString("\n\nSave Directory (Optional)\n")
	body.WriteString(styles.BorderWindow.Render(m.saveDir.View()))

	help := m.help.View(addPromptKeys)
	height := m.height - gloss.Height(help) - 1

	content := gloss.Place(
		m.width, height,
		gloss.Center, gloss.Center,
		body.String(),
	)

	return fullscreen.Render(content + help + "\n")
}

// initialAddPrompt creates the initial state of the add prompt view,
// creating text boxes and accepting data from the main model.
func initialAddPrompt(w, h int, dir *string, theme *styles.ColorTheme, help *help.Model, parent *model) addTorrent {
	torrent := textinput.New()
	torrent.Width = 32
	saveDir := torrent

	saveDir.SetValue(*dir)
	saveDir.Blur()
	torrent.Focus()

	return addTorrent{
		width:  w,
		height: h,

		theme: theme,
		help:  help,

		torrent: torrent,
		saveDir: saveDir,

		main: parent,
	}
}
