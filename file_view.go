package main

import (
	"fmt"
	"strings"

	"github.com/anacrolix/torrent"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

func fileView(files *[]*torrent.File, w *int) string {
	width := *w / 3
	box := gloss.NewStyle().
		BorderForeground(gloss.Color("86")).
		Border(border).
		Width(width)

	var list strings.Builder
	for i, f := range *files {
		done := uint64(f.BytesCompleted())
		total := uint64(f.Length())
		percentage := int(float64(done) / float64(total) * 100)

		download := fmt.Sprintf(
			"%s/%s (%d%%)",
			humanize.Bytes(done),
			humanize.Bytes(total),
			percentage,
		)

		padding := width - gloss.Width(download) - 2

		name := []rune(f.DisplayPath())
		initial := len(name)
		for padding-gloss.Width(string(name)) < 5 {
			if index := len(name) - 1; index > 0 {
				name = name[:index]
			} else {
				break
			}
		}
		if initial > len(name) {
			name[len(name)-1] = '…'
		}

		padding -= gloss.Width(string(name))
		if padding < 0 {
			padding = 0
		}

		newline := "\n"
		if i == len(*files)-1 {
			newline = ""
		}

		list.WriteString(
			fmt.Sprintf(
				" %s%s%s %s",
				string(name),
				strings.Repeat(" ", padding),
				download,
				newline,
			),
		)
	}

	return box.Render(list.String())
}
