package full

import (
	"strconv"
	"strings"

	"github.com/smmr-software/mabel/internal/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"

	"github.com/anacrolix/torrent"
)

type portStartupFailure struct {
	width, height int
	input         textinput.Model
	main          *model
}

// Init starts ticking to refresh the UI without user interaction.
func (m portStartupFailure) Init() tea.Cmd {
	return tick()
}

// Update responds to messages by refreshing and resizing the view. It
// responds to three kinds of key presses: the client quits on Ctrl+C,
// numeric keys and backspace are passed to the text field, while other
// keys attempt to bind the client to the port provided.
func (m portStartupFailure) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - styles.Fullscreen.GetHorizontalBorderSize()
		m.height = msg.Height - styles.Fullscreen.GetHorizontalBorderSize()

		updated, _ := m.main.Update(msg)
		if mdl, ok := updated.(model); ok {
			m.main = &mdl
		}

		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "backspace":
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		default:
			prt, err := strconv.Atoi(m.input.Value())
			if err != nil {
				return m, reportError(err)
			}
			port := uint(prt)

			config := genMabelConfig(&port, m.main.logging, m.main.encrypt)
			client, err := torrent.NewClient(config)
			if err != nil {
				return m, reportError(err)
			}

			m.main.client = client
			m.main.clientConfig = config
			m.main.width = m.width
			m.main.height = m.height

			return m.main, nil
		}
	case tickMsg:
		return m, tick()
	default:
		return m, nil
	}
}

// View renders a text field for a user to provide an alternative port
// when port binding fails on startup.
func (m portStartupFailure) View() string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(styles.Fullscreen)

	var body strings.Builder
	body.WriteString(styles.Bold.Render("Port Binding Failure"))
	body.WriteString("\nplease provide an unused port number for the client to bind with\n\n")
	body.WriteString(styles.Window.Render(m.input.View()))

	return fullscreen.Render(
		gloss.Place(
			m.width, m.height,
			gloss.Center, gloss.Center,
			body.String(),
		),
	)
}

// initialPortStartupFailure accepts a pointer to the main model,
// creates a textbox, and returns the initial model state.
func initialPortStartupFailure(parent *model) portStartupFailure {
	input := textinput.New()
	input.Width = 32
	input.Focus()

	return portStartupFailure{input: input, main: parent}
}
