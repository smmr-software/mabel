package mini

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/smmr-software/mabel/internal/shared"

	"github.com/adrg/xdg"

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
	client           *torrent.Client
}

func genMabelConfig() *torrent.ClientConfig {
	config := torrent.NewDefaultClientConfig()
	config.Logger = log.Default
	config.Logger.Handlers = []log.Handler{log.DiscardHandler}
	config.Seed = true

	metadataDirectory := os.TempDir()
	if metadataStorage, err := storage.NewDefaultPieceCompletionForDir(metadataDirectory); err != nil {
		config.DefaultStorage = storage.NewMMap("")
	} else {
		config.DefaultStorage = storage.NewMMapWithCompletion("", metadataStorage)
	}

	return config
}

func initialModel(t, dir *string) (model, error) {
	client, err := torrent.NewClient(genMabelConfig())
	if err != nil {
		log.Fatal(err)
	}
	m := model{
		torrent: t,
		saveDir: dir,
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
	cmd, _, _, err := shared.AddTorrent(m.torrent, m.saveDir, m.client)
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
	case shared.TorrentDownloadStarted, tickMsg:
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
		download = shared.DownloadStats(t, true)
		upload = shared.UploadStats(t)
		peers = shared.PeerStats(t)
		meta = fmt.Sprintf(
			"%s | %s | %s",
			download,
			upload,
			peers,
		)

		bar = shared.ProgressBar(t, &m.width)
	}

	spacer := m.width - gloss.Width(meta)
	name := shared.TruncateForMinimumSpacing(t.Name(), &spacer, 5)

	return fmt.Sprintf("%s\n%s\n", name+strings.Repeat(" ", spacer)+meta, bar)
}

func Execute(t *string) {
	dir := xdg.UserDirs.Download
	model, err := initialModel(t, &dir)
	if err != nil {
		log.Fatal(err)
	}

	if err := tea.NewProgram(model).Start(); err != nil {
		log.Fatal(err)
	}
}
