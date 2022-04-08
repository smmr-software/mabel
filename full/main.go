package full

import (
	"log"

	"github.com/smmr-software/mabel/internal/styles"
	"github.com/smmr-software/mabel/internal/torrent"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		torrent.AddTorrents(
			m.startupTorrents,
			m.dir, m.client,
			m.list, m.theme,
		),
		tick(),
	)
}

func Execute(torrents *[]string, dir *string, port *uint, theme *styles.ColorTheme) {
	model, err := initialModel(torrents, dir, port, theme)
	if err != nil {
		log.Fatal(err)
	}

	if err := tea.NewProgram(model, tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
