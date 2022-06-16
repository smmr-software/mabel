package full

import (
	"strings"

	"github.com/smmr-software/mabel/internal/styles"

	gloss "github.com/charmbracelet/lipgloss"
)

// View returns the UI model in its current state as a string to be
// displayed to the user.
func (m model) View() string {
	if m.err != nil {
		return errorView(&m)
	} else if m.portStartupFailure.enabled {
		return portStartupFailureView(&m)
	} else if m.addPrompt.enabled {
		return addPromptView(&m)
	} else {
		return mainView(&m)
	}
}

// portStartupFailureView renders the screen when the port binding
// fails and the user must provide a different port.
func portStartupFailureView(m *model) string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(styles.BorderWindow)

	var body strings.Builder
	body.WriteString(styles.Bold.Render("Port Binding Failure"))
	body.WriteString("\nplease provide an unused port number for the client to bind with\n\n")
	body.WriteString(styles.BorderWindow.Render(m.portStartupFailure.port.View()))

	return fullscreen.Render(
		gloss.Place(
			m.width, m.height,
			gloss.Center, gloss.Center,
			body.String(),
		),
	)
}

// addPromptView renders the screen when a new torrent is being added.
func addPromptView(m *model) string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(styles.BorderWindow)

	var body strings.Builder
	body.WriteString("Add Torrent\n")
	body.WriteString(styles.BorderWindow.Render(m.addPrompt.torrent.View()))
	body.WriteString("\n\nSave Directory (Optional)\n")
	body.WriteString(styles.BorderWindow.Render(m.addPrompt.saveDir.View()))

	help := m.help.View(addPromptKeys)
	height := m.height - gloss.Height(help) - 1

	content := gloss.Place(
		m.width, height,
		gloss.Center, gloss.Center,
		body.String(),
	)

	return fullscreen.Render(content + help + "\n")
}

// errorView renders the screen when an error occurs, presenting the
// error and allowing the user to return home.
func errorView(m *model) string {
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

// mainView renders the main screen with the torrent list.
func mainView(m *model) string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(styles.BorderWindow)

	var content string
	if torrents := m.client.Torrents(); len(torrents) > 0 {
		content = m.list.View()
	} else {
		content = "You have no torrents!"
	}

	help := m.help.View(homeKeys)
	height := m.height - gloss.Height(help) - 1

	content = gloss.Place(
		m.width, height,
		gloss.Center, gloss.Center,
		content,
	)

	return fullscreen.Render(content + help + "\n")
}
