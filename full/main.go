// Package full generates the full client for multiple torrents.
package full

import (
	"fmt"
	"os"

	"github.com/smmr-software/mabel/internal/styles"
	"github.com/smmr-software/mabel/internal/torrent"

	tea "github.com/charmbracelet/bubbletea"
)

// Init starts the UI and adds startup torrents.
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

// Execute creates the initial model and a Bubble Tea program, and
// aborts the client if that fails.
func Execute(torrents *[]string, dir *string, port *uint, logging, encrypt *bool, theme *styles.ColorTheme, keys CustomKeyMap) {
	keys.ToKeys()
	model, err := initialModel(torrents, dir, port, logging, encrypt, theme)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
