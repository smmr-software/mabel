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
)

const interval = 500 * time.Millisecond

type tickMsg time.Time
type model struct {
	width            int
	torrent, saveDir *string
	theme            *styles.ColorTheme
	client           *torrent.Client
}

func genMabelConfig(port *uint) *torrent.ClientConfig {
	config := torrent.NewDefaultClientConfig()
	config.Logger = log.Default
	config.Logger.Handlers = []log.Handler{log.DiscardHandler}
	config.Seed = true
	config.ListenPort = int(*port)

	metadataDirectory := os.TempDir()
	if metadataStorage, err := storage.NewDefaultPieceCompletionForDir(metadataDirectory); err != nil {
		config.DefaultStorage = storage.NewMMap("")
	} else {
		config.DefaultStorage = storage.NewMMapWithCompletion("", metadataStorage)
	}

	return config
}

func initialModel(t, dir *string, port *uint, theme *styles.ColorTheme) (model, error) {
	client, err := torrent.NewClient(genMabelConfig(port))
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

func tick() tea.Cmd {
	return tea.Tick(time.Duration(interval), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	cmd, err := trrnt.AddTorrent(m.torrent, m.saveDir, m.client, nil)
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

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
	name := utils.TruncateForMinimumSpacing(t.Name(), &spacer, 5)

	return fmt.Sprintf("%s\n%s\n", name+strings.Repeat(" ", spacer)+meta, bar)
}

func Execute(t, dir *string, port *uint, theme *styles.ColorTheme) {
	model, err := initialModel(t, dir, port, theme)
	if err != nil {
		log.Fatal(err)
	}

	if err := tea.NewProgram(model).Start(); err != nil {
		log.Fatal(err)
	}
}
