package stats

import (
	"github.com/smmr-software/mabel/internal/styles"

	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/progress"
	gloss "github.com/charmbracelet/lipgloss"
)

func ProgressBar(t *torrent.Torrent, width *int, theme *styles.ColorTheme) string {
	var (
		done      = t.BytesCompleted()
		total     = t.Length()
		percent   = float64(done) / float64(total)
		gradient1 = theme.Gradient1.Dark
		gradient2 = theme.Gradient2.Dark
	)

	if !gloss.HasDarkBackground() {
		gradient1 = theme.Gradient1.Light
		gradient2 = theme.Gradient2.Light
	}
	prog := progress.New(progress.WithGradient(gradient1, gradient2), progress.WithoutPercentage())

	if width != nil {
		prog.Width = *width
	}

	return prog.ViewAs(percent)
}
