package full

import (
	"os"
	"strings"

	"github.com/smmr-software/mabel/internal/list"
	"github.com/smmr-software/mabel/internal/styles"

	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"

	"github.com/charmbracelet/bubbles/help"
	clist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"

	"github.com/adrg/xdg"
)

type model struct {
	width, height   int
	startupTorrents *[]string
	dir             *string
	logging         *bool
	theme           *styles.ColorTheme

	client       *torrent.Client
	clientConfig *torrent.ClientConfig

	list *clist.Model
	help *help.Model

	err error
}

// genMabelConfig configures the torrent client (seeding, listening
// port, log directory, etc.)
func genMabelConfig(port *uint, logging *bool) *torrent.ClientConfig {
	config := torrent.NewDefaultClientConfig()
	config.Logger = log.Default
	config.Logger.Handlers = []log.Handler{log.DiscardHandler}
	config.Seed = true
	config.ListenPort = int(*port)

	if *logging {
		config.Debug = true
		if path, err := xdg.StateFile("mabel/client.log"); err == nil {
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

// genList returns a pointer to a new Bubbles list model for torrents
// on the main page.
func genList() *clist.Model {
	list := clist.New(make([]clist.Item, 0), list.ItemDelegate{}, 0, 0)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)
	return &list
}

// initialModel creates the model for the full client.
func initialModel(torrents *[]string, dir *string, port *uint, logging *bool, theme *styles.ColorTheme) (tea.Model, error) {
	config := genMabelConfig(port, logging)
	client, err := torrent.NewClient(config)
	hlp := help.New()

	m := model{
		startupTorrents: torrents,
		dir:             dir,
		logging:         logging,
		theme:           theme,

		client:       client,
		clientConfig: config,

		list: genList(),
		help: &hlp,
	}

	// Check for port startup failure or other error
	if err != nil {
		msg := err.Error()
		switch {
		case strings.HasPrefix(msg, "subsequent listen"), strings.HasPrefix(msg, "first listen"):
			return initialPortStartupFailure(&m), nil
		default:
			return model{}, err
		}
	}

	return m, nil
}
