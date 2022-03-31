package stats

import (
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/progress"
)

func ProgressBar(t *torrent.Torrent, width *int) string {
	var (
		done    = t.BytesCompleted()
		total   = t.Length()
		percent = float64(done) / float64(total)
		prog    = progress.New(progress.WithDefaultGradient(), progress.WithoutPercentage())
	)

	if width != nil {
		prog.Width = *width
	}

	return prog.ViewAs(percent)
}
