package full

import (
	"fmt"
	"strings"
	"time"

	"github.com/smmr-software/mabel/internal/shared"

	"github.com/acarl005/stripansi"
	gloss "github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.err != nil {
		return errorView(&m)
	} else if m.portStartupFailure.enabled {
		return portStartupFailureView(&m)
	} else if m.addPrompt.enabled {
		return addPromptView(&m)
	} else if m.viewingTorrentDetails {
		return torrentDetailView(&m)
	} else {
		return mainView(&m)
	}
}

func portStartupFailureView(m *model) string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(borderWindow)

	var body strings.Builder
	body.WriteString(bold.Render("Port Binding Failure"))
	body.WriteString("\nplease provide an unused port number for the client to bind with\n\n")
	body.WriteString(borderWindow.Render(m.portStartupFailure.port.View()))

	return fullscreen.Render(
		gloss.Place(
			m.width, m.height,
			gloss.Center, gloss.Center,
			body.String(),
		),
	)
}

func addPromptView(m *model) string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(borderWindow)

	var body strings.Builder
	body.WriteString("Add Torrent\n")
	body.WriteString(borderWindow.Render(m.addPrompt.torrent.View()))
	body.WriteString("\n\nSave Directory (Optional)\n")
	body.WriteString(borderWindow.Render(m.addPrompt.saveDir.View()))

	help := m.help.View(addPromptKeys)
	height := m.height - gloss.Height(help) - 1

	content := gloss.Place(
		m.width, height,
		gloss.Center, gloss.Center,
		body.String(),
	)

	return fullscreen.Render(content + help + "\n")
}

func errorView(m *model) string {
	popupWidth := m.width / 3
	popupHeight := m.height / 4
	padding := m.height / 16

	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(borderWindow).
		BorderForeground(errorRed)
	popupWindow := gloss.NewStyle().
		Width(popupWidth).
		Height(popupHeight).
		Padding(0, padding).
		Inherit(borderWindow).
		BorderForeground(errorRed)
	header := gloss.NewStyle().Bold(true)

	popup := popupWindow.Render(gloss.Place(
		popupWidth-padding*2, popupHeight,
		gloss.Center, gloss.Center,
		header.Render("Error")+"\n"+m.err.Error(),
	))

	help := tooltip.Render("press any key to return home")
	height := m.height - gloss.Height(help) - 1

	content := gloss.Place(
		m.width, height,
		gloss.Center, gloss.Center,
		popup,
	)

	return fullscreen.Render(content + help + "\n")
}

func torrentDetailView(m *model) string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(borderWindow)

	selected := m.list.SelectedItem().(Item)
	t := selected.Self

	info := t.Info()
	stats := t.Stats()
	files := t.Files()

	done := t.BytesCompleted()
	upload := stats.BytesWritten.Int64()

	icon := ""
	if info.IsDir() {
		icon = ""
	}

	filesDesc := "file"
	if len(files) > 1 {
		filesDesc += "s"
	}

	ratioDesc := "N/A"
	if done != 0 {
		ratioDesc = fmt.Sprintf("%0.2f", float64(upload)/float64(done))
	}

	created := ""
	if tm := selected.Created; !tm.IsZero() && tm != time.Now() {
		created = fmt.Sprintf("\n\nCreated on %s", tm.Format("02 Jan 2006"))
	}
	with := "\n\n"
	if prog := selected.Program; prog != "" && prog != "go.torrent" {
		with = fmt.Sprintf(" with %s\n\n", prog)
	}
	comment := ""
	if com := selected.Comment; com != "" && com != "dynamic metainfo from client" {
		comment = fmt.Sprintf("%s\n\n", com)
	}

	var body strings.Builder
	body.WriteString(bold.Render(stripansi.Strip(t.Name())) + "\n\n")
	body.WriteString(fmt.Sprintf("%s%s%s", created, with, comment))
	body.WriteString(shared.ProgressBar(t, nil))
	body.WriteString(
		fmt.Sprintf(
			"\n\n%s  %d %s | %s\n\n",
			icon, len(files), filesDesc,
			shared.PeerStats(t),
		),
	)
	body.WriteString(
		fmt.Sprintf(
			"%s | %s | %s ratio\n\n",
			shared.DownloadStats(t, true),
			shared.UploadStats(t),
			ratioDesc,
		),
	)
	if len(files) > 1 {
		body.WriteString(
			fmt.Sprintf(
				"\n\nContent\n%s",
				fileView(&files, &m.width, &m.height),
			),
		)
	}

	help := tooltip.Render("press any key to return home")
	height := m.height - gloss.Height(help) - 1

	content := gloss.Place(
		m.width, height,
		gloss.Center, gloss.Center,
		body.String(),
	)

	return fullscreen.Render(content + help + "\n")
}

func mainView(m *model) string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(borderWindow)

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
