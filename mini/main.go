// Package mini generates the mini client for downloading single
// torrents.
package mini

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/smmr-software/mabel/internal/stats"
	"github.com/smmr-software/mabel/internal/styles"
	trrnt "github.com/smmr-software/mabel/internal/torrent"
	"github.com/smmr-software/mabel/internal/utils"

	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"

	"github.com/adrg/xdg"
)

const interval = 500 * time.Millisecond

type tickMsg time.Time
type model struct {
	width            int
	torrent, saveDir *string
	theme            *styles.ColorTheme
	client           *torrent.Client
}

// genMabelConfig configures the torrent client (seeding, listening
// port, log directory, etc.)
func genMabelConfig(port *uint, logging, encrypt *bool) *torrent.ClientConfig {
	config := torrent.NewDefaultClientConfig()
	config.Logger = log.Default
	config.Logger.Handlers = []log.Handler{log.DiscardHandler}
	config.Seed = true
	config.ListenPort = int(*port)
	config.HeaderObfuscationPolicy.RequirePreferred = *encrypt

	if *logging {
		config.Debug = true
		if path, err := xdg.StateFile("mabel/mini.log"); err == nil {
			if file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err == nil {
				handler := log.DefaultHandler
				handler.W = file
				config.Logger.Handlers = []log.Handler{handler}
			}
		}
	}

	metadataDirectory := os.TempDir()
	if metadataStorage, err := storage.NewDefaultPieceCompletionForDir(metadataDirectory); err != nil {
		config.DefaultStorage = storage.NewMMap("")
	} else {
		config.DefaultStorage = storage.NewMMapWithCompletion("", metadataStorage)
	}

	return config
}

// initialModel creates the model for the mini client. If the torrent
// cannot be generated, the client aborts.
func initialModel(t, dir *string, port *uint, logging, encrypt *bool, theme *styles.ColorTheme) (model, error) {
	client, err := torrent.NewClient(genMabelConfig(port, logging, encrypt))
	if err != nil {
		log.Fatal(err)
	}
	m := model{
		torrent: t,
		saveDir: dir,
		theme:   theme,
		client:  client,
	}
	return m, nil
}

// tick refreshes the UI every half a second in order to update
// download progress.
// Note: could be improved to update by interval of download progress,
// rather than a time interval.
func tick() tea.Cmd {
	return tea.Tick(time.Duration(interval), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init starts the UI and adds startup torrents. If the torrrent cannot
// be generated, the client aborts.
func (m model) Init() tea.Cmd {
	cmd, err := trrnt.AddTorrent(m.torrent, m.saveDir, m.client, nil, m.theme)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

// Update responds to torrent progress, window size changes, and user
// keyboard messages and updates the UI model accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.client.Close()
			return m, tea.Quit
		}
	case trrnt.TorrentDownloadStarted, tickMsg:
		return m, tick()
	}
	return m, nil
}

// View prints the UI model with the download stats and progress bar.
func (m model) View() string {
	var download, upload, peers, meta, bar string

	t := m.client.Torrents()[0]
	info := t.Info()

	if info == nil {
		meta = "getting torrent info..."
	} else {
		download = stats.Download(t, true)
		upload = stats.Upload(t)
		peers = stats.Peers(t)
		meta = fmt.Sprintf(
			"%s | %s | %s",
			download,
			upload,
			peers,
		)

		bar = stats.ProgressBar(t, &m.width, m.theme)
	}

	spacer := m.width - gloss.Width(meta)
	name := utils.TruncateForMinimumPadding(t.Name(), &spacer, 5)

	return gloss.NewStyle().Width(m.width).Render(fmt.Sprintf("%s\n%s\n", name+strings.Repeat(" ", spacer)+meta, bar))
}

// Execute creates the initial model and a Bubble Tea program, and
// aborts the client if that fails.
func Execute(t, dir *string, port *uint, logging, encrypt *bool, theme *styles.ColorTheme) {
	model, err := initialModel(t, dir, port, logging, encrypt, theme)
	if err != nil {
		log.Fatal(err)
	}

	if err := tea.NewProgram(model).Start(); err != nil {
		log.Fatal(err)
	}
}
