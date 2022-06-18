package full

import (
	"github.com/smmr-software/mabel/internal/styles"

	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type errorScreen struct {
	width, height int
	err           error
	theme         *styles.ColorTheme
	main          *model
}

func (m errorScreen) Init() tea.Cmd {
	return tick()
}

// Update responds to messages by refreshing and resizing the view. It
// responds to two kinds of key presses: the client quits on Ctrl+C,
// while every other key returns the user to the main view.
func (m errorScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		if msg.Type == tea.KeyCtrlC {
			m.main.client.Close()
			return m, tea.Quit
		}
		return m.main, nil
	case tickMsg:
		return m, tick()
	default:
		return m, nil
	}
}

// View presents an error to the user in a centered modal
func (m errorScreen) View() string {
	popupWidth := m.width / 3
	popupHeight := m.height / 4
	padding := m.height / 16

	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(styles.BorderWindow).
		BorderForeground(m.theme.Error)
	popupWindow := gloss.NewStyle().
		Width(popupWidth).
		Height(popupHeight).
		Padding(0, padding).
		Inherit(styles.BorderWindow).
		BorderForeground(m.theme.Error)
	header := gloss.NewStyle().Bold(true)

	popup := popupWindow.Render(gloss.Place(
		popupWidth-padding*2, popupHeight,
		gloss.Center, gloss.Center,
		header.Render("Error")+"\n"+m.err.Error(),
	))

	tooltip := gloss.NewStyle().Foreground(m.theme.Tooltip).Padding(0, 2)
	help := tooltip.Render("press any key to return home")
	height := m.height - gloss.Height(help) - 1

	content := gloss.Place(
		m.width, height,
		gloss.Center, gloss.Center,
		popup,
	)

	return fullscreen.Render(content + help + "\n")
}
