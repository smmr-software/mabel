package full

import (
	"time"

	"github.com/smmr-software/mabel/internal/list"
	"github.com/smmr-software/mabel/internal/styles"

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

// Update responds to window size changes, user keyboard messages,
// client errors, and updates the UI model accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // Change to fit new window size
		width := msg.Width - styles.Fullscreen.GetHorizontalBorderSize()
		height := msg.Height - styles.Fullscreen.GetHorizontalBorderSize()

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
		case key.Matches(msg, homeKeys.quit):
			m.client.Close()
			return m, tea.Quit
		case key.Matches(msg, homeKeys.help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, homeKeys.add):
			return initialAddTorrent(
				m.width, m.height,
				m.dir, m.theme,
				m.help, &m,
			), nil
		case key.Matches(msg, homeKeys.details):
			if t, ok := m.list.SelectedItem().(list.Item); ok && t.Self.Info() != nil {
				return torrentDetails{
					width:  m.width,
					height: m.height,
					item:   &t,
					theme:  m.theme,
					main:   &m,
				}, nil
			}
		case key.Matches(msg, homeKeys.up):
			m.list.CursorUp()
		case key.Matches(msg, homeKeys.down):
			m.list.CursorDown()
		case key.Matches(msg, homeKeys.next):
			m.list.Paginator.NextPage()
		case key.Matches(msg, homeKeys.prev):
			m.list.Paginator.PrevPage()
		case key.Matches(msg, homeKeys.delete):
			zero := list.Item{}
			if t, ok := m.list.SelectedItem().(list.Item); ok && t != zero {
				t.Self.Drop()
				m.list.RemoveItem(m.list.Index())
			}
			if m.list.Index() == len(m.list.Items()) {
				m.list.CursorUp()
			}
		case key.Matches(msg, homeKeys.deselect):
			m.list.ResetSelected()
		}
	case tickMsg:
		return m, tick()
	case mabelError: // Deal with error
		return errorScreen{
			width:  m.width,
			height: m.height,
			err:    msg,
			theme:  m.theme,
			main:   &m,
		}, nil
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
