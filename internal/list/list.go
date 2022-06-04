// Package list defines the styles and renders the Bubbles list for
// torrents on the main screen.
package list

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/smmr-software/mabel/internal/stats"
	"github.com/smmr-software/mabel/internal/styles"
	"github.com/smmr-software/mabel/internal/utils"

	"github.com/acarl005/stripansi"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

// Define the torrent items that go in the torrent list
type Item struct {
	Self    *torrent.Torrent
	Theme   *styles.ColorTheme
	Added   time.Time
	Created time.Time
	Comment string
	Program string
}

// FilterValue currently has no function, as torrent filtering is not
// currently supported in the main screen.
func (i Item) FilterValue() string { return "" }

// Title returns the name of the torrent.
func (i Item) Title() string { return stripansi.Strip(i.Self.Name()) }

// Description returns the torrent info, including the download,
// upload, and peers stats.
func (i Item) Description() string {
	t := i.Self
	info := t.Info()

	if info == nil {
		return "getting torrent info..."
	}

	var download string
	if t.BytesMissing() != 0 {
		download = stats.Download(t, false)
	} else {
		download = "done!"
	}

	return fmt.Sprintf(
		"%s | %s | %s", download,
		stats.Upload(t),
		stats.Peers(t),
	)
}

type ItemDelegate struct{}

// Height returns the height setting for a list entry.
func (d ItemDelegate) Height() int { return 2 }

// Spacing returns the spacing setting for a list entry.
func (d ItemDelegate) Spacing() int { return 0 }

// Update updates the list model.
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

// Render renders a list entry for the torrent list with the color
// theme styling.
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		i = listItem.(Item)

		leftPadding = 2
		topMargin   = 1

		entry = gloss.NewStyle().Width(m.Width()).PaddingLeft(leftPadding).MarginTop(topMargin)
		title = gloss.NewStyle().Foreground(i.Theme.Primary)
		desc  = gloss.NewStyle().Foreground(i.Theme.Light)

		meta   = i.Description()
		spacer = m.Width() - gloss.Width(meta) - leftPadding
		name   = utils.TruncateForMinimumSpacing(i.Title(), &spacer, 5)
	)

	if index == m.Index() {
		entry = entry.
			Border(gloss.NormalBorder(), false, false, false, true).
			BorderForeground(i.Theme.Dark).
			PaddingLeft(1)
		name = title.Render(name)
		meta = desc.Render(meta)
	}

	fmt.Fprintf(w, entry.Render(name+strings.Repeat(" ", spacer)+meta))
}
