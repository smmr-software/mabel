package full

import (
	"log"

	"github.com/smmr-software/mabel/internal/torrent"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		torrent.AddTorrents(
			m.startupTorrents,
			m.dir, m.client, m.list,
		),
		tick(),
	)
}

func Execute(torrents *[]string, dir *string, port *uint) {
	model, err := initialModel(torrents, dir, port)
	if err != nil {
		log.Fatal(err)
	}

	if err := tea.NewProgram(model, tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
