package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

type item struct {
	self    *torrent.Torrent
	added   time.Time
	created time.Time
	comment string
	program string
}

func (i item) FilterValue() string { return "" }
func (i item) Title() string       { return i.self.Name() }
func (i item) Description() string {
	t := i.self
	info := t.Info()
	stats := t.Stats()

	if info == nil {
		return "getting torrent info..."
	}

	var download string
	if t.BytesMissing() != 0 {
		download = fmt.Sprintf(
			"%s/%s ↓",
			humanize.Bytes(uint64(t.BytesCompleted())),
			humanize.Bytes(uint64(t.Length())),
		)
	} else {
		download = "done!"
	}

	return fmt.Sprintf(
		"%s | %s ↑ | %d/%d peers",
		download,
		humanize.Bytes(uint64(stats.BytesWritten.Int64())),
		stats.ActivePeers,
		stats.TotalPeers,
	)
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 2 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i := listItem.(item)
	entry := gloss.NewStyle().
		Width(m.Width()).
		Border(gloss.NormalBorder(), false, false, true)

	selected := ""
	if index == m.Index() {
		selected = "* "
	}

	meta := i.Description()
	spacer := m.Width() - gloss.Width(selected) - gloss.Width(meta)
	name := truncateForMinimumSpacing(i.Title(), &spacer, 5)

	fmt.Fprintf(
		w,
		entry.Render(
			gloss.JoinHorizontal(
				gloss.Center,
				selected,
				string(name),
				strings.Repeat(" ", spacer),
				meta,
			),
		),
	)
}
