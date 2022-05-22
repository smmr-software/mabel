package full

import (
	"strconv"
	"time"

	"github.com/smmr-software/mabel/internal/list"
	"github.com/smmr-software/mabel/internal/styles"
	trrnt "github.com/smmr-software/mabel/internal/torrent"

	torrent "github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const interval = 500 * time.Millisecond

type tickMsg time.Time
type mabelError error

// reportError takes an error and returns it to Bubble Tea to be
// displayed to the user.
func reportError(err error) tea.Cmd {
	return func() tea.Msg {
		return mabelError(err)
	}
}

// Update responds to port startup failure, window size changes, user
// keyboard messages, if torrent detail view is requested, or if a
// torrent is being added and updates the UI model accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // Change to fit new window size
		width := msg.Width - styles.BorderWindow.GetHorizontalBorderSize()
		height := msg.Height - styles.BorderWindow.GetHorizontalBorderSize()

		m.width = width
		m.help.Width = width
		m.list.SetWidth(int(float64(width) * 0.9))
		m.height = height
		m.list.SetHeight(int(float64(height) * 0.9))
		return m, nil
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC: // Quit client
			if m.client != nil {
				m.client.Close()
			}
			return m, tea.Quit
		case m.err != nil: // Dismiss error display
			m.err = nil
			return m, nil
		case m.portStartupFailure.enabled: // Deal with port startup failure
			return portStartupFailureKeyPress(&m, &msg)
		case m.addPrompt.enabled: // Switch to add prompt screen
			return addPromptKeyPress(&m, &msg)
		case m.viewingTorrentDetails: // Close torrent details screen if open
			m.viewingTorrentDetails = false
			return m, nil
		default: // Deal with default user keyboard messages
			return defaultKeyPress(&m, &msg)
		}
	case tickMsg:
		return m, tick()
	case mabelError: // Deal with error
		m.err = msg
		return m, nil
	default:
		return m, nil
	}
}

// portStartupFailureKeyPress allows the user to choose a port if the
// attempted port binding fails.
func portStartupFailureKeyPress(m *model, msg *tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "backspace":
		var cmd tea.Cmd
		m.portStartupFailure.port, cmd = m.portStartupFailure.port.Update(*msg)
		return m, cmd
	default:
		prt, err := strconv.Atoi(m.portStartupFailure.port.Value())
		if err != nil {
			return m, reportError(err)
		}
		port := uint(prt)

		config := genMabelConfig(&port, m.logging)
		client, err := torrent.NewClient(config)
		if err != nil {
			return m, reportError(err)
		}

		m.client = client
		m.clientConfig = config
		m.portStartupFailure.enabled = false

		return m, nil
	}
}

// addPromptKeyPress handles user key presses when adding a torrent.
func addPromptKeyPress(m *model, msg *tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(*msg, addPromptKeys.quit):
		m.addPrompt = initialAddPrompt(m.dir)
		return m, nil
	case key.Matches(*msg, addPromptKeys.forward):
		if m.addPrompt.dir {
			input := m.addPrompt.torrent.Value()
			dir := m.addPrompt.saveDir.Value()

			m.addPrompt = initialAddPrompt(m.dir)

			cmd, err := trrnt.AddTorrent(&input, &dir, m.client, m.list, m.theme)
			if err != nil {
				return m, reportError(err)
			}

			return m, cmd
		} else {
			m.addPrompt.torrent.Blur()
			m.addPrompt.saveDir.Focus()
			m.addPrompt.dir = true
			return m, nil
		}
	case key.Matches(*msg, addPromptKeys.back):
		if m.addPrompt.dir {
			m.addPrompt.saveDir.Blur()
			m.addPrompt.torrent.Focus()
			m.addPrompt.dir = false
		}
		return m, nil
	default:
		var cmd tea.Cmd
		if m.addPrompt.dir {
			m.addPrompt.saveDir, cmd = m.addPrompt.saveDir.Update(*msg)
		} else {
			m.addPrompt.torrent, cmd = m.addPrompt.torrent.Update(*msg)
		}
		return m, cmd
	}
}

// defaultKeyPress allows the user to interact with the home page and
// quit the client.
func defaultKeyPress(m *model, msg *tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(*msg, homeKeys.quit):
		m.client.Close()
		return m, tea.Quit
	case key.Matches(*msg, homeKeys.help):
		m.help.ShowAll = !m.help.ShowAll
	case key.Matches(*msg, homeKeys.add):
		m.addPrompt.torrent.Focus()
		m.addPrompt.enabled = true
	case key.Matches(*msg, homeKeys.details):
		if t, ok := m.list.SelectedItem().(list.Item); ok && t.Self.Info() != nil {
			m.viewingTorrentDetails = true
		}
	case key.Matches(*msg, homeKeys.up):
		m.list.CursorUp()
	case key.Matches(*msg, homeKeys.down):
		m.list.CursorDown()
	case key.Matches(*msg, homeKeys.next):
		m.list.Paginator.NextPage()
	case key.Matches(*msg, homeKeys.prev):
		m.list.Paginator.PrevPage()
	case key.Matches(*msg, homeKeys.delete):
		zero := list.Item{}
		if t, ok := m.list.SelectedItem().(list.Item); ok && t != zero {
			t.Self.Drop()
			m.list.RemoveItem(m.list.Index())
		}
		if m.list.Index() == len(m.list.Items()) {
			m.list.CursorUp()
		}
	case key.Matches(*msg, homeKeys.deselect):
		m.list.ResetSelected()
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
