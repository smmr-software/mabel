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

func (m errorScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - styles.BorderWindow.GetHorizontalBorderSize()
		m.height = msg.Height - styles.BorderWindow.GetHorizontalBorderSize()
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
