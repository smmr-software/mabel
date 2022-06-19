package full

import (
	"fmt"
	"strings"
	"time"

	"github.com/smmr-software/mabel/internal/list"
	"github.com/smmr-software/mabel/internal/stats"
	"github.com/smmr-software/mabel/internal/styles"

	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"

	"github.com/acarl005/stripansi"
)

type torrentDetails struct {
	width, height int
	item          *list.Item
	theme         *styles.ColorTheme
	main          *model
}

// Init starts ticking to refresh the UI without user interaction.
func (m torrentDetails) Init() tea.Cmd {
	return tick()
}

// Update responds to messages by refreshing and resizing the view. It
// responds to two kinds of key presses: the client quits on Ctrl+C,
// while every other key returns the user to the main view.
func (m torrentDetails) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - styles.BorderWindow.GetHorizontalBorderSize()
		m.height = msg.Height - styles.BorderWindow.GetHorizontalBorderSize()

		updated, _ := m.main.Update(msg)
		if mdl, ok := updated.(model); ok {
			m.main = &mdl
		}

		return m, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.main.client.Close()
			return m, tea.Quit
		}
		return *m.main, nil
	case tickMsg:
		return m, tick()
	}
	return m, nil
}

// View renders the details of a specific torrent, including download
// progress, metainfo, name, and files.
func (m torrentDetails) View() string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(styles.BorderWindow)

	t := m.item.Self

	info := t.Info()
	stts := t.Stats()
	files := t.Files()

	done := t.BytesCompleted()
	upload := stts.BytesWritten.Int64()

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
	if tm := m.item.Created; !tm.IsZero() && tm != time.Now() {
		created = fmt.Sprintf("\n\nCreated on %s", tm.Format("02 Jan 2006"))
	}
	with := "\n\n"
	if prog := m.item.Program; prog != "" && prog != "go.torrent" {
		with = fmt.Sprintf(" with %s\n\n", prog)
	}
	comment := ""
	if com := m.item.Comment; com != "" && com != "dynamic metainfo from client" {
		comment = fmt.Sprintf("%s\n\n", com)
	}

	var body strings.Builder
	body.WriteString(styles.Bold.Render(stripansi.Strip(t.Name())) + "\n\n")
	body.WriteString(fmt.Sprintf("%s%s%s", created, with, comment))
	body.WriteString(stats.ProgressBar(t, nil, m.theme))
	body.WriteString(
		fmt.Sprintf(
			"\n\n%s  %d %s | %s\n\n",
			icon, len(files), filesDesc,
			stats.Peers(t),
		),
	)
	body.WriteString(
		fmt.Sprintf(
			"%s | %s | %s ratio\n\n",
			stats.Download(t, true),
			stats.Upload(t),
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

	tooltip := gloss.NewStyle().Foreground(m.theme.Tooltip).Padding(0, 2)
	help := tooltip.Render("press any key to return home")
	height := m.height - gloss.Height(help) - 1

	content := gloss.Place(
		m.width, height,
		gloss.Center, gloss.Center,
		body.String(),
	)

	return fullscreen.Render(content + help + "\n")
}
