package full

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/smmr-software/mabel/internal/stats"
	"github.com/smmr-software/mabel/internal/utils"

	"github.com/acarl005/stripansi"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type Item struct {
	Self    *torrent.Torrent
	Added   time.Time
	Created time.Time
	Comment string
	Program string
}

func (i Item) FilterValue() string { return "" }
func (i Item) Title() string       { return stripansi.Strip(i.Self.Name()) }
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

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 2 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		i = listItem.(Item)

		leftPadding = 2
		topMargin   = 1

		entry = gloss.NewStyle().Width(m.Width()).PaddingLeft(leftPadding).MarginTop(topMargin)
		title = gloss.NewStyle().Foreground(primaryBlue)
		desc  = gloss.NewStyle().Foreground(lightBlue)

		meta   = i.Description()
		spacer = m.Width() - gloss.Width(meta) - leftPadding
		name   = utils.TruncateForMinimumSpacing(i.Title(), &spacer, 5)
	)

	if index == m.Index() {
		entry = entry.
			Border(gloss.NormalBorder(), false, false, false, true).
			BorderForeground(darkBlue).
			PaddingLeft(1)
		name = title.Render(name)
		meta = desc.Render(meta)
	}

	fmt.Fprintf(w, entry.Render(name+strings.Repeat(" ", spacer)+meta))
}
