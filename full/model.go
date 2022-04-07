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
	"github.com/charmbracelet/bubbles/textinput"
)

type model struct {
	width, height   int
	startupTorrents *[]string
	dir             *string
	theme           *styles.ColorTheme

	client       *torrent.Client
	clientConfig *torrent.ClientConfig

	list *clist.Model
	help *help.Model

	addPrompt             *modelAddPrompt
	portStartupFailure    *portStartupFailure
	viewingTorrentDetails bool

	err error
}

type modelAddPrompt struct {
	enabled bool
	dir     bool
	torrent textinput.Model
	saveDir textinput.Model
}

type portStartupFailure struct {
	enabled bool
	port    textinput.Model
}

func initialAddPrompt(dir *string) *modelAddPrompt {
	torrent := textinput.New()
	torrent.Width = 32
	saveDir := torrent

	saveDir.SetValue(*dir)
	saveDir.Blur()

	s := modelAddPrompt{
		enabled: false,
		dir:     false,
		torrent: torrent,
		saveDir: saveDir,
	}
	return &s
}

func initialPortStartupFailure() *portStartupFailure {
	input := textinput.New()
	input.Width = 32
	input.Focus()

	port := portStartupFailure{port: input}
	return &port
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

func genList() *clist.Model {
	list := clist.New(make([]clist.Item, 0), list.ItemDelegate{}, 0, 0)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)
	return &list
}

func initialModel(torrents *[]string, dir *string, port *uint, theme *styles.ColorTheme) (model, error) {
	config := genMabelConfig(port)
	client, err := torrent.NewClient(config)
	hlp := help.New()

	m := model{
		startupTorrents: torrents,
		dir:             dir,
		theme:           theme,

		client:       client,
		clientConfig: config,

		list: genList(),
		help: &hlp,

		addPrompt:          initialAddPrompt(dir),
		portStartupFailure: initialPortStartupFailure(),
	}

	if err != nil {
		msg := err.Error()
		switch {
		case strings.HasPrefix(msg, "subsequent listen"), strings.HasPrefix(msg, "first listen"):
			m.portStartupFailure.enabled = true
		default:
			return model{}, err
		}
	}
	return m, nil
}
